# Manual deploy runbook — devops-app on k3s (huy)

Step 2 of the assignment: apply the app manifests by hand with `kubectl`, watch the
rollout, and reach the pods. Cluster build and kubeconfig are in
[`cluster-runbook.md`](cluster-runbook.md). Manifests live in `k8s/app/`.

The image is pinned by digest:
`ghcr.io/huyle-gradion/devops-app-huy@sha256:99ab49d854c336b0379e38d37b0f2fb35e6df27caca2a15239ce886dade018fa`.

## 0. Point kubectl at the cluster

```sh
export KUBECONFIG="$HOME/.kube/huy-k3s.yaml"
kubectl config current-context   # expect huy-k3s
```

The kubeconfig stays on the workstation, mode 0600. Never commit it.

## 1. Namespace first

`kubectl apply` of a namespaced object fails if its namespace does not exist, and
`--dry-run=server` needs a live namespace too. Apply the namespace on its own first:

```sh
kubectl apply -f k8s/app/namespace.yaml
```

## 2. Server-side dry-run

Ask the API server to validate the objects without persisting them:

```sh
kubectl apply --dry-run=server -f k8s/app/
```

Expect `configured`/`created (server dry run)` lines and no errors. A validation error
here means the manifest is wrong; fix the file before continuing.

## 3. Apply

```sh
kubectl apply -f k8s/app/
```

## 4. Watch the rollout

```sh
kubectl -n devops-app rollout status deployment/devops-app --timeout=180s
```

Expect `deployment "devops-app" successfully rolled out`. Some pods start out
`Terminating` or `Init` briefly — wait for the rollout rather than reading the first
`get pods`.

## 5. Inspect reality

```sh
kubectl get nodes -o wide
kubectl -n devops-app get all
kubectl -n devops-app get pods -l app=devops-app -o wide
```

`get nodes -o wide` shows the actual physical hosts and their addresses; the app pod
column on the last command shows which node each replica landed on.

## 6. Reach the app

Service forward (good for a single request; it picks one backing pod, so repeated
responses do not prove load balancing):

```sh
kubectl -n devops-app port-forward svc/devops-app 8080:8080
```

Then, in another terminal:

```sh
curl -i http://127.0.0.1:8080/health                 # 200 ok
curl -i http://127.0.0.1:8080/test                   # X-App-Env: cluster-huy, pod name in the body
```

Per-pod forward, to address each replica directly. Read the pod names from step 5 and
use a free local port per pod:

```sh
kubectl -n devops-app port-forward pod/<first-pod-name>  8081:8080
kubectl -n devops-app port-forward pod/<second-pod-name> 8082:8080
```

The `/test` body is `Hello world from container-<POD_NAME>`, so the two forwards should
report the two different pod names.

## 7. Expose it off the machine (quick tunnel)

Keep a forward alive, then start a Cloudflare quick tunnel against it:

```sh
cloudflared tunnel --url http://127.0.0.1:8080
```

Open the printed `https://<random>.trycloudflare.com/test` from a device on mobile data.
Only the app is exposed; the quick URL is ephemeral, so recreating the tunnel prints a
new URL. If a local `cloudflared` config interferes, point the command at an isolated
config path instead of deleting unrelated config.

## Single-node note (topology spread)

The stage plan specified a hard spread — `whenUnsatisfiable: DoNotSchedule` with
`minDomains: 2` — which assumes a two-host cluster. This cluster is currently
single-node (the second physical host is not attached yet), so this manifest uses
`whenUnsatisfiable: ScheduleAnyway` with `maxSkew: 1` over `kubernetes.io/hostname`.
That keeps both replicas schedulable today and still prefers spreading them once a
second host joins — but it is best-effort, not a guarantee. Do not describe this as
two-node placement until a second host exists and the pods actually land on both.

`minDomains` is omitted, not forgotten. The Kubernetes API only accepts `minDomains`
together with `whenUnsatisfiable: DoNotSchedule`; with `ScheduleAnyway` the API server
rejects it, and with `DoNotSchedule` it would force the second replica to stay
`Pending` on a single node. When a second physical host is attached and a hard spread
is wanted, switch this constraint back to:

```yaml
whenUnsatisfiable: DoNotSchedule
minDomains: 2
```

Re-run `kubectl apply --dry-run=server -f k8s/app/` before applying that change.

## 8. Monitoring stack (Stage 09B)

Prometheus + Grafana run as `kube-prometheus-stack` release **`huy-monitoring`** in
the **`monitoring`** namespace, scraping the app through a ServiceMonitor. The pinned
chart version, the `helm show values` source and the component image digests are
recorded in [`versions.md`](versions.md).

`k8s/monitoring/values.yaml` is **Helm input only**. It is not a Kubernetes object —
never `kubectl apply` the directory that contains it (only `k8s/monitoring/servicemonitor.yaml`
is a real object). Apply it exclusively through `helm upgrade --install -f`.

CRDs and the chart must exist **before** the ServiceMonitor is applied, because the
`ServiceMonitor` kind is defined by the chart's Prometheus Operator CRDs.

```sh
export KUBECONFIG="$HOME/.kube/huy-k3s.yaml"

# 1. Chart repo + pin (read-only)
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm show values prometheus-community/kube-prometheus-stack --version 92.3.0   # reference only

# 2. Install/upgrade the stack (namespace created on first run). Values are Helm
#    input; -f points at the file, the directory is never kubectl-applied.
helm upgrade --install huy-monitoring prometheus-community/kube-prometheus-stack \
  --version 92.3.0 \
  --namespace monitoring --create-namespace \
  -f k8s/monitoring/values.yaml \
  --wait --timeout 10m

# 3. Only now, with the CRDs present, apply the scrape target
kubectl apply -f k8s/monitoring/servicemonitor.yaml

# 4. Inspect reality
kubectl -n monitoring get pods,svc,pvc
```

The `--wait` above waits for the release's workloads to become Ready; PVCs should be
`Bound` (k3s local-path provisioner). After a few 15s scrapes, the target should be
`Up`; confirm in the Prometheus UI (`Targets`) or via
`kubectl -n monitoring get servicemonitor` then the Prometheus `targets` page.

### Discoverable Service names come from this release

The Services and PVCs that Prometheus/Grafana expose for UI forwarding and health are
named after **this** release (`huy-monitoring-…`) and the Operator-created
`prometheus-operated` Service — they are not guessed from, or shared with, any other
chart release. Read the actual names from `kubectl -n monitoring get svc` before
port-forwarding, rather than assuming another instance's names:

```sh
kubectl -n monitoring get svc
kubectl -n monitoring port-forward svc/huy-monitoring-grafana 3000:80
```

The `kubernetes.io/metadata.name: monitoring` label the namespace carries (added
automatically by Kubernetes) is what the release's Prometheus uses in
`serviceMonitorNamespaceSelector`; the ServiceMonitor carries `release: huy-monitoring`
to match `serviceMonitorSelector`. A wrong label on either side yields no targets
silently, so verify the target is actually `Up` rather than assuming discovery worked.

## 9. Workload dashboard and persistence proof (Stage 09C)

The dashboard source of truth is committed at [`dashboards/huy-workload.json`](../../dashboards/huy-workload.json).
It has five panels — request rate, error ratio (`percentunit`), p95 latency (`seconds`),
per-pod CPU (`cores`) and per-pod memory (`bytes`) — a **5s refresh**, and it reads the
Prometheus datasource through a `${datasource}` template variable, so **no other
instance's datasource UID is baked in**. The two resource panels come from
kubelet/cAdvisor and are keyed by pod; the three request panels come from the app's own
`/metrics` registry (`http_requests_total`, `http_request_duration_seconds`) for
`GET /test` in namespace `devops-app`.

### 9.1 Import and save it in Grafana

Forward the Grafana UI (the Service name comes from this release — read it first, do not
assume) and import the committed JSON:

```sh
export KUBECONFIG="$HOME/.kube/huy-k3s.yaml"
kubectl -n monitoring get svc                       # read the actual Grafana Service name
kubectl -n monitoring port-forward svc/huy-monitoring-grafana 3000:80
```

In the browser at <http://127.0.0.1:3000>:

1. **Dashboards → New → Import**.
2. Upload `dashboards/huy-workload.json` (or paste its contents).
3. When prompted, pick the **Prometheus datasource of the `huy-monitoring` release**
   for the `datasource` variable. That selection is what the panels use — the JSON
   itself carries no foreign UID, so importing on your instance binds to *your*
   Prometheus, not someone else's.
4. **Save**, and put it in a folder you own (for example a `huy` folder).

Saving here is deliberate: it is the object the persistence proof below must survive.

### 9.2 Send real traffic first, then wait

Panels will be empty or flat until the app actually serves workload. Drive real
requests through the app (from step 6's forward, or per-pod forwards), then be patient:

```sh
# success + invalid + intentional failure + added latency
for i in $(seq 1 30); do curl -s -o /dev/null http://127.0.0.1:8080/test; done
for i in $(seq 1 10); do curl -s -o /dev/null "http://127.0.0.1:8080/test?fail=1"; done
for i in $(seq 1 5);  do curl -s -o /dev/null "http://127.0.0.1:8080/test?delay_ms=1001"; done
for i in $(seq 1 10); do curl -s -o /dev/null "http://127.0.0.1:8080/test?delay_ms=300"; done
```

The ServiceMonitor scrapes every **15s** and every panel rates over a **1m** window, so
wait for **multiple 15s scrapes plus a full 1m rate window** after the traffic before
deciding whether a panel is correct. An empty result in that first minute is waiting
time, not a broken dashboard — and it is never a reason to fabricate a value or to
suppress a panel.

With only one Service forward, only one replica receives traffic; the port-forward picks
one backing pod. To prove both resource panels, drive each pod directly (step 6's
per-pod forwards) or send enough load that both replicas are involved, then compare the
two pods in the legend.

### 9.3 Persistence proof — delete own pods, keep PVCs and StatefulSets

The point is that a **saved** dashboard and **stored** history survive pod replacement
because Grafana and Prometheus persist to their own PVCs. The proof deletes exactly
**one** running pod of each and waits for the chart to replace it. It **never** deletes
a PVC, a PersistentVolume, or a StatefulSet, and never scales anything down — those
would manufacture success and can destroy the data being proven.

List the pods and read the exact names first:

```sh
kubectl -n monitoring get pods
```

**Grafana first.** Delete one observed Grafana pod by name, wait for the replacement to
be `Ready`, then reopen the dashboard you saved in 9.1:

```sh
kubectl -n monitoring delete pod <one-huy-monitoring-grafana-pod>
kubectl -n monitoring get pods -w        # Ctrl-C once the replacement is Running/Ready
```

The saved dashboard must still be there, in its folder, with its five panels — proof it
lived in Grafana's PVC, not in the old pod's memory.

**Then Prometheus.** Delete one observed Prometheus pod by name; note the current time
window before you do, so you can compare the same range afterwards:

```sh
kubectl -n monitoring get pods                     # find the Prometheus pod
kubectl -n monitoring delete pod <the-prometheus-pod>
kubectl -n monitoring get pods -w                  # wait for Running/Ready
```

After it returns, reload the dashboard over the **same time range you looked at before**
the deletion. Earlier request/error/latency history must still render — proof the
time-series came from Prometheus's PVC (`retention: 24h`), not from the pod it just
lost.

**What a failure looks like, honestly:** if the Prometheus PVC or the Grafana PVC were
never `Bound`, the StatefulSets would not have started; if history is missing only
*before* the restart but present after, that is a retention/storage issue to diagnose
with `kubectl -n monitoring describe pvc` and events — not a dashboard problem to paper
over. Do not delete PVCs or StatefulSets to force a green result. The `local-path` `RWO`
volumes are node-bound, not replicated or HA: a replacement pod only keeps its data while
its node and volume survive.

### 9.4 Unrun here

This worker authored the JSON and these steps but has **no cluster access**, so the
import, save, live panel observation and the pod-deletion persistence proof were **not**
run here. Sisyphus/owner performs them and captures the evidence.
