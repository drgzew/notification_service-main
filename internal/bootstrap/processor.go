package bootstrap

import (
	processors "notification_service/internal/services/processors"
	notifications "notification_service/internal/services/notifications"
)

func InitNotificationProcessor(notificationService notifications.NotificationServiceInterface) *processors.NotificationProcessor {
	return processors.NewNotificationProcessor(notificationService)
}
