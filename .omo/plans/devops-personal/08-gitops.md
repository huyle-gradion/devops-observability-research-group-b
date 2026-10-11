# Stage 08: Reviewed digest PR and both GitOps sync modes

## Governing artifacts

[Spec](../../../.scratch/assignment-3-personal/spec.md), [ticket 08](../../../.scratch/assignment-3-personal/issues/08-gitops.md), bonus 1/research b, overview Approach 8, [execution contract](00-execution-contract.md), [GitHub credential policy](../../../docs/agents/github-access.md).

## Entry gate and routing

Stage 07 accepted with retained manual/cross-host/tunnel captures. Preliminary complex-critical: reconciliation ownership, CI write credential and cross-system deploy contract. Read current pinned [Argo install](https://argo-cd.readthedocs.io/en/stable/getting_started/), [Application schema](https://argo-cd.readthedocs.io/en/stable/user-guide/application-specification/), [auto sync](https://argo-cd.readthedocs.io/en/stable/user-guide/auto_sync/) and kubectl local set image docs. No whole-project dispatch.

## Allowed scope and branches

Source stage/huy/08-gitops from accepted origin/personal/huy; worktree per contract. Files k8s/gitops/{project,application}.yaml, k8s/bootstrap/app-namespace.yaml, removal of obsolete k8s/app source copies after handoff, .github/workflows/ci.yml, two useful app startup-log attributes, own GitOps runbook/docs/captures.

Create env/huy as a clean desired-state branch in the SAME authorized public repo with README and k8s/app/{configmap,deployment,service}.yaml only. Sisyphus prepares its authorized initial commit/branch and separate persistent ENV_ROOT, then env edits on stage-env/huy/08-gitops. Do not include Namespace in Argo-managed path; manual namespace becomes source bootstrap. No accidental copies of app source/workflow into env branch. CI proposals derive from env/huy and target it through reviewed PRs. Commit source/desired-state changes separately and record both refs/SHAs.

Authorized operations: pinned Argo install in argocd, own AppProject/Application, app devops-app reconciliation and same-token repo PR settings/branches only. No CI kubeconfig/PAT in cluster, broad cluster-resource whitelist, Argo public tunnel, another member's state or imperative app apply after handoff.

Commit scope: source bootstrap/project/Application/CI proposal and named startup attributes, with obsolete source manifest removal; separate env commits contain only accepted app desired state. Review source and env diffs independently; no credential files.

## Bounded substeps

### 08A: Install and transfer ownership without deleting the running app

Install verified pinned Argo release manifest into argocd, wait actual controller/server readiness, record resources/release. Use local UI forward svc/argocd-server 8443:443. Owner retrieves local initial password privately; don't capture Secret value.

AppProject huy permits only actual repo URL and destination in-cluster/devops-app, no cluster-scoped resource whitelist. Application devops-app-huy uses that project, repoURL actual public repo, targetRevision env/huy, path k8s/app, destination https://kubernetes.default.svc namespace devops-app. Initially omit automated sync.

Seed env branch with EXACT manually accepted digest/config/security/spread/service. Sync/adopt existing namespaced objects, prove /test unchanged, then remove duplicate app manifests from personal/huy and retain only k8s/bootstrap/app-namespace.yaml. No resource deletion needed for transfer. Record source/desired-state branch authority and stop using old manual app apply. Bootstrap and monitoring remain source-managed outside this Application.

Exit gate: healthy manual-mode Application owns existing app, env/huy is sole desired state, manual proof remains, namespace cannot be pruned by this app.

### 08B: Fully green CI proposes one reviewed digest change

Add propose-deploy job needs publish and trusted source-push only. Built-in token permissions remain {}; same .env-derived PROJECT_GITHUB_TOKEN authorizes checkout, Git HTTPS and gh PR APIs explicitly. No alternate GITHUB_TOKEN/PAT, SSH key or cluster credential. PAT scopes are not narrowed by workflow permissions; require reviewed trusted workflow, explicit project branch/file bounds and secrets only in needed steps. Owner confirms PR creation/settings permit this account.

Validate nonempty successful publish digest, source SHA and repo/package identity. Checkout current env/huy; branch deploy/huy-FULL_SOURCE_SHA. Use kubectl set image --local to change ONLY app container image in k8s/app/deployment.yaml, via temp output then reviewed replacement. No regex replacement across all manifests. Check actual desired digest first; if already same, no commit/PR. A new SHA creates one proposal. Rerun with existing branch/open PR checks it out and adds a fast-forward scoped commit if digest differs; no duplicate PR, reset, force push, branch deletion or silent conflict resolution. Report an incompatible branch/conflict for human review.

Open PR base env/huy with actual source SHA, green run URL, digest and verification evidence. Never auto-merge. Owner/Sisyphus reviews exact diff and current relevant green source before authorizing merge. If multiple source runs are in flight, explicitly choose the intended newest verified artifact; do not blindly merge an older outstanding PR. env/huy/deploy branch events do not trigger source publication; PAT-created PRs can trigger workflows, so verify branch/path filters.

Exit gate: actual green source build proposes correctly scoped digest PR; rerun does not duplicate; failed/skipped publish cannot propose.

### 08C: Manual full loop

Useful source change: add resolved app_env attribute to initial startup log. Preserve every response. Integrate through normal reviewed source path -> CI check/build/sign/smoke -> digest PR -> reviewed desired-state merge. Observe OutOfSync and unchanged running image BEFORE sync. Manually click Sync, wait app rollout, inspect new digest and new pod IDs, request /test/config header. Capture UI status/before/after/run/PR/source/env SHAs.

Exit gate: actual source-to-image-to-env-to-manual-sync-to-new-pods loop, not merely a replica edit without CI evidence.

### 08D: Automatic loop, drift and rollback

Set Application syncPolicy.automated enabled true, selfHeal true, prune true, allowEmpty false; apply ONLY bootstrap Application. Useful source change: add actual runtime.GOARCH startup arch attribute. Repeat fully green source/digest-PR/reviewed-merge; watch controller auto-sync without UI click, verify new pods/digest/response. Capture distinct loop from manual proof.

Scale own app Deployment to 3 while Git says 2; observe self-heal restoring 2. Do not mask desired-state error with kubectl apply. Revert desired-state image-change commit in Git to previous valid digest through reviewed env PR and observe controller rollback. Preserve ConfigMap/security changes from later commits; revert only intended image-change commit. Do not kubectl rollout undo against Argo.

## Verification and evidence

```sh
kubectl -n argocd get application devops-app-huy -o yaml
kubectl -n devops-app get deployment devops-app -o jsonpath='{.spec.template.spec.containers[?(@.name=="app")].image}{"\n"}'
kubectl -n devops-app get pods -l app=devops-app -o wide
kubectl -n devops-app rollout status deployment/devops-app --timeout=180s
python3 scripts/project-github.py gh pr list --repo "$REPO" --base env/huy
```

Sisyphus independently verifies adoption/ownership, actual pipeline PR, both sync modes, no default credential fallback, drift and Git rollback. Save own bonus-1 captures with source/run/digest/env commit and pod/UI linkage; update runbook and requirements. Expected default polling delays are observed, not replaced with fake immediate success. Keep admin UI private and password scrubbed.

## Learning gate

Owner reviews a digest PR, demonstrates manual/automatic distinction, creates safe drift and explains reconciliation and Git rollback. Explain push versus pull credentials, Argo versus Flux, why reviewed merge is the approval boundary, and why CI publication alone is not CD.

## Completion and CP2 gate

Steps 2/3 and bonuses 1/2 must actually work before CP2 Mon 17:00, with capture buffer before 16:00. Missing token permission/worker/agent loop is explicit blocker. Manual local digest PR can preserve an operational Argo demo if automation is blocked, but cannot be mislabeled completed CI-to-PR automation; obtain owner resolution without swapping credentials. Retain your project, never reset to a teammate's selected baseline.

## Sisyphus handoff prompt

```text
Coordinate Stage 08/ticket 08 after retained manual Step 2/3 and two-host proof.
Prepare separate source/env stage branches/worktrees. Preliminary complex-critical.
Do not implement. Dispatch bounded ownership-transfer, proposal-job, manual-loop
and automatic-loop subtasks, verifying each before next. Same project token only;
CI holds no cluster credential and PRs are reviewed, never auto-merged. Adopt the
existing app without deleting it; env/huy is the only app desired-state copy.
Independently verify actual source/CI/digest/PR/UI/pod loops, self-heal and Git
rollback plus owner's learning. Commit only authorized scope and preserve CP2
proof, private admin access and source/desired-state identity.
```
