package main

import (
	"context"
	"flag"
	"log"

	"notification_service/config"
	api "notification_service/internal/api/service_api"
	"notification_service/internal/bootstrap"
)

func main() {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()
	if *configPath == "" {
		log.Fatal("config path is required: use --config")
	}
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	notificationStorage := bootstrap.InitPGStorage(cfg)
	notificationProducer := bootstrap.InitNotificationProducer(cfg)
	notificationService := bootstrap.InitNotificationService(ctx, notificationStorage, notificationProducer)
	notificationProcessor := bootstrap.InitNotificationProcessor(notificationService)
	notificationConsumer := bootstrap.InitNotificationConsumer(cfg, notificationProcessor)

	go func() {
		if err := notificationConsumer.Start(ctx); err != nil {
			log.Fatalf("NotificationConsumer crashed: %v", err)
		}
	}()

	notificationAPI := api.InitNotificationServiceAPI(notificationService)

	log.Println("Notification API started on :8080")
	if err := notificationAPI.Run(":8080"); err != nil {
		log.Fatalf("Failed to start API: %v", err)
	}
}