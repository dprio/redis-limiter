package createorder

import (
	"context"

	"github.com/dprio/redis-limiter/internal/domain"
	"github.com/google/uuid"
)

type createOrder struct {
	repository      OrderRepository
	eventDispatcher EventDispatcher
	eventCreator    EventCreator
}

func New(repository OrderRepository, eventDispatcher EventDispatcher, eventCreator EventCreator) UseCase {
	return &createOrder{
		repository:      repository,
		eventDispatcher: eventDispatcher,
		eventCreator:    eventCreator,
	}
}

func (co *createOrder) Execute(ctx context.Context, input Input) (Output, error) {
	order, err := domain.NewOrder(uuid.NewString(), input.Price, input.Tax)
	if err != nil {
		return Output{}, err
	}

	if err = order.CaluculateFinalPeice(); err != nil {
		return Output{}, err
	}

	savedOrder, err := co.repository.Save(ctx, order)
	if err != nil {
		return Output{}, err
	}

	out := NewOutput(savedOrder)

	if err = co.eventDispatcher.Dispatch(ctx, co.eventCreator.Create(out)); err != nil {
		return Output{}, err
	}

	return out, nil
}
