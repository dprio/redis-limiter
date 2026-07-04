package resolvers

import "github.com/dprio/redis-limiter/internal/usecase"

type Resolvers struct {
	OrderResolver *Resolver
}

func NewGraphQLResolvers(useCases *usecase.UseCases) *Resolvers {
	return &Resolvers{
		OrderResolver: &Resolver{
			createOrderUseCase: useCases.CreateOrderUseCase,
			getOrdersUseCase:   useCases.GetOrdersUseCase,
		},
	}

}
