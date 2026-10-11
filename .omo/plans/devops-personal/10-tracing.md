# Stage 10: Real backend traces, context-correlated logs and timing

## Governing artifacts

[Spec](../../../.scratch/assignment-3-personal/spec.md), [ticket 10](../../../.scratch/assignment-3-personal/issues/10-tracing.md), bonuses 3/5 and research e,f, overview Approach 10, [execution contract](00-execution-contract.md).

## Entry gate and routing

Stage 09 accepted actual dashboard and GitOps path. Preliminary complex-critical: SDK/Collector/viewer contract and environment rollout. Read current Go OTel trace/OTLPHTTP/otelhttp APIs via official docs/Context7, Collector component schema/changelog and [Jaeger getting started](https://www.jaegertracing.io/docs/2.22/getting-started/). Resolve compatible stable module/image versions; planning observed Jaeger2.22.0, not a perpetual latest claim.

## Allowed scope and branches

Source stage/huy/10-tracing from accepted origin/personal/huy, matching worktree. Files app/telemetry.go, app/main.go/main_test.go, go.mod/go.sum, k8s/tracing/{collector,jaeger}.yaml, .github/workflows/ci.yml only runtime-smoke environment additions, own runbook/versions/lessons and bonus3/5 evidence. Separate stage-env/huy/10-tracing for actual ConfigMap/env annotation update; image through normal digest PR.

Operations monitoring namespace tracing resources and Argo-managed app only. No central log backend, OTLP metric migration, legacy jaeger exporter, public trace/admin endpoint, credentials in ConfigMap, or production durable-storage claim. Commit scope SDK/config/JSON logging/timing/declarations/docs only.

## Bounded substeps

### 10A: Validate Collector and actual viewer

Resolve/pin otel/opentelemetry-collector-contrib image digest and verify its validate --config command plus modern otlp_http component before installation. Configuration uses OTLP HTTP receiver4318, memory_limiter check1s/limit192Mi/spike32Mi, batch and exporter otlp_http/jaeger to http://jaeger.monitoring.svc.cluster.local:4318. traces pipeline receivers[otlp], processors[memory_limiter,batch], exporters[otlp_http/jaeger]. No obsolete jaeger exporter/14250 or unverified deprecated alias.

Jaeger single-process current stable 2.22.0 image cr.jaegertracing.io/jaegertracing/jaeger pinned by actual digest; configure its documented transient memory mode. ClusterIP Service jaeger4318/16686 and Service otel-collector4318. Collector resources requests100m/128Mi limits500m/256Mi; Jaeger100m/128Mi limits1CPU/512Mi. Use actual image/schema to validate defaults, not a fake config accepted only by a different version. Deployment/Service/ConfigMap manifests grouped per component; no secrets.

Apply namespace monitoring if not already present, validate configuration using exact image locally or owned cluster validation job, then apply only these objects and wait Ready. Capture real validation/ready logs. Exit gate: actual compatible receiver/exporter path exists. No UI public tunnel.

### 10B: SDK lifecycle and configuration-before-image rollout

Define initTelemetry(ctx context.Context, id string) (func(context.Context) error, error). main passes resolved runtime ID; configure resource service.name devops-app-huy, service.instance.id id; OTLPHTTP exporter, batch processor and always-on sampling for this small learning workload. Set global tracer provider and TraceContext+Baggage propagation once. SDK init failure is startup failure; actual asynchronous export errors are surfaced, not suppressed. Flush/shutdown with bounded context on normal termination, integrated with actual HTTP shutdown, not a dead deferred shutdown behind os.Exit.

Set OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector.monitoring.svc.cluster.local:4318 through devops-app-config envFrom. Before merging the new instrumented image's env PR, merge reviewed ConfigMap addition and update Deployment pod-template annotation config-revision to the configuration change's UTC timestamp. Argo restarts pods so environment changes reach process. ConfigMap edit alone does not refresh env. Existing image ignores extra endpoint safely; this makes new-image startup deterministic.

Update local/CI smoke environment to pass a syntactically valid endpoint. Local CI with no reachable Collector proves only HTTP artifact startup; exporter errors are honest and not counted as trace proof. Cluster request/viewer evidence is mandatory. No fake exporter/no-op telemetry fallback.

Exit gate: correct SDK/config lifecycle and actual deployed app continues existing contracts through signed CI/desired-state rollout.

### 10C: Handler context, structured log and Server-Timing

Wrap GET /test using otelhttp.NewHandler(inner, "GET /test"). Inner workload/logging/metrics see the wrapper's request context. Create child test-work span around real selected delay/work, end it on success/error and record intentional error/status accurately. Propagation accepts real upstream traceparent for Stage 11. Do not trace /health,/metrics,/otel/v1/traces or SDK export traffic as workload.

Use slog JSON stdout; obtain valid trace/span context from r.Context, log request completed with trace_id/span_id, bounded route, actual status, measured duration and runtime identity. Never generate a fixed ID or use unrelated startup context. This is OTel-context log correlation, not deployment of OTel Logs SDK/Loki. Existing metric labels remain bounded; trace IDs are log fields, not metric labels.

Set Server-Timing app;dur=measured-ms on /test BEFORE first response WriteHeader/Write. A deferred header change after writing response does not work. Preserve exact response bodies/config headers and existing metric behavior. Duration reflects real handling; invalid requests must not fake selected delay. Focused tests cover meaningful header/error/context boundaries if needed, not log wording or forwarded mocks. Verify actual browser Network timing as runtime proof.

## Verification and evidence

```sh
(cd app && go vet ./... && go test ./...)
curl -i 'http://127.0.0.1:8080/test?delay_ms=200'
curl -i 'http://127.0.0.1:8080/test?delay_ms=200&fail=1'
kubectl -n devops-app logs -l app=devops-app --prefix --since=5m
kubectl -n monitoring port-forward svc/jaeger 16686:16686
```

Use live deployed app forward from Stage 07, then private Jaeger UI. Find ACTUAL delayed request under devops-app-huy service, inspect server/child spans and timing, match trace_id and logged server span_id to the same request in viewer. Capture successful and controlled-error status without confusing separate requests. Server-Timing measured value is visible in Chrome Network and readable via performance resource entry serverTiming for same-origin fetch later.

Sisyphus independently verifies Collector/viewer config/image, full app export, matching log IDs, timing/header placement and existing metrics. Save bonus3 viewer/config/request capture and bonus5 log/viewer match, each with digest/pod/time. Avoid whole Secret/config dumps. Jaeger restart loses history; regenerate a real trace, retain committed captures and state this limitation.

## Learning gate

Owner sends delayed/error requests and locates their actual spans/logs, explains parent/child/W3C context, sampling/batching, Collector routing versus raw exporter, Semantic Conventions and per-language signal maturity. Explain stdout correlation versus central log ingestion, header timing versus SDK traces, and transient storage.

## Completion and blockers

Viewer readiness or SDK init is not exported-trace proof. Diagnose endpoint/env rollout, protocol/port, SDK errors, Collector logs and viewer query/time range. No obsolete exporter or swallowed error shortcut. Missing current component blocks until a compatible supported pinned image validates. Backend trace/log/timing gates pass before browser stage.

## Sisyphus handoff prompt

```text
Coordinate Stage 10/ticket 10 after Stage 09. Prepare isolated source/env branches,
preliminary complex-critical. Do not implement. Dispatch exact-image Collector/
Jaeger validation, SDK/config rollout and inner-context logging/timing subtasks
with verification between. Preserve app/metric contracts, private UIs and same
project token policy. ConfigMap/env rollout precedes new image; no fake exporter.
Independently send real traffic and match viewer/log IDs plus measured timing,
require owner learning, and commit authorized SDK/manifests/scrubbed proof only.
Do not call readiness or a local unreachable-exporter smoke successful tracing.
```
