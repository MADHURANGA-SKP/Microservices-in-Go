package main

import (
	pb "common/api"
	"context"
	"errors"
)

var orders = make([]*pb.Order, 0)

type store struct {
	
}

func NewStore() *store {
	return &store{}
}

func (s *store) Create(ctx context.Context, p *pb.CreateOrderRequest, item []*pb.Item) (string, error) {
	id := "42"
	orders = append(orders, &pb.Order{
		ID: id,
		CustomerID: p.CustomerID,
		Status: "pending",
		Items: item,
		PaymentLink: "",
	})

	return id, nil
}

func (s *store) Get(ctx context.Context, id, customerID string) (*pb.Order, error){
	 for _, o := range orders {
		if o.ID == id && o.CustomerID == customerID {
			return o, nil
		}
	 }

	 return nil, errors.New("order not found")
}

func (s *store) Update(ctx context.Context, id string, newOrder *pb.Order) error {
	for i, order := range orders {
		if order.ID == id {
			orders[i].Status = newOrder.Status
			orders[i].PaymentLink =newOrder.PaymentLink


			return nil
		}
	}

	return nil
}
