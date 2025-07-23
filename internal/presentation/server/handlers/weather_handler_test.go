package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"weather-forecast/internal/domain/models"
	infraerrors "weather-forecast/internal/infrastructure/errors"
	apierrors "weather-forecast/internal/presentation/errors"

	"weather-forecast/internal/presentation/httperrors"

	mock_logger "weather-forecast/internal/infrastructure/logger/mocks"
	mock_handlers "weather-forecast/internal/presentation/server/handlers/mock"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func errToString(errMsg map[string]any) string {

	bytes, err := json.Marshal(errMsg)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

type mockWeatherServiceBehavior func(s *mock_handlers.MockWeatherService, logger *mock_logger.MockLogger, city GetWeatherRequest)

func TestWeatherHandler_Get(t *testing.T) {

	testTable := []struct {
		name                 string
		inputBody            string
		inputCity            GetWeatherRequest
		mockBehavior         mockWeatherServiceBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Successful",
			inputBody: `{"city":"Kyiv"}`,
			inputCity: GetWeatherRequest{
				City: "Kyiv",
			},
			mockBehavior: func(s *mock_handlers.MockWeatherService, logger *mock_logger.MockLogger, city GetWeatherRequest) {
				weather := models.Weather{
					Temperature: 5.1,
					Humidity:    80,
					Description: "Cloudy",
				}
				s.EXPECT().GetWeatherByCity(gomock.Any(), city.City).Return(&weather, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"temperature":5.1,"humidity":80,"description":"Cloudy"}`,
		},
		{
			name:      "Invalid input",
			inputBody: `{"city":"АБВГ"}`,
			mockBehavior: func(s *mock_handlers.MockWeatherService, logger *mock_logger.MockLogger, city GetWeatherRequest) {
				logger.EXPECT().Warnf(gomock.Any(), gomock.Any())
			},
			expectedStatusCode:   httperrors.New(apierrors.ErrInvalidRequest).Status(),
			expectedResponseBody: errToString(httperrors.New(apierrors.ErrInvalidRequest).Body()),
		},
		{
			name:      "City not found",
			inputBody: `{"city":"ABCD"}`,
			inputCity: GetWeatherRequest{
				City: "ABCD",
			},
			mockBehavior: func(s *mock_handlers.MockWeatherService, logger *mock_logger.MockLogger, city GetWeatherRequest) {
				s.EXPECT().GetWeatherByCity(gomock.Any(), city.City).Return(nil, infraerrors.ErrCityNotFound)
			},
			expectedStatusCode:   httperrors.New(infraerrors.ErrCityNotFound).Status(),
			expectedResponseBody: errToString(httperrors.New(infraerrors.ErrCityNotFound).Body()),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			loggerMock := mock_logger.NewMockLogger(ctrl)
			weatherServiceMock := mock_handlers.NewMockWeatherService(ctrl)
			testCase.mockBehavior(weatherServiceMock, loggerMock, testCase.inputCity)
			weatherHandler := NewWeatherHandler(weatherServiceMock, loggerMock)

			router := gin.New()

			router.GET("/weather", weatherHandler.Get)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/weather", bytes.NewBufferString(testCase.inputBody))

			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())

		})
	}
}
