package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	for _, word := range []string{string(message)} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)

	return err
}