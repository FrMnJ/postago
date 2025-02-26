package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/FrMnJ/postago/src/config"
	"github.com/FrMnJ/postago/src/email"
	"github.com/redis/go-redis/v9"
)

type MailQueue struct {
	RedisClient *redis.Client
}

func NewMailQueue() (*MailQueue, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%d",
			config.AppConfig.RedisConfig.Addr,
			config.AppConfig.RedisConfig.Port,
		),
		Password: config.AppConfig.RedisConfig.Passwd,
		DB:       config.AppConfig.RedisConfig.DB, // use default DB
	})

	log.Println("Ping...")
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	log.Println("Pong!")
	return &MailQueue{
		RedisClient: rdb,
	}, nil
}

func (mq *MailQueue) ReturnToQueue(info string) error {
	ctx := context.Background()
	_, err := mq.RedisClient.RPush(ctx, config.AppConfig.QueueName, info).Result()
	return err
}

func (mq *MailQueue) MainLoop() {
	ctx := context.Background()
	emailService := email.NewGmailEmailServiceAdapter()
	for {
		info, err := mq.RedisClient.LPop(ctx, config.AppConfig.QueueName).Result()
		log.Println("Info:", info)
		if err == redis.Nil {
			log.Println("Queue is empty, waiting for 1 minute...")
			time.Sleep(1 * time.Minute)
			continue
		} else if err != nil {
			fmt.Println(err)
			time.Sleep(15 * time.Second)
			continue
		}

		var infoMap map[string]interface{}
		err = json.Unmarshal([]byte(info), &infoMap)
		if err != nil {
			log.Println(err)
			time.Sleep(15 * time.Second)
			continue
		}

		log.Println("Sending email to:", infoMap["toEmail"])
		log.Println("Subject:", infoMap["subject"])
		log.Println("Template:", infoMap["templateName"])
		log.Println("Data:", infoMap["data"])

		if err := emailService.SendEmail(
			infoMap["toEmail"].(string),
			infoMap["subject"].(string),
			infoMap["templateName"].(string),
			infoMap["data"].(map[string]interface{}),
		); err != nil {
			log.Println(err)
			if err := mq.ReturnToQueue(info); err != nil {
				log.Println(err)
				continue
			}
			time.Sleep(15 * time.Second)
			continue
		}
	}
}
