package kafka

import (
	"context"
	"strings"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	w *kafka.Writer
}

func NewProducer(brokersCSV, topic string) *Producer {
	brokers := strings.Split(brokersCSV, ",")
	return &Producer{
		w: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireAll,
			Async:        false,
		},
	}
}

func (p *Producer) Send(ctx context.Context, key, value []byte) error {
	return p.w.WriteMessages(ctx, kafka.Message{Key: key, Value: value})
}

func (p *Producer) Close() error { return p.w.Close() }
