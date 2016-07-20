package env_config

import (
	"github.com/caarlos0/env"
	"log"
)

type Config struct {
	TelegramToken       string `env:"TELEGRAM_TOKEN"`
	BoardPollingTimeout int    `env:"BOARD_POLLING_TIMEOUT" envDefault:"5"`
	IpAddress           string `env:"IP_ADDRESS"`
	Port                int    `env:"PORT" envDefault:"8080"`
}

func ReadConfig() Config {
	config := Config{}
	err := env.Parse(&config)
	if err != nil {
		log.Panic(err)
	}

	return config
}
