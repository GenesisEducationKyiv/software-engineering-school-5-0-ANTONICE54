package services

import (
	"context"
	"testing"
	"time"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/apperrors"
	mock_logger "weather-forecast/internal/infrastructure/logger/mocks"
	mock_services "weather-forecast/internal/infrastructure/services/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type (
	mockSubscribeBehavior func(
		uc *mock_services.MockSubscriptionUseCase,
		tokenManager *mock_services.MockTokenManager,
		notificationService *mock_services.MockNotificationServiceI,
		email,
		city,
		frequency,
		token string,
	)
	mockConfirmBehavior func(
		uc *mock_services.MockSubscriptionUseCase,
		tokenManager *mock_services.MockTokenManager,
		notificationService *mock_services.MockNotificationServiceI,
		token string,
	)
)

func TestSubsctiptionService_Subscribe(t *testing.T) {

	subscriptionTemplate := models.Subscription{
		ID:        1,
		Email:     "test@test.com",
		Frequency: models.Hourly,
		City:      "Kyiv",
		Token:     "59d29860-39fa-4c9b-845a-3e91eab42e4b",
		Confirmed: false,
		CreatedAt: time.Now(),
	}

	testTable := []struct {
		name           string
		token          string
		email          string
		city           string
		frequency      string
		mockBehavior   mockSubscribeBehavior
		expectedError  error
		expectedResult *models.Subscription
	}{
		{
			name:      "Successful",
			token:     subscriptionTemplate.Token,
			email:     subscriptionTemplate.Email,
			city:      subscriptionTemplate.City,
			frequency: string(subscriptionTemplate.Frequency),
			mockBehavior: func(
				uc *mock_services.MockSubscriptionUseCase,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationServiceI,
				email,
				city,
				frequency,
				token string) {
				tokenManager.EXPECT().Generate(gomock.Any()).Return(token)
				passedSubscription := models.Subscription{
					Email:     email,
					City:      city,
					Frequency: models.Frequency(frequency),
					Token:     token,
				}

				uc.EXPECT().Subscribe(gomock.Any(), passedSubscription).Return(&subscriptionTemplate, nil)
				notificationService.EXPECT().SendConfirmation(gomock.Any(), passedSubscription.Email, passedSubscription.Token, passedSubscription.Frequency)

			},
			expectedResult: &subscriptionTemplate,
			expectedError:  nil,
		},
		{
			name:      "Already subscribed",
			token:     subscriptionTemplate.Token,
			email:     subscriptionTemplate.Email,
			city:      subscriptionTemplate.City,
			frequency: string(subscriptionTemplate.Frequency),
			mockBehavior: func(
				uc *mock_services.MockSubscriptionUseCase,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationServiceI,
				email,
				city,
				frequency,
				token string) {
				tokenManager.EXPECT().Generate(gomock.Any()).Return(token)
				passedSubscription := models.Subscription{
					Email:     email,
					City:      city,
					Frequency: models.Frequency(frequency),
					Token:     token,
				}

				uc.EXPECT().Subscribe(gomock.Any(), passedSubscription).Return(nil, apperrors.AlreadySubscribedError)

			},
			expectedResult: nil,
			expectedError:  apperrors.AlreadySubscribedError,
		},
		{
			name:      "Database error",
			token:     subscriptionTemplate.Token,
			email:     subscriptionTemplate.Email,
			city:      subscriptionTemplate.City,
			frequency: string(subscriptionTemplate.Frequency),
			mockBehavior: func(
				uc *mock_services.MockSubscriptionUseCase,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationServiceI,
				email,
				city,
				frequency,
				token string) {
				tokenManager.EXPECT().Generate(gomock.Any()).Return(token)
				passedSubscription := models.Subscription{
					Email:     email,
					City:      city,
					Frequency: models.Frequency(frequency),
					Token:     token,
				}

				uc.EXPECT().Subscribe(gomock.Any(), passedSubscription).Return(nil, apperrors.DatabaseError)

			},
			expectedResult: nil,
			expectedError:  apperrors.DatabaseError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			loggerMock := mock_logger.NewMockLogger(ctrl)
			tokenManagerMock := mock_services.NewMockTokenManager(ctrl)
			notificationServiceMock := mock_services.NewMockNotificationServiceI(ctrl)
			subscriptionUCMock := mock_services.NewMockSubscriptionUseCase(ctrl)

			testCase.mockBehavior(subscriptionUCMock, tokenManagerMock, notificationServiceMock, testCase.email, testCase.city, testCase.frequency, testCase.token)

			subscriptionService := NewSubscriptionService(subscriptionUCMock, tokenManagerMock, notificationServiceMock, loggerMock)
			res, err := subscriptionService.Subscribe(context.Background(), testCase.email, testCase.frequency, testCase.city)

			assert.Equal(t, testCase.expectedResult, res)
			assert.Equal(t, testCase.expectedError, err)

		})
	}

}

func TestSubscriptionService_Confirm(t *testing.T) {

	subscriptionTemplate := models.Subscription{
		ID:        1,
		Email:     "test@test.com",
		Frequency: models.Hourly,
		City:      "Kyiv",
		Token:     "59d29860-39fa-4c9b-845a-3e91eab42e4b",
		Confirmed: true,
	}

	testTable := []struct {
		name          string
		token         string
		mockBehavior  mockConfirmBehavior
		expectedError error
	}{
		{
			name:  "Successful",
			token: "59d29860-39fa-4c9b-845a-3e91eab42e4b",
			mockBehavior: func(
				uc *mock_services.MockSubscriptionUseCase,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationServiceI,
				token string,
			) {
				tokenManager.EXPECT().Validate(gomock.Any(), token).Return(true)
				uc.EXPECT().Confirm(gomock.Any(), token).Return(&subscriptionTemplate, nil)
				notificationService.EXPECT().SendConfirmed(gomock.Any(), subscriptionTemplate.Email, subscriptionTemplate.Token, subscriptionTemplate.Frequency)
			},
			expectedError: nil,
		},
		{
			name:  "Invalid token",
			token: "xxxx",
			mockBehavior: func(
				uc *mock_services.MockSubscriptionUseCase,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationServiceI,
				token string,
			) {
				tokenManager.EXPECT().Validate(gomock.Any(), token).Return(false)
			},
			expectedError: apperrors.InvalidTokenError,
		},
		{
			name:  "Token not found",
			token: "59d29860-39fa-3c2b-245a-3e9152b54e",
			mockBehavior: func(
				uc *mock_services.MockSubscriptionUseCase,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationServiceI,
				token string,
			) {
				tokenManager.EXPECT().Validate(gomock.Any(), token).Return(true)
				uc.EXPECT().Confirm(gomock.Any(), token).Return(nil, apperrors.TokenNotFoundError)

			},
			expectedError: apperrors.TokenNotFoundError,
		},
		{
			name:  "Database error",
			token: "59d29860-39fa-3c2b-245a-3e9152b54e",
			mockBehavior: func(
				uc *mock_services.MockSubscriptionUseCase,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationServiceI,
				token string,
			) {
				tokenManager.EXPECT().Validate(gomock.Any(), token).Return(true)
				uc.EXPECT().Confirm(gomock.Any(), token).Return(nil, apperrors.DatabaseError)

			},
			expectedError: apperrors.DatabaseError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			loggerMock := mock_logger.NewMockLogger(ctrl)
			tokenManagerMock := mock_services.NewMockTokenManager(ctrl)
			notificationServiceMock := mock_services.NewMockNotificationServiceI(ctrl)
			subscriptionUCMock := mock_services.NewMockSubscriptionUseCase(ctrl)

			testCase.mockBehavior(subscriptionUCMock, tokenManagerMock, notificationServiceMock, testCase.token)

			subscriptionService := NewSubscriptionService(subscriptionUCMock, tokenManagerMock, notificationServiceMock, loggerMock)

			err := subscriptionService.Confirm(context.Background(), testCase.token)

			assert.Equal(t, testCase.expectedError, err)

		})
	}
}

func TestSubscriptionService_Unsubscribe(t *testing.T) {

	testTable := []struct {
		name          string
		token         string
		mockBehavior  mockConfirmBehavior
		expectedError error
	}{
		{
			name:  "Successful",
			token: "59d29860-39fa-4c9b-845a-3e91eab42e4b",
			mockBehavior: func(
				uc *mock_services.MockSubscriptionUseCase,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationServiceI,
				token string,
			) {
				tokenManager.EXPECT().Validate(gomock.Any(), token).Return(true)
				uc.EXPECT().Unsubscribe(gomock.Any(), token).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:  "Invalid token",
			token: "xxxx",
			mockBehavior: func(
				uc *mock_services.MockSubscriptionUseCase,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationServiceI,
				token string,
			) {
				tokenManager.EXPECT().Validate(gomock.Any(), token).Return(false)
			},
			expectedError: apperrors.InvalidTokenError,
		},
		{
			name:  "Database error",
			token: "2812b8c0-44bh-aq9b-889y-3ey1oab42e4b",
			mockBehavior: func(
				uc *mock_services.MockSubscriptionUseCase,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationServiceI,
				token string,
			) {
				tokenManager.EXPECT().Validate(gomock.Any(), token).Return(true)
				uc.EXPECT().Unsubscribe(gomock.Any(), token).Return(apperrors.DatabaseError)
			},
			expectedError: apperrors.DatabaseError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			loggerMock := mock_logger.NewMockLogger(ctrl)
			tokenManagerMock := mock_services.NewMockTokenManager(ctrl)
			notificationServiceMock := mock_services.NewMockNotificationServiceI(ctrl)
			subscriptionUCMock := mock_services.NewMockSubscriptionUseCase(ctrl)

			testCase.mockBehavior(subscriptionUCMock, tokenManagerMock, notificationServiceMock, testCase.token)

			subscriptionService := NewSubscriptionService(subscriptionUCMock, tokenManagerMock, notificationServiceMock, loggerMock)

			err := subscriptionService.Unsubscribe(context.Background(), testCase.token)

			assert.Equal(t, testCase.expectedError, err)

		})
	}

}
