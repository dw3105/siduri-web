#!/usr/bin/env python3
"""Every repo path named in an agent-facing file must exist.

    python3 tools/paths_exist.py AGENTS-FILLABLE.md AGENTS-PROHIBITIONS.md

`AGENTS.md` is the one file every Codex worker reads, and it named four paths
that did not exist -- two spec files renamed out from under it, and two documents
never written. Nothing refused it. A worker's first act is to open those.

Stdlib only. Exit 1 on a path that is named and absent.
"""
from __future__ import annotations

import pathlib
import re
import sys

#: A repo path must carry a directory. A bare `routes.go` in prose is a
#: reference to a file whose directory the reader already knows -- treating it as
#: a path reports four false positives and buries the real one.
#: A citation carries a line reference -- `docs/site-requirements.md:336`. An
#: earlier version required the closing backtick straight after the path, so
#: every cited path was invisible, which is every path that matters.
PATH = re.compile(r"`((?:docs|internal|content|tools|skills|cmd|web|static|mk)/[\w./-]+?)(?::\d+(?:-\d+)?)?`")


def named_paths(text: str) -> set[str]:
    return {m.group(1).rstrip(".,;:") for m in PATH.finditer(text)}


def main(argv: list[str]) -> int:
    if not argv:
        print("usage: paths_exist.py <file>...", file=sys.stderr)
        return 2
    missing: list[tuple[str, str]] = []
    for name in argv:
        text = pathlib.Path(name).read_text(encoding="utf-8")
        for path in sorted(named_paths(text)):
            # `docs/adr/0002` names an ADR by number; the file carries a slug.
            if pathlib.Path(path).exists():
                continue
            parent = pathlib.Path(path).parent
            if parent.is_dir() and any(c.name.startswith(pathlib.Path(path).name) for c in parent.iterdir()):
                continue
            if True:
                missing.append((name, path))
    if not missing:
        print(f"paths_exist: {len(argv)} file(s), every named path resolves.")
        return 0
    print("paths_exist: named path does not exist", file=sys.stderr)
    for name, path in missing:
        print(f"  {name}: {path}", file=sys.stderr)
    return 1


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
