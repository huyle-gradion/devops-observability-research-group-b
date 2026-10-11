# Assignment 3 personal execution specification

Status: ready-for-agent

## Authority

Governing assignment: [Assignment 3 Group B](../../assignment-3-group-b-devops-observability-research.md), sections 3.1, 3.2, 4, 6, and 7. Personal execution overview: `local://devops-personal-learning-plan.md`. Agent ownership follows the owner's conversation approving Prometheus planning, Sisyphus execution control/QA, and worker implementation.

The overview is a high-level reference. Repository-persistent stage plans are indexed at [the personal plan pack](../../.omo/plans/devops-personal/README.md). This specification records personal decisions, not a replacement for the assignment's research requirements. No implementation is authorized or claimed by writing these documents.

## Personal decisions

- One personal complete project, retained across checkpoints. Partners manage their own versions. A separate-host worker allocation is the only partner dependency planned here.
- Complete required Steps 1–5 AND bonuses 1–8 before Wednesday Oct 14, 2026. CP1 Monday Oct 12 09:00; CP2 Monday Oct 12 17:00; CP3 Tuesday Oct 13 15:00. All times use the owner's local presentation timezone, to be recorded at execution. Wednesday review time is unknown; Tuesday evening completion and Wednesday 08:00 internal readiness are conservative targets.
- Thursday repair/rehearsal; live half-way review Friday Oct 16; final stakeholder material Oct 26 and live final presentation Friday Oct 30. No bonus is deferred to after Oct 16.
- Public repository, own machines, free tiers, no company production access or managed Kubernetes. No paid services, enterprise ingress, paging, production RUM or policy-enforced supply-chain deployment gate.
- Go standard-library API, native amd64 workstation, multi-platform amd64/arm64 image. The assignment's Apple/M1 assumption does not describe this workstation; report observed architecture and never claim an unperformed arm64 run.
- `GET /health` -> 200 `ok\n`; normal `GET /test` -> 200 `Hello world from container-<runtime-id>\n`. Resolve identity once from nonempty POD_NAME or real hostname. Kubernetes pod name is not a CRI container ID. APP_ENV defaults to local, is exposed by X-App-Env, and is injected from ConfigMap in cluster.
- Bounded `/test` query controls provide real latency/errors. Later add `/metrics`, exact root browser page, and exact POST `/otel/v1/traces`; no legacy /healthz or /api/work aliases.
- Non-root distroless static image, public GHCR, SHA tags and deployment by digest, inspected SBOMs and signed provenance. Trivy is a reported learning exercise, not an invented enforcement policy.
- Authenticated GitHub/GHCR access uses only the project account huyle-gradion via local .env GITHUB_ACCESS_TOKEN. Trusted CI uses PROJECT_GITHUB_TOKEN secret populated from that same value. No cached personal identity or built-in GITHUB_TOKEN fallback. See [credential policy](../../docs/agents/github-access.md); the missing workflow scope blocks new workflow pushes until owner resolution. OIDC signing remains separate from API authentication.
- Own k3s server Ubuntu VM plus isolated authorized worker VM on another physical host. Private Tailscale connectivity. Single-server SQLite, not HA. Namespaces devops-app, monitoring, argocd; context huy-k3s. No local two-VM substitute for bonus 2.
- First manually deploy with kubectl apply; subsequently make env/huy the only app desired-state branch. Argo CD manual and automatic loops; CI proposes a reviewed digest PR, never holds cluster credentials.
- Prometheus client metrics plus kubelet/cAdvisor pod resources, kube-prometheus-stack, PVC-backed Prometheus/Grafana. Collector OTLP HTTP to Jaeger transient storage; standard JSON logs correlated through OTel context.
- Minimal real HTML/JavaScript frontend using Web SDK and web-vitals, joined browser/backend trace, real Chrome measurement. Label browser SDK maturity and controlled-session measurement limitations accurately.

## Agent and learner gates

Roles, routing and live worker-role composition are authoritative in [multi-agent workflow](../../docs/agents/opencode-workflow.md); session-scoped transport follows [worker automation](../../docs/agents/worker-automation.md). Prometheus plans/documents, Sisyphus owns execution/QA/authorized commits, workers implement. Owner learning and operation approvals remain separate gates.

Every stage has technical and learning acceptance. The owner must reproduce an operation and explain observations. Agent-written findings do not certify understanding. No credentials, kubeconfigs, join tokens, authorization headers or private model context in public captures.

## Historical references

The original `.scratch/assignment-3/` map/tickets remain historical context. Their macOS/k3d setup, /healthz, fake Secret, main/latest CI triggers and optional-after-halfway bonuses conflict with confirmed personal decisions. Do not execute those alternatives for this personal project or change partners' historical material. PROJECT_PLAN.md is reference only. New personal tickets belong in this directory's issues/.

## Acceptance boundaries

Required technical and research outcomes are traced individually in the plan index. Every required step and bonus needs the owner's real run/response/capture, not another project, a diagram substituted for runtime proof, or an assumed output. Research covers ALL sections 3.1a–f with primary citations and moving-claim rechecks on Oct 16 and Oct 30.

Personal completion does not imply whole-team completion. The combined write-up, shared cluster/demo selection, recorded member split, every member presenting, and both Friday live presentations remain team obligations. Record the owner's contribution and obtain integration agreement; do not assign partners' work.

## External prerequisites

Unverified worker VM address, SSH/console access, installation permission, resource allocation, different physical host and tailnet permission must be supplied before remote operations. Docker was observed inactive/inaccessible; fix authorized access without chmod on its socket. No implementation worker may improvise a replacement topology, broaden credentials, or declare these outcomes passed without evidence.
