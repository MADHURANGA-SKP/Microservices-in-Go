package main

import (
	pb "common/api"
	"context"
)

type OrdersService interface{
	CreateOrder(context.Context) error
	ValidateOrder(context.Context, *pb.CreateOrderRequest) error
}

type OrderStore interface{
	Create(context.Context) error
}