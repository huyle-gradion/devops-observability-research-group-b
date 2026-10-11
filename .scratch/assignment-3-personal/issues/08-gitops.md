# 08: Both GitOps modes and reviewed CI digest updates

Type: task
Status: ready-for-agent
Blocked by: 07

## Objective and authority

Argo CD watches env/huy as sole app desired state; CI proposes reviewed digest PRs without cluster credentials. Governing spec: [personal spec](../spec.md). Execution: [Stage 08](../../../.omo/plans/devops-personal/08-gitops.md).

## Acceptance

Distinct actual source-to-image-to-desired-state-to-rollout loops for manual and automatic sync, UI/CLI evidence, drift restoration and Git-based rollback; owner explains the approval boundary. No duplicate authoritative manifests on source branch.

## Comments

Planning only; completion before CP2 Mon Oct 12 17:00.
