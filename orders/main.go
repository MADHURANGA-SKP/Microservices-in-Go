package main

import (
	"common"
	"common/discovery"
	"common/discovery/consul"
	"context"
	"fmt"
	"log"
	"net"
	"orders/gateway"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	serviceName = "orders"
	grpcAddr = common.EnvString("GRPC_ADDR", "localhost:2000")
	consulAddr = common.EnvString("CUNSL_ADDR","localhost:8500")
	kafkaPort    = common.EnvString("KAFKA_PORT", "localhost:29092")
	mongoUser   = common.EnvString("MONGO_DB_USER", "root")
	mongoPass   = common.EnvString("MONGO_DB_PASS", "example")
	mongoAddr   = common.EnvString("MONGO_DB_HOST", "localhost:27017")
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

	ch, _ , err := ConnectToKafka(kafkaPort, serviceName)
	if err != nil {
		panic(err)
	}


	// mongo db conn
	uri := fmt.Sprintf("mongodb://%s:%s@%s", mongoUser, mongoPass, mongoAddr)
	mongoClient, err := connectToMongoDB(uri)
	if err != nil {
		logger.Fatal("failed to connect to mongo db", zap.Error(err))
	}

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}
	defer l.Close()

	gateway := gateway.NewGateway(registry)
	store := NewStore(mongoClient)
	svc := NewService(store, gateway)
	svcWithLogging := NewLoggingMiddleware(svc)

	NewGRPCHandler(grpcServer, svcWithLogging, ch) 

	consumer := NewConsumer(svc)
	go consumer.Connect(serviceName, kafkaPort)

	logger.Info("Starting HTTP server", zap.String("port", grpcAddr))

	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}

func connectToMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	return client, err
}