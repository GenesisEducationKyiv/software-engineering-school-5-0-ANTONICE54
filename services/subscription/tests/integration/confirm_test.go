package integration

import (
	"context"
	"encoding/json"
	"subscription-service/internal/domain/models"
	"subscription-service/tests/mocks/publisher"
	"testing"
	"time"
	"weather-forecast/pkg/events"
	"weather-forecast/pkg/proto/subscription"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func assertConfirmedEventPublished(t *testing.T, publisher *publisher.MockEventPublisher, expectedEmail, expectedToken string, expectedFrequency models.Frequency) {
	t.Helper()

	eventList := publisher.GetPublishedEvents()
	require.Len(t, eventList, 1)

	lastEvent := eventList[0]
	assert.Equal(t, events.ConfirmedEmail, lastEvent.EventType)

	var confirmedEvent events.ConfirmedEvent
	err := json.Unmarshal(lastEvent.RawData, &confirmedEvent)
	require.NoError(t, err)

	assert.Equal(t, expectedEmail, confirmedEvent.Email)
	assert.Equal(t, string(expectedFrequency), confirmedEvent.Frequency)
	assert.Equal(t, expectedToken, confirmedEvent.Token)
}

func TestConfirm_Success(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := setupDB(t)
	subscriptionHandler, mockPublisher := setupHandler(db)
	unconfirmedSubscription := models.Subscription{
		Email:     "test@gmail.com",
		City:      "Kyiv",
		Frequency: models.Daily,
		Confirmed: false,
		Token:     "ce5bc383-2820-4358-a2af-c038382e617b",
	}
	err := db.Create(&unconfirmedSubscription).Error
	require.NoError(t, err)

	requestBody := &subscription.ConfirmRequest{
		Token: unconfirmedSubscription.Token,
	}

	resp, err := subscriptionHandler.Confirm(ctx, requestBody)
	require.NoError(t, err)
	assert.IsType(t, &emptypb.Empty{}, resp)

	var confirmedSubscription models.Subscription

	err = db.Where("id = ?", unconfirmedSubscription.ID).First(&confirmedSubscription).Error
	require.NoError(t, err)
	assert.True(t, confirmedSubscription.Confirmed)

	assertConfirmedEventPublished(t, mockPublisher, confirmedSubscription.Email, unconfirmedSubscription.Token, unconfirmedSubscription.Frequency)

}

func TestConfirm_InvalidToken(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := setupDB(t)
	subscriptionHandler, _ := setupHandler(db)

	requestBody := &subscription.ConfirmRequest{
		Token: "invalidToken",
	}

	resp, err := subscriptionHandler.Confirm(ctx, requestBody)
	require.Error(t, err)
	assert.IsType(t, &emptypb.Empty{}, resp)
	grpcStatus, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, grpcStatus.Code())
	assert.Contains(t, grpcStatus.Message(), "invalid token")

}
