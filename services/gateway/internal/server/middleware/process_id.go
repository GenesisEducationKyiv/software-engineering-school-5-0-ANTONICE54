package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const processIDKey = "process_id"

func ProcessIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		processID := uuid.New().String()
		ctx.Set("process_id", processID)
		ctx.Next()
	}
}
