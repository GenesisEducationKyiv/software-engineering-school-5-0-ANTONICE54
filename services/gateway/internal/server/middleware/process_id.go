package middleware

import (
	"weather-forecast/pkg/ctxutil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ProcessIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		processID := uuid.New().String()
		ctx.Set(ctxutil.ProcessIDKey.String(), processID)
		ctx.Next()
	}
}
