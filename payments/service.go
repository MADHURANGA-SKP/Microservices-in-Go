package main

import (
	pb "common/api"
	"context"
	"payments/processor"
)

type Service struct {
	processor processor.PaymentProcesser
}

func NewService(processor processor.PaymentProcesser) *Service {
	return &Service{processor}
}

func (s *Service) CreatePayments(ctx context.Context, o *pb.Order) (string, error){
	link, err := s.processor.CreatePaymentLink(o)
	if err != nil {
		return "", err
	}

	//update order with the link

	return link, nil
}