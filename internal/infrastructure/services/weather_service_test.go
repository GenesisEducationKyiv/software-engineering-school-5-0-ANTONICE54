package services

import (
	"context"
	"testing"
	"weather-forecast/internal/domain/models"
	infraerrors "weather-forecast/internal/infrastructure/errors"

	mock_logger "weather-forecast/internal/infrastructure/logger/mocks"
	mock_services "weather-forecast/internal/infrastructure/services/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockGetWeatherByCityBehavior func(p *mock_services.MockWeatherProvider, city string)

func TestWeatherService_GetWeatherByCity(t *testing.T) {

	testTable := []struct {
		name            string
		city            string
		mockBehavior    mockGetWeatherByCityBehavior
		expectedWeather *models.Weather
		expectedError   error
	}{
		{
			name: "Successful",
			city: "Kyiv",
			mockBehavior: func(p *mock_services.MockWeatherProvider, city string) {
				weather := models.Weather{
					Temperature: 21.5,
					Humidity:    40,
					Description: "Partly cloudy",
				}
				p.EXPECT().GetWeatherByCity(gomock.Any(), city).Return(&weather, nil)
			},
			expectedWeather: &models.Weather{
				Temperature: 21.5,
				Humidity:    40,
				Description: "Partly cloudy",
			},
			expectedError: nil,
		},
		{
			name: "City not found",
			city: "Jitomyr",
			mockBehavior: func(p *mock_services.MockWeatherProvider, city string) {
				p.EXPECT().GetWeatherByCity(gomock.Any(), city).Return(nil, infraerrors.CityNotFoundError)
			},
			expectedWeather: nil,
			expectedError:   infraerrors.CityNotFoundError,
		},
		{
			name: "Invalid city",
			city: "",
			mockBehavior: func(p *mock_services.MockWeatherProvider, city string) {
				p.EXPECT().GetWeatherByCity(gomock.Any(), city).Return(nil, infraerrors.GetWeatherError)
			},
			expectedWeather: nil,
			expectedError:   infraerrors.GetWeatherError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			loggerMock := mock_logger.NewMockLogger(ctrl)
			weatherProviderMock := mock_services.NewMockWeatherProvider(ctrl)
			weatherService := NewWeatherService(weatherProviderMock, loggerMock)

			testCase.mockBehavior(weatherProviderMock, testCase.city)

			res, err := weatherService.GetWeatherByCity(context.Background(), testCase.city)

			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedWeather, res)
		})
	}

}
