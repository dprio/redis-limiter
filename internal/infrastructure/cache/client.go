package cache

import (
	"context"
	"errors"
	"time"
)

var ErrKeyNotFound = errors.New("key not found")

type Client interface {
	Get(ctx context.Context, key string) (any, error)
	Increment(ctx context.Context, key string, ttl time.Duration) (int64, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
}
