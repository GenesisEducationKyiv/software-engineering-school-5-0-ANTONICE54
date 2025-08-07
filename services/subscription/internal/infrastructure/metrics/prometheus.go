package metrics

import (
	"fmt"
	"weather-forecast/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Prometheus struct {
	subscriptionsCreated   prometheus.Counter
	subscriptionsConfirmed prometheus.Counter
	subscriptionsDeleted   prometheus.Counter
	logger                 logger.Logger
}

func NewPrometheus(logger logger.Logger) *Prometheus {
	p := &Prometheus{
		subscriptionsCreated: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "subscriptions_created_total",
			Help: "Total number of successfully created subscriptions",
		}),
		subscriptionsConfirmed: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "subscriptions_confirmed_total",
			Help: "Total number of confirmed subscriptions",
		}),
		subscriptionsDeleted: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "subscriptions_deleted_total",
			Help: "Total number of deleted subscriptions (unsubscribes)",
		}),
		logger: logger,
	}

	prometheus.MustRegister(p.subscriptionsCreated, p.subscriptionsConfirmed, p.subscriptionsDeleted)

	return p
}

func (p *Prometheus) StartMetricsServer(port string) {
	router := gin.New()
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	addr := fmt.Sprintf(":%s", port)
	p.logger.Infof("Starting subscription metrics server on %s", addr)

	if err := router.Run(addr); err != nil {
		p.logger.Fatalf("Metrics server start: %s", err.Error())
	}
}

func (p *Prometheus) RecordSubscriptionCreated() {
	p.subscriptionsCreated.Inc()
}

func (p *Prometheus) RecordSubscriptionConfirmed() {
	p.subscriptionsConfirmed.Inc()
}

func (p *Prometheus) RecordSubscriptionDeleted() {
	p.subscriptionsDeleted.Inc()
}
