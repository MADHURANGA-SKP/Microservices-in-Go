package main

import (
	pb "common/api"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	confluentinc "github.com/confluentinc/confluent-kafka-go/kafka"
)

type consumer struct {
	service OrdersService
}

func NewConsumer(service OrdersService) *consumer {
	return &consumer{service}
}

// func ConnectToKafka(broker string) (sarama.Consumer, error) {
//     config := sarama.NewConfig()
// 	config.Consumer.Return.Errors = true
// 	conn, err := sarama.NewConsumer([]string{broker}, config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println("conn", conn)
// 	return conn, nil
// }

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

func(o *consumer) Connect(topic, broker string, ptn int32) (confluentinc.Consumer, error) {
	c, err := confluentinc.NewConsumer(&confluentinc.ConfigMap{
		"bootstrap.servers": broker,
	})

	if err != nil {
		panic(err)
	}

	err = c.SubscribeTopics([]string{topic, "^aRegex.*[Tt]opic"}, nil)

	if err != nil {
		panic(err)
	}

	sigchan := make(chan os.Signal, 1)
// 	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	Count := 0
	// Get signal for finish
	doneCh := make(chan struct{})
	// msgChan := make(chan  *sarama.ConsumerMessage)
	go func() {
		
		for {
			select {
			case sig := <-sigchan:
				fmt.Printf("Caught signal %v: terminating\n", sig)
				close(doneCh)
				return
			default:
				ev := c.Poll(100)
				switch e := ev.(type) {
				case *confluentinc.Message:
					odr := pb.Order{}
					if err := json.Unmarshal(e.Value, &odr); err != nil {
						log.Printf("failed to unmarshal order: %v", err)
						continue
					}

					_, err := o.service.UpdateOrder(context.Background(), &odr)
					if err != nil {
						log.Printf("failed to update order: %v", err)
						continue
					}
				case confluentinc.Error:
					fmt.Fprintf(os.Stderr, "Error: %v\n", e)
				}
			}
		}
	}()

	<- doneCh
	fmt.Println("Processed", Count, "")
	if err := c.Close(); err != nil {
		panic(err)
	}

	return confluentinc.Consumer{}, nil
}