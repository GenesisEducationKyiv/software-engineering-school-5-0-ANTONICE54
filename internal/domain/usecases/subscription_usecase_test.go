package usecases

import (
	"context"
	"testing"
	"time"
	"weather-forecast/internal/domain/models"
	mock_usecases "weather-forecast/internal/domain/usecases/mocks"
	"weather-forecast/internal/infrastructure/apperrors"
	mock_logger "weather-forecast/internal/infrastructure/logger/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type (
	mockSubscriveBehavior   func(repo *mock_usecases.MockSubscriptionRepository, logger *mock_logger.MockLogger, subscription models.Subscription)
	mockConfirmBehavior     func(repo *mock_usecases.MockSubscriptionRepository, logger *mock_logger.MockLogger, token string)
	mockUnsubscribeBehavior func(repo *mock_usecases.MockSubscriptionRepository, logger *mock_logger.MockLogger, token string)
)

func TestSubscriptionUsecase_Subscribe(t *testing.T) {
	createdTime := time.Now()

	testTable := []struct {
		name               string
		passedSubscription models.Subscription
		mockBehavior       mockSubscriveBehavior
		expectedResult     *models.Subscription
		expectedError      error
	}{
		{
			name: "Successful",
			passedSubscription: models.Subscription{
				Email:     "test@gmail.com",
				City:      "Kyiv",
				Frequency: models.Hourly,
				Token:     "5a57909b-b980-465f-a72a-606560278fa2",
			},
			mockBehavior: func(repo *mock_usecases.MockSubscriptionRepository, logger *mock_logger.MockLogger, subscription models.Subscription) {
				createdSubsc := models.Subscription{
					ID:        1,
					Email:     subscription.Email,
					City:      subscription.City,
					Frequency: subscription.Frequency,
					Token:     subscription.Token,
					Confirmed: false,
					CreatedAt: createdTime,
				}
				repo.EXPECT().GetByEmail(gomock.Any(), subscription.Email).Return(nil, nil)
				repo.EXPECT().Create(gomock.Any(), subscription).Return(&createdSubsc, nil)
			},
			expectedResult: &models.Subscription{
				ID:        1,
				Email:     "test@gmail.com",
				City:      "Kyiv",
				Frequency: models.Hourly,
				Token:     "5a57909b-b980-465f-a72a-606560278fa2",
				Confirmed: false,
				CreatedAt: createdTime,
			},
			expectedError: nil,
		},
		{
			name: "Already subscribed",
			passedSubscription: models.Subscription{
				Email:     "test@gmail.com",
				City:      "Kyiv",
				Frequency: models.Hourly,
				Token:     "5a57909b-b980-465f-a72a-606560278fa2",
			},
			mockBehavior: func(repo *mock_usecases.MockSubscriptionRepository, logger *mock_logger.MockLogger, subscription models.Subscription) {

				receivedSubsc := models.Subscription{
					Email:     "test@gmail.com",
					City:      "Kyiv",
					Frequency: models.Hourly,
					Token:     "5a57909b-b980-465f-a72a-606560278fa2",
					Confirmed: false,
					CreatedAt: createdTime,
				}

				repo.EXPECT().GetByEmail(gomock.Any(), subscription.Email).Return(&receivedSubsc, nil)
				logger.EXPECT().Warnf(gomock.Any(), gomock.Any())
			},
			expectedResult: nil,
			expectedError:  apperrors.AlreadySubscribedError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			loggerMock := mock_logger.NewMockLogger(ctrl)
			subscRepoMock := mock_usecases.NewMockSubscriptionRepository(ctrl)
			subscUC := NewSubscriptionUseCase(subscRepoMock, loggerMock)

			testCase.mockBehavior(subscRepoMock, loggerMock, testCase.passedSubscription)

			res, err := subscUC.Subscribe(context.Background(), testCase.passedSubscription)

			assert.Equal(t, testCase.expectedResult, res)
			assert.Equal(t, testCase.expectedError, err)

		})
	}

}

func TestSubscriptionUsecase_Confirm(t *testing.T) {
	createdTime := time.Now()

	subscriptionTemplate := models.Subscription{
		ID:        1,
		Email:     "test@gmail.com",
		City:      "Kyiv",
		Frequency: models.Hourly,
		Token:     "5a57909b-b980-465f-a72a-606560278fa2",
		Confirmed: true,
		CreatedAt: createdTime,
	}

	testTable := []struct {
		name           string
		passedToken    string
		mockBehavior   mockConfirmBehavior
		expectedResult *models.Subscription
		expectedError  error
	}{
		{
			name:        "Successful",
			passedToken: "5a57909b-b980-465f-a72a-606560278fa2",
			mockBehavior: func(repo *mock_usecases.MockSubscriptionRepository, logger *mock_logger.MockLogger, token string) {
				receivedSubscription := subscriptionTemplate
				receivedSubscription.Confirmed = false
				repo.EXPECT().GetByToken(gomock.Any(), token).Return(&receivedSubscription, nil)
				updatedSubscription := receivedSubscription
				updatedSubscription.Confirmed = true
				repo.EXPECT().Update(gomock.Any(), receivedSubscription).Return(&updatedSubscription, nil)
			},
			expectedResult: &subscriptionTemplate,
			expectedError:  nil,
		},
		{
			name:        "Not found",
			passedToken: "52a579r9b-b980-465f-a72a-606565254fa2",
			mockBehavior: func(repo *mock_usecases.MockSubscriptionRepository, logger *mock_logger.MockLogger, token string) {

				repo.EXPECT().GetByToken(gomock.Any(), token).Return(nil, nil)
			},
			expectedResult: nil,
			expectedError:  apperrors.TokenNotFoundError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			loggerMock := mock_logger.NewMockLogger(ctrl)
			subscRepoMock := mock_usecases.NewMockSubscriptionRepository(ctrl)
			subscUC := NewSubscriptionUseCase(subscRepoMock, loggerMock)

			testCase.mockBehavior(subscRepoMock, loggerMock, testCase.passedToken)

			res, err := subscUC.Confirm(context.Background(), testCase.passedToken)

			assert.Equal(t, testCase.expectedResult, res)
			assert.Equal(t, testCase.expectedError, err)

		})
	}

}

func TestSubscriptionUsecase_Unsubscribe(t *testing.T) {
	createdTime := time.Now()

	subscriptionTemplate := models.Subscription{
		ID:        1,
		Email:     "test@gmail.com",
		City:      "Kyiv",
		Frequency: models.Hourly,
		Token:     "5a57909b-b980-465f-a72a-606560278fa2",
		Confirmed: true,
		CreatedAt: createdTime,
	}

	testTable := []struct {
		name          string
		passedToken   string
		mockBehavior  mockUnsubscribeBehavior
		expectedError error
	}{
		{
			name:        "Successful",
			passedToken: "5a57909b-b980-465f-a72a-606560278fa2",
			mockBehavior: func(repo *mock_usecases.MockSubscriptionRepository, logger *mock_logger.MockLogger, token string) {
				receivedSubscription := subscriptionTemplate
				repo.EXPECT().GetByToken(gomock.Any(), token).Return(&receivedSubscription, nil)
				repo.EXPECT().DeleteByToken(gomock.Any(), receivedSubscription.Token).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:        "Not found",
			passedToken: "52a579r9b-b980-465f-a72a-606565254fa2",
			mockBehavior: func(repo *mock_usecases.MockSubscriptionRepository, logger *mock_logger.MockLogger, token string) {

				repo.EXPECT().GetByToken(gomock.Any(), token).Return(nil, nil)
			},
			expectedError: apperrors.TokenNotFoundError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			loggerMock := mock_logger.NewMockLogger(ctrl)
			subscRepoMock := mock_usecases.NewMockSubscriptionRepository(ctrl)
			subscUC := NewSubscriptionUseCase(subscRepoMock, loggerMock)

			testCase.mockBehavior(subscRepoMock, loggerMock, testCase.passedToken)

			err := subscUC.Unsubscribe(context.Background(), testCase.passedToken)

			assert.Equal(t, testCase.expectedError, err)

		})
	}

}
