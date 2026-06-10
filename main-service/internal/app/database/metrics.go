package database

import "github.com/prometheus/client_golang/prometheus"

var (
	dflBuckets = []float64{10, 100, 500, 1000, 2500, 5000}
)

const (
	patternQueryDurations = "db_query_duration_milliseconds"
)

type MetricsCollector struct {
	queryDuration *prometheus.HistogramVec
}

var metricsCollector MetricsCollector

func NewDBPrometheusMetrics(registry *prometheus.Registry, serviceName string) {
	metricsCollector.queryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        patternQueryDurations,
			Help:        "A histogram of the SQL query durations, partitioned by clause and status.",
			ConstLabels: prometheus.Labels{"service": serviceName},
			Buckets:     dflBuckets,
		},
		[]string{"clause", "success"},
	)
	registry.MustRegister(metricsCollector.queryDuration)
}

// getClause obtains clause from SQL query.
func getClause(sql string) string {
	if len(sql) > 5 {
		return sql[:6]
	}
	return ""
}
