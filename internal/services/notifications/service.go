package notifications

import (
	"context"
	"time"

	"notification_service/internal/models"
	notificationproducer "notification_service/internal/kafka"
)

type NotificationStorage interface {
	UpdateNotificationStatus(ctx context.Context, statuses []*models.NotificationStatus) error
	GetNotificationStatusByIDs(ctx context.Context, IDs []string) ([]*models.NotificationStatus, error)
}

type NotificationServiceInterface interface {
	Handle(ctx context.Context, n *models.Notification) error
}

type NotificationService struct {
	notificationStorage NotificationStorage
	producer            *notificationproducer.NotificationProducer
	batchSize           int
	timeout             time.Duration
}

func NewNotificationService(ctx context.Context, notificationStorage NotificationStorage, batchSize int, timeout time.Duration, producer *notificationproducer.NotificationProducer) *NotificationService {
	return &NotificationService{
		notificationStorage: notificationStorage,
		producer:            producer,
		batchSize:           batchSize,
		timeout:             timeout,
	}
}

func (s *NotificationService) Handle(ctx context.Context, n *models.Notification) error {
	return s.SendNotification(ctx, n)
}

func (s *NotificationService) SendNotification(ctx context.Context, n *models.Notification) error {
	status := &models.NotificationStatus{
		NotificationID: n.ID,
		Status:         "PENDING",
		SentAt:         nil,
	}

	if err := s.notificationStorage.UpdateNotificationStatus(ctx, []*models.NotificationStatus{status}); err != nil {
		return err
	}

	err := s.mockSend(n)

	if err != nil {
		status.Status = "FAILED"
		status.ErrorMessage = err.Error()
	} else {
		status.Status = "SENT"
		now := time.Now()
		status.SentAt = &now
	}
	return s.notificationStorage.UpdateNotificationStatus(ctx, []*models.NotificationStatus{status})
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
				SentAt:         nil,
			}
		}

		if err := s.notificationStorage.UpdateNotificationStatus(ctx, batchStatuses); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(s.timeout):
		}

		for j, n := range batch {
			err := s.mockSend(n)
			status := batchStatuses[j]
			if err != nil {
				status.Status = "FAILED"
				status.ErrorMessage = err.Error()
			} else {
				status.Status = "SENT"
				now := time.Now()
				status.SentAt = &now
			}
		}

		if err := s.notificationStorage.UpdateNotificationStatus(ctx, batchStatuses); err != nil {
			return err
		}
	}
	return nil
}

func (s *NotificationService) GetStatus(ctx context.Context, notificationID string) (*models.NotificationStatus, error) {
	statuses, err := s.notificationStorage.GetNotificationStatusByIDs(ctx, []string{notificationID})
	if err != nil {
		return nil, err
	}
	if len(statuses) == 0 {
		return nil, nil
	}
	return statuses[0], nil
}

func (s *NotificationService) GetStatuses(ctx context.Context, ids []string) ([]*models.NotificationStatus, error) {
	return s.notificationStorage.GetNotificationStatusByIDs(ctx, ids)
}

func (s *NotificationService) mockSend(n *models.Notification) error {
	return nil
}
