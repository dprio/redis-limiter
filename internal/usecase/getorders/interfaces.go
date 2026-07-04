package getorders

import (
	"context"

	"github.com/dprio/redis-limiter/internal/domain"
)

type (
	OrderRepository interface {
		GetAll(ctx context.Context) ([]domain.Order, error)
	}

	UseCase interface {
		Execute(ctx context.Context) ([]domain.Order, error)
	}
)
