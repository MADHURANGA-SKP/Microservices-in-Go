package processor

import (
	pb "common/api"
)

type PaymentProcesser interface {
	CreatePaymentLink(*pb.Order) (string,error)
}