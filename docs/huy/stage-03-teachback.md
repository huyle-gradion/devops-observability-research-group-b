# Stage 03 teachback — non-root multi-platform container (huyle-gradion)

Post-implementation teaching. Written 2026-10-11. Evidence: `captures/huy/step-1/container.md`.

## 1. What Stage 03 delivered

`app/Dockerfile` (multi-stage, non-root, shell-less), root `.dockerignore`, and a
`Container` section in `README.md`, with both base images pinned by **index digest**.
Image: `linux/amd64`, `user=65532:65532`, `entrypoint=["/app"]`, serves `/health` and
`/test` identically to the native binary.

## 2. What went wrong / friction

- **`sudo` again:** the daemon socket isn't reachable unprivileged, so the build/run were
  owner-run. (Recorded as a human dependency, not a defect.)
- **No `go.sum`:** a `COPY app/go.sum*` glob would have failed the build (no matches), so
  the dependency-metadata copy is `app/go.mod` only until real deps exist. A trap that
  silently breaks builds once you add a `go.sum` and forget to copy it.
- **Cross-platform default:** without `--platform=$BUILDPLATFORM` the builder stage would
  try to run under emulation (QEMU) for a foreign target. Setting the builder to the build
  platform and cross-compiling avoids QEMU entirely for Go.
- **`.dockerignore` must not starve the build:** `app/` and `app/go.mod` are build inputs;
  broad patterns that match them break the build. Verified by fnmatch against the six
  `app/*` paths.
- **Attestation by default:** buildx emitted an extra attestation manifest alongside the
  image — a preview of Stage 05's provenance work.

## 3. Lessons / assignment-required knowledge (research 3.1c)

- **What a container standardizes (§3.1c).** The same static binary runs the same way
  everywhere; the image bundles the app plus its userland contract. Multi-stage builds
  separate the heavy **builder** (toolchain) from the tiny **runtime**.
- **Layers and caching.** Each Dockerfile instruction is a cached layer; a change
  invalidates that layer and everything after it. Copying `go.mod` (dependency metadata)
  before the source keeps the module layer warm across source edits — the rebuild reused
  every layer (`CACHED`, 32.4s → 1.8s). **Order instructions from least- to most-changing.**
- **Multi-platform without emulation.** `$BUILDPLATFORM` = where the builder runs
  (native amd64 here); `$TARGETOS`/`$TARGETARCH` = the image's target. Go cross-compiles a
  **static** (`CGO_ENABLED=0`) binary for the target, so no QEMU and no dynamic loader.
- **Why deleting a secret in a later layer is unsafe.** Image layers are **additive**: if a
  secret is `COPY`d/`RUN`-written in an early layer, a later `RUN rm` only adds a
  whiteout — the bytes remain in the earlier layer and are recoverable from the image.
  Fix: never put secrets in the image; use build secrets / runtime injection.
- **Why a minimal image alone is not security.** Distroless removes shell and package
  manager (smaller attack surface, no `apt`, no interactive entry), but it does not fix
  application vulnerabilities, run-as-root, misconfiguration, or supply-chain integrity.
  Those need **non-root** (`USER 65532`), **pinned digests**, and later SBOM/attestation
  gates — each a separate control.
- **Why `/health` must be checked externally.** Distroless has no shell, so an in-image
  `HEALTHCHECK` shell command can't exist; liveness/readiness are HTTP probes from
  Kubernetes calling `/health` (Stage 06–07). Tags are mutable; **digests are immutable**,
  which is why we pin `@sha256:…`.
- **Manifest vs native run.** A multi-arch **manifest list** advertising arm64 is not a
  native arm64 execution test. Stage 03 runs native amd64 only; arm64 is a Stage 04 CI
  publication claim, never a local arm64 run.

## 4. How to run and inspect it

```sh
cd /home/thearesia/Projects/devops-worktrees/03-container
sudo docker buildx build --platform linux/amd64 --load -f app/Dockerfile -t devops-app-huy:container-proof .
sudo docker run -d --name huy-container-proof -p 127.0.0.1:18080:8080 -e APP_ENV=container-proof devops-app-huy:container-proof
sudo docker image inspect devops-app-huy:container-proof --format '{{.Os}}/{{.Architecture}} user={{.Config.User}} entrypoint={{json .Config.Entrypoint}}'
curl -i http://127.0.0.1:18080/health
curl -i http://127.0.0.1:18080/test
sudo docker run --rm --entrypoint /bin/sh devops-app-huy:container-proof   # expected: no /bin/sh
sudo docker stop huy-container-proof && sudo docker rm huy-container-proof
```

**What you can do with it:** run the exact deployed artifact locally with different
`APP_ENV` values; flip `delay_ms`/`fail` and watch the container behave identically to the
native process; prove there is no shell. This image is what CI (Stage 04) builds and
publishes to GHCR.

## 5. Fastest rebuild runbook

1. Decide bases and **resolve their index digests** from the registry (Docker Hub
   `library/golang`; `gcr.io/distroless/static-debian13:nonroot`) — pin `@sha256:…`, never
   a bare tag in the committed Dockerfile.
2. Write `app/Dockerfile`: `# syntax=docker/dockerfile:1`; builder
   `FROM --platform=$BUILDPLATFORM golang:1.27@<digest> AS build` with `ARG TARGETOS/TARGETARCH`,
   `COPY app/go.mod ./` then `COPY app/ ./`, `RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/app .`;
   runtime `FROM gcr.io/distroless/static-debian13:nonroot@<digest>`, `COPY --from=build /out/app /app`,
   `USER 65532:65532`, `EXPOSE 8080`, `ENTRYPOINT ["/app"]`.
3. Root `.dockerignore` excluding `.git`, `.omo`, `.scratch`, `captures`, `docs`, `.env*`,
   kubeconfig/keys, `node_modules`, build output — **keep `app/` and `app/go.mod`**.
4. Build with `-f app/Dockerfile .` (root context) and `--platform linux/amd64 --load`.
5. Run on a spare port; inspect `user`/`entrypoint`/arch; curl `/health` and `/test`; prove
   no shell; rebuild to show `CACHED`.

**Traps:** bare `:latest`/tags instead of digests; `COPY app/go.sum*` with no `go.sum`;
copying the whole repo (secrets/`.env`/`.omo`) into the context; running the builder under
QEMU by omitting `--platform=$BUILDPLATFORM`; a shell-based `HEALTHCHECK` in distroless;
claiming an arm64 native run from a manifest.
