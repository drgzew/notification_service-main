package notifications_api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"notification_service/internal/models"
	"notification_service/internal/services/notifications"
)

func InitNotificationServiceAPI(service *notifications.NotificationService) *gin.Engine {
	r := gin.Default()

	r.POST("/notifications", func(c *gin.Context) {
		var req struct {
			Recipient string `json:"recipient"`
			Channel   string `json:"channel"`
			Message   string `json:"message"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		notification := &models.Notification{
			ID:        uuid.NewString(),
			Recipient: req.Recipient,
			Channel:   req.Channel,
			Message:   req.Message,
		}
		now := time.Now()
		notification.CreatedAt = &now

		if err := service.SendNotification(c.Request.Context(), notification); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok", "id": notification.ID})
	})

	go func() {
		if err := r.Run(":8080"); err != nil {
			panic("Ошибка запуска Notification API: " + err.Error())
		}
	}()

	return r
}