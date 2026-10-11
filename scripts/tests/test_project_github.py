#!/usr/bin/env python3
"""Offline unit tests for scripts/project-github.py.

All credentials and child commands here are synthetic. The helper is exercised
as a real subprocess with fake ``gh``/``git``/``docker`` executables first on
PATH, so no network, no real token and no real GitHub/Docker access is used.

Covered: missing/empty/duplicate token rejection, wrong-account and
correct-account identity checks, child environment credential stripping and
git credential-helper boundary injection, registry-login/sync-ci-secret stdin
delivery, config isolation, and secret non-disclosure in every output channel.
"""

import os
import stat
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

HELPER = Path(__file__).resolve().parents[1] / "project-github.py"

SYNTHETIC_TOKEN = "ghp_synthetic_NOT_A_REAL_TOKEN_0123456789abcdef"
EXPECTED_LOGIN = "huyle-gradion"
WRONG_LOGIN = "someone-else"
AUTHORIZED_REPO = "huyle-gradion/devops-observability-research-group-b"

FAKE_GH = """#!/bin/sh
{
  echo "argv=$*"
  echo "GH_TOKEN=${GH_TOKEN-<unset>}"
  echo "GITHUB_TOKEN=${GITHUB_TOKEN-<unset>}"
  echo "GH_DEBUG=${GH_DEBUG-<unset>}"
  echo "GIT_CONFIG_COUNT=${GIT_CONFIG_COUNT-<unset>}"
  env | grep '^GIT_CONFIG_KEY_' | sort
  env | grep '^GIT_CONFIG_VALUE_' | sort
} >> "$FAKE_GH_CAPTURE"
if [ "$1" = "api" ] && [ "$2" = "user" ]; then
  printf '%s' "${FAKE_GH_LOGIN-}"
  exit "${FAKE_GH_EXIT-0}"
fi
if [ "$1" = "secret" ] && [ "$2" = "set" ]; then
  cat > "$FAKE_GH_STDIN"
fi
exit 0
"""

FAKE_GIT = """#!/bin/sh
{
  echo "argv=$*"
  echo "GH_TOKEN=${GH_TOKEN-<unset>}"
  echo "GIT_CONFIG_COUNT=${GIT_CONFIG_COUNT-<unset>}"
  env | grep '^GIT_CONFIG_KEY_' | sort
  env | grep '^GIT_CONFIG_VALUE_' | sort
} >> "$FAKE_GIT_CAPTURE"
exit 0
"""

FAKE_DOCKER = """#!/bin/sh
{
  echo "argv=$*"
  echo "GH_TOKEN=${GH_TOKEN-<unset>}"
} >> "$FAKE_DOCKER_CAPTURE"
cat > "$FAKE_DOCKER_STDIN"
exit 0
"""


class HelperTestCase(unittest.TestCase):
    def setUp(self) -> None:
        self._tmp = tempfile.TemporaryDirectory()
        self.root = Path(self._tmp.name)
        self.fakebin = self.root / "bin"
        self.fakebin.mkdir()
        self.envfile = self.root / "synthetic.env"
        self.gh_capture = self.root / "gh.capture"
        self.git_capture = self.root / "git.capture"
        self.docker_capture = self.root / "docker.capture"
        self.gh_stdin = self.root / "gh.stdin"
        self.docker_stdin = self.root / "docker.stdin"
        self.write_env(f"{'GITHUB_ACCESS_TOKEN'}={SYNTHETIC_TOKEN}\n")
        self.write_fake("gh", FAKE_GH)
        self.write_fake("git", FAKE_GIT)
        self.write_fake("docker", FAKE_DOCKER)

    def tearDown(self) -> None:
        self._tmp.cleanup()

    def write_env(self, text: str) -> None:
        self.envfile.write_text(text, encoding="utf-8")

    def write_fake(self, name: str, body: str) -> None:
        path = self.fakebin / name
        path.write_text(body, encoding="utf-8")
        path.chmod(0o755)

    def run_helper(self, *args, login=None, extra_env=None):
        env = dict(os.environ)
        env["PATH"] = str(self.fakebin) + os.pathsep + env.get("PATH", "")
        env["PROJECT_GITHUB_ENV_FILE"] = str(self.envfile)
        env["FAKE_GH_CAPTURE"] = str(self.gh_capture)
        env["FAKE_GIT_CAPTURE"] = str(self.git_capture)
        env["FAKE_DOCKER_CAPTURE"] = str(self.docker_capture)
        env["FAKE_GH_STDIN"] = str(self.gh_stdin)
        env["FAKE_DOCKER_STDIN"] = str(self.docker_stdin)
        if login is not None:
            env["FAKE_GH_LOGIN"] = login
        if extra_env:
            env.update(extra_env)
        return subprocess.run(
            [sys.executable, str(HELPER), *args],
            env=env,
            capture_output=True,
            text=True,
        )

    def read_capture(self, path: Path) -> str:
        return path.read_text(encoding="utf-8") if path.exists() else ""

    def assert_no_token(self, result, msg: str = "secret leaked") -> None:
        self.assertNotIn(SYNTHETIC_TOKEN, result.stdout, msg)
        self.assertNotIn(SYNTHETIC_TOKEN, result.stderr, msg)


class TokenLoadingTests(HelperTestCase):
    def test_accepts_single_nonempty_token(self) -> None:
        result = self.run_helper("verify", login=EXPECTED_LOGIN)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(result.stdout.strip(), EXPECTED_LOGIN)
        self.assert_no_token(result)

    def test_rejects_missing_token(self) -> None:
        self.write_env("SOME_OTHER_KEY=value\n")
        result = self.run_helper("verify", login=EXPECTED_LOGIN)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("GITHUB_ACCESS_TOKEN", result.stderr)
        self.assert_no_token(result)
        self.assertFalse(self.gh_capture.exists(), "gh must not run without a token")

    def test_rejects_empty_token(self) -> None:
        self.write_env("GITHUB_ACCESS_TOKEN=\n")
        result = self.run_helper("verify", login=EXPECTED_LOGIN)
        self.assertNotEqual(result.returncode, 0)
        self.assert_no_token(result)
        self.assertFalse(self.gh_capture.exists())

    def test_rejects_duplicate_token(self) -> None:
        self.write_env(f"GITHUB_ACCESS_TOKEN={SYNTHETIC_TOKEN}\nGITHUB_ACCESS_TOKEN=second\n")
        result = self.run_helper("verify", login=EXPECTED_LOGIN)
        self.assertNotEqual(result.returncode, 0)
        self.assert_no_token(result)

    def test_rejects_missing_env_file(self) -> None:
        self.envfile.unlink()
        result = self.run_helper("verify", login=EXPECTED_LOGIN)
        self.assertNotEqual(result.returncode, 0)
        self.assert_no_token(result)

    def test_rejects_relative_env_file(self) -> None:
        result = self.run_helper(
            "verify",
            login=EXPECTED_LOGIN,
            extra_env={"PROJECT_GITHUB_ENV_FILE": "relative.env"},
        )
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("absolute", result.stderr)

    def test_value_is_not_interpolated(self) -> None:
        # $HOME must stay literal; no shell expansion of the value.
        self.write_env("GITHUB_ACCESS_TOKEN=$HOME\n")
        result = self.run_helper("verify", login=EXPECTED_LOGIN)
        self.assertEqual(result.returncode, 0, result.stderr)
        capture = self.read_capture(self.gh_capture)
        self.assertIn("GH_TOKEN=$HOME", capture)
        self.assertNotIn("GH_TOKEN=/", capture)


class IdentityCheckTests(HelperTestCase):
    def test_accepts_authorized_account(self) -> None:
        result = self.run_helper("verify", login=EXPECTED_LOGIN)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(result.stdout.strip(), EXPECTED_LOGIN)
        self.assert_no_token(result)

    def test_rejects_wrong_account(self) -> None:
        result = self.run_helper("verify", login=WRONG_LOGIN)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("identity check failed", result.stderr)
        self.assertIn(WRONG_LOGIN, result.stderr)
        self.assert_no_token(result)

    def test_rejects_empty_login(self) -> None:
        result = self.run_helper("verify", login="")
        self.assertNotEqual(result.returncode, 0)
        self.assert_no_token(result)

    def test_gh_api_failure_is_reported(self) -> None:
        result = self.run_helper(
            "verify", login=EXPECTED_LOGIN, extra_env={"FAKE_GH_EXIT": "1"}
        )
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("gh api user failed", result.stderr)
        self.assert_no_token(result)


class GhPassthroughTests(HelperTestCase):
    def test_sets_token_and_strips_inherited_credentials(self) -> None:
        result = self.run_helper(
            "gh",
            "repo",
            "view",
            AUTHORIZED_REPO,
            extra_env={"GITHUB_TOKEN": "inherited-builtin", "GH_DEBUG": "1"},
        )
        self.assertEqual(result.returncode, 0, result.stderr)
        capture = self.read_capture(self.gh_capture)
        self.assertIn(f"GH_TOKEN={SYNTHETIC_TOKEN}", capture)
        self.assertIn("GITHUB_TOKEN=<unset>", capture)
        self.assertIn("GH_DEBUG=<unset>", capture)
        self.assertIn("argv=repo view " + AUTHORIZED_REPO, capture)
        self.assert_no_token(result)

    def test_gh_carries_git_credential_boundary(self) -> None:
        result = self.run_helper("gh", "repo", "view", AUTHORIZED_REPO)
        self.assertEqual(result.returncode, 0, result.stderr)
        capture = self.read_capture(self.gh_capture)
        self.assertIn("GIT_CONFIG_KEY_1=credential.helper", capture)
        self.assertIn("GIT_CONFIG_VALUE_1=!gh auth git-credential", capture)


class GitBoundaryTests(HelperTestCase):
    def test_git_uses_process_scoped_credential_boundary(self) -> None:
        result = self.run_helper("git", "fetch", "origin")
        self.assertEqual(result.returncode, 0, result.stderr)
        capture = self.read_capture(self.git_capture)
        self.assertIn(f"GH_TOKEN={SYNTHETIC_TOKEN}", capture)
        self.assertIn("GIT_CONFIG_COUNT=4", capture)
        # Reset inherited helpers...
        self.assertIn("GIT_CONFIG_KEY_0=credential.helper", capture)
        self.assertIn("GIT_CONFIG_VALUE_0=", capture)
        # ...then use gh's helper for HTTPS.
        self.assertIn("GIT_CONFIG_KEY_1=credential.helper", capture)
        self.assertIn("GIT_CONFIG_VALUE_1=!gh auth git-credential", capture)
        self.assertIn("GIT_CONFIG_KEY_2=url.https://github.com/.insteadOf", capture)
        self.assert_no_token(result)

    def test_git_requires_arguments(self) -> None:
        result = self.run_helper("git")
        self.assertNotEqual(result.returncode, 0)
        self.assertFalse(self.git_capture.exists())


class RegistryLoginTests(HelperTestCase):
    def test_feeds_token_via_stdin_with_isolated_config(self) -> None:
        config_dir = self.root / "docker-config"
        result = self.run_helper("registry-login", "--config", str(config_dir))
        self.assertEqual(result.returncode, 0, result.stderr)
        argv = self.read_capture(self.docker_capture)
        argv_line = next(
            line for line in argv.splitlines() if line.startswith("argv=")
        )
        self.assertIn(
            f"argv=login ghcr.io --username {EXPECTED_LOGIN} --password-stdin "
            f"--config {config_dir}",
            argv,
        )
        self.assertNotIn(SYNTHETIC_TOKEN, argv_line, "token must not be in argv")
        self.assertEqual(
            self.read_capture(self.docker_stdin).strip(), SYNTHETIC_TOKEN
        )
        self.assertTrue(config_dir.is_dir())
        mode = stat.S_IMODE(config_dir.stat().st_mode)
        self.assertEqual(mode, 0o700)
        self.assert_no_token(result)

    def test_rejects_existing_config_dir(self) -> None:
        config_dir = self.root / "docker-config"
        config_dir.mkdir()
        result = self.run_helper("registry-login", "--config", str(config_dir))
        self.assertNotEqual(result.returncode, 0)
        self.assertFalse(self.docker_capture.exists(), "docker must not run")

    def test_rejects_relative_config_dir(self) -> None:
        result = self.run_helper("registry-login", "--config", "relative-dir")
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("absolute", result.stderr)


class SyncCiSecretTests(HelperTestCase):
    def test_feeds_token_via_stdin_without_token_argument(self) -> None:
        result = self.run_helper("sync-ci-secret", "--repo", AUTHORIZED_REPO)
        self.assertEqual(result.returncode, 0, result.stderr)
        argv = self.read_capture(self.gh_capture)
        argv_line = next(
            line for line in argv.splitlines() if line.startswith("argv=")
        )
        self.assertIn("argv=secret set PROJECT_GITHUB_TOKEN --repo " + AUTHORIZED_REPO, argv)
        self.assertNotIn(SYNTHETIC_TOKEN, argv_line, "token must not be in argv")
        self.assertEqual(self.read_capture(self.gh_stdin).strip(), SYNTHETIC_TOKEN)
        self.assert_no_token(result)

    def test_rejects_url_repo(self) -> None:
        result = self.run_helper(
            "sync-ci-secret", "--repo", "https://github.com/huyle-gradion/x"
        )
        self.assertNotEqual(result.returncode, 0)
        self.assertFalse(self.gh_capture.exists(), "gh must not run on invalid repo")

    def test_rejects_flag_like_repo(self) -> None:
        result = self.run_helper("sync-ci-secret", "--repo", "--admin")
        self.assertNotEqual(result.returncode, 0)
        self.assertFalse(self.gh_capture.exists())


class DisclosureTests(HelperTestCase):
    def test_token_absent_from_all_outputs(self) -> None:
        config_dir = self.root / "docker-config"
        commands = [
            ("verify",),
            ("gh", "api", "user"),
            ("git", "remote", "-v"),
            ("registry-login", "--config", str(config_dir)),
            ("sync-ci-secret", "--repo", AUTHORIZED_REPO),
        ]
        for args in commands:
            with self.subTest(command=args[0]):
                result = self.run_helper(*args, login=EXPECTED_LOGIN)
                self.assert_no_token(result)


if __name__ == "__main__":
    unittest.main(verbosity=2)
