# Stage 11: Real browser page, joined trace and Web Vital

## Governing artifacts

[Spec](../../../.scratch/assignment-3-personal/spec.md), [ticket 11](../../../.scratch/assignment-3-personal/issues/11-browser.md), bonuses 6/8 and research f, overview Approach 11, [execution contract](00-execution-contract.md).

## Entry gate and routing

Stage 10 actual backend trace/log and correct endpoint accepted. Preliminary complex-critical: public ingestion proxy and cross-browser/backend propagation. Read [current OTel browser guide](https://opentelemetry.io/docs/languages/js/getting-started/browser/), SDK/instrumentation package APIs/peer dependencies via Context7, [web-vitals](https://github.com/GoogleChrome/web-vitals) and [W3C Trace Context](https://www.w3.org/TR/trace-context/). Browser instrumentation is experimental; do not call it uniformly stable.

## Allowed scope and branch

stage/huy/11-browser from accepted origin/personal/huy, matching worktree. Files web/package.json/package-lock.json/src/main.js, app/static/index.html, generated app/static/app.js, app/main.go/main_test.go and focused proxy.go if needed, app/Dockerfile, root .dockerignore, .github/workflows/ci.yml for browser build prerequisites, own docs/versions/bonus6/8 captures. Env endpoint is already configured; any adjustment uses separate reviewed stage-env/huy/11-browser.

No React/Vite product, TypeScript migration, auth redesign, arbitrary forward proxy, separate public Collector, production RUM platform, paid service or token in browser bundle. Commit scope: minimal real frontend/proxy/build and evidence. Whether to commit generated app.js is decided once: ignore it as build output and generate before every embed build; commit source and lockfile. Docker context must explicitly allow the generated artifact after npm build.

## Bounded substeps

### 11A: Functional real page and deterministic embed build

Plain HTML heading, Run test button, response/trace/measurement areas. Go serves exact root using current ServeMux exact-root pattern and only intended bundle path via go:embed. Unknown routes still 404. No mocked response. Bundle main.js with esbuild browser/IIFE/ES2015 into app/static/app.js. npm ci + npm run build is reproducible from pinned compatible packages/lockfile.

Add native BUILDPLATFORM Node26 build stage pinned by actual multi-platform index digest, npm ci/build from root web source, copy generated app/static output into Go stage BEFORE compilation. Native node bundle is architecture-independent; Go still compiles for TARGETOS/TARGETARCH. Update local/CI check and PR native build ordering to produce assets before Go embed. setup-node uses an explicit anonymous token setting in unprivileged checks; no implicit GitHub identity. No placeholder bundle/HTML to make embed compile.

```sh
(cd web && npm ci && npm run build && npm ls)
(cd app && go vet ./... && go test ./...)
```

Exit gate: actual page/button displays live /test response, generated bundle loads, existing API and unknown-path behavior remain correct, native/multi-platform image builds with real assets.

### 11B: Bounded same-origin OTLP route

Expose only POST /otel/v1/traces to configured private Collector /v1/traces. Validate OTEL_EXPORTER_OTLP_ENDPOINT syntax/allowed http(s) target at startup; absent/invalid target is an explicit error once proxy is enabled. Request input cannot choose upstream. Max body1MiB, request context/transport deadline5s, method405, over-limit413, unavailable upstream502. Remove incoming credential cookies/Authorization from this telemetry proxy, preserve valid OTLP content type and report real upstream result. No workload tracing/metrics for export route.

Enforce body limit BEFORE upstream dispatch, including chunked/unknown-length requests: bounded read with MaxBytesReader and explicit size error handling, then reuse validated bytes for proxy. A streaming proxy that turns over-limit errors into generic502 fails acceptance. Do not add arbitrary retry/auth/rate-limit subsystem or claim authenticated production ingestion; public demo ingestion has that explicit limitation.

Focused checks exercise actual wrong method, oversize/chunked body, unavailable upstream and valid transport against a local receiver if needed, verifying consumer status and body/target restrictions. No source-text/forwarding/mock-echo tests. Sisyphus's actual cluster browser export is the positive runtime proof.

Exit gate: bounded route behaves accurately; Collector/admin port remains private.

### 11C: Web SDK and one joined trace

Pin compatible current packages: @opentelemetry/api, core, resources, sdk-trace-base, sdk-trace-web, context-zone, instrumentation, instrumentation-fetch, instrumentation-xml-http-request, instrumentation-document-load, exporter-trace-otlp-http and web-vitals, plus esbuild. Use actual peer constraints, not an old version list or nonexistent propagator package.

resourceFromAttributes service.name devops-web-huy. WebTracerProvider constructor spanProcessors [BatchSpanProcessor(exporter)]; provider.register ZoneContextManager; register document-load/fetch/XHR instrumentation. Same-origin exporter URL window.location.origin + /otel/v1/traces. Ignore exact export URL in fetch/XHR to prevent recursive export spans. Verify current propagation API and actual traceparent injection; no hardcoded ID.

Button starts active test-click span and invokes fetch /test?delay_ms=100. Use Promise chains so browser context is preserved; do not rely on native async/await Zone preservation. Display actual response and trace ID; record real HTTP/exception status; end span on both success/error and flush provider after the click span ends. Do not show success before response or leave failed span open.

Exit gate: Chrome network request has traceparent; ONE viewer trace includes browser click/fetch and backend GET /test/test-work with correct trace ID and parent-child relationships, also matching backend JSON log. Independent browser/server trace screenshots are insufficient.

### 11D: Real Web Vital and browser performance explanation

Register onLCP/onINP/onCLS with reportAllChanges true. Show/log real name/value/measurement ID and correct unit, plus page URL, browser/context and capture time. Use Chrome, load actual tunneled own page, click and hide tab to finalize metrics. At least one actual supported callback measurement is required; zero CLS can be real, not a placeholder. No invented number or copied public-site result in place of own page proof.

Inspect same-origin /test Server-Timing in Network and PerformanceResourceTiming.serverTiming. Explain Navigation/Resource Timing and PerformanceObserver, lab versus real-user sessions and CrUX field aggregation. Controlled local session is not a production RUM cohort or CrUX population statistic.

## Verification and evidence

Sisyphus opens actual deployed page in Chrome/browser tooling, clicks, observes live result/console/network, captures screenshot plus traceparent/viewer/log linkage and measurement context, then closes verification tab. Owner supplies phone current URL if tunnel changed. Confirm bundle/SDK export absence of recursion and API/security error boundaries; no UI-only mock proof.

Save bonus6 joined-trace/network/log and bonus8 actual-vital captures. Record module/npm/image/action pins, source/digest/time, experimental SDK limitation and transient viewer storage. Native smoke and full signed CI/desired-state rollout must still pass after embed/proxy change.

## Learning gate

Owner triggers and locates joined trace, identifies traceparent/parentage and explains async context, same-origin proxy versus exposing Collector, Performance APIs/Server-Timing, LCP/INP/CLS units and field/lab distinctions. Owner reproduces one real measurement; worker output is not learning certification.

## Completion and blockers

No verification on fake frontend or static ID. Diagnose Console, request headers, exported batches, Go extraction, Collector and viewer range; flush only after spans end. If callback doesn't fire, exercise appropriate Chrome lifecycle and record support, not invent data. Wrong peer/API package requires checking current docs, not compatibility shim proliferation.

## Sisyphus handoff prompt

```text
Coordinate Stage 11/ticket 11 after backend trace/log proof. Prepare isolated stage
worktree, preliminary complex-critical. Do not implement. Dispatch real page/build,
bounded proxy, Web SDK propagation, then Web Vital subtasks with checks between.
Preserve exact API/metrics, token-free browser/check jobs and private Collector.
Independently exercise actual deployed Chrome UI, proxy boundaries and ONE joined
trace with network/header/log linkage, plus real measurement/context. Require
owner's propagation/performance learning gate. Commit authorized source/lockfile/
build changes and scrubbed evidence; no fake bundle/value or production RUM claim.
```
