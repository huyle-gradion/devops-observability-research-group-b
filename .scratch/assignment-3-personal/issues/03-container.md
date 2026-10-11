# 03: Non-root multi-platform image

Type: task
Status: ready-for-agent
Blocked by: 02

## Objective and authority

Build a native runnable distroless image using architecture-correct Go cross-compilation and pinned multi-platform bases. Governing spec: [personal spec](../spec.md). Execution: [Stage 03](../../../.omo/plans/devops-personal/03-container.md).

## Acceptance

Actual native /health and /test responses, configured non-root identity and no-shell proof; cache exercise; multi-platform output verified later in CI. Owner explains layers, cache, secret retention and manifest-versus-native-run distinction.

## Comments

Planning only; don't use a shell-based healthcheck in distroless.
