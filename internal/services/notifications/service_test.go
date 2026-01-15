package notifications_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"notification_service/internal/models"
	"notification_service/internal/services/notifications"
	mock_notifications "notification_service/internal/services/notifications/mocks"
	mock_kafka "notification_service/internal/kafka/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func ptrTime(t time.Time) *time.Time {
	return &t
}

type NotificationServiceTestSuite struct {
	suite.Suite
	service  *notifications.NotificationService
	storage  *mock_notifications.NotificationStorage
	producer *mock_kafka.NotificationProducerInterface
	ctx      context.Context
}

func (s *NotificationServiceTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.storage = &mock_notifications.NotificationStorage{}
	s.producer = &mock_kafka.NotificationProducerInterface{}
	s.service = notifications.NewNotificationService(s.ctx, s.storage, s.producer, 10, 5*time.Second)
}

func (s *NotificationServiceTestSuite) TestSendNotification_Success() {
	now := time.Now()
	n := &models.Notification{
		ID:        "1",
		CreatedAt: &now,
		Message:   "Test message",
	}

	s.storage.On("AddNotification", s.ctx, n).Return(nil)
	s.storage.On("UpdateNotificationStatus", s.ctx, mock.Anything).Return(nil)
	s.producer.On("SendNotification", s.ctx, n).Return(nil)
	s.producer.On("SendNotificationStatus", s.ctx, mock.Anything).Return(nil)

	err := s.service.SendNotification(s.ctx, n)
	s.NoError(err)

	s.storage.AssertCalled(s.T(), "AddNotification", s.ctx, n)
	s.producer.AssertCalled(s.T(), "SendNotification", s.ctx, n)
	s.producer.AssertCalled(s.T(), "SendNotificationStatus", s.ctx, mock.Anything)
}

func (s *NotificationServiceTestSuite) TestSendNotification_ProducerFails() {
	now := time.Now()
	n := &models.Notification{
		ID:        "2",
		CreatedAt: &now,
		Message:   "Fail test",
	}

	s.storage.On("AddNotification", s.ctx, n).Return(nil)
	s.storage.On("UpdateNotificationStatus", s.ctx, mock.Anything).Return(nil)
	s.producer.On("SendNotification", s.ctx, n).Return(errors.New("kafka error"))
	s.producer.On("SendNotificationStatus", s.ctx, mock.Anything).Return(nil)

	err := s.service.SendNotification(s.ctx, n)
	s.NoError(err)

	s.producer.AssertCalled(s.T(), "SendNotification", s.ctx, n)
}

func TestNotificationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationServiceTestSuite))
}