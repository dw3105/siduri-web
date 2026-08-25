#!/usr/bin/env python3
"""Guard the pre-P3 output against an empty article set and sales artifacts.

    python3 tools/pre_p3.py dist
    python3 tools/pre_p3.py --selftest

The article-page requirement must measure a real page. An empty ``dist/`` is
not evidence that the negative assertion holds, so this exits non-zero before
scanning for forbidden artifacts when no article page exists.

Stdlib only. Exit 1 on an empty article set or a forbidden artifact.
"""
from __future__ import annotations

import argparse
import pathlib
import re
import tempfile


FORBIDDEN = re.compile(r"/services|(?:€|\$)\s*[0-9]|<form(?:\s|>|$)", re.IGNORECASE)


class GuardError(Exception):
    """A pre-P3 assertion cannot be proven for this build."""


def article_pages(root: pathlib.Path) -> list[pathlib.Path]:
    """Return article pages only, in deterministic order."""
    journal = root / "journal"
    if not journal.is_dir():
        return []
    return sorted(
        path for path in journal.rglob("index.html") if path.is_file() and path.parent != journal
    )


def validate(root: pathlib.Path) -> list[pathlib.Path]:
    """Validate *root* and return the article pages that were measured."""
    pages = article_pages(root)
    if not pages:
        raise GuardError(
            f"pre-P3 guard refused: no article pages found in {root}; "
            "build the preview with drafts to inspect article coverage"
        )

    violations: list[tuple[pathlib.Path, str]] = []
    for path in sorted(candidate for candidate in root.rglob("*") if candidate.is_file()):
        data = path.read_bytes()
        if b"\x00" in data[:8192]:
            continue
        text = data.decode("utf-8", errors="ignore")
        match = FORBIDDEN.search(text)
        if match:
            violations.append((path, match.group(0)))
    if violations:
        details = "; ".join(f"{path}: {match!r}" for path, match in violations)
        raise GuardError(f"pre-P3 guard found a forbidden artifact: {details}")
    return pages


def selftest() -> int:
    with tempfile.TemporaryDirectory(prefix="siduri-pre-p3-") as directory:
        root = pathlib.Path(directory)
        try:
            validate(root)
        except GuardError as error:
            assert "no article pages" in str(error)
        else:
            raise AssertionError("empty output was accepted")

        page = root / "journal" / "fixture" / "index.html"
        page.parent.mkdir(parents=True)
        page.write_text("<main><article>fixture</article></main>\n", encoding="utf-8")
        assert len(validate(root)) == 1

        (root / "other.html").write_text("<p>/services/</p>\n", encoding="utf-8")
        try:
            validate(root)
        except GuardError as error:
            assert "forbidden artifact" in str(error)
        else:
            raise AssertionError("forbidden artifact was accepted")

        (root / "binary.png").write_bytes(b"\x89PNG\x00/services\x003\x00")
        (root / "other.html").unlink()
        assert len(validate(root)) == 1
    print("pre_p3 selftest: 4 cases pass")
    return 0


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("root", nargs="?", type=pathlib.Path)
    parser.add_argument("--selftest", action="store_true")
    args = parser.parse_args()
    if args.selftest:
        return selftest()
    if args.root is None:
        parser.error("the build output directory is required")
    try:
        pages = validate(args.root)
    except GuardError as error:
        print(error)
        return 1
    build_kind = "preview" if args.root.name in {".dev-dist", "preview"} else "published"
    print(f"pre-P3 guard: measured {len(pages)} article page(s) in {args.root} ({build_kind} build)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
