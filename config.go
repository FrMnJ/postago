package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var AppConfig = InitConfig()

type Config struct {
	RedisAddr   string
	RedisPasswd string
	RedisDB     int
	RedisPort   int64

	SmtpDomain string
	SmtpPort   int
	MailPasswd string
	Mail       string
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
	return &Config{
		RedisAddr:   viper.GetString("redis.address"),
		RedisPort:   viper.GetInt64("redis.port"),
		RedisPasswd: viper.GetString("redis.passwd"),
		RedisDB:     viper.GetInt("redis.db"),

		Mail:       viper.GetString("mail.account.address"),
		MailPasswd: viper.GetString("mail.account.passwd"),
		SmtpDomain: viper.GetString("mail.server.domain"),
		SmtpPort:   viper.GetInt("mail.server.port"),
	}
}
