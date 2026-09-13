package ratelimiter

import (
	"context"
	"errors"
	"log"

	"github.com/dprio/redis-limiter/internal/domain"
	"github.com/dprio/redis-limiter/internal/infrastructure/cache"
	"github.com/dprio/redis-limiter/internal/infrastructure/config"
	"github.com/dprio/redis-limiter/internal/infrastructure/db"
)

var (
	ErrRateLimitted   = errors.New("rate limit exceeded")
	ErrBlockingFailed = errors.New("failed to block the key")
	ErrInvalidAPIKey  = errors.New("invalid api key")
)

type LimiterGateway interface {
	Check(ctx context.Context, clientIP string, apiKey string) error
}

type rateLimiter struct {
	cacheClient     cache.Client
	tokenRepository db.TokenRepository
	cfg             *config.RateLimiter
}

func NewLimiterGateway(cacheClient cache.Client, tokenRepository db.TokenRepository, cfg *config.RateLimiter) LimiterGateway {
	return &rateLimiter{
		cacheClient:     cacheClient,
		tokenRepository: tokenRepository,
		cfg:             cfg,
	}
}

// Check decide qual limite aplicar: se apiKey for informada, ela tem
// prioridade sobre o IP. Se apiKey estiver vazia, aplica limite por IP.
func (l *rateLimiter) Check(ctx context.Context, clientIP, apiKey string) error {
	key, maxRequests, err := l.resolveLimit(ctx, clientIP, apiKey)
	if err != nil {
		return err
	}

	if l.isBlocked(ctx, key) {
		log.Printf("Key %s is blocked", key)
		return ErrRateLimitted
	}

	rateLimitKey := "rate_limit:" + key

	count, err := l.cacheClient.Increment(ctx, rateLimitKey, l.cfg.WindowDuration)
	if err != nil {
		return err
	}

	log.Printf("Key: %s, Count: %d, Max: %d", key, count, maxRequests)
	if count > maxRequests {
		if blockErr := l.block(ctx, key); blockErr != nil {
			log.Printf("failed to block key %s: %v", key, blockErr)
		}
		return ErrRateLimitted
	}

	return nil
}

func (l *rateLimiter) resolveLimit(ctx context.Context, clientIP string, apiKey string) (key string, maxRequests int64, err error) {
	if apiKey == "" {
		return "ip:" + clientIP, l.cfg.IPMaxRequests, nil
	}

	token, err := l.tokenRepository.Get(ctx, apiKey)
	if err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			return "", 0, ErrInvalidAPIKey
		}
		return "", 0, err
	}

	return "token:" + token.Hash, token.TotalRequests, nil
}

func (l *rateLimiter) isBlocked(ctx context.Context, key string) bool {
	blockedKey := "blocked:" + key
	if _, err := l.cacheClient.Get(ctx, blockedKey); err != nil && errors.Is(err, cache.ErrKeyNotFound) {
		return false
	}
	return true
}

func (l *rateLimiter) block(ctx context.Context, key string) error {
	blockedKey := "blocked:" + key

	if err := l.cacheClient.Set(ctx, blockedKey, "blocked", l.cfg.BlockDuration); err != nil {
		return ErrBlockingFailed
	}
	return nil
}
