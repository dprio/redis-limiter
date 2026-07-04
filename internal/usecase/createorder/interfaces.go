package createorder

import (
	"context"

	"github.com/dprio/redis-limiter/internal/domain"
	"github.com/dprio/redis-limiter/pkg/events"
)

type (
	OrderRepository interface {
		Save(ctx context.Context, order *domain.Order) (domain.Order, error)
	}

	EventCreator interface {
		Create(payload any) events.Event
	}

	EventDispatcher interface {
		Dispatch(ctx context.Context, event events.Event) error
	}

	UseCase interface {
		Execute(ctx context.Context, input Input) (Output, error)
	}
)
