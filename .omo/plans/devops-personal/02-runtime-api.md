# Stage 02: Runtime API and consumer-visible behavior

## Governing artifacts

[Spec](../../../.scratch/assignment-3-personal/spec.md), [ticket 02](../../../.scratch/assignment-3-personal/issues/02-runtime-api.md), assignment Step 1, overview Approach 3, [execution contract](00-execution-contract.md).

## Entry gate and routing

Stage 01 accepted; working Go installation recorded. Preliminary routine-verified: bounded standard-library API with deterministic behavior tests. Escalate if existing code/contracts conflict. Research Go net/http and current ServeMux method/path behavior before editing; do not add a framework.

## Allowed scope and branch

Base origin/personal/huy at Stage 01 accepted SHA; branch stage/huy/02-runtime-api; worktree PROJECTS_ROOT/devops-worktrees/02-runtime-api. Files app/go.mod, app/main.go, app/main_test.go, README runtime instructions, docs/huy/lessons.md and own Step 1 API evidence. Commit scope: API and behavior tests/docs only. No Docker, CI, telemetry, Kubernetes, frontend or fake credential. Preserve existing content if source appeared since planning.

## Bounded substeps

### 02A: Runtime resolution and handler

Use Go 1.27 and only net/http, httptest, slog and required stdlib utilities. Module path is derived from actual REPO and recorded, not a fabricated example. Define newHandler(id, appEnv string) http.Handler. main resolves nonempty POD_NAME or os.Hostname once and fails on hostname/startup failure, reads APP_ENV once with local default, logs initial runtime_id and binds :8080. Fixed port suffices; no configuration framework. Handler receives resolved values explicitly.

| Request | Required response |
|---|---|
| GET /health | 200, text/plain, exact ok newline |
| GET /test | 200, text/plain, exact Hello world from container-ID newline, X-App-Env resolved value |
| GET /test?delay_ms=N | Same success after integer N milliseconds, allowed 0–1000, absent means 0 |
| GET /test?fail=1 | 500, exact intentional test failure newline, after selected delay |
| Invalid delay or fail parameter | 400, exact invalid test parameters newline; validate before delaying |
| Unknown path | 404, including legacy /api/work and /healthz |

fail accepts absent, 0 or 1; an explicitly empty fail is invalid. An explicitly empty delay_ms is invalid; default applies only when absent. Use strict integer parsing, reject out-of-range values; no randomness or environment-baked fake ID. HTTP method routing must not expose workload on POST accidentally. State Go's automatic HEAD behavior rather than relying on it for the GET proof. `/` remains 404 until Stage 11.

Exit gate: code builds, normal/invalid/error branches follow the same explicit contract, and startup configuration is separate from request data.

### 02B: Behavior tests and actual process run

Table-driven httptest checks cover exact health/test bodies, content type, injected identity, APP_ENV, unknown path, non-GET workload rejection, valid fail=0/1, invalid fail, malformed/negative/>1000 delay, accepted zero/max boundary. Avoid wall-clock assertions; one accepted max-delay test can wait its bounded duration without asserting scheduler timing. Do not test source text, log wording or incidental implementation defaults. Permanent tests must catch consumer-visible regressions.

```sh
(cd app && go vet ./... && go test ./...)
(cd app && APP_ENV=api-proof POD_NAME= go run .)
```

In a separate terminal:

```sh
curl -i http://127.0.0.1:8080/health
curl -i http://127.0.0.1:8080/test
curl -i 'http://127.0.0.1:8080/test?delay_ms=100&fail=1'
curl -i 'http://127.0.0.1:8080/test?delay_ms=1001'
curl -i http://127.0.0.1:8080/unknown
```

Compare the real test ID with hostname for this non-container local process and header with api-proof. POD_NAME is cleared for this exercise so inherited shell values cannot masquerade as runtime identity. Stop only your launched process after capture. No unavailable dependency is silently mocked.

## Verification and evidence

Worker reports focused vet/tests and actual responses if authorized to launch. Sisyphus independently runs the actual process and verifies status/body/header and input boundaries after receiving the result, not merely a not-throw test. Save scrubbed test/response evidence to captures/huy/step-1/api.md, producing source SHA after authorized integration, and update own lessons/runtime README. Exact response bytes remain unit-tested; terminal capture is actual operation proof.

## Learning gate

Owner changes APP_ENV at startup and observes its header, requests success/invalid/failure paths, and explains process hostname versus container hostname versus pod name. Explain why /health proves the process can respond but does not promise every external exporter/backend is healthy.

## Completion and blockers

Build/test/runtime/learning gates pass before Stage 03. If port 8080 is occupied, identify its owner; do not kill an unrelated process. Stop your prior proof instance or obtain an authorized free port adjustment for the exercise, then retain deployment contract 8080. No alternate endpoint aliases, fake ID or more product features.

## Sisyphus handoff prompt

```text
Coordinate Stage 02 from .omo/plans/devops-personal/02-runtime-api.md and its
personal ticket. Confirm Stage 01, then prepare stage/huy/02-runtime-api from the
accepted source SHA in its isolated worktree. Preliminary routine-verified.
Do not implement. Compose a worker prompt from the current Worker role and the
bounded 02A/02B scope, preserving exact bodies, identity/configuration and input
boundaries. No Docker/CI/telemetry or extra API features. Receive the handoff,
independently exercise the actual server plus behavior checks, and require the
owner's learning gate. Commit/integrate only with authorization; record runtime
evidence tied to the source SHA. Do not close on compilation/tests alone.
```
