package main

// import (
// 	"common/api"
// 	inmemRegistry "common/discovery/inmem"
// 	"context"
// 	"payments/gateway"
// 	"payments/processor/inmem"
// 	"testing"
// )

// func TestService(t *testing.T){
// 	processor := inmem.NewInmem()
// 	registry := inmemRegistry.NewRegistry()

// 	gateway := gateway.NewGateway(registry)
// 	svc := NewService(processor, gateway)

// 	t.Run("shoud create a payment link", func(t *testing.T) {
// 		link, err := svc.CreatePayments(context.Background(), &api.Order{})
// 		if err != nil {
// 			t.Errorf("createPayment() error = %v, ", err)
// 		}
// 		if link == "" {
// 			t.Error("createPayment link u=is empty")
// 		}

// 	})

// }