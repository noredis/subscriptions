package main

import (
	"log"
	"os"

	"github.com/noredis/subscriptions/internal/common/config"
	"github.com/rs/zerolog"
)

func main() {
	cfg := config.MustLoad()
	_ = cfg
	log.Printf("config readed successfully")

	logger := setupLogger(cfg.Logger)
	logger.Info().Msg("logger configured")

	// db

	// http
}

func setupLogger(cfg config.Logger) *zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		log.Fatalf("unable to parse log level")
	}

	logger = logger.Level(level)

	return &logger
}
