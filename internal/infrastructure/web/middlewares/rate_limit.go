package middlewares

import (
	"net/http"
	"strings"

	"github.com/dprio/redis-limiter/internal/infrastructure/ratelimiter"
)

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

		if err := rlm.limiter.Check(r.Context(), clientIP); err != nil {
			http.Error(w,
				"you have reached the maximum number of requests or actions allowed within a certain time frame",
				http.StatusTooManyRequests,
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
