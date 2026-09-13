package cache

import (
	"context"
	"errors"
	"time"

	"github.com/dprio/redis-limiter/internal/infrastructure/config"
	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	redis *redis.Client
}

func NewRedisClient(cfg *config.Redis) Client {
	rds := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := rds.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return &redisClient{redis: rds}
}

func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.redis.Get(ctx, key).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return "", ErrKeyNotFound
	case err != nil:
		return "", err
	default:
		return val, nil
	}
}

func (r *redisClient) Increment(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	pipe := r.redis.Pipeline()
	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, ttl)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return incrCmd.Val(), nil
}

func (r *redisClient) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	result := r.redis.Set(ctx, key, value, ttl)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}
