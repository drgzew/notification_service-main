package bootstrap

import (
	"context"
	"time"

	"notification_service/internal/services/notifications"
	"notification_service/internal/storage/pgstorage"
	"notification_service/internal/kafka"
)
func InitNotificationService(ctx context.Context, storage *pgstorage.PGStorage, producer *kafka.NotificationProducer) *notifications.NotificationService {
	return notifications.NewNotificationService(ctx, storage, producer, 10, 1*time.Second)
}