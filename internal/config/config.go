package config

import (
	"context"
	"fmt"
	"github.com/iamstep4ik/TestTaskOzonBank/internal/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"time"
)

type Config struct {
	Server struct {
		Port string `envconfig:"SERVER_PORT" `
		Host string `envconfig:"SERVER_HOST"`
	} `envconfig:"SERVER"`
	Database struct {
		Host     string `envconfig:"DB_HOST"`
		Port     string `envconfig:"DB_PORT" `
		Name     string `envconfig:"DB_NAME" `
		User     string `envconfig:"DB_USER" `
		Password string `envconfig:"DB_PASSWORD"`
		SSLMode  string `envconfig:"DB_SSLMODE"`
	} `envconfig:"DATABASE"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Error("Failed to load configuration from environment variables", zap.Error(err))
	}
	return &cfg, nil
}

func (c *Config) ConnectDatabase(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name, cfg.Database.SSLMode)
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Error("Failed to parse database connection string", zap.Error(err))
	}
	dbpool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Error("Failed to create database connection pool", zap.Error(err))
	}
	err = dbpool.Ping(ctx)
	if err != nil {
		log.Error("Failed to ping database", zap.Error(err))
	}
	return dbpool, nil
}
