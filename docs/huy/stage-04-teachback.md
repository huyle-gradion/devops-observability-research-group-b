# Stage 04 teachback — CI publication and the CI/CD mental model (huyle-gradion)

Post-implementation teaching. Written 2026-10-11. Evidence: `captures/huy/step-1/ci/ci.md`.

## 1. What Stage 04 delivered

`.github/workflows/ci.yml` (`name: ci`) with three jobs:

| Job | Runs on | Purpose |
|---|---|---|
| `check` | every push/PR | anonymous checkout, `go vet`+`go test`, native `docker build` — **no secret** |
| `publish` | push to `personal/huy` only | multi-platform build + push to `ghcr.io/huyle-gradion/devops-app-huy`, full-SHA tag, SBOM, provenance |
| `smoke` | after `publish` | pull the exact digest, run it, assert `/health` and `/test` |

Real runs: PR run (check **pass**, publish/smoke **skipped**) and push run (check/publish/smoke
all **success**), publishing digest `sha256:737f62c8…`.

## 2. What went wrong / friction

- **Worker had no network** to `api.github.com`, so it could not re-verify the action pins.
  Sisyphus verified all five against GitHub release commits. Lesson: a sandboxed worker may
  not be able to prove supply-chain facts — verify pins where you *do* have network.
- **setup-go cache miss:** it looked for `go.mod` at the repo root, but the module is in
  `app/`. Non-fatal warning; the fix is `cache-dependency-path: app/go.sum` (or disable).
- **`actions/checkout` is not used:** we fetch the change with plain anonymous Git so no
  default GitHub token is ever injected into untrusted code.
- **Runner label drift:** `ubuntu-latest` migrates to Ubuntu 26 on 2026-10-19 — a reminder
  that "latest" is a moving target; note it in the write-up.

## 3. Lessons / assignment-required knowledge (research 3.1a)

- **What CI/CD is and what each stage gates (§3.1a).** *Continuous Integration* = every
  change is built and tested automatically so broken code is caught before merge;
  *Continuous Delivery/Deployment* = the tested artifact is shipped. The classic stages:
  **build** (does it compile?), **test** (does it behave?), **lint/scan** (is it clean/safe?),
  **deploy** (ship it). Each stage is a **gate** — a failing gate stops the pipeline.
- **GitHub Actions anatomy.** Trigger (`on:` push/pull_request) → **jobs** (parallel unless
  `needs:`) → **steps** (shell `run:` or `uses:` actions) on a **runner** (`ubuntu-latest`).
  `needs:` builds the DAG; `permissions:` scopes the built-in `GITHUB_TOKEN`.
- **Why a fork PR cannot publish.** Fork PRs run with a read-only token and **no secrets**.
  Our `publish`/`smoke` jobs are additionally gated to a push on `personal/huy`, so a PR —
  fork or same-repo — never receives `PROJECT_GITHUB_TOKEN`. That is the whole point of the
  event boundary: untrusted code must never see the credential.
- **Least privilege.** Top-level `permissions: {}` denies the built-in token everywhere;
  `publish` grants only `contents: read` + `packages: write`. We pass the project PAT
  **explicitly** (`github-token:` on build-push, password on login) instead of relying on
  any implicit token.
- **2025–2026 supply-chain additions (§3.1a).** *Immutable Actions* publish actions as OCI
  packages to `ghcr.io` so a pinned version cannot be changed later (rollout moving — check
  the GitHub changelog). *Artifact Attestations* are build provenance in SLSA/in-toto format
  signed with Sigstore, verifiable after the build. Our `build-push` already emits an SBOM
  and `provenance: mode=max` — the raw material Stage 05 turns into signed, verified proof.
- **Pin by digest, not tag.** Tags (`:latest`, `:v7`) are mutable; we tag with the **full
  source SHA** and identify the artifact by **index digest** `sha256:…`.
- **Publication ≠ runtime.** A successful push only proves the image *uploaded*. The separate
  `smoke` job proves the pushed digest actually starts and answers — don't conflate them.

## 4. How to run / read it

```sh
# what the workflow does, locally
cd app && go vet ./... && go test ./...
# inspect the real runs
python3 scripts/project-github.py gh run list --repo huyle-gradion/devops-observability-research-group-b
python3 scripts/project-github.py gh run view <run-id> --repo huyle-gradion/devops-observability-research-group-b --log
# anonymous package proof (after the package is made public)
python3 scripts/project-github.py gh pr checks 1 --repo huyle-gradion/devops-observability-research-group-b
```

**What you can do with it:** every push to `personal/huy` now produces a tagged,
multi-platform image automatically; open a PR and watch `publish` skip. The digest is the
artifact identity later stages deploy (Stage 08 proposes it by digest).

## 5. Fastest rebuild runbook

1. Resolve and **pin each action to a verified release commit** (not a major tag).
2. `name: ci`; `on:` push to the integration branch and `pull_request` targeting it, both
   path-filtered to the app/build inputs; `permissions: {}` at the top.
3. `check`: anonymous `git init`/`fetch` checkout; `setup-go` with `token: ''`; `go vet`+`go test`;
   native `docker build` (no login). No secret.
4. `publish`: `needs: check`; `if` push to the integration branch; `packages: write`; buildx +
   login (username = account, password = secret) + build-push with `platforms`, `push: true`,
   full-SHA tag, `sbom: true`, `provenance: mode=max`, explicit `github-token`; output `digest`.
5. `smoke`: `needs: publish`; pull `image@${{ needs.publish.outputs.digest }}`, run on
   loopback, poll `/health`, assert `/test`, always clean up the container.

**Traps:** moving tags instead of pinned commits; giving PR jobs the secret; relying on the
implicit `github.token` instead of the explicit PAT; tagging `:latest`; claiming publication
proves runtime; assuming a public repo means a public package.
