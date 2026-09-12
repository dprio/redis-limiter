package middlewares

import "github.com/dprio/redis-limiter/internal/infrastructure/ratelimiter"

type Middlewares struct {
	RateLimitMiddleware *RateLimitMiddleware
}

func New(limiterGateway ratelimiter.LimiterGateway) *Middlewares {
	return &Middlewares{
		RateLimitMiddleware: newRateLimitMiddleware(limiterGateway),
	}
}
