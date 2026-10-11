# Stage 07: Manual deploy, real two-node traffic and public tunnel

## Governing artifacts

[Spec](../../../.scratch/assignment-3-personal/spec.md), [ticket 07](../../../.scratch/assignment-3-personal/issues/07-manual-deploy.md), Steps 2/3 and bonus 2, overview Approach 7, [execution contract](00-execution-contract.md).

## Entry gate and routing

Stage 05 accepted green verified public digest and Stage 06 two-host cluster accepted. Preliminary complex-critical: cluster deployment/security context and public exposure. Read current Kubernetes Deployment/Service/Downward API/topology-spread docs and [quick tunnels](https://developers.cloudflare.com/cloudflare-one/networks/connectors/cloudflare-tunnel/do-more-with-tunnels/trycloudflare/). Use actual context huy-k3s; no production target.

## Allowed scope and branch

stage/huy/07-manual-deploy from accepted origin/personal/huy, matching worktree. Files k8s/app/{namespace,configmap,deployment,service}.yaml, docs/huy/deploy-runbook.md, versions/lessons and own Step 2/3/bonus 2 captures. Authorized operations: app namespace devops-app, own temporary diagnostic pods and localhost forward/tunnel. No kubeconfig/PAT in Git, ingress/LB/public admin ports, Helm/GitOps yet, arbitrary Secret or other cluster namespaces. Commit scope: app manifests and manual/external proof.

## Bounded substeps

### 07A: Declarative app resources and server validation

Namespace devops-app. Deployment/Service devops-app; container app; matching stable app: devops-app label. Use verified IMAGE_REF index digest, two replicas, requests CPU100m/memory64Mi and limits CPU500m/memory128Mi. HTTP readiness/liveness /health on container 8080 with bounded initial startup/probe timings; no shell exec. runAsNonRoot true, user/group65532, no privilege escalation, read-only root, drop ALL capabilities and RuntimeDefault seccomp.

POD_NAME from fieldRef metadata.name; ConfigMap devops-app-config contains APP_ENV cluster-huy and envFrom.configMapRef. No unused simulated API key. ClusterIP Service8080->8080, port name http. Topology spread maxSkew1/minDomains2, hostname key, DoNotSchedule and same app label selector. Verify support in the actual pinned Kubernetes version; don't silently drop the two-domain requirement.

Apply namespace first so server-side dry-run of namespaced resources can succeed on a new cluster:

```sh
export KUBECONFIG="$HOME/.kube/huy-k3s.yaml"
kubectl config current-context
kubectl apply -f k8s/app/namespace.yaml
kubectl apply --dry-run=server -f k8s/app/
kubectl apply -f k8s/app/
kubectl -n devops-app rollout status deployment/devops-app --timeout=180s
kubectl get nodes -o wide
kubectl -n devops-app get all
kubectl -n devops-app get pods -l app=devops-app -o wide
```

Exit gate: two Ready app pods, one on each actual physical-host node, configured security/probes/env and expected digest. Capture manual kubectl apply now; final Argo state is not evidence that manual apply happened.

### 07B: Real networking, identity and self-healing

```sh
kubectl -n devops-app port-forward svc/devops-app 8080:8080
```

Separate terminal: curl /health and /test, observe APP_ENV cluster-huy and actual pod-name body. Record pod list alongside response. Request each pod through separate pod-specific forwards on different local ports; use observed pod names, not invented ones. Service port-forward chooses one pod, so repeated responses are NOT load-balancing proof.

Resolve/pin official curl container digest for temporary diagnostic pod. From a diagnostic pod explicitly placed on huy-worker, curl the actual app pod IP on huy-server:8080/test and devops-app.devops-app.svc.cluster.local:8080/test. Repeat with opposite placement. Use explicit namespace/nodeName and temporary names owned by the stage. Capture real responses, DNS/pod-IP result and placement; remove diagnostic pods after proof. No mock or ordinary local curl substituted for pod networking.

Delete ONE observed own app pod and observe Deployment replacement returning two Ready pods on both nodes. Do not delete Deployment/namespace/PVC or change spread constraints. Exit gate: bidirectional cross-node Pod/Service connectivity, both pod identities/config response and actual controller self-healing.

### 07C: Off-machine request

Keep the app forward alive and run:

```sh
cloudflared tunnel --url http://127.0.0.1:8080
```

Owner opens actual current HTTPS tunnel /test on phone with mobile data, not workstation curl masquerading as off-machine proof. Capture URL/status/body, device context and time without personal phone/account data. Only app is exposed. All dashboard/GitOps/trace viewer/API endpoints remain private. Do not publish Tailscale join keys or kubeconfig for demo access.

Exit gate: actual phone response through tunnel. Quick URL is ephemeral; runbook recreates tunnel/updates URL after restart. If user cloudflared config interferes with quick tunnel, diagnose actual installed config and use an isolated config path; do not delete unrelated config.

## Verification and evidence

Sisyphus independently verifies manifest dry-run, actual rollout, digest/security/env, physical-host pod spread, cross-node tests and self-heal. Owner-performed phone evidence is received and checked, not invented. Save nodes/get-all/apply/pods/config/network outputs under captures/huy/step-2/ and bonus-2/, phone/tunnel under step-3/. Keep manual evidence immutable when Stage 08 moves manifest authority.

## Learning gate

Owner applies the manifests, traces Deployment->ReplicaSet->Pod and Service selectors, explains ConfigMap versus image config/Secret, identifies both actual hosts, deletes own pod and explains reconciliation. Explain Namespace logical separation versus security/network isolation, readiness versus liveness, and public tunnel versus private mesh.

## Completion and blockers

Required Steps 2/3 and bonus 2 need every actual proof above. If one node is unavailable, fix real prerequisite; no fake two-machine claim or weaker placement. Resource errors use describe/events, then authorized budget correction. Captured failure does not silently close requested all-bonus scope. Preserve app running for Stage 08.

## Sisyphus handoff prompt

```text
Coordinate Stage 07/ticket 07 after verified artifact and two-host cluster gates.
Prepare isolated stage branch/worktree, preliminary complex-critical. Do not
implement. Dispatch bounded manifest, connectivity and tunnel tasks only in the
named app namespace/own temporary resources. Require manual apply before Argo,
real ConfigMap header/pod identity, hardened two-node placement and bidirectional
Pod/Service responses. Receive owner phone/mobile-data evidence; do not simulate
it. Independently verify/review and require learning. Commit authorized manifests
and scrubbed captures; keep deployment running and manual evidence for handoff.
```
