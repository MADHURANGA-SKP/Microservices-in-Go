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

func(s *Service) CreateOrder(context.Context) error {
	return nil
}

func (s *Service) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) error {
	if len(p.Items) == 0 {
		return common.ErrNoItems
	}

	mergedItems := mergeItemsQuantities(p.Items)
	log.Print(mergedItems)


	return nil
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