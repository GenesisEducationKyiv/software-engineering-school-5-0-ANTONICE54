package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"weather-forecast/pkg/proto/weather"
	"weather-service/internal/config"
	"weather-service/internal/domain/models"
	"weather-service/internal/domain/usecases"
	"weather-service/internal/infrastructure/clients/openweather"
	"weather-service/internal/infrastructure/clients/weatherapi"
	"weather-service/internal/infrastructure/providers"

	stub_logger "weather-forecast/pkg/stubs/logger"
	"weather-service/internal/presentation/server/handlers"

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

const (
	testAPIKey = "testAPIKey"
)

type (
	Cacher interface {
		providers.CacheReader
		providers.CacheWriter
	}

	MockServer struct {
		*httptest.Server
		shouldBeCalled bool
		wasCalled      bool
	}
)

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

	cfg := &config.Config{
		WeatherAPIURL:  weatherAPIURLMock,
		WeatherAPIKey:  testAPIKey,
		OpenWeatherURL: openWeatherURLMock,
		OpenWeatherKey: testAPIKey,
	}
	weatherAPIClient := weatherapi.NewClient(cfg, client, stubLogger)
	weatherAPIProvider := providers.NewWeatherAPIProvider(weatherAPIClient, stubLogger)
	weatherAPILink := providers.NewWeatherLink(weatherAPIProvider)

	openWeatherClient := openweather.NewClient(cfg, client, stubLogger)
	openWeatherProvider := providers.NewOpenWeatherProvider(openWeatherClient, stubLogger)
	openWeatherLink := providers.NewWeatherLink(openWeatherProvider)

	weatherAPILink.SetNext(openWeatherLink)

	weatherService := usecases.NewWeatherService(weatherAPILink, stubLogger)
	weatherHandler := handlers.NewWeatherHandler(weatherService, stubLogger)
	return weatherHandler
}

func setupWeatherHandlerWithCache(cacher Cacher, metrics providers.MetricsRecorder, weatherAPIURLMock, openWeatherURLMock string) *handlers.WeatherHandler {
	stubLogger := stub_logger.New()
	client := &http.Client{}
	cfg := &config.Config{
		WeatherAPIURL:  weatherAPIURLMock,
		WeatherAPIKey:  testAPIKey,
		OpenWeatherURL: openWeatherURLMock,
		OpenWeatherKey: testAPIKey,
	}

	weatherAPIClient := weatherapi.NewClient(cfg, client, stubLogger)
	weatherAPIProvider := providers.NewWeatherAPIProvider(weatherAPIClient, stubLogger)
	cacheableWeatherAPIProvider := providers.NewCacheDecorator(weatherAPIProvider, cacher, metrics, stubLogger)
	weatherAPILink := providers.NewWeatherLink(cacheableWeatherAPIProvider)

	openWeatherClient := openweather.NewClient(cfg, client, stubLogger)
	openWeatherProvider := providers.NewOpenWeatherProvider(openWeatherClient, stubLogger)
	cacheableOpenWeatherProvider := providers.NewCacheDecorator(openWeatherProvider, cacher, metrics, stubLogger)
	openWeatherLink := providers.NewWeatherLink(cacheableOpenWeatherProvider)

	cacheProvider := providers.NewCacheWeather(cacher, metrics, stubLogger)
	cacheProviderLink := providers.NewWeatherLink(cacheProvider)

	cacheProviderLink.SetNext(weatherAPILink)
	weatherAPILink.SetNext(openWeatherLink)

	weatherService := usecases.NewWeatherService(cacheProviderLink, stubLogger)
	weatherHandler := handlers.NewWeatherHandler(weatherService, stubLogger)
	return weatherHandler
}

func assertWeatherResponse(t *testing.T, response *weather.GetWeatherResponse, expectedWeather models.Weather) {
	t.Helper()

	assert.Equal(t, expectedWeather.Temperature, response.Temperature)
	assert.Equal(t, expectedWeather.Humidity, int(response.Humidity))
	assert.Equal(t, expectedWeather.Description, response.Description)
}
