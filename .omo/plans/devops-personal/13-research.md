# Stage 13: Actual pipeline diagram, complete research and lessons

## Governing artifacts

[Spec](../../../.scratch/assignment-3-personal/spec.md), [ticket 13](../../../.scratch/assignment-3-personal/issues/13-research.md), ALL research3.1a–f, bonus7 and assignment4/6, overview Approach12, [execution contract](00-execution-contract.md).

## Entry gate and routing

Research/document preparation may start after Stage01; final acceptance requires Stages02–12 real evidence. Preliminary routine-verified for public-doc synthesis; escalate architecture changes, uncertain security claims or private source context. No zero-cost preference inferred. Use primary sources; external research support is allowed only AFTER plan-pack completion/execution authorization. Prometheus owns unresolved decision changes, owner owns final explanations.

## Allowed scope and branch

stage/huy/13-research from current accepted origin/personal/huy, matching worktree. Files docs/huy/{findings,lessons,requirements,versions}.md, pipeline.mmd and rendered diagram, personal research/evidence README links and presentation outlines. No technical implementation, changing acceptance, copying another project's captures or assigning partner work. Commit scope: checked citations/own learning/evidence synthesis. Integrate docs without overwriting concurrent source changes; this stage alone may overlap implementation with separate file ownership.

## Bounded substeps

### 13A: Claim registry and source checks

For each claim record topic, exact claim, primary URL/section, checked date, relevant version/status, observed confirmation/correction/uncertainty and own-system evidence where applicable. Findings are in owner's words; no invented stable-since date, quoted third-party endorsement or undocumented vendor guarantee. Make literal brief claims a verification checklist, not assertions to repeat.

Minimum checked coverage:

| Topic | Required concrete study |
|---|---|
| a CI/CD | Build/test/lint/deploy gates; Actions trigger/job/step/runner/needs/event anatomy; CI/delivery/deployment; SLSA/in-toto/Sigstore and actual attestation; Immutable Actions rollout/OCI claim independently checked; short GitLab CI/Jenkins comparison |
| b GitOps/access | Push cluster credentials versus pull controller; Argo deep dive/Flux comparison; reviewed desired state/manual/auto/drift/revert; exact source-image-digest-PR connection; cloudflared/ngrok public tunnels versus Tailscale private mesh and applicable free tiers |
| c Containers | Isolation/runtime/layers/cache; non-root/minimal/distroless/no shell; secret-layer persistence; external health versus probes; actual installed Buildx SBOM syntax and inventory; Trivy/Grype/Scout gates and scan limits; Docker Hardened Images provenance/SBOM/VEX claims verified in Docker docs, not adoption required |
| d Kubernetes | Pod/ReplicaSet/Deployment/Service/Namespace relationships and control loops; scheduling/scaling/rolling update/self-heal; ConfigMap and Secret env/file injection; base64 not encryption-by-default and git danger; emptyDir/pod filesystem/PV/PVC/RWO versus RWX; local tool storage/default provisioner/version limits; control plane/worker/single/multi-node; actual extra objects ServiceMonitor/CRDs/AppProject/Application and why |
| e Observability | Traces/metrics/logs versus plain logs; context correlation/W3C propagation/vendor neutrality; EACH selected language SDK's signal maturity; Collector current major/release/changelog and any verified config stabilization; batching/routing/fan-out explained without pretending demo exports to multiple viewers; versioned Semantic Conventions; Prometheus pull/discovery/time series/PromQL; Grafana workload plus kubelet/cAdvisor; k6 traffic/threshold behavior |
| f Browser | RUM versus synthetic strengths; PerformanceObserver/Navigation/Resource Timing; LCP/INP/CLS definitions/units and web-vitals; CrUX field aggregation versus Chrome DevTools/Trace Viewer lab; current Web SDK/fetch/XHR/document-load APIs/maturity; actual traceparent joined path; Server-Timing/PerformanceServerTiming versus full SDK |

Primary starting sources: [Actions](https://docs.github.com/en/actions), [GitHub changelog](https://github.blog/changelog/), [Docker](https://docs.docker.com/), [Kubernetes concepts](https://kubernetes.io/docs/concepts/), [k3s storage](https://docs.k3s.io/storage), [minikube storage](https://minikube.sigs.k8s.io/docs/handbook/persistent_volumes/), [kind](https://kind.sigs.k8s.io/), [k3d](https://k3d.io/), [Argo](https://argo-cd.readthedocs.io/), [Flux](https://fluxcd.io/flux/), [OpenGitOps](https://opengitops.dev/), [OTel status](https://opentelemetry.io/status/), [Collector](https://opentelemetry.io/docs/collector/), [Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/), [Prometheus](https://prometheus.io/docs/), [Grafana/k6](https://grafana.com/docs/k6/latest/), [web-vitals](https://github.com/GoogleChrome/web-vitals), [CrUX](https://developer.chrome.com/docs/crux), [DevTools](https://developer.chrome.com/docs/devtools/performance) and [W3C](https://www.w3.org/TR/trace-context/). Cite precise pages, not only these landing URLs.

Keep uncertain brief claims explicit: Immutable Releases GA does not establish Immutable Actions as OCI-on-GHCR; local provisioner RWO does not mean Kubernetes can never support RWX through another provider; namespace alone is not security isolation. Recheck exact selected/versioned defaults. No blanket unverified maturity claims. Compare current documentation to planning snapshots and record changes.

Exit gate: every required subtopic has primary check and explanation/correction; moving claims have Oct16 and Oct30 recheck entries, not future checks marked done.

### 13B: Diagram the actual operated system

Create docs/huy/pipeline.mmd from actual Stage12 state. Nodes/edges: own source commit -> Actions checks -> multiarch image/SBOM/signed attestation -> GHCR digest -> reviewed env PR -> env/huy -> Argo -> devops-app namespace across two physical hosts -> workstation Service port-forward -> public tunnel -> browser. Monitoring namespace alongside: app scrape->Prometheus->Grafana, app traces->Collector->Jaeger, browser export through bounded Go proxy to Collector; log-to-trace correlation and k6 private target. Show historical manual kubectl apply as a labeled earlier proof path, not competing current ownership.

Draw private/public boundary and credential locations: local .env/project CI secret and OIDC signing separate, private kubeconfig/join token, no credential path from CI into cluster. Do not embed token/keys/secret values or claim Argo manages monitoring when it only manages app. Names/ports must match observed deployment; local-stateful-storage/transient Jaeger limitations noted. No invented future ingress/cloud service.

Render Mermaid with available tooling or viewer and visually inspect readability/edge accuracy. Rendering is diagram proof only; it cannot replace required runtime captures. Save source and usable render alongside write-up. Exit gate: diagram matches installed components/routes/ownership and owner can trace one request/change.

### 13C: Evidence ledger and personal synthesis

Link every Step1–5/bonus1–8 to own actual capture, source SHA/run/index digest/env commit, command/device/time and observed result. Check all files exist and no token/header/credential dump; evidence-only commits can differ from producing source but must not misattribute it. Add actual problems/lessons, chosen-versus-alternative rationale and operational limits. Progressive commits throughout, not one final pasted archive.

Research and technical ledger states remain separate. Required Friday live occurrence remains future until Stage14; Tuesday rehearsal is not its substitute. Record your personal topic/step contributions and provide material for combined team write-up/shared cluster/demo/member presentations. Team selection and other members' work are an external coordination gate, not assigned here.

## Verification and evidence

Sisyphus independently compares the written coverage checklist to assignment3.1a–f and section6; follows primary citations and dated moving-status entries; checks each proof link and artifact identity; renders/reads actual diagram. Check Markdown local links and diagram syntax using available tooling without inventing successful render. Owner reviews prose against actual experience. No content/wording unit tests.

## Learning gate

Owner explains all six topics with links to own experience, walks diagram from commit to user and telemetry, compares the rejected alternatives and identifies uncertain/limited claims. Focus on reasons and observed failure/fix, not memorized product lists. Agent drafts do not count as the owner's words until reviewed/reworked and explained.

## Completion and blockers

Close by Tue evening only when all subtopics/diagram/evidence are accurate and accepted. Unknown rollout/vendor detail is labeled unverified with exact attempted sources; it cannot erase a research topic. Missing actual proof remains incomplete, not replaced by a diagram. Preserve final team/presentation obligations separately.

## Sisyphus handoff prompt

```text
Coordinate Stage13/ticket13. Start public-doc research after Stage01, but final
acceptance needs actual Stages02–12 evidence. Prepare docs-only isolated branch;
preliminary routine-verified unless risky claim/context appears. Do not implement.
Dispatch bounded a–f source checks, actual-system diagram and evidence synthesis,
using current Worker role/credential policy. Verify each claim against primary
pages and every own capture/identity; render/read diagram. Require owner review
and six-topic explanation. Commit only authorized cited docs/render/scrubbed
links; don't repeat unverified brief claims or assign partners' work. Keep Oct16/
Oct30 moving-claim and live-presentation gates explicit and future until executed.
```
