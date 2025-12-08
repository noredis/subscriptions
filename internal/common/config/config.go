package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App    App
	Logger Logger
}

type App struct {
	Env string `envconfig:"APP_ENV" default:"dev"`
}

type Logger struct {
	Level string `envconfig:"LOG_LEVEL" default:"debug"`
}

func MustLoad() *Config {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("unable to load config: %v", err)
	}

	return &cfg
}
