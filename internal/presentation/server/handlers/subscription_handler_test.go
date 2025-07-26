package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	domainerr "weather-forecast/internal/domain/errors"
	"weather-forecast/internal/domain/models"
	infraerrors "weather-forecast/internal/infrastructure/errors"

	mock_logger "weather-forecast/internal/infrastructure/logger/mocks"
	mock_handlers "weather-forecast/internal/presentation/server/handlers/mock"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type (
	mockSubscribeBehavior   func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, request SubscribeRequest)
	mockConfirmBehavior     func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, token string)
	mockUnsubscribeBehavior func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, token string)
)

func TestSubscriptionHandler_Subscribe(t *testing.T) {

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
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: errToString("invalid request format"),
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

				s.EXPECT().Subscribe(gomock.Any(), request.Email, request.Frequency, request.City).Return(nil, domainerr.ErrAlreadySubscribed)
			},
			expectedStatusCode:   http.StatusConflict,
			expectedResponseBody: errToString(domainerr.ErrAlreadySubscribed.Error()),
		},
		{
			name:      "Database error",
			inputBody: `{"email":"test@test.com","city":"Kyiv","frequency":"hourly"}`,
			inputData: SubscribeRequest{
				Email:     "test@test.com",
				City:      "Kyiv",
				Frequency: "hourly",
			},
			mockBehavior: func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, request SubscribeRequest) {
				logger.EXPECT().Warnf(gomock.Any(), gomock.Any())
				s.EXPECT().Subscribe(gomock.Any(), request.Email, request.Frequency, request.City).Return(nil, infraerrors.ErrDatabase)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: errToString("internal server error"),
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
				s.EXPECT().Confirm(gomock.Any(), token).Return(domainerr.ErrInvalidToken)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: errToString(domainerr.ErrInvalidToken.Error()),
		},
		{
			name:  "Token not found",
			token: "59a29260-39fa-4c9b-845a-4a23bb342e4b",
			mockBehavior: func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, token string) {
				s.EXPECT().Confirm(gomock.Any(), token).Return(domainerr.ErrTokenNotFound)
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: errToString(domainerr.ErrTokenNotFound.Error()),
		},
		{
			name:  "Database error",
			token: "59d29860-39fa-4c9b-845a-3e91eab42e4b",
			mockBehavior: func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, token string) {
				logger.EXPECT().Warnf(gomock.Any(), gomock.Any())
				s.EXPECT().Confirm(gomock.Any(), token).Return(infraerrors.ErrDatabase)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: errToString("internal server error"),
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

	testTable := []struct {
		name                 string
		token                string
		mockBehavior         mockUnsubscribeBehavior
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
				s.EXPECT().Unsubscribe(gomock.Any(), token).Return(domainerr.ErrInvalidToken)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: errToString(domainerr.ErrInvalidToken.Error()),
		},
		{
			name:  "Token not found",
			token: "59a29260-39fa-4c9b-845a-4a23bb342e4b",
			mockBehavior: func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, token string) {
				s.EXPECT().Unsubscribe(gomock.Any(), token).Return(domainerr.ErrTokenNotFound)
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: errToString(domainerr.ErrTokenNotFound.Error()),
		},
		{
			name:  "Database error",
			token: "59d29860-39fa-4c9b-845a-3e91eab42e4b",
			mockBehavior: func(s *mock_handlers.MockSubsctiptionService, logger *mock_logger.MockLogger, token string) {
				logger.EXPECT().Warnf(gomock.Any(), gomock.Any())
				s.EXPECT().Unsubscribe(gomock.Any(), token).Return(infraerrors.ErrDatabase)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: errToString("internal server error"),
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
