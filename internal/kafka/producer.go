package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type NotificationProducer struct {
	writer *kafka.Writer
}

func NewNotificationProducer(brokers []string, topic string) *NotificationProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &NotificationProducer{writer: writer}
}

func (p *NotificationProducer) Send(ctx context.Context, key string, value []byte) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
	}
	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}
	return nil
}

func (p *NotificationProducer) Close() error {
	return p.writer.Close()
}
