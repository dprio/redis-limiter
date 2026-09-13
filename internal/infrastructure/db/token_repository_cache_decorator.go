package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dprio/redis-limiter/internal/domain"
	"github.com/dprio/redis-limiter/internal/infrastructure/cache"
)

type tokenRepositoryCacheDecorator struct {
	tokenRepository TokenRepository
	cacheClient     cache.Client
}

func NewTokenRepositoryCacheDecorator(tokenRepository TokenRepository, cacheClient cache.Client) TokenRepository {
	return &tokenRepositoryCacheDecorator{
		tokenRepository: tokenRepository,
		cacheClient:     cacheClient,
	}
}

func (repo *tokenRepositoryCacheDecorator) Save(ctx context.Context, totalRequests int64) (domain.Token, error) {
	token, err := repo.tokenRepository.Save(ctx, totalRequests)
	if err != nil {
		return domain.Token{}, err
	}

	if err := repo.writeCache(ctx, token); err != nil {
		// escrita no banco já aconteceu com sucesso — não falha a operação por causa do cache,
		// só loga (idealmente com um logger injetado)
		fmt.Printf("erro ao popular cache no save: %v\n", err)
	}

	return token, nil
}

func (repo *tokenRepositoryCacheDecorator) Update(ctx context.Context, token string, totalRequests int64) error {
	if err := repo.tokenRepository.Update(ctx, token, totalRequests); err != nil {
		return err
	}

	t := domain.Token{
		Hash:          token,
		TotalRequests: totalRequests,
	}

	if err := repo.writeCache(ctx, t); err != nil {
		fmt.Printf("erro ao atualizar cache no update: %v\n", err)
	}

	return nil
}

func (repo *tokenRepositoryCacheDecorator) Get(ctx context.Context, token string) (domain.Token, error) {
	cacheKey := tokenCacheKey(token)

	cached, err := repo.cacheClient.Get(ctx, cacheKey)
	if err == nil {
		var t domain.Token
		if unmarshalErr := json.Unmarshal([]byte(cached), &t); unmarshalErr == nil {
			return t, nil
		}
		// cache corrompido -> ignora e segue pro banco
	}

	// cache miss (ou erro no cache) -> busca no MySQL
	t, err := repo.tokenRepository.Get(ctx, token)
	if err != nil {
		return domain.Token{}, err
	}

	if writeErr := repo.writeCache(ctx, t); writeErr != nil {
		fmt.Printf("erro ao popular cache no get: %v\n", writeErr)
	}

	return t, nil
}

func (repo *tokenRepositoryCacheDecorator) writeCache(ctx context.Context, t domain.Token) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return repo.cacheClient.Set(ctx, tokenCacheKey(t.Hash), string(data), 0*time.Second)
}

func tokenCacheKey(token string) string {
	return fmt.Sprintf("token:%s", token)
}
