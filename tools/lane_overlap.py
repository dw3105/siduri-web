#!/usr/bin/env python3
"""Refuse a wave whose lane tasks name the same file twice.

    python3 tools/lane_overlap.py ~/siduri-tasks/A*.md
    python3 tools/lane_overlap.py --selftest

`FO-03` says a lane owns only new files and a file two lanes need belongs to the
integrator. `FO-12` said that a shared line found while writing a lane task is a
seam defect -- and fired today only because somebody noticed. Noticing is the
variable under load at eight tasks.

A filename appearing in two task files is `FO-03` restated as a grep. Non-empty
output means the wave is not ready.

Stdlib only. Exit 1 on overlap, so it can sit in a gate.
"""
from __future__ import annotations

import re
import sys

#: A path is a backticked token carrying a slash, or one of the root files a
#: lane can actually own. A bare `_templ.go` is a generated-suffix mentioned in
#: prose, not a file -- requiring the slash is what tells those apart.
PATH = re.compile(r"`([A-Za-z0-9_.-]*/[A-Za-z0-9_./-]+|_headers|wrangler\.toml|\.gitignore)`")

#: Files every task legitimately cites -- the contract, the agent file, the seam
#: they all register against. Naming these is reading, never owning.
SHARED_BY_DESIGN = {
    "docs/site-requirements.md",
    "docs/comments-requirements.md",
    "AGENTS.md",
    "CLAUDE.md",
    "internal/site/routes.go",
    "internal/site/build.go",
    "internal/site/content.go",
    "internal/site/markdown.go",
    "Makefile",
    "go.mod",
    "go.sum",
}


def paths_in(text: str) -> set[str]:
    found = set()
    for match in PATH.finditer(text):
        token = match.group(1).rstrip("/")
        if token and token not in SHARED_BY_DESIGN:
            found.add(token)
    return found


def overlaps(tasks: dict[str, str]) -> dict[str, list[str]]:
    """Map each file named by more than one task to the tasks naming it."""
    owners: dict[str, list[str]] = {}
    for name, text in sorted(tasks.items()):
        for path in paths_in(text):
            owners.setdefault(path, []).append(name)
    return {p: n for p, n in sorted(owners.items()) if len(n) > 1}


def selftest() -> int:
    clean = {
        "A1": "Own `internal/site/article_a1.go` and read `docs/site-requirements.md`.",
        "A2": "Own `internal/site/tags_a2.go` and read `docs/site-requirements.md`.",
    }
    assert overlaps(clean) == {}, overlaps(clean)

    dirty = dict(clean, A3="Also edit `internal/site/tags_a2.go` for the feed.")
    assert overlaps(dirty) == {"internal/site/tags_a2.go": ["A2", "A3"]}, overlaps(dirty)

    # A file shared by design is read by every lane and is not an overlap.
    everyone = {"A1": "read `Makefile`", "A2": "read `Makefile`"}
    assert overlaps(everyone) == {}, overlaps(everyone)

    # A generated-file suffix named in prose is not a path. This false positive
    # fired on the first real run, against W0 and W1 both saying `_templ.go`.
    prose = {"A1": "commit the `_templ.go` files", "A2": "commit the `_templ.go` files"}
    assert overlaps(prose) == {}, overlaps(prose)

    # A root file a lane really can own is still caught.
    roots = {"A8": "own `wrangler.toml`", "A7": "also `wrangler.toml`"}
    assert overlaps(roots) == {"wrangler.toml": ["A7", "A8"]}, overlaps(roots)

    print("lane_overlap selftest: 5 cases pass")
    return 0


def main(argv: list[str]) -> int:
    if "--selftest" in argv:
        return selftest()
    if not argv:
        print("usage: lane_overlap.py <task.md>...", file=sys.stderr)
        return 2
    tasks = {}
    for path in argv:
        with open(path, encoding="utf-8") as handle:
            tasks[path.rsplit("/", 1)[-1].removesuffix(".md")] = handle.read()
    found = overlaps(tasks)
    if not found:
        print(f"lane_overlap: {len(tasks)} tasks, no file named by two. Wave is cuttable.")
        return 0
    print(f"lane_overlap: {len(found)} file(s) named by more than one task -- seam defect, not task problem", file=sys.stderr)
    for path, names in found.items():
        print(f"  {path}  <- {', '.join(names)}", file=sys.stderr)
    return 1


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
