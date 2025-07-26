package integration

import (
	"encoding/json"
	"net/http"
	"testing"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/providers"
	"weather-forecast/internal/presentation/server/handlers"
	"weather-forecast/tests/integration/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	weatherAPISuccessResponse := providers.GetWeatherSuccessResponse{
		Current: providers.GetWeatherCurrentResponse{
			TempC: testWeather.Temperature,
			Condition: providers.GetWeatherConditionResponse{
				Text: testWeather.Description,
			},
			Humidity: testWeather.Humidity,
		},
	}

	weatherAPIServerMock := setupWeatherAPIMock(t, weatherAPISuccessResponse, http.StatusOK, city)

	openWeatherAPIServcerMock := setupOpenWeatherMock(t, nil, 0, "", false)

	weatherHandler := setupWeatherHandler(weatherAPIServerMock.URL, openWeatherAPIServcerMock.URL)
	router := setupWeatherRouter(weatherHandler)

	requestBody := handlers.GetWeatherRequest{
		City: city,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	resp := requestWeather(router, body)

	assertWeatherResponse(t, resp, testWeather)

}

func TestGetWeather_WeatherAPI_Failed(t *testing.T) {
	city := "Kyiv"

	openWeatherSuccessResponse := providers.GetOpenWeatherSuccessResponse{
		Weather: []providers.GetOpenWeatherDescriptionResponse{
			{Description: testWeather.Description},
		},
		Main: providers.GetOpenWeatherMainResponse{
			Temperature: testWeather.Temperature,
			Humidity:    testWeather.Humidity,
		},
	}

	weatherAPIErrorResponseBody := providers.GetWeatherErrorResponse{
		Error: providers.GetWeatherErrorDetails{
			Code:    9999,
			Message: "Internal application error.",
		},
	}

	weatherAPIServerMock := setupWeatherAPIMock(t, weatherAPIErrorResponseBody, http.StatusBadRequest, city)

	openWeatherAPIServcerMock := setupOpenWeatherMock(t, openWeatherSuccessResponse, http.StatusOK, city, true)

	weatherHandler := setupWeatherHandler(weatherAPIServerMock.URL, openWeatherAPIServcerMock.URL)
	router := setupWeatherRouter(weatherHandler)

	requestBody := handlers.GetWeatherRequest{
		City: city,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)
	resp := requestWeather(router, body)

	assertWeatherResponse(t, resp, testWeather)

}

func TestGetWeather_CacheFlow(t *testing.T) {
	city := "Kyiv"

	cache := testutils.NewInMemoryCache()
	metrics := testutils.NewInMemoryMetrics()

	weatherAPISuccessResponse := providers.GetWeatherSuccessResponse{
		Current: providers.GetWeatherCurrentResponse{
			TempC: testWeather.Temperature,
			Condition: providers.GetWeatherConditionResponse{
				Text: testWeather.Description,
			},
			Humidity: testWeather.Humidity,
		},
	}

	weatherAPIServerMock := setupWeatherAPIMock(t, weatherAPISuccessResponse, http.StatusOK, city)
	openWeatherAPIServerMock := setupOpenWeatherMock(t, nil, 0, "", false)
	weatherHandler := setupWeatherHandlerWithCache(cache, metrics, weatherAPIServerMock.URL, openWeatherAPIServerMock.URL)
	router := setupWeatherRouter(weatherHandler)

	requestBody := handlers.GetWeatherRequest{City: city}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	resp := requestWeather(router, body)
	assertWeatherResponse(t, resp, testWeather)

	hits, misses, errors := metrics.Stats()
	assert.Equal(t, 0, hits)
	assert.Equal(t, 1, misses)
	assert.Equal(t, 0, errors)
	assert.Equal(t, 1, cache.Size())

	resp = requestWeather(router, body)
	assertWeatherResponse(t, resp, testWeather)

	hits, misses, errors = metrics.Stats()
	assert.Equal(t, 1, hits)
	assert.Equal(t, 1, misses)
	assert.Equal(t, 0, errors)
	assert.Equal(t, 1, cache.Size())

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
		expectedCode          int
		expectedErrorBody     string
	}{
		{
			name:                  "Invalid City",
			city:                  "123",
			weatherAPIResp:        nil,
			weatherAPICode:        0,
			openWeatherResp:       nil,
			openWeatherCode:       0,
			shouldCallOpenWeather: false,
			expectedCode:          http.StatusBadRequest,
			expectedErrorBody:     `{"error":"invalid request format"}`,
		},
		{
			name: "City Not Found",
			city: "Odeca",
			weatherAPIResp: providers.GetWeatherErrorResponse{
				Error: providers.GetWeatherErrorDetails{
					Code:    weatherAPINotFoundErrorCode,
					Message: "No matching location found.",
				},
			},
			weatherAPICode: http.StatusBadRequest,
			openWeatherResp: providers.GetOpenWeatherErrorResponse{
				Cod:     openWeatherNotFoundErrorCode,
				Message: "city not found",
			},
			openWeatherCode:       http.StatusNotFound,
			shouldCallOpenWeather: true,
			expectedCode:          http.StatusNotFound,
			expectedErrorBody:     `{"error":"there is no city with such name"}`,
		},
		{
			name: "Both Providers Fail",
			city: "TestCity",
			weatherAPIResp: providers.GetWeatherErrorResponse{
				Error: providers.GetWeatherErrorDetails{
					Code:    9999,
					Message: "Internal application error.",
				},
			},
			weatherAPICode: http.StatusBadRequest,
			openWeatherResp: providers.GetOpenWeatherErrorResponse{
				Cod:     "500",
				Message: "internal server error",
			},

			openWeatherCode:       http.StatusInternalServerError,
			shouldCallOpenWeather: true,

			expectedCode:      http.StatusInternalServerError,
			expectedErrorBody: `{"error":"failed to get weather"}`,
		},
		{
			name:                  "Timout Exceeded",
			city:                  "Odeca",
			weatherAPIResp:        nil,
			weatherAPICode:        0,
			openWeatherResp:       nil,
			openWeatherCode:       0,
			shouldCallOpenWeather: false,
			expectedCode:          http.StatusInternalServerError,
			expectedErrorBody:     `{"error":"failed to get weather"}`,
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

			router := setupWeatherRouter(weatherHandler)

			requestBody := handlers.GetWeatherRequest{
				City: testCase.city,
			}
			body, err := json.Marshal(requestBody)
			require.NoError(t, err)

			resp := requestWeather(router, body)

			assert.Equal(t, testCase.expectedCode, resp.Code)
			assert.Equal(t, testCase.expectedErrorBody, resp.Body.String())

		})
	}

}
