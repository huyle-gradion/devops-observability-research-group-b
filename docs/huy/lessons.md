# Lessons learned (personal project — huyle-gradion)

Owner: huy. Created: 2026-10-11 (Stage 01, substep 01C).

Each lesson uses **expected / observed / explanation / evidence**. Only observations
that actually happened are recorded. Where the cause is not yet known, the explanation
is `pending` rather than guessed. Agent-written text does not certify understanding;
the owner's own reproduction and explanation is the learning gate.

## Lesson 1 — Docker engine not reachable → now active

| Field | Content |
|---|---|
| Expected | Docker engine is running and the worker can exercise it with the authorized operations. |
| Observed | Initially at check time the docker **client** (29.7.2) was present but the **engine** was inactive. After the owner ran the authorized fix, `systemctl is-active docker` = `active` and the engine responds: server Engine 29.7.2, API 1.55; `sudo docker run --rm hello-world` ran successfully (amd64). |
| Explanation | The authorized fix (`sudo systemctl start docker`), run by the owner, brought the engine up. Engine and container run are now verified. |
| Evidence | `sudo docker version` (server Engine 29.7.2, API 1.55, containerd v2.3.5, runc 1.5.1, linux/amd64) and `sudo docker run --rm hello-world` (pulled `library/hello-world`, digest `sha256:5e23090353324d887c48ad5e5c56d294eab81588df9605b07d1afe895f9cc8f8`, printed `Hello from Docker!`). Recorded 2026-10-11. |

## Lesson 2 — `sudo -n` requires a password

| Field | Content |
|---|---|
| Expected | Non-interactive sudo (`sudo -n true`) succeeds, allowing scripted privileged steps. |
| Observed | `sudo -n true` fails: "sudo: a password is required" (exit 1). This still holds after the Docker fix: the engine is now active, but the user is not in the `docker` group, so privileged docker access continues to require an explicit password-bearing `sudo docker` (or the service must be started by the owner). |
| Explanation | This is a human sudo gate, not a worker-solvable permission. The owner ran the authorized engine start; the gate for privileged commands remains. |
| Evidence | `sudo -n true` exit code 1, message recorded 2026-10-11. `sudo docker version` succeeded only because the owner supplied a password; a non-sudo `docker run` returned a socket permission error (see Actual errors). No socket/permission changes were attempted. |

## Lesson 3 — Project credential lacks `workflow` scope

| Field | Content |
|---|---|
| Expected | The project credential can push new GitHub Actions workflow files so Stage 04 CI can be published. |
| Observed | The project credential does not advertise the `workflow` scope. |
| Explanation | GitHub requires the `workflow` scope to add or update Actions workflow files (except identical content/path on another branch). Known, owner-owned; not solvable by the worker without broadening credentials, which is forbidden. |
| Evidence | `docs/agents/github-access.md` records the observed scopes (`repo`, `write:packages`, broad admin/delete; no `workflow`). Recorded 2026-10-11. No fallback token was minted. |

## Lesson 4 — the local Docker image store cannot carry SBOM/provenance

| Field | Content |
|---|---|
| Expected | A `docker/build-push-action` step with `sbom: true`, `provenance: mode=max`, `push: false` and `load: true` builds the image locally and emits its SBOM and provenance. |
| Observed | Not run in CI yet, but BuildKit rejects that combination. Attaching attestations needs an image exporter that supports them; the `docker` exporter behind `--load`/`type=docker` does not, so buildx fails with `docker exporter does not currently support exporting manifest lists` unless the daemon exposes the OCI importer. |
| Explanation | The buildx source (`build/opt.go`) only converts a `docker` export to an OCI export when the daemon advertises `OCIImporter`; otherwise any attestation request is an error. The default GitHub runner image store does not guarantee that. The six-job build therefore uses `outputs: type=oci,dest=...` for the no-push build (which does support attestations) and lets the smoke job build its own runnable copy with plain `docker build`. The signed, published attestations stay the ones the `publish` job emits. |
| Evidence | buildx validation read in `build/opt.go` (`supportAttestations` / `docker exporter does not currently support exporting manifest lists`); see also docker/buildx issues #1714 and docker/build-push-action #1135. Recorded 2026-10-12. |

## Actual errors

Raw failures encountered while doing this work, kept separate from the structured
lessons above. Add entries as they occur; do not delete an entry because it was fixed.

| Date | Command / action | Error / result | Resolution |
|---|---|---|---|
| 2026-10-11 | `sudo -n true` | exit 1, "sudo: a password is required" | remains; use explicit password-bearing sudo for privileged commands |
| 2026-10-11 | `docker info` | engine inactive / not responding | resolved: owner ran `sudo systemctl start docker`; engine now active (29.7.2, API 1.55) |
| 2026-10-11 | `docker run` (no sudo) | "permission denied while trying to connect to the docker API at unix:///var/run/docker.sock" — `/var/run/docker.sock` is `root:docker` mode 660 and the user is not in the `docker` group | use explicit `sudo docker` for this assignment; no socket chmod and no persistent group change |
| 2026-10-12 | `docker build -f app/Dockerfile -t devops-app-huy:ci-local .` | same "permission denied" on the docker API socket; `sudo -n true` still exits 1 | local image build/run checks could not be run on the workstation; they remain for the GitHub runner and the owner's password-bearing reproduction |
