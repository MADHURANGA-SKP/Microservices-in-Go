package main

import (
	pb "common/api"
	"common/kafka"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
)

type PaymentHTTPHAndler struct{
	consumer sarama.Consumer
}

func NewPaymentHTTPHandler(consumer sarama.Consumer) *PaymentHTTPHAndler {
	return &PaymentHTTPHAndler{consumer}
}

func (h *PaymentHTTPHAndler) registerRoutes(router *http.ServeMux) {
	router.HandleFunc("/webhook", h.handleCheckoutWebhook)
}

func (h *PaymentHTTPHAndler) handleCheckoutWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), endpointStripeSecret)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if session.PaymentStatus == "paid" {
			log.Printf("Payment for Checkout Session %v succeeded!", session.ID)

			orderID := session.Metadata["orderID"]
			customerID := session.Metadata["customerID"]

			_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			o := &pb.Order{
				ID:          orderID,
				CustomerID:  customerID,
				Status:      "paid",
				PaymentLink: "",
			}
			fmt.Println("------------------------------------------")
			marshalledOrder, err := json.Marshal(o)
			if err != nil {
				log.Fatal(err.Error())
			}

			// tr := otel.Tracer("amqp")
			// amqpContext, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - publish - %s", broker.OrderPaidEvent))
			// defer messageSpan.End()

			// headers := broker.InjectAMQPHeaders(amqpContext)

			// publish a message
			// h.channel.PublishWithContext(ctx, broker.OrderPaidEvent, "", false, false, amqp.Publishing{
			// 	ContentType:  "application/json",
			// 	Body:         marshalledOrder,
			// 	DeliveryMode: amqp.Persistent,
			// 	// Headers:      headers,
			// })

			err = kafka.PushOrderToQueue(serviceName, kafkaPort, marshalledOrder)
			if err != nil {
				log.Fatal(err)
			}
			data := []string{string(marshalledOrder)}
			fmt.Println("Message published order.paid",data )
		}
	}

	w.WriteHeader(http.StatusOK)
}

