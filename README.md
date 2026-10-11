# devops-observability-research-group-b

Personal DevOps and observability project for Assignment 3 Group B. It covers a small app, its container image, CI and supply-chain checks, a k3s cluster, GitOps deployment, metrics, tracing, browser telemetry, and load testing.

Current state: planning documents only. Implementation stages have not been run yet.

- [Personal plan index](.omo/plans/devops-personal/README.md)
- [Personal specification](.scratch/assignment-3-personal/spec.md)
- [Multi-agent workflow](docs/agents/opencode-workflow.md)
- [GitHub access policy](docs/agents/github-access.md)

## Runtime

The runtime API is a stdlib-only Go service in [`app/`](app/). The runtime identity
is resolved once from a nonempty `POD_NAME`, otherwise the hostname; `APP_ENV`
defaults to `local`. Run and exercise it locally:

```sh
cd app
APP_ENV=local go run .     # binds :8080
curl -i http://127.0.0.1:8080/health
curl -i http://127.0.0.1:8080/test
```

`GET /health` returns `ok`. `GET /test` returns `Hello world from container-<runtime-id>`
and echoes the resolved `APP_ENV` in `X-App-Env`. `GET /test?delay_ms=<0-1000>` adds
bounded latency and `GET /test?fail=1` returns an intentional 500; invalid parameters
are rejected with 400 before any delay.

## Container

[`app/Dockerfile`](app/Dockerfile) builds the image with the repository root as the
build context (a `$BUILDPLATFORM` Go builder cross-compiles a static binary for the
target platform; the runtime is the non-root, shell-less
`gcr.io/distroless/static-debian13:nonroot` image, running as UID/GID `65532` with
entrypoint `["/app"]` and listening on `:8080`).

Build and run locally (owner-run: the local Docker daemon needs `sudo`):

```sh
sudo docker buildx build --platform linux/amd64 --load \
  -f app/Dockerfile -t devops-app-huy:container-proof .
sudo docker run --rm -p 127.0.0.1:18080:8080 -e APP_ENV=container-proof \
  devops-app-huy:container-proof
# from another shell:
curl -i http://127.0.0.1:18080/health
curl -i http://127.0.0.1:18080/test
```

The runtime image contains no shell or package manager, so the container has no
in-image `HEALTHCHECK`; health is checked externally against `/health` (for example
with the HTTP probes added in the cluster stage).

## CI

[`.github/workflows/ci.yml`](.github/workflows/ci.yml) runs on pushes to
`personal/huy` and on pull requests targeting it, both filtered to `app/**`,
`web/**`, the workflow file and `.dockerignore`.

- `check` receives no project secret. It checks out the change anonymously
  (`git init` + public fetch of the push SHA or the PR merge ref), runs
  `go vet ./...` and `go test ./...` in `app/`, then a native
  `docker build -f app/Dockerfile -t ci-check .`. It never logs in or pushes.
- `publish` (push to `personal/huy` only) builds and pushes the multi-platform
  (`linux/amd64`, `linux/arm64`) image `ghcr.io/huyle-gradion/devops-app-huy` with
  a full-SHA tag, SBOM and maximum provenance, and exposes the published digest.
- `smoke` (`needs: publish`) pulls that exact digest, runs it on loopback with
  `APP_ENV=ci-proof`, and asserts `/health` is `200 ok` and `/test` returns
  `200 Hello world from container-<hostname>` with header `X-App-Env: ci-proof`.

Pull requests never publish: `publish` and `smoke` are gated on a push to
`personal/huy`, and the workflow's default `permissions` are empty. Only the
trusted `publish`/`smoke` steps receive the project token. Additionally, `publish`
now attaches signed SLSA build provenance for the pushed digest via
`actions/attest-build-provenance`, signed with the workflow's short-lived OIDC
identity (`id-token: write`), which is separate from the project PAT. Action pins
are recorded in [`docs/huy/versions.md`](docs/huy/versions.md).
