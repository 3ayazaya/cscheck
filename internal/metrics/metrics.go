package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type GaugeMetrics struct {
	name    string
	help    string
	Metrics prometheus.Gauge
}

type GaugeVecMetrics struct {
	name    string
	help    string
	Metrics *prometheus.GaugeVec
}

func NewGaugeMetrics(name string, help string) *GaugeMetrics {
	return &GaugeMetrics{
		name: name,
		help: help,
		Metrics: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: name,
				Help: help,
			},
		),
	}
}
func NewGaugeVecMetrics(name string, help string) *GaugeVecMetrics {
	return &GaugeVecMetrics{
		name: name,
		help: help,
		Metrics: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: name,
				Help: help,
			},
			[]string{
				"name",
			},
		),
	}
}
