#!/usr/bin/env python3
"""Check built output for client-side storage and cookie-banner markup.

    python3 tools/storagescan.py dist
    python3 tools/storagescan.py --selftest

This checks emitted files for the storage APIs named by the site contract,
Set-Cookie in a meta tag or ``_headers``, and cookie-banner vocabulary in
class/id attributes.  A ``<script>`` tag by itself is not a finding.

Stdlib only.  Exit 1 on a malformed/missing output tree or a finding.
"""
from __future__ import annotations

import argparse
import pathlib
import re
import shutil
import sys
import tempfile
from dataclasses import dataclass


class ScanError(Exception):
    """The guard could not measure the requested tree."""


@dataclass(frozen=True)
class Violation:
    path: pathlib.Path
    line: int
    marker: str

    def message(self) -> str:
        return f"storagescan: {self.path}:{self.line}: forbidden {self.marker}"


STORAGE_PATTERNS = (
    ("localStorage", re.compile(r"\blocalStorage\b")),
    ("sessionStorage", re.compile(r"\bsessionStorage\b")),
    ("indexedDB", re.compile(r"\bindexedDB\b")),
    ("document.cookie", re.compile(r"\bdocument\.cookie\b")),
)
SET_COOKIE = re.compile(r"\bSet-Cookie\b", re.IGNORECASE)
META_SET_COOKIE = re.compile(r"<meta\b[^>]*\bSet-Cookie\b[^>]*>", re.IGNORECASE | re.DOTALL)
ATTRIBUTE = re.compile(
    r"\b(?:class|id)\s*=\s*(?:\"(?P<double>[^\"]*)\"|'(?P<single>[^']*)'|(?P<bare>[^\s>]+))",
    re.IGNORECASE | re.DOTALL,
)
COOKIE_BANNER = re.compile(
    r"(?:cookie[\s_-]*(?:consent|banner|bar|notice|preferences?|settings?|popup|modal)|"
    r"(?:gdpr|ccpa)[\s_-]*(?:consent|banner|notice|preferences?|settings?)|"
    r"consent[\s_-]*(?:banner|notice|popup)|privacy[\s_-]*banner)",
    re.IGNORECASE,
)


def _files(root: pathlib.Path) -> list[pathlib.Path]:
    if not root.is_dir():
        raise ScanError(f"storagescan: dist directory not found: {root}")
    return sorted(path for path in root.rglob("*") if path.is_file())


def _read_text(path: pathlib.Path) -> str:
    try:
        return path.read_bytes().decode("utf-8", errors="replace")
    except OSError as error:
        raise ScanError(f"storagescan: cannot read {path}: {error}") from error


def _line_at(text: str, offset: int) -> int:
    return text.count("\n", 0, offset) + 1


def scan(dist: pathlib.Path) -> list[Violation]:
    violations: list[Violation] = []
    for path in _files(dist):
        text = _read_text(path)
        lines = text.splitlines()
        if not lines:
            lines = [""]

        for line_number, line in enumerate(lines, 1):
            for marker, pattern in STORAGE_PATTERNS:
                if pattern.search(line):
                    violations.append(Violation(path, line_number, marker))
            if path.name == "_headers" and SET_COOKIE.search(line):
                violations.append(Violation(path, line_number, "Set-Cookie in _headers"))

        for match in META_SET_COOKIE.finditer(text):
            violations.append(Violation(path, _line_at(text, match.start()), "Set-Cookie in meta tag"))

        for match in ATTRIBUTE.finditer(text):
            value = match.group("double") or match.group("single") or match.group("bare") or ""
            banner = COOKIE_BANNER.search(value)
            if banner:
                violations.append(
                    Violation(path, _line_at(text, match.start()), f"cookie-banner markup ({banner.group(0)})")
                )
    return violations


def _case(dist: pathlib.Path, relative: str, content: str, marker: str) -> str:
    shutil.rmtree(dist)
    dist.mkdir()
    path = dist / relative
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content, encoding="utf-8")
    violations = scan(dist)
    assert violations and any(marker in violation.marker for violation in violations), marker
    return violations[0].message()


def selftest() -> int:
    with tempfile.TemporaryDirectory(prefix="siduri-storagescan-") as directory:
        dist = pathlib.Path(directory) / "dist"
        dist.mkdir()
        (dist / "index.html").write_text(
            "<!doctype html><html><body><script type=\"application/ld+json\">{}</script></body></html>\n",
            encoding="utf-8",
        )
        assert scan(dist) == [], "clean output or JSON-LD script was rejected"
        red_cases = [
            ("localStorage", "storage.html", '<script>localStorage.setItem("x", "y")</script>\n', "localStorage"),
            ("sessionStorage", "storage.html", '<script>sessionStorage.removeItem("x")</script>\n', "sessionStorage"),
            ("indexedDB", "storage.html", '<script>indexedDB.open("db")</script>\n', "indexedDB"),
            ("document.cookie", "storage.html", '<script>document.cookie = "x=y"</script>\n', "document.cookie"),
            ("Set-Cookie header", "_headers", "/*\n  Set-Cookie: siduri=1\n*/\n", "Set-Cookie in _headers"),
            ("Set-Cookie meta", "meta.html", '<meta http-equiv="Set-Cookie" content="siduri=1">\n', "Set-Cookie in meta"),
            ("cookie banner", "banner.html", '<div id="gdpr-consent" class="cookie-banner"></div>\n', "cookie-banner markup"),
        ]
        observed = [
            f"{label} -> {_case(dist, relative, content, marker)}"
            for label, relative, content, marker in red_cases
        ]
    print("storagescan selftest: clean case and script allowance pass; " + "; ".join(observed))
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
    print(f"storagescan: checked emitted files in {args.dist}; no storage or cookie-banner markers found")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
