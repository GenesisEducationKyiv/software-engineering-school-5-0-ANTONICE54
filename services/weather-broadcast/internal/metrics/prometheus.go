package metrics

import (
	"fmt"
	"time"
	"weather-forecast/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prometheus struct {
	broadcastDuration *prometheus.HistogramVec
	logger            logger.Logger
}

func NewPrometheus(logger logger.Logger) *Prometheus {
	p := &Prometheus{
		broadcastDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "weather_broadcast_duration_seconds",
				Help:    "Time taken to complete weather broadcast",
				Buckets: []float64{1, 5, 10, 30, 60, 120, 300, 600},
			},
			[]string{"frequency"},
		),
		logger: logger,
	}

	prometheus.MustRegister(p.broadcastDuration)

	return p
}

func (m *Prometheus) StartMetricsServer(port string) {
	router := gin.New()
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	addr := fmt.Sprintf(":%s", port)
	m.logger.Infof("Starting weather-broadcast metrics server on %s", addr)

	if err := router.Run(addr); err != nil {
		m.logger.Fatalf("Metrics server start: %s", err.Error())
	}
}

func (p *Prometheus) RecordBroadcastDuration(frequency string, duration time.Duration) {
	p.broadcastDuration.WithLabelValues(frequency).Observe(duration.Seconds())
}
