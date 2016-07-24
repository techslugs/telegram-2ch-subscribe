package env_config

import (
	"github.com/caarlos0/env"
	"log"
)

type Config struct {
	TelegramToken            string `env:"TELEGRAM_TOKEN"`
	BoardPollingTimeout      int    `env:"BOARD_POLLING_TIMEOUT" envDefault:"5"`
	BoardsListPollingTimeout int    `env:"BOARD_POLLING_TIMEOUT" envDefault:"300"`
	MongoURL                 string `env:"MONGO_URL" envDefault:"127.0.0.1"`
	MongoDatabase            string `env:"MONGO_DATABASE" envDefault:"telegram-2ch-subscribe"`
	IpAddress                string `env:"IP_ADDRESS"`
	Port                     int    `env:"PORT"`
}

func ReadConfig() Config {
	config := Config{}
	err := env.Parse(&config)
	if err != nil {
		log.Panic(err)
	}

	if config.IpAddress != "" && config.Port != 0 {
		return config
	}

	openshiftConfig := ReadOpenshiftConfig()
	if config.IpAddress == "" {
		config.IpAddress = openshiftConfig.IpAddress
	}
	if config.Port == 0 {
		config.Port = openshiftConfig.Port
	}

	return config
}
