package kafka

import (
	"context"
	"fmt"
	"time"
	"log"

	"github.com/segmentio/kafka-go"
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
		Async:        false, // синхронная отправка
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	return &NotificationProducer{writer: writer, topic: topic}
}

func (p *NotificationProducer) Send(ctx context.Context, key string, value []byte) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
	}
	fmt.Printf("Sending message with key=%s to topic=%s\n", key, p.writer.Topic) 
	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		fmt.Printf("Failed to send message: %v\n", err) 
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}
	fmt.Println("Message sent successfully") 
	return nil
}

func (p *NotificationProducer) Close() error {
	log.Printf("Closing Kafka writer for topic %s\n", p.topic)
	return p.writer.Close()
}