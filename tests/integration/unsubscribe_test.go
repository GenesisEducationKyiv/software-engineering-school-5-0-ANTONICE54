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
	db := setupDB(t)
	subscriptionHandler := setupSubscriptionHandler(db)
	router := setupUnsubscribeRouter(subscriptionHandler)

	expectedResponseBody := `{"message":"Unsubscribed successfuly."}`
	unsubscribeSubscription := models.Subscription{
		Email:     "test@gmail.com",
		City:      "Kyiv",
		Frequency: models.Daily,
		Confirmed: true,
		Token:     "59d29860-39fa-4c9b-845a-3e91eab42e4b",
	}
	err := db.Create(&unsubscribeSubscription).Error
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/unsubscribe/"+unsubscribeSubscription.Token, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponseBody, w.Body.String())

	res := db.Where("id = ?", unsubscribeSubscription.ID).Find(&models.Subscription{})
	require.Equal(t, int64(0), res.RowsAffected)

}

func TestUnsubscribe_InvalidToken(t *testing.T) {
	db := setupDB(t)
	subscriptionHandler := setupSubscriptionHandler(db)
	router := setupUnsubscribeRouter(subscriptionHandler)

	expectedResponseBody := `{"error":"invalid token"}`
	requestToken := "invalidToken"

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/unsubscribe/"+requestToken, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, expectedResponseBody, w.Body.String())

}

func TestUnsubscribe_TokenNotFound(t *testing.T) {
	db := setupDB(t)
	subscriptionHandler := setupSubscriptionHandler(db)
	router := setupUnsubscribeRouter(subscriptionHandler)

	expectedResponseBody := `{"error":"there is no subscription with such token"}`
	requestToken := "59d29860-39fa-4c9b-845a-3e91eab42e4b"

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/unsubscribe/"+requestToken, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, expectedResponseBody, w.Body.String())

}
