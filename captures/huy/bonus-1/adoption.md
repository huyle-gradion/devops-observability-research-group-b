# Stage 08 — GitOps adoption evidence (huy)

Observed **2026-10-12**. Partial Stage 08 (08A done; 08B/C/D pending).

- **Argo CD:** v3.5.4 (pinned release `manifests/install.yaml`) in namespace `argocd`; all
  controller/server/repo-server/redis/dex/notifications pods Running. Admin password read from
  the `argocd-initial-admin-secret` Secret **privately** (never printed).
- **AppProject** `huy`: sourceRepo = the project repo; destination = in-cluster `devops-app`;
  `clusterResourceWhitelist: []` (no cluster-scoped resources).
- **Application** `devops-app-huy`: repoURL = the public repo, `targetRevision: env/huy`,
  `path: k8s/app`, destination `https://kubernetes.default.svc` / `devops-app`; **manual**
  sync (no `automated` block).
- **Desired-state authority:** `env/huy` @ `35c31c5` (README + `k8s/app/{configmap,deployment,service}.yaml`
  only). The source branch keeps only `k8s/bootstrap/app-namespace.yaml` + `k8s/gitops/*`.
- **Adoption:** the Application reported `OutOfSync` / `Healthy` against the manually-applied
  objects, then a **manual sync** reconfigured ConfigMap/Service/Deployment → `Synced` / `Healthy`.
  `/test` unchanged: `200`, `X-App-Env: cluster-huy`, `Hello world from container-…-scwv5`.

No cluster credentials in CI; the Argo admin password and any forward are local-only and not
captured.
