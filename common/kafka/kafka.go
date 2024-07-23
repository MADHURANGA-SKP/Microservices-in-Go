package kfk

import (
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
// 	config := sarama.NewConfig()
// 	config.Producer.Return.Successes = true
// 	config.Producer.RequiredAcks = sarama.WaitForAll
// 	config.Producer.Retry.Max = 5

// 	return sarama.NewSyncProducer(brokers, config)
// }

// func PushOrderToQueue(topic, broker string, message []byte) error {
// 	brokers := []string{broker}
// 	// Create connection
// 	producer, err := ConnectProducer(brokers)
// 	if err != nil {
// 		return err
// 	}

// 	defer producer.Close()

// 	showdata := string(message)
// 	fmt.Println("----------------------------------------------------")
// 	log.Printf("stored data are %s:", showdata)
// 	fmt.Println("----------------------------------------------------")
// 	// Create a new message
// 	msg := &sarama.ProducerMessage{
// 		Topic: topic,
// 		Value: sarama.StringEncoder(message),
// 	}

// 	// Send message
// 	partition, offset, err := producer.SendMessage(msg)
// 	if err != nil {
// 		return err
// 	}

// 	log.Printf("Order is stored in topic(%s)/partition(%d)/offset(%d)\n",
// 		topic,
// 		partition,
// 		offset)

// 	return nil
// }

func PushOrderToQueue(topic, broker string, message []byte) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	defer p.Close()

	fmt.Println("producer ", p)

	// Event handling
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
			case kafka.Error:
				fmt.Printf("Error: %v\n", ev)
			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()

	for {
		fmt.Println("Producer message ", message)

		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          message,
		}, nil)

		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrQueueFull {
				// Producer queue is full, wait 1s for messages to be delivered then try again.
				time.Sleep(time.Second)
				continue
			}
			fmt.Printf("Failed to produce message: %v\n", err)
			return err
		}
		break
	}

	// Flush and close the producer and the events channel
	p.Flush(10000)
	p.Close()

	return err 
}