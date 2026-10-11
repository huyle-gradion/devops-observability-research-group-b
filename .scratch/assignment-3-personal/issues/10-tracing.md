# 10: Real backend trace and correlated request log

Type: task
Status: ready-for-agent
Blocked by: 09

## Objective and authority

Go OTel traces through validated Collector OTLP HTTP to Jaeger, with trace-context JSON stdout logging and Server-Timing. Governing spec: [personal spec](../spec.md). Execution: [Stage 10](../../../.omo/plans/devops-personal/10-tracing.md).

## Acceptance

Actual delayed request/child span in viewer, matching valid trace/span IDs in its log, measured Server-Timing, working exporter configuration and honest transient-storage/SDK maturity statements.

## Comments

Planning only; no central log backend or durable Jaeger claim.
