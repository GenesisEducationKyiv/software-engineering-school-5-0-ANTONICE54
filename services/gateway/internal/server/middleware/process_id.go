package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const processIDKey = "process_id"

func ProcessIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		processID := uuid.New().String()
		c.Set("process_id", processID)
		c.Next()
	}
}
