# syntax=docker/dockerfile:1

# Stage 03A: architecture-correct, non-root, shell-less runtime image for the
# stdlib Go API in app/. The repository root is the build context, so build with
# the Dockerfile path:
#   sudo docker buildx build --platform linux/amd64 --load \
#     -f app/Dockerfile -t devops-app-huy:container-proof .

# --- Builder: runs natively on the build platform, cross-compiles to the target ---
FROM --platform=$BUILDPLATFORM golang:1.27@sha256:e432b43af23a9328d56a7c499be0476810aa344acbcf65fc7c455d4ff5a40602 AS build

# Automatic platform args supplied by BuildKit; no architecture is hardcoded.
ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

# Dependency metadata first so the module layer stays cached across source edits.
# app/go.sum is added here once the module graph gains dependencies; there is no
# go.sum today, so it is intentionally omitted rather than copied with a glob
# that would fail the build.
COPY app/go.mod ./

# Application source (main.go and tests).
COPY app/ ./

# Static binary for the requested target platform. CGO is disabled so the result
# needs no dynamic loader and runs on the distroless static runtime.
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/app .

# --- Runtime: distroless static, non-root identity, exec-form entrypoint ---
FROM gcr.io/distroless/static-debian13:nonroot@sha256:e2e927ec666bae08560abb3c55d0659eceabb657f56b6782ab500a9fc7f555e3

COPY --from=build /out/app /app

# distroless :nonroot is UID/GID 65532. Explicit numeric form documents intent.
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/app"]
