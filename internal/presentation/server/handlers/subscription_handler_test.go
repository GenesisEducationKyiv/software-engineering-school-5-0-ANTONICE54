package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/apperrors"
	mock_logger "weather-forecast/internal/infrastructure/logger/mocks"
	mock_handlers "weather-forecast/internal/presentation/server/handlers/mock"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptionHandler_Subscribe(t *testing.T) {
	type mockSubscribeBehavior func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, request SubscribeRequest)

	testTable := []struct {
		name                 string
		inputBody            string
		inputData            SubscribeRequest
		mockBehavior         mockSubscribeBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Successful",
			inputBody: `{"email":"test@test.com","city":"Kyiv","frequency":"hourly"}`,
			inputData: SubscribeRequest{
				Email:     "test@test.com",
				City:      "Kyiv",
				Frequency: "hourly",
			},
			mockBehavior: func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, request SubscribeRequest) {
				subscription := &models.Subscription{
					Email:     request.Email,
					Frequency: models.Frequency(request.Frequency),
					City:      request.City,
				}
				s.EXPECT().Subscribe(gomock.Any(), request.Email, request.Frequency, request.City).Return(subscription, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"Subscription successful. Confirmation email sent."}`,
		},
		{
			name:      "Invalid input",
			inputBody: `{"email":"test","city":"Kyiv123","frequency":"monthly"}`,
			inputData: SubscribeRequest{
				Email:     "test",
				City:      "Kyiv123",
				Frequency: "monthly",
			},
			mockBehavior: func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, request SubscribeRequest) {
				logger.EXPECT().Warnf(gomock.Any(), gomock.Any())
			},
			expectedStatusCode:   apperrors.InvalidRequestError.Status(),
			expectedResponseBody: errToString(apperrors.InvalidRequestError.JSON()),
		},
		{
			name:      "Already subscribed",
			inputBody: `{"email":"test@test.com","city":"Kyiv","frequency":"hourly"}`,
			inputData: SubscribeRequest{
				Email:     "test@test.com",
				City:      "Kyiv",
				Frequency: "hourly",
			},
			mockBehavior: func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, request SubscribeRequest) {

				s.EXPECT().Subscribe(gomock.Any(), request.Email, request.Frequency, request.City).Return(nil, apperrors.AlreadySubscribedError)
			},
			expectedStatusCode:   apperrors.AlreadySubscribedError.Status(),
			expectedResponseBody: errToString(apperrors.AlreadySubscribedError.JSON()),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			loggerMock := mock_logger.NewMockLogger(ctrl)
			subscriptionServiceMock := mock_handlers.NewMockSubsctiptionService(ctrl)
			testCase.mockBehavior(subscriptionServiceMock, loggerMock, testCase.inputData)
			subscriptionHandler := NewSubscriptionHandler(subscriptionServiceMock, loggerMock)

			router := gin.New()

			router.POST("/subscribe", subscriptionHandler.Subscribe)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/subscribe", bytes.NewBufferString(testCase.inputBody))

			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())

		})
	}
}

func TestSubscriptionHandler_Confirm(t *testing.T) {
	type mockConfirmBehavior func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, token string)

	testTable := []struct {
		name                 string
		token                string
		mockBehavior         mockConfirmBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:  "Successful",
			token: "59d29860-39fa-4c9b-845a-3e91eab42e4b",
			mockBehavior: func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, token string) {
				s.EXPECT().Confirm(gomock.Any(), token).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"Subscription confirmed."}`,
		},
		{
			name:  "Invalid token",
			token: "xxxxx",
			mockBehavior: func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, token string) {
				s.EXPECT().Confirm(gomock.Any(), token).Return(apperrors.InvalidTokenError)
			},
			expectedStatusCode:   apperrors.InvalidTokenError.Status(),
			expectedResponseBody: errToString(apperrors.InvalidTokenError.JSON()),
		},
		{
			name:  "Token not found",
			token: "59a29260-39fa-4c9b-845a-4a23bb342e4b",
			mockBehavior: func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, token string) {
				s.EXPECT().Confirm(gomock.Any(), token).Return(apperrors.TokenNotFoundError)
			},
			expectedStatusCode:   apperrors.TokenNotFoundError.Status(),
			expectedResponseBody: errToString(apperrors.TokenNotFoundError.JSON()),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			loggerMock := mock_logger.NewMockLogger(ctrl)
			subscriptionServiceMock := mock_handlers.NewMockSubsctiptionService(ctrl)
			testCase.mockBehavior(subscriptionServiceMock, loggerMock, testCase.token)
			subscriptionHandler := NewSubscriptionHandler(subscriptionServiceMock, loggerMock)

			router := gin.New()

			router.GET("/confirm/:token", subscriptionHandler.Confirm)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/confirm/%s", testCase.token), nil)

			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())

		})
	}

}

func TestSubscriptionHandler_Unsubcribe(t *testing.T) {
	type mockConfirmBehavior func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, token string)

	testTable := []struct {
		name                 string
		token                string
		mockBehavior         mockConfirmBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:  "Successful",
			token: "59d29860-39fa-4c9b-845a-3e91eab42e4b",
			mockBehavior: func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, token string) {
				s.EXPECT().Unsubscribe(gomock.Any(), token).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"Unsubscribed successfuly."}`,
		},
		{
			name:  "Invalid token",
			token: "xxxxx",
			mockBehavior: func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, token string) {
				s.EXPECT().Unsubscribe(gomock.Any(), token).Return(apperrors.InvalidTokenError)
			},
			expectedStatusCode:   apperrors.InvalidTokenError.Status(),
			expectedResponseBody: errToString(apperrors.InvalidTokenError.JSON()),
		},
		{
			name:  "Token not found",
			token: "59a29260-39fa-4c9b-845a-4a23bb342e4b",
			mockBehavior: func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, token string) {
				s.EXPECT().Unsubscribe(gomock.Any(), token).Return(apperrors.TokenNotFoundError)
			},
			expectedStatusCode:   apperrors.TokenNotFoundError.Status(),
			expectedResponseBody: errToString(apperrors.TokenNotFoundError.JSON()),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			loggerMock := mock_logger.NewMockLogger(ctrl)
			subscriptionServiceMock := mock_handlers.NewMockSubsctiptionService(ctrl)
			testCase.mockBehavior(subscriptionServiceMock, loggerMock, testCase.token)
			subscriptionHandler := NewSubscriptionHandler(subscriptionServiceMock, loggerMock)

			router := gin.New()

			router.GET("/unsubscribe/:token", subscriptionHandler.Unsubscribe)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/unsubscribe/%s", testCase.token), nil)

			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())

		})
	}
}
