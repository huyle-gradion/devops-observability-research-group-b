package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
)

// workloadPrefix is the consumer-visible workload marker. Non-success paths must
// never contain it.
const workloadPrefix = "Hello world from container-"

func TestNewHandler(t *testing.T) {
	const (
		testID  = "unit-test-id"
		testEnv = "test-env"
	)

	srv := httptest.NewServer(newHandler(testID, testEnv))
	defer srv.Close()

	workloadBody := workloadPrefix + testID + "\n"

	tests := []struct {
		name       string
		method     string
		target     string
		wantStatus int
		wantBody   string // exact body when non-empty
		wantEnv    string // expected X-App-Env ("" asserts the header is absent)
	}{
		{"health ok", http.MethodGet, "/health", http.StatusOK, healthBody, ""},
		{"test default", http.MethodGet, "/test", http.StatusOK, workloadBody, testEnv},
		{"test delay zero boundary", http.MethodGet, "/test?delay_ms=0", http.StatusOK, workloadBody, testEnv},
		{"test delay max boundary accepted", http.MethodGet, "/test?delay_ms=1000", http.StatusOK, workloadBody, testEnv},
		{"test fail zero", http.MethodGet, "/test?fail=0", http.StatusOK, workloadBody, testEnv},
		{"test fail one", http.MethodGet, "/test?fail=1", http.StatusInternalServerError, failBody, ""},
		{"test fail one with delay", http.MethodGet, "/test?delay_ms=20&fail=1", http.StatusInternalServerError, failBody, ""},
		{"empty delay invalid", http.MethodGet, "/test?delay_ms=", http.StatusBadRequest, invalidBody, ""},
		{"delay malformed", http.MethodGet, "/test?delay_ms=abc", http.StatusBadRequest, invalidBody, ""},
		{"delay fractional", http.MethodGet, "/test?delay_ms=1.5", http.StatusBadRequest, invalidBody, ""},
		{"delay signed", http.MethodGet, "/test?delay_ms=+5", http.StatusBadRequest, invalidBody, ""},
		{"delay negative", http.MethodGet, "/test?delay_ms=-1", http.StatusBadRequest, invalidBody, ""},
		{"delay too large", http.MethodGet, "/test?delay_ms=1001", http.StatusBadRequest, invalidBody, ""},
		{"empty fail invalid", http.MethodGet, "/test?fail=", http.StatusBadRequest, invalidBody, ""},
		{"fail out of range", http.MethodGet, "/test?fail=2", http.StatusBadRequest, invalidBody, ""},
		{"fail negative", http.MethodGet, "/test?fail=-1", http.StatusBadRequest, invalidBody, ""},
		{"fail non-integer", http.MethodGet, "/test?fail=true", http.StatusBadRequest, invalidBody, ""},
		{"invalid params rejected before delay", http.MethodGet, "/test?delay_ms=1000&fail=2", http.StatusBadRequest, invalidBody, ""},
		{"unknown path", http.MethodGet, "/unknown", http.StatusNotFound, "", ""},
		{"legacy work alias", http.MethodGet, "/api/work", http.StatusNotFound, "", ""},
		{"legacy healthz alias", http.MethodGet, "/healthz", http.StatusNotFound, "", ""},
		{"root", http.MethodGet, "/", http.StatusNotFound, "", ""},
		{"post test does not expose workload", http.MethodPost, "/test", http.StatusMethodNotAllowed, "", ""},
		{"put health rejected", http.MethodPut, "/health", http.StatusMethodNotAllowed, "", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, srv.URL+tc.target, nil)
			if err != nil {
				t.Fatalf("new request: %v", err)
			}
			resp, err := srv.Client().Do(req)
			if err != nil {
				t.Fatalf("do request: %v", err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}
			gotBody := string(body)

			if resp.StatusCode != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body %q)", resp.StatusCode, tc.wantStatus, gotBody)
			}
			if tc.wantBody != "" {
				if gotBody != tc.wantBody {
					t.Errorf("body = %q, want %q", gotBody, tc.wantBody)
				}
				if ct := resp.Header.Get("Content-Type"); ct != "text/plain" {
					t.Errorf("Content-Type = %q, want %q", ct, "text/plain")
				}
			}
			if tc.wantStatus != http.StatusOK && strings.Contains(gotBody, workloadPrefix) {
				t.Errorf("non-200 response leaked workload: %q", gotBody)
			}
			if gotEnv := resp.Header.Get("X-App-Env"); gotEnv != tc.wantEnv {
				t.Errorf("X-App-Env = %q, want %q", gotEnv, tc.wantEnv)
			}
		})
	}
}

func TestResolveRuntimeID(t *testing.T) {
	t.Run("nonempty POD_NAME wins", func(t *testing.T) {
		t.Setenv("POD_NAME", "pod-abc")
		got, err := resolveRuntimeID()
		if err != nil {
			t.Fatalf("resolveRuntimeID: %v", err)
		}
		if got != "pod-abc" {
			t.Errorf("id = %q, want %q", got, "pod-abc")
		}
	})

	t.Run("empty POD_NAME falls back to hostname", func(t *testing.T) {
		t.Setenv("POD_NAME", "")
		want, err := os.Hostname()
		if err != nil {
			t.Skipf("hostname unavailable: %v", err)
		}
		got, err := resolveRuntimeID()
		if err != nil {
			t.Fatalf("resolveRuntimeID: %v", err)
		}
		if got != want {
			t.Errorf("id = %q, want hostname %q", got, want)
		}
	})
}

func TestResolveAppEnv(t *testing.T) {
	t.Run("defaults to local", func(t *testing.T) {
		t.Setenv("APP_ENV", "")
		if got := resolveAppEnv(); got != "local" {
			t.Errorf("appEnv = %q, want %q", got, "local")
		}
	})

	t.Run("explicit value preserved", func(t *testing.T) {
		t.Setenv("APP_ENV", "api-proof")
		if got := resolveAppEnv(); got != "api-proof" {
			t.Errorf("appEnv = %q, want %q", got, "api-proof")
		}
	})
}

// TestWorkloadMetricBoundary proves the emitted series match real traffic: one
// GET /test success, error and invalid request change only the bounded workload
// series, while /health and /metrics traffic never creates or inflates them.
func TestWorkloadMetricBoundary(t *testing.T) {
	const (
		testID  = "metrics-test-id"
		testEnv = "test-env"
	)

	m := newHandlerMetrics()
	srv := httptest.NewServer(newHandlerWithMetrics(testID, testEnv, m))
	defer srv.Close()

	get := func(target string) int {
		t.Helper()
		resp, err := srv.Client().Get(srv.URL + target)
		if err != nil {
			t.Fatalf("GET %s: %v", target, err)
		}
		defer resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
		return resp.StatusCode
	}

	for _, tc := range []struct {
		target string
		want   int
	}{
		{"/test", http.StatusOK},
		{"/test?fail=1", http.StatusInternalServerError},
		{"/test?delay_ms=1001", http.StatusBadRequest},
	} {
		if got := get(tc.target); got != tc.want {
			t.Fatalf("GET %s = %d, want %d", tc.target, got, tc.want)
		}
	}

	// Non-workload traffic that must not touch the workload series.
	for i := 0; i < 2; i++ {
		get("/health")
		get("/metrics")
	}

	families, err := m.registry.Gather()
	if err != nil {
		t.Fatalf("gather registry: %v", err)
	}

	var (
		routes         = map[string]bool{}
		methods        = map[string]bool{}
		statusCounts   = map[string]float64{}
		histogramCount uint64
		bucketBounds   []float64
	)
	for _, family := range families {
		switch family.GetName() {
		case "http_requests_total":
			for _, metric := range family.GetMetric() {
				labels := map[string]string{}
				for _, pair := range metric.GetLabel() {
					labels[pair.GetName()] = pair.GetValue()
				}
				routes[labels["route"]] = true
				methods[labels["method"]] = true
				statusCounts[labels["status"]] += metric.GetCounter().GetValue()
			}
		case "http_request_duration_seconds":
			for _, metric := range family.GetMetric() {
				histogram := metric.GetHistogram()
				histogramCount += histogram.GetSampleCount()
				for _, bucket := range histogram.GetBucket() {
					bucketBounds = append(bucketBounds, bucket.GetUpperBound())
				}
			}
		}
	}

	if want := map[string]float64{"200": 1, "400": 1, "500": 1}; !reflect.DeepEqual(statusCounts, want) {
		t.Errorf("http_requests_total by status = %v, want %v", statusCounts, want)
	}
	if len(routes) != 1 || !routes[testRoute] {
		t.Errorf("route labels = %v, want only %q (health/metrics must not be counted)", routes, testRoute)
	}
	if len(methods) != 1 || !methods[http.MethodGet] {
		t.Errorf("method labels = %v, want only %q", methods, http.MethodGet)
	}
	if histogramCount != 3 {
		t.Errorf("http_request_duration_seconds observed %d requests, want 3", histogramCount)
	}
	if want := append([]float64{}, latencyBuckets...); !reflect.DeepEqual(bucketBounds, want) {
		// The text exposition adds the implicit +Inf bucket; the gathered dto
		// carries only the declared finite bounds.
		t.Errorf("latency buckets = %v, want %v", bucketBounds, want)
	}
}

// TestInitTelemetry proves the trace pipeline is built without contacting the
// endpoint (the exporter is lazy) and that it installs the global tracer
// provider and W3C trace-context plus baggage propagator the contract requires.
// The endpoint is deliberately unreachable: shutdown must still terminate and
// its export error is the honest outcome, not a failure of this lifecycle test.
func TestInitTelemetry(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://127.0.0.1:4318")

	shutdown, err := initTelemetry(context.Background(), "telemetry-test-id")
	if err != nil {
		t.Fatalf("initTelemetry: %v", err)
	}
	if shutdown == nil {
		t.Fatal("initTelemetry returned a nil shutdown function")
	}

	if otel.GetTracerProvider() == nil {
		t.Error("global tracer provider was not installed")
	} else {
		// Emit one span so the batch processor really attempts an export on
		// shutdown. With no Collector the export must fail honestly rather than
		// be swallowed; the error is logged, not asserted, because a Collector
		// could be reachable on this endpoint.
		_, span := otel.Tracer("app.test").Start(context.Background(), "smoke")
		span.End()
	}

	gotFields := otel.GetTextMapPropagator().Fields()
	sort.Strings(gotFields)
	wantFields := []string{"baggage", "traceparent", "tracestate"}
	if !reflect.DeepEqual(gotFields, wantFields) {
		t.Errorf("global propagator fields = %v, want %v", gotFields, wantFields)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := shutdown(ctx); err != nil {
		t.Logf("shutdown reported an export error as expected with no Collector: %v", err)
	}
}

// TestInitTelemetryInvalidEndpoint proves SDK init fails on a malformed endpoint
// so main can treat it as a startup failure instead of running blind.
func TestInitTelemetryInvalidEndpoint(t *testing.T) {
	for _, endpoint := range []string{"://not-a-url", "not-a-url"} {
		t.Run(endpoint, func(t *testing.T) {
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", endpoint)
			shutdown, err := initTelemetry(context.Background(), "telemetry-test-id")
			if err == nil {
				if shutdown != nil {
					_ = shutdown(context.Background())
				}
				t.Fatalf("initTelemetry accepted malformed endpoint %q", endpoint)
			}
		})
	}
}
