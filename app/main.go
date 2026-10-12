// Command app is the Assignment 3 group B runtime API.
//
// It exposes the consumer-visible /health and /test contract while keeping
// startup configuration separate from request data. Prometheus instrumentation
// is handler-owned (see metrics.go) and never changes a response.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	listenAddr    = ":8080"
	defaultAppEnv = "local"
	maxDelayMs    = 1000

	// serverShutdownTimeout bounds draining in-flight requests so a slow client
	// cannot hang termination.
	serverShutdownTimeout = 5 * time.Second

	// telemetryShutdownTimeout bounds the final span flush so an unreachable or
	// slow Collector cannot hang termination.
	telemetryShutdownTimeout = 5 * time.Second

	healthBody  = "ok\n"
	invalidBody = "invalid test parameters\n"
	failBody    = "intentional test failure\n"

	// testHandlerSpanName is the otelhttp server span for the wrapped /test
	// route. Its span context is what the inner workload logs, so the logged
	// span_id matches the server span a viewer shows for the request.
	testHandlerSpanName = "GET /test"

	// testWorkSpanName is the child span around the real selected delay/work.
	testWorkSpanName = "test-work"
)

func main() {
	os.Exit(run())
}

// run wires telemetry and the HTTP server and returns the process exit code. On
// a stop signal it drains the server first so in-flight spans are complete and
// then flushes the tracer provider with a bounded context. The flush happens
// explicitly before run returns (and before os.Exit), so a deferred shutdown can
// never sit behind os.Exit and silently drop spans.
func run() int {
	// Structured logs are JSON on stdout so the platform can parse them and
	// OTel-context correlation fields (trace_id/span_id) are machine-readable.
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	id, err := resolveRuntimeID()
	if err != nil {
		slog.Error("failed to resolve runtime identity", "error", err)
		return 1
	}
	appEnv := resolveAppEnv()
	slog.Info("runtime API starting", "runtime_id", id, "app_env", appEnv, "addr", listenAddr, "pid", os.Getpid(), "arch", runtime.GOARCH)

	shutdownTelemetry, err := initTelemetry(context.Background(), id)
	if err != nil {
		slog.Error("failed to initialize telemetry", "error", err)
		return 1
	}

	srv := &http.Server{Addr: listenAddr, Handler: newHandler(id, appEnv)}

	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.ListenAndServe() }()

	stopCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	exitCode := 0
	select {
	case err := <-serveErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server stopped", "error", err)
			exitCode = 1
		}
	case <-stopCtx.Done():
		slog.Info("stop signal received")
	}

	// Drain in-flight requests before flushing so their spans are complete.
	httpCtx, cancelHTTP := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancelHTTP()
	if err := srv.Shutdown(httpCtx); err != nil {
		slog.Error("http server shutdown failed", "error", err)
	}

	// Flush and stop the tracer provider. Shutdown and export errors are
	// surfaced honestly on the structured log stream but do not by themselves
	// turn a normal stop into a process failure; export errors are expected when
	// no Collector is reachable (local/CI smoke). Only a genuine serve failure
	// sets exitCode.
	telemetryCtx, cancelTelemetry := context.WithTimeout(context.Background(), telemetryShutdownTimeout)
	defer cancelTelemetry()
	if err := shutdownTelemetry(telemetryCtx); err != nil {
		slog.Error("telemetry shutdown failed", "error", err)
	}

	return exitCode
}

// resolveRuntimeID resolves the runtime identity once: a nonempty POD_NAME
// wins, otherwise the real hostname is used. A hostname lookup failure is fatal.
func resolveRuntimeID() (string, error) {
	if pod := os.Getenv("POD_NAME"); pod != "" {
		return pod, nil
	}
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("resolve hostname: %w", err)
	}
	return hostname, nil
}

// resolveAppEnv reads APP_ENV once, defaulting to local when unset or empty.
func resolveAppEnv() string {
	if env := os.Getenv("APP_ENV"); env != "" {
		return env
	}
	return defaultAppEnv
}

// newHandler builds the router for the resolved runtime identity and app
// environment. Both values are injected explicitly and are never read from the
// environment per request.
func newHandler(id, appEnv string) http.Handler {
	return newHandlerWithMetrics(id, appEnv, newHandlerMetrics())
}

// newHandlerWithMetrics builds the router around an explicit metrics owner. The
// serving path uses newHandler, which creates one handlerMetrics per handler;
// tests may pass their own to inspect the same registry that /metrics serves.
func newHandlerWithMetrics(id, appEnv string, m *handlerMetrics) http.Handler {
	mux := http.NewServeMux()

	// "GET" patterns match GET and HEAD. Go's ServeMux therefore answers HEAD
	// automatically by running these handlers and suppressing the body; the GET
	// contract below is proven with explicit GET requests, not that behavior.
	// Non-GET methods on the same path are rejected by ServeMux with 405 before
	// the handler runs, so the workload response is never exposed on POST.
	//
	// Only GET /test is instrumented with metrics. /health and /metrics stay out
	// of the workload series, and the middleware skips HEAD so it is never
	// recorded as GET.
	//
	// otelhttp wraps only /test, outermost so the metrics middleware and the
	// workload handler both see the server span's request context. /health, the
	// /metrics scrape and the exporter's own traffic therefore create no
	// workload span. An upstream W3C traceparent is honored through the global
	// propagator, which is how the browser stage joins its trace.
	mux.HandleFunc("GET /health", handleHealth)
	mux.Handle("GET /test", otelhttp.NewHandler(m.instrumentTest(handleTest(id, appEnv)), testHandlerSpanName))
	mux.Handle("GET /metrics", m.handler())
	return mux
}

// tracer resolves a tracer from the process-global provider at call time. It is
// deliberately not cached in a package variable: the provider is installed by
// initTelemetry (or a test) after package initialization, and a tracer captured
// earlier would bind the no-op provider and silently drop every span.
func tracer() trace.Tracer {
	return otel.Tracer(serviceName)
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeText(w, http.StatusOK, healthBody)
}

// handleTest serves the bounded workload contract inside the otelhttp server
// span created by the wrapper, so r.Context() carries that span's context.
//
// Correlation: trace_id/span_id are read from the server span context already
// on r.Context() before the child span starts, so the logged span_id matches the
// server span a viewer shows for the request. The child test-work span wraps the
// real selected delay/work and is ended on every path by defer.
//
// Timing: Server-Timing and every other response header are set before the
// single writeText, because a header change after WriteHeader is never sent. The
// measured duration covers real handling only, so a rejected request does not
// report the delay it never ran.
func handleTest(id, appEnv string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverCtx := trace.SpanContextFromContext(r.Context())

		_, span := tracer().Start(r.Context(), testWorkSpanName)
		defer span.End()

		start := time.Now()
		status, body, env := testResponse(r, id, appEnv)
		duration := time.Since(start)

		span.SetAttributes(attribute.Int("http.response.status_code", status))
		if status == http.StatusInternalServerError {
			// The failure is deliberate workload behavior, so record it as a real
			// error with the reason instead of a silent success.
			span.SetStatus(codes.Error, failBody)
		}

		w.Header().Set("Server-Timing", "app;dur="+strconv.FormatInt(duration.Milliseconds(), 10))
		if env != "" {
			w.Header().Set("X-App-Env", env)
		}
		writeText(w, status, body)

		slog.InfoContext(r.Context(), "request completed",
			"trace_id", serverCtx.TraceID().String(),
			"span_id", serverCtx.SpanID().String(),
			"route", testRoute,
			"status", status,
			"duration_ms", duration.Milliseconds(),
			"runtime_id", id,
		)
	}
}

// testResponse applies the bounded /test parameters and returns the response to
// write. Every parameter is validated before the requested delay is applied, so
// an invalid request is rejected immediately and never runs (or reports) the
// selected delay.
func testResponse(r *http.Request, id, appEnv string) (status int, body, env string) {
	query := r.URL.Query()

	delayMs := 0
	if query.Has("delay_ms") {
		// ParseUint rejects signs, whitespace, fractions and an explicitly empty
		// value; the range check rejects values above the bound.
		n, err := strconv.ParseUint(query.Get("delay_ms"), 10, 64)
		if err != nil || n > maxDelayMs {
			return http.StatusBadRequest, invalidBody, ""
		}
		delayMs = int(n)
	}

	fail := false
	if query.Has("fail") {
		// absent, "0" or "1" only; an explicitly empty fail is invalid.
		switch query.Get("fail") {
		case "0":
		case "1":
			fail = true
		default:
			return http.StatusBadRequest, invalidBody, ""
		}
	}

	if delayMs > 0 {
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}

	if fail {
		return http.StatusInternalServerError, failBody, ""
	}
	return http.StatusOK, "Hello world from container-" + id + "\n", appEnv
}

// writeText writes an exact text/plain response body. Callers set any extra
// headers (for example X-App-Env) before calling.
func writeText(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}
