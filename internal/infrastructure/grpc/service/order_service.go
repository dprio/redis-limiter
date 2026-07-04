package service

import (
	"context"

	"github.com/dprio/redis-limiter/internal/infrastructure/grpc/pb"
	"github.com/dprio/redis-limiter/internal/usecase/createorder"
	"github.com/dprio/redis-limiter/internal/usecase/getorders"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase createorder.UseCase
	GetOrdersUseCase   getorders.UseCase
}

func NewOrderService(createOrderUseCase createorder.UseCase, getOrdersUseCase getorders.UseCase) pb.OrderServiceServer {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		GetOrdersUseCase:   getOrdersUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := createorder.Input{
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(ctx, dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) GetOrders(ctx context.Context, _ *pb.Empty) (*pb.OrdersResponse, error) {
	output, err := s.GetOrdersUseCase.Execute(ctx)
	if err != nil {
		return nil, err
	}

	ordersResponse := make([]*pb.CreateOrderResponse, len(output))
	for i, out := range output {
		ordersResponse[i] = &pb.CreateOrderResponse{
			Id:         out.ID,
			Price:      float32(out.Price),
			Tax:        float32(out.Tax),
			FinalPrice: float32(out.FinalPrice),
		}
	}

	return &pb.OrdersResponse{
		Orders: ordersResponse,
	}, nil
}
