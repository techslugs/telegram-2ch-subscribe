package env_config

import (
	"github.com/caarlos0/env"
	"log"
)

type OpenshiftConfig struct {
	IpAddress string `env:"OPENSHIFT_GO_IP" envDefault:"127.0.0.1"`
	Port      int    `env:"OPENSHIFT_GO_PORT" envDefault:"8080"`
}

func ReadOpenshiftConfig() OpenshiftConfig {
	config := OpenshiftConfig{}
	err := env.Parse(&config)
	if err != nil {
		log.Panic(err)
	}

	return config
}
