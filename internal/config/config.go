package config

import (
	"context"
	"fmt"
	"time"

	"github.com/iamstep4ik/TestTaskOzonBank/internal/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
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
	Redis struct {
		Addr     string `envconfig:"REDIS_ADDR"`
		Password string `envconfig:"REDIS_PASSWORD"`
		DB       int    `envconfig:"REDIS_DB"`
	} `envconfig:"REDIS"`
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
func (c *Config) ConnectRedis(ctx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Error("Failed to connect to Redis", zap.Error(err))
		return nil, err
	}

	return rdb, nil
}
