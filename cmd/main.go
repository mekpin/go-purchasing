package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"go-purchasing/config"
	"go-purchasing/model"
	"log"
	"time"

	"github.com/rs/xid"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func init() {

}

func main() {
	fmt.Println("kafka publish test")
	// Initialize the Kafka writer
	mechanism := plain.Mechanism{
		Username: config.Common.KafkaUsername,
		Password: config.Common.KafkaPassword,
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
		TLS: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{config.Common.KafkaUrl},
		Topic:    "mekpin-perlu-belanja",
		Balancer: &kafka.LeastBytes{},
		Dialer:   dialer,
	})

	cart := model.Cart{
		Id:          xid.New(),
		NamaItem:    "bawang",
		Jumlah:      10,
		HargaSatuan: 1000,
		HargaTotal:  10 * 1000,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	cart2 := model.Cart{
		Id:          xid.New(),
		NamaItem:    "garem",
		Jumlah:      1,
		HargaSatuan: 500,
		HargaTotal:  500,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	// Convert the Cart instance to JSON
	messageBody, err := json.Marshal(cart)
	if err != nil {
		log.Fatal("Error encoding message body:", err)
	}
	// Convert the Cart instance to JSON
	messageBody2, err := json.Marshal(cart2)
	if err != nil {
		log.Fatal("Error encoding message body:", err)
	}

	// Create a Kafka message
	message := kafka.Message{
		Value: messageBody,
	}
	message2 := kafka.Message{
		Value: messageBody2,
	}
	// Publish the message to the Kafka topic
	err = writer.WriteMessages(context.Background(), message, message2)
	if err != nil {
		log.Fatal("Error publishing message:", err)
	}

	fmt.Println("Message published successfully!")

	// Close the Kafka writer
	err = writer.Close()
	if err != nil {
		log.Fatal("Error closing Kafka writer:", err)
	}
}
