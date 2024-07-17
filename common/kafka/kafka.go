package kafka

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return sarama.NewSyncProducer(brokers, config)
}

func PushOrderToQueue(topic, broker string, message []byte) error {
	brokers := []string{broker}
	// Create connection
	producer, err := ConnectProducer(brokers)
	if err != nil {
		return err
	}

	defer producer.Close()

	showdata := string(message)
	fmt.Println("----------------------------------------------------")
	log.Printf("stored data are %s:", showdata)
	fmt.Println("----------------------------------------------------")
	// Create a new message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	// Send message
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Order is stored in topic(%s)/partition(%d)/offset(%d)\n",
		topic,
		partition,
		offset)

	return nil
}



// func ConnectToKafka(ctx context.Context, topic, broker string) {
// 	config := sarama.NewConfig()
// 	config.Consumer.Return.Errors = true

// 	// 1. Create a new consumer and start it.
// 	worker, err := sarama.NewConsumer([]string{broker}, config)
// 	if err != nil {
// 		panic(err)
// 	}

// 	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Consumer started ")

// 	// 2. Handle OS signals - used to stop the process.
// 	sigchan := make(chan os.Signal, 1)
// 	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
// 	msgCnt := 0
// 	// 3. Create a Goroutine to run the consumer / worker.
// 	doneCh := make(chan struct{})
// 	go func() {
// 		for {
// 			select {
// 			case err := <-consumer.Errors():
// 				fmt.Println(err)
// 			case msg := <-consumer.Messages():
// 				msgCnt++
// 				fmt.Printf("Received order Count %d: | Topic(%s) | Message(%s) \n", msgCnt, string(msg.Topic), string(msg.Value))
// 			case <-ctx.Done():
// 				fmt.Println("Interrupt is detected")
// 				doneCh <- struct{}{}
// 			}
// 		}
// 	}()

// 	<-doneCh
// 	fmt.Println("Processed", msgCnt, "messages")

// }