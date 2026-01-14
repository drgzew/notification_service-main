package notifications_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"notification_service/internal/models"
	"notification_service/internal/services/notifications"
	"notification_service/internal/services/notifications/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

//go:generate mockgen -destination=mocks/mock_notification_storage.go -package=mocks notification_service/internal/services/notifications NotificationStorage
//go:generate mockgen -destination=mocks/mock_notification_producer.go -package=mocks notification_service/internal/kafka NotificationProducerInterface

type NotificationServiceTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	mockStorage  *mocks.MockNotificationStorage
	mockProducer *mocks.MockNotificationProducerInterface
	service      *notifications.NotificationService
	ctx          context.Context
}

func (s *NotificationServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockStorage = mocks.NewMockNotificationStorage(s.ctrl)
	s.mockProducer = mocks.NewMockNotificationProducerInterface(s.ctrl)
	s.ctx = context.Background()

	s.service = notifications.NewNotificationService(s.ctx, s.mockStorage, s.mockProducer, 2, 10*time.Millisecond)
}

func (s *NotificationServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *NotificationServiceTestSuite) TestSendNotification_Success() {
	n := &models.Notification{ID: "1", Message: "Hello"}

	s.mockStorage.EXPECT().UpdateNotificationStatus(s.ctx, gomock.Any()).Return(nil)
	s.mockStorage.EXPECT().AddNotification(s.ctx, n).Return(nil)
	s.mockProducer.EXPECT().SendNotification(s.ctx, n).Return(nil)
	s.mockStorage.EXPECT().UpdateNotificationStatus(s.ctx, gomock.Any()).Return(nil)

	err := s.service.SendNotification(s.ctx, n)
	s.NoError(err)
}

func (s *NotificationServiceTestSuite) TestSendNotification_ProducerFails() {
	n := &models.Notification{ID: "1", Message: "Hello"}

	s.mockStorage.EXPECT().UpdateNotificationStatus(s.ctx, gomock.Any()).Return(nil)
	s.mockStorage.EXPECT().AddNotification(s.ctx, n).Return(nil)
	s.mockProducer.EXPECT().SendNotification(s.ctx, n).Return(errors.New("kafka error"))
	s.mockStorage.EXPECT().UpdateNotificationStatus(s.ctx, gomock.Any()).Return(nil)

	err := s.service.SendNotification(s.ctx, n)
	s.NoError(err)
}

func (s *NotificationServiceTestSuite) TestSendBatch_Success() {
	n1 := &models.Notification{ID: "1", Message: "Hello1"}
	n2 := &models.Notification{ID: "2", Message: "Hello2"}
	notifications := []*models.Notification{n1, n2}

	s.mockStorage.EXPECT().UpdateNotificationStatus(s.ctx, gomock.Any()).Return(nil).Times(2)
	s.mockProducer.EXPECT().SendNotification(s.ctx, n1).Return(nil)
	s.mockProducer.EXPECT().SendNotification(s.ctx, n2).Return(nil)

	err := s.service.SendBatch(s.ctx, notifications)
	s.NoError(err)
}

func (s *NotificationServiceTestSuite) TestSendBatch_PartialFailure() {
	n1 := &models.Notification{ID: "1", Message: "Hello1"}
	n2 := &models.Notification{ID: "2", Message: "Hello2"}
	notifications := []*models.Notification{n1, n2}

	s.mockStorage.EXPECT().UpdateNotificationStatus(s.ctx, gomock.Any()).Return(nil).Times(2)
	s.mockProducer.EXPECT().SendNotification(s.ctx, n1).Return(errors.New("fail1"))
	s.mockProducer.EXPECT().SendNotification(s.ctx, n2).Return(nil)

	err := s.service.SendBatch(s.ctx, notifications)
	s.NoError(err)
}

func TestNotificationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationServiceTestSuite))
}