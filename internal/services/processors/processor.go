package processors

import (
	"context"
	"fmt"
	"notification_service/internal/models"
	"notification_service/internal/services/notifications"
)

type NotificationProcessor struct {
	service notifications.NotificationServiceInterface
}

func NewNotificationProcessor(s notifications.NotificationServiceInterface) *NotificationProcessor {
	return &NotificationProcessor{service: s}
}

func (p *NotificationProcessor) HandleNotification(ctx context.Context, notification *models.Notification) error {
	fmt.Printf("Обрабатываем уведомление ID=%s для %s: %s\n",
		notification.ID, notification.Recipient, notification.Message)

	if p.service != nil {
		return p.service.Handle(ctx, notification)
	}

	return nil
}

func (p *NotificationProcessor) HandleNotificationStatus(ctx context.Context, status *models.NotificationStatus) error {
	fmt.Printf("Обрабатываем статус уведомления ID=%s: %s\n",
		status.NotificationID, status.Status)
	return nil // пока просто логируем
}
