# Stage 01 teachback — workspace, prerequisites, credentials (huyle-gradion)

Audience: the owner. Purpose: after Stage 01 was implemented, learn **what happened,
what broke, what to carry forward, and how to rebuild it fast**. Written 2026-10-11.

Scoring reminder from the assignment (§6): proof of work 40%, research 40%,
presentation 20%. "A real pipeline run you captured yourself is worth at least as
much as a well-designed slide deck." Stage 01 exists to make that possible safely.

---

## 1. What Stage 01 actually delivered

| Item | Where |
|---|---|
| Public personal repo | `huyle-gradion/devops-observability-research-group-b` |
| Credential helper + tests | `scripts/project-github.py`, `scripts/tests/test_project_github.py` |
| Context import (plan/spec/agents/assignment) | repo root, `.omo/`, `.scratch/`, `docs/agents/` |
| Evidence ledger | `docs/huy/requirements.md` (25 rows, all `not-run`) |
| Versions / lessons / findings | `docs/huy/versions.md`, `lessons.md`, `findings.md` |
| Secret-exclusion boundary | `.gitignore` |
| Setup capture | `captures/huy/step-1/setup.md` |
| Integration refs | `personal/huy` = `stage/huy/01-workspace` = `c62d6ea`; `main` @ `cfc71a2` |

## 2. What went wrong (honest friction log)

1. **Docker appeared "broken" but was only not started.** The engine was `inactive`.
   Worse, the first run attempt omitted `sudo` and failed with
   `permission denied … /var/run/docker.sock`. The socket is `root:docker` mode 660 and
   the user is not in the `docker` group, so unprivileged `docker` cannot talk to the
   daemon. **Fix:** `sudo systemctl start docker` (already installed) then explicit
   `sudo docker …`. We deliberately did **not** `chmod` the socket or join the group.
2. **`.gitignore` accidentally hid evidence.** A blanket `*.log` rule silently ignored
   `.log` captures under `captures/`. The assignment (§4, §6) requires **every capture
   committed next to the write-up**. **Fix:** `!captures/**` re-includes captures after
   the transient-log rule. Caught by review, not by the worker.
3. **Model/transport drift.** The dispatch header declared one model; the harness ran a
   different one. Lesson: verify the model at dispatch, and record the actual value in
   the handoff instead of assuming the header was honoured.
4. **Wrong cached GitHub identity was present.** `gh` was logged in as a personal
   account (`huyhuyhuy8s`). Using it would have been the wrong identity. The helper
   forces the project account (`huyle-gradion`) via the `.env` token and fails closed
   otherwise.
5. **Destination repo did not exist.** The plan anticipated this: an empty public repo
   receives exactly one bootstrap seed, then branches are created from it. Do not
   invent repository existence — verify it.
6. **Commit author identity.** The global git identity was the cached personal account.
   We set the **repo-local** identity to the project account so commits are not
   misattributed.

## 3. Lessons / assignment-required knowledge for this stage

- **Public repo + incremental history (§3.2, §4, §6).** Git history must show
  incremental work, "not one final paste". Every capture is committed in the repo.
  This is *why* Stage 01 sets up commits, `.gitignore` and a `captures/` tree.
- **Credential separation.** A GitHub token is a secret. In a public repo, anything
  committed is permanent — rotating a leaked token does **not** remove it from history.
  So: token only in a private `.env` (never committed), read in memory, `GH_TOKEN`
  overridden per child process, fail closed unless the identity is the project account.
  A committed token is harvested by scanners within minutes.
- **Scope matters.** The token needs `repo` (+ `write:packages` for GHCR) and, for
  editing Actions workflow files, `workflow`. Ours lacks `workflow`, which **blocks
  Stage 04** until the owner fixes it on the same account. Never mint a fallback token
  or switch identity to work around a missing scope.
- **Docker privilege is root-equivalent.** The daemon runs as root; socket access lets
  a caller mount the host and run privileged containers. Granting that permanently
  (group membership) or via a world-writable socket is a security regression. A scoped
  `sudo docker` is the safe, reversible choice.
- **amd64 vs arm64 gotcha (§3.2 step 1).** The assignment assumes Apple/M1 (arm64)
  laptops while GitHub Actions builds `linux/amd64`; an image that builds in CI may fail
  locally. **Our deviation, recorded honestly:** this workstation is native
  **amd64**, so that mismatch is reduced here — but we must never claim an arm64 run
  we did not perform.
- **Free tiers only (§7).** Own machines, public GHCR, no paid cloud / managed K8s.

## 4. How to run and inspect it

```sh
cd /home/thearesia/Projects/devops-worktrees/01-workspace

# 1. Prove the credential helper works (prints only the login; never the token)
python3 scripts/project-github.py verify                 # -> huyle-gradion

# 2. Run the offline tests yourself (synthetic tokens, fake gh/git/docker)
python3 -m unittest discover -s scripts/tests -p 'test_*.py'

# 3. See the failure that made sudo necessary
docker run --rm hello-world                              # -> permission denied on the socket
sudo docker run --rm hello-world                         # -> Hello from Docker!

# 4. Inspect repo state and history
git log --oneline --decorate -n 5
git status --short
python3 scripts/project-github.py git ls-remote --heads origin
```

**What you can do with it:** any authenticated GitHub/GHCR operation for this project
goes through the helper (`verify`, `gh …`, `git …`, `registry-login`, `sync-ci-secret`),
so the token never lands in shell history or argv; the repo is the source of truth for
all later stages; `docs/huy/requirements.md` is the checklist you update as outcomes
get real evidence.

## 5. Fastest rebuild runbook (adapted from the lessons above)

Do it in this order; each step is the lesson from §2/§3 applied.

1. **Prereqs:** confirm `git`, `go`, `gh`, `python3`, `docker` exist; start Docker the
   authorized way — `sudo systemctl start docker`; verify `sudo docker run --rm
   hello-world`. Do **not** chmod the socket or change group membership.
2. **Verify identity before anything remote:** load the token from the private `.env`
   and check `gh api user --jq .login` == the project account; refuse to continue
   otherwise. Never use a cached login.
3. **Create/confirm the public repo** (personal account). If it does not exist, create
   it public with a single README seed; remember the seed SHA.
4. **Clone as a clean sibling** and create branches from the seed SHA: the integration
   branch (`personal/huy`) and a stage branch (`stage/huy/01-workspace`); add the stage
   worktree under a `devops-worktrees/` parent.
5. **Set the repo-local commit identity** to the project account (avoid the cached
   personal identity).
6. **Add the credential helper + tests**; run the offline tests; then run the live
   read-only smoke (`verify`).
7. **Import context files** (plan pack, spec/tickets, docs/agents, assignment doc,
   GLOSSARY) preserving relative paths; add the AGENTS.md pointer and a short README.
8. **Create the ledgers** (`requirements.md` with one row per required outcome + bonuses
   + research + team/presentation, all `not-run`), `versions.md`, `lessons.md`,
   `findings.md`, and `.gitignore`.
9. **Guard the boundary:** `.gitignore` must ignore `.env`/`.env.*` and kubeconfigs but
   must **not** ignore `captures/`. Verify with `git check-ignore -q`.
10. **Commit explicitly** (list paths; no broad `git add`), integrate the stage branch
    into the integration branch (fast-forward), and push. Record the exact SHAs.

**Traps to avoid:** committing a token; letting `*.log`/broad patterns hide captures;
checking out a branch that is already checked out elsewhere; broad `git add`; claiming
an arm64 run you did not do; using the cached GitHub account.
