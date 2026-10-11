# 06: Authorized two-physical-host k3s cluster

Type: task
Status: needs-info
Blocked by: 01

## Objective and authority

Own Ubuntu server VM plus isolated worker VM on another physical host, connected privately and running the same stable k3s release. Governing spec: [personal spec](../spec.md). Execution: [Stage 06](../../../.omo/plans/devops-personal/06-cluster.md).

## Missing prerequisites

Actual remote worker address/access, allocation resources, different physical host, installation permission and tailnet authorization. Repo tools cannot supply these. Local preparation may proceed within authorized bounds; remote installation is blocked until supplied.

## Acceptance

Two Ready nodes and physical-host mapping; private API/kubeconfig/token boundaries; Stage 07 proves real cross-node app connectivity and pod placement to close bonus 2.

## Comments

Planning only; a second local VM cannot satisfy this ticket.
