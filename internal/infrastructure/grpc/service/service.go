package service

import (
	"github.com/dprio/redis-limiter/internal/infrastructure/grpc/pb"
	"github.com/dprio/redis-limiter/internal/usecase"
)

type GRPCServices struct {
	OrderService pb.OrderServiceServer
}

func NewGRPCServices(useCases *usecase.UseCases) *GRPCServices {
	return &GRPCServices{
		OrderService: NewOrderService(useCases.CreateOrderUseCase, useCases.GetOrdersUseCase),
	}
}
