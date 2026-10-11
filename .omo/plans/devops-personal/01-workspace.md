# Stage 01: Workspace, prerequisites and learning evidence

## Governing artifacts

[Personal spec](../../../.scratch/assignment-3-personal/spec.md), [ticket 01](../../../.scratch/assignment-3-personal/issues/01-workspace.md), assignment sections 4/6, and [execution contract](00-execution-contract.md). This covers overview Approach 1–2 and starts the evidence discipline used by every later stage.

## Entry gate and routing

Owner authorizes this stage and names the permitted account/destination. No existing source implementation is assumed. Preliminary class complex-critical: public publication, credentials and privileged service access. Human gates are sudo, GitHub ownership/public visibility, timezone, remote allocation request and model-policy consent. Sisyphus owns final classification.

## Allowed scope and branch

Only new selected SOURCE_ROOT: README.md, .gitignore, AGENTS.md instruction pointer if absent, plan pack/personal tickets/spec, assignment document/GLOSSARY, docs/agents guidance including workflow/automation/credential policy, scripts/project-github.py, docs/huy/{requirements,versions,lessons,findings}.md and contract execution-state paths. Copy only named planning/context artifacts from DOC_ROOT, never .env, .codegraph, plugin files, other notes or historical tracker. The tiny credential helper is the only code in this setup stage; it implements the documented authentication boundary, not product features.

Sisyphus uses the authorized project account huyle-gradion through docs/agents/github-access.md. Choose existing authorized public team repo if its actual URL is supplied; otherwise proposed huyle-gradion/devops-observability-research-group-b. Do not use cached huyhuyhuy8s credentials or infer repository existence. Clone into a clean sibling, not DOC_ROOT. An empty repo may receive one authorized README seed on main, a bootstrap exception only; then personal/huy and stage/huy/01-workspace from its seed SHA. For a nonempty repo, branch from actual authorized default ref and preserve contents. Record base ref/SHA; worktree PROJECTS_ROOT/devops-worktrees/01-workspace. Commit scope: named personal planning/context/ledger files plus reviewed credential helper.

Forbidden: application code, installers executing without permissions, initialization of DOC_ROOT, config/preset changes, publishing credentials, modifying another member's files, force pushes or blanket directory copies.

## Bounded substeps

### 01A: Authorized prerequisite operations

Owner or specifically authorized operator starts the installed Docker service, then proves actual API/run access. Use explicit sudo docker for this assignment; do not chmod the socket or automatically change persistent group membership.

```sh
sudo systemctl start docker
sudo docker version
sudo docker run --rm hello-world
# GitHub identity is checked through the credential loader below, not cached gh auth.
```

Do not recite auth tokens in output. Known failure was Docker inactive/socket denial; these commands are the authorized fix and post-fix exercise, not a repeat diagnosis to dispute that observation. Exit gate: engine responds and its hello-world binary runs. If sudo is unavailable, record exact missing access and continue only independent docs/hosted work; local runtime stays blocked.

Inventory git, go, docker/buildx, gh, curl, node/npm, cloudflared, Tailscale. Record actual versions and architecture. Install kubectl/Helm/Lima/QEMU/k6 only when their later stage requires them and permission exists. Use official package/release checksums and stable versions, no curl-to-root installer accepted without reviewing/pinning. Remote worker access is requested now, not after CP1: actual VM access, distinct physical host, resources and install/tailnet permission. Do not assign partners' projects.

Before authenticated GitHub operations, implement/verify scripts/project-github.py from the credential policy on a local prepared bootstrap branch/worktree. Sisyphus may prepare local Git refs before remote authentication; no publication occurs until identity checks pass. Require .env GITHUB_ACCESS_TOKEN, safe data parsing without sourcing/executing file contents, child GH_TOKEN override and wrong-account/missing-token failure. For Git use process-scoped HTTPS gh credential helper; never global token persistence. Registry login uses stdin and temporary private Docker config. Focused helper tests use synthetic tokens/child commands and assert rejection/secret non-disclosure, not real credentials in fixtures. Sisyphus's live smoke is authenticated read-only /user -> huyle-gradion. Do not print raw .env.

The observed token lacks workflow scope. Owner must resolve that on the same project account before Stage 04 can push new workflow files. Record excessive admin/delete scopes as risk; no agent rotation or fallback. Owner-authorized setup of Actions secret PROJECT_GITHUB_TOKEN uses the same token via stdin; verify secret NAME existence only, never retrieve/print its value. Do not place it in PR jobs.

### 01B: Destination and branch isolation

Sisyphus establishes authorized public repo, clean clone, protected refs and stage worktree. Check actual instructions in an existing repo first. Import the named plan/spec/ticket/context files so future sessions resolve local links without local:// access. If AGENTS.md is absent in the destination, add a minimal pointer to the personal spec, plan index and local tracker conventions; if it exists, preserve it and ask about conflicts rather than overwriting.

Exit gate: owner can open the public repo; worker's branch/worktree/base are recorded; personal/huy exists; no implementation is on a shared/protected branch. Commit/merge are Sisyphus's authorized actions, not a worker action.

### 01C: Evidence and learning ledgers

Create docs/huy/requirements.md with separate rows for Steps 1–5, bonuses 1–8, research a–f and team/presentation obligations. Fields: requirement, plan/stage, technical state, learning state, producing source SHA, artifact/run/digest, evidence path, observed result, checked date and blocker. Initial state is not-run, not passed.

versions.md records tool versions/architecture, later pins/digests/chart values and source URLs. findings.md is organized a–f with honest pending sections, not fabricated findings. lessons.md records expected/observed/explanation/evidence and actual errors. Put secrets/local kubeconfig/env files and transient node_modules/build output in .gitignore. Do NOT ignore captures or committed reproducibility artifacts such as lockfiles/dashboard JSON. Empty capture folders need no .gitkeep; create them when evidence exists.

Exit gate: every outcome has a ledger row, no secret-bearing files are staged, and links resolve in SOURCE_ROOT.

## Verification and evidence

Sisyphus independently observes Docker run, public repo/branch existence, exact worktree/ref and ledger/context portability. Commands below are execution-only, run from the selected worktree:

```sh
git rev-parse --show-toplevel
git branch --show-current
git status --short
python3 scripts/project-github.py gh repo view "$REPO" --json nameWithOwner,url,isPrivate,defaultBranchRef
```

Expected: correct isolated worktree and stage ref, public=false for isPrivate, known authorized changes only. Save scrubbed Docker/tool/architecture results and repo URL to captures/huy/step-1/setup.md; never include credential output. Record the meeting timezone without inventing Wednesday's time.

## Learning gate

Owner runs the Docker exercise and explains why Docker API/group access is root-equivalent, what is safe to publish, and why the personal integration branch is not a worker editing branch. Owner identifies all three local checkpoint times and the separate-host dependency.

## Completion and blockers

Technical gate passes only for observed local Docker/repo/branch/ledger results. Remote access may remain explicitly pending for Stage 06; do not label bonus 2 complete. Missing local permission/public repo access is recorded with exact scope, not replaced by an imaginary result. Do not install the whole stack as a shortcut.

## Sisyphus handoff prompt

```text
Coordinate Stage 01 using .omo/plans/devops-personal/01-workspace.md and
.scratch/assignment-3-personal/issues/01-workspace.md. Read the execution contract.
Do not implement. Confirm owner-authorized account/destination and operations.
Prepare the named isolated branch/worktree before worker edits. Preliminary class
complex-critical. Separate owner sudo/auth operations from docs-only worker work.
Dispatch 01B/01C only within their named destination files, using the current Worker
role plus separator and your eight-field bounded prompt. Preserve DOC_ROOT and
all existing user/partner changes. Independently verify public repo, branch safety,
Docker runtime proof and ledger portability; require the owner's learning gate.
Commit/integrate only when authorized. Return actual evidence and pending worker
allocation. Do not begin application work until Stage 01's required gates pass.
```
