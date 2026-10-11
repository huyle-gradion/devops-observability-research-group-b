# Stage 03: Non-root architecture-correct container

## Governing artifacts

[Spec](../../../.scratch/assignment-3-personal/spec.md), [ticket 03](../../../.scratch/assignment-3-personal/issues/03-container.md), assignment Step 1 and research 3.1c, overview Approach 4, [execution contract](00-execution-contract.md).

## Entry gate and routing

Stage 02 accepted and local Docker access exercised. Preliminary complex-critical because base-image trust, public build context/secrets and execution privileges are security-sensitive. Read [Docker multi-platform builds](https://docs.docker.com/build/building/multi-platform/), [Dockerfile reference](https://docs.docker.com/reference/dockerfile/), and [distroless](https://github.com/GoogleContainerTools/distroless) before worker editing. Resolve actual stable tags and index digests at execution; no placeholder digest in committed Dockerfile.

## Allowed scope and branch

Base accepted origin/personal/huy after Stage 02; branch stage/huy/03-container; matching contract worktree. Files app/Dockerfile, root .dockerignore, container instructions in README/docs/huy/versions.md/lessons.md and Step 1 container captures. Commit scope: container build/runtime contract only. No source endpoint changes, CI, cluster or credential files. Later Stage 11 may add the browser build stage.

## Bounded substeps

### 03A: Multi-stage build and safe context

Repository root is build context; Dockerfile is app/Dockerfile. Native builder FROM --platform=$BUILDPLATFORM golang:1.27@resolved-multiarch-index AS build. Declare TARGETOS and TARGETARCH; compile CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/app . in copied app source. Keep dependency metadata copy before source when modules exist to preserve cache. No hardcoded GOARCH=amd64 in the committed build.

Runtime FROM gcr.io/distroless/static-debian13:nonroot@resolved-index. COPY compiled app, USER 65532:65532, EXPOSE 8080, exec-form ENTRYPOINT ["/app"]. The target binary is static; do not choose a dynamic runtime unnecessarily. Image must not require writeable filesystem or shell startup.

.dockerignore excludes .git, .omo, .scratch, captures, .env*, credential/private-key/kubeconfig files, node_modules and local output. Do not exclude go.mod/go.sum or later intentional app/static browser bundle. No secrets through COPY/ARG/ENV. No shell/curl HEALTHCHECK in distroless; externally call /health and use Kubernetes HTTP probes later.

Exit gate: resolved digest pins and build context match source files; no shell/credential introduced.

### 03B: Actual native run and image inspection

Choose a new proof container name owned by this stage; check collision before creation. Run locally on the observed amd64 host:

```sh
sudo docker buildx build --platform linux/amd64 --load \
  -f app/Dockerfile -t devops-app-huy:container-proof .
sudo docker run -d --name huy-container-proof \
  -p 127.0.0.1:18080:8080 -e APP_ENV=container-proof \
  devops-app-huy:container-proof
sudo docker inspect huy-container-proof --format '{{.Config.Hostname}}'
sudo docker image inspect devops-app-huy:container-proof \
  --format '{{.Os}}/{{.Architecture}} user={{.Config.User}} entrypoint={{json .Config.Entrypoint}}'
curl -i http://127.0.0.1:18080/health
curl -i http://127.0.0.1:18080/test
sudo docker run --rm --entrypoint /bin/sh devops-app-huy:container-proof
```

Last command is EXPECTED to fail because no /bin/sh exists; save its real error/exit code and do not fold it into a must-pass shell chain. Match test identity with the running container hostname, not full docker inspect .Id. Inspect USER=65532:65532 and amd64. This proves configured user and native behavior; don't call it an exhaustive runtime security audit.

Exit gate: actual status/body/header, non-root image setting, exec-form startup and no-shell evidence. Stop/remove only huy-container-proof after capture.

### 03C: Cache and architecture learning experiment

Rebuild unchanged source and inspect real cache output; make a useful scoped source/comment-free change only if necessary, or use a throwaway copy to demonstrate cache invalidation without changing app behavior. Preserve working source, remove experiment copy after evidence. Do not fabricate timing improvements.

Full amd64/arm64 publication is Stage 04. Cross-compiling Go from native BUILDPLATFORM does not require QEMU for the Go binary. Future native Node builder stays on BUILDPLATFORM too. Local VM virtualization uses QEMU independently; those are different requirements. Record that an arm64 manifest is not an arm64 native execution test.

## Verification and evidence

Worker focused checks: Docker build/native response and inspection. Sisyphus repeats the runnable image path independently and inspects build context/no secret inclusion before authorized commit. Save build/cache outputs, response and expected no-shell failure in captures/huy/step-1/container.md. Record both pinned base INDEX digests and official source URLs in versions.md. Remove only own proof resources; do not prune shared Docker images/volumes.

## Learning gate

Owner builds/runs the image and explains native versus target platform, layers/cache, why deleting a secret in a later layer is unsafe, why minimal image alone is not security, and why /health must be checked externally. Demonstrate passing runtime configuration without rebuilding.

## Completion and blockers

No stage completion on Dockerfile plausibility alone. Architecture failures require fixing compiler target/index selection, not using an amd64-only image on arm64. If distroless tag availability changed, verify supported official equivalent and report a contract decision to Prometheus/owner rather than silently downgrading to shell-bearing runtime.

## Sisyphus handoff prompt

```text
Coordinate Stage 03 and personal ticket 03. Confirm Stage 02 plus Docker runtime
access; prepare stage/huy/03-container in its isolated worktree. Preliminary
complex-critical. Do not implement. Research current official multi-platform and
distroless sources, then dispatch 03A/03B bounded prompts using current Worker role.
Require real digest pins, root context, target-correct static Go build, non-root
exec-form runtime and no shell healthcheck. No CI/cluster or app contract changes.
Receive results, independently run the native image, inspect identity/user/context,
collect expected no-shell/cache evidence and owner's learning demonstration.
Commit/integrate only when authorized; do not claim arm64 was executed.
```
