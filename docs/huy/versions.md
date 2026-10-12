# Tool versions, architecture and pins (personal project — huyle-gradion)

Recorded by worker during Stage 01 (substep 01C). Observed on the workstation on
**2026-10-11**. Nothing in the "To be appended later" section is known yet; do not
fill it with guesses.

## Observed workstation (2026-10-11)

| Item | Observed value | How observed |
|---|---|---|
| Date of observation | 2026-10-11 | `date` |
| Kernel / platform | Linux | `uname -s` |
| CPU architecture | `x86_64` (native amd64) | `uname -m` |
| git | 2.55.0 | `git --version` |
| go | go1.27.2 linux/amd64 | `go version` |
| docker (client) | 29.7.2 (build a7dcaa6fdb), API 1.55, Go go1.26.5, built 2026-08-06, linux/amd64 | `docker --version`; `sudo docker version` |
| docker engine (server) | **active**; Engine 29.7.2, API 1.55 (minimum 1.40), Go go1.26.5, linux/amd64 | `sudo docker version` (owner-run 2026-10-11) |
| containerd | v2.3.5 | `sudo docker version` |
| runc | 1.5.1 | `sudo docker version` |
| gh | 2.102.0 (2026-09-30) | `gh --version` |
| curl | 8.22.0 (x86_64-pc-linux-gnu) | `curl --version` |
| node | v26.10.0 | `node --version` |
| npm | 11.19.1 (present) | `npm --version` |
| cloudflared | 2026.8.3 (built 20260831-21:22:47) | `cloudflared --version` |
| tailscale | 1.102.3 | `tailscale version` |
| python3 | 3.14.7 | `python3 --version` |

### Not yet installed at observation time

| Tool | State | Needed by |
|---|---|---|
| kubectl | not found | Stages 06–08 |
| helm | not found | Stages 06, 09 (kube-prometheus-stack) |
| k6 | not found | Stages 12, 14 |
| limactl | not found | Stage 06 (if selected) |

Install later stages' tools only when that stage requires them and permission exists,
using official package/release checksums and stable versions (Stage 01 plan).

## Environment facts not yet known

| Item | Value |
|---|---|
| Presentation timezone | pending (owner) |
| Wednesday 14 Oct review time | pending (owner) — do not invent a time |

## To be appended later (pins, digests, chart values, source URLs)

This section is intentionally empty of values. It will be filled as later stages
produce real pins, image digests, chart versions and source URLs. Each entry must
name the producing stage, the observed value and the primary source URL that was
actually checked.

| Item | Value | Source URL | Producing stage | Checked date |
|---|---|---|---|---|
| Image digest (index) | `sha256:737f62c891b4f7fcebc3bf5c77e004ece31f12c4ee7854efe3eae21dc02fe430` | https://github.com/huyle-gradion/devops-observability-research-group-b/pkgs/container/devops-app-huy | 04 | 2026-10-11 |
| Builder base index digest (`golang:1.27`) | `sha256:e432b43af23a9328d56a7c499be0476810aa344acbcf65fc7c455d4ff5a40602` | https://hub.docker.com/_/golang | 03 | 2026-10-11 |
| Runtime base index digest (`gcr.io/distroless/static-debian13:nonroot`) | `sha256:e2e927ec666bae08560abb3c55d0659eceabb657f56b6782ab500a9fc7f555e3` | https://github.com/GoogleContainerTools/distroless | 03 | 2026-10-11 |
| kube-prometheus-stack chart version | pending | pending | 09 | pending |
| OpenTelemetry Collector version | pending | pending | 10 | pending |
| Go module pins | pending | pending | 02 | pending |
| GitHub Actions pinned revisions | see "GitHub Actions pins" below | per-action release pages | 04 | 2026-10-11 |
| OTel JS Web SDK package version | pending | pending | 11 | pending |
| k3s version | pending | pending | 06 | pending |
| Argo CD or Flux version | Argo CD v3.5.4 (pinned `manifests/install.yaml`) | https://github.com/argoproj/argo-cd/releases/tag/v3.5.4 | 08 (08A) | 2026-10-12 |

Both Stage 03 bases above are multi-platform OCI image indexes that include
`linux/amd64` and `linux/arm64`; the digests were confirmed live against the
registries (anonymous manifest requests) on 2026-10-11. Stage 03 exercises only the
native amd64 build and run. CI multi-platform publication and arm64 native execution
are verified later (Stage 04), not here.

## GitHub Actions pins (Stage 04 — 2026-10-11)

`.github/workflows/ci.yml` pins every action to a reviewed release commit, never a
moving major tag. Source URLs are the action release pages.

| Action | Version | Commit SHA | Source URL |
|---|---|---|---|
| `actions/checkout` | v7.0.1 | `3d3c42e5aac5ba805825da76410c181273ba90b1` | https://github.com/actions/checkout/releases/tag/v7.0.1 |
| `actions/setup-go` | v7.0.0 | `b7ad1dad31e06c5925ef5d2fc7ad053ef454303e` | https://github.com/actions/setup-go/releases/tag/v7.0.0 |
| `docker/setup-buildx-action` | v4.4.1 | `f87e5991a6d7451dcb8d9637bfbc97413f497069` | https://github.com/docker/setup-buildx-action/releases/tag/v4.4.1 |
| `docker/login-action` | v4.6.0 | `dbcb813823bdd20940b903addbd779551569679f` | https://github.com/docker/login-action/releases/tag/v4.6.0 |
| `docker/build-push-action` | v7.4.0 | `c3c9e263c25d99ce0380d002d59b67737d91b0dc` | https://github.com/docker/build-push-action/releases/tag/v7.4.0 |

These SHAs were supplied by the owner-controlled execution directive for Stage 04;
they were not independently re-resolved here (no outbound network to the GitHub API
from the worker sandbox). Sisyphus should confirm each release commit before push.

The Stage 04 workflow uses only four of these pins: `actions/setup-go`,
`docker/setup-buildx-action`, `docker/login-action` and `docker/build-push-action`.
`actions/checkout` is recorded here for reference but is not referenced by
`.github/workflows/ci.yml`, which checks out the public change with plain anonymous
Git so no default GitHub token is ever injected.

## GitHub Actions pins (Stage 05 — signed provenance addendum)

Stage 05 substep 05A adds one pinned action to the `publish` job. The stage plan
named the generic `actions/attest`, but that action is a generic attestation
primitive that requires the caller to supply a predicate; the action actually
executed is **`actions/attest-build-provenance`**, which generates SLSA build
provenance signed with the workflow's short-lived OIDC identity (`id-token: write`),
not the generic `actions/attest`. (As of v4 the wrapper is a thin layer over
`actions/attest`; v4.2.2 pins `actions/attest@508db95dd578ae2727ebd6217d5ba78e4fbda05d`
internally.)

| Action | Version | Commit SHA | Source URL |
|---|---|---|---|
| `actions/attest-build-provenance` | v4.2.2 | `4d101475d8b20a2381f78447822ac1eab6504dd8` | https://github.com/actions/attest-build-provenance/releases/tag/v4.2.2 |

The v4.2.2 commit SHA above was independently confirmed live on 2026-10-11 via an
anonymous public `GET /repos/actions/attest-build-provenance/commits/v4.2.2`
(returned `4d101475d8b20a2381f78447822ac1eab6504dd8`, a validly-signed release commit).

## GitHub Actions pins (Stage 04b — six-job refactor, 2026-10-12)

The six-job `.github/workflows/ci.yml` (`lint -> compile -> test -> build -> smoke -> publish`)
uses the same five pinned actions as Stages 04/05; no new action was added. The Stage 04 note
above asked Sisyphus to confirm each release commit before push. That confirmation was made
live on **2026-10-12** with anonymous `GET /repos/<owner>/<repo>/commits/<tag>` (no credential),
and every tag resolved to the pinned commit:

| Action | Version | Commit SHA | Anonymous verification |
|---|---|---|---|
| `actions/setup-go` | v7.0.0 | `b7ad1dad31e06c5925ef5d2fc7ad053ef454303e` | `commits/v7.0.0` → same SHA |
| `docker/setup-buildx-action` | v4.4.1 | `f87e5991a6d7451dcb8d9637bfbc97413f497069` | `commits/v4.4.1` → same SHA |
| `docker/login-action` | v4.6.0 | `dbcb813823bdd20940b903addbd779551569679f` | `commits/v4.6.0` → same SHA |
| `docker/build-push-action` | v7.4.0 | `c3c9e263c25d99ce0380d002d59b67737d91b0dc` | `commits/v7.4.0` → same SHA |
| `actions/attest-build-provenance` | v4.2.2 | `4d101475d8b20a2381f78447822ac1eab6504dd8` | `commits/v4.2.2` → same SHA |

`actions/checkout` is still not referenced: every job uses plain anonymous Git. The build job
adds `sbom: true` and `provenance: mode=max` on `docker/build-push-action` with an
`outputs: type=oci,dest=...` archive (no push, no login); the local `--load` path cannot carry
attestations, so the OCI exporter is used instead.

| Item | Value | Producing stage | Checked date |
|---|---|---|---|
| Trivy scanner (pinned download in `smoke`) | `0.75.0` (`trivy_0.75.0_Linux-64bit.tar.gz`, checksum-verified) | 04b | 2026-10-12 |
| Trivy vulnerability DB snapshot (local check) | `UpdatedAt 2026-10-12 01:14:57 UTC` | 04b | 2026-10-12 |

## GitOps digest proposal — Stage 08B (2026-10-12)

Substep 08B adds one job, `propose-deploy`, to `.github/workflows/ci.yml`. It runs
only after a fully green trusted push and proposes a reviewed digest PR against
`env/huy`. It adds **no new GitHub Action**: every step is a plain `run:` shell
step, so the five action pins already recorded for Stages 04/05/04b are unchanged
and are not repeated here.

### Pinned tool added by 08B

`propose-deploy` installs a pinned `kubectl` and uses it only for a local,
server-free manifest edit (`kubectl set image --local -f ... -o yaml`):

| Item | Value | Source URL | Producing stage | Checked date |
|---|---|---|---|---|
| `kubectl` version | `v1.37.1` | https://dl.k8s.io/release/v1.37.1/bin/linux/amd64/kubectl | 08B | 2026-10-12 |
| `kubectl` linux/amd64 sha256 | `65691ff77eb6fa44c908b77a1082c9f092c3b9733b5cefabec0d1104890e21a8` | https://dl.k8s.io/release/v1.37.1/bin/linux/amd64/kubectl.sha256 | 08B | 2026-10-12 |

The pin matches the cluster server minor: k3s `v1.37.1+k3s1`
([`cluster-runbook.md`](cluster-runbook.md)). The checksum was confirmed live on
2026-10-12 by fetching the official `.sha256` over the public HTTPS release path
(no credential) and by comparing the workstation client binary, which hashed to
the same value. The workflow verifies the download against `kubectl.sha256`
before installing.

Confirmation requested before push, as for the Stage 04/05/04b pins: the
`propose-deploy` job installs only this one new third-party binary and adds no
`uses:` line, so no action commit SHA needs re-resolution for 08B.
