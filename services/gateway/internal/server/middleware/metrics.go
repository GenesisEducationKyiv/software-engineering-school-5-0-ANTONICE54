package middleware

import (
	"time"

	"weather-forecast/gateway/internal/metrics"

	"github.com/gin-gonic/gin"
)

func MetricsMiddleware(metric metrics.MetricRecorder) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		path := c.FullPath()

		metric.RecordRequest(path, c.Request.Method, duration)
	}
}
