package main

import (
	pb "common/api"
	kfk "common/kafka"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrdersService
	consumer *kafka.Consumer
}

func NewGRPCHandler(grpcServer *grpc.Server, service OrdersService, consumer *kafka.Consumer ){
	handler := &grpcHandler{
		service: service,
		consumer: consumer,
	}

	pb.RegisterOrderServiceServer(grpcServer, handler)

}


func (h *grpcHandler) UpdateOrder(ctx context.Context, p *pb.Order) (*pb.Order, error) {
	return h.service.UpdateOrder(ctx, p)
}

func (h *grpcHandler) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error){
	return h.service.GetOrder(ctx, p)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {	
	log.Printf("New order recived! order %v", p)
	
	items, err := h.service.ValidateOrder(ctx, p)
	if err != nil {
		return nil, err
	}

	o, err := h.service.CreateOrder(ctx, p, items)
	if err != nil {
		return nil, err
	}

	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	fmt.Println("marshalledOrder",marshalledOrder)
	err = kfk.PushOrderToQueue(serviceName, kafkaPort, marshalledOrder)
	if err != nil {
		return nil, err
	}

	return o, nil
}

