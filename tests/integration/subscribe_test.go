package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/domain/usecases"
	"weather-forecast/internal/infrastructure/database"
	stub_logger "weather-forecast/internal/infrastructure/logger/stub"
	"weather-forecast/internal/infrastructure/repositories"
	"weather-forecast/internal/infrastructure/services"
	stub_services "weather-forecast/internal/infrastructure/services/stubs"
	"weather-forecast/internal/infrastructure/token"
	"weather-forecast/internal/presentation/server/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	require.NoError(t, err)

	database.RunMigration(db)

	return db
}

func setupSubscriptionHandler(db *gorm.DB, stubMailer services.NotificationServiceI) *handlers.SubscriptionHandler {

	stubLogger := stub_logger.New()
	tokenManager := token.NewUUIDManager()

	subscRepo := repositories.NewSubscriptionRepository(db, stubLogger)
	subscUC := usecases.NewSubscriptionUseCase(subscRepo, stubLogger)
	subscService := services.NewSubscriptionService(subscUC, tokenManager, stubMailer, stubLogger)
	subscHandler := handlers.NewSubscriptionHandler(subscService, stubLogger)

	return subscHandler
}

func setupSubscribeRouter(handler *handlers.SubscriptionHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/subscribe", handler.Subscribe)
	return router
}

func TestSubscribe_Success(t *testing.T) {
	db := setupDB(t)
	stubMailer := stub_services.NewStubMailer()
	subscriptionHandler := setupSubscriptionHandler(db, stubMailer)
	router := setupSubscribeRouter(subscriptionHandler)

	requestBody := handlers.SubscribeRequest{
		Email:     "test@gmail.com",
		City:      "Kyiv",
		Frequency: "daily",
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)

	expectedResponseBody := `{"message":"Subscription successful. Confirmation email sent."}`
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponseBody, w.Body.String())

	subscFromDB := models.Subscription{}
	err = db.Where("email = ?", requestBody.Email).First(&subscFromDB).Error
	require.NoError(t, err)
	assert.Equal(t, requestBody.City, subscFromDB.City)
	assert.Equal(t, requestBody.Frequency, string(subscFromDB.Frequency))
	assert.False(t, subscFromDB.Confirmed)
	assert.Equal(t, requestBody.Email, subscFromDB.Email)
	assert.Len(t, stubMailer.SentConfirmations, 1)
	assert.Equal(t, requestBody.Email, stubMailer.SentConfirmations[0].Email)
	assert.EqualValues(t, requestBody.Frequency, stubMailer.SentConfirmations[0].Frequency)

}

func TestSubscribe_AlreadySubscribed(t *testing.T) {
	db := setupDB(t)
	stubMailer := stub_services.NewStubMailer()
	subscriptionHandler := setupSubscriptionHandler(db, stubMailer)
	router := setupSubscribeRouter(subscriptionHandler)

	requestBody := handlers.SubscribeRequest{
		Email:     "test@gmail.com",
		City:      "Kyiv",
		Frequency: "daily",
	}
	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(body))
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(body))
	router.ServeHTTP(w2, req2)
	expectedResponseBody := `{"error":"email already subscribed"}`
	assert.Equal(t, http.StatusConflict, w2.Code)
	assert.Equal(t, expectedResponseBody, w2.Body.String())
	assert.Len(t, stubMailer.SentConfirmations, 1)
	assert.Equal(t, requestBody.Email, stubMailer.SentConfirmations[0].Email)
	assert.EqualValues(t, requestBody.Frequency, stubMailer.SentConfirmations[0].Frequency)

}

func TestSubscribe_InvalidInput(t *testing.T) {
	db := setupDB(t)
	stubMailer := stub_services.NewStubMailer()
	subscriptionHandler := setupSubscriptionHandler(db, stubMailer)
	router := setupSubscribeRouter(subscriptionHandler)
	errorMessage := `{"error":"invalid request"}`
	testTable := []struct {
		name                string
		request             handlers.SubscribeRequest
		expecttedStatusCode int
		expectedMessage     string
	}{
		{
			name: "Invalid email",
			request: handlers.SubscribeRequest{
				Email:     "abc",
				City:      "Lviv",
				Frequency: "daily",
			},
			expecttedStatusCode: http.StatusBadRequest,
			expectedMessage:     errorMessage,
		},
		{
			name: "Invalid city",
			request: handlers.SubscribeRequest{
				Email:     "test@gmail.com",
				City:      "",
				Frequency: "daily",
			},
			expecttedStatusCode: http.StatusBadRequest,
			expectedMessage:     errorMessage,
		},
		{
			name: "Invalid frequency",
			request: handlers.SubscribeRequest{
				Email:     "test@gmail.com",
				City:      "Lviv",
				Frequency: "yearly",
			},
			expecttedStatusCode: http.StatusBadRequest,
			expectedMessage:     errorMessage,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			body, err := json.Marshal(testCase.request)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/subscribe", bytes.NewBuffer(body))
			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.expecttedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedMessage, w.Body.String())
		})
	}
}
