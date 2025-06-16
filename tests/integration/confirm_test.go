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

func setupConfirmRouter(handler *handlers.SubscriptionHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/confirm/:token", handler.Confirm)
	return router
}

func TestConfirm_Success(t *testing.T) {
	db := setupDB(t)
	subscriptionHandler := setupSubscriptionHandler(db)
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

}

func TestConfirm_InvalidToken(t *testing.T) {
	db := setupDB(t)
	subscriptionHandler := setupSubscriptionHandler(db)
	router := setupConfirmRouter(subscriptionHandler)
	expectedResponseBody := `{"error":"invalid token"}`
	requestToken := "invalidToken"

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/confirm/"+requestToken, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, expectedResponseBody, w.Body.String())

}

func TestConfirm_TokenNotFound(t *testing.T) {
	db := setupDB(t)
	subscriptionHandler := setupSubscriptionHandler(db)
	router := setupConfirmRouter(subscriptionHandler)
	expectedResponseBody := `{"error":"there is no subscription with such token"}`
	requestToken := "59a29260-39fa-4c9b-845a-4a23bb342e4b"

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/confirm/"+requestToken, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, expectedResponseBody, w.Body.String())

}
