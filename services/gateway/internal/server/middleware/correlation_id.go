package middleware

import (
	"weather-forecast/pkg/ctxutil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CorrelationIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		correlationID := uuid.New().String()
		ctx.Set(ctxutil.CorrelationIDKey.String(), correlationID)
		ctx.Next()
	}
}
