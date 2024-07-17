package main

import (
	"common"
	"common/discovery"
	"common/discovery/consul"
	"context"
	"log"
	"net"
	"orders/gateway"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	serviceName = "orders"
	grpcAddr = common.EnvString("GRPC_ADDR", "localhost:2000")
	consulAddr = common.EnvString("CUNSL_ADDR","localhost:8500")
	kafkaPort    = common.EnvString("KAFKA_PORT", "localhost:29092")
	// amqpUser    = common.EnvString("RABBITMQ_USER", "guest")
	// amqpPass    = common.EnvString("RABBITMQ_PASS", "guest")
	// amqpHost    = common.EnvString("RABBITMQ_HOST", "localhost")
	// amqpPort    = common.EnvString("RABBITMQ_PORT", "5672")
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(ctx ,instanceID, serviceName); err != nil {
				logger.Error("Failed to health check", zap.Error(err))
			}
			time.Sleep(time.Second * 1)
		}
	}()

	_ , err = ConnectToKafka(kafkaPort)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	defer l.Close()

	store := NewStore()
	gateway := gateway.NewGateway(registry)
	svc := NewService(store, gateway)

	NewGRPCHandler(grpcServer, svc) 

	log.Println("GRPC server started at : ", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}