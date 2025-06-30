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

const (
	testAPIKey = "testAPIKey"
)

var (
	testWeather = models.Weather{
		Temperature: 22.5,
		Humidity:    64,
		Description: "Partly cloudy",
	}
	weatherAPISuccessResponse = providers.GetWeatherSuccessResponse{
		Current: providers.GetWeatherCurrentResponse{
			TempC: testWeather.Temperature,
			Condition: providers.GetWeatherConditionResponse{
				Text: testWeather.Description,
			},
			Humidity: testWeather.Humidity,
		},
	}

	openWeatherSuccessResponse = providers.GetOpenWeatherSuccessResponse{
		Weather: []providers.GetOpenWeatherDescriptionResponse{
			{Description: testWeather.Description},
		},
		Main: providers.GetOpenWeatherMainResponse{
			Temperature: testWeather.Temperature,
			Humidity:    testWeather.Humidity,
		},
	}
)

type MockServer struct {
	*httptest.Server
	shouldBeCalled bool
	wasCalled      bool
}

func newMockServer(t *testing.T, responseBody any, statusCode int, expectedQuery string, shouldBeCalled bool) *MockServer {
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
		if m.shouldBeCalled && !m.wasCalled {
			t.Errorf("Expected mock server to be called, but it wasn't")
		}
		if !m.shouldBeCalled && m.wasCalled {
			t.Errorf("Mock server was called, but shouldn't have been")
		}
		mock.Server.Close()
	})

	return mock
}

func setupWeatherAPIMock(t *testing.T, responseBody any, statusCode int, city string) *MockServer {
	t.Helper()
	expectedQuery := fmt.Sprintf("key=%s&q=%s", TEST_API_KEY, city)
	return newMockServer(t, responseBody, statusCode, expectedQuery, true)
}

func setupOpenWeatherMock(t *testing.T, responseBody any, statusCode int, city string, shouldBeCalled bool) *MockServer {
	t.Helper()
	expectedQuery := ""
	if shouldBeCalled {
		expectedQuery = fmt.Sprintf("q=%s&appid=%s&units=metrics", city, TEST_API_KEY)
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

	weatherService := services.NewWeatherService(weatherAPILink, stubLogger)
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

	city := "Kyiv"

	responseBody := providers.GetWeatherSuccessResponse{
		Current: providers.GetWeatherCurrentResponse{
			TempC: testWeather.Temperature,
			Condition: providers.GetWeatherConditionResponse{
				Text: testWeather.Description,
			},
			Humidity: testWeather.Humidity,
		},
	}

	weatherAPIServerMock := setupWeatherAPIMock(t, responseBody, http.StatusOK, city)
	defer weatherAPIServerMock.Close()

	openWeatherAPIServcerMock := setupOpenWeatherMock(t, nil, 0, "", false)
	defer openWeatherAPIServcerMock.Close()

	weatherHandler := setupWeatherHandler(weatherAPIServerMock.URL, openWeatherAPIServcerMock.URL)
	router := setupWeatherRouter(weatherHandler)

	requestBody := handlers.GetWeatherRequest{
		City: city,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/weather", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	assertWeatherResponse(t, w, testWeather)

}

func TestGetWeather_WeatherAPI_Failed(t *testing.T) {
	city := "Kyiv"

	weatherAPIResponseBody := providers.GetWeatherErrorResponse{
		Error: providers.GetWeatherErrorDetails{
			Code:    9999,
			Message: "Internal application error.",
		},
	}

	openWeatherResponseBody := providers.GetOpenWeatherSuccessResponse{
		Weather: []providers.GetOpenWeatherDescriptionResponse{
			{
				Description: testWeather.Description,
			},
		},
		Main: providers.GetOpenWeatherMainResponse{
			Temperature: testWeather.Temperature,
			Humidity:    testWeather.Humidity,
		},
	}

	weatherAPIServerMock := setupWeatherAPIMock(t, weatherAPIResponseBody, http.StatusBadRequest, city)
	defer weatherAPIServerMock.Close()

	openWeatherAPIServcerMock := setupOpenWeatherMock(t, openWeatherResponseBody, http.StatusOK, city, true)
	defer openWeatherAPIServcerMock.Close()

	weatherHandler := setupWeatherHandler(weatherAPIServerMock.URL, openWeatherAPIServcerMock.URL)
	router := setupWeatherRouter(weatherHandler)

	requestBody := handlers.GetWeatherRequest{
		City: city,
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/weather", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	assertWeatherResponse(t, w, testWeather)

}

func TestGetWeather_ErrorScenarios(t *testing.T) {

	const weatherAPIErrorCode = 1006

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
					Code:    weatherAPIErrorCode,
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

				openWeatherAPIServcerMock := setupOpenWeatherMock(t, nil, 0, "", false)
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

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/weather", bytes.NewBuffer(body))
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedCode, w.Code)
			assert.Equal(t, testCase.expectedErrorBody, w.Body.String())

		})
	}

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
