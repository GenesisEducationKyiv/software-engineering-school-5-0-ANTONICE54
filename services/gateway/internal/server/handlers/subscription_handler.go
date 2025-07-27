package handlers

import (
	"context"
	"net/http"
	"weather-forecast/gateway/internal/errors"
	"weather-forecast/pkg/logger"

	"github.com/gin-gonic/gin"
)

type (
	SubscriptionClient interface {
		Subscribe(ctx context.Context, info SubscribeRequest) error
		Confirm(ctx context.Context, token string) error
		Unsubscribe(ctx context.Context, token string) error
	}

	SubscriptionHandler struct {
		subscriptionClient SubscriptionClient
		logger             logger.Logger
	}

	SubscribeRequest struct {
		Email     string `json:"email" binding:"required,email"`
		City      string `json:"city" binding:"required"`
		Frequency string `json:"frequency" binding:"required,oneof=hourly daily"`
	}
)

func NewSubscriptionHandler(subscriptionClient SubscriptionClient, logger logger.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionClient: subscriptionClient,
		logger:             logger,
	}
}

func (h *SubscriptionHandler) Subscribe(ctx *gin.Context) {
	var req SubscribeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Failed to unmarshal request: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	err := h.subscriptionClient.Subscribe(ctx, req)

	if err != nil {
		httpErr := errors.NewHTTPFromGRPC(err, h.logger)
		ctx.JSON(httpErr.StatusCode, httpErr.Body)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription successful. Confirmation email sent."})

}

func (h *SubscriptionHandler) Confirm(ctx *gin.Context) {
	token := ctx.Param("token")

	err := h.subscriptionClient.Confirm(ctx, token)

	if err != nil {
		httpErr := errors.NewHTTPFromGRPC(err, h.logger)
		ctx.JSON(httpErr.StatusCode, httpErr.Body)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription confirmed."})

}

func (h *SubscriptionHandler) Unsubscribe(ctx *gin.Context) {
	token := ctx.Param("token")

	err := h.subscriptionClient.Unsubscribe(ctx, token)

	if err != nil {
		httpErr := errors.NewHTTPFromGRPC(err, h.logger)
		ctx.JSON(httpErr.StatusCode, httpErr.Body)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfuly."})

}
