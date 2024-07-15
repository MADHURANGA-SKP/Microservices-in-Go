package main

import (
	pb "common/api"
	"context"
)

type PaymentsService interface {
	CreatePayments(context.Context, *pb.Order) (string, error)
}