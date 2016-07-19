package env_config

import (
	"github.com/caarlos0/env"
	"log"
)

type Config struct {
	TelegramToken string `env:"TELEGRAM_TOKEN"`
	Port          int    `env:"PORT" envDefault:"8080"`
	IpAddress     string `env:"IP_ADDRESS"`
	Board         string `env:"BOARD"`
}

func ReadConfig() Config {
	config := Config{}
	err := env.Parse(&config)
	if err != nil {
		log.Panic(err)
	}

	return config
}
