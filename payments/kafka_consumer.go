package main

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)


type consumer struct {
	service PaymentsService
}

func NewConsumer(service PaymentsService) *consumer {
	return &consumer{service}
}

func(o *consumer) Connect(topic, broker string)  (*kafka.Consumer, error){
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":         "orders",
		"go.application.rebalance.enable": true,
	})
	if err != nil {
		fmt.Printf("ann error occred ----- %v", err)
		return nil, err
	}
	
    err = consumer.SubscribeTopics([]string{"orders"}, nil)
	if err != nil {
		fmt.Printf("Failed to get subscribed details: %s", err)
	}
	
	defer consumer.Close()

	fmt.Println("pulling...")
	for {
		message, err := consumer.ReadMessage(-1)
		if err != nil {
			fmt.Printf("pulling error: %v\n", err.Error())
			return nil ,err
		} else {
			fmt.Printf("message %s: %s \n", message.TopicPartition, string(message.Value))

			for _, header := range message.Headers {
				fmt.Printf("header %s - [%s:%s]\n", message.TopicPartition, header.Key, string(header.Value))
			}

			return consumer,nil
		}
	}

}

