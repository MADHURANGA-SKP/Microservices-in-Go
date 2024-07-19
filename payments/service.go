package main

import (
	pb "common/api"
	"context"
	"fmt"
	"payments/gateway"
	"payments/processor"
)

type Service struct {
	processor processor.PaymentProcesser
	gateway gateway.OrdersGateway
}

func NewService(processor processor.PaymentProcesser, gateway gateway.OrdersGateway) *Service {
	return &Service{processor, gateway}
}

func (s *Service) CreatePayments(ctx context.Context, o *pb.Order) (string, error){
	fmt.Printf("user--CreatePaymentLink-------------------------\n\n%s\n\n",o)
	link, err := s.processor.CreatePaymentLink(o)
	if err != nil {
		return "", err
	}

	//update order with the link
	err = s.gateway.UpdateOrderAfterPaymentLink(ctx, o.ID, link)
	if err != nil {
		return "", err
	}
	
	return link, nil
}