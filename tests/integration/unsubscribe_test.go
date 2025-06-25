package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/presentation/server/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUnsubscribeRouter(handler *handlers.SubscriptionHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/unsubscribe/:token", handler.Unsubscribe)
	return router
}

func TestUnsubscribe_Success(t *testing.T) {
	expectedResponseBody := `{"message":"Unsubscribed successfuly."}`
	unsubscribeSubscription := models.Subscription{
		Email:     "test@gmail.com",
		City:      "Kyiv",
		Frequency: models.Daily,
		Confirmed: true,
		Token:     "59d29860-39fa-4c9b-845a-3e91eab42e4b",
	}

	db := setupDB(t)
	subscriptionHandler := setupSubscriptionHandler(db)
	router := setupUnsubscribeRouter(subscriptionHandler)

	err := db.Create(&unsubscribeSubscription).Error
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/unsubscribe/"+unsubscribeSubscription.Token, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponseBody, w.Body.String())

	res := db.Where("id = ?", unsubscribeSubscription.ID).Find(&models.Subscription{})
	require.NoError(t, res.Error)
	require.Equal(t, int64(0), res.RowsAffected)

}

func TestUnsubscribe_ErrorScenarios(t *testing.T) {
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
			subscriptionHandler := setupSubscriptionHandler(db)
			router := setupUnsubscribeRouter(subscriptionHandler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/unsubscribe/"+testCase.requestToken, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
