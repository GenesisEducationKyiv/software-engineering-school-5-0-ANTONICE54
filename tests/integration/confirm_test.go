package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"weather-forecast/internal/domain/models"
	stub_services "weather-forecast/internal/infrastructure/services/stubs"
	"weather-forecast/internal/presentation/server/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupConfirmRouter(handler *handlers.SubscriptionHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/confirm/:token", handler.Confirm)
	return router
}

func TestConfirm_Success(t *testing.T) {
	db := setupDB(t)
	stubMailer := stub_services.NewStubMailer()
	subscriptionHandler := setupSubscriptionHandler(db, stubMailer)
	router := setupConfirmRouter(subscriptionHandler)
	expectedResponseBody := `{"message":"Subscription confirmed."}`
	unconfirmedSubscription := models.Subscription{
		Email:     "test@gmail.com",
		City:      "Kyiv",
		Frequency: models.Daily,
		Confirmed: false,
		Token:     "59d29860-39fa-4c9b-845a-3e91eab42e4b",
	}
	err := db.Create(&unconfirmedSubscription).Error
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/confirm/"+unconfirmedSubscription.Token, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponseBody, w.Body.String())

	var confirmedSubscription models.Subscription

	err = db.Where("id = ?", unconfirmedSubscription.ID).First(&confirmedSubscription).Error
	require.NoError(t, err)
	assert.True(t, confirmedSubscription.Confirmed)
	assert.Len(t, stubMailer.SentConfirmeds, 1)
	assert.Equal(t, unconfirmedSubscription.Email, stubMailer.SentConfirmeds[0].Email)
	assert.EqualValues(t, unconfirmedSubscription.Frequency, stubMailer.SentConfirmeds[0].Frequency)

}

func TestConfirm_ErrorScenarios(t *testing.T) {
	testTable := []struct {
		name                 string
		requestToken         string
		expectedResponseBody string
		expectedCode         int
	}{
		{
			name:                 "Invalid Token",
			requestToken:         "invalidToken",
			expectedResponseBody: `{"error":"invalid token"}`,
			expectedCode:         http.StatusBadRequest,
		},
		{
			name:                 "Token Not Found",
			requestToken:         "59d29860-39fa-4c9b-845a-3e91eab42e4b",
			expectedResponseBody: `{"error":"there is no subscription with such token"}`,
			expectedCode:         http.StatusNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			db := setupDB(t)
			stubMailer := stub_services.NewStubMailer()
			subscriptionHandler := setupSubscriptionHandler(db, stubMailer)
			router := setupConfirmRouter(subscriptionHandler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/confirm/"+testCase.requestToken, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
