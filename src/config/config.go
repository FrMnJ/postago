package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var AppConfig = InitConfig()

type RedisConfig struct {
	Addr   string
	Passwd string
	DB     int
	Port   int64
}

type MailConfig struct {
	Account  string
	Passwd   string
	SmtpPort int
	SmtpHost string
}

type Config struct {
	RedisConfig
	MailConfig

	QueueName string
}

func GetBaseProjectPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	for !strings.HasSuffix(dir, "postago") {
		dir = filepath.Dir(dir)
	}
	return dir
}

func configureViper() {
	viper.SetConfigName("conf")
	viper.AddConfigPath(GetBaseProjectPath())
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file")
	}
}

func InitConfig() *Config {
	configureViper()
	redisConfig := RedisConfig{
		Addr:   viper.GetString("redis.address"),
		Port:   viper.GetInt64("redis.port"),
		Passwd: viper.GetString("redis.passwd"),
		DB:     viper.GetInt("redis.db"),
	}

	mailConfig := MailConfig{
		Account:  viper.GetString("mail.account.address"),
		Passwd:   viper.GetString("mail.account.passwd"),
		SmtpPort: viper.GetInt("mail.server.port"),
		SmtpHost: viper.GetString("mail.server.domain"),
	}

	return &Config{
		RedisConfig: redisConfig,
		MailConfig:  mailConfig,
		QueueName:   viper.GetString("queue.name"),
	}
}
