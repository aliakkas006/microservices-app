package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aliakkas006/email-service/email"
	"github.com/aliakkas006/email-service/kafka"
	"github.com/aliakkas006/email-service/models"
	"github.com/joho/godotenv"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	brokers := os.Getenv("BOOTSTRAP_SERVERS")
	topic := os.Getenv("TOPIC_NAME")
	group := os.Getenv("GROUP_ID")
	if brokers == "" || topic == "" || group == "" {
		log.Fatal("BOOTSTRAP_SERVERS, TOPIC_NAME, and GROUP_ID must be set")
	}

	// Initialize Kafka consumer
	consumer := kafka.NewConsumer(brokers, topic, group)
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		log.Println("Shutting down gracefully...")
		cancel()
	}()

	log.Println("🚀 Email service started, waiting for todo events...")

	// Run consumer
	err := consumer.Run(ctx, func(key, value []byte) error {
		var todo models.Todo
		if err := json.Unmarshal(value, &todo); err != nil {
			log.Printf("⚠ Failed to unmarshal todo: %v", err)
			return nil
		}

		if err := email.SendEmail(todo); err != nil {
			log.Printf("⚠ Failed to send email: %v", err)
		}
		return nil
	})

	if err != nil && ctx.Err() == nil {
		log.Fatalf("Consumer error: %v", err)
	}
}
