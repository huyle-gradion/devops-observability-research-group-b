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
| kube-prometheus-stack chart version | `92.3.0` (appVersion `v0.94.1`) | https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack | 09 (09B) | 2026-10-12 |
| OpenTelemetry Collector version | `0.162.0` (contrib distribution; core module `v0.162.0`) | https://github.com/open-telemetry/opentelemetry-collector-releases/releases/tag/v0.162.0 | 10 (10A) | 2026-10-12 |
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

## Monitoring stack pins — Stage 09B (2026-10-12)

Substep 09B pins the `kube-prometheus-stack` chart and resolves the component
image digests for the `huy-monitoring` release. Observed and resolved on
**2026-10-12** from the workstation.

### Helm and chart pin

| Item | Value | How observed |
|---|---|---|
| `helm` version | `v3.22.0` (`GitCommit 144ca65f8501953fa8b41cd1d37c7223051c85b7`, `GoVersion go1.26.8`) | `helm version` |
| Chart repository | `prometheus-community` = `https://prometheus-community.github.io/helm-charts` | `helm repo add prometheus-community https://prometheus-community.github.io/helm-charts` |
| Pinned chart version | `92.3.0` | `helm search repo prometheus-community/kube-prometheus-stack -l` (latest stable) |
| Chart `appVersion` (kube-prometheus / Prometheus Operator) | `v0.94.1` | `helm show chart prometheus-community/kube-prometheus-stack --version 92.3.0` |
| Values source | `helm show values prometheus-community/kube-prometheus-stack --version 92.3.0` (6219 lines, saved and inspected; not committed) | `helm show values ... --version 92.3.0` |

The chart version is a moving target; it is pinned here and passed explicitly as
`--version 92.3.0` on every Helm command, so a later `helm repo update` cannot
silently change what gets installed.

### Component image digests (multi-arch manifest-list digests)

Resolved anonymously on 2026-10-12 with
`docker buildx imagetools inspect <image>:<tag>`. Each value is the digest of the
**manifest list** the chart's image tag resolves to (the index that serves both
`linux/amd64` and `linux/arm64`), not a per-platform digest. These are what the
pinned chart actually pulls; the chart values reference the tags above, so the
digests are recorded for supply-chain reference, not substituted into objects.

| Component | Image (chart tag) | Manifest-list digest |
|---|---|---|
| Prometheus | `quay.io/prometheus/prometheus:v3.15.0-distroless` | `sha256:b2a413d5a03ea6a76782a508d1c7947440bba3b973931a25676e278431891b01` |
| Prometheus Config Reloader | `quay.io/prometheus-operator/prometheus-config-reloader:v0.94.1` | `sha256:06b52bd4dbe3ed6dd5905aadaf8b9987d9da4c31d37bb86deaadae68cee2d26b` |
| Prometheus Operator | `quay.io/prometheus-operator/prometheus-operator:v0.94.1` | `sha256:7c88d4e7bae63bd8d0f8da986054337c23b4472fa8d3c6817f5bc4b8407a2d6a` |
| Grafana | `docker.io/grafana/grafana:13.2.3-distroless` | `sha256:202e5d5b3f8499e70ccac7991eea6d7da72d2790eb2e0597e7987acf724e7180` |
| Node Exporter | `quay.io/prometheus/node-exporter:v1.12.1-distroless` | `sha256:8c9bac11973b94b59be88d6e11fee4429aa743c8846cdc75d65b18db33f6a106` |
| kube-state-metrics | `registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.20.0` | `sha256:42cfe3723a5f058171c627537fb57a3ea0f26e4380fa18555a95cb1a1b4cfc5b` |

Alertmanager is disabled in `k8s/monitoring/values.yaml`, so no Alertmanager image
digest is recorded.

Local, cluster-free validation of the values (`helm template`, no install) confirms
the rendered Prometheus CR resources/retention/storage, the Grafana resources and
1Gi PVC, and that only the apiserver/coredns/kubelet/operator/prometheus/grafana/
kube-state-metrics/node-exporter ServiceMonitors are generated (scheduler,
controller-manager, etcd and kube-proxy are absent as intended).

**Unrun:** the actual `helm upgrade --install huy-monitoring ...` and the
`kubectl apply -f k8s/monitoring/servicemonitor.yaml` were **not** run here — this
worker has no kubeconfig and no cluster access. Sisyphus/owner runs the real install;
these digests are what that install will pull.

## Tracing stack pins — Stage 10A (2026-10-12)

Substep 10A pins the OpenTelemetry Collector and the Jaeger viewer by **tag + digest**.
Both digests were resolved and confirmed **anonymously** on 2026-10-12 against their
registries (`docker-content-digest` from a manifest `HEAD`, no credential), and they
match the references supplied for the stage.

| Component | Image | Tag | Digest (multi-arch index) | Registry source |
|---|---|---|---|---|
| Collector (contrib) | `ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-contrib` | `0.162.0` | `sha256:39923a8e431bd1f57be82411999d389fcfe40857492e4365456d97a4c1f74be6` | https://github.com/open-telemetry/opentelemetry-collector-releases/pkgs/container/opentelemetry-collector-contrib |
| Jaeger (all-in-one) | `jaegertracing/jaeger` | `2.22.0` | `sha256:836b967b5abcf1f7758066aa6b348619267a2f1459f72d472840570f9bd88a17` | https://hub.docker.com/r/jaegertracing/jaeger |

The same Jaeger image is mirrored at `cr.jaegertracing.io/jaegertracing/jaeger:2.22.0`
(the reference the Jaeger *getting started* doc uses); `k8s/tracing/jaeger.yaml` pins the
Docker Hub name plus digest.

Jaeger runs its **embedded all-in-one** configuration
(`cmd/jaeger/internal/all-in-one.yaml`, compiled into the binary): OTLP/HTTP on
`0.0.0.0:4318`, the UI on `16686`, and **transient in-memory** trace storage
(`jaeger_storage` → `memory`, `max_traces: 100000`). No PVC and no durable-storage claim.

### Collector component / module versions

`otelcol-contrib --version` reports `0.162.0`. In this version the OTLP/HTTP **exporter**
type is `otlp_http` (canonical) with `otlphttp` as its **deprecated alias** — the exporter
module is `go.opentelemetry.io/collector/exporter/otlphttpexporter v0.162.0` (reported by
`otelcol-contrib components`). The obsolete Jaeger-native `jaeger` **exporter** does
**not** exist in 0.162.0; the manifest uses `otlp` (receiver) → `otlp_http/jaeger`
(exporter), never the removed exporter or its old gRPC-era endpoint.

### Collector `validate` evidence (this worker, 2026-10-12)

No local docker daemon was reachable (`sudo` needs a password), so the config was
validated against the pinned **release binary** instead of the container image — same
version, checksum-verified against the published `.sha256`:

```text
otelcol-contrib_0.162.0_linux_amd64.tar.gz
  published sha256 fcc063749f730f8c21fe29f2d340ff174f5f1c5885bd3156fb6c985a3036fcc3  (matched)
otelcol-contrib version 0.162.0

# config extracted from the ConfigMap in k8s/tracing/collector.yaml
./otelcol-contrib validate --config /tmp/opencode/otel-collector-config.yaml   # exit 0
```

`validate` accepted the extracted ConfigMap config (exit 0, no output). Negative control:
a config using the obsolete `jaeger` exporter fails with
`'exporters' unknown type: "jaeger" ... (valid values: [... otlp_http ... otlphttp ...])`,
so the check discriminates rather than rubber-stamping.

**Unrun:** the image-form `docker run ... validate` (runbook §10.1) and every cluster step
(`kubectl apply`, rollout status, port-forward) were **not** run here — no docker daemon
and no cluster access. Sisyphus/owner runs them and captures the cluster evidence.

## Research claim registry — versions and status (Stage 13A, 2026-10-13)

This table records the observed **version/status** that the Stage 13A findings claims rest
on ([`findings.md`](findings.md)). Every value was read from a committed project artifact
(this repo) on 2026-10-13, or recorded by the producing stage named below; the Source URL
is the canonical release/registry page for that version, not a re-fetched landing page.
The primary **documentation pages** for each claim (with section anchors and checked dates)
live in `findings.md`, not here.

| Item | Observed value / status | Where observed | Source URL | Checked |
|---|---|---|---|---|
| OTel Collector (contrib distribution) | `0.162.0` (core module `v0.162.0`); `otlp_http` canonical, `otlphttp` deprecated alias | `k8s/tracing/collector.yaml`; Stage 10A | https://github.com/open-telemetry/opentelemetry-collector-releases/releases/tag/v0.162.0 | 2026-10-12 |
| OpenTelemetry Go SDK | `go.opentelemetry.io/otel` `v1.47.0`; `otlptracehttp` `v1.47.0`; `otelhttp` `v0.72.0`; semconv `v1.43.0` | `app/go.mod` (read 2026-10-13) | https://pkg.go.dev/go.opentelemetry.io/otel@v1.47.0 | 2026-10-13 |
| OpenTelemetry JS / browser packages | `@opentelemetry/sdk-trace-web` `2.12.0`; `api` `1.9.1`; `context-zone` `2.12.0`; `resources` `2.12.0`; `sdk-trace-base` `2.12.0`; `exporter-trace-otlp-http` `0.223.0`; `instrumentation` `0.223.0`; `instrumentation-fetch` `0.223.0`; `instrumentation-xml-http-request` `0.223.0`; `instrumentation-document-load` `0.68.0` | `web/package.json` (read 2026-10-13); **EXPERIMENTAL** browser instrumentation | https://www.npmjs.com/package/@opentelemetry/sdk-trace-web | 2026-10-13 |
| `web-vitals` JS library | `6.2.3` | `web/package.json` (read 2026-10-13) | https://www.npmjs.com/package/web-vitals | 2026-10-13 |
| `esbuild` (browser bundle) | `0.28.2` | `web/package.json` (read 2026-10-13) | https://www.npmjs.com/package/esbuild | 2026-10-13 |
| Prometheus | chart tag `quay.io/prometheus/prometheus:v3.15.0-distroless` | `docs/huy/versions.md` Stage 09B | https://quay.io/repository/prometheus/prometheus?tab=tags | 2026-10-12 |
| Grafana | chart tag `docker.io/grafana/grafana:13.2.3-distroless` | `docs/huy/versions.md` Stage 09B | https://hub.docker.com/r/grafana/grafana | 2026-10-12 |
| `kube-prometheus-stack` chart | `92.3.0` (appVersion `v0.94.1`) | `docs/huy/versions.md` Stage 09B | https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack | 2026-10-12 |
| k3s | `v1.37.1+k3s1` (single-node server) | `docs/huy/cluster-runbook.md`; `captures/huy/step-2/cluster.md` | https://github.com/k3s-io/k3s/releases/tag/v1.37.1%2Bk3s1 | 2026-10-11 |
| kubectl | `v1.37.1` | `docs/huy/versions.md` Stage 08B | https://dl.k8s.io/release/v1.37.1/bin/linux/amd64/kubectl.sha256 | 2026-10-12 |
| k6 | `v2.3.0` | `docs/huy/load-runbook.md`; `captures/huy/step-5/run-metadata.md` | https://github.com/grafana/k6/releases/tag/v2.3.0 | 2026-10-12 |
| Argo CD | `v3.5.4` (pinned `manifests/install.yaml`) | `docs/huy/versions.md` Stage 08A; `captures/huy/bonus-1/adoption.md` | https://github.com/argoproj/argo-cd/releases/tag/v3.5.4 | 2026-10-12 |
| Jaeger | `2.22.0` all-in-one (transient in-memory storage) | `docs/huy/versions.md` Stage 10A; `k8s/tracing/jaeger.yaml` | https://hub.docker.com/r/jaegertracing/jaeger | 2026-10-12 |
| Trivy | `0.75.0` | `docs/huy/versions.md` Stage 04b; `captures/huy/bonus-4/trivy-image.json` | https://github.com/aquasecurity/trivy/releases/tag/v0.75.0 | 2026-10-12 |
| Go toolchain | `go1.27.2 linux/amd64` | `docs/huy/versions.md` Stage 01 | https://go.dev/dl/ | 2026-10-11 |
| Docker client + engine | `29.7.2` (API 1.55; containerd `v2.3.5`, runc `1.5.1`) | `docs/huy/versions.md` Stage 01 | https://docs.docker.com/engine/release-notes/ | 2026-10-11 |
| Node.js / npm | `v26.10.0` / `11.19.1` | `docs/huy/versions.md` Stage 01 | https://nodejs.org/en/blog | 2026-10-11 |
| cloudflared | `2026.8.3` | `docs/huy/versions.md` Stage 01 | https://github.com/cloudflare/cloudflared/releases | 2026-10-11 |
| Tailscale | `1.102.3` | `docs/huy/versions.md` Stage 01 | https://tailscale.com/changelog | 2026-10-11 |

### Moving-claim status carried by the registry

Two findings rows are moving claims, not settled facts. Their current status is recorded in
`findings.md` with explicit **recheck** entries dated **16 Oct 2026** and **30 Oct 2026**
(both `pending` until performed):

- **Immutable Actions rollout / OCI-on-GHCR claim** — see `findings.md` §a. The GitHub
  changelog is the primary source; "Immutable Releases GA" does not by itself establish
  Immutable Actions as OCI packages on GHCR.
- **OpenTelemetry per-signal maturity** — see `findings.md` §e. The authoritative source is
  <https://opentelemetry.io/status/>, which changes over time.

No row above asserts a vendor "stable since" date or an endorsement; each value is the
version actually observed in this project's artifacts.
