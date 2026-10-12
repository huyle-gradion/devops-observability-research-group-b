# Stage 04b teachback — the 6-job CI pipeline (huyle-gradion)

Post-implementation teaching. Written 2026-10-12. Evidence: `captures/huy/step-1/ci/ci-6jobs.md`.

## 1. What changed

`.github/workflows/ci.yml` went from 3 jobs (`check`/`publish`/`smoke`) to **6 ordered jobs**,
each gating the next via `needs`:

```
lint -> compile -> test -> build -> smoke -> publish
```

The team's "8 steps" map onto these jobs (lint=1, compile=2, test=3, build=4, local smoke+scan=5,
publish+attest+digest-smoke=6). `permissions: {}` at the top; only `publish` gets `id-token: write`
and the PAT.

## 2. What went wrong / friction (real)

- **First run failed at `smoke` (and `publish` skipped) — the gate worked.** The pinned-Trivy
  download saved the tarball as `/tmp/trivy.tar.gz` while the checksums file names
  `trivy_0.75.0_Linux-64bit.tar.gz`, so `sha256sum -c` couldn't find it. Fixed by keeping the
  original filename.
- **`mise` `gh` shim broke git auth:** `~/.local/bin/gh` prints setup noise on stdout, which
  corrupted the `gh auth git-credential` output → "invalid credential line" and failed pushes.
  Fix: put the real binary dir ahead on `PATH` for the credential helper.
- **`--load` cannot carry SBOM/provenance:** the local `docker` exporter rejects attestations
  ("docker exporter does not currently support exporting manifest lists"); the `build` job uses the
  **OCI exporter** (`outputs: type=oci,dest=…`) instead, and `smoke` builds its own runnable copy.
- **PR-after-merge gives no diff:** once `personal/huy` contains the branch, you can't open a PR —
  a separate small branch is needed to exercise the PR path.

## 3. Lessons / assignment knowledge

- **Each job is a gate.** The `needs` chain means a broken `lint`/`compile`/`test`/`build`/`smoke`
  stops the pipeline before the secret-bearing `publish` runs — proven live: run `38212715813` failed
  `smoke` and skipped `publish`.
- **Least privilege is structural, not aspirational.** `permissions: {}` + per-job grants + the
  `publish`-only `if: push && refs/heads/personal/huy` mean a PR (fork or same-repo) can *never*
  reach the PAT or OIDC.
- **Split jobs give clearer signal.** `gofmt`/`vet` (lint), a build-only `compile`, the behavior
  `test`, an image `build` with SBOM/provenance, a local `smoke`+scan, and the trusted `publish`
  each fail in isolation — easier to debug and to present than one monolithic job.
- **Pinned binaries need their exact filename for checksums.** A renamed download breaks
  `sha256sum -c`; download under the release name.
- **Attestations need a capable exporter.** BuildKit can attach SBOM/provenance to an OCI/registry
  export, not to the local `docker` image store.

## 4. How to read it / verify

```sh
gh run list --repo huyle-gradion/devops-observability-research-group-b
gh run view <id> --repo huyle-gradion/devops-observability-research-group-b --json jobs \
  --jq '.jobs[] | "\(.name): \(.conclusion)"'
```

**What you can do with it:** push to `personal/huy` → the whole six-job pipeline runs and, if green,
publishes a multi-arch image with signed provenance; open a PR → `publish` skips. The job boundaries
are the "team language" reviewers expect.

## 5. Fastest rebuild runbook

1. `permissions: {}` top-level; triggers `push`/`pull_request` to the integration branch, path-filtered.
2. Six jobs with `needs`: lint (`gofmt`+`vet`) → compile (`CGO_ENABLED=0 go build`) → test (`go test`)
   → build (`build-push` no push, `outputs: type=oci`, `sbom`+`provenance`) → smoke (local build, run
   on loopback, curl asserts, pinned Trivy report) → publish (`if` push to the integration branch;
   login + multi-arch push + `attest-build-provenance` + digest re-smoke; the only job with the PAT).
3. Anonymous checkout in every job; `setup-go` with `token: ''`; keep the PAT + `id-token: write`
   only on `publish`; pin every action to a verified commit.

**Traps:** PR jobs touching the secret; a renamed checksum download; `--load` with attestations;
treating a green push as proof a PR would publish.
