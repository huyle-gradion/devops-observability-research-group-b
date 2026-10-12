# Stage 07 teachback — manual deploy, config injection and public exposure (huyle-gradion)

Post-implementation teaching. Written 2026-10-12. Evidence: `captures/huy/step-2/deploy.md`,
`captures/huy/step-3/tunnel.md`, `docs/huy/deploy-runbook.md`.

## 1. What Stage 07 delivered

`k8s/app/{namespace,configmap,deployment,service}.yaml` applied **manually** with `kubectl`
into namespace `devops-app`: a hardened 2-replica Deployment (digest-pinned image, ConfigMap
`APP_ENV=cluster-huy`, Downward-API `POD_NAME`, HTTP probes, non-root/read-only/ALL-drop),
a ClusterIP Service, real per-pod identity responses, controller **self-healing**, and an
**off-machine phone request** through a `cloudflared` quick tunnel.

## 2. What went wrong / friction (real)

- **`minDomains` conflict:** the dispatch asked for `whenUnsatisfiable: ScheduleAnyway` **and**
  `minDomains: 2`, but Kubernetes only allows `minDomains` with `DoNotSchedule` — the API would
  reject it. The worker caught this and used `ScheduleAnyway` (no `minDomains`).
- **Single node can't satisfy a hard 2-replica spread:** with `maxSkew 1` + `DoNotSchedule`,
  the 2nd replica on one node is `Pending`. `ScheduleAnyway` (soft) is required until the second
  host exists. Bonus 2 stays deferred — no two-node claim.
- **The cluster VM died** between stages (QEMU process gone); the k3s state persisted on disk, so
  relaunching the VM brought the node back `Ready`.
- **`pkill -f '…'` self-match:** a `pkill -f` matched its own shell command line and killed the
  shell. Use `pkill -x <name>` or kill by pid.

## 3. Lessons / assignment-required knowledge (Steps 2/3, research 3.1d/b)

- **Deployment → ReplicaSet → Pod.** A Deployment owns a ReplicaSet which owns Pods; deleting a
  Pod makes the ReplicaSet create a replacement (observed: 2/2 → delete → new Pod → 2/2). That is
  the declarative control loop in action.
- **Service selectors.** A Service is a stable ClusterIP that load-balances over the Pods matching
  its `selector` (`app: devops-app`); you need it because Pod IPs change.
- **ConfigMap vs image vs Secret.** `APP_ENV` lives in a **ConfigMap** and is injected via
  `envFrom`, so config changes without rebuilding the image (`X-App-Env: cluster-huy` proves it).
  A **Secret** is only base64-encoded by default, which is why we don't put credentials in git;
  we injected none here.
- **Downward API.** `POD_NAME` comes from `fieldRef: metadata.name` — the pod name, not the
  hostname/container ID — and each pod reports its own name.
- **Namespace = logical partition, not a security boundary by itself.** `devops-app` separates
  naming/RBAC/quota scope; it does not, by default, isolate network traffic. Real isolation needs
  NetworkPolicies/RBAC.
- **Readiness vs liveness probes.** Readiness gates Service endpoints (don't send traffic yet);
  liveness restarts a stuck container. Both use HTTP `/health` (no shell — distroless).
- **Hardened `securityContext`.** `runAsNonRoot`, `runAsUser/Group 65532`, `readOnlyRootFilesystem`,
  `allowPrivilegeEscalation: false`, `drop ALL`, `seccompProfile: RuntimeDefault` — each closes a
  class of container escape/escalation.
- **Public tunnel vs private mesh.** `cloudflared` quick tunnel = a public HTTPS URL to one service
  (fits a demo); Tailscale = a private mesh (fits a team-internal service). Only the app was exposed.
- **Manual before GitOps.** `kubectl apply` by hand first (this stage), so the GitOps agent (Stage 08)
  can be shown to reproduce the same state — a manual apply now is *not* evidence GitOps happened.

## 4. How to run / inspect it

```sh
export KUBECONFIG=$HOME/.kube/huy-k3s.yaml
kubectl apply -f k8s/app/namespace.yaml && kubectl apply -f k8s/app/
kubectl -n devops-app rollout status deployment/devops-app
kubectl -n devops-app get pods -o wide
kubectl -n devops-app port-forward svc/devops-app 8080:8080 &   # then curl
curl -i http://127.0.0.1:8080/test        # X-App-Env: cluster-huy + pod name
cloudflared tunnel --url http://127.0.0.1:8080   # open the URL on a phone
```

**What you can do with it:** reach your app on a public URL and prove config/identity end to end;
watch self-healing; this is the state Stage 08 hands to Argo CD.

## 5. Fastest rebuild runbook

1. Namespace first; then ConfigMap; then Deployment (labels `app: devops-app`); then Service.
2. Deployment: `replicas: 2`, image by digest, `envFrom` the ConfigMap, `POD_NAME` from the
   Downward API, HTTP `/health` readiness+liveness, resources, hardened `securityContext`.
3. `topologySpreadConstraints`: `ScheduleAnyway` on a single node; switch to `DoNotSchedule`
   (+ optional `minDomains`) only when ≥2 nodes exist.
4. `kubectl apply --dry-run=server` then apply; `rollout status`; verify per-pod `/test` and
   `X-App-Env`; delete a pod to show self-heal.
5. `cloudflared tunnel --url http://127.0.0.1:8080`; open the URL off-machine (phone).

**Traps:** `minDomains` with `ScheduleAnyway` (invalid); a hard spread on one node (Pending 2nd
pod); putting config in the image instead of a ConfigMap; exposing dashboards/APIs on the public
tunnel; treating a manual apply as GitOps; `pkill -f` matching your own shell.
