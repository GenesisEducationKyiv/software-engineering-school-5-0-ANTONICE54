package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"weather-forecast/internal/domain/models"
	stub_logger "weather-forecast/internal/infrastructure/logger/stub"
	"weather-forecast/internal/infrastructure/providers"
	"weather-forecast/internal/infrastructure/services"
	"weather-forecast/internal/presentation/server/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TEST_API_KEY = "testAPIKey"

func setupWeatherHandler(testAPIURL, testAPIKey string) *handlers.WeatherHandler {
	stubLogger := stub_logger.New()
	weatherProvider := providers.NewWeatherProvider(testAPIURL, testAPIKey, &http.Client{}, stubLogger)
	weatherService := services.NewWeatherService(weatherProvider, stubLogger)
	weatherHandler := handlers.NewWeatherHandler(weatherService, stubLogger)
	return weatherHandler
}

func setupWeatherRouter(handler *handlers.WeatherHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/weather", handler.Get)
	return router
}

func setupWeatherAPIMock(t *testing.T, responseBody any, statusCode int, city string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/", r.URL.Path)
		assert.Equal(t, fmt.Sprintf("key=%s&q=%s", TEST_API_KEY, city), r.URL.RawQuery)

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(responseBody)

	}))
}

func TestGetWeather_Success(t *testing.T) {

	city := "Kyiv"
	weather := models.Weather{
		Temperature: 22.5,
		Humidity:    64,
		Description: "Partly cloudy",
	}
	responseBody := providers.GetWeatherSuccessResponse{
		Current: providers.GetWeatherCurrentResponse{
			TempC: weather.Temperature,
			Condition: providers.GetWeatherConditionResponse{
				Text: weather.Description,
			},
			Humidity: weather.Humidity,
		},
	}

	weatherAPIServerMock := setupWeatherAPIMock(t, responseBody, http.StatusOK, city)
	defer weatherAPIServerMock.Close()
	weatherHandler := setupWeatherHandler(weatherAPIServerMock.URL, TEST_API_KEY)
	router := setupWeatherRouter(weatherHandler)

	requestBody := handlers.GetWeatherRequest{
		City: city,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/weather", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response handlers.GetWeatherResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, weather.Temperature, response.Temperature)
	assert.Equal(t, weather.Humidity, response.Humidity)
	assert.Equal(t, weather.Description, response.Description)

}

func TestGetWeather_ErrorScenarios(t *testing.T) {
	testTable := []struct {
		name              string
		city              string
		apiResponseBody   any
		apiResponseCode   int
		expectedCode      int
		expectedErrorBody string
	}{
		{
			name:              "Invalid City",
			city:              "123",
			apiResponseBody:   nil,
			apiResponseCode:   0,
			expectedCode:      http.StatusBadRequest,
			expectedErrorBody: `{"error":"invalid request"}`,
		},
		{
			name: "City Not Found",
			city: "Odeca",
			apiResponseBody: providers.GetWeatherErrorResponse{
				Error: providers.GetWeatherErrorDetails{
					Code:    providers.LOCATION_NOT_FOUND,
					Message: "No matching location found.",
				},
			},
			apiResponseCode:   http.StatusBadRequest,
			expectedCode:      http.StatusNotFound,
			expectedErrorBody: `{"error":"there is no city with such name"}`,
		},
		{
			name: "Weather API Error",
			city: "Odeca",
			apiResponseBody: providers.GetWeatherErrorResponse{
				Error: providers.GetWeatherErrorDetails{
					Code:    9999,
					Message: "Internal application error.",
				},
			},
			apiResponseCode:   http.StatusBadRequest,
			expectedCode:      http.StatusInternalServerError,
			expectedErrorBody: `{"error":"failed to get weather"}`,
		},
		{
			name:              "Timout Exceeded",
			city:              "Odeca",
			apiResponseBody:   nil,
			apiResponseCode:   0,
			expectedCode:      http.StatusInternalServerError,
			expectedErrorBody: `{"error":"failed to get weather"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			var weatherHandler *handlers.WeatherHandler
			if testCase.apiResponseBody != nil {
				weatherAPIServerMock := setupWeatherAPIMock(t, testCase.apiResponseBody, testCase.apiResponseCode, testCase.city)
				defer weatherAPIServerMock.Close()
				weatherHandler = setupWeatherHandler(weatherAPIServerMock.URL, TEST_API_KEY)

			} else {
				weatherHandler = setupWeatherHandler("", TEST_API_KEY)
			}

			router := setupWeatherRouter(weatherHandler)

			requestBody := handlers.GetWeatherRequest{
				City: testCase.city,
			}
			body, err := json.Marshal(requestBody)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/weather", bytes.NewBuffer(body))
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedCode, w.Code)
			assert.Equal(t, testCase.expectedErrorBody, w.Body.String())

		})
	}

}
