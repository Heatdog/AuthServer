package redisStorage

import (
	"context"
	"fmt"
	"time"

	"github.com/Heaterdog/AuthServer/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, cfg *config.RedisSettings) (*redis.Client, error) {
	time.Sleep(time.Duration(cfg.TimePrepare) * time.Second)
	host := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	storage := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: cfg.Password,
		DB:       0,
	})

	if _, err := storage.Ping(ctx).Result(); err != nil {
		return nil, err
	}
	return storage, nil
}
