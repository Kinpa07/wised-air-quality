package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	dflBuckets = []float64{300, 1200, 5000}
)

var metricsRegistered bool

const (
	patternReqsName     = "http_requests_pattern_total"
	patternLatencyName  = "http_requests_pattern_duration_milliseconds"
	patternInFlightName = "http_requests_pattern_in_flight"
)

// Middleware is a handler that exposes prometheus metrics for the number of requests,
// the latency and the response size, partitioned by status code, method and HTTP path.
type Middleware struct {
	reqs     *prometheus.CounterVec
	latency  *prometheus.HistogramVec
	inFlight *prometheus.GaugeVec
}

// NewPrometheusPatternMiddleware returns a new prometheus Middleware handler that groups requests by the chi routing pattern.
// EX: /users/{firstName} instead of /users/bob
func NewPrometheusPatternMiddleware(name string, registry prometheus.Registerer, buckets ...float64) func(next http.Handler) http.Handler {
	var m Middleware
	defer func() {
		metricsRegistered = true
	}()

	if !metricsRegistered {
		m.reqs = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        patternReqsName,
				Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path (with patterns).",
				ConstLabels: prometheus.Labels{"service": name},
			},
			[]string{"code", "method", "path"},
		)
		registry.MustRegister(m.reqs)

		if len(buckets) == 0 {
			buckets = dflBuckets
		}
		m.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        patternLatencyName,
			Help:        "How long it took to process the request, partitioned by status code, method and HTTP path (with patterns).",
			ConstLabels: prometheus.Labels{"service": name},
			Buckets:     buckets,
		},
			[]string{"code", "method", "path"},
		)
		registry.MustRegister(m.latency)

		m.inFlight = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name:        patternInFlightName,
			Help:        "How many requests service is actually serving, partitioned by method.",
			ConstLabels: prometheus.Labels{"service": name},
		},
			[]string{"method"},
		)
		registry.MustRegister(m.inFlight)
	}
	return m.patternHandler
}

func (c Middleware) patternHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		c.inFlight.WithLabelValues(r.Method).Inc()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		rctx := chi.RouteContext(r.Context())
		routePattern := strings.Join(rctx.RoutePatterns, "")
		routePattern = strings.Replace(routePattern, "/*/", "/", -1)

		c.inFlight.WithLabelValues(r.Method).Dec()
		c.reqs.WithLabelValues(http.StatusText(ww.Status()), r.Method, routePattern).Inc()
		c.latency.WithLabelValues(http.StatusText(ww.Status()), r.Method, routePattern).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	}
	return http.HandlerFunc(fn)
}

func InitGoRuntimeMetrics(registry prometheus.Registerer) {
	if !metricsRegistered {
		registry.MustRegister(collectors.NewBuildInfoCollector())
		registry.MustRegister(collectors.NewGoCollector(
			collectors.WithGoCollections(collectors.GoRuntimeMetricsCollection),
		))
	}
}
