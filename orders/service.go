package main

import (
	"common"
	pb "common/api"
	"context"
	"log"
	"orders/gateway"
)

type Service struct {
	store OrderStore
	gateway gateway.StockGateway
}

func NewService(store OrderStore, gateway gateway.StockGateway) *Service {
	return &Service{store,gateway}
}

func (s *Service) UpdateOrder(ctx context.Context, o *pb.Order) (*pb.Order, error) {
	err := s.store.Update(ctx, o.ID, o)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (s *Service) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	return s.store.Get(ctx, p.OrderID, p.CustomerID)
}

func(s *Service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {
	id , err := s.store.Create(ctx, p, items)
	if err != nil {
		return nil, err
	}

	o := &pb.Order{
		ID: id,
		CustomerID: p.CustomerID,
		Status: "pennding",
		Items: items,
	}

	return o,nil
}

func (s *Service) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error) {
	if len(p.Items) == 0 {
		return nil, common.ErrNoItems
	}

	mergedItems := mergeItemsQuantities(p.Items)
	log.Print(mergedItems)

	//validate with stock service 
	var itemsWithPrice []*pb.Item
	for _, i := range mergedItems {
		itemsWithPrice = append(itemsWithPrice, &pb.Item{
			PriceID: "price_1PcmBtRwz2fOefiNbNnE4VrB",
			ID: i.ID,
			Quantity: i.Quantity,
		})
	}

	return itemsWithPrice, nil
}

func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity,0)

	for _, item := range items {
		found := false
		for _, findItem := range merged {
			if findItem.ID == item.ID {
				findItem.Quantity += item.Quantity
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, item)
		}
	}
	
	return merged
}