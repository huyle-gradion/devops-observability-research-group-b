# Stage 04b — CI refactor to six jobs (huy)

Prepared by the implementation worker on `stage/huy/04b-ci-6jobs` (base `origin/personal/huy`
at `bba9f82dfb741d8635ba421108979b12c8ca8a63`). Observed **2026-10-12**. No secrets.

This capture records local preparation only. It is **not** proof that a real workflow run
passed; real PR and push runs are observed by Sisyphus after integration.

## Why this refactor

The Stage 04 workflow had three jobs (`check`, `publish`, `smoke`) with publication and the
digest smoke both after one combined check. The team's agreed pipeline has eight named steps;
this refactor keeps that step language while expressing it as six ordered GitHub Actions jobs
with an explicit `needs` chain, so a reviewer can see where each gate sits.

## Team eight-step language inside six jobs

| Team step | Job | What it gates |
|---|---|---|
| 1 lint | `lint` | `gofmt` + `go vet` clean |
| 2 compile | `compile` | the binary compiles (`CGO_ENABLED=0 go build`) |
| 3 test | `test` | `go test` locks exact `/health` and `/test` bytes and headers |
| 4 build | `build` | image builds with no login/no push; SBOM + provenance emitted; base digests recorded |
| 5 local smoke + scan | `smoke` | local image runs on loopback and answers the contract; Trivy report |
| 6 publish + attest + digest smoke | `publish` | multi-arch push, signed provenance, digest re-smoke |

`needs` chain: `lint -> compile -> test -> build -> smoke -> publish`. Any failing gate stops
the chain.

## Triggers and permissions

- `on.push` and `on.pull_request` both target `personal/huy`, path-filtered to `app/**`,
  `web/**`, `.github/workflows/ci.yml`, `.dockerignore`. No `pull_request_target`; no trigger
  on `env/huy`.
- Top-level `permissions: {}` denies the built-in token everywhere.
- `publish` alone adds `contents: read`, `packages: write`, `id-token: write`, and alone
  receives `secrets.PROJECT_GITHUB_TOKEN`. `lint`/`compile`/`test`/`build`/`smoke` receive no
  project secret.
- `publish` is guarded by `if: github.event_name == 'push' && github.ref == 'refs/heads/personal/huy'`.

## Per-job behaviour

- **lint** — anonymous Git checkout (`git init`/`fetch`, no `actions/checkout`, no persisted
  credentials), `actions/setup-go` with `token: ''`, `gofmt -l` fails on any unformatted file,
  then `go vet ./...`. The default setup-go cache path is fixed with `cache-dependency-path: app/go.mod`.
- **compile** — anonymous checkout + setup-go, `CGO_ENABLED=0 go build -o /tmp/ci-compile-app .`
  in `app/`. Proves the binary builds before the tests run.
- **test** — anonymous checkout + setup-go, `go test ./...` in `app/`. `app/main_test.go` asserts
  the exact `/health` body (`ok\n`), the exact `/test` body including the runtime id, the
  `X-App-Env` header, and the status codes.
- **build** — anonymous checkout + `docker/setup-buildx-action`, then `docker/build-push-action`
  with `context: .`, `file: app/Dockerfile`, `platforms: linux/amd64`, `push: false`,
  `load` unset, `outputs: type=oci,dest=/tmp/ci-image.tar`, `sbom: true`,
  `provenance: mode=max`, `github-token: ''`. A shell step prints every `FROM` line of
  `app/Dockerfile` and fails if any base image is not pinned by full `@sha256:` digest.
- **smoke** — anonymous checkout, `docker build -f app/Dockerfile -t devops-app-huy:ci-local .`
  (no login, no push), install pinned Trivy `0.75.0` with checksum verification, run the image
  on `127.0.0.1:18080` with `APP_ENV=ci-local`, poll `/health` for up to 30s, assert `200 ok`
  plus the exact `/test` body (hostname) plus `X-App-Env: ci-local`, then a report-only Trivy
  scan. A `trap` always stops and removes only its own container (`ci-local-smoke`).
- **publish** — anonymous checkout + buildx + `docker/login-action` (username `huyle-gradion`,
  password `secrets.PROJECT_GITHUB_TOKEN`), `docker/build-push-action` with
  `platforms: linux/amd64,linux/arm64`, `push: true`, tag `ghcr.io/huyle-gradion/devops-app-huy:${{ github.sha }}`,
  OCI `source`/`revision` labels, `sbom: true`, `provenance: mode=max`, explicit
  `github-token: ${{ secrets.PROJECT_GITHUB_TOKEN }}`; then `actions/attest-build-provenance`
  with `subject-name` without tag, `subject-digest` from the build output, `push-to-registry: true`,
  `create-storage-record: false`; then pulls `${IMAGE_NAME}@${digest}` and repeats the same
  health/test/header asserts. The digest is exposed as the job output.

## Design decision — build output uses the OCI exporter, not `--load`

`docker/build-push-action` exposes `sbom` and `provenance`, but BuildKit cannot attach
attestations to an image exported for the local Docker daemon (`load: true` / `type=docker`).
With an attestation requested and no OCI-capable image store, buildx fails with
`docker exporter does not currently support exporting manifest lists` (buildx validation in
`build/opt.go`; see also docker/buildx issues #1714, #1135). Routing the no-push build through
`type=oci` is the supported path that actually emits the SBOM and max provenance, so the build
job uses an OCI archive output and the smoke job builds its own runnable copy with plain
`docker build`. This keeps every job self-contained (each does its own anonymous checkout) and
keeps all attestation inputs declared. The published, signed attestations remain the ones the
`publish` job produces.

## Action pins (verified live 2026-10-12)

Each `uses:` is pinned to a full release commit with the release name in a comment. Pins were
re-resolved this session with anonymous `GET /repos/<owner>/<repo>/commits/<tag>` (no credential).

| Action | Version | Commit SHA |
|---|---|---|
| `actions/setup-go` | v7.0.0 | `b7ad1dad31e06c5925ef5d2fc7ad053ef454303e` |
| `docker/setup-buildx-action` | v4.4.1 | `f87e5991a6d7451dcb8d9637bfbc97413f497069` |
| `docker/login-action` | v4.6.0 | `dbcb813823bdd20940b903addbd779551569679f` |
| `docker/build-push-action` | v7.4.0 | `c3c9e263c25d99ce0380d002d59b67737d91b0dc` |
| `actions/attest-build-provenance` | v4.2.2 | `4d101475d8b20a2381f78447822ac1eab6504dd8` |

`actions/checkout` is intentionally not used: every job fetches the change with plain anonymous
Git so no default GitHub token is injected into untrusted code.

## Base image digests recorded by the build job

`app/Dockerfile` (unchanged) pins, both multi-platform OCI indexes:

```
FROM --platform=$BUILDPLATFORM golang:1.27@sha256:e432b43af23a9328d56a7c499be0476810aa344acbcf65fc7c455d4ff5a40602 AS build
FROM gcr.io/distroless/static-debian13:nonroot@sha256:e2e927ec666bae08560abb3c55d0659eceabb657f56b6782ab500a9fc7f555e3
```

## Local verification performed

| Command | Result |
|---|---|
| `python3` `yaml.safe_load` on `.github/workflows/ci.yml` | loads; triggers, `{}` permissions, job order and `needs` chain confirmed |
| `bash -n` on every `run:` block | all 16 blocks parse clean |
| `gofmt -l .` in `app/` | exit 0, no files listed |
| `go vet ./...` in `app/` | exit 0 |
| `go test ./...` in `app/` | exit 0 (`ok .../app 1.040s`) |
| `CGO_ENABLED=0 go build -o /tmp/ci-compile-app .` in `app/` | exit 0; static ELF x86-64 produced |
| `trivy image --exit-code 0 --no-progress --severity HIGH,CRITICAL <public digest>` | exit 0; report-only; DB `UpdatedAt 2026-10-12 01:14:57 UTC` recorded via `trivy --version` |

The Trivy command was exercised against the already-public digest
`ghcr.io/huyle-gradion/devops-app-huy@sha256:737f62c891b4f7fcebc3bf5c77e004ece31f12c4ee7854efe3eae21dc02fe430`
to validate the command shape and DB-date capture; the workflow scans its own local
`devops-app-huy:ci-local` tag instead.

## Unrun checks (not performed here)

- `docker build` / `docker buildx build --load` / `docker run` / container curl asserts: **blocked
  on this workstation**. The Docker daemon socket is `root:docker` mode 660, the user is not in the
  `docker` group, and `sudo` requires a password (`sudo -n true` exit 1). No socket/permission change
  was attempted. These must run on the GitHub runner and the owner's reproduction.
- Real PR run and source-push run: **unrun** — no push to `personal/huy` was permitted for this
  worker; observation belongs to Sisyphus after integration.
- arm64 native execution: **unrun** — the workflow publishes an arm64 manifest only; no arm64 host
  is available here.

## Tracked facts to carry forward

- Image identity is the index digest from `publish`; tags are the full `github.sha`, never `latest`.
- The project credential's missing `workflow` scope still blocks pushing a changed workflow file to
  GitHub from the worker; the branch `stage/huy/04b-ci-6jobs` and its docs remain local until
  Sisyphus publishes on an authorized path.
