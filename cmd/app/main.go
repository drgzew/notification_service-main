package main

import (
	"context"
	"fmt"
	"os"

	"notification_service/config"
	api "notification_service/internal/api/service_api"
	"notification_service/internal/bootstrap"

	swaggerFiles "github.com/swaggo/files"
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
	notificationAPI := api.InitNotificationServiceAPI(notificationService)

	// Swagger UI
	notificationAPI.GET("/swagger/*any", api.WrapSwaggerHandler(swaggerFiles.Handler))

	go func() {
		fmt.Println("HTTP server running at :8080")
		if err := notificationAPI.Run(":8080"); err != nil {
			panic("Ошибка запуска HTTP сервера: " + err.Error())
		}
	}()

	bootstrap.AppRun(ctx, notificationService, notificationConsumer)
}

