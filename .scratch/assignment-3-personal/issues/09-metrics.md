# 09: Actual app metrics and persistent workload dashboard

Type: task
Status: ready-for-agent
Blocked by: 08

## Objective and authority

Prometheus client instrumentation plus cluster pod resources, monitored in a separate namespace using kube-prometheus-stack. Governing spec: [personal spec](../spec.md). Execution: [Stage 09](../../../.omo/plans/devops-personal/09-metrics.md).

## Acceptance

Actual scrape target Up, real request/error/latency and per-pod CPU/memory panels, working dashboard JSON, PVC-backed persistence across pod replacement. No fabricated values or stock dashboard substituted for own metrics.

## Comments

Planning only; Stage 12 adds load-driven proof.
