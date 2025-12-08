package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Logger interface {
	Printf(format string, v ...any)
}

func New(
	ctx context.Context,
	dsn string,
	attempts int,
	delay time.Duration,
	logger Logger,
) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := pingWithRetry(ctx, pool, attempts, delay, logger); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func pingWithRetry(
	ctx context.Context,
	pool *pgxpool.Pool,
	attempts int,
	delay time.Duration,
	logger Logger,
) error {
	var err error

	for i := 1; i <= attempts; i++ {
		if err = pool.Ping(ctx); err == nil {
			return nil
		}

		if i < attempts {
			logger.Printf("ping to db failed (%d/%d): %v. Retrying in %s...\n", i, attempts, err, delay)

			select {
			case <-ctx.Done():
				return fmt.Errorf("context cancelled: %w", ctx.Err())
			case <-time.After(delay):
			}
		}
	}

	return fmt.Errorf("failed to ping database: %w", err)
}
