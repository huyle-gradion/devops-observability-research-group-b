# Stage 03 — container build/runtime evidence (huy)

Produced by Sisyphus from the owner-run privileged commands. Observed **2026-10-11**.
No secrets. Base SHA `b9a1c6dd0719630100f97c55bf812bc6a7f70b5d`. Owner-run because the
local Docker daemon requires `sudo`.

## Pinned bases

| Base | Index digest | Source |
|---|---|---|
| builder `golang:1.27` | `sha256:e432b43af23a9328d56a7c499be0476810aa344acbcf65fc7c455d4ff5a40602` | https://hub.docker.com/_/golang |
| runtime `gcr.io/distroless/static-debian13:nonroot` | `sha256:e2e927ec666bae08560abb3c55d0659eceabb657f56b6782ab500a9fc7f555e3` | https://github.com/GoogleContainerTools/distroless |

Both are multi-platform OCI indexes including `linux/amd64` and `linux/arm64`
(confirmed live by anonymous manifest requests).

## Build

`sudo docker buildx build --platform linux/amd64 --load -f app/Dockerfile -t devops-app-huy:container-proof .`

- First build: `[+] Building 32.4s (15/15) FINISHED`; builder resolved `golang:1.27@sha256:e432b43a…`,
  runtime resolved `…static-debian13:nonroot@sha256:e2e927e…`; compile step
  `RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app .`; exported
  image manifest list `sha256:d391019bc70b581955df122b8328e0a7f38ab9a9073a1788dfd7483edeecac84`.
- Rebuild (unchanged source): `[+] Building 1.8s (15/15) FINISHED` with every step
  `CACHED` (dockerfile config, both base `FROM`s, `WORKDIR`, both `COPY`s, the `go build`,
  and the runtime `COPY`) → layer cache reused.

## Run and inspect

- Container `690283fcc08395112a30c11749e76f69d31cd19f8d5946d7473639d9f0125a5d`;
  `-p 127.0.0.1:18080:8080 -e APP_ENV=container-proof`.
- `docker inspect … --format '{{.Config.Hostname}}'` → `690283fcc083`.
- `docker image inspect …` → `linux/amd64 user=65532:65532 entrypoint=["/app"]`.

| Request | Result |
|---|---|
| `GET /health` | `200 OK`, `text/plain`, `ok` (len 3) |
| `GET /test` | `200 OK`, `text/plain`, `X-App-Env: container-proof`, `Hello world from container-690283fcc083` (len 40) |

The `/test` identity matches the running container hostname (`690283fcc083`).

## No-shell proof (expected failure, recorded, not folded into a passing chain)

`sudo docker run --rm --entrypoint /bin/sh devops-app-huy:container-proof` →
`exec: "/bin/sh": stat /bin/sh: no such file or directory`, **exit 127**. The distroless
static runtime carries no shell.

## Cleanup

Only this stage's proof container/image were stopped and removed. No shared images or
volumes pruned.

## Boundaries

- Stage 03 exercised the **native amd64** build/run only. CI multi-platform publication
  and arm64 native execution are Stage 04; an arm64 manifest is not an arm64 native run.
- No in-image `HEALTHCHECK` (distroless); health is checked externally against `/health`.
