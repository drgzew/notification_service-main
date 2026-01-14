package main

import (
	"context"
	"log"
	"flag"
	"fmt"

	"notification_service/config"
	api "notification_service/internal/api/service_api"
	"notification_service/internal/bootstrap"

	swaggerFiles "github.com/swaggo/files"
)

func main() {
    configPath := flag.String("config", "", "path to config file")
    flag.Parse()

    if *configPath == "" {
        log.Fatal("config path is required: use --config")
    }

	ctx := context.Background()

    cfg, err := config.LoadConfig(*configPath)
    if err != nil {
        log.Fatal(err)
    }

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

