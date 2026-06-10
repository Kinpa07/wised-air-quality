package worker

import "github.com/prometheus/client_golang/prometheus"

const (
	jobsProcessedName         = "jobs_processed_total"
	jobsProcessedTotalLatency = "jobs_processed_duration_milliseconds"
	jobsEnqueuedName          = "jobs_enqueued_total"
)

var (
	dflBuckets = []float64{300, 1200, 5000}
)

var processedJobs *prometheus.CounterVec
var processedJobsLatency *prometheus.HistogramVec
var enqueuedJobs *prometheus.CounterVec

func newWorkerMetrics(name string, registry prometheus.Registerer, buckets ...float64) {
	processedJobs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        jobsProcessedName,
			Help:        "How many jobs were processed, divided by name.",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"name"},
	)
	//Discard error
	registry.Register(processedJobs)

	processedJobsLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        jobsProcessedTotalLatency,
		Help:        "How long it took to process the job, partitioned by name.",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		[]string{"name"},
	)
	//Discard error
	registry.Register(processedJobsLatency)

	enqueuedJobs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        jobsEnqueuedName,
			Help:        "How many jobs were enqueued, divided by name.",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"name"},
	)
	//Discard error
	registry.Register(enqueuedJobs)
}
