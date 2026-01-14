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
	storage NotificationStorage
	producer *kafka.NotificationProducer
	batchSize int
	timeout time.Duration
}

func NewNotificationService(
	ctx context.Context, storage NotificationStorage, producer *kafka.NotificationProducer, batchSize int, timeout time.Duration) *NotificationService {
	return &NotificationService{
		storage: storage, producer: producer, batchSize: batchSize, timeout: timeout}
}

func (s *NotificationService) SendNotification(ctx context.Context, n *models.Notification) error {

	status := &models.NotificationStatus{
		NotificationID: n.ID,
		Status:         "PENDING",
	}
	if err := s.storage.UpdateNotificationStatus(ctx, []*models.NotificationStatus{status}); err != nil {
		return err
	}

	if err := s.storage.AddNotification(ctx, n); err != nil {
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
		status.SentAt = &now
		log.Printf("Notification ID=%s sent successfully", n.ID)
	}

	return s.storage.UpdateNotificationStatus(ctx, []*models.NotificationStatus{status})
}

func (s *NotificationService) SendBatch(ctx context.Context, notifications []*models.Notification) error {
	total := len(notifications)
	for i := 0; i < total; i += s.batchSize {
		end := i + s.batchSize
		if end > total {
			end = total
		}
		batch := notifications[i:end]

		batchStatuses := make([]*models.NotificationStatus, len(batch))
		for j, n := range batch {
			batchStatuses[j] = &models.NotificationStatus{
				NotificationID: n.ID,
				Status:         "PENDING",
			}
		}

		if err := s.storage.UpdateNotificationStatus(ctx, batchStatuses); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(s.timeout):
		}

		for j, n := range batch {
			err := s.producer.SendNotification(ctx, n)
			status := batchStatuses[j]

			if err != nil {
				status.Status = "FAILED"
				status.ErrorMessage = err.Error()
				log.Printf("Failed to send notification ID=%s: %v", n.ID, err)
			} else {
				status.Status = "SENT"
				now := time.Now()
				status.SentAt = &now
				log.Printf("Notification ID=%s sent successfully", n.ID)
			}
		}

		if err := s.storage.UpdateNotificationStatus(ctx, batchStatuses); err != nil {
			return err
		}
	}
	return nil
}

func (s *NotificationService) Handle(ctx context.Context, n *models.Notification) error {
    return s.SendNotification(ctx, n)
}