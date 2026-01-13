package notifications

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"notification_service/internal/models"
	mocks "notification_service/internal/services/notifications/mocks"


)

type NotificationServiceTestSuite struct {
	suite.Suite
	service   *NotificationService
	mockStore *mocks.NotificationStorage
}

func (suite *NotificationServiceTestSuite) SetupTest() {
	suite.mockStore = &mocks.NotificationStorage{}
	suite.service = &NotificationService{
		notificationStorage: suite.mockStore,
		batchSize:           2,
		timeout:             10 * time.Millisecond,
	}
}

func (suite *NotificationServiceTestSuite) TestSendNotification_Success() {
	ctx := context.Background()
	n := &models.Notification{ID: "1", Recipient: "test@example.com", Message: "Hello"}

	suite.mockStore.On("UpdateNotificationStatus", ctx, mock.Anything).Return(nil).Twice()

	err := suite.service.SendNotification(ctx, n)
	assert.NoError(suite.T(), err)

	suite.mockStore.AssertNumberOfCalls(suite.T(), "UpdateNotificationStatus", 2)
}

func (suite *NotificationServiceTestSuite) TestSendBatch_Success() {
	ctx := context.Background()
	notifs := []*models.Notification{
		{ID: "1", Recipient: "a", Message: "m1"},
		{ID: "2", Recipient: "b", Message: "m2"},
		{ID: "3", Recipient: "c", Message: "m3"},
	}

	suite.mockStore.On("UpdateNotificationStatus", ctx, mock.Anything).Return(nil).Times(4)

	err := suite.service.SendBatch(ctx, notifs)
	assert.NoError(suite.T(), err)
	suite.mockStore.AssertExpectations(suite.T())
}

func (suite *NotificationServiceTestSuite) TestGetStatus() {
	ctx := context.Background()
	status := &models.NotificationStatus{NotificationID: "1", Status: "SENT"}
	suite.mockStore.On("GetNotificationStatusByIDs", ctx, []string{"1"}).Return([]*models.NotificationStatus{status}, nil)

	got, err := suite.service.GetStatus(ctx, "1")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "SENT", got.Status)
}

func (suite *NotificationServiceTestSuite) TestGetStatuses() {
	ctx := context.Background()
	statuses := []*models.NotificationStatus{
		{NotificationID: "1", Status: "SENT"},
		{NotificationID: "2", Status: "FAILED"},
	}
	suite.mockStore.On("GetNotificationStatusByIDs", ctx, []string{"1", "2"}).Return(statuses, nil)

	got, err := suite.service.GetStatuses(ctx, []string{"1", "2"})
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), got, 2)
}

func TestNotificationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationServiceTestSuite))
}