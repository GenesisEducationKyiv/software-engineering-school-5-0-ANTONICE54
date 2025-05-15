package handlers

import (
	"net/http"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/apperrors"
	"weather-forecast/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
)

type (
	SubsctiptionServiceI interface {
		Subscribe(subscription models.Subscription) (*models.Subscription, error)
		Confirm(token string) error
		Unsubscribe(token string) error
	}

	SubscriptionHandler struct {
		subscriptionService SubsctiptionServiceI
		logger              logger.Logger
	}

	SubscribeInput struct {
		Email     string `json:"email" binding:"required,email"`
		City      string `json:"city" binding:"required,alpha"`
		Frequency string `json:"frequency" binding:"required,oneof=hourly daily"`
	}
)

func NewSubscriptionHandler(subscriptionService SubsctiptionServiceI, logger logger.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		logger:              logger,
	}
}

func (h *SubscriptionHandler) Subscribe(ctx *gin.Context) {
	var req SubscribeInput
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Failed to unmarshal request: %s", err.Error())
		ctx.JSON(apperrors.InvalidRequestError.Status(), apperrors.InvalidRequestError.JSONMessage)
		return
	}

	subscription := models.Subscription{
		Email:     req.Email,
		Frequency: models.Frequency(req.Frequency),
		City:      req.City,
	}

	_, err := h.subscriptionService.Subscribe(subscription)

	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			ctx.JSON(appErr.Status(), appErr.JSONMessage)
			return
		}
		ctx.JSON(apperrors.InternalError.Status(), apperrors.InternalError.JSONMessage)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription successful. Confirmation email sent."})

}

func (h *SubscriptionHandler) Confirm(ctx *gin.Context) {
	token := ctx.Param("token")

	err := h.subscriptionService.Confirm(token)

	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			ctx.JSON(appErr.Status(), appErr.JSONMessage)
			return
		}
		ctx.JSON(apperrors.InternalError.Status(), apperrors.InternalError.JSONMessage)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription confirmed."})

}

func (h *SubscriptionHandler) Unsubscribe(ctx *gin.Context) {
	token := ctx.Param("token")

	err := h.subscriptionService.Unsubscribe(token)

	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			ctx.JSON(appErr.Status(), appErr.JSONMessage)
			return
		}
		ctx.JSON(apperrors.InternalError.Status(), apperrors.InternalError.JSONMessage)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfuly."})

}
