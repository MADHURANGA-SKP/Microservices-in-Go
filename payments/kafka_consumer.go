package main

import (
	pb "common/api"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)


type consumer struct {
	service PaymentsService
}

func NewConsumer(service PaymentsService) *consumer {
	return &consumer{service}
}

func ConnectToKafka(broker, topic string) (*kafka.Consumer, string, error) {
    consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"broker.address.family": "v4",
		"group.id":         "payments-group",
		"auto.offset.reset": "earliest",
		"session.timeout.ms":    6000,
		"enable.auto.offset.store": false,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created Consumer %v\n", consumer)

	return consumer,topic,nil
}

func(o *consumer) Connect(topic, broker string) {
	consumer, _, err := ConnectToKafka(broker, topic)
	if err != nil {
		fmt.Println("failed to connect to kafka : ", err)
	}
	
    err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("Failed to get subscribed details: %s", err)
	}
	
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to subscribe to topics: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created Consumer 2%v\n", consumer)

	run := true

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				// Process the message received.
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}

				odr := pb.Order{}
				if err := json.Unmarshal(e.Value, &odr); err != nil {
					log.Printf("failed to unmarshal order: %v", err)
						continue
				}
				paymentLink, err := o.service.CreatePayments(context.Background(), &odr)
				if err != nil {
					log.Printf("failed to create payment: %v", err)
					continue
				}
				log.Printf("payment link created \n\n :%s\n\n",paymentLink)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
			default:
				fmt.Printf("Ignored %v\n", e)
		}
		}
	}
	fmt.Printf("Closing consumer\n")
	consumer.Close()
	
}