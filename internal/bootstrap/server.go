package bootstrap

import (
	"context"
	"fmt"
	"log/slog"

	"notification_service/internal/api/service_api"
	"notification_service/internal/services/notifications"
	"notification_service/internal/kafka"
)

func AppRun(
	ctx context.Context,
	service *notifications.NotificationService,
	consumer *kafka.NotificationConsumer,
) {
	go func() {
		slog.Info("NotificationConsumer started")
		if err := consumer.Start(ctx); err != nil && err != context.Canceled {
			panic(fmt.Errorf("NotificationConsumer crashed: %v", err))
		}
	}()

	notifications_api.InitNotificationServiceAPI(service)
	slog.Info("Notification API started on :8080")

	<-ctx.Done()
	slog.Info("Shutdown signal received, stopping Notification Service")
}
