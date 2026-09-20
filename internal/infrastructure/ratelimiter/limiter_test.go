package ratelimiter

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dprio/redis-limiter/internal/domain"
	"github.com/dprio/redis-limiter/internal/infrastructure/cache"
	"github.com/dprio/redis-limiter/internal/infrastructure/config"
	mockcache "github.com/dprio/redis-limiter/mocks/cache"
	mockdb "github.com/dprio/redis-limiter/mocks/db"
)

const (
	testClientIP       = "192.168.0.1"
	testAPIKey         = "minha-api-key"
	testTokenHash      = "abc123hash"
	testIPMaxRequests  = int64(10)
	testWindowDuration = time.Minute
	testBlockDuration  = 5 * time.Minute
)

var errInfra = errors.New("connection refused")

func newTestConfig() *config.RateLimiter {
	return &config.RateLimiter{
		WindowDuration: testWindowDuration,
		IPMaxRequests:  testIPMaxRequests,
		BlockDuration:  testBlockDuration,
	}
}

func TestRateLimiter_Check(t *testing.T) {
	ipKey := "ip:" + testClientIP
	tokenKey := "token:" + testTokenHash

	tests := []struct {
		name       string
		clientIP   string
		apiKey     string
		setupMocks func(cacheClient *mockcache.MockClient, tokenRepo *mockdb.MockTokenRepository)
		wantErr    error
	}{
		// ---------- Caminho por IP ----------
		{
			name:     "por IP: dentro do limite",
			clientIP: testClientIP,
			apiKey:   "",
			setupMocks: func(c *mockcache.MockClient, _ *mockdb.MockTokenRepository) {
				c.EXPECT().Get(mock.Anything, "blocked:"+ipKey).Return("", cache.ErrKeyNotFound)
				c.EXPECT().Increment(mock.Anything, "rate_limit:"+ipKey, testWindowDuration).Return(int64(1), nil)
			},
			wantErr: nil,
		},
		{
			name:     "por IP: exatamente no limite ainda passa",
			clientIP: testClientIP,
			apiKey:   "",
			setupMocks: func(c *mockcache.MockClient, _ *mockdb.MockTokenRepository) {
				c.EXPECT().Get(mock.Anything, "blocked:"+ipKey).Return("", cache.ErrKeyNotFound)
				c.EXPECT().Increment(mock.Anything, "rate_limit:"+ipKey, testWindowDuration).Return(testIPMaxRequests, nil)
			},
			wantErr: nil,
		},
		{
			name:     "por IP: excedeu o limite e bloqueia a chave",
			clientIP: testClientIP,
			apiKey:   "",
			setupMocks: func(c *mockcache.MockClient, _ *mockdb.MockTokenRepository) {
				c.EXPECT().Get(mock.Anything, "blocked:"+ipKey).Return("", cache.ErrKeyNotFound)
				c.EXPECT().Increment(mock.Anything, "rate_limit:"+ipKey, testWindowDuration).Return(testIPMaxRequests+1, nil)
				c.EXPECT().Set(mock.Anything, "blocked:"+ipKey, "blocked", testBlockDuration).Return(nil)
			},
			wantErr: ErrRateLimitted,
		},
		{
			name:     "por IP: chave ja bloqueada nao incrementa contador",
			clientIP: testClientIP,
			apiKey:   "",
			setupMocks: func(c *mockcache.MockClient, _ *mockdb.MockTokenRepository) {
				c.EXPECT().Get(mock.Anything, "blocked:"+ipKey).Return("blocked", nil)
			},
			wantErr: ErrRateLimitted,
		},

		// ---------- Caminho por Token ----------
		{
			name:     "por token: dentro do limite usa a chave do hash",
			clientIP: testClientIP,
			apiKey:   testAPIKey,
			setupMocks: func(c *mockcache.MockClient, tr *mockdb.MockTokenRepository) {
				tr.EXPECT().Get(mock.Anything, testAPIKey).
					Return(domain.Token{Hash: testTokenHash, TotalRequests: 100}, nil)
				c.EXPECT().Get(mock.Anything, "blocked:"+tokenKey).Return("", cache.ErrKeyNotFound)
				c.EXPECT().Increment(mock.Anything, "rate_limit:"+tokenKey, testWindowDuration).Return(int64(1), nil)
			},
			wantErr: nil,
		},
		{
			name:     "por token: limite do token sobrescreve o do IP",
			clientIP: testClientIP,
			apiKey:   testAPIKey,
			setupMocks: func(c *mockcache.MockClient, tr *mockdb.MockTokenRepository) {
				// token permite 100; o limite de IP (10) ja teria sido estourado com count=50
				tr.EXPECT().Get(mock.Anything, testAPIKey).
					Return(domain.Token{Hash: testTokenHash, TotalRequests: 100}, nil)
				c.EXPECT().Get(mock.Anything, "blocked:"+tokenKey).Return("", cache.ErrKeyNotFound)
				c.EXPECT().Increment(mock.Anything, "rate_limit:"+tokenKey, testWindowDuration).Return(int64(50), nil)
			},
			wantErr: nil,
		},
		{
			name:     "por token: excedeu o limite do token",
			clientIP: testClientIP,
			apiKey:   testAPIKey,
			setupMocks: func(c *mockcache.MockClient, tr *mockdb.MockTokenRepository) {
				tr.EXPECT().Get(mock.Anything, testAPIKey).
					Return(domain.Token{Hash: testTokenHash, TotalRequests: 100}, nil)
				c.EXPECT().Get(mock.Anything, "blocked:"+tokenKey).Return("", cache.ErrKeyNotFound)
				c.EXPECT().Increment(mock.Anything, "rate_limit:"+tokenKey, testWindowDuration).Return(int64(101), nil)
				c.EXPECT().Set(mock.Anything, "blocked:"+tokenKey, "blocked", testBlockDuration).Return(nil)
			},
			wantErr: ErrRateLimitted,
		},
		{
			name:     "por token: token inexistente retorna ErrInvalidAPIKey",
			clientIP: testClientIP,
			apiKey:   testAPIKey,
			setupMocks: func(_ *mockcache.MockClient, tr *mockdb.MockTokenRepository) {
				tr.EXPECT().Get(mock.Anything, testAPIKey).
					Return(domain.Token{}, domain.ErrTokenNotFound)
			},
			wantErr: ErrInvalidAPIKey,
		},
		{
			name:     "por token: erro de infra no repositorio propaga o erro original",
			clientIP: testClientIP,
			apiKey:   testAPIKey,
			setupMocks: func(_ *mockcache.MockClient, tr *mockdb.MockTokenRepository) {
				tr.EXPECT().Get(mock.Anything, testAPIKey).
					Return(domain.Token{}, errInfra)
			},
			wantErr: errInfra,
		},

		// ---------- Erros de infraestrutura ----------
		{
			name:     "increment falha propaga o erro",
			clientIP: testClientIP,
			apiKey:   "",
			setupMocks: func(c *mockcache.MockClient, _ *mockdb.MockTokenRepository) {
				c.EXPECT().Get(mock.Anything, "blocked:"+ipKey).Return("", cache.ErrKeyNotFound)
				c.EXPECT().Increment(mock.Anything, "rate_limit:"+ipKey, testWindowDuration).Return(int64(0), errInfra)
			},
			wantErr: errInfra,
		},
		{
			name:     "falha ao bloquear ainda retorna ErrRateLimitted",
			clientIP: testClientIP,
			apiKey:   "",
			setupMocks: func(c *mockcache.MockClient, _ *mockdb.MockTokenRepository) {
				c.EXPECT().Get(mock.Anything, "blocked:"+ipKey).Return("", cache.ErrKeyNotFound)
				c.EXPECT().Increment(mock.Anything, "rate_limit:"+ipKey, testWindowDuration).Return(testIPMaxRequests+1, nil)
				c.EXPECT().Set(mock.Anything, "blocked:"+ipKey, "blocked", testBlockDuration).Return(errInfra)
			},
			wantErr: ErrRateLimitted,
		},
		{
			name:     "fail-closed: erro no cache ao checar bloqueio trata como bloqueado",
			clientIP: testClientIP,
			apiKey:   "",
			setupMocks: func(c *mockcache.MockClient, _ *mockdb.MockTokenRepository) {
				c.EXPECT().Get(mock.Anything, "blocked:"+ipKey).Return("", errInfra)
			},
			wantErr: ErrRateLimitted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			cacheClient := mockcache.NewMockClient(t)
			tokenRepo := mockdb.NewMockTokenRepository(t)

			tt.setupMocks(cacheClient, tokenRepo)

			limiter := NewLimiterGateway(cacheClient, tokenRepo, newTestConfig())

			err := limiter.Check(ctx, tt.clientIP, tt.apiKey)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				return
			}
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
