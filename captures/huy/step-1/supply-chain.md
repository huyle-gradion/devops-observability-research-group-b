# Stage 05 — supply-chain evidence (huy)

Produced by Sisyphus from real runs/registries. Observed **2026-10-11**. No secrets.

Producing source: `personal/huy` @ `cd0e7d58e97587b53e617ee37abe4193d0d80ce3`
Run: `38119874453` (push) — `check`/`publish`/`smoke` all success, including the
new `Attest build provenance` step.

Artifact: `ghcr.io/huyle-gradion/devops-app-huy@sha256:0716851ad477e43ec0436514986eab47f3af49b59bccf049e2f6dd0a8fd706ca`

## 05A — signed provenance

- `publish` adds `id-token: write`; step `actions/attest-build-provenance@4d101475…`
  attests the pushed digest (`subject-name` without tag, `push-to-registry: true`,
  `create-storage-record: false`, explicit `PROJECT_GITHUB_TOKEN`).
- Attestation uploaded to the registry and to the repo
  (attestation `54753252`); Sigstore Rekor entry `3202502156`.

## 05B — SBOM + cryptographic verification

- Both platform SPDX SBOMs extracted with `docker buildx imagetools inspect` (no sudo
  needed) → `captures/huy/bonus-4/sbom-amd64.json`, `sbom-arm64.json`
  (SPDX-2.3, syft v1.51.0 / buildkit v0.33.1, 9 packages, distinct document
  namespaces + md5).
- `gh attestation verify oci://…@sha256:0716851a…` (anonymous) → **exit 0**. Verified
  identity: signer workflow
  `…/.github/workflows/ci.yml@refs/heads/personal/huy`, issuer
  `token.actions.githubusercontent.com`, trigger `push`, source SHA `cd0e7d58…`,
  subject digest `sha256:0716851a…`, SLSA predicate `https://slsa.dev/provenance/v1`.
  Saved: `captures/huy/bonus-4/attestation-verification.json`.

## 05C (partial) — vulnerability scan

- Trivy `0.75.0`, DB `mirror.gcr.io/aquasec/trivy-db:2`, scanned the public digest
  anonymously. Detected OS `debian 13.7` (distroless), 6 OS packages + 1 Go binary.
  **Result: 0 vulnerabilities.**
- Report: `captures/huy/bonus-4/trivy-image.json`.
- Interpretation: a clean result on a distroless static image is expected (tiny package
  surface) and is **not** proof of universal safety — it reflects the scanner, DB snapshot
  and date. This is a report, not an enforced gate (the assignment makes scanning optional).

## Boundaries

- arm64 is a published manifest + arm64 SBOM, **not** a native arm64 execution test.
- SBOM ≠ vulnerability scan ≠ signed provenance: three different guarantees.
