package main

import (
	pb "common/api"
	"context"
	"log"

	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrdersService
}

func NewGRPCHandler(grpcServer *grpc.Server, service OrdersService){
	handler := &grpcHandler{
		service: service,
	}

	pb.RegisterOrderServiceServer(grpcServer, handler)

}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {	
	log.Printf("New order recived! order %v", p)
	o := &pb.Order{
		ID: "42",
	}

	return o, nil
}