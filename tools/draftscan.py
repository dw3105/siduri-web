#!/usr/bin/env python3
"""Check that draft post content is absent from a built site.

    python3 tools/draftscan.py dist
    python3 tools/draftscan.py --selftest

The source posts are deliberately part of this check.  A clean output tree
does not prove the check works when there are no drafts, so the self-test builds
an isolated source/output tree and makes each kind of draft leak go red.

Stdlib only.  Exit 1 on a malformed source tree or a draft leak.
"""
from __future__ import annotations

import argparse
import html
import pathlib
import re
import shutil
import sys
import tempfile
from dataclasses import dataclass


class ScanError(Exception):
    """The guard could not measure the requested tree."""


@dataclass(frozen=True)
class DraftPost:
    source: pathlib.Path
    title: str
    slug: str
    summary: str
    plain_summary: str
    body_fragments: tuple[str, ...]


@dataclass(frozen=True)
class Violation:
    path: pathlib.Path
    line: int
    field: str
    value: str

    def message(self) -> str:
        return f"draftscan: {self.path}:{self.line}: draft {self.field} leaked ({self.value!r})"


def _unquote(value: str) -> str:
    if len(value) >= 2 and value[0] == value[-1] and value[0] in {"'", '"'}:
        return value[1:-1]
    return value


def _frontmatter(path: pathlib.Path) -> tuple[dict[str, str], str]:
    try:
        source = path.read_text(encoding="utf-8")
    except OSError as error:
        raise ScanError(f"draftscan: cannot read {path}: {error}") from error
    lines = source.replace("\r\n", "\n").split("\n")
    if not lines or lines[0].strip() != "---":
        raise ScanError(f"draftscan: {path}: frontmatter must start with ---")

    end = next((index for index, line in enumerate(lines[1:], 1) if line.strip() == "---"), -1)
    if end < 0:
        raise ScanError(f"draftscan: {path}: frontmatter has no closing ---")

    fields: dict[str, str] = {}
    index = 1
    while index < end:
        line = lines[index].strip()
        index += 1
        if not line or line.startswith("#"):
            continue
        if ":" not in line:
            raise ScanError(f"draftscan: {path}: malformed frontmatter line {line!r}")
        key, value = line.split(":", 1)
        key = key.strip()
        value = value.strip()
        if key == "tags" and not value:
            while index < end and lines[index].strip().startswith("-"):
                index += 1
            continue
        fields[key] = _unquote(value)
    return fields, "\n".join(lines[end + 1 :]).strip()


def _normalise(value: str) -> str:
    return re.sub(r"\s+", " ", value).strip()


def _body_fragments(body: str) -> tuple[str, ...]:
    """Return rendered-text-sized pieces of Markdown body content."""
    fragments: list[str] = []
    for block in re.split(r"\n\s*\n", body.replace("\r\n", "\n")):
        rendered_lines: list[str] = []
        for line in block.split("\n"):
            stripped = line.strip()
            if stripped.startswith("```"):
                continue
            stripped = re.sub(r"^#{1,6}\s+", "", stripped)
            stripped = re.sub(r"^[-*+]\s+", "", stripped)
            stripped = re.sub(r"!\[([^]]*)\]\([^)]*\)", r"\1", stripped)
            stripped = re.sub(r"\[([^]]+)\]\([^)]*\)", r"\1", stripped)
            stripped = stripped.replace("**", "").replace("__", "")
            stripped = stripped.replace("*", "").replace("_", "").replace("`", "")
            if stripped:
                rendered_lines.append(stripped)
        fragment = _normalise(" ".join(rendered_lines))
        if fragment:
            fragments.append(fragment)
    return tuple(fragments)


def load_drafts(posts_dir: pathlib.Path) -> list[DraftPost]:
    if not posts_dir.is_dir():
        raise ScanError(f"draftscan: posts directory not found: {posts_dir}")

    drafts: list[DraftPost] = []
    for path in sorted(posts_dir.glob("*.md")):
        fields, body = _frontmatter(path)
        if fields.get("draft", "").strip().lower() != "true":
            continue
        required = ("title", "slug", "summary", "plain_summary")
        missing = [field for field in required if not fields.get(field, "").strip()]
        if missing:
            raise ScanError(f"draftscan: {path}: draft is missing {', '.join(missing)}")
        drafts.append(
            DraftPost(
                source=path,
                title=fields["title"],
                slug=fields["slug"],
                summary=fields["summary"],
                plain_summary=fields["plain_summary"],
                body_fragments=_body_fragments(body),
            )
        )
    return drafts


def _files(root: pathlib.Path) -> list[pathlib.Path]:
    if not root.is_dir():
        raise ScanError(f"draftscan: dist directory not found: {root}")
    return sorted(path for path in root.rglob("*") if path.is_file())


def _read_text(path: pathlib.Path) -> str:
    try:
        return path.read_bytes().decode("utf-8", errors="replace")
    except OSError as error:
        raise ScanError(f"draftscan: cannot read {path}: {error}") from error


def _line_for(text: str, needles: tuple[str, ...]) -> int:
    for line_number, line in enumerate(text.splitlines(), 1):
        if any(needle and needle in line for needle in needles):
            return line_number
    return 1


def _needle_variants(value: str) -> tuple[str, ...]:
    return tuple(dict.fromkeys((value, html.escape(value, quote=False))))


def _visible_text(text: str) -> str:
    without_tags = re.sub(r"<[^>]*>", " ", text)
    return _normalise(html.unescape(without_tags))


def scan(dist: pathlib.Path, posts_dir: pathlib.Path | None = None) -> list[Violation]:
    if posts_dir is None:
        posts_dir = pathlib.Path(__file__).resolve().parent.parent / "content" / "posts"
    drafts = load_drafts(posts_dir)
    paths = _files(dist)
    violations: list[Violation] = []

    for path in paths:
        relative_path = path.relative_to(dist).as_posix()
        text = _read_text(path)
        visible = _visible_text(text)
        for draft in drafts:
            if draft.slug and draft.slug in relative_path:
                violations.append(Violation(path, 1, "slug in path", draft.slug))

            fields = (
                ("title", draft.title),
                ("summary", draft.summary),
                ("plain_summary", draft.plain_summary),
            )
            for field, value in fields:
                needles = _needle_variants(value)
                if any(needle in text or needle in visible for needle in needles):
                    violations.append(Violation(path, _line_for(text, needles), field, value))

            for fragment in draft.body_fragments:
                if fragment in visible:
                    violations.append(
                        Violation(path, _line_for(text, (fragment,)), "body text", fragment)
                    )
                    break
    return violations


def _write_draft(path: pathlib.Path) -> None:
    path.write_text(
        """---
title: Selftest Draft Title
slug: selftest-draft-slug
date: 2026-08-27
summary: Selftest draft summary
plain_summary: Selftest draft plain summary
tags:
  - method
draft: true
---

Selftest draft body sentinel.
""",
        encoding="utf-8",
    )


def selftest() -> int:
    with tempfile.TemporaryDirectory(prefix="siduri-draftscan-") as directory:
        root = pathlib.Path(directory)
        posts = root / "content" / "posts"
        dist = root / "dist"
        posts.mkdir(parents=True)
        dist.mkdir()
        _write_draft(posts / "draft.md")

        (dist / "index.html").write_text("<p>Published content only.</p>\n", encoding="utf-8")
        assert scan(dist, posts) == [], "clean output was rejected"
        red_cases = [
            ("slug path", dist / "journal" / "selftest-draft-slug" / "index.html", "<p>clean</p>\n", "slug in path"),
            ("title", dist / "title.html", "<h1>Selftest Draft Title</h1>\n", "title"),
            ("summary", dist / "summary.html", "<p>Selftest draft summary</p>\n", "summary"),
            ("plain_summary", dist / "plain.html", "<p>Selftest draft plain summary</p>\n", "plain_summary"),
            ("body", dist / "body.html", "<p>Selftest draft body sentinel.</p>\n", "body text"),
        ]
        observed: list[str] = []
        for label, path, content, field in red_cases:
            shutil.rmtree(dist)
            dist.mkdir()
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_text(content, encoding="utf-8")
            violations = scan(dist, posts)
            assert violations and any(item.field == field for item in violations), label
            observed.append(f"{label} -> {violations[0].message()}")

    print("draftscan selftest: clean case pass; " + "; ".join(observed))
    return 0


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("dist", nargs="?", type=pathlib.Path)
    parser.add_argument("--selftest", action="store_true")
    args = parser.parse_args()
    if args.selftest:
        return selftest()
    if args.dist is None:
        parser.error("the build output directory is required")
    try:
        violations = scan(args.dist)
    except ScanError as error:
        print(error, file=sys.stderr)
        return 1
    if violations:
        for violation in violations:
            print(violation.message(), file=sys.stderr)
        return 1
    drafts = load_drafts(pathlib.Path(__file__).resolve().parent.parent / "content" / "posts")
    print(f"draftscan: checked {len(drafts)} draft post(s) against {args.dist}; no draft content found")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
