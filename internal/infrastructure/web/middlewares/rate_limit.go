package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dprio/redis-limiter/internal/infrastructure/ratelimiter"
)

const apiKeyHeader = "API_KEY"

type RateLimitMiddleware struct {
	limiter ratelimiter.LimiterGateway
}

func newRateLimitMiddleware(limiter ratelimiter.LimiterGateway) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiter: limiter,
	}
}

func (rlm *RateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := strings.Split(r.RemoteAddr, ":")[0]
		apiKey := r.Header.Get(apiKeyHeader)

		if err := rlm.limiter.Check(r.Context(), clientIP, apiKey); err != nil {
			status := http.StatusTooManyRequests
			msg := "you have reached the maximum number of requests or actions allowed within a certain time frame"

			if errors.Is(err, ratelimiter.ErrInvalidAPIKey) {
				status = http.StatusUnauthorized
				msg = "invalid api key"
			}

			http.Error(w, msg, status)
			return
		}

		next.ServeHTTP(w, r)
	})
}
