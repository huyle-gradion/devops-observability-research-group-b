# Stage 09: Real workload metrics and PVC-backed dashboard

## Governing artifacts

[Spec](../../../.scratch/assignment-3-personal/spec.md), [ticket 09](../../../.scratch/assignment-3-personal/issues/09-metrics.md), Step 4/research d,e, overview Approach 9, [execution contract](00-execution-contract.md).

## Entry gate and routing

Stage 08 accepted automatic deploy loop. Preliminary complex-critical: instrumentation and cross-system scrape/storage contract. Read current prometheus/client_golang, Prometheus Operator ServiceMonitor and pinned kube-prometheus-stack values/schema before editing. Current dependencies/docs must be verified through official sources/Context7 at execution, not copied from an old dashboard blog.

## Allowed scope and branches

Source stage/huy/09-metrics from accepted origin/personal/huy, matching worktree. Files app/main.go, new app/metrics.go only if cohesion warrants it, app/main_test.go, go.mod/go.sum, k8s/monitoring/{values,servicemonitor}.yaml, dashboards/huy-workload.json, docs/runbook/versions/lessons and Step 4 captures. Deployment digest updates via reviewed CI PR to env/huy only. If resource/config edits are necessary, separate stage-env/huy/09-metrics PR; no resurrected source k8s/app authority.

Authorized operations only Helm huy-monitoring in monitoring, namespace bootstrap and own persisted workload pods/PVCs. No Alertmanager/paging, central logs, trace SDK yet, other cluster or deleting PVCs/data. Commit scope: app instrumentation, monitoring declarations/dashboard/docs.

## Bounded substeps

### 09A: Accurate application instrumentation

Use prometheus/client_golang with handler-owned custom Registry, avoiding package-global registration collisions in tests. CounterVec http_requests_total{method,route,status}; HistogramVec http_request_duration_seconds{method,route}, seconds buckets .005,.01,.025,.05,.1,.25,.5,1,2.5. Register once per serving handler, not each request. /metrics via promhttp.HandlerFor that registry.

Record actual GET /test responses, including 200/400/500, and actual elapsed seconds. Exclude /health, /metrics, root/static and future export proxy; exclude non-GET requests rather than disguising HEAD/POST as GET. Labels only bounded GET, /test and actual bounded status; no query strings, pod/trace IDs, raw URLs. Preserve exact API bodies/config behavior.

Focused checks: existing behavior tests plus one meaningful workload-metric boundary test if needed: actual success/error/invalid traffic changes correct series while health/scrapes do not inflate workload. No tests of wrapper wiring or text layout. Run actual process, send workload and read /metrics. Integrate through green signed CI/digest PR/auto-sync.

Exit gate: actual emitted metrics match observed traffic/status; no scrape-driven feedback/cardinality explosion.

### 09B: Monitoring namespace, chart values and discovery

Pin compatible kube-prometheus-stack chart version and component images; record helm show values source and actual image digests. Install release huy-monitoring in monitoring. Explicit local learning budgets: Prometheus requests200m/512Mi limits1CPU/1Gi, retention24h, local-path RWO2Gi PVC; Grafana requests100m/128Mi limits500m/256Mi, local-path RWO1Gi persistence. Check scheduling/usage and adjust through authorized evidence-based configuration, not pretend these are minimum guarantees.

Keep kubelet/cAdvisor, node exporter and kube-state-metrics. Disable Alertmanager and unavailable k3s scheduler/controller-manager/etcd/kube-proxy scrape jobs, documenting actual disabled chart fields. Do not hide app failures by disabling its target. Set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues false; explicit serviceMonitorSelector.matchLabels.release huy-monitoring and serviceMonitorNamespaceSelector.matchLabels.kubernetes.io/metadata.name monitoring.

ServiceMonitor lives in monitoring, release label huy-monitoring, selects app: devops-app Service in namespaceSelector.matchNames devops-app; endpoint named Service port http, path /metrics, interval15s. Install CRDs/chart before applying ServiceMonitor. values.yaml is Helm input, NOT a Kubernetes object: never kubectl apply the entire directory containing it.

```sh
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm show values prometheus-community/kube-prometheus-stack --version "$MONITORING_CHART_VERSION"
helm upgrade --install huy-monitoring prometheus-community/kube-prometheus-stack \
  --version "$MONITORING_CHART_VERSION" --namespace monitoring --create-namespace \
  -f k8s/monitoring/values.yaml --wait --timeout 10m
kubectl apply -f k8s/monitoring/servicemonitor.yaml
kubectl -n monitoring get pods,svc,pvc
```

Exit gate: workloads Ready, PVCs Bound, own app scrape target Up with observed labels and emitted series. Discover actual generated Service names for private UI forwards; don't guess another chart release's names.

### 09C: Custom dashboard and persistence

Create/import/save a workload dashboard with 5s refresh, Prometheus datasource, actual request/error/latency and per-pod CPU/memory panels. Export working JSON to dashboards/huy-workload.json; avoid hardcoding another instance's datasource UID when importing. Error ratio panel unit is percentunit; latency seconds, CPU cores, memory bytes. Include time range/namespace/source labels.

Use these baseline queries, adjusting only selectors to VERIFIED actual target labels:

```promql
sum(rate(http_requests_total{namespace="devops-app",route="/test"}[1m]))

sum(rate(http_requests_total{namespace="devops-app",route="/test",status=~"4..|5.."}[1m]))
/ clamp_min(sum(rate(http_requests_total{namespace="devops-app",route="/test"}[1m])), 0.001)

histogram_quantile(0.95, sum by (le) (
  rate(http_request_duration_seconds_bucket{namespace="devops-app",route="/test"}[1m])
))

sum by (pod) (rate(container_cpu_usage_seconds_total{
  namespace="devops-app",container="app",image!=""}[1m]))

sum by (pod) (container_memory_working_set_bytes{
  namespace="devops-app",container="app",image!=""})
```

Send actual success/error/delay traffic; wait for multiple 15s scrapes and 1m rate window. Idle/no samples is not fake zero or a chart bug to suppress. Resource panels must identify each app pod even if Service port-forward loads only one. Save dashboard to Grafana before testing persistence.

Delete only one observed own Grafana pod and wait replacement, verify saved dashboard remains. Repeat for own Prometheus pod and verify earlier time-series history after restart. Never delete PVC or StatefulSet to manufacture success. local-path RWO is bound to its node, not replicated/HA.

## Verification and evidence

Sisyphus independently observes raw emitted metrics, target discovery/Up, chart config/PVCs, actual workload panels, exported dashboard JSON and pod-replacement persistence. Run go vet/test plus actual deployed traffic. Save Step 4 captures with time range/namespace/source digest and real data. Stage 12 proves sustained load; this stage is not complete on Helm success or stock dashboard screenshots alone.

## Learning gate

Owner changes traffic and explains rate versus counter, histogram p95, actual error ratio, bounded label cardinality, pull scraping and ServiceMonitor discovery. Identify app-emitted metrics versus kubelet/cAdvisor, per-pod resource differences, emptyDir versus PVC and why namespace/RWO don't guarantee isolation/HA.

## Completion and blockers

No invented panels/data or broadening selectors to unrelated workloads. Pending PVC/targets/permissions require describe/events and actual label/network diagnostics. Helm values are applied only through Helm. Persist observed fixes and leave dashboard/stack running for traces/load.

## Sisyphus handoff prompt

```text
Coordinate Stage 09/ticket 09 after accepted GitOps. Prepare source stage worktree
and separate env PR if needed. Preliminary complex-critical. Do not implement.
Dispatch accurate app metrics, pinned monitoring/discovery, then custom dashboard
and persistence as bounded tasks with checks between. No app apply outside Argo,
no kubectl apply of Helm values.yaml, paging or fabricated metrics. Independently
verify emitted status/time series, actual target Up, own panels/per-pod resources
and PVC survival. Require owner's PromQL/discovery/storage learning gate. Commit
only authorized instrumentation/declarations/exported JSON/scrubbed proof.
```
