package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App App
}

type App struct {
	Env string `envconfig:"APP_ENV" default:"dev"`
}

func MustLoad() *Config {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	return &cfg
}
