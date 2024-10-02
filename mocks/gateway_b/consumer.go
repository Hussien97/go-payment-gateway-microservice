package main

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func StartKafkaConsumer(topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"kafka-like:9092"},
		Topic:    topic,
		GroupID:  "gateway_b_group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	defer func() {
		if err := r.Close(); err != nil {
			log.Fatal("Failed to close reader:", err)
		}
	}()

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			continue
		}
		log.Printf("Received message on Gateway B: %s\n", string(msg.Value))

		ProcessTransaction(msg.Value)
	}
}
