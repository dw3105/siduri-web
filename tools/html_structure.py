#!/usr/bin/env python3
"""Validate generated HTML with a nesting-aware stdlib parser.

    python3 tools/html_structure.py dist
    python3 tools/html_structure.py --selftest

Tag counts cannot distinguish a close tag belonging to the wrong opener. This
check keeps the open-element stack, rejects mismatched closes, and rejects a
nested ``article`` so an article page has one owner for its page-level article
element. Comment cards use ``div`` containers inside that article.

Stdlib only. Exit 1 on a structural mismatch.
"""
from __future__ import annotations

import argparse
from html.parser import HTMLParser
import pathlib


VOID_ELEMENTS = {
    "area",
    "base",
    "br",
    "col",
    "embed",
    "hr",
    "img",
    "input",
    "link",
    "meta",
    "param",
    "source",
    "track",
    "wbr",
}


class StructureError(ValueError):
    """Generated markup is not structurally balanced."""


class StackParser(HTMLParser):
    def __init__(self) -> None:
        super().__init__(convert_charrefs=False)
        self.stack: list[str] = []

    def location(self) -> str:
        line, column = self.getpos()
        return f"line {line}, column {column + 1}"

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        del attrs
        tag = tag.lower()
        if tag in VOID_ELEMENTS:
            return
        if tag == "article" and "article" in self.stack:
            raise StructureError(f"nested <article> at {self.location()}")
        self.stack.append(tag)

    def handle_startendtag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        del attrs
        if tag.lower() not in VOID_ELEMENTS:
            raise StructureError(f"self-closing non-void <{tag}> at {self.location()}")

    def handle_endtag(self, tag: str) -> None:
        tag = tag.lower()
        if not self.stack:
            raise StructureError(f"unexpected </{tag}> at {self.location()}")
        expected = self.stack[-1]
        if expected != tag:
            raise StructureError(
                f"mismatched </{tag}> at {self.location()}; expected </{expected}>"
            )
        self.stack.pop()


def validate_html(path: pathlib.Path) -> None:
    parser = StackParser()
    try:
        parser.feed(path.read_text(encoding="utf-8"))
        parser.close()
    except (OSError, UnicodeError) as error:
        raise StructureError(f"{path}: {error}") from error
    if parser.stack:
        raise StructureError(f"{path}: unclosed <{parser.stack[-1]}>")


def html_files(paths: list[pathlib.Path]) -> list[pathlib.Path]:
    files: set[pathlib.Path] = set()
    for path in paths:
        if path.is_dir():
            files.update(candidate for candidate in path.rglob("*.html") if candidate.is_file())
        elif path.is_file() and path.suffix == ".html":
            files.add(path)
        else:
            raise StructureError(f"HTML input does not exist or is not an HTML file: {path}")
    return sorted(files)


def selftest() -> int:
    def accepts(source: str) -> None:
        parser = StackParser()
        parser.feed(source)
        parser.close()
        assert not parser.stack

    def refuses(source: str, expected: str) -> None:
        parser = StackParser()
        try:
            parser.feed(source)
            parser.close()
        except StructureError as error:
            assert expected in str(error), (expected, error)
        else:
            raise AssertionError("malformed HTML was accepted")

    accepts("<main><article><div><div></div></div></article></main>")
    refuses(
        "<main><article><div><section></div></section></article></main>",
        "mismatched",
    )
    refuses(
        "<main><article><div><article><div></div></article></div></article></main>",
        "nested <article>",
    )
    print("html_structure selftest: 3 cases pass")
    return 0


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("paths", nargs="*", type=pathlib.Path)
    parser.add_argument("--selftest", action="store_true")
    args = parser.parse_args()
    if args.selftest:
        return selftest()
    if not args.paths:
        parser.error("at least one HTML file or directory is required")
    try:
        files = html_files(args.paths)
        for path in files:
            validate_html(path)
    except StructureError as error:
        print(f"html_structure: {error}")
        return 1
    print(f"html_structure: validated {len(files)} HTML file(s)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
