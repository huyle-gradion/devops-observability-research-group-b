# Agent instructions

## Personal project

- Personal specification: `.scratch/assignment-3-personal/spec.md`
- Personal plan index (stages 01–14): `.omo/plans/devops-personal/README.md`
- Execution contract: `.omo/plans/devops-personal/00-execution-contract.md`

## Local tracker conventions

Issues and specs live as local markdown files under `.scratch/`. See `docs/agents/issue-tracker.md`.

Triage roles (`needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`): see `docs/agents/triage-labels.md`.

Domain docs: `GLOSSARY.md` (single-context layout). See `docs/agents/domain.md`.

## Multi-agent workflow and GitHub access

- Roles and dispatch: `docs/agents/opencode-workflow.md`, `docs/agents/worker-automation.md`
- Before any authenticated GitHub or GHCR operation, read `docs/agents/github-access.md`. Use only the project account's `.env` `GITHUB_ACCESS_TOKEN`. Never print its value or fall back to cached credentials.
