package integration

import (
	"context"
	"net/http"
	"testing"
	"time"
	"weather-forecast/pkg/proto/weather"
	stub_logger "weather-forecast/pkg/stubs/logger"
	"weather-service/internal/domain/models"
	"weather-service/internal/infrastructure/cache"
	"weather-service/internal/infrastructure/clients/openweather"
	"weather-service/internal/infrastructure/clients/weatherapi"
	"weather-service/internal/presentation/server/handlers"
	"weather-service/tests/integration/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	testWeather = models.Weather{
		Temperature: 22.5,
		Humidity:    64,
		Description: "Partly cloudy",
	}
)

func TestGetWeather_Success(t *testing.T) {

	city := "Kyiv"
	weatherAPISuccessResponse := weatherapi.WeatherSuccessResponse{
		Current: weatherapi.WeatherCurrentResponse{
			TempC: testWeather.Temperature,
			Condition: weatherapi.WeatherConditionResponse{
				Text: testWeather.Description,
			},
			Humidity: testWeather.Humidity,
		},
	}

	weatherAPIServerMock := setupWeatherAPIMock(t, weatherAPISuccessResponse, http.StatusOK, city)

	openWeatherAPIServcerMock := setupOpenWeatherMock(t, nil, 0, "", false)

	weatherHandler := setupWeatherHandler(weatherAPIServerMock.URL, openWeatherAPIServcerMock.URL)

	requestBody := &weather.GetWeatherRequest{
		City: city,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := weatherHandler.GetWeather(ctx, requestBody)
	require.NoError(t, err)

	assertWeatherResponse(t, resp, testWeather)

}

func TestGetWeather_WeatherAPI_Failed(t *testing.T) {
	city := "Kyiv"

	openWeatherSuccessResponse := openweather.OpenWeatherSuccessResponse{
		Weather: []openweather.OpenWeatherDescriptionResponse{
			{Description: testWeather.Description},
		},
		Main: openweather.OpenWeatherMainResponse{
			Temperature: testWeather.Temperature,
			Humidity:    testWeather.Humidity,
		},
	}

	weatherAPIErrorResponseBody := weatherapi.WeatherErrorResponse{
		Error: weatherapi.WeatherErrorDetails{
			Code:    9999,
			Message: "Internal application error.",
		},
	}

	weatherAPIServerMock := setupWeatherAPIMock(t, weatherAPIErrorResponseBody, http.StatusBadRequest, city)

	openWeatherAPIServcerMock := setupOpenWeatherMock(t, openWeatherSuccessResponse, http.StatusOK, city, true)

	weatherHandler := setupWeatherHandler(weatherAPIServerMock.URL, openWeatherAPIServcerMock.URL)

	requestBody := &weather.GetWeatherRequest{
		City: city,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := weatherHandler.GetWeather(ctx, requestBody)
	require.NoError(t, err)

	assertWeatherResponse(t, resp, testWeather)

}

func TestGetWeather_RedisCacheFlow(t *testing.T) {
	city := "Kyiv"

	testRedis := testutils.SetupTestRedis(t)
	redisCache, err := cache.NewRedis(testRedis.ConnectionString(), stub_logger.New())
	require.NoError(t, err)
	metrics := testutils.NewInMemoryMetrics()

	weatherAPISuccessResponse := weatherapi.WeatherSuccessResponse{
		Current: weatherapi.WeatherCurrentResponse{
			TempC: testWeather.Temperature,
			Condition: weatherapi.WeatherConditionResponse{
				Text: testWeather.Description,
			},
			Humidity: testWeather.Humidity,
		},
	}

	weatherAPIServerMock := setupWeatherAPIMock(t, weatherAPISuccessResponse, http.StatusOK, city)
	openWeatherAPIServerMock := setupOpenWeatherMock(t, nil, 0, "", false)
	weatherHandler := setupWeatherHandlerWithCache(redisCache, metrics, weatherAPIServerMock.URL, openWeatherAPIServerMock.URL)

	requestBody := &weather.GetWeatherRequest{City: city}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := weatherHandler.GetWeather(ctx, requestBody)
	require.NoError(t, err)

	assertWeatherResponse(t, resp, testWeather)

	hits, misses, errors := metrics.Stats()
	assert.Equal(t, 0, hits)
	assert.Equal(t, 1, misses)
	assert.Equal(t, 0, errors)
	assert.Equal(t, 1, testRedis.Size())

	resp, err = weatherHandler.GetWeather(ctx, requestBody)
	require.NoError(t, err)
	assertWeatherResponse(t, resp, testWeather)

	hits, misses, errors = metrics.Stats()
	assert.Equal(t, 1, hits)
	assert.Equal(t, 1, misses)
	assert.Equal(t, 0, errors)
	assert.Equal(t, 1, testRedis.Size())

}

func TestGetWeather_ErrorScenarios(t *testing.T) {

	const weatherAPINotFoundErrorCode = 1006
	const openWeatherNotFoundErrorCode = "404"

	testTable := []struct {
		name                  string
		city                  string
		weatherAPIResp        interface{}
		weatherAPICode        int
		openWeatherResp       interface{}
		openWeatherCode       int
		shouldCallOpenWeather bool
		expectedCode          codes.Code
		expectedErrorBody     string
	}{
		// {
		// 	name:                  "Invalid City",
		// 	city:                  "123",
		// 	weatherAPIResp:        nil,
		// 	weatherAPICode:        0,
		// 	openWeatherResp:       nil,
		// 	openWeatherCode:       0,
		// 	shouldCallOpenWeather: false,
		// 	expectedCode:          codes.InvalidArgument,
		// 	expectedErrorBody:     "invalid request",
		// },
		{
			name: "City Not Found",
			city: "Odeca",
			weatherAPIResp: weatherapi.WeatherErrorResponse{
				Error: weatherapi.WeatherErrorDetails{
					Code:    weatherAPINotFoundErrorCode,
					Message: "No matching location found.",
				},
			},
			weatherAPICode: http.StatusBadRequest,
			openWeatherResp: openweather.OpenWeatherErrorResponse{
				Cod:     openWeatherNotFoundErrorCode,
				Message: "city not found",
			},
			openWeatherCode:       http.StatusNotFound,
			shouldCallOpenWeather: true,
			expectedCode:          codes.NotFound,
			expectedErrorBody:     "there is no city with such name",
		},
		{
			name: "Both Providers Fail",
			city: "TestCity",
			weatherAPIResp: weatherapi.WeatherErrorResponse{
				Error: weatherapi.WeatherErrorDetails{
					Code:    9999,
					Message: "Internal application error.",
				},
			},
			weatherAPICode: http.StatusBadRequest,
			openWeatherResp: openweather.OpenWeatherErrorResponse{
				Cod:     "500",
				Message: "internal server error",
			},

			openWeatherCode:       http.StatusInternalServerError,
			shouldCallOpenWeather: true,

			expectedCode:      codes.Internal,
			expectedErrorBody: "failed to get weather",
		},
		{
			name:                  "Timout Exceeded",
			city:                  "Odeca",
			weatherAPIResp:        nil,
			weatherAPICode:        0,
			openWeatherResp:       nil,
			openWeatherCode:       0,
			shouldCallOpenWeather: false,
			expectedCode:          codes.Internal,
			expectedErrorBody:     "failed to get weather",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			var weatherHandler *handlers.WeatherHandler
			if testCase.weatherAPIResp != nil {
				weatherAPIServerMock := setupWeatherAPIMock(t, testCase.weatherAPIResp, testCase.weatherAPICode, testCase.city)
				defer weatherAPIServerMock.Close()

				openWeatherAPIServcerMock := setupOpenWeatherMock(t, testCase.openWeatherResp, testCase.openWeatherCode, testCase.city, testCase.shouldCallOpenWeather)
				defer openWeatherAPIServcerMock.Close()

				weatherHandler = setupWeatherHandler(weatherAPIServerMock.URL, openWeatherAPIServcerMock.URL)

			} else {
				weatherHandler = setupWeatherHandler("", "")
			}

			requestBody := &weather.GetWeatherRequest{
				City: testCase.city,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := weatherHandler.GetWeather(ctx, requestBody)
			require.Error(t, err)
			assert.Nil(t, resp)

			grpcStatus, ok := status.FromError(err)
			require.True(t, ok)

			assert.Equal(t, testCase.expectedCode, grpcStatus.Code())
			assert.Contains(t, grpcStatus.Message(), testCase.expectedErrorBody)
		})
	}

}
