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
