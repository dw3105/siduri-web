#!/usr/bin/env python3
"""Measure the static site's NFR-1 budgets.

The command deliberately measures bytes instead of printing a generic success:

    python3 tools/budget.py --dist dist --build-seconds 1.24
    python3 tools/budget.py --selftest

Compression is Python's gzip level 9 with a zero mtime.  That is a stable,
conservative local measurement; it is not a claim that Cloudflare uses the
same compressor.  One KiB is printed as one KB here, matching the site's
budget language and keeping the unit beside every limit.
"""
from __future__ import annotations

import argparse
import gzip
import html.parser
import os
import pathlib
import re
import tempfile
from urllib.parse import urlsplit


JS_BUDGET = 30 * 1024
CSS_BUDGET = 15 * 1024
HTML_BUDGET = 60 * 1024
BUILD_BUDGET = 15.0


class ResourceParser(html.parser.HTMLParser):
    """Collect browser-fetching attributes, not ordinary anchor links."""

    RESOURCE_TAGS = {
        "script": {"src"},
        "img": {"src", "srcset"},
        "iframe": {"src"},
        "video": {"src", "poster"},
        "audio": {"src"},
        "source": {"src", "srcset"},
        "object": {"data"},
        "embed": {"src"},
        "input": {"src"},
    }

    def __init__(self) -> None:
        super().__init__(convert_charrefs=True)
        self.urls: list[str] = []

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        attrs_by_name = {name.lower(): value for name, value in attrs}
        names = self.RESOURCE_TAGS.get(tag.lower(), set())
        for name in names:
            value = attrs_by_name.get(name)
            if value:
                self.urls.extend(value.split(",") if name == "srcset" else [value])
        if tag.lower() == "link":
            rel = {part.lower() for part in (attrs_by_name.get("rel") or "").split()}
            if rel & {"stylesheet", "preload", "modulepreload", "icon"}:
                value = attrs_by_name.get("href")
                if value:
                    self.urls.append(value)


def gzip_size(data: bytes) -> int:
    return len(gzip.compress(data, compresslevel=9, mtime=0))


def article_pages(dist: pathlib.Path) -> list[pathlib.Path]:
    pages = []
    for path in dist.rglob("index.html"):
        relative = path.relative_to(dist).parts
        if len(relative) == 3 and relative[0] == "journal" and relative[2] == "index.html":
            pages.append(path)
    return sorted(pages)


def external_resource_count(pages: list[pathlib.Path], dist: pathlib.Path) -> int:
    urls: set[str] = set()
    for page in pages:
        parser = ResourceParser()
        parser.feed(page.read_text(encoding="utf-8", errors="replace"))
        for value in parser.urls:
            parsed = urlsplit(value.strip().split()[0])
            if parsed.scheme in {"http", "https"} and parsed.netloc:
                urls.add(value.strip())
            elif value.strip().startswith("//"):
                urls.add(value.strip())

    # CSS can introduce requests even when the HTML is clean.
    for css in sorted(dist.rglob("*.css")):
        text = css.read_text(encoding="utf-8", errors="replace")
        for match in re.finditer(r"url\(\s*['\"]?([^)'\"]+)", text, re.IGNORECASE):
            value = match.group(1).strip()
            parsed = urlsplit(value)
            if parsed.scheme in {"http", "https"} and parsed.netloc:
                urls.add(value)
            elif value.startswith("//"):
                urls.add(value)
    return len(urls)


def measure(dist: pathlib.Path, build_seconds: float | None) -> tuple[list[str], dict[str, object]]:
    failures: list[str] = []
    if not dist.is_dir():
        return [f"distribution directory does not exist: {dist}"], {}

    pages = article_pages(dist)
    all_html = sorted(dist.rglob("*.html"))
    css_files = sorted(dist.rglob("*.css"))
    js_files = sorted(dist.rglob("*.js"))
    js_bytes = sum(gzip_size(path.read_bytes()) for path in js_files)
    css_bytes = sum(gzip_size(path.read_bytes()) for path in css_files)
    html_bytes = {path: path.stat().st_size for path in pages}
    third_party = external_resource_count(all_html, dist)

    print("Performance budgets (gzip -9; 1 KB = 1024 bytes):")
    print(f"  JavaScript: {js_bytes / 1024:.2f} KB gzipped / {JS_BUDGET / 1024:.2f} KB budget")
    print(f"  CSS: {css_bytes / 1024:.2f} KB gzipped / {CSS_BUDGET / 1024:.2f} KB budget")
    if not pages:
        print("  HTML: 0 pages measured / at least 1 article page required")
        failures.append("no article pages found under dist/journal/<slug>/index.html")
    else:
        for path, size in html_bytes.items():
            print(f"  HTML {path.relative_to(dist)}: {size / 1024:.2f} KB / {HTML_BUDGET / 1024:.2f} KB budget")
    print(f"  Third-party first-load requests: {third_party} requests / 0 requests budget")
    if build_seconds is None:
        print("  Full build: not measured / < 15.00 s budget")
        failures.append("full build duration is required (--build-seconds)")
    else:
        print(f"  Full build: {build_seconds:.2f} s / < {BUILD_BUDGET:.2f} s budget")

    if js_bytes > JS_BUDGET:
        failures.append(f"JavaScript budget breached: {js_bytes} > {JS_BUDGET} bytes gzipped")
    if css_bytes > CSS_BUDGET:
        failures.append(f"CSS budget breached: {css_bytes} > {CSS_BUDGET} bytes gzipped")
    for path, size in html_bytes.items():
        if size > HTML_BUDGET:
            failures.append(f"HTML budget breached: {path} is {size} > {HTML_BUDGET} bytes")
    if third_party:
        failures.append(f"third-party request budget breached: {third_party} request(s)")
    if build_seconds is not None and build_seconds >= BUILD_BUDGET:
        failures.append(f"full build budget breached: {build_seconds:.2f} s >= {BUILD_BUDGET:.2f} s")

    return failures, {
        "pages": pages,
        "js": js_bytes,
        "css": css_bytes,
        "html": html_bytes,
        "third_party": third_party,
    }


def selftest() -> int:
    """Prove the happy path and every budget's failure path."""
    with tempfile.TemporaryDirectory(prefix="siduri-budget-") as raw:
        root = pathlib.Path(raw)
        (root / "journal/example").mkdir(parents=True)
        (root / "journal/example/index.html").write_text("<html><body>article</body></html>\n", encoding="utf-8")
        (root / "site.css").write_text("body { color: black; }\n", encoding="utf-8")
        failures, _ = measure(root, 1.0)
        assert not failures, f"baseline unexpectedly failed: {failures}"

        (root / "app.js").write_bytes(os.urandom(JS_BUDGET + 2048))
        failures, _ = measure(root, 1.0)
        assert any("JavaScript budget breached" in failure for failure in failures)
        (root / "app.js").unlink()

        (root / "site.css").write_bytes(os.urandom(CSS_BUDGET + 2048))
        failures, _ = measure(root, 1.0)
        assert any("CSS budget breached" in failure for failure in failures)
        (root / "site.css").write_text("body { color: black; }\n", encoding="utf-8")

        (root / "journal/example/index.html").write_bytes(os.urandom(HTML_BUDGET + 2048))
        failures, _ = measure(root, 1.0)
        assert any("HTML budget breached" in failure for failure in failures)
        (root / "journal/example/index.html").write_text("<html><body>article</body></html>\n", encoding="utf-8")

        (root / "journal/example/index.html").write_text(
            '<html><body><script src="https://cdn.example.test/app.js"></script></body></html>',
            encoding="utf-8",
        )
        failures, _ = measure(root, 1.0)
        assert any("third-party request budget breached" in failure for failure in failures)
        (root / "journal/example/index.html").write_text("<html><body>article</body></html>\n", encoding="utf-8")

        failures, _ = measure(root, BUILD_BUDGET)
        assert any("full build budget breached" in failure for failure in failures)
    print("budget selftest: baseline plus JavaScript, CSS, HTML, third-party, and build-time breaches pass")
    return 0


def main(argv: list[str]) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--dist", type=pathlib.Path, default=pathlib.Path("dist"))
    parser.add_argument("--build-seconds", type=float)
    parser.add_argument("--selftest", action="store_true")
    args = parser.parse_args(argv)
    if args.selftest:
        return selftest()
    failures, _ = measure(args.dist, args.build_seconds)
    if failures:
        print("Budget failures:", file=os.sys.stderr)
        for failure in failures:
            print(f"  {failure}", file=os.sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main(os.sys.argv[1:]))
