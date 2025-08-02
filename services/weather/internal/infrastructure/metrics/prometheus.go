package metrics

import (
	"fmt"
	"weather-forecast/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type (
	Prometheus struct {
		cacheHits   prometheus.Counter
		cacheMisses prometheus.Counter
		cacheErrors prometheus.Counter
		logger      logger.Logger
	}
)

func NewPrometheus(logger logger.Logger) *Prometheus {

	metricManager := &Prometheus{
		cacheHits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "weather_cache_hits_total",
			Help: "Total number of cache hits",
		}),
		cacheMisses: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "weather_cache_misses_total",
			Help: "Total number of cache misses",
		}),
		cacheErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "weather_cache_errors_total",
			Help: "Total number of cache errors",
		}),
		logger: logger,
	}

	prometheus.MustRegister(
		metricManager.cacheHits,
		metricManager.cacheMisses,
		metricManager.cacheErrors,
	)

	return metricManager

}

func (m *Prometheus) StartMetricsServer(port string) {
	router := gin.New()
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	addr := fmt.Sprintf(":%s", port)
	m.logger.Infof("Starting weather metrics server on %s", addr)

	if err := router.Run(addr); err != nil {
		m.logger.Fatalf("Metrics server start: %s", err.Error())
	}
}

func (m *Prometheus) RecordCacheHit() {
	m.cacheHits.Inc()
}

func (m *Prometheus) RecordCacheMiss() {
	m.cacheMisses.Inc()
}

func (m *Prometheus) RecordCacheError() {
	m.cacheErrors.Inc()
}
