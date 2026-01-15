package bootstrap

import (
	"fmt"
	"notification_service/config"

	kafka "notification_service/internal/kafka"
	processors "notification_service/internal/services/processors"
)


func InitNotificationConsumer(cfg *config.Config, processor *processors.NotificationProcessor) *kafka.NotificationConsumer {
	brokers := []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)}
	return kafka.NewNotificationConsumer(processor, brokers, []string{cfg.Kafka.EventNotificationTopic, cfg.Kafka.NotificationStatusTopic}, "notification-service")
}
