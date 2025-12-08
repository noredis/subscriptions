package config

import (
	"fmt"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App    App
	Logger Logger
	DB     DB
}

type App struct {
	Env  string `envconfig:"APP_ENV" default:"dev"`
	Port int    `envconfig:"APP_PORT" default:"8080"`
}

type Logger struct {
	Level string `envconfig:"LOG_LEVEL" default:"debug"`
}

type DB struct {
	User            string        `envconfig:"DB_USER" required:"true"`
	Password        string        `envconfig:"DB_PASSWORD" required:"true"`
	Host            string        `envconfig:"DB_HOST" default:"localhost"`
	Port            int           `envconfig:"DB_PORT" default:"5432"`
	Name            string        `envconfig:"DB_NAME" required:"true"`
	MaxConns        int           `envconfig:"DB_MAX_CONNS" default:"4"`
	MinConns        int           `envconfig:"DB_MIN_CONNS" default:"1"`
	MaxConnLifetime time.Duration `envconfig:"DB_MAX_CONN_LIFETIME" default:"1h"`
	MaxConnIdleTime time.Duration `envconfig:"DB_MAX_CONN_IDLE_TIME" default:"15m"`
	Attempts        int           `envconfig:"DB_CONN_ATTEMPTS" default:"5"`
	Delay           time.Duration `envconfig:"DB_CONN_DELAY" default:"3s"`
}

func (dbCfg *DB) DSN() string {
	dbPoolParams := fmt.Sprintf(
		"pool_max_conns=%d&pool_min_conns=%d&pool_max_conn_lifetime=%s&pool_max_conn_idle_time=%s",
		dbCfg.MaxConns,
		dbCfg.MinConns,
		dbCfg.MaxConnLifetime,
		dbCfg.MaxConnIdleTime,
	)

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?%s",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Name,
		dbPoolParams,
	)
}

func MustLoad() *Config {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return &cfg
}
