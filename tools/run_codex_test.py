#!/usr/bin/env python3
"""Black-box tests for Siduri's bounded Codex runner."""

from __future__ import annotations

import json
import os
from pathlib import Path
import signal
import subprocess
import tempfile
import textwrap
import unittest


ROOT = Path(__file__).resolve().parents[1]
RUNNER = ROOT / "tools" / "run_codex.sh"


class RunCodexTest(unittest.TestCase):
    def setUp(self) -> None:
        self.tempdir = tempfile.TemporaryDirectory()
        self.root = Path(self.tempdir.name)
        self.home = self.root / "home"
        self.home.mkdir()
        self.bin = self.root / "bin"
        self.bin.mkdir()
        self.npm_cache = self.root / "npm-cache"
        self.npm_cache.mkdir()
        self._write_stub_commands()

        self.environment = os.environ.copy()
        self.environment.update(
            {
                "HOME": str(self.home),
                "PATH": f"{self.bin}{os.pathsep}/usr/bin:/bin",
                "STUB_ARGV": str(self.root / "stub-argv.json"),
                "STUB_STDIN": str(self.root / "stub-stdin.txt"),
                "STUB_MODE": "success",
                "STUB_NPM_CACHE": str(self.npm_cache),
                "POLL_S": "0.05",
            }
        )

    def tearDown(self) -> None:
        self.tempdir.cleanup()

    def _write_stub_commands(self) -> None:
        codex = self.bin / "codex"
        codex.write_text(
            textwrap.dedent(
                """
                #!/usr/bin/env python3
                import json
                import os
                import signal
                import sys
                import time

                with open(os.environ["STUB_ARGV"], "w", encoding="utf-8") as handle:
                    json.dump(sys.argv[1:], handle)
                with open(os.environ["STUB_STDIN"], "w", encoding="utf-8") as handle:
                    handle.write(sys.stdin.read())
                if os.environ.get("STUB_PID"):
                    with open(os.environ["STUB_PID"], "w", encoding="utf-8") as handle:
                        handle.write(str(os.getpid()))

                mode = os.environ.get("STUB_MODE", "success")
                if mode == "deadline":
                    signal.signal(signal.SIGTERM, lambda *_args: None)
                    while True:
                        time.sleep(0.05)
                if mode == "stall":
                    time.sleep(float(os.environ.get("STUB_RUN_S", "1")))
                if mode == "exit7":
                    print("stub exited with seven", flush=True)
                    raise SystemExit(7)
                print("stub completed", flush=True)
                """
            ).lstrip(),
            encoding="utf-8",
        )
        codex.chmod(0o755)

        npm = self.bin / "npm"
        npm.write_text(
            "#!/bin/sh\n"
            "test \"$1\" = config && test \"$2\" = get && test \"$3\" = cache || exit 2\n"
            "printf '%s\\n' \"$STUB_NPM_CACHE\"\n",
            encoding="utf-8",
        )
        npm.chmod(0o755)

    def _git_worktree(self) -> tuple[Path, Path]:
        repo = self.root / "repo"
        worktree = self.root / "worktree"
        repo.mkdir()
        subprocess.run(["git", "-C", str(repo), "init", "-q"], check=True)
        subprocess.run(
            ["git", "-C", str(repo), "config", "user.email", "runner@test"],
            check=True,
        )
        subprocess.run(
            ["git", "-C", str(repo), "config", "user.name", "runner-test"],
            check=True,
        )
        (repo / "README").write_text("throwaway\n", encoding="utf-8")
        subprocess.run(["git", "-C", str(repo), "add", "README"], check=True)
        subprocess.run(["git", "-C", str(repo), "commit", "-qm", "initial"], check=True)
        subprocess.run(
            ["git", "-C", str(repo), "worktree", "add", "-q", "-b", "runner-test", str(worktree)],
            check=True,
        )
        gitdir = Path(
            subprocess.run(
                ["git", "-C", str(worktree), "rev-parse", "--git-dir"],
                check=True,
                capture_output=True,
                text=True,
            ).stdout.strip()
        )
        if not gitdir.is_absolute():
            gitdir = (worktree / gitdir).resolve()
        return worktree, gitdir

    def _task(self) -> Path:
        task = self.root / "task.md"
        task.write_text("test prompt\n", encoding="utf-8")
        return task

    def _run(self, *arguments: str, timeout: float = 10) -> subprocess.CompletedProcess[str]:
        return subprocess.run(
            [str(RUNNER), *arguments],
            cwd=ROOT,
            env=self.environment,
            capture_output=True,
            text=True,
            timeout=timeout,
        )

    def _artefact_dir(self, label: str) -> Path:
        return self.home / "siduri-wt"

    def test_missing_deadline_exits_two(self) -> None:
        worktree, _gitdir = self._git_worktree()
        completed = self._run("--label", "missing-deadline", str(worktree), str(self._task()))
        self.assertEqual(completed.returncode, 2, completed.stdout + completed.stderr)
        self.assertIn("usage:", completed.stdout + completed.stderr)

    def test_stub_run_writes_artefacts_and_all_four_add_dirs(self) -> None:
        worktree, gitdir = self._git_worktree()
        task = self._task()
        label = "run-codex-artefacts"
        completed = self._run("--deadline", "5", "--label", label, str(worktree), str(task))
        output = completed.stdout + completed.stderr
        self.assertEqual(completed.returncode, 0, output)

        artefacts = self._artefact_dir(label)
        for suffix in ("log", "status", "pid"):
            self.assertTrue((artefacts / f"{label}.{suffix}").exists())
        self.assertIn("stub completed", (artefacts / f"{label}.log").read_text(encoding="utf-8"))

        argv = json.loads((self.root / "stub-argv.json").read_text(encoding="utf-8"))
        add_dirs = [argv[index + 1] for index, value in enumerate(argv[:-1]) if value == "--add-dir"]
        expected = {
            str(self.home / ".cache" / "go-build"),
            str(self.home / "go" / "pkg" / "mod"),
            str(self.npm_cache),
            str(gitdir),
        }
        self.assertEqual(set(add_dirs), expected)
        self.assertEqual(add_dirs.count(str(gitdir)), 1)
        self.assertEqual(argv[argv.index("-C") + 1], str(worktree))
        self.assertEqual(argv[argv.index("--color") + 1], "never")
        self.assertEqual(argv[-1], "-")
        self.assertEqual((self.root / "stub-stdin.txt").read_text(encoding="utf-8"), "test prompt\n")

    def test_child_exit_status_is_returned_directly(self) -> None:
        worktree, _gitdir = self._git_worktree()
        self.environment["STUB_MODE"] = "exit7"
        completed = self._run("--deadline", "5", "--label", "run-codex-exit7", str(worktree), str(self._task()))
        self.assertEqual(completed.returncode, 7, completed.stdout + completed.stderr)

    def test_deadline_kill_leaves_audit_copy(self) -> None:
        worktree, _gitdir = self._git_worktree()
        label = "run-codex-deadline"
        self.environment.update(
            {
                "STUB_MODE": "deadline",
                "KILL_AFTER_S": "0.1",
                "STUB_PID": str(self.root / "stub.pid"),
            }
        )
        completed = self._run(
            "--deadline", "0.5", "--label", label, str(worktree), str(self._task()), timeout=5
        )
        output = completed.stdout + completed.stderr
        self.assertIn(completed.returncode, (124, 137), output)
        artefacts = self._artefact_dir(label)
        audit = artefacts / "audit" / "runs" / label
        self.assertTrue((audit / f"{label}.log").exists())
        self.assertTrue((audit / f"{label}.status").exists())
        self.assertIn("TIMED OUT", (audit / f"{label}.status").read_text(encoding="utf-8"))

        pid_path = self.root / "stub.pid"
        if pid_path.exists():
            pid = int(pid_path.read_text(encoding="utf-8"))
            try:
                os.kill(pid, 0)
            except ProcessLookupError:
                pass
            else:
                os.kill(pid, signal.SIGKILL)

    def test_stall_report_and_kill(self) -> None:
        worktree, _gitdir = self._git_worktree()
        label = "run-codex-stall"
        self.environment.update(
            {
                "STUB_MODE": "stall",
                "STUB_RUN_S": "4",
                "STALL_S": "1",
                "STALL_KILL_S": "2",
                "STALL_TERM_WAIT_S": "0.05",
            }
        )
        completed = self._run(
            "--deadline", "5", "--label", label, str(worktree), str(self._task()), timeout=5
        )
        output = completed.stdout + completed.stderr
        self.assertNotEqual(completed.returncode, 0, output)
        status = (self._artefact_dir(label) / f"{label}.status").read_text(encoding="utf-8")
        self.assertIn("STALLED", status)
        self.assertIn("STALL-KILL", status)


if __name__ == "__main__":
    unittest.main()
