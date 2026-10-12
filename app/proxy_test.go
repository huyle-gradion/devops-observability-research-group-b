package main

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// receiverHit records one upstream request the local httptest receiver saw, so
// tests assert on real forwarded traffic rather than on proxy internals.
type receiverHit struct {
	method      string
	path        string
	query       string
	contentType string
	cookie      string
	authz       string
	proxyAuthz  string
	body        string
}

// newRecordingReceiver stands up a local receiver that records every request it
// handles and replies with the configured status/body. It is the real upstream
// for proxy tests; the proxy under test is pointed at it through the production
// newHandlerWithProxy constructor.
func newRecordingReceiver(t *testing.T, status int, respBody, respContentType string) (*httptest.Server, *[]receiverHit) {
	t.Helper()
	hits := &[]receiverHit{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("receiver read body: %v", err)
		}
		*hits = append(*hits, receiverHit{
			method:      r.Method,
			path:        r.URL.Path,
			query:       r.URL.RawQuery,
			contentType: r.Header.Get("Content-Type"),
			cookie:      r.Header.Get("Cookie"),
			authz:       r.Header.Get("Authorization"),
			proxyAuthz:  r.Header.Get("Proxy-Authorization"),
			body:        string(body),
		})
		if respContentType != "" {
			w.Header().Set("Content-Type", respContentType)
		}
		w.WriteHeader(status)
		_, _ = io.WriteString(w, respBody)
	}))
	t.Cleanup(srv.Close)
	return srv, hits
}

// proxyTargetFor derives the exact target the proxy dials from a receiver URL:
// the receiver base plus the fixed OTLP traces path.
func proxyTargetFor(t *testing.T, base string) *url.URL {
	t.Helper()
	u, err := url.Parse(base)
	if err != nil {
		t.Fatalf("parse receiver URL %q: %v", base, err)
	}
	u.Path = otlpTracesPath
	return u
}

// TestOTLPProxyMethodBoundary proves only POST is accepted on the one registered
// route: a GET is rejected with 405 and never reaches the upstream.
func TestOTLPProxyMethodBoundary(t *testing.T) {
	receiver, hits := newRecordingReceiver(t, http.StatusOK, "accepted", "application/x-protobuf")
	srv := httptest.NewServer(newHandlerWithProxy("proxy-id", "test-env", newHandlerMetrics(), newOTLPTraceProxy(proxyTargetFor(t, receiver.URL))))
	defer srv.Close()

	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			req, err := http.NewRequest(method, srv.URL+otelProxyRoute, nil)
			if err != nil {
				t.Fatalf("new request: %v", err)
			}
			resp, err := srv.Client().Do(req)
			if err != nil {
				t.Fatalf("%s %s: %v", method, otelProxyRoute, err)
			}
			defer resp.Body.Close()
			_, _ = io.Copy(io.Discard, resp.Body)

			if resp.StatusCode != http.StatusMethodNotAllowed {
				t.Fatalf("status = %d, want 405", resp.StatusCode)
			}
			if allow := resp.Header.Get("Allow"); !strings.Contains(allow, http.MethodPost) {
				t.Errorf("Allow = %q, want to include POST", allow)
			}
		})
	}

	if len(*hits) != 0 {
		t.Errorf("upstream received %d requests on rejected methods, want 0", len(*hits))
	}
}

// TestOTLPProxyBodyLimit proves a known-length body over 1 MiB is rejected with
// 413 before any upstream dispatch, and that the limit is a real byte bound:
// exactly 1 MiB is accepted.
func TestOTLPProxyBodyLimit(t *testing.T) {
	receiver, hits := newRecordingReceiver(t, http.StatusOK, "accepted", "application/x-protobuf")
	srv := httptest.NewServer(newHandlerWithProxy("proxy-id", "test-env", newHandlerMetrics(), newOTLPTraceProxy(proxyTargetFor(t, receiver.URL))))
	defer srv.Close()

	oversize := bytes.Repeat([]byte("a"), otelProxyMaxBody+1)
	resp, err := srv.Client().Post(srv.URL+otelProxyRoute, "application/x-protobuf", bytes.NewReader(oversize))
	if err != nil {
		t.Fatalf("oversize POST: %v", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversize status = %d, want 413", resp.StatusCode)
	}
	if len(*hits) != 0 {
		t.Fatalf("upstream received %d requests for an over-limit body, want 0 (limit must precede dispatch)", len(*hits))
	}

	exact := bytes.Repeat([]byte("a"), otelProxyMaxBody)
	resp2, err := srv.Client().Post(srv.URL+otelProxyRoute, "application/x-protobuf", bytes.NewReader(exact))
	if err != nil {
		t.Fatalf("boundary POST: %v", err)
	}
	body2, _ := io.ReadAll(resp2.Body)
	resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("exact-limit status = %d, want 200 (body %q)", resp2.StatusCode, body2)
	}
	if len(*hits) != 1 {
		t.Fatalf("upstream received %d requests at the exact limit, want 1", len(*hits))
	}
	if len((*hits)[0].body) != otelProxyMaxBody {
		t.Errorf("forwarded body = %d bytes, want exactly %d", len((*hits)[0].body), otelProxyMaxBody)
	}
}

// TestOTLPProxyChunkedBodyLimit proves a chunked/unknown-length oversize body is
// rejected with 413 (a real size error), not masked as a generic 502. The body
// is sent with ContentLength -1 so net/http uses chunked transfer encoding.
func TestOTLPProxyChunkedBodyLimit(t *testing.T) {
	receiver, hits := newRecordingReceiver(t, http.StatusOK, "accepted", "application/x-protobuf")
	srv := httptest.NewServer(newHandlerWithProxy("proxy-id", "test-env", newHandlerMetrics(), newOTLPTraceProxy(proxyTargetFor(t, receiver.URL))))
	defer srv.Close()

	oversize := bytes.Repeat([]byte("b"), otelProxyMaxBody+4096)
	req, err := http.NewRequest(http.MethodPost, srv.URL+otelProxyRoute, io.NopCloser(bytes.NewReader(oversize)))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.ContentLength = -1 // force chunked/unknown length
	req.Header.Set("Content-Type", "application/x-protobuf")

	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("chunked oversize POST: %v", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusRequestEntityTooLarge {
		t.Fatalf("chunked oversize status = %d, want 413 (not 502)", resp.StatusCode)
	}
	if len(*hits) != 0 {
		t.Fatalf("upstream received %d requests for a chunked over-limit body, want 0", len(*hits))
	}
}

// TestOTLPProxyUnreachableUpstream proves an unusable upstream yields 502. A
// listener is opened and closed so the target port is a real, refused endpoint.
func TestOTLPProxyUnreachableUpstream(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	addr := ln.Addr().String()
	ln.Close()

	target, err := url.Parse("http://" + addr)
	if err != nil {
		t.Fatalf("parse target: %v", err)
	}
	target.Path = otlpTracesPath

	srv := httptest.NewServer(newHandlerWithProxy("proxy-id", "test-env", newHandlerMetrics(), newOTLPTraceProxy(target)))
	defer srv.Close()

	resp, err := srv.Client().Post(srv.URL+otelProxyRoute, "application/x-protobuf", strings.NewReader(`{"resourceSpans":[]}`))
	if err != nil {
		t.Fatalf("POST to unreachable upstream: %v", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("status = %d, want 502", resp.StatusCode)
	}
}

// TestOTLPProxyForwardsAndStripsCredentials proves the real forward path: the
// body and OTLP Content-Type reach the upstream unchanged, incoming Cookie,
// Authorization and Proxy-Authorization never do, and the upstream's status and
// body are what the client sees.
func TestOTLPProxyForwardsAndStripsCredentials(t *testing.T) {
	const (
		payload    = `{"resourceSpans":[{"scopeSpans":[]}]}`
		upstreamCT = "application/x-protobuf"
	)
	receiver, hits := newRecordingReceiver(t, http.StatusAccepted, "collector-ack", upstreamCT)
	srv := httptest.NewServer(newHandlerWithProxy("proxy-id", "test-env", newHandlerMetrics(), newOTLPTraceProxy(proxyTargetFor(t, receiver.URL))))
	defer srv.Close()

	req, err := http.NewRequest(http.MethodPost, srv.URL+otelProxyRoute, strings.NewReader(payload))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", upstreamCT)
	req.Header.Set("Cookie", "session=super-secret")
	req.Header.Set("Authorization", "Bearer super-secret")
	req.Header.Set("Proxy-Authorization", "Basic super-secret")

	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()
	gotBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("status = %d, want the real upstream 202", resp.StatusCode)
	}
	if string(gotBody) != "collector-ack" {
		t.Errorf("body = %q, want the real upstream body", gotBody)
	}
	if ct := resp.Header.Get("Content-Type"); ct != upstreamCT {
		t.Errorf("response Content-Type = %q, want upstream %q", ct, upstreamCT)
	}

	if len(*hits) != 1 {
		t.Fatalf("upstream received %d requests, want 1", len(*hits))
	}
	hit := (*hits)[0]
	if hit.method != http.MethodPost {
		t.Errorf("upstream method = %q, want POST", hit.method)
	}
	if hit.body != payload {
		t.Errorf("upstream body = %q, want %q", hit.body, payload)
	}
	if hit.contentType != upstreamCT {
		t.Errorf("upstream Content-Type = %q, want preserved %q", hit.contentType, upstreamCT)
	}
	if hit.cookie != "" {
		t.Errorf("Cookie leaked to upstream: %q", hit.cookie)
	}
	if hit.authz != "" {
		t.Errorf("Authorization leaked to upstream: %q", hit.authz)
	}
	if hit.proxyAuthz != "" {
		t.Errorf("Proxy-Authorization leaked to upstream: %q", hit.proxyAuthz)
	}
}

// TestOTLPProxyUpstreamIsFixed proves the upstream is chosen only by the
// configured target, never by request input. Redirect-style headers and a query
// parameter that a generic forward proxy might honor are ignored: the upstream
// still receives exactly the configured /v1/traces path with no query.
func TestOTLPProxyUpstreamIsFixed(t *testing.T) {
	receiver, hits := newRecordingReceiver(t, http.StatusOK, "ok", "application/x-protobuf")
	srv := httptest.NewServer(newHandlerWithProxy("proxy-id", "test-env", newHandlerMetrics(), newOTLPTraceProxy(proxyTargetFor(t, receiver.URL))))
	defer srv.Close()

	req, err := http.NewRequest(http.MethodPost, srv.URL+otelProxyRoute+"?target="+url.QueryEscape(receiver.URL+"/evil")+"&upstream=1.2.3.4", strings.NewReader("payload"))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("Host", "attacker.example")
	req.Header.Set("X-Forwarded-Host", "attacker.example")
	req.Header.Set("X-Original-URL", "/evil")
	req.Host = "attacker.example"

	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if len(*hits) != 1 {
		t.Fatalf("upstream received %d requests, want 1", len(*hits))
	}
	hit := (*hits)[0]
	if hit.path != otlpTracesPath {
		t.Errorf("upstream path = %q, want the fixed %q", hit.path, otlpTracesPath)
	}
	if hit.query != "" {
		t.Errorf("upstream query = %q, want none (input must not shape the upstream)", hit.query)
	}
	if hit.body != "payload" {
		t.Errorf("upstream body = %q, want %q", hit.body, "payload")
	}
}

// TestResolveOTLPTarget proves the startup validation boundary: a valid http(s)
// base resolves to the fixed /v1/traces target, while a relative URL or a
// non-http(s) scheme is an explicit error.
func TestResolveOTLPTarget(t *testing.T) {
	t.Run("valid base resolves the fixed trace path", func(t *testing.T) {
		t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://otel-collector.monitoring.svc.cluster.local:4318")
		got, err := resolveOTLPTarget()
		if err != nil {
			t.Fatalf("resolveOTLPTarget: %v", err)
		}
		if want := "http://otel-collector.monitoring.svc.cluster.local:4318" + otlpTracesPath; got.String() != want {
			t.Errorf("target = %q, want %q", got.String(), want)
		}
	})

	t.Run("valid base with subpath appends the trace path", func(t *testing.T) {
		t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "https://collector.example/otlp")
		got, err := resolveOTLPTarget()
		if err != nil {
			t.Fatalf("resolveOTLPTarget: %v", err)
		}
		if want := "https://collector.example/otlp" + otlpTracesPath; got.String() != want {
			t.Errorf("target = %q, want %q", got.String(), want)
		}
	})

	t.Run("default base is used when unset", func(t *testing.T) {
		t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
		got, err := resolveOTLPTarget()
		if err != nil {
			t.Fatalf("resolveOTLPTarget: %v", err)
		}
		if want := defaultOTLPEndpoint + otlpTracesPath; got.String() != want {
			t.Errorf("target = %q, want %q", got.String(), want)
		}
	})

	for _, endpoint := range []string{
		"://not-a-url",
		"/relative/path",
		"otel-collector:4318",
		"ftp://collector.example:4318",
	} {
		t.Run("rejected: "+endpoint, func(t *testing.T) {
			t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", endpoint)
			if got, err := resolveOTLPTarget(); err == nil {
				t.Fatalf("resolveOTLPTarget accepted %q as %q", endpoint, got)
			}
		})
	}
}
