package notifications_api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "notification_service/docs" 
	files "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"notification_service/internal/models"
	"notification_service/internal/services/notifications"
)

type notificationsRequest struct {
	Recipient string `json:"recipient" binding:"required"`
	Message   string `json:"message" binding:"required"`
}

type notificationsResponse struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func InitNotificationServiceAPI(service *notifications.NotificationService) *gin.Engine {
	router := gin.Default()

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	// POST /notifications
	router.POST("/notifications", func(c *gin.Context) {
		createNotificationHandler(c, service)
	})

	// Для проверки сервера
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return router
}

// createNotificationHandler 
// @Summary Создать уведомление
// @Description Создает новое уведомление для указанного получателя
// @Tags notifications
// @Accept json
// @Produce json
// @Param notification body notificationsRequest true "Notification payload"
// @Success 201 {object} notificationsResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /notifications [post]
func createNotificationHandler(c *gin.Context, service *notifications.NotificationService) {
	var req notificationsRequest
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

	c.JSON(http.StatusCreated, notificationsResponse{
		Status: "ok",
		ID:     notification.ID,
	})
}