# Stage 04 — CI publication evidence (huy)

Produced by Sisyphus from real GitHub Actions runs. Observed **2026-10-11**. No secrets.

Repo: `huyle-gradion/devops-observability-research-group-b`.
Workflow: `.github/workflows/ci.yml` (`name: ci`).
Producing source commit (merge to `personal/huy`): `9e85849f78b4d0140b35a52b9005dc5d6d5fdd08`.

## Run 1 — pull request (must not publish)

- Run id `38119260804`, event `pull_request`, head `stage/huy/04-ci` @ `6077fea9`.
- URL: https://github.com/huyle-gradion/devops-observability-research-group-b/actions/runs/38119260804
- Jobs: `check` **success** (1m6s) — anonymous checkout, `go vet`/`go test`, native `docker build`;
  `publish` **skipped**; `smoke` **skipped**.
- PR checks view: `check pass`, `publish skipping`, `smoke skipping`.

Exit gate 04A satisfied: the PR check succeeds and publication is skipped.

## Run 2 — source push to personal/huy (publish + smoke)

- Run id `38119356408`, event `push`, head branch `personal/huy`, head SHA `9e85849f78b4d0140b35a52b9005dc5d6d5fdd08`.
- URL: https://github.com/huyle-gradion/devops-observability-research-group-b/actions/runs/38119356408
- Jobs: `check` success, `publish` success, `smoke` success.

Published artifact:
- Image: `ghcr.io/huyle-gradion/devops-app-huy`
- Tag: `sha=9e85849f78b4d0140b35a52b9005dc5d6d5fdd08`
- **Index digest:** `sha256:737f62c891b4f7fcebc3bf5c77e004ece31f12c4ee7854efe3eae21dc02fe430`
- Manifest-list export also recorded `sha256:737f62c8…`; SBOM + `provenance: mode=max` enabled.

Smoke (native amd64 runner, exact digest):

```
smoke ok: hostname ccee2028fff6, digest sha256:737f62c891b4f7fcebc3bf5c77e004ece31f12c4ee7854efe3eae21dc02fe430
```

The smoke job asserted `/health` = `200 ok`, `/test` = `200` body
`Hello world from container-ccee2028fff6` with header `X-App-Env: ci-proof`.

## Package visibility

`GHCR` package `devops-app-huy` is currently **private**. The owner must set it public so
the anonymous digest pull (04C) can be recorded. Public repo does not imply public package.

## Action pins (verified against GitHub releases)

| Action | Version | Commit |
|---|---|---|
| actions/setup-go | v7.0.0 | `b7ad1dad31e06c5925ef5d2fc7ad053ef454303e` |
| docker/setup-buildx-action | v4.4.1 | `f87e5991a6d7451dcb8d9637bfbc97413f497069` |
| docker/login-action | v4.6.0 | `dbcb813823bdd20940b903addbd779551569679f` |
| docker/build-push-action | v7.4.0 | `c3c9e263c25d99ce0380d002d59b67737d91b0dc` |

`actions/checkout` is intentionally **not used**: the change is fetched with anonymous Git
so no default GitHub token is injected.

## Run annotations (non-fatal)

- setup-go: `Restore cache failed: Dependencies file is not found …` — the module lives in
  `app/`, so caching looked at the repo root. Harmless; fixable with `cache-dependency-path`.
- Runner label notice: `ubuntu-latest` migrates to Ubuntu 26 from 2026-10-19.

## Boundaries

- Publication and smoke prove the **amd64** image runs natively on the runner; the arm64
  manifest is published but **not** a native arm64 execution test.
- Publication does not by itself prove smoke passed; the smoke job is the separate gate.
