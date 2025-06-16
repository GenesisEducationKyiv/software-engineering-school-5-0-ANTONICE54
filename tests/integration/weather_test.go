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

func TestGetWeather_Success(t *testing.T) {
	testAPIKey := "testAPIKey"
	city := "Kyiv"
	weatherInfo := models.Weather{
		Temperature: 22.5,
		Humidity:    64,
		Description: "Partly cloudy",
	}
	weatherAPIServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/", r.URL.Path)
		assert.Equal(t, fmt.Sprintf("key=%s&q=%s", testAPIKey, city), r.URL.RawQuery)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"current": map[string]interface{}{
				"temp_c": weatherInfo.Temperature,
				"condition": map[string]string{
					"text": weatherInfo.Description,
				},
				"humidity": weatherInfo.Humidity,
			}})

	}))
	defer weatherAPIServerMock.Close()
	weatherHandler := setupWeatherHandler(weatherAPIServerMock.URL, testAPIKey)
	router := setupWeatherRouter(weatherHandler)

	requestBody := map[string]string{
		"city": city,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/weather", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response struct {
		Temperature float64 `json:"temperature"`
		Humidity    int     `json:"humidity"`
		Description string  `json:"description"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, weatherInfo.Temperature, response.Temperature)
	assert.Equal(t, weatherInfo.Humidity, response.Humidity)
	assert.Equal(t, weatherInfo.Description, response.Description)

}

func TestGetWeather_InvalidCity(t *testing.T) {
	expectedResponseBody := `{"error":"invalid request"}`

	city := "123"
	weatherHandler := setupWeatherHandler("", "")
	router := setupWeatherRouter(weatherHandler)

	requestBody := map[string]string{
		"city": city,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/weather", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, expectedResponseBody, w.Body.String())

}

func TestGetWeather_CityNotFound(t *testing.T) {
	expectedResponseBody := `{"error":"there is no city with such name"}`

	testAPIKey := "testAPIKey"
	city := "Odeca"
	weatherAPIServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/", r.URL.Path)
		assert.Equal(t, fmt.Sprintf("key=%s&q=%s", testAPIKey, city), r.URL.RawQuery)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    1006,
				"message": "No matching location found.",
			}})

	}))
	defer weatherAPIServerMock.Close()
	weatherHandler := setupWeatherHandler(weatherAPIServerMock.URL, testAPIKey)
	router := setupWeatherRouter(weatherHandler)

	requestBody := map[string]string{
		"city": city,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/weather", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, expectedResponseBody, w.Body.String())

}

func TestGetWeather_WeatherAPIError(t *testing.T) {
	expectedResponseBody := `{"error":"failed to get weather"}`

	testAPIKey := "testAPIKey"
	city := "Odeca"
	weatherAPIServerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/", r.URL.Path)
		assert.Equal(t, fmt.Sprintf("key=%s&q=%s", testAPIKey, city), r.URL.RawQuery)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{
				"code":    9999,
				"message": "Internal application error.",
			}})

	}))
	defer weatherAPIServerMock.Close()
	weatherHandler := setupWeatherHandler(weatherAPIServerMock.URL, testAPIKey)
	router := setupWeatherRouter(weatherHandler)

	requestBody := map[string]string{
		"city": city,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/weather", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedResponseBody, w.Body.String())

}

func TestGetWeather_TimeoutExceeded(t *testing.T) {
	expectedResponseBody := `{"error":"failed to get weather"}`

	city := "Odeca"
	weatherHandler := setupWeatherHandler("", "")
	router := setupWeatherRouter(weatherHandler)

	requestBody := map[string]string{
		"city": city,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/weather", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expectedResponseBody, w.Body.String())

}
