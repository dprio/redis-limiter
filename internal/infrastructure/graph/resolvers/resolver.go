package resolvers

import (
	"github.com/dprio/redis-limiter/internal/usecase/createorder"
	"github.com/dprio/redis-limiter/internal/usecase/getorders"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	createOrderUseCase createorder.UseCase
	getOrdersUseCase   getorders.UseCase
}
