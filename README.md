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
