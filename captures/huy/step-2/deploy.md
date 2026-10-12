# Stage 07 — manual deploy + networking evidence (huy)

Produced by Sisyphus against the live `huy-k3s` cluster. Observed **2026-10-12**.
No secrets. Base SHA `4a18dbbd0d6f700ef73f6324f46f7208ab08afb4`.

Cluster: single-node k3s v1.37.1+k3s1 (VM on omarchy). Namespace `devops-app`.
Image: `ghcr.io/huyle-gradion/devops-app-huy@sha256:99ab49d8…` (from the 6-job push run).

## Step 2 — manual `kubectl apply`

```
kubectl apply -f k8s/app/namespace.yaml      # namespace first
kubectl apply --dry-run=server -f k8s/app/   # all four resources valid
kubectl apply -f k8s/app/                    # configmap + deployment + service created
kubectl -n devops-app rollout status deployment/devops-app --timeout=180s
  -> deployment "devops-app" successfully rolled out
```

`kubectl -n devops-app get all`:

```
pod/devops-app-6f77b9498b-g25tx   1/1   Running
pod/devops-app-6f77b9498b-scwv5   1/1   Running
service/devops-app   ClusterIP   10.43.24.138   8080/TCP
deployment.apps/devops-app   2/2   2   2
replicaset.apps/devops-app-6f77b9498b   2   2   2
```

`kubectl get nodes -o wide`: `huy-k3s  Ready  control-plane  v1.37.1+k3s1` (single node).

## ConfigMap injection and pod identity

Service `port-forward` + per-pod `port-forward`:

| Source | `/health` | `/test` body | `X-App-Env` |
|---|---|---|---|
| Service (8080) | `200 ok` | `Hello world from container-devops-app-6f77b9498b-g25tx` | `cluster-huy` |
| pod `…g25tx` | `200 ok` | `Hello world from container-devops-app-6f77b9498b-g25tx` | `cluster-huy` |
| pod `…scwv5` | `200 ok` | `Hello world from container-devops-app-6f77b9498b-scwv5` | `cluster-huy` |

`APP_ENV=cluster-huy` comes from the **ConfigMap** (not the image); each pod reports its
**own POD_NAME** via the Downward API (not the hostname).

## Self-healing

Deleted `…g25tx`; the Deployment's ReplicaSet created `…zwq2z` and readiness returned the
set to **2/2 Ready**. Controller reconciliation observed.

## Step 3 — public tunnel (owner phone pending)

`cloudflared tunnel --url http://127.0.0.1:8080` (quick tunnel) → URL:
`https://debian-performances-truly-facing.trycloudflare.com` (ephemeral). Workstation
sanity request through the tunnel returned `HTTP/2 200`. The off-machine proof is the
owner's **phone/mobile-data** request (captured separately under `captures/huy/step-3/`).

## Single-node deviation

Bonus 2 (two physical hosts, one pod per host, cross-node pod-IP/Service responses) and
the plan's hard `DoNotSchedule` spread are **deferred** until the macOS Linux VM joins. The
manifest uses `whenUnsatisfiable: ScheduleAnyway` (valid; tolerates one node). `minDomains`
is omitted because the API only allows it with `DoNotSchedule`. No two-node claim is made.
