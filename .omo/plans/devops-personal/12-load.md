# Stage 12: Controlled k6 load and live observability proof

## Governing artifacts

[Spec](../../../.scratch/assignment-3-personal/spec.md), [ticket 12](../../../.scratch/assignment-3-personal/issues/12-load.md), Steps 4/5/research e, overview Approach 12, [execution contract](00-execution-contract.md).

## Entry gate and routing

Stages 09–11 accepted with actual target/dashboard, trace/log and browser proofs. Preliminary routine-verified: bounded test client against own private service. Escalate if target/permission is ambiguous, resource loss or secret handling arises. Read current [k6 scenarios](https://grafana.com/docs/k6/latest/using-k6/scenarios/), response callback, checks and tagged thresholds documentation before implementing.

## Allowed scope and branch

stage/huy/12-load from accepted origin/personal/huy, matching worktree. Files load-tests/test.js, docs/huy/load-runbook.md/lessons.md/requirements.md and Step4/5/bonus evidence captures. No app/collector/dashboard metric fabrication, token use in request headers, paid load service, unbounded stress or unrelated/public targets. Commit scope: controlled client/runbook and actual proof. No new permanent suite to test k6 wiring.

## Bounded substeps

### 12A: Deterministic traffic and truthful thresholds

Read BASE_URL from environment, require explicit own authorized target. One ramping-VU workload:30s ramp to10 VUs,60s hold10,30s ramp down; bounded graceful completion. Each iteration sends nine /test?delay_ms=50 expected200 then one /test?delay_ms=200&fail=1 expected500. Success checks body prefix Hello world from container- and expected status; controlled-error checks exact intentional test failure newline and500. No random errors or invented backend timing.

Use custom request tag case:success versus case:controlled-error, not scenario:success. scenario is k6's system scenario tag; the overview's conceptual distinction must not collide with it. The actual workload remains ONE scenario so nine-to-one traffic is preserved. Set http.setResponseCallback(http.expectedStatuses(200,500)); thresholds checks rate==1 and http_req_failed{case:success} rate<0.01. Wrong success500 still fails strict status/body checks even when callback classifies500 expected. Application counters STILL classify500 as app errors. No misleading global <5% threshold on deliberately ~10% app errors.

Exit gate: script uses bounded own target and correct case/status/body semantics. Focused worker check is actual run against deployed app, not source-text assertions.

### 12B: Real private load and synchronized dashboard

Install verified stable official k6 binary/package when authorized, record version. Keep app service port-forward localhost8080. Browser/tunnel proof is separate; deterministic load goes to private forward rather than depending on ephemeral tunnel/provider limits.

```sh
k6 version
k6 run --env BASE_URL=http://127.0.0.1:8080 load-tests/test.js
kubectl -n devops-app get pods -l app=devops-app -o wide
kubectl -n monitoring get pods
```

Open private Grafana custom dashboard with5s refresh, view current time range and all request/error/latency/resource panels. Warm up until multiple15s scrape intervals/1m rate window exist; record baseline, ramp/steady and cooldown. Capture k6 stdout/exit status, actual traffic count/tags/checks and same-time dashboard. Expected error ratio near10% over steady interval, not necessarily exact during ramp/window transitions. CPU/memory display real per-pod samples; don't promise a large spike from a tiny sleeping handler or fabricate one.

Service port-forward loads one selected pod connection, so one pod may be busy while the other is quieter. Pod resource panels should honestly show both. If forward's pod rolls during experiment, recreate own forward and rerun; do not count transport errors from abandoned forwarding as a successful demonstration. No source/env deployments during the accepted timed run.

Exit gate: actual k6 checks and success unexpected-failure threshold pass, dashboard workload series reacts during/cooldown, measured resources visible, tracing/logging remain operational. If thresholds fail, diagnose real status/body/forward/export/resource issue before acceptance.

### 12C: CP3 packet and repeatable live script

Retain earlier backend/browser/measurement captures; generate fresh traces after any Jaeger restart. Assemble CP3 before Tue14:00, present at15:00. Owner runs k6 while explaining live panels and checks. Save the active chart time range, source/index digest, replicas/node placement, k6 version, invocation/time and exit status. Do not label a recording as required Friday live demo.

Document exact start/access/run/observe/stop sequence and troubleshooting: service forward lifecycle, private Grafana service discovery, target Up, CPU/memory versus request metrics, expected500 interpretation and sample-window lag. Stop only own load/client/forward when needed; retain cluster/services/data. Repeat live Wed review and Fri Oct16 in Stage14.

## Verification and evidence

Sisyphus independently runs actual script, observes dashboard while traffic occurs, checks histogram/error counts against requests and reviews target authorization/no secret. Save captures/huy/step-5/k6-output.log and synchronized Step4/5 screenshots/recording with actual run metadata. No theoretical/chart-only evidence, bare nonempty tests, or CLI output copied from docs.

## Learning gate

Owner launches test and explains VUs versus request rate, deliberately expected500 versus operational failure, why checks and http_req_failed differ, p95 estimation/time windows, one-pod forwarding versus Service balancing, and app versus resource metrics. Owner diagnoses one actual failure if encountered; don't manufacture a production incident.

## Completion and blockers

CP3 technical scope is Step4/5 rehearsal and bonuses3/5/6/8 demonstrated, but Step5's Oct16 live occurrence remains future/not-run until delivered. No missing panel/bonus downgraded. If OOM occurs, inspect actual limits/usage then get authorized adjustment and repeat planned10-VU profile; never drop required components for green screenshots.

## Sisyphus handoff prompt

```text
Coordinate Stage12/ticket12 after metric/trace/browser gates. Prepare isolated
stage/huy/12-load worktree; preliminary routine-verified unless new risk. Do not
implement. Dispatch bounded k6/client/runbook work with custom case tags, nine-to-
one expected responses and truthful thresholds. Target only own private forwarded
app. Independently run real load while observing own Grafana panels and resource
samples, collect synchronized proof and require owner interpretation. Commit
only authorized client/docs/scrubbed captures. CP3 is Tue15:00; Friday live demo
must remain a separate future gate, not pre-certified from this recording.
```
