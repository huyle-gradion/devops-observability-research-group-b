package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
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
