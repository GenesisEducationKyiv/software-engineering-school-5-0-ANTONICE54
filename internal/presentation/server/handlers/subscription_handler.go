package handlers

import (
	"context"
	"errors"
	"net/http"
	domainerr "weather-forecast/internal/domain/errors"
	infraerrors "weather-forecast/internal/infrastructure/errors"

	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
)

type (
	SubsctiptionService interface {
		Subscribe(ctx context.Context, email, frequency, city string) (*models.Subscription, error)
		Confirm(ctx context.Context, token string) error
		Unsubscribe(ctx context.Context, token string) error
	}

	SubscriptionHandler struct {
		subscriptionService SubsctiptionService
		logger              logger.Logger
	}

	SubscribeRequest struct {
		Email     string `json:"email" binding:"required,email"`
		City      string `json:"city" binding:"required,alpha"`
		Frequency string `json:"frequency" binding:"required,oneof=hourly daily"`
	}
)

func NewSubscriptionHandler(subscriptionService SubsctiptionService, logger logger.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		logger:              logger,
	}
}

func (h *SubscriptionHandler) Subscribe(ctx *gin.Context) {
	var req SubscribeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Failed to unmarshal request: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	_, err := h.subscriptionService.Subscribe(ctx, req.Email, req.Frequency, req.City)

	if err != nil {
		h.handleSubscribeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription successful. Confirmation email sent."})

}

func (h *SubscriptionHandler) handleSubscribeError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, domainerr.ErrAlreadySubscribed):
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})

	case errors.Is(err, domainerr.ErrInvalidFrequency):
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	case errors.Is(err, infraerrors.ErrDatabase):
		h.logger.Warnf("Database error during subscription: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})

	case errors.Is(err, infraerrors.ErrInternal):
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	default:
		h.logger.Warnf("Unexpected error during subscription: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

}

func (h *SubscriptionHandler) Confirm(ctx *gin.Context) {
	token := ctx.Param("token")

	err := h.subscriptionService.Confirm(ctx, token)

	if err != nil {
		h.handleConfirmError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription confirmed."})

}

func (h *SubscriptionHandler) handleConfirmError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, domainerr.ErrTokenNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

	case errors.Is(err, infraerrors.ErrDatabase):
		h.logger.Warnf("Database error during subscription: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})

	case errors.Is(err, infraerrors.ErrInternal):
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	case errors.Is(err, domainerr.ErrInvalidToken):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	default:
		h.logger.Warnf("Unexpected error during confirmation: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

}

func (h *SubscriptionHandler) Unsubscribe(ctx *gin.Context) {
	token := ctx.Param("token")

	err := h.subscriptionService.Unsubscribe(ctx, token)

	if err != nil {
		h.handleUnsubscribeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfuly."})

}

func (h *SubscriptionHandler) handleUnsubscribeError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, domainerr.ErrTokenNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

	case errors.Is(err, domainerr.ErrInvalidToken):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	case errors.Is(err, infraerrors.ErrDatabase):
		h.logger.Warnf("Database error during subscription: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})

	case errors.Is(err, infraerrors.ErrInternal):
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	default:
		h.logger.Warnf("Unexpected error during canceling subscription: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

}
