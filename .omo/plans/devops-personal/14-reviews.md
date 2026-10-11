# Stage 14: Operated-system rehearsal and supervisor delivery

## Governing artifacts

[Spec](../../../.scratch/assignment-3-personal/spec.md), [ticket 14](../../../.scratch/assignment-3-personal/issues/14-reviews.md), assignment Steps1–5/bonuses1–8/sections4,6,8 and overview review schedule, [execution contract](00-execution-contract.md).

## Entry gate and routing

Stages01–13 technically accepted and owner learning evidence recorded before full-system gate. Preliminary complex-critical because live rollout/restart, full-system dependencies and weak proof if presentations are pre-certified. Human dependencies: actual review timezone/time, presentation access, phone device and team integration agreement. Wednesday time is unspecified; 08:00 is internal readiness, not invented meeting.

## Allowed scope and branch

stage/huy/14-reviews from latest accepted origin/personal/huy, matching worktree. Files docs/huy/{operations-runbook,requirements,lessons,findings,versions}.md, slides/huy/{cp1,cp2,cp3,oct14,oct16,oct30}.md as actual usable Markdown decks, own captures and one useful app/main.go startup diagnostic for whole-loop rehearsal if authorized. Runtime image update follows normal source CI and env PR, never bypasses it. Separate env stage branch only for reviewed necessary config correction.

Permitted operations: own authorized VM/network/service access restart, private forwards/tunnel, own app/tracing/monitoring checks, bounded load and approved source rollout. No destructive teardown, namespace/PVC deletion, surprise worker-host reboot, token rotation, global Docker prune, component removal or broader topology/architecture change. A repair outside this scope is a bounded new worker prompt with Prometheus/owner decision, not opportunistic implementation by Sisyphus.

Commit scope: personal runbook/decks/findings/ledger and scrubbed review evidence; the separately authorized one-attribute source diagnostic and necessary reviewed env correction have distinct commits. No unrelated repairs or future presentation claims.

## Bounded substeps

### 14A: Restart/access runbook and operational check

Write exact commands from actual successful installation for starting named local VM, Tailscale/k3s authorized services, kubectl context, own app/tunnel forward and private Argo/Grafana/Jaeger access. Discover installed chart Services, not guessed names. Include how to read initial/local passwords privately without captures and how to use the verified project credential helper for repo/registry calls.

From a stopped access session, recreate own forwards/tunnel and browser access; if VM/services were actually stopped, use authorized start commands. Do not stop/reboot a peer VM merely to manufacture proof; obtain explicit restart permission first. Verify two Ready nodes, two correctly placed app pods, app/config response, scrape targets/dashboard, Argo sync and fresh trace/log/browser measurement. Prometheus/Grafana retain PVC data after own pod replacement; Jaeger history is transient and fresh requests regenerate it. Explain single-server/node-local storage availability limits.

```sh
export KUBECONFIG="$HOME/.kube/huy-k3s.yaml"
kubectl config current-context
kubectl get nodes -o wide
kubectl -n devops-app get all
kubectl -n devops-app get pods -l app=devops-app -o wide
kubectl -n monitoring get pods,svc,pvc
kubectl -n argocd get application devops-app-huy
```

Exit gate: actual system/access can be operated from runbook, no secret included, current private/public endpoints recorded and fresh telemetry works. A recording is backup, not live proof replacement.

### 14B: Full change-to-user walkthrough

If a new source change is needed to rehearse full loop, approve adding actual runtime.Version() as startup go_version attribute, preserving response behavior. Worker edits only that bounded diagnostic on prepared stage branch; Sisyphus independently checks/integrates when authorized. Follow real source push -> checks -> multiarch/SBOM/signature/smoke -> verified digest -> reviewed env PR/merge -> automatic Argo rollout -> new pod IDs on two physical hosts -> actual phone/tunnel response -> browser joined trace/log -> Web Vital -> k6/live dashboard.

Do not attribute fresh image to a docs-only commit or present empty workflow trigger as successful publication. Reverify the exact produced digest, not previous artifact's attestation. Wait rollout/scrape warmup before accepted load; keep source/env changes frozen during timed k6. Manual apply/manual-sync remain proven by earlier retained evidence, not erased by current automated state.

Exit gate: owner can conduct actual complete walkthrough and answer questions from own results. Retain backup capture with source/env/run/digest/time, not fabricated terminal outputs.

### 14C: Checkpoint and Wednesday review

Each CP deck starts with live component/evidence, then design/observed failure/fix/lesson:

- CP1 MonOct12 09:00: runnable native public image, actual CI checks/publication and bonus4, container/security explanations; freeze proof before08:15.
- CP2 MonOct12 17:00: manual proof/public request/two physical hosts plus both GitOps loops; capture buffer by16:00.
- CP3 TueOct13 15:00: own dashboard under load, real backend/browser/log/vital proofs; buffer by14:00.
- Tue evening: all13 technical outcomes and researcha–f/diagram accepted; no bonus deferred.
- WedOct14 before actual supervisor time: repeat whole-project operation and explanations; internal08:00 readiness. Record actual feedback/questions and uncovered failures, not a claim meeting occurred ahead of time.

Deck format is usable Markdown slides with slide separators and actual relative evidence links, export with an available agreed presentation tool if needed. Do not make deck completion depend on an unchosen paid editor. Owner rehearses slide rendering and live-window switching, not just file existence.

### 14D: Friday half-way and final obligations

ThuOct15 address actual Wednesday findings through bounded authorized repairs and repeat affected checks. FriOct16 deliver live presentation with running k6 and real-time Grafana; demonstrate complete personal scope and recheck every moving claim against current primary docs. Update checked date/status and say what changed since earlier evidence. Keep backup recording for connectivity interruption but label when live proof genuinely fails.

Supply personal write-up, captures and deck to the agreed combined team artifact. Record actual topic/step ownership, chosen shared cluster/demo and member presentation contribution; don't assign partners' tasks or silently impose your cluster. Personal technical completeness cannot certify every member presented. This is an explicit external team acceptance item.

By Oct26 finalize stakeholder material from feedback. Extend deck/findings and repeat moving-claim checks for FriOct30 live final presentation. Both Oct16 and Oct30 occurrences must be recorded AFTER delivery, not closed as future checkboxes. Keep incremental history/captures; don't destroy own cluster before reviews or replace it with another project.

## Verification and evidence

Sisyphus independently repeats runbook access/actual responses, artifact verification, rollout and telemetry/load paths. Review all ledger rows and source/evidence provenance, private credential boundaries, actual deck rendering and links. Capture live demonstration dates and owner contribution only when observed. Record known limits: broad PAT risk, experimental browser SDK, ephemeral tunnel, transient Jaeger, local-path RWO and single control plane. No claim enterprise production-grade resilience.

## Learning gate

Owner explains whole commit-to-user/telemetry path, compares chosen tools with researched alternatives, interprets live metrics and traces, identifies credential/signing identities and local-storage limits, and describes actual lessons. Supervisor questions become learning/repair entries, not worker guesses about what was asked.

## Completion and future gates

Wed personal readiness requires all five steps/eight bonuses already operated plus complete cited research. Stage14 has dated sub-gates; readiness and a future live presentation are different states. Technical blockers remain explicit with commands/attempts/missing prerequisite; scope/deadline changes require owner approval. Assignment final completion additionally requires combined team deliverables/member split and both actual Friday presentations.

## Sisyphus handoff prompt

```text
Coordinate Stage14/ticket14 only after Stages01–13 accepted. Prepare named stage
worktree; preliminary complex-critical for live operations. Do not implement.
Dispatch bounded runbook/deck or diagnostic-change work, respecting exact host/
namespace permissions and same project-token policy. Independently operate actual
system and verify ordered CI-to-user/telemetry/load proof plus owner learning.
Keep historical manual evidence and current artifact provenance. Record readiness
by Tue evening/Wed internal08:00, actual feedback after review, live Oct16 k6 and
future Oct26/30 updates separately. Commit only authorized docs/source/evidence;
no destruction, bonus deferral, unverified meeting claim or partner assignment.
```
