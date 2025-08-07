package metrics

import (
	"fmt"
	"time"
	"weather-forecast/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type (
	Prometheus struct {
		requestCount    *prometheus.CounterVec
		requestDuration *prometheus.HistogramVec
		logger          logger.Logger
	}
)

func NewPrometheus(logger logger.Logger) *Prometheus {
	p := &Prometheus{
		requestCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Number of HTTP requests",
			},
			[]string{"path", "method"},
		),
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"path", "method"},
		),
		logger: logger,
	}

	prometheus.MustRegister(p.requestCount, p.requestDuration)

	return p
}

func (m *Prometheus) StartMetricsServer(port string) {
	router := gin.New()
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	addr := fmt.Sprintf(":%s", port)
	m.logger.Infof("Starting gateway metrics server on %s", addr)

	if err := router.Run(addr); err != nil {
		m.logger.Fatalf("Metrics server start: %s", err.Error())
	}
}

func (p *Prometheus) RecordRequest(path, method string, duration time.Duration) {
	p.requestCount.WithLabelValues(path, method).Inc()
	p.requestDuration.WithLabelValues(path, method).Observe(duration.Seconds())
}
