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
	service OrdersService
}

func NewConsumer(service OrdersService) *consumer {
	return &consumer{service}
}

func ConnectToKafka(broker, topic string) (*kafka.Consumer, string, error) {
    consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"broker.address.family": "v4",
		"group.id":         "orders-group",
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

// func(o *consumer) Connect(topic, broker string, ptn int32){
// 	worker, err := ConnectToKafka(broker)
// 	if err != nil {
// 		panic(err)
// 	}

// 	consumer, err := worker.ConsumePartition(topic , ptn, sarama.OffsetOldest)
// 	if err != nil {
// 		panic(err)
// 	}

// 	sigchan := make(chan os.Signal, 1)
// 	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

// 	Count := 0
// 	// Get signal for finish
// 	doneCh := make(chan struct{})
// 	// msgChan := make(chan  *sarama.ConsumerMessage)
// 	go func() {
		
// 		for {
// 			select {
// 			case err := <-consumer.Errors():
// 				fmt.Println(err)
// 			case  broker := <-consumer.Messages():
// 				// Count++
// 				// fmt.Printf("Received  Broker key %s: | Topic(%s) | Message(%s) \n", string(broker.Key), string(broker.Topic), string(broker.Value))
// 				// msgChan <- broker
// 				odr := pb.Order{}
// 				if err := json.Unmarshal(broker.Value, &odr); err != nil {
// 					log.Printf("failed to unmarshal order: %v", err)
// 					continue
// 				}

// 				_, err := o.service.UpdateOrder(context.Background(), &odr)
// 				if err != nil {
// 					log.Printf("failed to update order: %v", err)
// 					continue
// 				}
// 			case <-sigchan:
// 				fmt.Println("Interrupt is detected")
// 			}
// 		}
// 	}()

// 	<- doneCh
// 	fmt.Println("Processed", Count, "")
// 	if err := worker.Close(); err != nil {
// 		panic(err)
// 	}
// }

func(o *consumer) Connect(topic, broker string){
	consumer, _, err := ConnectToKafka(broker, topic)
	if err != nil {
		fmt.Println("failed to connect to kafka : ", err)
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
					data, err := o.service.UpdateOrder(context.Background(), &odr)
					if err != nil {
							log.Printf("failed to update order: %v", err)
							continue
					}

					fmt.Println("data -------------------- :", data)

					_, err = consumer.StoreMessage(e)
					if err != nil {
						fmt.Fprintf(os.Stderr, "%% Error storing offset after message %s:\n",
							e.TopicPartition)
					}
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