package bootstrap

import (
	"context"
	"fmt"
	"log"
	
	"notification_service/internal/api/service_api"
	"notification_service/internal/kafka"
	"notification_service/internal/services/notifications"
)

func AppRun(ctx context.Context, service *notifications.NotificationService, consumer *kafka.NotificationConsumer) {

	go func() {
		log.Println("NotificationConsumer started")
		if err := consumer.Start(ctx); err != nil && err != context.Canceled {
			panic(fmt.Errorf("NotificationConsumer crashed: %v", err))
		}
	}()

	r := notifications_api.InitNotificationServiceAPI(service)
	log.Println("Notification API started on :8080")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start API: %v", err)
	}

	<-ctx.Done()
	log.Println("Shutdown signal received, stopping Notification Service")
}
