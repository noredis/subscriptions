package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/noredis/subscriptions/docs"
	"github.com/noredis/subscriptions/internal/application/appservice"
	"github.com/noredis/subscriptions/internal/common/config"
	"github.com/noredis/subscriptions/internal/domain/service"
	"github.com/noredis/subscriptions/internal/infrastructure/repository"
	"github.com/noredis/subscriptions/internal/presentation/http/handlers"
	"github.com/noredis/subscriptions/internal/presentation/http/middlewares"
	"github.com/noredis/subscriptions/pkg/postgres"
	"github.com/noredis/subscriptions/pkg/rules"
	"github.com/noredis/subscriptions/pkg/validatorext"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func main() {
	app := NewApp()

	if err := app.Init(); err != nil {
		log.Fatal().Err(err).Msg("failed to init app")
	}

	if err := app.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to start app")
	}
}

type App struct {
	cfg      *config.Config
	db       *pgxpool.Pool
	logger   *zerolog.Logger
	fiberApp *fiber.App
}

func NewApp() *App {
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

	return &App{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}
}

func (app *App) Init() error {
	app.fiberApp = fiber.New()

	app.fiberApp.Use(recover.New())
	app.fiberApp.Use(middlewares.Logging(app.logger))

	heartbeatHandler := handlers.NewHeartbeatHandler()
	heartbeatHandler.Register(app.fiberApp)

	validate := validator.New()
	if err := validate.RegisterValidation("date_format", rules.DateFormat); err != nil {
		return err
	}

	validate.RegisterTagNameFunc(validatorext.FieldTag)

	subscriptionRepo := repository.NewSubscriptionRepository(app.db)
	subscriptionService := appservice.NewSubscriptionService(validate, subscriptionRepo)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService, app.logger)
	subscriptionHandler.Register(app.fiberApp)
	log.Printf("VALIDATOR BEFORE: %#v\n", validate)

	calculator := service.NewCostCalculator()
	costService := appservice.NewCostService(validate, subscriptionRepo, calculator)
	costHandler := handlers.NewCostHandler(app.logger, costService)
	costHandler.Register(app.fiberApp)

	app.fiberApp.Get("/swagger/*", fiberSwagger.WrapHandler)

	return nil
}

func (app *App) Start() error {
	app.logger.Info().Msgf("app starting on port %d", app.cfg.App.Port)

	go func() {
		if err := app.fiberApp.Listen(fmt.Sprintf(":%d", app.cfg.App.Port)); err != nil {
			app.logger.Fatal().Err(err).Msg("failed to start app")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quit

	app.logger.Info().Msg("received shutdown signal")

	app.db.Close()
	app.logger.Info().Msg("database connection closed")

	return app.Shutdown()
}

func (app *App) Shutdown() error {
	app.logger.Info().Msg("shutting down...")

	if err := app.fiberApp.Shutdown(); err != nil {
		app.logger.Error().Err(err).Msg("fiber shutdown failed")
	}

	app.logger.Info().Msg("shutdown completed")
	return nil
}

func setupLogger(cfg config.Logger) *zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		log.Fatal().Msg("failed to parse log level")
	}

	logger = logger.Level(level)

	return &logger
}
