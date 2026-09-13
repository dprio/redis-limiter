package usecase

import (
	"github.com/dprio/redis-limiter/internal/domain/eventtype"
	"github.com/dprio/redis-limiter/internal/infrastructure/db"
	"github.com/dprio/redis-limiter/internal/usecase/createorder"
	"github.com/dprio/redis-limiter/internal/usecase/createtoken"
	"github.com/dprio/redis-limiter/internal/usecase/getorders"
	"github.com/dprio/redis-limiter/pkg/events"
)

type UseCases struct {
	CreateOrderUseCase createorder.UseCase
	GetOrdersUseCase   getorders.UseCase
	CreateTokenUseCase createtoken.UseCase
}

func New(dbs *db.DBs, eventDispatcher events.EventDispatcherInterface) *UseCases {
	return &UseCases{
		CreateOrderUseCase: buildCreateOrderUseCase(dbs, eventDispatcher),
		GetOrdersUseCase:   getorders.New(dbs.OrderRepository),
		CreateTokenUseCase: createtoken.NewCreateTokenUseCase(dbs.TokenRepository),
	}
}

func buildCreateOrderUseCase(dbs *db.DBs, eventDispatcher createorder.EventDispatcher) createorder.UseCase {
	eventCreator := events.NewEventCreator(eventtype.OrderCreated)
	return createorder.New(dbs.OrderRepository, eventDispatcher, eventCreator)
}
