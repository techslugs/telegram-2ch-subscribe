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
	MongoDatabase            string `env:"MONGO_DATABASE" envDefault:"telegram-2ch-news-bot"`
	IpAddress                string `env:"IP_ADDRESS" envDefault:"127.0.0.1"`
	Port                     int    `env:"PORT" envDefault:"8080"`
}

func ReadConfig() Config {
	config := Config{}
	err := env.Parse(&config)
	if err != nil {
		log.Panic(err)
	}

	return config
}
