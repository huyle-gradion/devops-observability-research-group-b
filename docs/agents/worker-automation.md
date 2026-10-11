# Owner-enabled worker dispatch

## Default and enablement

Default is owner relay: Sisyphus returns the composed worker prompt, the owner delivers it, and the owner returns the worker result. Writing plans does not enable dispatch. No driver is installed or implemented by this document.

Direct dispatch requires an explicit owner instruction for THIS session naming the worker's OpenCode-slim session, permitted stage/subtask, permitted transport and readback route. The owner signals when a new worker is needed. Do not infer a session from history, adopt another project's driver, start a new worker silently, or reuse expired session permission.

## Composition and transmission

Sisyphus reads the live Worker (Hephaestus) section from [workflow](opencode-workflow.md), appends exactly `---`, then the bounded prompt with [execution-contract header](../../.omo/plans/devops-personal/00-execution-contract.md). Supply only named plan/context scope. No .env, token, private kubeconfig or restricted content in the message or driver log. Project GitHub authentication follows [credential policy](github-access.md), never worker-session/cached account credentials.

Before sending, confirm prepared branch/worktree, base SHA, worker class/model policy, operation permission and acceptance criteria. A driver may transport and read results; it does not expand scope or bypass QA. If the approved transport is unavailable, report it and return to owner relay, not a fabricated dispatch.

## Structured handoff

Write `.omo/handoffs/huy/NN-subtask.md` before ending a directly dispatched worker turn. If an owner-approved existing transport requires a machine-readable schema, use that schema and preserve all fields below. No transport-specific schema has been established here yet.

Required fields:

- Stage/subtask and governing issue/spec/plan paths.
- Branch/worktree, base ref and SHA, authorized files/operations.
- Changed files and implementation summary.
- Producing source SHA, image digest and run URL when applicable.
- Actual commands with exit codes/results and evidence paths.
- External operations/resources created or changed.
- Failed/blocked checks, unrun checks and outstanding decisions.
- Explicit statement that the worker has not committed or performed final QA.

Scrub credentials before persistence. Do not include environment dumps or authentication headers. Handoff existence is not acceptance.

## Sisyphus readback and close

Read the actual worker-session result AND handoff artifact, reconcile discrepancies, independently verify/review, require owner learning evidence, and record acceptance. Report missing/read-failed result without declaring implementation complete. Commit/integrate only when owner authorized and checks meet the plan, including post-integration CI/deployment evidence. Reclassification and stop rules remain in the workflow.
