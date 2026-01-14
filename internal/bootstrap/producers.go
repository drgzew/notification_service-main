package bootstrap

import (
	"fmt"
	"notification_service/config"
	"notification_service/internal/kafka"
)

func InitNotificationProducer(cfg *config.Config) *kafka.NotificationProducer {
    brokers := []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)}
    fmt.Printf("[Kafka] Initializing EventNotification producer with brokers: %v\n", brokers)
    return kafka.NewNotificationProducer(brokers, cfg.Kafka.EventNotificationTopic)
}

func InitNotificationStatusProducer(cfg *config.Config) *kafka.NotificationProducer {
    brokers := []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)}
    fmt.Printf("[Kafka] Initializing NotificationStatus producer with brokers: %v\n", brokers)
    return kafka.NewNotificationProducer(brokers, cfg.Kafka.NotificationStatusTopic)
}
