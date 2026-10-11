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
| Image digest (index) | pending | pending | 03–04 | pending |
| Base image digest | pending | pending | 03 | pending |
| kube-prometheus-stack chart version | pending | pending | 09 | pending |
| OpenTelemetry Collector version | pending | pending | 10 | pending |
| Go module pins | pending | pending | 02 | pending |
| GitHub Actions pinned revisions | pending | pending | 04 | pending |
| OTel JS Web SDK package version | pending | pending | 11 | pending |
| k3s version | pending | pending | 06 | pending |
| Argo CD or Flux version | pending | pending | 08 | pending |
