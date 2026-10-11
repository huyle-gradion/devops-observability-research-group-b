# Research findings (Assignment 3.1a–3.1f) — personal project

Owner: huy. Created: 2026-10-11 (Stage 01, substep 01C).

This file is organized by the assignment's research topics 3.1a–3.1f. It is a
placeholder structure: no findings, citations, versions or vendor claims have been
authored yet. Every claim in the final write-up must be tied to a page that was
actually checked, and moving claims (OTel signal maturity, Immutable Actions rollout)
must be rechecked on Fri 16 Oct and Fri 30 Oct rather than stated as settled.

Do not invent findings or citations. Fill each section during the cited stage(s).

## a) CI/CD fundamentals & one real system

pending — to be completed in Stage 13 (and cited stages 04–05).

Scope when authored: what a CI/CD pipeline is; build/test/lint/deploy stages and what
each gates; GitHub Actions workflow YAML anatomy (triggers, jobs, steps, runners);
2025–2026 supply-chain additions (Immutable Actions rollout status, Artifact
Attestations in SLSA/in-toto with Sigstore); a short GitLab CI and Jenkins comparison.
All claims cited to primary docs; rollout status checked against the GitHub changelog.

## b) GitOps & continuous delivery

pending — to be completed in Stage 13 (and cited stages 07–08).

Scope when authored: push vs pull/GitOps models; Argo CD and Flux as the standard
agents (pick one and go deep); git as single source of truth; reviewable/rollbackable
desired state; drift detection/correction; manual vs automatic sync; how the image tag
is updated after a CI build; cloudflared/ngrok/Tailscale for reaching a service from
outside, and which fits a public demo vs a team-internal service (free tiers).

## c) Containers & production-ready images

pending — to be completed in Stage 13 (and cited stages 03, 05).

Scope when authored: what a container is and why it standardizes deployment; Dockerfile
layers and caching; production readiness (distroless/minimal base, no secrets in
layers, health checks, build-time SBOM, vulnerability-scan gate); the `--sbom` vs
`buildx` flag change against the installed version; Docker Hardened Images and what it
really includes (provenance, SBOM, VEX) checked in Docker's own docs.

## d) Orchestration — Kubernetes core concepts

pending — to be completed in Stage 13 (and cited stages 06–09).

Scope when authored: why an orchestrator (self-healing, scheduling, scaling, rolling
updates); Pod, ReplicaSet, Deployment, Service, Namespace, ConfigMap/Secret, Volume/PVC
and how they relate; why a Secret is only base64-encoded by default and what that means
for storing one in git; pod-lifetime volume vs PVC; local-cluster storage specifics and
limits (minikube hostPath, k3s `local-path`, kind/k3d `standard`; ReadWriteOnce-only);
the declarative model; control plane vs worker nodes and single- vs multi-node; the
chosen local cluster tool and why.

## e) Observability with OpenTelemetry, Prometheus & Grafana

pending — to be completed in Stage 13 (and cited stages 09–10, 12).

Scope when authored: OpenTelemetry and its three signals; how it differs from reading
raw logs; honest per-signal maturity from the OTel status page; the Collector and what
it adds over a raw SDK export; current Collector major version and recently stabilized
features (e.g. file-based config) from the changelog; Semantic Conventions versioning;
Prometheus pull-based scraping, series, PromQL; Grafana dashboards; service discovery in
Kubernetes; a useful workload dashboard (rate, errors, latency, per-pod CPU/memory);
kube-prometheus-stack; app-emitted vs cluster-provided (kubelet/cAdvisor) metrics; k6
for realistic load.

## f) Frontend & browser observability — the same signals, from the user's tab

pending — to be completed in Stage 13 (and cited stage 11).

Scope when authored: Real User Monitoring vs synthetic/lab testing; browser Performance
APIs (`PerformanceObserver`, Navigation Timing, Resource Timing); Core Web Vitals
(LCP/INP/CLS) and why Chrome made them standard; the `web-vitals` JS library; CrUX
field data vs DevTools Performance / Trace Viewer lab data, and which tool gives which;
the OpenTelemetry Web SDK (`@opentelemetry/sdk-trace-web` + auto-instrumentation for
fetch/XHR/document-load) and W3C Trace Context (`traceparent`) linking browser click,
backend and response; the Server-Timing API and `PerformanceServerTiming`. Verify
current package/API names against the OTel JS docs, not old blog posts.
