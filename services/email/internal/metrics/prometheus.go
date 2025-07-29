package metrics

import (
	"email-service/internal/mappers"
	"fmt"
	"weather-forecast/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type (
	Prometheus struct {
		emailsSentTotal   *prometheus.CounterVec
		emailsFailedTotal *prometheus.CounterVec
		logger            logger.Logger
	}
)

func NewPrometheus(logger logger.Logger) *Prometheus {
	p := &Prometheus{
		emailsSentTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "emails_sent_total",
				Help: "Total number of successfully sent emails",
			},
			[]string{"subject_type"},
		),
		emailsFailedTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "emails_failed_total",
				Help: "Total number of failed email sends",
			},
			[]string{"subject_type"},
		),
		logger: logger,
	}

	prometheus.MustRegister(p.emailsSentTotal, p.emailsFailedTotal)

	return p
}

func (m *Prometheus) StartMetricsServer(port string) {
	router := gin.New()
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	addr := fmt.Sprintf(":%s", port)
	m.logger.Infof("Starting email metrics server on %s", addr)

	if err := router.Run(addr); err != nil {
		m.logger.Fatalf("Metrics server start: %s", err.Error())
	}
}

func (p *Prometheus) RecordEmailSuccess(subject string) {
	subjectType := mappers.SubjectToSubjectType(subject)
	p.emailsSentTotal.WithLabelValues(subjectType).Inc()
}

func (p *Prometheus) RecordEmailFail(subject string) {
	subjectType := mappers.SubjectToSubjectType(subject)
	p.emailsFailedTotal.WithLabelValues(subjectType).Inc()
}
