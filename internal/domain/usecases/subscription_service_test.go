package usecases

import (
	"context"
	"testing"
	"weather-forecast/internal/domain/models"
	mock_services "weather-forecast/internal/domain/usecases/mocks"
	"weather-forecast/internal/infrastructure/apperrors"
	mock_logger "weather-forecast/internal/infrastructure/logger/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type (
	mockSubscribeBehavior func(
		repo *mock_services.MockSubscriptionRepository,
		tokenManager *mock_services.MockTokenManager,
		notificationService *mock_services.MockNotificationSender,
		email,
		city,
		frequency,
		token string,
	)
	mockConfirmBehavior func(
		repo *mock_services.MockSubscriptionRepository,
		tokenManager *mock_services.MockTokenManager,
		notificationService *mock_services.MockNotificationSender,
		token string,
	)
)

var (
	subscriptionTemplate = models.Subscription{
		ID:        1,
		Email:     "test@test.com",
		Frequency: models.Hourly,
		City:      "Kyiv",
		Token:     "59d29860-39fa-4c9b-845a-3e91eab42e4b",
		Confirmed: true,
	}
)

func TestSubsctiptionService_Subscribe(t *testing.T) {

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
				repo *mock_services.MockSubscriptionRepository,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationSender,
				email,
				city,
				frequency,
				token string) {
				repo.EXPECT().GetByEmail(gomock.Any(), email).Return(nil, nil)
				tokenManager.EXPECT().Generate(gomock.Any()).Return(token)
				passedSubscription := models.Subscription{
					Email:     email,
					City:      city,
					Frequency: models.Frequency(frequency),
					Token:     token,
				}
				repo.EXPECT().Create(gomock.Any(), passedSubscription).Return(&subscriptionTemplate, nil)
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
				repo *mock_services.MockSubscriptionRepository,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationSender,
				email,
				city,
				frequency,
				token string) {
				passedSubscription := models.Subscription{
					Email:     email,
					City:      city,
					Frequency: models.Frequency(frequency),
					Token:     token,
				}
				repo.EXPECT().GetByEmail(gomock.Any(), passedSubscription.Email).Return(&passedSubscription, nil)

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
				repo *mock_services.MockSubscriptionRepository,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationSender,
				email,
				city,
				frequency,
				token string) {

				repo.EXPECT().GetByEmail(gomock.Any(), email).Return(nil, nil)

				tokenManager.EXPECT().Generate(gomock.Any()).Return(token)
				passedSubscription := models.Subscription{
					Email:     email,
					City:      city,
					Frequency: models.Frequency(frequency),
					Token:     token,
				}

				repo.EXPECT().Create(gomock.Any(), passedSubscription).Return(nil, apperrors.DatabaseError)

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
			notificationServiceMock := mock_services.NewMockNotificationSender(ctrl)
			subscriptionRepo := mock_services.NewMockSubscriptionRepository(ctrl)

			testCase.mockBehavior(subscriptionRepo, tokenManagerMock, notificationServiceMock, testCase.email, testCase.city, testCase.frequency, testCase.token)

			subscriptionService := NewSubscriptionService(subscriptionRepo, tokenManagerMock, notificationServiceMock, loggerMock)
			res, err := subscriptionService.Subscribe(context.Background(), testCase.email, testCase.frequency, testCase.city)

			assert.Equal(t, testCase.expectedResult, res)
			assert.Equal(t, testCase.expectedError, err)

		})
	}

}

func TestSubscriptionService_Confirm(t *testing.T) {

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
				repo *mock_services.MockSubscriptionRepository,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationSender,
				token string,
			) {
				receivedSubsc := subscriptionTemplate
				receivedSubsc.Confirmed = false
				tokenManager.EXPECT().Validate(gomock.Any(), token).Return(true)
				repo.EXPECT().GetByToken(gomock.Any(), token).Return(&receivedSubsc, nil)
				passedToUpdate := receivedSubsc
				passedToUpdate.Confirmed = true
				repo.EXPECT().Update(gomock.Any(), passedToUpdate).Return(&passedToUpdate, nil)
				notificationService.EXPECT().SendConfirmed(gomock.Any(), subscriptionTemplate.Email, subscriptionTemplate.Token, subscriptionTemplate.Frequency)
			},
			expectedError: nil,
		},
		{
			name:  "Invalid token",
			token: "xxxx",
			mockBehavior: func(
				repo *mock_services.MockSubscriptionRepository,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationSender,
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
				repo *mock_services.MockSubscriptionRepository,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationSender,
				token string,
			) {
				tokenManager.EXPECT().Validate(gomock.Any(), token).Return(true)
				repo.EXPECT().GetByToken(gomock.Any(), token).Return(nil, nil)

			},
			expectedError: apperrors.TokenNotFoundError,
		},
		{
			name:  "Database error getting token",
			token: "59d29860-39fa-3c2b-245a-3e9152b54e",
			mockBehavior: func(
				repo *mock_services.MockSubscriptionRepository,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationSender,
				token string,
			) {
				tokenManager.EXPECT().Validate(gomock.Any(), token).Return(true)
				repo.EXPECT().GetByToken(gomock.Any(), token).Return(nil, apperrors.DatabaseError)

			},
			expectedError: apperrors.DatabaseError,
		},
		{
			name:  "Database error updating subsc",
			token: "59d29860-39fa-3c2b-245a-3e9152b54e",
			mockBehavior: func(
				repo *mock_services.MockSubscriptionRepository,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationSender,
				token string,
			) {
				tokenManager.EXPECT().Validate(gomock.Any(), token).Return(true)
				receivedSubscription := subscriptionTemplate
				receivedSubscription.Confirmed = false
				repo.EXPECT().GetByToken(gomock.Any(), token).Return(&receivedSubscription, nil)
				passedToUpdate := receivedSubscription
				passedToUpdate.Confirmed = true
				repo.EXPECT().Update(gomock.Any(), passedToUpdate).Return(nil, apperrors.DatabaseError)

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
			notificationServiceMock := mock_services.NewMockNotificationSender(ctrl)
			subscriptionRepo := mock_services.NewMockSubscriptionRepository(ctrl)

			testCase.mockBehavior(subscriptionRepo, tokenManagerMock, notificationServiceMock, testCase.token)

			subscriptionService := NewSubscriptionService(subscriptionRepo, tokenManagerMock, notificationServiceMock, loggerMock)

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
				repo *mock_services.MockSubscriptionRepository,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationSender,
				token string,
			) {
				tokenManager.EXPECT().Validate(gomock.Any(), token).Return(true)
				receivedSubscription := subscriptionTemplate
				repo.EXPECT().GetByToken(gomock.Any(), token).Return(&receivedSubscription, nil)
				repo.EXPECT().DeleteByToken(gomock.Any(), receivedSubscription.Token).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:  "Invalid token",
			token: "xxxx",
			mockBehavior: func(
				repo *mock_services.MockSubscriptionRepository,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationSender,
				token string,
			) {
				tokenManager.EXPECT().Validate(gomock.Any(), token).Return(false)
			},
			expectedError: apperrors.InvalidTokenError,
		},
		{
			name:  "Database error getting token",
			token: "2812b8c0-44bh-aq9b-889y-3ey1oab42e4b",
			mockBehavior: func(
				repo *mock_services.MockSubscriptionRepository,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationSender,
				token string,
			) {
				tokenManager.EXPECT().Validate(gomock.Any(), token).Return(true)
				repo.EXPECT().GetByToken(gomock.Any(), token).Return(nil, apperrors.DatabaseError)
			},
			expectedError: apperrors.DatabaseError,
		},
		{
			name:  "Database error deleting subsc",
			token: "2812b8c0-44bh-aq9b-889y-3ey1oab42e4b",
			mockBehavior: func(
				repo *mock_services.MockSubscriptionRepository,
				tokenManager *mock_services.MockTokenManager,
				notificationService *mock_services.MockNotificationSender,
				token string,
			) {
				receivedSubscription := subscriptionTemplate
				receivedSubscription.Token = token
				tokenManager.EXPECT().Validate(gomock.Any(), token).Return(true)
				repo.EXPECT().GetByToken(gomock.Any(), token).Return(&receivedSubscription, nil)
				repo.EXPECT().DeleteByToken(gomock.Any(), receivedSubscription.Token).Return(apperrors.DatabaseError)

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
			notificationServiceMock := mock_services.NewMockNotificationSender(ctrl)
			subscriptionRepo := mock_services.NewMockSubscriptionRepository(ctrl)

			testCase.mockBehavior(subscriptionRepo, tokenManagerMock, notificationServiceMock, testCase.token)

			subscriptionService := NewSubscriptionService(subscriptionRepo, tokenManagerMock, notificationServiceMock, loggerMock)

			err := subscriptionService.Unsubscribe(context.Background(), testCase.token)

			assert.Equal(t, testCase.expectedError, err)

		})
	}

}
