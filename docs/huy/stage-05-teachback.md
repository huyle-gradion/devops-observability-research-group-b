# Stage 05 teachback — SBOM, signed provenance and the red/green gate (huyle-gradion)

Post-implementation teaching. Written 2026-10-11. Evidence: `captures/huy/step-1/supply-chain.md`,
`captures/huy/bonus-4/`.

## 1. What Stage 05 delivered

- Signed **SLSA build provenance** on publish (`actions/attest-build-provenance`, OIDC keyless).
- Both platform **SPDX SBOMs** extracted; provenance **verified** with expected identity/source/digest.
- A **deliberate red→green** experiment proving a failing check skips publication.
- A real **Trivy** scan of the published digest.

## 2. What went wrong / friction

- **The negative experiment didn't fail at first.** The behavior test compares `/health`
  against the app's `healthBody` **constant**, so changing the constant is invisible to the
  test. Only a real wire change (hardcoding a different body in the handler) tripped it.
  Lesson: tests that assert against the same constant they're supposed to protect can't
  catch a change to that constant's value — pin the literal bytes.
- **Sudo-created temp dir:** a `docker --config` dir created by `sudo docker` couldn't be
  `rm -rf`'d as the user; needs `sudo rm`.
- **Action naming:** the plan said `actions/attest`, but that action is a *generic* attestation
  primitive needing a predicate; the correct action for signed provenance is
  `actions/attest-build-provenance` (same inputs). Clarified and pinned.
- **BuildKit provenance ≠ GitHub workflow provenance.** The image already carried a BuildKit
  SBOM/provenance manifest; the GitHub-signed attestation is a *separate* object with a
  different guarantee.

## 3. Lessons / assignment-required knowledge (research 3.1a/3.1c, bonus 4)

- **SBOM (Software Bill of Materials).** A machine-readable inventory of everything in the
  image (OS packages + language deps), generated at build time (here by BuildKit via syft)
  in SPDX or CycloneDX. It answers "what's inside?" — not "is it safe?".
- **Vulnerability scan.** Trivy matches package versions against a vulnerability database and
  reports CVEs + fix availability. It is a **report, not a guarantee**: distroless has a tiny
  surface (our scan was **0 vulns**), but scans can miss app-logic bugs, may lag the DB, and
  are only as good as the DB snapshot/date. A clean scan is not "vulnerability-free".
- **Signed provenance.** An in-toto/SLSA **attestation** that says *which* workflow identity
  built *which* digest from *which* source commit; signed **keylessly** with Sigstore (Fulcio
  cert + Rekor transparency log). Verification checks the **expected identity/workflow/source
  SHA/digest** — a valid signature alone is not enough; you must pin *who* is allowed to sign.
- **OIDC identity vs PAT.** The **PAT** (`PROJECT_GITHUB_TOKEN`) authorizes API/registry calls;
  the **OIDC token** (`id-token: write`) is the short-lived signing identity from GitHub Actions.
  They are separate: you cannot sign with the PAT, and you don't call the API with OIDC.
- **Digest identity.** The attestation subject is `name + sha256:…`; tags are mutable, digests
  are not — always reference the digest.
- **Publication is gated.** A failing `check` **skips** `publish`/`smoke` (proven red→green):
  broken code can never reach the registry from `personal/huy`.
- **Package vs repo visibility are independent.** A public repo does not make its GHCR package
  public (we set the package public explicitly to allow the anonymous pull).
- **SBOM ≠ scan ≠ provenance** — three different questions: *what's in it*, *what's known-bad*,
  *who built it from what*.

## 4. How to run / verify it

```sh
IMG=ghcr.io/huyle-gradion/devops-app-huy@sha256:0716851ad477e43ec0436514986eab47f3af49b59bccf049e2f6dd0a8fd706ca
# SBOM (no sudo needed)
docker buildx imagetools inspect "$IMG" --format '{{ json (index .SBOM "linux/amd64").SPDX }}'
# signed provenance
python3 scripts/project-github.py gh attestation verify "oci://$IMG" \
  --repo huyle-gradion/devops-observability-research-group-b \
  --signer-workflow huyle-gradion/devops-observability-research-group-b/.github/workflows/ci.yml \
  --source-digest cd0e7d58e97587b53e617ee37abe4193d0d80ce3
# vulnerability scan
trivy image "$IMG"
```

**What you can do with it:** prove *to a third party* that this exact image came from your
repo/workflow/commit and list exactly what's inside it, then show the gate that stops broken
code from ever being published.

## 5. Fastest rebuild runbook

1. In `publish`, add `id-token: write`; after build-push add
   `actions/attest-build-provenance@<pin>` with `subject-name` (no tag), `subject-digest:
   ${{ steps.build.outputs.digest }}`, `push-to-registry: true`, `create-storage-record: false`,
   explicit `github-token`.
2. Keep SBOM/provenance on in build-push (`sbom: true`, `provenance: mode=max`).
3. After a green run: extract both platform SPDX SBOMs with `imagetools inspect`.
4. Verify: `gh attestation verify` with `--signer-workflow` + `--source-digest`.
5. Scan with pinned Trivy; record scanner + DB versions and date.
6. Prove the gate: push a deliberately failing revision to the integration branch, confirm
   `check` fails and `publish` skips, then restore green.

**Traps:** treating a valid signature as "trusted" without pinning the expected signer;
thinking the PAT signs; assuming a clean scan means safe; assuming a public repo makes the
package public; leaving the red revision live; tests that compare against the same constant
they should guard.
