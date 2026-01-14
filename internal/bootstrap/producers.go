package bootstrap

import (
	"fmt"
	"notification_service/config"
	"notification_service/internal/kafka"
)

func InitNotificationProducer(cfg *config.Config) *kafka.NotificationProducer {
	brokers := []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)}
	return kafka.NewNotificationProducer(brokers, cfg.Kafka.EventNotificationTopic)
}
