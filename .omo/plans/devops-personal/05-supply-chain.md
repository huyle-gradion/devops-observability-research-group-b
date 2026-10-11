# Stage 05: SBOM, signed provenance and CP1 security proof

## Governing artifacts

[Spec](../../../.scratch/assignment-3-personal/spec.md), [ticket 05](../../../.scratch/assignment-3-personal/issues/05-supply-chain.md), Step 1/bonus 4/research a,c, overview Approach 5, [execution contract](00-execution-contract.md), [credential policy](../../../docs/agents/github-access.md).

## Entry gate and routing

Stage 04 actual publication/smoke accepted; project token's workflow scope and CI secret resolved by owner. Preliminary complex-critical: signing/identity, credential scope and deliberately failing source revision. Read [actions/attest](https://github.com/actions/attest), [SBOM inspection](https://docs.docker.com/build/metadata/attestations/sbom/), [CLI verification](https://cli.github.com/manual/gh_attestation_verify) and official Trivy docs. Keep current API/signing facts distinct from unverified assignment Immutable Actions claims.

## Allowed scope and branch

Base accepted origin/personal/huy; stage/huy/05-supply-chain in isolated worktree. Files .github/workflows/ci.yml, a tightly bounded temporary app/main.go response mutation for the negative exercise, docs/huy versions/lessons/CP1 findings and captures/huy/step-1/, bonus-4/. Scope includes project-account CI/API operations via the same .env-derived secret only. No key files, unrelated repo/package changes, policy engine, cluster deployment or invented scan exceptions. Commit scope: signing workflow/evidence; separate owner-authorized negative and restoring commits, not normal unverified source integration.

## Bounded substeps

### 05A: Signed workflow provenance

After Buildx publication, call verified full-SHA-pinned actions/attest with subject-name IMAGE_NAME without tag, subject-digest build output, push-to-registry true, github-token secrets.PROJECT_GITHUB_TOKEN explicitly. Use create-storage-record false to avoid assuming unverified PAT permission for the current optional artifact metadata endpoint. Registry access uses the same token. Grant publish id-token: write for keyless signing; keep other built-in token permissions disabled. PAT API access and OIDC workflow identity are separate; do not attempt to sign using the PAT itself.

Native digest smoke stays part of the successful publish job. Expose its successful digest only after all signing/smoke steps pass. A BuildKit provenance predicate is not, by itself, signed GitHub workflow provenance. Capture actual bundle/attestation link and signer/source identity.

Exit gate: real source-push run finishes attestation and smoke successfully, signed object is retrievable, no implicit github.token API fallback.

### 05B: Inspect and cryptographically verify own artifact

Resolve REPO, SOURCE_SHA, RUN_ID and IMAGE_REF from the actual green source build. Save outputs in the named capture folders, not token-bearing console dumps. Authenticate registry reads using only the credential helper's temporary private Docker config if needed; anonymously readable content does not require a substitute credential.

```sh
REGISTRY_CONFIG=$(mktemp -d)
python3 scripts/project-github.py registry-login --config "$REGISTRY_CONFIG"
sudo docker --config "$REGISTRY_CONFIG" buildx imagetools inspect "$IMAGE_REF" \
  --format '{{ json (index .SBOM "linux/amd64").SPDX }}'
sudo docker --config "$REGISTRY_CONFIG" buildx imagetools inspect "$IMAGE_REF" \
  --format '{{ json (index .SBOM "linux/arm64").SPDX }}'
DOCKER_CONFIG="$REGISTRY_CONFIG" python3 scripts/project-github.py gh attestation verify "oci://$IMAGE_REF" \
  --repo "$REPO" --signer-workflow "$REPO/.github/workflows/ci.yml" \
  --source-digest "$SOURCE_SHA" --format json
```

Both platform JSON outputs must contain actual SPDX inventory, not null or a flag screenshot. Verification must succeed with expected GitHub Actions issuer/workflow, repo, producing source SHA and artifact digest. CLI OCI verification needs registry read authentication according to its docs; integrate the helper's private config in the verifier's process, not root's cached identity. Record exact auth/config method and scrub/remove temporary credentials after the session. A failed verification blocks bonus 4 signed proof.

### 05C: Deliberate red/fixed CI and scan report

Obtain explicit authorization for a temporary negative source push. Worker makes the bounded response change without changing tests. Sisyphus checks that the planned negative test actually fails, then may commit/push that isolated experiment to personal/huy solely to capture check failure and publish skipped. No env update/deploy is allowed. Restore immediately in a separate source commit, wait for fixed green signing/smoke and retain both real run URLs. This is the documented exception to normal green-source integration, not a loophole for unrelated failing checks. PR-only red run does not prove source-push publication gating; capture event type.

Run pinned/stable official Trivy against the successful image, with project-token registry auth only if needed. Save actual report, exit status, scanner/database versions and checked date. If download/auth/database fails, record the failure; do not call it vulnerability-free. Interpret severities/fix availability and limitations. The assignment makes scan enforcement optional; this stage reports, not invents a hard gate/allowlist.

## Verification and evidence

Sisyphus independently inspects SBOMs, runs verifier and re-pulls/runs exact green digest on native workstation. Reuse Stage 03 runtime inspection with published artifact and anonymous configuration for the public-pull proof. Capture no-shell/non-root/cache evidence, workflow event/gates, red/fixed results, SPDX files and verification JSON linked to the SAME successful source/digest. Use captures/huy/bonus-4/{sbom-amd64.json,sbom-arm64.json,attestation-verification.json} and own scan report. Captures are committed incrementally only after credential scrub/review.

## Learning gate

Owner verifies the artifact and explains SBOM versus vulnerability scan versus signed provenance, why verification checks expected identity/digest, which token authorizes API calls versus which identity signs, what a scan can miss, and why public package visibility is independent of source visibility. Owner shows the red/fixed gate and actual runnable image, not just YAML.

## Completion and CP1 gate

Before Mon 08:15, require Steps 02–05 actual acceptance, signed/SBOM evidence, green public runnable digest and owner's explanations. At CP1 Mon 09:00 show the live native artifact and evidence. A later source change requires the producing SHA/digest to remain identified; never attribute a previous image to a capture-only commit. No arm64 run claim without an actual machine.

## Sisyphus handoff prompt

```text
Coordinate Stage 05 and ticket 05 after Stage 04. Prepare the named branch/worktree;
preliminary complex-critical. Do not implement. Use the project's .env-derived
PROJECT_GITHUB_TOKEN explicitly for all GitHub/GHCR calls and OIDC only for signing.
Dispatch bounded signing and inspection subtasks; require current pinned action,
actual two-platform SPDX and expected identity/source/digest verification. Obtain
separate owner approval before the negative source-push experiment; publish/deploy
must stay skipped on red and restored source must go green. Independently verify,
review actual scan evidence and require owner learning. Commit only authorized
scope, freeze the runnable CP1 packet before 08:15, and retain real failure evidence.
```
