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
	router := gin.Default()

	// POST /notifications — создать новое уведомление
	router.POST("/notifications", func(c *gin.Context) {
		var req struct {
			Recipient string `json:"recipient" binding:"required"`
			Message   string `json:"message" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request payload: " + err.Error(),
			})
			return
		}

		now := time.Now()
		notification := &models.Notification{
			ID:        uuid.NewString(),
			Recipient: req.Recipient,
			Message:   req.Message,
			CreatedAt: &now,
		}

		if err := service.SendNotification(c.Request.Context(), notification); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to send notification: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status": "ok",
			"id":     notification.ID,
		})
	})

	return router
}

func WrapSwaggerHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}