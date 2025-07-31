package integration

import (
	"context"
	"subscription-service/internal/domain/models"
	"subscription-service/tests/mocks/publisher"
	"testing"
	"time"
	protoevents "weather-forecast/pkg/proto/events"
	"weather-forecast/pkg/proto/subscription"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func assertUnsubscribedEventPublished(t *testing.T, publisher *publisher.MockEventPublisher, expectedEmail, expectedCity string, expectedFrequency models.Frequency) {
	t.Helper()

	eventList := publisher.GetPublishedEvents()
	require.Len(t, eventList, 1)
	lastEvent := eventList[0]
	assert.Equal(t, "emails.unsubscribed", lastEvent.EventType)

	var unsubscribedEvent protoevents.UnsubscribedEvent
	err := proto.Unmarshal(lastEvent.RawData, &unsubscribedEvent)
	require.NoError(t, err)

	assert.Equal(t, expectedEmail, unsubscribedEvent.Email)
	assert.Equal(t, string(expectedFrequency), unsubscribedEvent.Frequency)
	assert.Equal(t, expectedCity, unsubscribedEvent.City)
}

func TestUnsubscribe_Success(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := setupDB(t)
	subscriptionHandler, mockPublisher := setupHandler(db)
	unsubscribeSubscription := models.Subscription{
		Email:     "test@gmail.com",
		City:      "Kyiv",
		Frequency: models.Daily,
		Confirmed: true,
		Token:     "ce5bc383-2820-4358-a2af-c038382e617b",
	}
	err := db.Create(&unsubscribeSubscription).Error
	require.NoError(t, err)

	requestBody := &subscription.UnsubscribeRequest{
		Token: unsubscribeSubscription.Token,
	}

	resp, err := subscriptionHandler.Unsubscribe(ctx, requestBody)
	require.NoError(t, err)
	assert.IsType(t, &emptypb.Empty{}, resp)

	res := db.Where("id = ?", unsubscribeSubscription.ID).Find(&models.Subscription{})
	require.NoError(t, res.Error)
	require.Equal(t, int64(0), res.RowsAffected)

	assertUnsubscribedEventPublished(t, mockPublisher, unsubscribeSubscription.Email, unsubscribeSubscription.City, unsubscribeSubscription.Frequency)

}

func TestUnsubscribe_InvalidToken(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := setupDB(t)
	subscriptionHandler, _ := setupHandler(db)

	requestBody := &subscription.UnsubscribeRequest{
		Token: "invalidToken",
	}

	resp, err := subscriptionHandler.Unsubscribe(ctx, requestBody)
	require.Error(t, err)
	assert.IsType(t, &emptypb.Empty{}, resp)
	grpcStatus, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, grpcStatus.Code())
	assert.Contains(t, grpcStatus.Message(), "invalid token")

}
