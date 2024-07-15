package main

import (
	"common"
	pb "common/api"
	"context"
	"log"
)

type Service struct {
	store OrderStore
}

func NewService(store OrderStore) *Service {
	return &Service{store}
}

func(s *Service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	items, err := s.ValidateOrder(ctx, p)
	if err != nil {
		return nil, err
	}

	o := &pb.Order{
		ID: "42",
		CustomerID: p.CustomerID,
		Status: "pending",
		Items: items,
	}

	return o, nil
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