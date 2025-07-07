package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/domain/usecases"
	stub_logger "weather-forecast/internal/infrastructure/logger/stub"
	"weather-forecast/internal/infrastructure/providers"
	"weather-forecast/internal/presentation/server/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAPIKey = "testAPIKey"
)

type MockServer struct {
	*httptest.Server
	shouldBeCalled bool
	wasCalled      bool
}

func newMockServer(t *testing.T, responseBody interface{}, statusCode int, expectedQuery string, shouldBeCalled bool) *MockServer {
	t.Helper()

	mock := &MockServer{
		shouldBeCalled: shouldBeCalled,
	}

	mock.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mock.wasCalled = true

		if expectedQuery != "" {
			assert.Equal(t, "/", r.URL.Path)
			assert.Equal(t, expectedQuery, r.URL.RawQuery)
		}

		w.WriteHeader(statusCode)
		if responseBody != nil {
			err := json.NewEncoder(w).Encode(responseBody)
			require.NoError(t, err)
		}
	}))

	t.Cleanup(func() {
		if mock.shouldBeCalled && !mock.wasCalled {
			t.Errorf("Expected mock server to be called, but it wasn't")
		}
		if !mock.shouldBeCalled && mock.wasCalled {
			t.Errorf("Mock server was called, but shouldn't have been")
		}
		mock.Close()
	})

	return mock
}

func setupWeatherAPIMock(t *testing.T, responseBody interface{}, statusCode int, city string) *MockServer {
	t.Helper()
	expectedQuery := fmt.Sprintf("key=%s&q=%s", testAPIKey, city)
	return newMockServer(t, responseBody, statusCode, expectedQuery, true)
}

func setupOpenWeatherMock(t *testing.T, responseBody interface{}, statusCode int, city string, shouldBeCalled bool) *MockServer {
	t.Helper()
	expectedQuery := ""
	if shouldBeCalled {
		expectedQuery = fmt.Sprintf("appid=%s&q=%s&units=metric", testAPIKey, city)
	}
	return newMockServer(t, responseBody, statusCode, expectedQuery, shouldBeCalled)
}

func setupWeatherHandler(weatherAPIURLMock, openWeatherURLMock string) *handlers.WeatherHandler {
	stubLogger := stub_logger.New()
	client := &http.Client{}

	weatherAPIProvider := providers.NewWeatherAPIProvider(weatherAPIURLMock, testAPIKey, client, stubLogger)
	weatherAPILink := providers.NewWeatherLink(weatherAPIProvider)

	openWeatherProvider := providers.NewOpenWeatherProvider(openWeatherURLMock, testAPIKey, client, stubLogger)
	openWeatherLink := providers.NewWeatherLink(openWeatherProvider)

	weatherAPILink.SetNext(openWeatherLink)

	weatherService := usecases.NewWeatherService(weatherAPILink, stubLogger)
	weatherHandler := handlers.NewWeatherHandler(weatherService, stubLogger)
	return weatherHandler
}

func setupWeatherRouter(handler *handlers.WeatherHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/weather", handler.Get)
	return router
}

func assertWeatherResponse(t *testing.T, w *httptest.ResponseRecorder, expectedWeather models.Weather) {
	t.Helper()

	assert.Equal(t, http.StatusOK, w.Code)

	var response handlers.GetWeatherResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedWeather.Temperature, response.Temperature)
	assert.Equal(t, expectedWeather.Humidity, response.Humidity)
	assert.Equal(t, expectedWeather.Description, response.Description)
}
