# Personal execution contract

## Authority and roles

Read [personal specification](../../../.scratch/assignment-3-personal/spec.md), [plan index](README.md), the assigned stage, [multi-agent workflow](../../../docs/agents/opencode-workflow.md) and [credential policy](../../../docs/agents/github-access.md) before dispatch. Assignment requirements remain authoritative; approved personal decisions add all bonuses and earlier checkpoints. Historical tickets are reference, not current instructions. The workflow owns roles/routing; this contract owns stage execution mechanics.

Planning readiness is not execution authorization or runtime proof. This pack contains no implemented app or active services. Do not dispatch until the complete pack's planning verification passes and the owner authorizes the selected stage. Never authorize all stages by assuming a previous approval of the overview enabled host changes or direct transport.

## Sisyphus stage procedure

1. Read the stage, governing ticket/spec and applicable repository instructions. Confirm stage dependencies have technical PASS and learning evidence, except explicitly allowed prerequisite preparation. Read actual code if it now exists; do not assume the planning-time empty repo remains empty.
2. Resolve the stage's live values from recorded execution state, not guesses. Check owner authorization for remote/public/privileged operations and for any external model's source-context policy. Resolve hooks/transport availability; when unavailable, perform the equivalent explicit research, fresh review, security checks and decision/failure recording. Record which ran; never claim installed hooks without execution.
3. Record base ref and SHA, prepare the named branch/worktree, select final class using precedence below, and author a bounded prompt. Run @librarian and @explore when research-first is active; otherwise obtain the same external/current-doc and actual-code evidence through available tools. Codegraph has zero source at planning time; do not reinitialize it. Once code exists, use the existing index/sync if available. Tool/model availability is checked at dispatch time, not fabricated in these plans.
4. Compose the worker message from the CURRENT `## Worker (Hephaestus)` section of `docs/agents/opencode-workflow.md`, verbatim, then exactly `---`, then the bounded prompt beginning with the eight-field header. Read that section at dispatch time; no stale cached role. Deliver via owner relay unless the owner explicitly enabled direct dispatch for a named session.
5. Receive the actual worker result. Independently run the stage's runtime acceptance and focused checks, review scope/security, and require the owner's learning demonstration. If a non-Sol handoff fails verification twice or new risk appears, reclassify to complex-critical. Do not confuse implementation checks with final QA.
6. Commit exact authorized files only after checks and review pass and commit authorization exists. Integrate the reviewed source change to personal/huy via authorized merge/PR, then observe any CI/GitOps behavior specified by the stage. A stage can require pre-commit checks AND post-integration live evidence; it is not closed before both.
   The owner-authorized Stage 05 negative experiment may temporarily push a deliberately failing test revision to prove publication is skipped; its expected failure is verification evidence, not deployable acceptance. Restore a fixed green revision immediately; never promote the red revision.
7. Store result/evidence, append ticket comments and decisions/failures, and close only the completed scope. Runtime observations determine success; an unmet external prerequisite is BLOCKED, not PASS. Remove temporary resources owned by this stage after required evidence, but retain the operated project.

No worker may continue to another stage on its own. Infrastructure installation and live demos are operations substeps with explicit host/namespace permissions. Owner-performed OAuth, phone requests, model consent and learning acceptance are human gates, not tasks to fake or automate away.

## Branch and workspace contract

The supplied directory `/home/thearesia/Projects/assignment-3` is the document source and currently not a Git checkout. Planning documents added here are intentional; do not initialize it or publish unrelated files. Stage 01 establishes the implementation destination separately.

Execution record variables:

| Name | Meaning |
|---|---|
| DOC_ROOT | Supplied document directory |
| SOURCE_ROOT | Clean selected implementation clone |
| REPO | Actual owner/repository, no URL prefix |
| PROJECTS_ROOT | /home/thearesia/Projects unless owner selects another parent |
| SOURCE_REF | origin/personal/huy integration ref |
| SOURCE_SHA | Full commit that produced the image |
| IMAGE_NAME | ghcr.io/huyle-gradion/devops-app-huy, or owner-approved personal collision replacement |
| IMAGE_REF | IMAGE_NAME@sha256:verified-index-digest |
| RUN_ID | Successful source-push workflow run, not a PR-only run |
| ENV_ROOT | Separate env/huy worktree |
| KUBECONFIG | $HOME/.kube/huy-k3s.yaml, mode 0600 |
| SERVER_IP / WORKER_IP | Actual private VM tailnet addresses |
| PROJECT_GITHUB_ENV_FILE | Private absolute .env locator; local key GITHUB_ACCESS_TOKEN, account huyle-gradion |
| PROJECT_GITHUB_TOKEN | Actions secret populated from that same .env token; never a literal value in state |

After Stage 01, SOURCE_ROOT is not DOC_ROOT. In SOURCE_ROOT, Sisyphus prepares branch `stage/huy/NN-slug` from the latest dependency-accepted `origin/personal/huy`, worktree `$PROJECTS_ROOT/devops-worktrees/NN-slug`. Record exact base SHA in the handoff. Do not reset an existing branch/worktree or take over uncommitted user changes; use a fresh suffix and record it. No worker implementations on main, master, personal/huy, env/huy, or another member's branch. No force push, history rewriting, broad git add, or unauthorized merge.

For env changes, prepare `stage-env/huy/NN-slug` from current `origin/env/huy` in a separate worktree. Reviewed PRs target env/huy. CI proposals use `deploy/huy-<full-source-sha>`. A worktree cannot check out a branch already checked out elsewhere; use distinct stage branches and leave the persistent ENV_ROOT as the reviewer view.

Source stages serialize integration because they share app/workflow files. Stage 06 prerequisite preparation may overlap CI once Stage 01 has established identity and authorization, but it cannot deploy unverified images. Stage 13 research can proceed alongside technical stages using its own docs-only branch, integrated without replacing others' changes. No parallel shared-file edits without a named integration owner.

When deploying remotely, select context huy-k3s and explicit namespace. Confirm each target machine and that it is an authorized isolated allocation. Never install over an existing cluster, overwrite unrelated kubeconfig, or delete data/PVCs to obtain green status.

## Worker classification

Use the authoritative routing table in [multi-agent workflow](../../../docs/agents/opencode-workflow.md). Owner-provided model names are checked at dispatch, not availability claims. Assignment free tiers do not imply free-safe model selection. Provider consent and credential exclusion apply regardless of routing.

## Bounded worker prompt structure

After the role and separator, Sisyphus supplies:

```text
WORKER CLASS: selected final class
PRESET: selected approved preset
MODEL: checked model identifier
VARIANT: configured value or —
RISK BASIS: stage risks and any newly found risks
GOVERNING ARTIFACT: personal issue path plus stage plan path
SCOPE: exact files and authorized operations for this subtask
VERIFICATION: focused implementation commands; runtime acceptance for Sisyphus

Task: one bounded subtask, not the whole project.
Governing artifact: current spec, issue, stage, accepted base SHA.
Allowed scope/files: exact source/env branch and files, host/context/namespace if relevant.
Required behavior / acceptance criteria: numbered observable requirements from the stage.
Forbidden changes: stage exclusions plus protected refs/resources.
Constraints / instructions: relevant repo skills, compatible pinned versions, no secret leakage.
Focused verification: applicable commands and evidence paths; distinguish actual results from unrun checks.
Required handoff: changed files, checks/exit codes, runtime evidence, versions/digests, blockers.

Subagents may be used at the worker's discretion; accountability stays with the worker.
No commits.
No final QA claims.
Stop and report on ambiguity or out-of-scope needs.
When direct dispatch is enabled, write the structured handoff artifact before ending the turn.
```

## Handoff and execution evidence

Owner-relay default, enablement, composition, structured handoff fields and readback rules are authoritative in [worker automation](../../../docs/agents/worker-automation.md). No driver or direct-dispatch session is established by this planning pack. Missing worker result cannot be called complete.

Sisyphus owns `.omo/execution/huy/status.md` with separate technical, learning, review, commit and post-integration evidence states. Keep authorization records local, with no secrets. Public captures go under `captures/huy/step-N/` and `bonus-N/`, linked from `docs/huy/requirements.md`; cite the source SHA even when evidence commits are later. Persist decisions/failed approaches under `.omo/decisions/huy/` and `.omo/failed-approaches/huy/`, or the configured memory tool if actually available. Do not assume OrionBooks scripts exist here.

