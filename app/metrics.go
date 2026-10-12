package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// testRoute is the only workload route instrumented. Keeping it a literal
// bounds the route label instead of exposing a raw request URL.
const testRoute = "/test"

// latencyBuckets are the fixed seconds buckets for the workload latency
// histogram.
var latencyBuckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5}

// handlerMetrics owns the metric instruments for exactly one serving handler.
// Each handler builds a private prometheus.Registry instead of using the
// package-global default, so building many handlers (for example in tests)
// never hits duplicate-registration panics and no state leaks between them.
type handlerMetrics struct {
	registry *prometheus.Registry
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
}

// newHandlerMetrics creates a private registry and registers the workload
// instruments on it exactly once. It is called once per serving handler, never
// per request.
func newHandlerMetrics() *handlerMetrics {
	registry := prometheus.NewRegistry()

	requests := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total GET /test responses by method, route and HTTP status code.",
	}, []string{"method", "route", "status"})

	duration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "GET /test request duration in seconds by method and route.",
		Buckets: latencyBuckets,
	}, []string{"method", "route"})

	registry.MustRegister(requests, duration)
	return &handlerMetrics{registry: registry, requests: requests, duration: duration}
}

// handler serves this handler's own registry. It is the only reader of the
// registry, so /metrics scrapes never add workload observations.
func (m *handlerMetrics) handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

// instrumentTest wraps the GET /test handler and records one observation per
// completed request. Only real GET requests are recorded: net/http routes HEAD
// to a GET pattern and the mux rejects other methods with 405 before the
// handler runs, so neither is mislabeled as GET. The status label is the status
// the handler actually wrote, which the /test contract bounds to 200/400/500.
func (m *handlerMetrics) instrumentTest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(rec, r)
		elapsed := time.Since(start).Seconds()

		m.requests.WithLabelValues(http.MethodGet, testRoute, strconv.Itoa(rec.status)).Inc()
		m.duration.WithLabelValues(http.MethodGet, testRoute).Observe(elapsed)
	})
}

// statusRecorder observes the status code a handler writes while forwarding
// every write unchanged, so clients see byte-identical responses.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
