# syntax=docker/dockerfile:1

# Stage 03A: architecture-correct, non-root, shell-less runtime image for the
# stdlib Go API in app/. The repository root is the build context, so build with
# the Dockerfile path:
#   sudo docker buildx build --platform linux/amd64 --load \
#     -f app/Dockerfile -t devops-app-huy:container-proof .

# --- Web: build the architecture-independent browser bundle natively ---------
# The browser bundle is plain JavaScript, identical for every target
# architecture, so this stage runs on $BUILDPLATFORM regardless of the image
# platform being built. Stage 11 makes the generated app/static/app.js a real
# build input for the Go stage's go:embed, so it must exist before go build.
FROM --platform=$BUILDPLATFORM node:26@sha256:32fa97f3363975684b08bf4e8a68a47c7905175cc20275b50b974bbd02aba731 AS web

WORKDIR /web

# Manifest + lockfile first so the npm layer stays cached across source edits.
# npm ci installs the exact locked tree; the committed lockfile makes the bundle
# reproducible.
COPY web/package.json web/package-lock.json ./
RUN npm ci

# The real browser source. The committed build script bundles src/main.js to
# ../app/static/app.js, so this uses the identical esbuild invocation as the
# local `npm run build`. The output directory is created first so the location
# does not depend on esbuild's mkdir behavior.
COPY web/ ./
RUN mkdir -p /app/static && npm run build

# --- Builder: runs natively on the build platform, cross-compiles to the target ---
FROM --platform=$BUILDPLATFORM golang:1.27@sha256:e432b43af23a9328d56a7c499be0476810aa344acbcf65fc7c455d4ff5a40602 AS build

# Automatic platform args supplied by BuildKit; no architecture is hardcoded.
ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

# Dependency metadata first so the module layer stays cached across source edits.
# app/go.sum arrives with the source copy below, before go build.
COPY app/go.mod ./

# Application source: main.go and tests, plus the committed static/index.html.
# The gitignored app/static/app.js is a build output and is copied from the web
# stage below, not taken from this context.
COPY app/ ./

# Generated browser artifact from the web stage, copied in BEFORE go build so
# go:embed can find app/static/app.js. This runs after the source copy so the
# freshly built bundle always wins over any stale local file in the context.
COPY --from=web /app/static/app.js /src/static/app.js

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
