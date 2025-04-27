package database

import (
	"context"
	"fmt"
	"time"
	"log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nurashi/AIForge/auth-service/config"
	"github.com/redis/go-redis/v9"
)

func InitPostgres(cfg *config.PostgresConfig) (*pgxpool.Pool, error) {

	dbpool, err := pgxpool.New(context.Background(), cfg.DSN()) 
	if err != nil {
		return nil, fmt.Errorf("ERROR with dbpool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = dbpool.Ping(ctx); err != nil {
		dbpool.Close() 
		return nil, fmt.Errorf("ERROR with ctx in postgres %w", err)
	}
	log.Println("connected to PostgreSQL - NICE") 
	return dbpool, nil
}

func InitRedis(cfg *config.RedisConfig) (*redis.Client, error) {

	opts := &redis.Options{
		Addr:     cfg.Addr(), 
		Password: cfg.Password,
		DB:       cfg.DB,
	}
	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("ERROR with ctx in redis %w", err)
	}
	log.Println("connected to Redis - NICE") 
	return rdb, nil
}
