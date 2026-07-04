package grpcserver

import (
	"net"

	"github.com/dprio/redis-limiter/internal/infrastructure/grpc/pb"
	"github.com/dprio/redis-limiter/internal/infrastructure/grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server *grpc.Server
}

func New(services *service.GRPCServices) *GRPCServer {
	serv := &GRPCServer{
		server: grpc.NewServer(),
	}

	pb.RegisterOrderServiceServer(serv.server, services.OrderService)
	reflection.Register(serv.server)

	return serv
}

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	return s.server.Serve(lis)
}
