# Personal DevOps execution plans

Start with [execution contract](00-execution-contract.md) and [personal specification](../../../.scratch/assignment-3-personal/spec.md). This is the complete detailed breakdown of `local://devops-personal-learning-plan.md`. Stage plans are self-contained with respect to architecture decisions; a later agent need not access that session-local URI.

This delivery is PLANNING ONLY. No worker was dispatched, no repository was published, no software installed, and no project/runtime acceptance claimed. Sisyphus may begin only after the owner authorizes execution. Writing every stage now avoids architecture-by-worker while keeping each worker assignment bounded.

Repository role authority is [multi-agent workflow](../../../docs/agents/opencode-workflow.md); dispatch transport is [worker automation](../../../docs/agents/worker-automation.md). [Credential policy](../../../docs/agents/github-access.md) supersedes the overview's older GITHUB_TOKEN/cached-account assumptions: .env GITHUB_ACCESS_TOKEN authenticates as huyle-gradion; trusted CI uses the same value in PROJECT_GITHUB_TOKEN. Missing workflow scope is a named execution prerequisite, not a reason to switch accounts.

## Stage order

| Stage | Plan | Depends on | Preliminary class | Deadline |
|---|---|---|---|---|
| 01 | [Workspace, prerequisites and learning ledger](01-workspace.md) | Owner authorization | complex-critical | Sun Oct 11 first work window |
| 02 | [Runtime API and behavior tests](02-runtime-api.md) | 01 | routine-verified | Sun afternoon |
| 03 | [Non-root multi-platform image](03-container.md) | 02 | complex-critical | Sun afternoon |
| 04 | [CI, GHCR publication and runtime smoke](04-ci.md) | 03 | complex-critical | Sun afternoon/evening |
| 05 | [SBOM, signed provenance and security proof](05-supply-chain.md) | 04 | complex-critical | Mon Oct 12 before 08:15; CP1 09:00 |
| 06 | [Two-physical-host k3s cluster](06-cluster.md) | 01 and actual worker access | complex-critical | Mon before deployment; preparation can start Sun |
| 07 | [Manual deployment and public tunnel](07-manual-deploy.md) | 05, 06 | complex-critical | Mon before 16:00 |
| 08 | [Both GitOps modes and CI digest PR](08-gitops.md) | 07 | complex-critical | Mon before 16:00; CP2 17:00 |
| 09 | [App metrics and persistent dashboard](09-metrics.md) | 08 | complex-critical | Mon evening |
| 10 | [Backend traces and correlated logs](10-tracing.md) | 09 | complex-critical | Mon evening/Tue morning |
| 11 | [Browser joined trace and Web Vitals](11-browser.md) | 10 | complex-critical | Tue before 14:00 |
| 12 | [Load test and live observability proof](12-load.md) | 09, 10, 11 | routine-verified | Tue before 14:00; CP3 15:00 |
| 13 | [Accurate diagram and all research](13-research.md) | Start after 01; close after 12 | routine-verified | Tue evening |
| 14 | [Rehearsal, presentations and final delivery](14-reviews.md) | 01–13 accepted | complex-critical | Wed Oct 14; Fri Oct 16; Oct 26/30 |

Each stage has bounded substeps. Dispatch ONE substep, receive its handoff, and verify before proceeding; a stage plan is not a request for one worker to race through all substeps. Source/app/workflow changes integrate serially. Parallel preparation is limited to actual independent work, notably worker access and docs research.

## Requirement-to-stage acceptance map

| Outcome | Implement/operate in | Proof/gate |
|---|---|---|
| Required Step 1 | 02–05 | Public repo, real tests/build/push, image runs; CI run/digest and response captures; CP1 |
| Required Step 2 | 06–07 | Manual kubectl apply, nodes/get all, isolated namespace, ConfigMap-visible response; CP2 |
| Required Step 3 | 07 | Actual phone/mobile-data request through current tunnel URL; CP2 |
| Required Step 4 | 09, 12 | App rate/errors/latency plus per-pod CPU/memory dashboard with actual traffic; CP3 |
| Required Step 5 | 12, 14 | Real k6 with Grafana moving; live repeat Oct 16, not substituted by Tuesday capture |
| Bonus 1 | 08 | Separate source→CI→desired-state→manual-sync and automatic-sync rollouts |
| Bonus 2 | 06–07 | Two actual physical hosts, Ready nodes, pods spread, pod-IP and Service cross-node responses |
| Bonus 3 | 10 | Actual request exported through Collector and visible in Jaeger |
| Bonus 4 | 05 | Both platform SPDX SBOMs and successful signed-provenance verification for own image |
| Bonus 5 | 10 | Actual JSON request log trace/span IDs match exported trace |
| Bonus 6 | 11 | Browser click/fetch and backend child spans share one trace and correct parentage |
| Bonus 7 | 13 | Actual-system Mermaid diagram includes CI, registry, both deploy paths, namespaces, tunnel/user and monitoring |
| Bonus 8 | 11 | Real supported LCP/INP/CLS number from own Chrome page, with context and units |
| Research 3.1a | 04–05, 13 | Actions anatomy/gates, supply-chain status, GitLab/Jenkins comparison with primary citations |
| Research 3.1b | 07–08, 13 | Argo deep dive/Flux comparison, image-update boundary, push/pull, access-tool free-tier comparison |
| Research 3.1c | 03, 05, 13 | Layers/cache/non-root/secrets/health, SBOM CLI, scanning, DHI provenance/SBOM/VEX verified claims |
| Research 3.1d | 06–09, 13 | Objects/reconciliation, Secret/ConfigMap, namespace limits, local storage/tool comparisons, cluster shape |
| Research 3.1e | 09–10, 12–13 | Signals and per-language maturity, Collector version/changelog, Semantic Conventions, discovery/PromQL/k6 |
| Research 3.1f | 11, 13 | RUM/synthetic, browser Performance APIs, CWV/CrUX/lab distinction, Web SDK and Server-Timing |
| Public Git, incremental commits, captures | 01 and every stage | Authorized scoped commits as work progresses; linked evidence tied to producing source SHA |
| App and monitoring namespaces | 07, 09–11 | devops-app and monitoring actual objects; private administrative services |
| Free tiers/no production access | 01, 06–08, 14 | Recorded accounts/host permission; no paid/cloud/production substitution |
| Wed Oct 14 full personal project | 14 | All 13 technical outcomes accepted by Tue evening and operational readiness repeated Wed |
| Fri Oct 16 and Oct 30 live decks | 13–14 | Personal deck, live system/k6; dated moving-claim recheck and actual presentation record |
| Team combined write-up/shared demo/member split | 13–14 | Owner contribution and agreed integration/presentation roles; remains external team gate, never claimed from personal proof alone |

## Dates and completion rules

CP1 Mon 12 Oct 09:00; CP2 Mon 12 Oct 17:00; CP3 Tue 13 Oct 15:00. These use the screenshot, not PROJECT_PLAN.md's conflicting CP2 time. The owner committed all bonuses by the Wed Oct 14 review. Record actual presentation timezone in Stage 01. Wednesday's review time must be confirmed; use 08:00 as internal readiness only, not an invented meeting time.

Technical PASS needs runtime evidence, independent QA and accepted source/environment integration where relevant. Learning PASS needs the owner's reproduction/explanation. A missing worker allocation, sudo authorization, OAuth or phone request is a named human dependency; no screenshot substitution or silent topology change.

## Initial Sisyphus prompt

```text
Act as Sisyphus for Assignment 3 Group B's personal DevOps project.
Read .omo/plans/devops-personal/00-execution-contract.md,
.omo/plans/devops-personal/README.md, .omo/plans/devops-personal/01-workspace.md,
.scratch/assignment-3-personal/spec.md and docs/agents/opencode-workflow.md.
The owner requested the complete plan pack before execution. Confirm execution
is now authorized; planning publication itself is not permission to dispatch.
Do not implement. Establish DOC_ROOT and the authorized SOURCE_ROOT/REPO.
Prepare the specified branch/worktree before worker editing. Use the owner's
routing table, check current risk/provider policy, and compose worker prompts
from the CURRENT Worker (Hephaestus) section of docs/agents/opencode-workflow.md
plus --- and the bounded header. Use only docs/agents/github-access.md credentials;
the authorized project account is huyle-gradion, never cached personal gh login.
Use owner relay unless direct dispatch is explicitly enabled for a named session.
Start with Stage 01's first actionable bounded substep. Record unresolved human
inputs. Verify each worker result independently, collect technical AND learning
evidence, and commit only when authorized. Do not start the next dependent stage
until its entry gates pass. Preserve all five required steps and all eight bonuses.
```

## Pack verification

Planning verification checks every referenced stage/ticket, dependency graph, internal document link, required outcome mapping, branch contract, Sisyphus prompt, learning gate and fenced shell syntax. It does not execute infrastructure commands or certify the future implementation. See [verification record](planning-verification.md) after this pack's authoring checks finish.
