package main

import (
	common "common"
	"common/discovery"
	"common/discovery/consul"
	"context"
	"gatway/gateway"
	"log"
	"net/http"
	"time"
)

var (
	serviceName = "gateway"
	httpAddr = common.EnvString("HTTP_ADDR",":8080")
	consulAddr = common.EnvString("CUNSL_ADDR","localhost:8500")
)

func main(){
	registry, err := consul.NewRegistory(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
		panic(err)
	}

	go func ()  {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("failed to health check")
			}
			time.Sleep(time.Second * 1)
		}
	} ()

	defer registry.DeRegister(ctx, instanceID, serviceName)

	mux := http.NewServeMux()


	OrdersGateway := gateway.NewGRPCGatway(registry)

	handler := NewHandler(OrdersGateway)
	handler.registerRoutes(mux)

	log.Printf("starting http server at %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("failed to start http server")
	}
}