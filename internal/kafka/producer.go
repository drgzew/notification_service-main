package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"notification_service/internal/models"
	"notification_service/config"
)

type NotificationProducerInterface interface {
    SendNotification(ctx context.Context, n *models.Notification) error
    SendNotificationStatus(ctx context.Context, status *models.NotificationStatus) error
}

type NotificationProducer struct {
    writerNotification *kafka.Writer
    writerStatus       *kafka.Writer
    cfg                *config.KafkaConfig
}

func NewNotificationProducer(brokers []string, cfg *config.KafkaConfig) *NotificationProducer {
    return &NotificationProducer{
        writerNotification: &kafka.Writer{
            Addr:  kafka.TCP(brokers...),
            Topic: cfg.EventNotificationTopic,
            Balancer: &kafka.LeastBytes{},
			Async: false,
        },
        writerStatus: &kafka.Writer{
            Addr:  kafka.TCP(brokers...),
            Topic: cfg.NotificationStatusTopic,
            Balancer: &kafka.LeastBytes{},
			Async: false,
        },
        cfg: cfg,
    }
}

func (p *NotificationProducer) SendNotification(ctx context.Context, n *models.Notification) error {
    data, _ := json.Marshal(n)
    return p.writerNotification.WriteMessages(ctx, kafka.Message{
        Key:   []byte(n.ID),
        Value: data,
    })
}

func (p *NotificationProducer) SendNotificationStatus(ctx context.Context, status *models.NotificationStatus) error {
    data, _ := json.Marshal(status)
    msg := kafka.Message{
        Key:   []byte(status.NotificationID),
        Value: data,
    }

    if err := p.writerStatus.WriteMessages(ctx, msg); err != nil {
        return err
    }
	
	return nil
}

func (p *NotificationProducer) Close() error {
    if err := p.writerNotification.Close(); err != nil {
        return err
    }
    if err := p.writerStatus.Close(); err != nil {
        return err
    }
    log.Printf("Kafka writers for topics %s and %s closed\n", p.cfg.EventNotificationTopic, p.cfg.NotificationStatusTopic)
    return nil
}