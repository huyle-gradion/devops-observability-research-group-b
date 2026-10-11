# Stage 02 — runtime API evidence (huy)

Produced by Sisyphus during Stage 02 (independent of the worker). Observed on
**2026-10-11**. No secrets. Base SHA `7cdd7181c8c918281bffd5c3e68ff26e8485f593`.

## Scope

`app/go.mod`, `app/main.go`, `app/main_test.go` (module
`github.com/huyle-gradion/devops-observability-research-group-b/app`), and a
`Runtime` section in `README.md`. Stdlib only (`net/http`, `httptest`, `slog`).

## Focused checks (Sisyphus re-run)

| Command | Result |
|---|---|
| `gofmt -l .` | empty (clean) |
| `go vet ./...` | exit 0 |
| `go test ./...` | `ok … app` (24 table cases + resolver tests) |

## Actual process run

Launched `APP_ENV=api-proof POD_NAME= go run .` (identity falls back to hostname
`omarchy` for a local non-container process). Observed:

| Request | Status | Body / header |
|---|---|---|
| `GET /health` | 200 | `text/plain`, `ok\n` (len 3) |
| `GET /test` | 200 | `text/plain`, `X-App-Env: api-proof`, `Hello world from container-omarchy\n` (len 35) |
| `GET /test?delay_ms=0` | 200 | success |
| `GET /test?fail=0` | 200 | success |
| `GET /test?fail=1` | 500 | `intentional test failure\n` |
| `GET /test?delay_ms=` | 400 | `invalid test parameters\n` (empty is invalid) |
| `GET /test?delay_ms=+5` | 400 | signed rejected |
| `GET /test?delay_ms=-1` | 400 | negative rejected |
| `GET /test?delay_ms=1001` | 400 | over max rejected |
| `GET /test?fail=` | 400 | empty fail invalid |
| `GET /test?fail=2` | 400 | out of range |
| `GET /healthz` | 404 | legacy alias absent |
| `GET /api/work` | 404 | legacy alias absent |
| `GET /` | 404 | reserved until Stage 11 |
| `POST /test` | 405 | workload not exposed on POST |
| `HEAD /health` | 200 | `Content-Length: 3`, no body (Go auto-HEAD; documented, not relied on) |

Server startup log: `INFO runtime API starting runtime_id=omarchy app_env=api-proof addr=:8080`.

Port `8080` was free before launch and free after cleanup; no unrelated process touched.

## Notes

- Success responses emit exactly `text/plain` (no charset); 404/405 are served by
  `ServeMux` and include `charset=utf-8` (not asserted).
- `/health` proves the process can respond; it does not assert that any external
  exporter or backend is healthy (see the Stage 02 teachback).
