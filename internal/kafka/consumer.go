package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"notification_service/internal/models"
)

type NotificationProcessor interface {
	HandleNotification(ctx context.Context, n *models.Notification) error
	HandleNotificationStatus(ctx context.Context, status *models.NotificationStatus) error
}

type NotificationConsumer struct {
	processor   NotificationProcessor
	kafkaBroker []string
	topics      []string
	groupID     string
}

func NewNotificationConsumer(processor NotificationProcessor, brokers []string, topics []string, groupID string) *NotificationConsumer {
	return &NotificationConsumer{
		processor:   processor,
		kafkaBroker: brokers,
		topics:      topics,
		groupID:     groupID,
	}
}

func (c *NotificationConsumer) Start(ctx context.Context) error {
	for _, topic := range c.topics {
		go c.consumeTopic(ctx, topic)
	}
	<-ctx.Done()
	return nil
}

func (c *NotificationConsumer) consumeTopic(ctx context.Context, topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.kafkaBroker,
		Topic:       topic,
		GroupID:     c.groupID + "-" + topic, // отдельный groupID на каждый топик
		StartOffset: kafka.FirstOffset,       // читаем с начала
	})
	defer func() {
		if err := reader.Close(); err != nil {
			log.Println("Error closing Kafka reader:", err)
		}
	}()

	log.Printf("Consumer started for topic: %s", topic)

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Consumer context cancelled, stopping:", topic)
				return
			}
			log.Printf("Error reading message from %s: %v", topic, err)
			time.Sleep(time.Second)
			continue
		}

		log.Printf("Kafka message received from %s: key=%s offset=%d", topic, string(msg.Key), msg.Offset)

		switch topic {
		case "EventNotificationTopic":
			var n models.Notification
			if err := json.Unmarshal(msg.Value, &n); err != nil {
				log.Println("Failed to unmarshal Notification:", err)
				continue
			}
			if err := c.processor.HandleNotification(ctx, &n); err != nil {
				log.Println("Error processing notification:", err)
			}

		case "NotificationStatusTopic":
			var status models.NotificationStatus
			if err := json.Unmarshal(msg.Value, &status); err != nil {
				log.Println("Failed to unmarshal NotificationStatus:", err)
				continue
			}
			if err := c.processor.HandleNotificationStatus(ctx, &status); err != nil {
				log.Println("Error processing notification status:", err)
			}

		default:
			log.Printf("Unknown topic: %s", topic)
		}
	}
}