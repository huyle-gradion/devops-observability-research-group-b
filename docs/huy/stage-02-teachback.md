# Stage 02 teachback — runtime API and behavior tests (huyle-gradion)

Post-implementation teaching for the owner. Written 2026-10-11. Companion evidence:
`captures/huy/step-1/api.md`.

## 1. What Stage 02 delivered

A stdlib-only Go service in `app/` implementing Assignment Step 1's app:

| File | Role |
|---|---|
| `app/main.go` | identity resolution, `newHandler(id, appEnv)`, `:8080` server |
| `app/main_test.go` | table-driven `httptest` behavior tests |
| `app/go.mod` | module `…/devops-observability-research-group-b/app`, `go 1.27` |

Contract: `GET /health` → `ok\n`; `GET /test` → `Hello world from container-<id>\n`
with `X-App-Env`; `?delay_ms=0..1000`; `?fail=1` → 500; bad params → 400 (before the
delay); unknown path/aliases → 404. Identity comes from `POD_NAME` or the hostname;
`APP_ENV` defaults to `local`; nothing is read from the environment per request.

## 2. What went wrong / friction (this stage was clean)

- **Functional: nothing.** The worker's code passed vet/tests first time; my independent
  live run matched every case.
- **Harness lesson:** my own `curl -X HEAD` check *hung* — `-X HEAD` makes curl wait for a
  body that never arrives. The correct tool is `curl -I`. A subtle trap that looks like a
  server hang but is a client mistake.
- **Content-Type nuance:** success bodies are exactly `text/plain`; `ServeMux`'s own
  404/405 responses are `text/plain; charset=utf-8`. The tests must assert exact bytes on
  the responses *we* control and not over-constrain the ones we don't.
- **Method routing:** Go 1.22+ `ServeMux` patterns like `"GET /test"` return **405** for a
  `POST` (not 404) and run the handler for `HEAD` automatically. So `POST /test` never
  reaches the workload — verified.
- **`+5` is not a valid unsigned int:** `strconv.ParseUint` rejects a leading sign, so
  `?delay_ms=+5` is a 400 (tested), not silently accepted.

## 3. Lessons / assignment-required knowledge for this stage

- **Exact consumer-visible contract (§3.2.1).** The grader wants a real toy API in a
  public repo. "Working" means the documented bytes/status/headers, not just "returns
  200". Table-driven tests lock the contract so later stages can't silently regress it.
- **Runtime identity (learning gate).** Three different things get confused:
  - **process hostname** = the local machine name (`omarchy` here);
  - **Kubernetes pod name** = `POD_NAME` (a *name*, not an ID);
  - **CRI container ID** = a container runtime identifier.
  We resolve **once at startup** (`POD_NAME` else hostname) and inject it — the pod name
  is not a container ID, and the identity must not change per request.
- **Config vs request data.** `APP_ENV` is read once at startup (default `local`) and
  surfaced as `X-App-Env`; in the cluster it comes from a **ConfigMap** (Stage 07). Keeping
  startup config separate from per-request input is a core design point.
- **`/health` semantics.** `200 ok` proves *this process can respond*. It does **not**
  prove every external dependency (exporter, trace backend, database) is healthy. A deep
  health check would call those; ours deliberately doesn't — say so honestly in the write-up.
- **Validate before acting.** All parameters are validated **before** the delay, so an
  invalid request cannot be used to make the server sleep (a mini denial-of-service
  surface). Same idea later for load/error injection.
- **amd64 vs arm64 (§3.2.1).** The assignment assumes Apple/M1 laptops while GitHub
  Actions builds `linux/amd64`; our workstation is native amd64, so we note that deviation
  and never claim an unperformed arm64 run.

## 4. How to run and inspect it

```sh
cd /home/thearesia/Projects/devops-worktrees/02-runtime-api/app
go test ./...                                   # behavior contract
APP_ENV=api-proof POD_NAME= go run .            # binds :8080, identity = hostname
# in another shell:
curl -i http://127.0.0.1:8080/health
curl -i http://127.0.0.1:8080/test
curl -i 'http://127.0.0.1:8080/test?delay_ms=250&fail=1'
curl -i 'http://127.0.0.1:8080/test?delay_ms=1001'   # 400
curl -I http://127.0.0.1:8080/health                 # HEAD (use -I, not -X HEAD)
```

**What you can do with it:** it is the workload all later stages deploy, containerize and
observe. Change `APP_ENV` at startup and watch the header change; use `delay_ms`/`fail` to
generate real latency/errors for the metrics and load-test stages.

## 5. Fastest rebuild runbook

1. `mkdir app`; `go mod init github.com/huyle-gradion/devops-observability-research-group-b/app`
   (Go 1.27). Stdlib only — no framework.
2. In `main.go`: resolve identity once (`POD_NAME` else `os.Hostname()`, fatal on failure),
   read `APP_ENV` once (default `local`), log `runtime_id`, bind `:8080`.
3. `newHandler(id, appEnv)` with `mux.HandleFunc("GET /health", …)` and
   `"GET /test"`; exact bodies `ok\n`, `Hello world from container-<id>\n`,
   `intentional test failure\n`, `invalid test parameters\n`; set `X-App-Env` on `/test`
   success; **validate params before sleeping**; `ParseUint` + `<=1000`.
4. Write table-driven `httptest` tests for every contract row (exact bodies, content type,
   header, 404/405, boundaries).
5. `gofmt -l .` clean, `go vet ./...` and `go test ./...` green; then run the server and
   `curl` each row (use `curl -I` for HEAD).

**Traps:** using a framework "for convenience"; reading env per request; treating a pod
name as a container ID; testing on a non-200 that it does *not* leak the workload marker;
`curl -X HEAD` hanging; `:8080` occupied (identify the owner, don't kill it).
