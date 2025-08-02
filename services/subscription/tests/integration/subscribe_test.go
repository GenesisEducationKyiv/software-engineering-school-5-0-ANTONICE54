package integration

import (
	"context"
	"subscription-service/internal/domain/models"
	"subscription-service/internal/domain/usecases"
	"subscription-service/internal/infrastructure/database"
	"subscription-service/internal/infrastructure/repositories"
	"subscription-service/internal/infrastructure/sender"
	"subscription-service/internal/infrastructure/token"
	"subscription-service/internal/presentation/mappers"
	"subscription-service/internal/presentation/server/handlers"
	"subscription-service/tests/mocks/publisher"
	protoevents "weather-forecast/pkg/proto/events"

	"testing"
	"time"

	"weather-forecast/pkg/proto/subscription"
	stub_logger "weather-forecast/pkg/stubs/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	require.NoError(t, err)

	database.RunMigration(db)

	return db
}

func setupHandler(db *gorm.DB) (*handlers.SubscriptionHandler, *publisher.MockEventPublisher) {

	stubLogger := stub_logger.New()
	tokenManager := token.NewUUIDManager()

	publisher := publisher.NewMockEventPublisher()
	sender := sender.NewEventSender(publisher, stubLogger)
	subscRepo := repositories.NewSubscriptionRepository(db, stubLogger)
	subscUC := usecases.NewSubscriptionService(subscRepo, tokenManager, sender, stubLogger)
	subscHandler := handlers.NewSubscriptionHandler(subscUC, stubLogger)

	return subscHandler, publisher
}

func assertSubscriptionEventPublished(t *testing.T, publisher *publisher.MockEventPublisher, expectedEmail string, expectedFrequency models.Frequency) {
	t.Helper()

	eventList := publisher.GetPublishedEvents()
	require.Len(t, eventList, 1)
	lastEvent := eventList[0]
	assert.Equal(t, "emails.subscription", lastEvent.EventType)

	var subscriptionEvent protoevents.SubscriptionEvent
	err := proto.Unmarshal(lastEvent.RawData, &subscriptionEvent)
	require.NoError(t, err)

	assert.Equal(t, expectedEmail, subscriptionEvent.Email)
	assert.Equal(t, string(expectedFrequency), subscriptionEvent.Frequency)
	assert.NotEmpty(t, subscriptionEvent.Token)
}

func TestSubscribe_Success(t *testing.T) {
	db := setupDB(t)

	subscriptionHandler, mockPublisher := setupHandler(db)

	requestBody := &subscription.SubscribeRequest{
		Email:     "test@gmail.com",
		City:      "Kyiv",
		Frequency: subscription.Frequency_DAILY,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := subscriptionHandler.Subscribe(ctx, requestBody)
	require.NoError(t, err)
	assert.IsType(t, &emptypb.Empty{}, resp)

	subscFromDB := models.Subscription{}
	err = db.Where("email = ?", requestBody.Email).First(&subscFromDB).Error
	require.NoError(t, err)
	assert.Equal(t, requestBody.City, subscFromDB.City)
	assert.Equal(t, mappers.ProtoToFrequency(requestBody.Frequency), subscFromDB.Frequency)
	assert.False(t, subscFromDB.Confirmed)
	assert.Equal(t, requestBody.Email, subscFromDB.Email)
	assertSubscriptionEventPublished(t, mockPublisher, requestBody.Email, mappers.ProtoToFrequency(requestBody.Frequency))
}

func TestSubscribe_AlreadySubscribed(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := setupDB(t)
	subscriptionHandler, mockPublisher := setupHandler(db)

	requestBody := &subscription.SubscribeRequest{
		Email:     "test@gmail.com",
		City:      "Kyiv",
		Frequency: subscription.Frequency_DAILY,
	}

	resp, err := subscriptionHandler.Subscribe(ctx, requestBody)
	require.NoError(t, err)
	assert.IsType(t, &emptypb.Empty{}, resp)
	assertSubscriptionEventPublished(t, mockPublisher, requestBody.Email, mappers.ProtoToFrequency(requestBody.Frequency))

	resp, err = subscriptionHandler.Subscribe(ctx, requestBody)

	assert.IsType(t, &emptypb.Empty{}, resp)
	grpcStatus, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.AlreadyExists, grpcStatus.Code())
	assert.Contains(t, grpcStatus.Message(), "email already subscribed")

}
