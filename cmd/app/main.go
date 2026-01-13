package main

import (
	"fmt"
	"os"
	"context"

	"notification_service/config"
	"notification_service/internal/bootstrap"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("Ошибка парсинга файла конфигурации: %v", err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	notificationStorage := bootstrap.InitPGStorage(cfg)
	notificationProducer := bootstrap.InitNotificationProducer(cfg)
	notificationService := bootstrap.InitNotificationService(ctx, notificationStorage, cfg, notificationProducer)
	notificationProcessor := bootstrap.InitNotificationProcessor(notificationService)
	notificationConsumer := bootstrap.InitNotificationConsumer(cfg, notificationProcessor)

	bootstrap.AppRun(ctx, notificationService, notificationConsumer)
}