package integration

import (
	"context"
	"email-service/internal/processors"
	"email-service/internal/services"
	"email-service/tests/mock/mailer"
	"weather-forecast/pkg/proto/events"

	"testing"
	stub_logger "weather-forecast/pkg/stubs/logger"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type (
	ExpectedEmail struct {
		Subject string
		Body    string
		SentTo  string
	}
)

func setupEventProcessor() (*processors.EventProcessor, *mailer.MockSMTPMailer) {
	stubLogger := stub_logger.New()
	mockMailer := mailer.NewMockSMTPMailer()

	emailBuilder := services.NewSimpleEmailBuild("https://test.example.com", stubLogger)
	notificationService := services.NewNotificationService(mockMailer, emailBuilder, stubLogger)
	eventProcessor := processors.NewEventProcessor(notificationService, stubLogger)

	return eventProcessor, mockMailer
}

func assertEmailMatches(t *testing.T, sentEmail, expectedEmail mailer.SentEmail) {
	t.Helper()
	assert.Equal(t, expectedEmail.Subject, sentEmail.Subject)
	assert.Equal(t, expectedEmail.Body, sentEmail.Body)
	assert.Equal(t, expectedEmail.SentTo, sentEmail.SentTo)
}
func Test_SubscriptionEvent(t *testing.T) {
	eventProcessor, mockMailer := setupEventProcessor()

	event := &events.SubscriptionEvent{
		Email:     "test@example.com",
		Frequency: "daily",
		Token:     "abc123",
	}

	expected := mailer.SentEmail{
		Subject: "Confirm your subscription",
		Body:    "You have signed up for an daily newsletter. \nPlease, use this token to confirm your subscription: abc123\nOr use this link: https://test.example.com/confirm/abc123",
		SentTo:  "test@example.com",
	}

	eventBody, err := proto.Marshal(event)
	require.NoError(t, err)

	ctx := context.Background()

	eventProcessor.Handle(ctx, "emails.subscription", eventBody)

	emails := mockMailer.GetSentEmails()
	require.Len(t, emails, 1)

	assertEmailMatches(t, emails[0], expected)
}

func Test_ConfirmedEvent(t *testing.T) {
	eventProcessor, mockMailer := setupEventProcessor()

	event := &events.ConfirmedEvent{
		Email:     "test@example.com",
		Frequency: "daily",
		Token:     "abc123",
	}

	expected := mailer.SentEmail{
		Subject: "Subscription confirmed",
		Body:    "Congratulations, you have successfully confirmed your daily subscription.\nYou can cancel your subscription using this token: abc123\nOr use this link: https://test.example.com/unsubscribe/abc123",
		SentTo:  "test@example.com",
	}

	eventBody, err := proto.Marshal(event)
	require.NoError(t, err)

	ctx := context.Background()

	eventProcessor.Handle(ctx, "emails.confirmed", eventBody)

	emails := mockMailer.GetSentEmails()
	require.Len(t, emails, 1)

	assertEmailMatches(t, emails[0], expected)
}

func Test_UnsubscribedEvent(t *testing.T) {
	eventProcessor, mockMailer := setupEventProcessor()

	event := &events.UnsubscribedEvent{
		Email:     "test@example.com",
		Frequency: "daily",
		City:      "Kyiv",
	}

	expected := mailer.SentEmail{
		Subject: "Subscription canceled",
		Body:    "You have successfully canceled your daily subscription for city Kyiv.",
		SentTo:  "test@example.com",
	}

	eventBody, err := proto.Marshal(event)
	require.NoError(t, err)

	ctx := context.Background()

	eventProcessor.Handle(ctx, "emails.unsubscribed", eventBody)

	emails := mockMailer.GetSentEmails()
	require.Len(t, emails, 1)

	assertEmailMatches(t, emails[0], expected)
}

func Test_WeatherSuccessEvent(t *testing.T) {
	eventProcessor, mockMailer := setupEventProcessor()

	weather := &events.Weather{
		Temperature: 54,
		Humidity:    54,
		Description: "Sunny",
	}

	event := &events.WeatherSuccessEvent{
		Email:   "test@example.com",
		City:    "Kyiv",
		Weather: weather,
	}

	expected := mailer.SentEmail{
		Subject: "Weather Update",
		Body:    "Here's the latest weather update for your city: Kyiv\nTemperature: 54.0°C\nHumidity: 54%\nDescription: Sunny",
		SentTo:  "test@example.com",
	}

	eventBody, err := proto.Marshal(event)
	require.NoError(t, err)

	ctx := context.Background()
	eventProcessor.Handle(ctx, "emails.weather.success", eventBody)

	emails := mockMailer.GetSentEmails()
	require.Len(t, emails, 1)
	assertEmailMatches(t, emails[0], expected)
}

func Test_WeatherErrorEvent(t *testing.T) {
	eventProcessor, mockMailer := setupEventProcessor()

	event := &events.WeatherErrorEvent{
		Email: "test@example.com",
		City:  "Kyiv",
	}

	expected := mailer.SentEmail{
		Subject: "Weather Update",
		Body:    "Sorry, there was an error retrieving weather in your city: Kyiv",
		SentTo:  "test@example.com",
	}

	eventBody, err := proto.Marshal(event)
	require.NoError(t, err)

	ctx := context.Background()
	eventProcessor.Handle(ctx, "emails.weather.error", eventBody)

	emails := mockMailer.GetSentEmails()
	require.Len(t, emails, 1)
	assertEmailMatches(t, emails[0], expected)
}
