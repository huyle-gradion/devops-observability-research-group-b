// Command app is the Assignment 3 group B runtime API.
//
// It is intentionally stdlib-only: the server exposes the consumer-visible
// /health and /test contract while keeping startup configuration separate from
// request data.
package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	listenAddr    = ":8080"
	defaultAppEnv = "local"
	maxDelayMs    = 1000

	healthBody  = "ok\n"
	invalidBody = "invalid test parameters\n"
	failBody    = "intentional test failure\n"
)

func main() {
	id, err := resolveRuntimeID()
	if err != nil {
		slog.Error("failed to resolve runtime identity", "error", err)
		os.Exit(1)
	}
	appEnv := resolveAppEnv()
	slog.Info("runtime API starting", "runtime_id", id, "app_env", appEnv, "addr", listenAddr, "pid", os.Getpid())

	if err := http.ListenAndServe(listenAddr, newHandler(id, appEnv)); err != nil {
		slog.Error("http server stopped", "error", err)
		os.Exit(1)
	}
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
	mux := http.NewServeMux()

	// "GET" patterns match GET and HEAD. Go's ServeMux therefore answers HEAD
	// automatically by running these handlers and suppressing the body; the GET
	// contract below is proven with explicit GET requests, not that behavior.
	// Non-GET methods on the same path are rejected by ServeMux with 405 before
	// the handler runs, so the workload response is never exposed on POST.
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /test", handleTest(id, appEnv))
	return mux
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeText(w, http.StatusOK, healthBody)
}

func handleTest(id, appEnv string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		delayMs := 0
		if query.Has("delay_ms") {
			// ParseUint rejects signs, whitespace, fractions and an explicitly
			// empty value; the range check rejects values above the bound.
			n, err := strconv.ParseUint(query.Get("delay_ms"), 10, 64)
			if err != nil || n > maxDelayMs {
				writeText(w, http.StatusBadRequest, invalidBody)
				return
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
				writeText(w, http.StatusBadRequest, invalidBody)
				return
			}
		}

		// Every parameter is validated before the requested delay is applied.
		if delayMs > 0 {
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		}

		if fail {
			writeText(w, http.StatusInternalServerError, failBody)
			return
		}

		w.Header().Set("X-App-Env", appEnv)
		writeText(w, http.StatusOK, "Hello world from container-"+id+"\n")
	}
}

// writeText writes an exact text/plain response body. Callers set any extra
// headers (for example X-App-Env) before calling.
func writeText(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}
