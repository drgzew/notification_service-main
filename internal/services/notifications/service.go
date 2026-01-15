package notifications

import (
	"context"
	"time"
	"log"

	"notification_service/internal/models"
	"notification_service/internal/kafka"
)

type NotificationStorage interface {
	UpdateNotificationStatus(ctx context.Context, statuses []*models.NotificationStatus) error
	GetNotificationStatusByIDs(ctx context.Context, IDs []string) ([]*models.NotificationStatus, error)
	AddNotification(ctx context.Context, n *models.Notification) error
}

type NotificationServiceInterface interface {
	Handle(ctx context.Context, n *models.Notification) error
}

type NotificationService struct {
	storage   NotificationStorage
	producer  kafka.NotificationProducerInterface
	batchSize int
	timeout   time.Duration
}

func NewNotificationService(
	ctx context.Context, storage NotificationStorage, producer kafka.NotificationProducerInterface, batchSize int, timeout time.Duration) *NotificationService {

	return &NotificationService{
		storage: storage, producer: producer, batchSize: batchSize, timeout: timeout,
	}
}

func (s *NotificationService) SendNotification(ctx context.Context, n *models.Notification) error {
    if err := s.storage.AddNotification(ctx, n); err != nil {
        return err
    }

    status := &models.NotificationStatus{
        NotificationID: n.ID,
        Status:         "PENDING",
        CreatedAt:      n.CreatedAt, // время создания уведомления
    }

    if err := s.storage.UpdateNotificationStatus(ctx, []*models.NotificationStatus{status}); err != nil {
        return err
    }

    err := s.producer.SendNotification(ctx, n)
    if err != nil {
        status.Status = "FAILED"
        status.ErrorMessage = err.Error()
        log.Printf("Failed to send notification ID=%s: %v", n.ID, err)
    } else {
        status.Status = "SENT"
        now := time.Now()
        status.SentAt = &now // Время фактической отправки
        log.Printf("Notification ID=%s sent successfully", n.ID)
    }

    if err := s.storage.UpdateNotificationStatus(ctx, []*models.NotificationStatus{status}); err != nil {
        return err
    }

    if err := s.producer.SendNotificationStatus(ctx, status); err != nil {
        log.Printf("Failed to send notification status to Kafka for ID=%s: %v", n.ID, err)
    }

    return nil
}

func (s *NotificationService) Handle(ctx context.Context, n *models.Notification) error {
	return s.SendNotification(ctx, n)
}