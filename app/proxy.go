// Bounded same-origin OTLP/HTTP ingestion proxy.
//
// The browser page cannot reach the private Collector directly, so it exports
// traces same-origin to POST /otel/v1/traces and this handler forwards them to
// the Collector configured by OTEL_EXPORTER_OTLP_ENDPOINT. It is deliberately
// narrow: one fixed route, one fixed upstream resolved once, a hard 1 MiB body
// bound enforced before any upstream dispatch, a 5s upstream deadline, and no
// request-controlled target. It adds no retry, authentication or rate limiting,
// so this ingestion path is a public demo surface with no authenticated
// production guarantee.

package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

const (
	// otelProxyRoute is the one same-origin ingestion route this service
	// exposes. It is not a general forward proxy and no other forwarding route
	// exists.
	otelProxyRoute = "/otel/v1/traces"

	// otelProxyMaxBody is the hard 1 MiB inbound bound. It is enforced by
	// reading the request through http.MaxBytesReader before any upstream call,
	// so an over-limit body, including a chunked or otherwise unknown-length
	// one, is rejected with 413 instead of occupying a Collector connection and
	// surfacing as a generic 502.
	otelProxyMaxBody = 1 << 20

	// otelProxyTimeout bounds one upstream exchange (dial, response headers,
	// whole request).
	otelProxyTimeout = 5 * time.Second
)

// otelTraceProxy forwards one bounded OTLP/HTTP traces POST to the fixed
// upstream captured at construction. No request input can change the target.
type otelTraceProxy struct {
	target string
	client *http.Client
}

// resolveOTLPTarget validates OTEL_EXPORTER_OTLP_ENDPOINT (falling back to the
// in-cluster default) and returns the fixed upstream traces URL: the configured
// base with the OTLP/HTTP traces signal path appended. Query, fragment and
// userinfo are cleared so only scheme/host/base-path survive. An absent,
// relative, non-http(s) or host-less endpoint is an error. run calls this at
// startup and refuses to serve on error, so the proxy fails fast rather than
// failing lazily per request.
func resolveOTLPTarget() (*url.URL, error) {
	raw := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if raw == "" {
		raw = defaultOTLPEndpoint
	}

	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("parse OTEL_EXPORTER_OTLP_ENDPOINT %q: %w", raw, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT %q must be an absolute http(s) URL", raw)
	}
	if u.Host == "" {
		return nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT %q must include a host", raw)
	}

	target := *u
	target.Path = path.Join(u.Path, otlpTracesPath)
	target.RawQuery = ""
	target.Fragment = ""
	target.User = nil
	return &target, nil
}

// newOTLPTraceProxy builds the ingestion handler from an already-resolved
// target. A nil target means the endpoint was invalid: run fails startup before
// serving, so this only occurs on a non-serving construction path, where the
// route stays registered and bounded (502) rather than disappearing.
func newOTLPTraceProxy(target *url.URL) http.Handler {
	if target == nil {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			writeText(w, http.StatusBadGateway, "collector unavailable\n")
		})
	}
	return &otelTraceProxy{
		target: target.String(),
		// An explicit transport with no ambient proxy: the upstream is the
		// configured cluster Collector only. Dial and response-header timeouts
		// plus the overall client timeout bound the round trip. No retries.
		client: &http.Client{
			Transport: &http.Transport{
				Proxy:                 nil,
				DialContext:           (&net.Dialer{Timeout: otelProxyTimeout}).DialContext,
				ResponseHeaderTimeout: otelProxyTimeout,
			},
			Timeout: otelProxyTimeout,
		},
	}
}

// ServeHTTP bounds the body, strips credential headers, and forwards the
// validated bytes to the fixed upstream. Every inbound value except
// Content-Type is ignored, so the request can never choose the upstream. The
// real upstream Content-Type, status and body are returned.
func (p *otelTraceProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeText(w, http.StatusMethodNotAllowed, "method not allowed\n")
		return
	}

	// Bound the body before touching the upstream. MaxBytesReader stops the read
	// at the limit for both sized and chunked/unknown-length bodies, so an
	// oversize payload is a 413 here and never reaches the transport.
	r.Body = http.MaxBytesReader(w, r.Body, otelProxyMaxBody)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			writeText(w, http.StatusRequestEntityTooLarge, "request body too large\n")
			return
		}
		writeText(w, http.StatusBadRequest, "cannot read request body\n")
		return
	}

	// Derive the upstream deadline from the caller's context so a client that
	// disconnects cancels the upstream work as well.
	ctx, cancel := context.WithTimeout(r.Context(), otelProxyTimeout)
	defer cancel()

	upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.target, bytes.NewReader(body))
	if err != nil {
		writeText(w, http.StatusInternalServerError, "build upstream request\n")
		return
	}
	// Only the OTLP content type is preserved. Cookie, Authorization and
	// Proxy-Authorization are never forwarded, so credential material cannot
	// ride telemetry traffic to the Collector.
	if contentType := r.Header.Get("Content-Type"); contentType != "" {
		upstreamReq.Header.Set("Content-Type", contentType)
	}

	resp, err := p.client.Do(upstreamReq)
	if err != nil {
		slog.ErrorContext(r.Context(), "otlp trace proxy upstream unavailable", "error", err)
		writeText(w, http.StatusBadGateway, "collector unavailable\n")
		return
	}
	defer resp.Body.Close()

	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}
