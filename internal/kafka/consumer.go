package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"notification_service/internal/models"
)

type NotificationProcessor interface {
	Handle(ctx context.Context, notification *models.Notification) error
}

type NotificationConsumer struct {
	processor   NotificationProcessor
	kafkaBroker []string
	topicName   string
}

func NewNotificationConsumer(processor NotificationProcessor, brokers []string, topic string) *NotificationConsumer {
	return &NotificationConsumer{
		processor:   processor,
		kafkaBroker: brokers,
		topicName:   topic,
	}
}

func (c *NotificationConsumer) Start(ctx context.Context) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.kafkaBroker,
		Topic:       c.topicName,
		StartOffset: kafka.FirstOffset,
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var notification models.Notification
		if err := json.Unmarshal(msg.Value, &notification); err != nil {
			log.Println("Failed to unmarshal message:", err)
			continue
		}

		if err := c.processor.Handle(ctx, &notification); err != nil {
			log.Println("Processor error:", err)
		}
	}
}