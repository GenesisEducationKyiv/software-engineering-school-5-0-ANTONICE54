package handlers

import (
	"context"
	"net/http"
	"weather-forecast/internal/domain/models"
	infraerrors "weather-forecast/internal/infrastructure/errors"
	"weather-forecast/internal/infrastructure/logger"
	apierrors "weather-forecast/internal/presentation/errors"
	httperrors "weather-forecast/internal/presentation/http_errors"
	"weather-forecast/pkg/apperrors"

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
		httpErr := httperrors.New(apierrors.InvalidRequestError)
		ctx.JSON(httpErr.Status(), httpErr.JSON())
		return
	}

	_, err := h.subscriptionService.Subscribe(ctx, req.Email, req.Frequency, req.City)

	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			httpErr := httperrors.New(appErr)
			ctx.JSON(httpErr.Status(), httpErr.JSON())
			return
		}
		httpErr := httperrors.New(infraerrors.InternalError)
		ctx.JSON(httpErr.Status(), httpErr.JSON())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription successful. Confirmation email sent."})

}

func (h *SubscriptionHandler) Confirm(ctx *gin.Context) {
	token := ctx.Param("token")

	err := h.subscriptionService.Confirm(ctx, token)

	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			httpErr := httperrors.New(appErr)
			ctx.JSON(httpErr.Status(), httpErr.JSON())
			return
		}
		httpErr := httperrors.New(infraerrors.InternalError)
		ctx.JSON(httpErr.Status(), httpErr.JSON())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription confirmed."})

}

func (h *SubscriptionHandler) Unsubscribe(ctx *gin.Context) {
	token := ctx.Param("token")

	err := h.subscriptionService.Unsubscribe(ctx, token)

	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			httpErr := httperrors.New(appErr)
			ctx.JSON(httpErr.Status(), httpErr.JSON())
			return
		}
		httpErr := httperrors.New(infraerrors.InternalError)
		ctx.JSON(httpErr.Status(), httpErr.JSON())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfuly."})

}
