# Requirements evidence ledger (personal project — huyle-gradion)

Owner: huy. Created: 2026-10-11 (Stage 01, substep 01C).

This ledger is the public record of every required outcome for Assignment 3 Group B's
personal DevOps and observability project. It is honest bookkeeping: a row exists so the
outcome can be tracked, not so it can be claimed.

Rules for this table:

- One row per outcome. Every row starts with **technical state = `not-run`** and
  **learning state = `not-run`**.
- State vocabulary: `not-run` → `in-progress` → `blocked` / `failed` / `passed`.
- `passed` requires real runtime evidence **and** the owner's reproduction/explanation.
  An agent-written result or a summary never sets `passed`.
- `—` means unknown/not yet produced; `pending` means expected but not yet available.
  Never invent a source SHA, digest, run id, capture path or observation.
- Cited evidence must be scrubbed of credentials, kubeconfigs, join tokens, authorization
  headers and private model context.
- Public captures live under `captures/huy/step-N/` and `captures/huy/bonus-N/` and are
  linked from the `Evidence path` column, tied to the producing source SHA.

Source of the acceptance map: `.omo/plans/devops-personal/README.md`
("Requirement-to-stage acceptance map"). Governing specification:
`.scratch/assignment-3-personal/spec.md`. Stage plan: `.omo/plans/devops-personal/01-workspace.md`.

| Requirement | Plan/Stage | Technical state | Learning state | Producing source SHA | Artifact/run/digest | Evidence path | Observed result | Checked date | Blocker |
|---|---|---|---|---|---|---|---|---|---|
| Required Step 1 — app, image, CI, public repo | 02–05 | not-run | not-run | pending | pending | pending | — | — | `PROJECT_GITHUB_TOKEN` lacks `workflow` scope (owner) — blocks new Actions workflow push from Stage 04 |
| Required Step 2 — cluster and manual deploy into own namespace | 06–07 | not-run | not-run | pending | pending | pending | — | — | worker VM allocation pending (owner) |
| Required Step 3 — expose service outside the machine, real remote request | 07 | not-run | not-run | pending | pending | pending | — | — | worker VM allocation pending (owner) |
| Required Step 4 — app metrics + pod CPU/memory dashboard | 09, 12 | not-run | not-run | pending | pending | pending | — | — | — |
| Required Step 5 — live k6 demo with Grafana moving | 12, 14 | not-run | not-run | pending | pending | pending | — | — | presentation timezone pending (owner) |
| Bonus 1 — GitOps automate deploy (manual + automatic sync) | 08 | not-run | not-run | pending | pending | pending | — | — | — |
| Bonus 2 — multi-machine cluster, pods on separate physical hosts | 06–07 | not-run | not-run | pending | pending | pending | — | — | worker VM allocation pending (owner) |
| Bonus 3 — real trace exported to a viewer | 10 | not-run | not-run | pending | pending | pending | — | — | — |
| Bonus 4 — SBOM and/or signed build attestation | 05 | not-run | not-run | pending | pending | pending | — | — | — |
| Bonus 5 — log correlated to a trace via OTel context | 10 | not-run | not-run | pending | pending | pending | — | — | — |
| Bonus 6 — browser-to-backend joined trace | 11 | not-run | not-run | pending | pending | pending | — | — | — |
| Bonus 7 — accurate pipeline diagram matching the built system | 13 | not-run | not-run | pending | pending | pending | — | — | — |
| Bonus 8 — one real Core Web Vital from a real page | 11 | not-run | not-run | pending | pending | pending | — | — | — |
| Research 3.1a — CI/CD fundamentals and one real system | 04–05, 13 | not-run | not-run | pending | pending | pending | — | — | — |
| Research 3.1b — GitOps and continuous delivery | 07–08, 13 | not-run | not-run | pending | pending | pending | — | — | — |
| Research 3.1c — containers and production-ready images | 03, 05, 13 | not-run | not-run | pending | pending | pending | — | — | — |
| Research 3.1d — orchestration, Kubernetes core concepts | 06–09, 13 | not-run | not-run | pending | pending | pending | — | — | — |
| Research 3.1e — OpenTelemetry, Prometheus, Grafana | 09–10, 12–13 | not-run | not-run | pending | pending | pending | — | — | — |
| Research 3.1f — frontend and browser observability | 11, 13 | not-run | not-run | pending | pending | pending | — | — | — |
| Public Git, incremental commits, committed captures | 01 and every stage | not-run | not-run | pending | pending | pending | — | — | — |
| App and monitoring namespaces (`devops-app`, `monitoring`) | 07, 09–11 | not-run | not-run | pending | pending | pending | — | — | — |
| Free tiers / no production access recorded | 01, 06–08, 14 | not-run | not-run | pending | pending | pending | — | — | owner account/host permission pending (owner) |
| Wed Oct 14 full personal project review | 14 | not-run | not-run | pending | pending | pending | — | — | Wednesday review time unconfirmed (owner) |
| Fri Oct 16 and Fri Oct 30 live decks | 13–14 | not-run | not-run | pending | pending | pending | — | — | presentation timezone pending (owner) |
| Team combined write-up / shared demo / member split | 13–14 | not-run | not-run | pending | pending | pending | — | — | external team gate (owner + partners) |

## Notes

- Row count: 25 (Required Steps 1–5, Bonuses 1–8, Research 3.1a–3.1f, plus the
  public-Git/incremental-commits/captures, namespace, free-tier, checkpoint,
  live-deck and team-obligation rows).
- The `Blocker` column records only named, already-observed human dependencies. It is
  not a claim that other rows have no prerequisites.
- `Required Step 3` and `Bonus 2` share the separate-physical-host worker dependency;
  it is not satisfied by a two-VM single-host substitute.
- Moving-claim research (OTel signal maturity, Immutable Actions rollout) must be
  rechecked on Fri 16 Oct and Fri 30 Oct, not stated as settled here.
