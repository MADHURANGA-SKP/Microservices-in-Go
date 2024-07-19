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

	"github.com/IBM/sarama"
)


type consumer struct {
	service PaymentsService
}

func NewConsumer(service PaymentsService) *consumer {
	return &consumer{service}
}

func ConnectToKafka(broker string) (sarama.Consumer, error) {
    config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	conn, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		return nil, err
	}
	
	fmt.Println("conn \n\n", conn)
	return conn, nil
}

func(o *consumer) Connect(topic, broker string, ptn int32){
	worker, err := ConnectToKafka(broker)
	if err != nil {
		panic(err)
	}

	consumer, err := worker.ConsumePartition("orders" , ptn, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	Count := 0
	// Get signal for finish
	doneCh := make(chan struct{})
	// msgChan := make(chan  *sarama.ConsumerMessage)
	go func() {
		
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case  broker := <-consumer.Messages():
				odr := pb.Order{}
				if err := json.Unmarshal(broker.Value, &odr); err != nil {
					log.Printf("failed to unmarshal order: %v", err)
					continue
				}
				paymentLink, err := o.service.CreatePayments(context.Background(), &odr)
				if err != nil {
					log.Printf("failed to create payment: %v", err)
					continue
				}
				log.Printf("payment link created \n\n :%s\n\n",paymentLink)
			case <-sigchan:
				fmt.Println("Interrupt is detected")
				close(doneCh)
				return
			}
		}
	}()
	<- doneCh
	fmt.Println("Processed", Count, "")
}

