package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"notification_service/internal/models"
)

type NotificationProducer struct {
	writer *kafka.Writer
	topic  string
}

func NewNotificationProducer(brokers []string, topic string) *NotificationProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{}, 
		RequiredAcks: kafka.RequireAll,     
		Async:        false,               
		WriteTimeout: 10 * time.Second,
	}
	return &NotificationProducer{
		writer: writer,
		topic:  topic,
	}
}

func (p *NotificationProducer) Send(ctx context.Context, key string, value []byte) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
	}

	log.Printf("Sending message to topic %s: key=%s value=%s\n", p.topic, key, string(value))

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("Failed to send message to topic %s: %v\n", p.topic, err)
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	log.Printf("Message successfully sent to topic %s: key=%s\n", p.topic, key)
	return nil
}

func (p *NotificationProducer) SendNotification(ctx context.Context, n *models.Notification) error {
	value, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}
	return p.Send(ctx, n.ID, value)
}

func (p *NotificationProducer) Close() error {
	log.Printf("Closing Kafka writer for topic %s\n", p.topic)
	return p.writer.Close()
}