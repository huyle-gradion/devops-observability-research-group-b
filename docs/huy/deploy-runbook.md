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
