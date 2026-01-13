package bootstrap

import (
	"context"

	"notification_service/config"
	"notification_service/internal/services/notifications"
	"notification_service/internal/storage/pgstorage"
	notificationproducer "notification_service/internal/kafka"
)

func InitNotificationService(ctx context.Context, storage *pgstorage.PGStorage, cfg *config.Config, producer *notificationproducer.NotificationProducer) *notifications.NotificationService {
    return notifications.NewNotificationService(ctx, storage, cfg.NotificationServiceSettings.NotificationBatchSize, cfg.NotificationServiceSettings.NotificationTimeout, producer)
}
