#!/usr/bin/env python3
"""Refuse a wave whose lanes would collide.

    python3 tools/lane_overlap.py ~/siduri-tasks/A*.md          # intent, before cutting
    python3 tools/lane_overlap.py --branches <base> lane/a1 ...  # behaviour, before merging
    python3 tools/lane_overlap.py --selftest

Two checks, and neither substitutes for the other.

**Intent**, over task text, before any worktree is cut. `FO-03`: a lane owns only
new files, and a file two lanes need belongs to the integrator. Cheap, and it
prevents.

**Behaviour**, over branches, before merge. A lane that writes a file it never
named is invisible to the first check and no regex closes that. `FO-08`: eight
branches is 28 pairs. This detects.

Two thresholds, because `FO-03` sets two. A lane-owned file is a defect when
**two** tasks name it. An integrator-owned file is a defect when **one** does --
naming it is already claiming a file the lane does not own.

Stdlib only. Exit 1 on any collision, so it can sit in a gate.
"""
from __future__ import annotations

import pathlib
import posixpath
import re
import subprocess
import sys

#: Enforced by a path check in `make check`; no lane can write them. Naming one
#: is reading. Suppressed entirely.
READ_ONLY_BY_GUARD = {
    "docs/site-requirements.md",
    "docs/comments-requirements.md",
    "docs/DECISIONS.md",
    "docs/FINDINGS.md",
    "AGENTS.md",
    "CLAUDE.md",
}

#: The integrator's, and **written by lanes routinely** unless refused. Any lane
#: running `go get` writes two of these; any lane wanting a target writes the
#: Makefile. Threshold is one, not two.
INTEGRATOR_OWNED = {
    "go.mod",
    "go.sum",
    "Makefile",
    "internal/site/routes.go",
    "internal/site/build.go",
    "internal/site/content.go",
    "internal/site/markdown.go",
}

#: Root files a lane can own whose names the rules below would otherwise reject.
#: The four Cloudflare ones all begin with an underscore, which is the test that
#: rejects `_templ.go`. `_worker.js` is the Pages advanced-mode entry point and
#: `_routes.json` decides which paths invoke the Worker; both are plausible for
#: the deploy lane and both were invisible.
KNOWN_ROOT_FILES = {"_headers", "_redirects", "_worker.js", "_routes.json", "Dockerfile"}

#: Backticked or bare, since a task can name a path either way.
TOKEN = re.compile(r"`([^`\s]+)`|(?<![\w/`])((?:\.{0,2}/)?[\w.-]+(?:/[\w.-]+)+)(?![\w`])")


def is_path(token: str) -> bool:
    """A generated-file suffix is not a path; a root config file is.

    `_templ.go` and `_test.go` are written in prose as suffixes and carry no
    directory. Requiring a slash rejects them -- and rejects `.golangci.yml` too,
    which is why the extension rule exists beside it.
    """
    if "/" in token:
        return True
    if token in KNOWN_ROOT_FILES or token in INTEGRATOR_OWNED:
        return True
    if token.startswith("_"):
        return False
    return "." in token.strip(".")


def normalise(token: str) -> str:
    """`./internal/site/x.go` and `internal/site/x.go` are one file.

    `normpath` alone is right. An earlier version chased the `./` with
    `.lstrip("./")`, which strips *any* leading dot or slash and turned
    `.golangci.yml` into `golangci.yml` -- a root config file quietly renamed by
    its own normaliser. The selftest case for root files caught it.
    """
    return posixpath.normpath(token)


#: A task states what it owns and, deliberately, what it must not touch. Both
#: are filenames in the same file. Scanning the whole text reports every
#: must-not-touch line as a claim -- eight lanes "claiming" `routes.go` when all
#: eight were told to keep away from it. Ownership is a section, so read it.
OWNS_HEADING = re.compile(r"^#+\s*Files you own\s*$", re.M)
NEXT_HEADING = re.compile(r"^#+\s", re.M)


def ownership_block(text: str) -> str:
    """The `Files you own` section, or the whole text when there is no such section."""
    match = OWNS_HEADING.search(text)
    if not match:
        return text
    rest = text[match.end():]
    following = NEXT_HEADING.search(rest)
    return rest[: following.start()] if following else rest


def paths_in(text: str) -> set[str]:
    found = set()
    for match in TOKEN.finditer(text):
        raw = (match.group(1) or match.group(2) or "").rstrip("/,.;:")
        if raw and is_path(raw):
            found.add(normalise(raw))
    return found


def collisions(tasks: dict[str, str]) -> tuple[dict[str, list[str]], dict[str, list[str]]]:
    """Return (lane-owned named twice, integrator-owned named at all)."""
    owners: dict[str, list[str]] = {}
    for name, text in sorted(tasks.items()):
        for path in paths_in(ownership_block(text)):
            if path in READ_ONLY_BY_GUARD:
                continue
            owners.setdefault(path, []).append(name)
    shared = {p: n for p, n in sorted(owners.items()) if p not in INTEGRATOR_OWNED and len(n) > 1}
    claimed = {p: n for p, n in sorted(owners.items()) if p in INTEGRATOR_OWNED}
    return shared, claimed


def branch_writes(base: str, branches: list[str]) -> dict[str, list[str]]:
    """Files actually written by more than one branch. Behaviour, not intent."""
    written: dict[str, list[str]] = {}
    for branch in branches:
        out = subprocess.run(
            # Three dots, not two. `a..b` compares two endpoints, so once one
            # lane merges into `base` every later branch appears to have written
            # the merged lane's files -- a wall of false collisions arriving
            # exactly during integration, the only time this runs. `a...b` is
            # changes on `b` since the merge base, which is the question asked.
            ["git", "diff", "--name-only", f"{base}...{branch}"],
            capture_output=True, text=True, check=True,
        ).stdout.split()
        for path in out:
            written.setdefault(path, []).append(branch)
    return {p: b for p, b in sorted(written.items()) if len(b) > 1}


def report(shared, claimed) -> int:
    if not shared and not claimed:
        return 0
    if claimed:
        print("lane_overlap: task naming an integrator-owned file -- threshold is one, not two", file=sys.stderr)
        for path, names in claimed.items():
            print(f"  {path}  <- {', '.join(names)}", file=sys.stderr)
    if shared:
        print("lane_overlap: file named by more than one task -- seam defect, not task problem", file=sys.stderr)
        for path, names in shared.items():
            print(f"  {path}  <- {', '.join(names)}", file=sys.stderr)
    return 1


def selftest() -> int:
    def check(tasks, want_shared, want_claimed, label):
        got = collisions(tasks)
        assert got == (want_shared, want_claimed), f"{label}: {got}"

    check({"A1": "own `internal/site/a1.go`", "A2": "own `internal/site/a2.go`"}, {}, {}, "clean")
    check({"A1": "own `internal/site/x.go`", "A2": "also `internal/site/x.go`"},
          {"internal/site/x.go": ["A1", "A2"]}, {}, "lane-owned twice")

    # Read-only-by-guard is cited by every task and is never a collision.
    check({"A1": "read `docs/site-requirements.md`", "A2": "read `docs/site-requirements.md`"}, {}, {}, "guarded")

    # A generated suffix named in prose is not a path. Fired on the first real
    # run, against W0 and W1 both saying `_templ.go`.
    check({"A1": "commit `_templ.go`", "A2": "commit `_templ.go`"}, {}, {}, "generated suffix")

    # Requiring a slash fixed that and opened this: a root config file has none.
    # The repair was proved against the case it fixed, not the case it opened.
    check({"A5": "own `.golangci.yml`", "A6": "own `.golangci.yml`"},
          {".golangci.yml": ["A5", "A6"]}, {}, "root file, no slash")

    # One file, two spellings.
    check({"B1": "own `internal/site/x.go`", "B2": "own `./internal/site/x.go`"},
          {"internal/site/x.go": ["B1", "B2"]}, {}, "normalised spelling")

    # Unbackticked path in prose still counts.
    check({"C1": "edit internal/site/plain_y.go here", "C2": "and internal/site/plain_y.go"},
          {"internal/site/plain_y.go": ["C1", "C2"]}, {}, "unbackticked")

    # Integrator-owned: ONE task naming it is already a defect.
    check({"A2": "run go get so `go.mod` changes"}, {}, {"go.mod": ["A2"]}, "integrator, single claim")
    check({"A2": "add a field to `internal/site/routes.go`", "A4": "also `internal/site/routes.go`"},
          {}, {"internal/site/routes.go": ["A2", "A4"]}, "integrator, two claims")

    # `_headers` starts with an underscore and is a real root file a lane owns.
    check({"A8": "own `_headers`", "A7": "also `_headers`"},
          {"_headers": ["A7", "A8"]}, {}, "underscore root file")

    # Every integrator-owned file must be recognised as a path, or its
    # threshold-one branch can never fire. `Makefile` failed this: no dot, no
    # slash, not underscore-prefixed, so `is_path` returned False and two tasks
    # claiming it reported cuttable. Asserting the invariant makes the omission
    # impossible rather than caught.
    # A task states what it owns and, deliberately, what it must not touch.
    # Scanning whole text reported all eight lanes "claiming" routes.go when all
    # eight had been told to keep away from it -- caught before worktrees were cut.
    forbidden = """### Files you own

`internal/site/a1.go`

### Files you must not touch

`internal/site/routes.go`, `go.mod`
"""
    check({"A1": forbidden, "A2": forbidden.replace("a1.go", "a2.go")}, {}, {}, "must-not-touch is not a claim")

    # A genuine claim inside the ownership block is still caught.
    claims = forbidden.replace("`internal/site/a1.go`", "`internal/site/a1.go`\n`go.mod`")
    check({"A1": claims}, {}, {"go.mod": ["A1"]}, "claim inside ownership block")

    unreachable = [p for p in sorted(INTEGRATOR_OWNED) if not is_path(p)]
    assert not unreachable, f"integrator-owned but unreachable by is_path: {unreachable}"

    branch_cases = _selftest_branches()
    print(f"lane_overlap selftest: 12 intent cases, 1 invariant, {branch_cases} branch cases pass")
    return 0


def _selftest_branches() -> int:
    """A base that has moved is the case that will actually occur."""
    import shutil, subprocess as sp, tempfile

    if not shutil.which("git"):
        return 0
    root = tempfile.mkdtemp(prefix="lane_overlap_selftest_")
    try:
        def git(*args):
            sp.run(["git", "-C", root, *args], check=True, capture_output=True, text=True)

        sp.run(["git", "init", "-q", "-b", "main", root], check=True)
        git("config", "user.name", "selftest")
        git("config", "user.email", "selftest@localhost")
        (pathlib.Path(root) / "base.txt").write_text("x\n")
        git("add", "-A"); git("commit", "-q", "-m", "base")
        base = sp.run(["git", "-C", root, "rev-parse", "HEAD"],
                      check=True, capture_output=True, text=True).stdout.strip()

        git("checkout", "-q", "-b", "lane/a")
        (pathlib.Path(root) / "a_only.go").write_text("a\n")
        git("add", "-A"); git("commit", "-q", "-m", "a")

        git("checkout", "-q", "-b", "lane/b", base)
        (pathlib.Path(root) / "b_only.go").write_text("b\n")
        git("add", "-A"); git("commit", "-q", "-m", "b")

        cwd = pathlib.Path.cwd()
        import os
        os.chdir(root)
        try:
            before = branch_writes(base, ["lane/a", "lane/b"])
            assert before == {}, f"disjoint branches before merge: {before}"

            # Move the base, which is what integration does.
            git("checkout", "-q", "main"); git("merge", "-q", "--no-edit", "lane/a")
            after = branch_writes("main", ["lane/b"])
            assert after == {}, f"moved base attributed merged work to lane/b: {after}"

            # And a genuine collision is still caught after the base moved.
            git("checkout", "-q", "lane/b")
            (pathlib.Path(root) / "a_only.go").write_text("b touched it\n")
            git("add", "-A"); git("commit", "-q", "-m", "b touches a_only")
            git("checkout", "-q", "-b", "lane/c", base)
            (pathlib.Path(root) / "a_only.go").write_text("c too\n")
            git("add", "-A"); git("commit", "-q", "-m", "c touches a_only")
            real = branch_writes(base, ["lane/b", "lane/c"])
            assert real == {"a_only.go": ["lane/b", "lane/c"]}, f"real collision missed: {real}"
        finally:
            os.chdir(cwd)
        return 3
    finally:
        shutil.rmtree(root, ignore_errors=True)


def main(argv: list[str]) -> int:
    if "--selftest" in argv:
        return selftest()
    if argv and argv[0] == "--branches":
        if len(argv) < 3:
            print("usage: lane_overlap.py --branches <base> <branch>...", file=sys.stderr)
            return 2
        found = branch_writes(argv[1], argv[2:])
        if not found:
            print(f"lane_overlap: {len(argv) - 2} branches, no file written by two. Merge order is free.")
            return 0
        print(f"lane_overlap: {len(found)} file(s) written by more than one branch", file=sys.stderr)
        for path, branches in found.items():
            print(f"  {path}  <- {', '.join(branches)}", file=sys.stderr)
        return 1
    if not argv:
        print("usage: lane_overlap.py <task.md>... | --branches <base> <branch>... | --selftest", file=sys.stderr)
        return 2
    tasks = {}
    for path in argv:
        with open(path, encoding="utf-8") as handle:
            tasks[path.rsplit("/", 1)[-1].removesuffix(".md")] = handle.read()
    shared, claimed = collisions(tasks)
    code = report(shared, claimed)
    if code == 0:
        print(f"lane_overlap: {len(tasks)} tasks, no collision. Wave is cuttable.")
    return code


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
