# Bonus 1 — CI proposes one reviewed digest PR against env/huy (huy)

Prepared by the implementation worker on `stage/huy/08b-propose-deploy` (base
`origin/personal/huy` at `4561d9d523d47547c2486ad28dab1f4b8107083f`). Observed
**2026-10-12**. No secrets; no cluster access.

This capture records local preparation and simulation only. It is **not** proof
that a real workflow run passed. The real `propose-deploy` run and the PR it opens
are pushed and observed by **Sisyphus** after integration (see "Unrun checks").

## Purpose

After the whole `lint..publish` chain is green on a trusted push to
`personal/huy`, propose exactly one **reviewed** desired-state change against
`env/huy`: a PR that moves only the app container image in
`k8s/app/deployment.yaml` to the digest `publish` just built, signed and
smoke-tested. The job opens a proposal; a human merge is the approval boundary.
CI holds no cluster credential and never auto-merges.

## Design

- **Additive.** The six existing jobs (`lint -> compile -> test -> build -> smoke
  -> publish`) are byte-for-byte unchanged. `propose-deploy` is appended as the
  7th job.
- **Gate.** `needs: publish` plus
  `if: github.event_name == 'push' && github.ref == 'refs/heads/personal/huy'`.
  A `pull_request` event never reaches it, and because a job whose `needs` job
  failed or was skipped is itself skipped, a failed/skipped `publish` cannot
  propose. No `always()` is used (that would defeat the skip semantics).
- **Credentials.** No `permissions:` block is added, so the job inherits the
  top-level `permissions: {}` (built-in token denied). Everything is authorized
  with the project PAT explicitly: `GH_TOKEN`/`GITHUB_TOKEN` env for `gh`, and a
  **process-scoped** Git HTTPS credential — an `http.https://github.com/.extraheader`
  set through git's `GIT_CONFIG_*` environment config, so the base64
  `x-access-token:<token>` value never appears in argv, a remote URL, or any file
  on disk, and dies with the step. No SSH key, no kubeconfig, no cluster
  credential.
- **Branch + PR idempotency.** Deploy branch `deploy/huy-<${ github.sha }>` (full
  source SHA). The job detects an existing branch/open PR via `git ls-remote` and
  `gh pr list`, checks the branch out, and adds a **fast-forward** commit only if
  the desired digest differs. An existing open PR is reused (never duplicated).
  If the branch is not fast-forwardable from `env/huy` (env moved on / diverged),
  the job stops with both tips printed for human review — never reset, rebase,
  merge, force-push, delete or silently resolve.
- **Validate before proposing.** Empty publish digest, non-`sha256:` digest, empty
  source SHA, or an unexpected repo/image identity all fail closed. Empty digest
  proposes nothing.
- **Manifest edit.** `kubectl set image --local -f k8s/app/deployment.yaml
  app=<IMAGE>@<DIGEST> -o yaml` runs against a temp file (no cluster), then the
  tracked file is replaced by a reviewed, single-line splice of the image scalar
  (see decision 1). A `git diff` gate requires exactly one file and a `+1/-1`
  change before committing. If the desired digest is already present on a fresh
  branch, the job makes no commit and no PR and exits 0.
- **PR.** `gh pr create --base env/huy --head deploy/huy-<sha>` with a body
  carrying the source SHA, the green run URL, the digest and the verification/
  attestation linkage. Never auto-merged; `env/huy` is never pushed directly.
- **Event/filter boundary.** Only a `push`/`pull_request` targeting `personal/huy`
  enters the workflow; `env/huy` and `deploy/*` are not in `on:`, and the paths
  filter excludes `k8s/**`. A PAT-created PR *can* trigger workflows, so those
  explicit branch+path filters — not the branch names — are the guard. Nothing in
  `on:` was changed.

## Pins

- **No new GitHub Action.** Every `propose-deploy` step is a plain `run:` shell
  step, so the five action pins from Stages 04/05/04b are unchanged.
- **kubectl** `v1.37.1`, from
  `https://dl.k8s.io/release/v1.37.1/bin/linux/amd64/kubectl`, verified against
  the release `kubectl.sha256`
  (`65691ff77eb6fa44c908b77a1082c9f092c3b9733b5cefabec0d1104890e21a8`) before
  installation. Matches the k3s server minor (`v1.37.1+k3s1`). Recorded in
  `docs/huy/versions.md`.

## The new job YAML (verbatim from `.github/workflows/ci.yml`)

```yaml

  # --- step 7: propose reviewed digest PR against env/huy ---------------------
  # Trusted source-push only, and only when the whole lint..publish chain is
  # green: `needs: publish` plus the same push/ref guard means a pull_request
  # event never reaches this job, and a failed or skipped `publish` skips this
  # job automatically (a job that needs a failed or skipped job is skipped too).
  # This is the second and last job that sees PROJECT_GITHUB_TOKEN. It holds NO
  # cluster credential -- no kubeconfig, no SSH key, no service-account token --
  # and it never auto-merges; a human merges the reviewed PR.
  #
  # Event/filter boundary (deliberate): only a push or PR that targets
  # personal/huy enters this workflow at all. env/huy and deploy/* are NOT in
  # `on:`, so a push to them cannot trigger this workflow, and the paths filter
  # (app/**, web/**, ci.yml, .dockerignore) excludes k8s/**. A PAT-created PR
  # CAN trigger other workflows, so those explicit branch+path filters -- not the
  # branch names -- are the real guard. Do not add env/huy or deploy/* to any
  # trigger, and do not add an `env/huy` source-build recursion.
  propose-deploy:
    name: propose-deploy
    needs: publish
    if: github.event_name == 'push' && github.ref == 'refs/heads/personal/huy'
    runs-on: ubuntu-latest
    env:
      DEPLOY_BRANCH: deploy/huy-${{ github.sha }}
    steps:
      - name: Validate publish digest, source SHA and identity
        env:
          DIGEST: ${{ needs.publish.outputs.digest }}
          SOURCE_SHA: ${{ github.sha }}
          REPO_SLUG: ${{ github.repository }}
        run: |
          set -euo pipefail
          if [ -z "${DIGEST}" ]; then
            echo "ERROR: publish digest is empty; refusing to propose anything" >&2
            exit 1
          fi
          case "${DIGEST}" in
            sha256:*) ;;
            *) echo "ERROR: digest '${DIGEST}' is not a sha256:... reference" >&2; exit 1 ;;
          esac
          if [ -z "${SOURCE_SHA}" ]; then
            echo "ERROR: source SHA is empty" >&2
            exit 1
          fi
          if [ "${REPO_SLUG}" != "huyle-gradion/devops-observability-research-group-b" ]; then
            echo "ERROR: unexpected repository '${REPO_SLUG}'" >&2
            exit 1
          fi
          if [ "${IMAGE_NAME}" != "ghcr.io/huyle-gradion/devops-app-huy" ]; then
            echo "ERROR: unexpected image identity '${IMAGE_NAME}'" >&2
            exit 1
          fi
          echo "inputs ok: repo=${REPO_SLUG} sha=${SOURCE_SHA} digest=${DIGEST}"

      - name: Install pinned kubectl (local manifest editing only)
        run: |
          set -euo pipefail
          KUBECTL_VERSION=v1.37.1
          base="https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/amd64"
          curl -fsSL -o /tmp/kubectl "${base}/kubectl"
          curl -fsSL -o /tmp/kubectl.sha256 "${base}/kubectl.sha256"
          ( cd /tmp && echo "$(cat kubectl.sha256)  kubectl" | sha256sum -c - )
          sudo install -m 0755 /tmp/kubectl /usr/local/bin/kubectl
          kubectl version --client=true

      - name: Prepare env/huy deploy branch and stage the one image change
        id: prepare
        env:
          GH_TOKEN: ${{ secrets.PROJECT_GITHUB_TOKEN }}
          GITHUB_TOKEN: ${{ secrets.PROJECT_GITHUB_TOKEN }}
          DIGEST: ${{ needs.publish.outputs.digest }}
          SOURCE_SHA: ${{ github.sha }}
          RUN_URL: ${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}
        run: |
          set -euo pipefail
          WANT_IMAGE="${IMAGE_NAME}@${DIGEST}"

          # Process-scoped Git HTTPS credential: a repo-scoped Authorization
          # header supplied through git's own environment config. It never lands
          # in argv, in a remote URL, or in any file on disk, and it dies with
          # this shell. GIT_CONFIG_* is cleared by the runner between steps.
          AUTH="AUTHORIZATION: basic $(printf '%s' "x-access-token:${GH_TOKEN}" | base64 | tr -d '\n')"
          export GIT_CONFIG_COUNT=1
          export GIT_CONFIG_KEY_0="http.https://github.com/.extraheader"
          export GIT_CONFIG_VALUE_0="${AUTH}"

          WORK="${RUNNER_TEMP:-/tmp}/env-propose"
          rm -rf "${WORK}"
          mkdir -p "${WORK}"
          cd "${WORK}"

          git init --quiet .
          git remote add origin "https://github.com/${GITHUB_REPOSITORY}.git"
          git fetch --no-tags --quiet origin \
            "refs/heads/env/huy:refs/remotes/origin/env/huy"

          branch_exists=false
          if git ls-remote --heads origin "refs/heads/${DEPLOY_BRANCH}" | grep -q .; then
            branch_exists=true
          fi

          if [ "${branch_exists}" = true ]; then
            git fetch --no-tags --quiet origin \
              "refs/heads/${DEPLOY_BRANCH}:refs/remotes/origin/${DEPLOY_BRANCH}"
            # Fast-forwardable only: env/huy's tip must be an ancestor of the
            # existing deploy branch. If env/huy moved on, or the branch
            # diverged, do NOT reset / rebase / merge / force-push -- stop so a
            # human can review the two tips.
            if ! git merge-base --is-ancestor \
                 refs/remotes/origin/env/huy "refs/remotes/origin/${DEPLOY_BRANCH}"; then
              echo "ERROR: ${DEPLOY_BRANCH} is not fast-forwardable from env/huy" >&2
              echo "  env/huy          $(git rev-parse refs/remotes/origin/env/huy)" >&2
              echo "  ${DEPLOY_BRANCH} $(git rev-parse refs/remotes/origin/${DEPLOY_BRANCH})" >&2
              echo "STOP: refusing to reset, force-push or silently resolve; needs human review" >&2
              exit 1
            fi
            git checkout --quiet -B "${DEPLOY_BRANCH}" "refs/remotes/origin/${DEPLOY_BRANCH}"
          else
            git checkout --quiet -B "${DEPLOY_BRANCH}" refs/remotes/origin/env/huy
          fi

          current_image="$(
            python3 -c 'import yaml; d=yaml.safe_load(open("k8s/app/deployment.yaml")); print([c for c in d["spec"]["template"]["spec"]["containers"] if c["name"] == "app"][0]["image"])'
          )"
          echo "current image on ${DEPLOY_BRANCH}: ${current_image}"

          if [ "${current_image}" = "${WANT_IMAGE}" ]; then
            echo "desired digest already present; no commit and no manifest change"
            echo "changed=false" >> "${GITHUB_OUTPUT}"
            if [ "${branch_exists}" = true ]; then
              echo "propose=true" >> "${GITHUB_OUTPUT}"
            else
              echo "propose=false" >> "${GITHUB_OUTPUT}"
            fi
            exit 0
          fi

          # kubectl is the authority for the target image string, but its YAML
          # output re-serialises the document (reorders fields, drops comments,
          # adds strategy/status). So take only the computed image scalar from
          # it and splice exactly that one value into the tracked file; every
          # other byte, including comments, is preserved. This is a scoped,
          # reviewed replacement of the one manifest -- never a regex across
          # manifests.
          kubectl set image --local -f k8s/app/deployment.yaml \
            "app=${WANT_IMAGE}" -o yaml > /tmp/deployment.kubectl.yaml

          python3 - "${WANT_IMAGE}" <<'PY'
          import copy
          import sys
          import yaml

          want = sys.argv[1]
          path = "k8s/app/deployment.yaml"

          with open(path) as fh:
              original_text = fh.read()
          original = yaml.safe_load(original_text)
          with open("/tmp/deployment.kubectl.yaml") as fh:
              kubectl_doc = yaml.safe_load(fh)

          def app_container(doc):
              containers = doc["spec"]["template"]["spec"]["containers"]
              matches = [c for c in containers if c["name"] == "app"]
              if len(matches) != 1:
                  raise SystemExit("expected exactly one container named 'app'")
              return matches[0]

          if app_container(kubectl_doc)["image"] != want:
              raise SystemExit("kubectl produced an unexpected image")

          old_image = app_container(original)["image"]
          lines = original_text.splitlines(keepends=True)
          hits = [i for i, line in enumerate(lines) if old_image in line]
          if len(hits) != 1:
              raise SystemExit("expected exactly one line carrying the old image, found %d" % len(hits))
          lines[hits[0]] = lines[hits[0]].replace(old_image, want)
          new_text = "".join(lines)

          checked = yaml.safe_load(new_text)
          if app_container(checked)["image"] != want:
              raise SystemExit("spliced image is wrong")

          before = copy.deepcopy(original)
          after = copy.deepcopy(checked)
          popped_before = app_container(before).pop("image")
          popped_after = app_container(after).pop("image")
          if before != after or popped_before != old_image or popped_after != want:
              raise SystemExit("splice changed more than the app container image")

          with open(path, "w") as fh:
              fh.write(new_text)
          print("staged image change: %s -> %s" % (old_image, want))
          PY

          git add k8s/app/deployment.yaml
          changed_files="$(git diff --cached --name-only)"
          if [ "${changed_files}" != "k8s/app/deployment.yaml" ]; then
            echo "ERROR: unexpected changed files:" >&2
            echo "${changed_files}" >&2
            exit 1
          fi
          added="$(git diff --cached --numstat -- k8s/app/deployment.yaml | cut -f1)"
          removed="$(git diff --cached --numstat -- k8s/app/deployment.yaml | cut -f2)"
          if [ "${added}" != "1" ] || [ "${removed}" != "1" ]; then
            echo "ERROR: expected a single-line image change, got +${added}/-${removed}" >&2
            git diff --cached -- k8s/app/deployment.yaml >&2
            exit 1
          fi

          git -c user.name="github-actions[bot]" \
              -c user.email="41898282+github-actions[bot]@users.noreply.github.com" \
              commit --quiet \
              -m "deploy: pin devops-app image to ${DIGEST}" \
              -m "Proposed by CI run ${RUN_URL} from source ${SOURCE_SHA}. Changes only the app container image in k8s/app/deployment.yaml; promoted to desired state by a reviewed PR merge into env/huy."

          echo "changed=true" >> "${GITHUB_OUTPUT}"
          echo "propose=true" >> "${GITHUB_OUTPUT}"

      - name: Push the scoped deploy branch (fast-forward only)
        if: steps.prepare.outputs.changed == 'true'
        env:
          GH_TOKEN: ${{ secrets.PROJECT_GITHUB_TOKEN }}
          GITHUB_TOKEN: ${{ secrets.PROJECT_GITHUB_TOKEN }}
        run: |
          set -euo pipefail
          AUTH="AUTHORIZATION: basic $(printf '%s' "x-access-token:${GH_TOKEN}" | base64 | tr -d '\n')"
          export GIT_CONFIG_COUNT=1
          export GIT_CONFIG_KEY_0="http.https://github.com/.extraheader"
          export GIT_CONFIG_VALUE_0="${AUTH}"
          cd "${RUNNER_TEMP:-/tmp}/env-propose"
          # Plain push, never --force: a non-fast-forward remote is rejected and
          # fails the job for human review instead of overwriting anything.
          git push origin "HEAD:refs/heads/${DEPLOY_BRANCH}"

      - name: Open or reuse exactly one digest PR against env/huy
        if: steps.prepare.outputs.propose == 'true'
        env:
          GH_TOKEN: ${{ secrets.PROJECT_GITHUB_TOKEN }}
          GITHUB_TOKEN: ${{ secrets.PROJECT_GITHUB_TOKEN }}
          DIGEST: ${{ needs.publish.outputs.digest }}
          SOURCE_SHA: ${{ github.sha }}
          RUN_URL: ${{ github.server_url }}/${{ github.repository }}/actions/runs/${{ github.run_id }}
        run: |
          set -euo pipefail
          REPO="${GITHUB_REPOSITORY}"

          existing="$(gh pr list --repo "${REPO}" --head "${DEPLOY_BRANCH}" \
            --base env/huy --state open --json number --jq '.[].number')"
          if [ -n "${existing}" ]; then
            echo "reusing existing open PR #${existing} for ${DEPLOY_BRANCH}; no duplicate"
            exit 0
          fi

          if ! git ls-remote --heads "https://github.com/${REPO}.git" \
               "refs/heads/${DEPLOY_BRANCH}" | grep -q .; then
            echo "remote branch ${DEPLOY_BRANCH} does not exist; nothing to propose"
            exit 0
          fi

          {
            echo "## Reviewed desired-state digest update"
            echo
            echo "- Source commit: \`${SOURCE_SHA}\`"
            echo "- Green trusted CI run: ${RUN_URL}"
            echo "- Published image: \`${IMAGE_NAME}@${DIGEST}\`"
            echo "- Changed file: \`k8s/app/deployment.yaml\` (app container image only)"
            echo
            echo "### Verification / attestation linkage"
            echo
            echo "The image was built, pushed and digest-smoke-tested by the \`publish\` job,"
            echo "and signed SLSA build provenance was attested there with"
            echo "\`actions/attest-build-provenance\` (keyless OIDC, registry push) under the"
            echo "same run linked above. Verify locally with:"
            echo
            echo '```'
            echo "gh attestation verify oci://${IMAGE_NAME}@${DIGEST} --repo ${REPO}"
            echo '```'
            echo
            echo "### Approval boundary"
            echo
            echo "CI opened this PR but never auto-merges and never writes \`env/huy\`"
            echo "directly. A human reviews this exact diff and merges it; that reviewed"
            echo "merge is the approval boundary that promotes the artifact to desired"
            echo "state. CI publication alone is not CD."
          } > /tmp/pr-body.md

          gh pr create --repo "${REPO}" --base env/huy --head "${DEPLOY_BRANCH}" \
            --title "deploy: propose devops-app image for ${SOURCE_SHA:0:12}" \
            --body-file /tmp/pr-body.md
```

## Decisions worth flagging

1. **kubectl's YAML output is lossy; only its image scalar is taken.** `kubectl
   set image --local ... -o yaml` re-serialises the document: it reorders fields,
   drops the manifest's comments, and adds `strategy: {}` and `status: {}`.
   Replacing the whole file with that output would violate "change ONLY the app
   container image" (it would strip the deliberate `minDomains` comment). So
   kubectl remains the authority for the target image string, and a checked
   Python splice replaces exactly that one scalar in the tracked file, preserving
   every other byte. This is a reviewed, single-file replacement — never a regex
   substitution across manifests.
2. **Process-scoped credential, not argv, not disk.** The requirements allow a
   process-scoped `http.https://github.com/.extraheader`. It is supplied through
   git's `GIT_CONFIG_COUNT`/`GIT_CONFIG_KEY_0`/`GIT_CONFIG_VALUE_0` environment
   config, which keeps the token out of `git ... -c` argv and out of any config
   file, and is not carried to the next step.
3. **`env: GITHUB_TOKEN:` is legal.** The `GITHUB_` reserved-prefix rule applies
   to secret/repository/org/environment *variable* names, not to step `env`
   maps; setting `GITHUB_TOKEN` in a step `env` block is standard (confirmed by
   many public workflows, e.g. PaddlePaddle/Paddle, FreeTubeApp/FreeTube). `gh`
   prefers `GH_TOKEN`; both are set for completeness.

## Local verification performed (2026-10-12)

| Check | Result |
|---|---|
| `yaml.safe_load` on `.github/workflows/ci.yml` | loads; jobs = `lint, compile, test, build, smoke, publish, propose-deploy` |
| `list(d['jobs'])` + `needs`/`if`/`permissions` of `propose-deploy` | `needs=publish`, `if=push && refs/heads/personal/huy`, `permissions` **absent** (top-level `{}` unchanged) |
| `bash -n` on every `run:` block (all jobs) | 21/21 parse clean |
| `grep -n PROJECT_GITHUB_TOKEN .github/workflows/ci.yml` | only `publish` and `propose-deploy` (plus explanatory comments) |
| `on:` block re-parse | unchanged: only `personal/huy` push/PR, same paths |
| kubectl pin | `.sha256` for v1.37.1 fetched; workstation client hashes to the same value |

### Simulated job logic against a real (local bare) remote

The `prepare` and `push` run blocks were executed against a local bare Git repo
seeded with the actual `env/huy` tree, standing in for `origin`:

| Scenario | Expected | Observed |
|---|---|---|
| Fresh SHA, `env/huy` differs | 1-line image commit, push to `deploy/huy-<sha>` | commit `+1/-1` in `k8s/app/deployment.yaml` only; comments preserved; pushed |
| Re-run same SHA, same digest | no commit, reuse PR | `changed=false propose=true`; remote branch tip unchanged |
| Re-run same SHA, new digest | fast-forward commit on top | `changed=true`; `b2a5f8e..cf833ac` (fast-forward) |
| Fresh SHA, `env/huy` already has digest | no commit, no branch, no PR | `changed=false propose=false`; no branch created |
| Existing branch, `env/huy` advanced | STOP, non-zero, remote untouched | exit 1 with both tips; branch unchanged |
| Validation guards (empty/`latest` digest, wrong repo, wrong image) | each fails closed | exit 1 with a specific error |

## Unrun checks (Sisyphus / real run)

- **Real `propose-deploy` run on a green push to `personal/huy`** — unpushed here;
  observe the opened PR and its exact diff.
- **Real PR reuse on a re-run** of the same source SHA (no duplicate).
- **Real divergence STOP** against an advanced `env/huy`.
- **PAT scope:** `repo` (present, per `docs/agents/github-access.md`) covers branch
  push and PR creation; commit touches only `k8s/**`, so the credential's missing
  `workflow` scope is not needed for the deploy branch (it is still needed for
  Sisyphus to push the changed `ci.yml` on the source branch).
- **Owner confirmation** that repository PR settings allow this account to create
  PRs targeting `env/huy`.
