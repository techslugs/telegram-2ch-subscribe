package env_config

import (
	"github.com/caarlos0/env"
	"log"
)

type Config struct {
	TelegramToken            string `env:"TELEGRAM_TOKEN"`
	BoardPollingTimeout      int    `env:"BOARD_POLLING_TIMEOUT" envDefault:"5"`
	BoardsListPollingTimeout int    `env:"BOARD_POLLING_TIMEOUT" envDefault:"300"`
	IpAddress                string `env:"IP_ADDRESS"`
	Port                     int    `env:"PORT"`
	MongoURL                 string `env:"MONGO_URL"`
	MongoDatabase            string `env:"MONGO_DATABASE"`
}

func ReadConfig() Config {
	config := Config{}
	err := env.Parse(&config)
	if err != nil {
		log.Panic(err)
	}

	openshiftConfig := ReadOpenshiftConfig()
	config.IpAddress = getStringValue(config.IpAddress, openshiftConfig.IpAddress)
	config.Port = getIntValue(config.Port, openshiftConfig.Port)
	config.MongoURL = getStringValue(config.MongoURL, openshiftConfig.MongoURL)
	config.MongoDatabase = getStringValue(config.MongoDatabase, openshiftConfig.MongoDatabase)

	return config
}

func getStringValue(original string, new string) string {
	if original == "" {
		return new
	}
	return original
}

func getIntValue(original int, new int) int {
	if original == 0 {
		return new
	}
	return original
}
