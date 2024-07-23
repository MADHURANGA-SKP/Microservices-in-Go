package main

import (
	"common"
	"common/discovery"
	"common/discovery/consul"
	"context"
	"log"
	"net"
	"net/http"
	"payments/gateway"
	stripeProcessor "payments/processor/stripe"
	"time"

	"github.com/stripe/stripe-go/v78"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	serviceName = "payment"
	grpcAddr = common.EnvString("GRPC_ADDR", "localhost:2001")
	httpAddr = common.EnvString("HTTP_ADDR", "localhost:8081")
	consulAddr = common.EnvString("CUNSL_ADDR","localhost:8500")
	// amqpUser    = common.EnvString("RABBITMQ_USER", "guest")
	// amqpPass    = common.EnvString("RABBITMQ_PASS", "guest")
	// amqpHost    = common.EnvString("RABBITMQ_HOST", "localhost")
	// amqpPort    = common.EnvString("RABBITMQ_PORT", "5672")
	stripeKey    = common.EnvString("STRIPE_KEY", "sk_test_51Pcm84Rwz2fOefiNmTQSI0dSY5eUIMv84XF9lg862TBJ8SzTfAbC5MuIfn6pWB4vSpufr5FJxOdiqUmZ1fIUdj6G006GK0lcq1")
	endpointStripeSecret = common.EnvString("STRIPE_ENDPOINT_SECRET", "whsec_32f21a7f8dc9d3af672bfab02ede0f6f22efd7fb91de868a3c47799d5855e706")
	kafkaPort    = common.EnvString("KAFKA_PORT", "localhost:29092")
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
	defer registry.DeRegister(ctx, instanceID, serviceName)

	stripe.Key = stripeKey
	// ch, close := broker.Connect(amqpUser,amqpPass, amqpHost, amqpPort)
	// defer func () {
	// 	close()
	// 	ch.Close()
	// } ()

	ch, _ , err := ConnectToKafka(kafkaPort, serviceName)
	if err != nil {
		panic(err)
	}

	stripeProcessor:= stripeProcessor.NewProcessor()
	gateway := gateway.NewGateway(registry)
	svc := NewService(stripeProcessor, gateway)
	

	mux := http.NewServeMux()
	httpServer := NewPaymentHTTPHandler(ch)
	httpServer.registerRoutes(mux)

	go func() {
		log.Printf("Stating http server at %s: ", httpAddr)
		if err := http.ListenAndServe(httpAddr, mux); err != nil {
			log.Fatal("faild to start http server")
		}
	}()

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	kafkaConsumer := NewConsumer(svc)
	go kafkaConsumer.Connect("orders", kafkaPort)

	log.Println("GRPC Server Started at ", grpcAddr)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}