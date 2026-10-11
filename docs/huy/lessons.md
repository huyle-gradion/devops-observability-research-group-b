# Lessons learned (personal project — huyle-gradion)

Owner: huy. Created: 2026-10-11 (Stage 01, substep 01C).

Each lesson uses **expected / observed / explanation / evidence**. Only observations
that actually happened are recorded. Where the cause is not yet known, the explanation
is `pending` rather than guessed. Agent-written text does not certify understanding;
the owner's own reproduction and explanation is the learning gate.

## Lesson 1 — Docker engine not reachable

| Field | Content |
|---|---|
| Expected | Docker engine is running and the worker can exercise it with the authorized operations. |
| Observed | At check time the docker **client** (29.7.2) is present but the **engine** is inactive. |
| Explanation | pending owner-run of the authorized fix (`sudo systemctl start docker`). |
| Evidence | `docker --version` reports 29.7.2; `docker info` returned client-only output with the engine not responding. Recorded 2026-10-11. |

## Lesson 2 — `sudo -n` requires a password

| Field | Content |
|---|---|
| Expected | Non-interactive sudo (`sudo -n true`) succeeds, allowing scripted privileged steps. |
| Observed | `sudo -n true` fails: "sudo: a password is required" (exit 1). |
| Explanation | pending owner-run; this is a human sudo gate, not a worker-solvable permission. |
| Evidence | `sudo -n true` exit code 1, message recorded 2026-10-11. No socket/permission changes were attempted. |

## Lesson 3 — Project credential lacks `workflow` scope

| Field | Content |
|---|---|
| Expected | The project credential can push new GitHub Actions workflow files so Stage 04 CI can be published. |
| Observed | The project credential does not advertise the `workflow` scope. |
| Explanation | GitHub requires the `workflow` scope to add or update Actions workflow files (except identical content/path on another branch). Known, owner-owned; not solvable by the worker without broadening credentials, which is forbidden. |
| Evidence | `docs/agents/github-access.md` records the observed scopes (`repo`, `write:packages`, broad admin/delete; no `workflow`). Recorded 2026-10-11. No fallback token was minted. |

## Actual errors

Raw failures encountered while doing this work, kept separate from the structured
lessons above. Add entries as they occur; do not delete an entry because it was fixed.

| Date | Command / action | Error / result | Resolution |
|---|---|---|---|
| 2026-10-11 | `sudo -n true` | exit 1, "sudo: a password is required" | pending owner-run |
| 2026-10-11 | `docker info` | engine inactive / not responding | pending owner-run of `sudo systemctl start docker` |
