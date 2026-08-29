#!/usr/bin/env python3
"""Refuse unapproved changes under docs/.

The two requirements files are deliberately not checked here.  The existing
contract checks (``tools/amendcheck.py`` via ``w2-amendcheck``, plus the
working-tree check in ``Makefile``) are authoritative for those paths.

The remaining records are append-only, existing ADRs are immutable, and a new
ADR is the only other file that may be added under docs/.  The committed range
judges only commits that are not already reachable from the published trunk;
the index and working tree are then checked separately so a change cannot be
hidden in one of those three states.
"""
from __future__ import annotations

import argparse
import dataclasses
import difflib
import os
import pathlib
import re
import stat
import subprocess
import sys
from collections.abc import Iterable


CONTRACT_FILES = frozenset(
    {
        "docs/site-requirements.md",
        "docs/comments-requirements.md",
    }
)
# This set is deliberately closed: a missing record file is never appendable,
# so this guard cannot be used to introduce a fourth record file.
RECORD_FILES = frozenset(
    {
        "docs/DECISIONS.md",
        "docs/FINDINGS.md",
        "docs/FRICTION.md",
    }
)
ADR_PATH = re.compile(r"^docs/adr/(?P<number>\d{4})-[^/]+\.md$")


class DocsGuardError(Exception):
    """The guard could not inspect the repository."""


@dataclasses.dataclass(frozen=True)
class FileState:
    mode: str
    data: bytes


def git_bytes(repo: pathlib.Path, arguments: list[str]) -> bytes:
    result = subprocess.run(
        ["git", *arguments],
        cwd=repo,
        capture_output=True,
        check=False,
    )
    if result.returncode:
        detail = result.stderr.decode("utf-8", errors="replace").strip()
        command = "git " + " ".join(arguments)
        raise DocsGuardError(f"{command} failed: {detail or 'exit ' + str(result.returncode)}")
    return result.stdout


def repository_root() -> pathlib.Path:
    result = subprocess.run(
        ["git", "rev-parse", "--show-toplevel"],
        capture_output=True,
        check=False,
        text=True,
    )
    if result.returncode:
        raise DocsGuardError("docsguard must run inside a Git repository")
    return pathlib.Path(result.stdout.strip()).resolve()


def tree_states(repo: pathlib.Path, revision: str) -> dict[str, FileState]:
    states: dict[str, FileState] = {}
    output = git_bytes(repo, ["ls-tree", "-r", "-z", revision, "--", "docs"])
    for entry in output.split(b"\0"):
        if not entry:
            continue
        try:
            header, raw_path = entry.split(b"\t", 1)
            mode, object_type, _object_id = header.split()
        except ValueError as error:
            raise DocsGuardError(f"cannot parse {revision} docs tree entry") from error
        if object_type != b"blob":
            continue
        path = raw_path.decode("utf-8", errors="surrogateescape")
        states[path] = FileState(
            mode.decode("ascii"),
            git_bytes(repo, ["show", f"{revision}:{path}"]),
        )
    return states


def index_states(repo: pathlib.Path) -> dict[str, FileState]:
    states: dict[str, FileState] = {}
    output = git_bytes(repo, ["ls-files", "--stage", "-z", "--", "docs"])
    for entry in output.split(b"\0"):
        if not entry:
            continue
        try:
            header, raw_path = entry.split(b"\t", 1)
            mode, object_id, stage = header.split()
        except ValueError as error:
            raise DocsGuardError("cannot parse the index docs entry") from error
        if stage != b"0":
            path = raw_path.decode("utf-8", errors="surrogateescape")
            raise DocsGuardError(f"unmerged index entry under docs/: {path}")
        # An intent-to-add entry has no blob yet.  The working-tree transition
        # below still sees its file and reports it normally.
        if set(object_id) == {ord("0")}:
            continue
        path = raw_path.decode("utf-8", errors="surrogateescape")
        states[path] = FileState(
            mode.decode("ascii"),
            git_bytes(repo, ["show", f":{path}"]),
        )
    return states


def working_file_state(path: pathlib.Path) -> FileState:
    file_stat = path.lstat()
    if stat.S_ISLNK(file_stat.st_mode):
        data = os.readlink(path).encode("utf-8", errors="surrogateescape")
        return FileState("120000", data)
    if stat.S_ISREG(file_stat.st_mode):
        mode = "100755" if file_stat.st_mode & stat.S_IXUSR else "100644"
        return FileState(mode, path.read_bytes())
    raise DocsGuardError(f"unsupported non-file entry under docs/: {path}")


def working_states(repo: pathlib.Path) -> dict[str, FileState]:
    root = repo / "docs"
    states: dict[str, FileState] = {}
    if not root.exists() and not root.is_symlink():
        return states

    for directory, directories, files in os.walk(root, followlinks=False):
        # os.walk puts symlinked directories in `directories`; treat them as
        # file-like entries so a path replacement cannot disappear from the
        # working-tree view.
        symlink_directories = [
            name for name in directories if (pathlib.Path(directory) / name).is_symlink()
        ]
        for name in symlink_directories:
            directories.remove(name)
            files.append(name)
        for name in files:
            path = pathlib.Path(directory) / name
            relative = path.relative_to(repo).as_posix()
            states[relative] = working_file_state(path)
    return states


def changed_paths(before: dict[str, FileState], after: dict[str, FileState]) -> list[str]:
    return sorted(path for path in before.keys() | after.keys() if before.get(path) != after.get(path))


def additions_only(before: bytes, after: bytes) -> bool:
    # A missing trailing newline makes the next append replace the final line;
    # that is intentionally refused rather than silently treated as an append.
    old_lines = before.splitlines(keepends=True)
    new_lines = after.splitlines(keepends=True)
    matcher = difflib.SequenceMatcher(None, old_lines, new_lines, autojunk=False)
    return all(tag in {"equal", "insert"} for tag, _i1, _i2, _j1, _j2 in matcher.get_opcodes())


def transition_label(label: str, path: str) -> str:
    return f"{label}: {path}"


def adr_number(path: str) -> str | None:
    match = ADR_PATH.fullmatch(path)
    return match.group("number") if match else None


def inspect_transition(
    before: dict[str, FileState],
    after: dict[str, FileState],
    label: str,
    base_paths: set[str],
    base_adr_numbers: set[str],
    delegated: set[str],
) -> list[str]:
    violations: list[str] = []
    for path in changed_paths(before, after):
        if path in CONTRACT_FILES:
            delegated.add(path)
            continue

        old = before.get(path)
        new = after.get(path)

        if path in RECORD_FILES:
            if old is None:
                violations.append(
                    f"{transition_label(label, path)} — record file was added; it must already exist"
                )
            elif new is None:
                violations.append(
                    f"{transition_label(label, path)} — record file was deleted"
                )
            elif old.mode != new.mode:
                violations.append(
                    f"{transition_label(label, path)} — record file mode changed"
                )
            elif not additions_only(old.data, new.data):
                violations.append(
                    f"{transition_label(label, path)} — existing content was changed; only added lines are permitted"
                )
            continue

        if ADR_PATH.fullmatch(path):
            if new is None:
                violations.append(
                    f"{transition_label(label, path)} — an ADR deletion is not permitted"
                )
            elif old is None and adr_number(path) in base_adr_numbers:
                number = adr_number(path)
                violations.append(
                    f"{transition_label(label, path)} — ADR number {number} already exists in the published trunk; duplicate ADR numbers are not permitted"
                )
            elif path in base_paths:
                violations.append(
                    f"{transition_label(label, path)} — existing ADRs are immutable"
                )
            # A new, non-colliding ADR may be added or edited before it is
            # committed.  Published-trunk paths remain immutable.
            continue

        if old is None:
            reason = "new docs files are not permitted"
        elif new is None:
            reason = "deletions under docs/ are not permitted"
        else:
            reason = "changes to existing docs files are not permitted"
        violations.append(f"{transition_label(label, path)} — {reason}")
    return violations


def published_trunk(repo: pathlib.Path) -> tuple[str, str]:
    for display, ref in (
        ("origin/main", "refs/remotes/origin/main"),
        ("main", "refs/heads/main"),
    ):
        try:
            revision = git_bytes(repo, ["rev-parse", "--verify", "--quiet", f"{ref}^{{commit}}"])
        except DocsGuardError:
            continue
        revision = revision.strip().decode("ascii", errors="replace")
        if revision:
            return display, revision
    raise DocsGuardError("cannot resolve published trunk: tried origin/main, then main")


def range_commits(repo: pathlib.Path, base: str, head: str, trunk: str) -> tuple[list[str], list[str]]:
    all_commits = git_bytes(repo, ["rev-list", "--reverse", f"{base}..{head}"])
    judged_commits = git_bytes(repo, ["rev-list", "--reverse", f"{base}..{head}", "--not", trunk])
    all_ids = all_commits.decode("ascii", errors="replace").splitlines()
    judged_ids = judged_commits.decode("ascii", errors="replace").splitlines()
    return all_ids, judged_ids


def assess(repo: pathlib.Path, base: str, head: str) -> tuple[bool, list[str]]:
    trunk_ref, trunk_revision = published_trunk(repo)
    all_commits, judged_commits = range_commits(repo, base, head, trunk_ref)
    judged_word = "commit" if len(judged_commits) == 1 else "commits"
    print(
        f"docsguard: scope — published trunk {trunk_ref}; "
        f"judged {len(judged_commits)} {judged_word} not reachable from {trunk_ref}"
    )
    out_of_scope = [commit for commit in all_commits if commit not in set(judged_commits)]
    if out_of_scope:
        if len(out_of_scope) <= 3:
            details = ", ".join(
                f"{commit[:12]} (already reachable from {trunk_ref}; out of scope)"
                for commit in out_of_scope
            )
            print(f"docsguard: scope — {details}")
        else:
            print(
                f"docsguard: scope — {len(out_of_scope)} commits already reachable from "
                f"{trunk_ref}; out of scope"
            )

    trunk_states = tree_states(repo, trunk_revision)
    live_head = git_bytes(repo, ["rev-parse", "--verify", "HEAD^{commit}"])
    live_head = live_head.strip().decode("ascii", errors="replace")
    head_states = tree_states(repo, live_head)
    index = index_states(repo)
    worktree = working_states(repo)
    delegated: set[str] = set()
    violations: list[str] = []
    base_paths = set(trunk_states)
    base_adr_numbers = {
        number
        for path in base_paths
        if (number := adr_number(path)) is not None
    }

    for commit in judged_commits:
        parent_states = tree_states(repo, f"{commit}^")
        commit_states = tree_states(repo, commit)
        violations.extend(
            inspect_transition(
                parent_states,
                commit_states,
                f"committed commit {commit[:12]}",
                base_paths,
                base_adr_numbers,
                delegated,
            )
        )
    violations.extend(
        inspect_transition(head_states, index, "index", base_paths, base_adr_numbers, delegated)
    )
    violations.extend(
        inspect_transition(
            index, worktree, "working tree", base_paths, base_adr_numbers, delegated
        )
    )
    if delegated:
        delegated_paths = ", ".join(sorted(delegated))
        print(
            "docsguard: PASS — "
            f"{delegated_paths} delegated to tools/amendcheck.py via w2-amendcheck "
            "and the existing Makefile working-tree check; those checks are authoritative"
        )
    return not violations, violations


def parse_args(arguments: Iterable[str]) -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("base", help="the merge-base revision")
    parser.add_argument("head", nargs="?", default="HEAD", help="the committed head revision")
    return parser.parse_args(list(arguments))


def main(arguments: Iterable[str]) -> int:
    args = parse_args(arguments)
    try:
        repo = repository_root()
        passed, violations = assess(repo, args.base, args.head)
    except DocsGuardError as error:
        print(f"docsguard: ERROR — {error}", file=sys.stderr)
        return 2

    if not passed:
        for violation in violations:
            print(f"docsguard: FAIL — {violation}", file=sys.stderr)
        return 1
    print("docsguard: PASS — no forbidden docs changes")
    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
