package services

import (
	mock_services "weather-forecast/internal/infrastructure/services/mocks"
)

type (
	mockSubscribeBehavior func(
		uc *mock_services.MockSubscriptionUseCase,
		tokenManager *mock_services.MockTokenManager,
		notificationService *mock_services.MockNotificationServiceI,
		email,
		city,
		frequency string,
	)
)
