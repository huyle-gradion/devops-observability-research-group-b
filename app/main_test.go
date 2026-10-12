package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
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

// TestStaticPage proves the real page and its generated bundle are served from
// exactly the two routes the browser contract names, with real embedded bytes
// and correct content types, while an unknown path still 404s and no other
// static file is exposed. The bundle is a build input (web/ npm run build), so a
// missing bundle is a compile error for this package, never a placeholder that
// still lets the test pass.
func TestStaticPage(t *testing.T) {
	srv := httptest.NewServer(newHandler("static-test-id", "test-env"))
	defer srv.Close()

	get := func(target string) (*http.Response, string) {
		t.Helper()
		resp, err := srv.Client().Get(srv.URL + target)
		if err != nil {
			t.Fatalf("GET %s: %v", target, err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("read %s: %v", target, err)
		}
		return resp, string(body)
	}

	t.Run("root serves the page", func(t *testing.T) {
		resp, body := get("/")
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want 200", resp.StatusCode)
		}
		if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
			t.Errorf("Content-Type = %q, want text/html", ct)
		}
		for _, want := range []string{
			"<h1", ">Run test<", `id="response"`, `id="trace-id"`, `id="web-vital"`, "/static/app.js",
		} {
			if !strings.Contains(body, want) {
				t.Errorf("page body missing %q", want)
			}
		}
	})

	t.Run("bundle is the exact embedded artifact", func(t *testing.T) {
		resp, body := get("/static/app.js")
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want 200", resp.StatusCode)
		}
		if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/javascript") {
			t.Errorf("Content-Type = %q, want text/javascript", ct)
		}
		want, err := staticFiles.ReadFile(staticBundlePath)
		if err != nil {
			t.Fatalf("read embedded bundle: %v", err)
		}
		if body != string(want) {
			t.Errorf("served bundle (%d bytes) != embedded bundle (%d bytes)", len(body), len(want))
		}
	})

	for _, target := range []string{"/unknown", "/static/index.html", "/static/other.js"} {
		t.Run("not served: "+target, func(t *testing.T) {
			resp, body := get(target)
			if resp.StatusCode != http.StatusNotFound {
				t.Fatalf("status = %d, want 404", resp.StatusCode)
			}
			if strings.Contains(body, workloadPrefix) {
				t.Errorf("404 leaked workload: %q", body)
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

// TestServerTimingHeader proves the measured Server-Timing header is present on
// every /test response and reflects real handling: a delayed request measures at
// least its delay, while a rejected request whose selected delay would be large
// reports a small duration instead of faking the delay it never ran.
func TestServerTimingHeader(t *testing.T) {
	const (
		testID  = "timing-test-id"
		testEnv = "test-env"
	)

	srv := httptest.NewServer(newHandler(testID, testEnv))
	defer srv.Close()

	get := func(target string) *http.Response {
		t.Helper()
		resp, err := srv.Client().Get(srv.URL + target)
		if err != nil {
			t.Fatalf("GET %s: %v", target, err)
		}
		defer resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
		return resp
	}

	t.Run("delay is measured", func(t *testing.T) {
		resp := get("/test?delay_ms=120")
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status = %d, want 200", resp.StatusCode)
		}
		if got := serverTimingMillis(t, resp); got < 120 {
			t.Errorf("Server-Timing dur = %d, want >= 120 (the real delay)", got)
		}
	})

	t.Run("rejected request does not fake the selected delay", func(t *testing.T) {
		resp := get("/test?delay_ms=1001")
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", resp.StatusCode)
		}
		if got := serverTimingMillis(t, resp); got >= 1000 {
			t.Errorf("Server-Timing dur = %d, want < 1000: a rejected request must not run the selected delay", got)
		}
	})

	t.Run("failure is still measured", func(t *testing.T) {
		resp := get("/test?fail=1")
		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("status = %d, want 500", resp.StatusCode)
		}
		_ = serverTimingMillis(t, resp)
		if env := resp.Header.Get("X-App-Env"); env != "" {
			t.Errorf("X-App-Env = %q, want absent on error", env)
		}
	})
}

// serverTimingMillis parses the app duration from the Server-Timing header,
// failing the test when the header is missing or malformed.
func serverTimingMillis(t *testing.T, resp *http.Response) int64 {
	t.Helper()
	const prefix = "app;dur="
	value := resp.Header.Get("Server-Timing")
	if !strings.HasPrefix(value, prefix) {
		t.Fatalf("Server-Timing = %q, want prefix %q", value, prefix)
	}
	ms, err := strconv.ParseInt(strings.TrimPrefix(value, prefix), 10, 64)
	if err != nil {
		t.Fatalf("Server-Timing duration %q: %v", value, err)
	}
	return ms
}

// TestTestRouteSpans proves /test runs inside the otelhttp server span, creates
// the test-work child span under it, honors an upstream W3C traceparent, and
// creates no spans for /health or the /metrics scrape.
func TestTestRouteSpans(t *testing.T) {
	const (
		testID  = "spans-test-id"
		testEnv = "test-env"
		// Canonical W3C trace-context example values.
		incomingTraceID = "0af7651916cd43dd8448eb211c80319c"
		incomingSpanID  = "b7ad6b7169203331"
	)

	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	prevProvider := otel.GetTracerProvider()
	prevPropagator := otel.GetTextMapPropagator()
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	t.Cleanup(func() {
		otel.SetTracerProvider(prevProvider)
		otel.SetTextMapPropagator(prevPropagator)
		_ = provider.Shutdown(context.Background())
	})

	srv := httptest.NewServer(newHandler(testID, testEnv))
	defer srv.Close()

	do := func(target, traceparent string) {
		t.Helper()
		req, err := http.NewRequest(http.MethodGet, srv.URL+target, nil)
		if err != nil {
			t.Fatalf("new request: %v", err)
		}
		if traceparent != "" {
			req.Header.Set("traceparent", traceparent)
		}
		resp, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("GET %s: %v", target, err)
		}
		defer resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	}

	do("/test?delay_ms=5", "00-"+incomingTraceID+"-"+incomingSpanID+"-01")
	do("/health", "")
	do("/metrics", "")

	ended := recorder.Ended()
	if len(ended) != 2 {
		t.Fatalf("ended spans = %d, want 2 (server + test-work): only /test may be traced", len(ended))
	}

	server := findSpan(t, ended, testHandlerSpanName)
	work := findSpan(t, ended, testWorkSpanName)

	if got := server.SpanContext().TraceID().String(); got != incomingTraceID {
		t.Errorf("server trace_id = %s, want upstream %s", got, incomingTraceID)
	}
	if parent := server.Parent(); !parent.IsRemote() || parent.SpanID().String() != incomingSpanID {
		t.Errorf("server parent span_id = %s (remote %t), want remote upstream %s", parent.SpanID(), parent.IsRemote(), incomingSpanID)
	}
	if work.Parent().SpanID() != server.SpanContext().SpanID() {
		t.Errorf("test-work parent span_id = %s, want server %s", work.Parent().SpanID(), server.SpanContext().SpanID())
	}
	if work.SpanContext().TraceID() != server.SpanContext().TraceID() {
		t.Errorf("test-work trace_id = %s, want server %s", work.SpanContext().TraceID(), server.SpanContext().TraceID())
	}
}

// findSpan returns the single ended span with the given name, failing otherwise.
func findSpan(t *testing.T, spans []sdktrace.ReadOnlySpan, name string) sdktrace.ReadOnlySpan {
	t.Helper()
	for _, span := range spans {
		if span.Name() == name {
			return span
		}
	}
	t.Fatalf("span %q not found among %d ended spans", name, len(spans))
	return nil
}
