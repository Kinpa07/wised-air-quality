package metrics

import "github.com/prometheus/client_golang/prometheus"

func NewMetrics() prometheus.Registerer {
	return prometheus.NewRegistry()
}
