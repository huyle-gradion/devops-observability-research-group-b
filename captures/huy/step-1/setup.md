# Stage 01 — workspace, prerequisites and setup evidence (huy)

Produced by Sisyphus during Stage 01. Observed on **2026-10-11**. No secrets,
kubeconfigs, authorization headers or credential output are recorded here.

## Repository and refs

| Item | Value |
|---|---|
| REPO | `huyle-gradion/devops-observability-research-group-b` |
| URL | https://github.com/huyle-gradion/devops-observability-research-group-b |
| Visibility | **public** (`isPrivate:false`), personal account (not the `gradionhq` org) |
| Default branch | `main` (seed `cfc71a2`) |
| Integration branch | `personal/huy` |
| Stage 01 editing branch | `stage/huy/01-workspace` |
| Stage 01 worktree | `/home/thearesia/Projects/devops-worktrees/01-workspace` |
| Stage 01 commit | `2759e1e89a2ecd656028edc3a5f040117e57709b` |
| DOC_ROOT | `/home/thearesia/Projects/assignment-3` (document source; not a Git checkout) |

Remote refs verified through the project credential helper: `main` @ `cfc71a2`,
`personal/huy` @ `2759e1e`, `stage/huy/01-workspace` @ `2759e1e`.

## Credential boundary (no secret values)

- Authenticated access uses `GITHUB_ACCESS_TOKEN` from the private `.env`
  (`PROJECT_GITHUB_ENV_FILE`), loaded in memory by `scripts/project-github.py`.
- Identity check (`gh api user --jq .login`) resolved to **huyle-gradion**; the
  helper fails closed on any other login and never falls back to a cached account.
- No token value was printed, captured or committed; the cached personal account is
  not used.

## Tool inventory and architecture (2026-10-11)

| Tool | Version / state |
|---|---|
| OS / arch | Linux, `x86_64` (native amd64) |
| git | 2.55.0 |
| go | go1.27.2 linux/amd64 |
| docker (client) | 29.7.2 (API 1.55, built 2026-08-06) |
| docker engine (server) | 29.7.2, API 1.55 (min 1.40) — active |
| containerd | v2.3.5 |
| runc | 1.5.1 |
| gh | 2.102.0 |
| curl | 8.22.0 |
| node | v26.10.0 |
| npm | 11.19.1 |
| cloudflared | 2026.8.3 |
| tailscale | 1.102.3 |
| python3 | 3.14.7 |
| kubectl / helm / k6 / limactl | not installed (needed by later stages) |

## Docker runtime proof (owner-run, privileged)

Docker was initially inactive. The owner ran the authorized fix
(`sudo systemctl start docker`); `systemctl is-active docker` then reported `active`.
Privileged access is owner-performed because `sudo -n` requires a password and the
user is not in the `docker` group — the authorized method is explicit `sudo docker`
(no socket chmod, no persistent group change).

```
$ sudo docker version
Client: 29.7.2, API 1.55, linux/amd64
Server: Engine 29.7.2, API 1.55 (minimum 1.40), linux/amd64
        containerd v2.3.5, runc 1.5.1, docker-init 0.19.0

$ sudo docker run --rm hello-world
Digest: sha256:5e23090353324d887c48ad5e5c56d294eab81588df9605b07d1afe895f9cc8f8
Hello from Docker!
```

Result: engine responds and its container run succeeds (`Hello from Docker!`, amd64).

## Checkpoints (owner local timezone pending)

| Checkpoint | Time |
|---|---|
| CP1 | Mon 12 Oct 09:00 |
| CP2 | Mon 12 Oct 17:00 |
| CP3 | Tue 13 Oct 15:00 |

Timezone and the Wednesday review time remain owner-provided; neither is invented here.

## Pending at capture time

- Owner learning gate (Docker privilege / public-vs-private / branch discipline).
- `workflow` scope on the same project credential (blocks Stage 04 workflow pushes).
- Remote worker-VM allocation (Stage 06 / bonus 2), timezone, team integration.
