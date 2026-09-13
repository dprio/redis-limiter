package createtoken

import (
	"context"
	"fmt"

	"github.com/dprio/redis-limiter/internal/infrastructure/db"
)

type (
	UseCase interface {
		Execute(ctx context.Context, input Input) (Output, error)
	}
)

type createTokenUseCase struct {
	tokenRepository db.TokenRepository
}

func NewCreateTokenUseCase(tokenRepository db.TokenRepository) UseCase {
	return &createTokenUseCase{
		tokenRepository: tokenRepository,
	}
}

func (uc *createTokenUseCase) Execute(ctx context.Context, input Input) (Output, error) {
	if input.TotalRequests <= 0 {
		return Output{}, ErrInvalidTotalRequests
	}

	token, err := uc.tokenRepository.Save(ctx, input.TotalRequests)
	if err != nil {
		return Output{}, fmt.Errorf("erro ao salvar token: %w", err)
	}

	return Output{
		Token:         token.Hash,
		TotalRequests: input.TotalRequests,
	}, nil
}
