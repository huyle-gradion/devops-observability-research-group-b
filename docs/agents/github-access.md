# Project GitHub credential policy

## Single credential authority

Owner instruction: all authenticated GitHub access uses the project account access token stored in `.env`. Planning securely observed key `GITHUB_ACCESS_TOKEN` in `/home/thearesia/Projects/assignment-3/.env`; authenticated GET /user returned `huyle-gradion`. Token value was not printed, copied or persisted. Earlier cached authentication as huyhuyhuy8s is NOT authority for this project.

| Context | Credential / identity |
|---|---|
| Local gh/API/Git HTTPS | .env GITHUB_ACCESS_TOKEN, mapped in child process to GH_TOKEN; expected account huyle-gradion |
| Local authenticated GHCR | Same .env token, stdin to docker login ghcr.io, username huyle-gradion |
| Trusted source-push CI | Actions secret PROJECT_GITHUB_TOKEN populated from the same .env token; expected project account |
| CI API/PR/attestation API | GH_TOKEN or action github-token/token input explicitly set from PROJECT_GITHUB_TOKEN |
| CI registry login | password PROJECT_GITHUB_TOKEN, username huyle-gradion, not github.actor |
| Public PR checks/public Argo reads/anonymous package proof | Anonymous reads, no project secret and no authenticated fallback |
| Attestation signature | Short-lived Actions OIDC identity with id-token: write, NOT another GitHub access token |

Do not use cached gh login, personal-account tokens, SSH keys, literal token URLs or built-in GITHUB_TOKEN for project GitHub/GHCR calls. Public anonymous reads and GitHub's internal runner transport are not alternate project-account identities. OIDC proves the signing workflow, whereas the PAT authorizes API/registry access; a PAT cannot replace OIDC signing.

## Secure local loading contract

Stage 01 will implement the bounded helper `scripts/project-github.py`; it does not exist during planning. Run authenticated CLI commands through it after verification. The helper reads the owner-selected absolute .env path, requires exactly one nonempty GITHUB_ACCESS_TOKEN, parses data without shell execution/interpolation, and never dumps values. The private locator is PROJECT_GITHUB_ENV_FILE, defaulting to the observed DOC_ROOT/.env; .env is not automatically copied into SOURCE_ROOT/worktrees.

Child environment overrides GH_TOKEN from that value, removes inherited GITHUB_TOKEN/alternate GitHub credentials and debug tracing, and uses no cached-account fallback. `gh api user --jq .login` is the identity check; fail closed unless huyle-gradion. For Git, use HTTPS plus process-scoped credential.helper reset and gh auth git-credential, not global persistence or tokens in remote URLs. Restrict remote host/repository to the authorized project. Set process-scoped Git configuration in the child environment too, so Git invoked inside gh uses the same credential boundary; never inherit a personal helper.

For Docker, helper subcommand `registry-login --config <new-private-config-dir>` sends the same token through stdin to docker login ghcr.io as huyle-gradion. Registry login does not require Docker daemon/root access; create this config as the user so gh attestation can read it. No token in argv/global login. Later sudo docker explicitly uses --config with that directory. Remove only the owned temporary auth config after the session. Never echo credentials, set -x, source executable .env content or emit credential-helper output into public logs. Record key/account/scope names only.

Public lookup can be anonymous; authenticated control operations go through this contract. If the helper is not yet implemented/verified, use an equivalent owner-approved in-memory loader; do not silently fall back to bare gh or persistent personal credentials. Missing/empty/invalid/wrong-account token blocks authenticated operations.

## CI secret and trusted-code boundary

Owner/Sisyphus with explicit authorization sets PROJECT_GITHUB_TOKEN from the in-memory .env value over stdin. Stage 01 helper subcommand `sync-ci-secret --repo <authorized-repo>` feeds that same value to gh secret set without a token argument; no secret-management operation is permitted without explicit authorization. GitHub secret names cannot start GITHUB_, so local key and Actions secret differ. CI cannot read workstation .env. Never commit/upload-artifact/Docker-copy the file. Rotation updates .env and that secret through the same controlled path.

Project secret is provided only to explicitly trusted source-push publication/propose-deploy steps, never to pull_request code, check job, model prompts or arbitrary package install scripts. PR checkout is anonymous public Git fetch; actions/setup-go and other helpers must explicitly disable default GitHub-token auth where supported. If an action cannot run without an implicit alternate token, choose a verified anonymous equivalent or report the blocker. Trusted checkout can use PROJECT_GITHUB_TOKEN with persist-credentials false; erase checkout credentials before build/test commands. Docker build-push action must set github-token explicitly rather than accepting its default. Minimal built-in workflow permissions are {} except publish id-token: write for signing. These permissions DO NOT narrow PAT scopes.

Pinned actions and reviewed trusted source reduce, but do not eliminate, the risk of passing a broad PAT to a runner. No workflow change from unreviewed PR code may acquire it. CI PR automation uses only project-owned named deploy branches and env/huy; a PAT-created PR CAN trigger other workflows, so keep explicit source-branch/path filters and no env/huy source-build recursion.

## Observed permissions and blocker

GET /user advertised repo, write:packages and broad admin/delete scopes. No workflow scope was advertised. [GitHub scope docs](https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/scopes-for-oauth-apps) require workflow to add/update Actions workflow files, except identical existing content/path on another branch. New CI publication is therefore blocked until the owner grants workflow to the SAME authorized project credential or replaces its value in .env with an appropriately scoped credential for the SAME project account and synchronizes the Actions secret. Agents do not mint a fallback token or switch identity.

The credential can reach resources beyond this assignment. Prohibit organization/enterprise administration, repo/package deletion, production access and unrelated repositories regardless of token capability. Record the risk; recommend owner narrowing unnecessary scopes on the same project account. Do not silently rotate/change it. Workflow permissions and repository path checks are not a cryptographic restriction on this broad PAT.

For this public project, repo/public_repo access, workflow for workflow edits, and write:packages for GHCR are relevant. [GHCR docs](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry) require PAT classic for PAT authentication; do not assume a fine-grained replacement can publish packages. [Attestation REST docs](https://docs.github.com/en/rest/repos/attestations) require repo write/read access; actions/attest's current [github-token input](https://github.com/actions/attest/blob/main/action.yml) accepts the explicit project PAT. Set create-storage-record false unless its additional API permission is verified for this credential; signed provenance/registry push remain required.

## Agent access boundaries and incident handling

Prometheus may inspect key names and authenticated identity/scopes to plan, without publishing/mutating GitHub. Sisyphus performs authorized branch/PR/secret setup and independent verification. Workers receive only the named credential locator and bounded operation permission, never the value in prompts. Reader/model tools must not open raw .env into their transcript.

Keep .env/.env.* out of git and Docker context; exception only a deliberately value-free example if ever requested. Protect file permissions locally. Never save secret-bearing CLI configs, auth headers, environment dumps or credential helper output in captures. If disclosure is detected, stop, notify owner, revoke/rotate through owner authority, and repair exposure before resuming. A permission failure is reported, not solved by broadening access or changing accounts.
