package ratelimiter

import (
	"context"
	"errors"
	"log"

	"github.com/dprio/redis-limiter/internal/infrastructure/cache"
	"github.com/dprio/redis-limiter/internal/infrastructure/config"
)

var (
	ErrRateLimitted   = errors.New("rate limit exceeded")
	ErrBlockingFailed = errors.New("failed to block the key")
)

type LimiterGateway interface {
	IsBlocked(ctx context.Context, key string) bool
	Check(ctx context.Context, key string) error
	Block(ctx context.Context, key string) error
}

type rateLimiter struct {
	cacheClient cache.Client
	cfg         *config.RateLimiter
}

func NewLimiterGateway(cacheClient cache.Client, cfg *config.RateLimiter) LimiterGateway {
	return &rateLimiter{
		cacheClient: cacheClient,
		cfg:         cfg,
	}
}

func (l *rateLimiter) IsBlocked(ctx context.Context, key string) bool {
	blockedKey := "blocked:" + key
	if _, err := l.cacheClient.Get(ctx, blockedKey); err != nil && errors.Is(err, cache.ErrKeyNotFound) {
		return false
	}
	return true
}

func (l *rateLimiter) Check(ctx context.Context, key string) error {

	if l.IsBlocked(ctx, key) {
		log.Printf("Key %s is blocked", key)
		return ErrRateLimitted
	}

	rateLimitKey := "rate_limit:" + key

	count, err := l.cacheClient.Increment(ctx, rateLimitKey, l.cfg.WindowDuration)
	if err != nil {
		return err
	}

	log.Printf("Key: %s, Count: %d", key, count)
	if count > l.cfg.IPMaxRequests {
		l.Block(ctx, key)
		return ErrRateLimitted
	}

	return nil
}

func (l *rateLimiter) Block(ctx context.Context, key string) error {
	blockedKey := "blocked:" + key

	r := l.cacheClient.Set(ctx, blockedKey, "blocked", l.cfg.BlockDuration)
	if r != nil {
		return ErrBlockingFailed
	}
	return nil
}
