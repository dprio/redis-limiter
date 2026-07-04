package getorders

import (
	"context"
	"fmt"

	"github.com/dprio/redis-limiter/internal/domain"
)

type getOrders struct {
	orderRepository OrderRepository
}

func New(orderRepository OrderRepository) UseCase {
	return &getOrders{
		orderRepository: orderRepository,
	}
}

func (u *getOrders) Execute(ctx context.Context) ([]domain.Order, error) {
	fmt.Println("Executing getOrders")
	return u.orderRepository.GetAll(ctx)
}
