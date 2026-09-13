package handlers

import (
	"github.com/dprio/redis-limiter/internal/infrastructure/web/handlers/orderhandler"
	"github.com/dprio/redis-limiter/internal/infrastructure/web/handlers/tokenhandler"
	"github.com/dprio/redis-limiter/internal/usecase"
)

type Handlers struct {
	CreateOrderHandler *orderhandler.OrderHandler
	CreateTokenHandler *tokenhandler.TokenHandler
}

func New(useCases usecase.UseCases) *Handlers {
	return &Handlers{
		CreateOrderHandler: orderhandler.New(useCases.CreateOrderUseCase, useCases.GetOrdersUseCase),
		CreateTokenHandler: tokenhandler.New(useCases.CreateTokenUseCase),
	}
}
