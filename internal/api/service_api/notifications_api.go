package notifications_api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"notification_service/internal/models"
	"notification_service/internal/services/notifications"
)

// InitNotificationServiceAPI создает роуты для NotificationService
func InitNotificationServiceAPI(service *notifications.NotificationService) *gin.Engine {
	r := gin.Default()

	// POST /notifications — создать новое уведомление
	// @Summary Create notification
	// @Description Sends a notification via NotificationService
	// @Tags notifications
	// @Accept  json
	// @Produce  json
	// @Param notification body models.Notification true "Notification payload"
	// @Success 201 {object} map[string]string
	// @Failure 400 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /notifications [post]
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

		c.JSON(http.StatusCreated, gin.H{"status": "ok", "id": notification.ID})
	})  

	return r
}

func WrapSwaggerHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}