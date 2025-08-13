package kafka

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

// Consumer struct
type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokersCSV, topic, group string) *Consumer {
	brokers := strings.Split(brokersCSV, ",")
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        group,
			MinBytes:       1,
			MaxBytes:       10e6,
			CommitInterval: time.Second, 
		}),
	}
}

// Run consumes messages and calls the handler
func (c *Consumer) Run(ctx context.Context, handler func(key, value []byte) error) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		if err := handler(msg.Key, msg.Value); err != nil {
			log.Printf("Error handling message: %v", err)
		}
	}
}

// Close closes the consumer
func (c *Consumer) Close() error {
	return c.reader.Close()
}
