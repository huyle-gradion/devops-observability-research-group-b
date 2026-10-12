# Reviewed GitOps digest loop — env/huy (huy)

How a green source push becomes a reviewed desired-state change, how the human
approval boundary sits in the middle, and how pull-based Argo CD turns a merged
desired state into running pods. Companion docs: [`cluster-runbook.md`](cluster-runbook.md)
(the k3s VM), [`deploy-runbook.md`](deploy-runbook.md) (the original manual apply).

## Two branches, two authorities

| Branch | Authority | Contents |
|---|---|---|
| `personal/huy` | application **source** | `app/`, `web/`, `k8s/bootstrap/app-namespace.yaml`, `k8s/gitops/*`, `.github/workflows/ci.yml` |
| `env/huy` | **desired state** for the app | `README.md` + `k8s/app/{configmap,deployment,service}.yaml` only |

Argo CD's `Application` `devops-app-huy` watches `env/huy`, path `k8s/app`, and
reconciles the `devops-app` namespace. The source branch no longer carries a copy
of the app manifests, so `env/huy` is the single source of truth for the app.
**CI never writes `env/huy`; it only proposes.**

## The loop

```
source push to personal/huy
  -> lint -> compile -> test -> build -> smoke -> publish   (green gate)
  -> propose-deploy: opens/updates deploy/huy-<source-sha>  (a PR, not a merge)
  -> HUMAN reviews the exact diff and merges into env/huy    (approval boundary)
  -> Argo CD sees env/huy changed, syncs the app
       - manual mode: a human clicks Sync
       - automated mode: the controller syncs on its own
  -> new pods roll out on the new image digest
```

`propose-deploy` runs **only** when: the event is a `push` to `personal/huy`, and
the whole `lint..publish` chain succeeded. It never runs on a pull request, and a
failed or skipped `publish` skips it automatically (`needs` semantics). It carries
the project PAT to create one branch/PR and holds **no cluster credential** — no
kubeconfig, no SSH key, no service-account token.

## How to read a digest PR diff

Open the PR against `env/huy` and look at the changed file:

```sh
REPO=huyle-gradion/devops-observability-research-group-b
gh pr list  --repo "$REPO" --base env/huy --state open
gh pr diff  --repo "$REPO" <number>
```

A well-formed digest PR changes **exactly one line** in
`k8s/app/deployment.yaml` — the `app` container's image — from the previous
digest to the new one:

```diff
       containers:
         - name: app
-          image: ghcr.io/huyle-gradion/devops-app-huy@sha256:<old>
+          image: ghcr.io/huyle-gradion/devops-app-huy@sha256:<new>
```

Check it before merging:

- `gh pr diff` shows one changed line, one file (`k8s/app/deployment.yaml`).
  Anything else (reordered fields, `strategy: {}`/`status: {}` added, comments
  dropped, other manifests touched) means the proposal was not scoped and must
  not be merged.
- The digest is a **full index digest**, not a tag; tags are never used as
  desired state.
- The PR body names the source commit, the green run URL, and the digest, and
  links the verification path. Confirm the digest in the PR equals the digest the
  green run published:

  ```sh
  gh attestation verify "oci://ghcr.io/huyle-gradion/devops-app-huy@<digest>" --repo "$REPO"
  ```

- If several source runs are in flight, do **not** blindly merge an older
  outstanding PR: choose the newest verified artifact you actually intend to
  deploy and check its digest against the run.

## Why the reviewed merge is the approval boundary

CI's job ends at "here is a candidate digest, with evidence." A merge is a
deliberate human action on a reviewed diff. That is the point at which a build
artifact becomes desired state, and it is the only write to `env/huy`. Because
the diff is one reviewed line, the reviewer can see exactly what will change in
the cluster before it changes. Auto-merge is deliberately absent: `propose-deploy`
never calls `gh pr merge`, so no green run can promote itself.

## Push vs pull credentials

- **Push side (CI):** `propose-deploy` authenticates to GitHub with the project
  PAT (`PROJECT_GITHUB_TOKEN`, from the same `.env` value; the CI secret is the
  same credential, not a new identity). It uses it for Git HTTPS via a
  process-scoped `http.https://github.com/.extraheader` and for the `gh` PR API
  via `GH_TOKEN`. It has no cluster access at all.
- **Pull side (Argo CD):** the cluster pulls `env/huy` — Argo CD inside the
  cluster fetches the repo and reconciles. There is no CI kubeconfig, no
  long-lived cluster credential handed to CI, and no workflow step talks to the
  API server. The registry is read by the kubelet on the cluster, not by CI.

Keeping the two apart means a compromised or over-broad CI token still cannot
touch the cluster, and the cluster never needs a GitHub write credential.

## CI publication alone is not CD

A successful `publish` proves the image was built, pushed, signed and
smoke-tested. It does **not** change `env/huy` and does **not** change any pod.
The registry now holds a new digest; desired state is unchanged; the running app
is unchanged. Deployment happens only after (1) a reviewed merge updates
`env/huy` and (2) Argo CD syncs. If nobody merges, nothing deploys — that gap is
the approval boundary, not a failure. Do not describe a green pipeline as
"deployed"; describe it as "published and proposed."

## Manual vs automatic sync

- **Manual (default for `devops-app-huy`):** merging the digest PR moves
  `env/huy`; the Application goes `OutOfSync`, and a human clicks **Sync**. This
  is the mode that shows the approval boundary most clearly.
- **Automatic (enabled for the bonus loop):** `syncPolicy.automated` with
  `selfHeal: true`, `prune: true`, `allowEmpty: false`. After the reviewed merge,
  the controller syncs on its own; there is still a reviewed merge first.

## Drift and rollback

- **Drift:** scale the Deployment by hand (`kubectl -n devops-app scale
  deployment/devops-app --replicas=3`) while Git says `2`; with `selfHeal` the
  controller restores `2`. Do not "fix" drift with `kubectl apply` — Git is the
  desired state.
- **Rollback:** revert the image change **in Git**. Find the digest change commit
  on `env/huy` and revert only that commit through a reviewed env PR:

  ```sh
  git checkout -b revert/huy-<digest> origin/env/huy
  git revert <image-change-commit>          # reverts the image line only
  git push origin revert/huy-<digest>       # then open + review the PR
  ```

  Preserve later ConfigMap/security commits; revert only the intended image
  commit. Do not `kubectl rollout undo` against an Argo-managed Deployment — the
  controller would just restore Git's value.

## If a proposal branch diverges

`propose-deploy` will add a fast-forward commit to an existing
`deploy/huy-<sha>` branch only if that branch still descends from `env/huy`. If
`env/huy` moved on (or the branch otherwise diverged), the job **stops** with a
non-zero exit and prints both tips for human review: it never resets, rebases,
force-pushes, deletes or silently resolves the branch. Investigate the two
commits and resolve by hand.

## Quick reference

```sh
REPO=huyle-gradion/devops-observability-research-group-b
# what is proposed?
gh pr list --repo "$REPO" --base env/huy --state open
gh pr diff --repo "$REPO" <number>
# what is deployed / desired?
kubectl -n argocd    get application devops-app-huy -o yaml
kubectl -n devops-app get deployment devops-app \
  -o jsonpath='{.spec.template.spec.containers[?(@.name=="app")].image}{"\n"}'
kubectl -n devops-app get pods -l app=devops-app -o wide
```
