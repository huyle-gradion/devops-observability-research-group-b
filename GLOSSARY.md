# DevOps & Observability Research

The shared vocabulary for the DevOps and Observability research project (Assignment 3, Group B).

## Language

### Pipeline & Infrastructure

**Proof of Work**:
The 5-stage running deployment and verification pipeline (commit → CI → image → deploy → cluster → dashboards) operated on team laptops.
_Avoid_: Demo script, toy example

**Distroless**:
A minimal container image containing only the compiled application binary and essential runtime dependencies, excluding package managers, shells, and operating system utilities.
_Avoid_: Alpine, slim image, bare metal

**GitOps**:
A deployment model where Git is the single source of truth for declared system state, and an in-cluster agent continuously reconciles cluster drift against the repository.
_Avoid_: Push-based deployment, script deploy

### Observability

**Signal**:
One of the standardized categories of telemetry data defined by OpenTelemetry: metrics, traces, or logs.
_Avoid_: Telemetry type, monitoring data

**Semantic Conventions**:
Standardized naming rules and attribute keys for telemetry signals to ensure cross-backend and cross-vendor portability.
_Avoid_: Custom tags, arbitrary labels

**Real User Monitoring (RUM)**:
Observability captured directly from end-users' browser sessions (e.g. Navigation Timing, Core Web Vitals) reflecting real-world conditions.
_Avoid_: Synthetic testing, lab benchmarks
