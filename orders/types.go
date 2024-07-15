package main

import (
	pb "common/api"
	"context"
)

type OrdersService interface{
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
}

type OrderStore interface{
	Create(context.Context) error
}