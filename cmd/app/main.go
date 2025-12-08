package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/noredis/subscriptions/internal/common/config"
	"github.com/noredis/subscriptions/pkg/postgres"
	"github.com/rs/zerolog"
)

func main() {
	cfg := config.MustLoad()
	log.Printf("config readed successfully")

	logger := setupLogger(cfg.Logger)
	logger.Info().Msg("logger configured")

	ctx := context.Background()
	db, err := postgres.New(ctx, cfg.DB.DSN(), cfg.DB.Attempts, cfg.DB.Delay, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}
	logger.Info().Msg("app successfully connected to db")
	_ = db

	app := fiber.New()
	app.Use(recover.New())
	logger.Info().Msgf("app starting on port %d", cfg.App.Port)

	if err := app.Listen(fmt.Sprintf(":%d", cfg.App.Port)); err != nil {
		logger.Fatal().Err(err).Msg("failed to start app")
	}
}

func setupLogger(cfg config.Logger) *zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		log.Fatalf("failed to parse log level")
	}

	logger = logger.Level(level)

	return &logger
}
