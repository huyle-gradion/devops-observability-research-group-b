#!/usr/bin/env python3
"""Bounded credential helper for the Assignment 3 personal project.

This is the single local entry point for authenticated GitHub/GHCR operations.
It is deliberately small and fail-closed:

* It reads an owner-selected absolute ``.env`` file (``PROJECT_GITHUB_ENV_FILE``)
  and requires exactly one nonempty ``GITHUB_ACCESS_TOKEN``. The file is parsed
  as ``KEY=VALUE`` text only: it is never sourced, executed or interpolated.
* Every child process receives a rebuilt environment in which ``GH_TOKEN`` is
  the only GitHub credential. Inherited ``GITHUB_TOKEN``/alternate credentials
  and Git/GH debug tracing are removed, so there is no cached-account fallback.
* ``verify`` runs ``gh api user --jq .login`` and fails closed unless the login
  is the authorized project account.
* ``git`` (and ``gh`` passthrough, so Git invoked inside gh) uses process-scoped
  ``credential.helper`` reset plus ``gh auth git-credential``. Nothing is
  persisted globally and no token is placed in a remote URL.
* ``registry-login`` and ``sync-ci-secret`` feed the token over stdin only; the
  token never appears in argv.

The token value is never printed, echoed or written to disk by this helper.

Authoritative policy: docs/agents/github-access.md.
"""

from __future__ import annotations

import argparse
import os
import re
import subprocess
import sys
from pathlib import Path

DEFAULT_ENV_FILE = "/home/thearesia/Projects/assignment-3/.env"
ENV_FILE_VAR = "PROJECT_GITHUB_ENV_FILE"
TOKEN_KEY = "GITHUB_ACCESS_TOKEN"
EXPECTED_LOGIN = "huyle-gradion"
REGISTRY = "ghcr.io"
DEFAULT_SECRET_NAME = "PROJECT_GITHUB_TOKEN"
GIT_CREDENTIAL_HELPER = "!gh auth git-credential"

# Inherited credentials and debug tracing that must never reach a child process.
STRIPPED_ENV = (
    "GITHUB_TOKEN",
    "GH_TOKEN",
    "GITHUB_ENTERPRISE_TOKEN",
    "GH_ENTERPRISE_TOKEN",
    "GITHUB_PAT",
    "GIT_ASKPASS",
    "SSH_ASKPASS",
    "GH_DEBUG",
    "GIT_TRACE",
    "GIT_TRACE_PACKET",
    "GIT_TRACE_CURL",
    "GIT_TRACE_CURL_NO_DATA",
    "GIT_CURL_VERBOSE",
    "GH_HOST",
)

# owner/name only: rejects URLs, hosts, dashes-as-flags and whitespace.
_REPO_RE = re.compile(r"^[A-Za-z0-9][A-Za-z0-9_.-]*/[A-Za-z0-9][A-Za-z0-9_.-]*$")


class HelperError(Exception):
    """Fail-closed condition. Messages never contain the token value."""


def parse_env_text(text: str) -> dict[str, list[str]]:
    """Parse ``KEY=VALUE`` lines as literal data.

    No shell sourcing, execution or interpolation happens: ``$VAR`` stays
    literal. Blank lines and ``#`` comments are skipped. Matching surrounding
    single/double quotes are stripped as literal text.
    """
    parsed: dict[str, list[str]] = {}
    for raw in text.splitlines():
        line = raw.strip()
        if not line or line.startswith("#"):
            continue
        key, sep, value = line.partition("=")
        if not sep:
            continue
        key = key.strip()
        if not key:
            continue
        value = value.strip()
        if len(value) >= 2 and value[0] == value[-1] and value[0] in ("'", '"'):
            value = value[1:-1]
        parsed.setdefault(key, []).append(value)
    return parsed


def load_token(env_file: str) -> str:
    """Return the single nonempty GITHUB_ACCESS_TOKEN from ``env_file``."""
    path = Path(env_file)
    if not path.is_absolute():
        raise HelperError(f"{ENV_FILE_VAR} must be an absolute path, got {env_file!r}")
    try:
        text = path.read_text(encoding="utf-8")
    except OSError as exc:
        raise HelperError(f"cannot read env file {path}: {exc.strerror or exc}") from None
    values = parse_env_text(text).get(TOKEN_KEY, [])
    nonempty = [value for value in values if value]
    if len(values) != 1 or not nonempty:
        raise HelperError(
            f"expected exactly one nonempty {TOKEN_KEY} in {path}; "
            f"found {len(values)} occurrence(s)"
        )
    return nonempty[0]


def base_child_env() -> dict[str, str]:
    """Rebuilt environment with inherited credentials and tracing removed."""
    env = dict(os.environ)
    for name in STRIPPED_ENV:
        env.pop(name, None)
    # Drop any inherited process-scoped git config so ours is the only one.
    for name in list(env):
        if (
            name == "GIT_CONFIG_COUNT"
            or name.startswith("GIT_CONFIG_KEY_")
            or name.startswith("GIT_CONFIG_VALUE_")
        ):
            env.pop(name, None)
    env["GIT_TERMINAL_PROMPT"] = "0"
    return env


def with_token(env: dict[str, str], token: str) -> dict[str, str]:
    env = dict(env)
    env["GH_TOKEN"] = token
    return env


def with_git_credential_boundary(env: dict[str, str]) -> dict[str, str]:
    """Process-scoped git config: reset helpers, then use gh's credential helper."""
    env = dict(env)
    pairs = [
        ("credential.helper", ""),
        ("credential.helper", GIT_CREDENTIAL_HELPER),
        ("url.https://github.com/.insteadOf", "git@github.com:"),
        ("url.https://github.com/.insteadOf", "ssh://git@github.com/"),
    ]
    env["GIT_CONFIG_COUNT"] = str(len(pairs))
    for index, (key, value) in enumerate(pairs):
        env[f"GIT_CONFIG_KEY_{index}"] = key
        env[f"GIT_CONFIG_VALUE_{index}"] = value
    return env


def _run(cmd: list[str], env: dict[str, str], **kwargs) -> subprocess.CompletedProcess:
    try:
        return subprocess.run(cmd, env=env, **kwargs)
    except FileNotFoundError:
        raise HelperError(f"required executable not found on PATH: {cmd[0]}") from None


def run_verify(token: str) -> int:
    env = with_token(base_child_env(), token)
    completed = _run(
        ["gh", "api", "user", "--jq", ".login"],
        env,
        capture_output=True,
        text=True,
    )
    if completed.returncode != 0:
        if completed.stderr:
            sys.stderr.write(completed.stderr)
        raise HelperError(f"gh api user failed with exit code {completed.returncode}")
    login = completed.stdout.strip()
    if login != EXPECTED_LOGIN:
        raise HelperError(
            f"identity check failed: resolved login {login!r} is not {EXPECTED_LOGIN!r}"
        )
    print(login)
    return 0


def run_gh(token: str, args: list[str]) -> int:
    if not args:
        raise HelperError("gh subcommand requires arguments")
    # Include the git boundary so Git invoked inside gh shares the credential.
    env = with_git_credential_boundary(with_token(base_child_env(), token))
    return _run(["gh", *args], env).returncode


def run_git(token: str, args: list[str]) -> int:
    if not args:
        raise HelperError("git subcommand requires arguments")
    env = with_git_credential_boundary(with_token(base_child_env(), token))
    return _run(["git", *args], env).returncode


def run_registry_login(token: str, config_dir: str) -> int:
    if not config_dir:
        raise HelperError("registry-login requires --config DIR")
    cfg = Path(config_dir)
    if not cfg.is_absolute():
        raise HelperError("--config must be an absolute path")
    if cfg.exists():
        raise HelperError(f"--config directory must be new but already exists: {cfg}")
    cfg.mkdir(mode=0o700, parents=True)
    os.chmod(cfg, 0o700)
    # Token goes over stdin only; docker does not need GH_TOKEN.
    cmd = [
        "docker",
        "login",
        REGISTRY,
        "--username",
        EXPECTED_LOGIN,
        "--password-stdin",
        "--config",
        str(cfg),
    ]
    return _run(cmd, base_child_env(), input=token + "\n", text=True).returncode


def run_sync_ci_secret(token: str, repo: str, secret_name: str) -> int:
    if not _REPO_RE.match(repo or ""):
        raise HelperError(f"--repo must be owner/name, got {repo!r}")
    env = with_token(base_child_env(), token)
    cmd = ["gh", "secret", "set", secret_name, "--repo", repo]
    # Value over stdin only; never a token argument.
    return _run(cmd, env, input=token + "\n", text=True).returncode


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="project-github.py",
        description="Bounded project GitHub/GHCR credential helper.",
    )
    sub = parser.add_subparsers(dest="command", required=True)

    sub.add_parser("verify", help=f"check gh identity is {EXPECTED_LOGIN}")

    p_gh = sub.add_parser("gh", help="run gh with the credential boundary")
    p_gh.add_argument("args", nargs=argparse.REMAINDER)

    p_git = sub.add_parser("git", help="run git with the credential boundary")
    p_git.add_argument("args", nargs=argparse.REMAINDER)

    p_reg = sub.add_parser("registry-login", help="docker login ghcr.io via stdin")
    p_reg.add_argument("--config", required=True, help="new private docker config dir")

    p_sync = sub.add_parser("sync-ci-secret", help="gh secret set via stdin")
    p_sync.add_argument("--repo", required=True, help="authorized owner/name")
    p_sync.add_argument("--secret-name", default=DEFAULT_SECRET_NAME)
    return parser


def main(argv: list[str] | None = None) -> int:
    parser = build_parser()
    args = parser.parse_args(argv)
    env_file = os.environ.get(ENV_FILE_VAR, DEFAULT_ENV_FILE)
    try:
        token = load_token(env_file)
        if args.command == "verify":
            return run_verify(token)
        if args.command == "gh":
            return run_gh(token, args.args)
        if args.command == "git":
            return run_git(token, args.args)
        if args.command == "registry-login":
            return run_registry_login(token, args.config)
        if args.command == "sync-ci-secret":
            return run_sync_ci_secret(token, args.repo, args.secret_name)
    except HelperError as exc:
        sys.stderr.write(f"project-github: {exc}\n")
        return 2
    parser.error(f"unhandled command: {args.command}")
    return 2


if __name__ == "__main__":
    raise SystemExit(main())
