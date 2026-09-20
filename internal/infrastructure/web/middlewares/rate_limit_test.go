package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dprio/redis-limiter/internal/infrastructure/ratelimiter"
	mockratelimiter "github.com/dprio/redis-limiter/mocks/ratelimiter"
)

func TestRateLimitMiddleware_Handle(t *testing.T) {
	tests := []struct {
		name       string
		remoteAddr string
		apiKeyHdr  string
		mockErr    error
		wantStatus int
		nextCalled bool
	}{
		{
			name:       "permite passar quando nao ha erro (por IP)",
			remoteAddr: "10.0.0.1:54321",
			apiKeyHdr:  "",
			mockErr:    nil,
			wantStatus: http.StatusOK,
			nextCalled: true,
		},
		{
			name:       "permite passar quando nao ha erro (por token)",
			remoteAddr: "10.0.0.1:54321",
			apiKeyHdr:  "minha-api-key",
			mockErr:    nil,
			wantStatus: http.StatusOK,
			nextCalled: true,
		},
		{
			name:       "retorna 429 quando rate limit excedido",
			remoteAddr: "10.0.0.1:54321",
			apiKeyHdr:  "",
			mockErr:    ratelimiter.ErrRateLimitted,
			wantStatus: http.StatusTooManyRequests,
			nextCalled: false,
		},
		{
			name:       "retorna 401 quando api key invalida",
			remoteAddr: "10.0.0.1:54321",
			apiKeyHdr:  "chave-invalida",
			mockErr:    ratelimiter.ErrInvalidAPIKey,
			wantStatus: http.StatusUnauthorized,
			nextCalled: false,
		},
		{
			name:       "retorna 429 para erro generico do limiter",
			remoteAddr: "10.0.0.1:54321",
			apiKeyHdr:  "",
			mockErr:    assert.AnError,
			wantStatus: http.StatusTooManyRequests,
			nextCalled: false,
		},
		{
			name:       "extrai IP corretamente removendo a porta",
			remoteAddr: "203.0.113.5:8080",
			apiKeyHdr:  "",
			mockErr:    nil,
			wantStatus: http.StatusOK,
			nextCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedIP := extractIP(tt.remoteAddr)

			limiter := mockratelimiter.NewMockLimiterGateway(t)
			limiter.EXPECT().
				Check(mock.Anything, expectedIP, tt.apiKeyHdr).
				Return(tt.mockErr)

			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			mw := newRateLimitMiddleware(limiter)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = tt.remoteAddr
			if tt.apiKeyHdr != "" {
				req.Header.Set(apiKeyHeader, tt.apiKeyHdr)
			}
			rec := httptest.NewRecorder()

			mw.Handle(next).ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, tt.nextCalled, nextCalled)
		})
	}
}

// extractIP replica a extração feita no middleware (strings.Split(r.RemoteAddr, ":")[0])
// só pra deixar explícito no teste qual IP é esperado na chamada ao mock.
func extractIP(remoteAddr string) string {
	for i, c := range remoteAddr {
		if c == ':' {
			return remoteAddr[:i]
		}
	}
	return remoteAddr
}

var _ = context.Background // evita import não usado caso context não seja referenciado diretamente
