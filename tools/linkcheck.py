#!/usr/bin/env python3
"""Check generated internal links, with an opt-in network half.

    python3 tools/linkcheck.py --dist dist             # internal links only
    python3 tools/linkcheck.py --dist dist --external  # CI: also HTTP links
    python3 tools/linkcheck.py --selftest

External checks are intentionally opt-in so a lane can work without a network.
The CI workflow passes ``--external``.  ``/services/`` is a known pending route
because the current article fixture names the later services surface before
its lane exists.
"""
from __future__ import annotations

import argparse
import html.parser
import pathlib
import tempfile
import urllib.error
import urllib.parse
import urllib.request


KNOWN_PENDING = {
    "/services/": "the current article fixture names the later services surface before its lane exists",
}


class LinkParser(html.parser.HTMLParser):
    def __init__(self) -> None:
        super().__init__(convert_charrefs=True)
        self.links: list[str] = []

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        attrs_by_name = {name.lower(): value for name, value in attrs}
        tag = tag.lower()
        if tag in {"a", "area"} and attrs_by_name.get("href"):
            self.links.append(attrs_by_name["href"] or "")
        if tag == "link" and attrs_by_name.get("href"):
            rel = {part.lower() for part in (attrs_by_name.get("rel") or "").split()}
            if rel & {"alternate", "icon", "modulepreload", "preload", "stylesheet"}:
                self.links.append(attrs_by_name["href"] or "")
        if tag in {"script", "img", "iframe", "video", "audio", "source", "object", "embed"}:
            for name in ("src", "data"):
                if attrs_by_name.get(name):
                    self.links.append(attrs_by_name[name] or "")


def html_files(dist: pathlib.Path) -> list[pathlib.Path]:
    return sorted(dist.rglob("*.html")) if dist.is_dir() else []


def internal_target(dist: pathlib.Path, page: pathlib.Path, value: str) -> pathlib.Path | None:
    parsed = urllib.parse.urlsplit(value)
    if parsed.scheme or parsed.netloc:
        return None
    path = urllib.parse.unquote(parsed.path)
    if not path:
        return None
    if path.startswith("/"):
        candidate = dist / path.lstrip("/")
    else:
        candidate = page.parent / path
    try:
        candidate = candidate.resolve()
        candidate.relative_to(dist.resolve())
    except ValueError:
        return candidate
    if candidate.is_dir() or path.endswith("/"):
        return candidate / "index.html"
    if candidate.exists():
        return candidate
    if candidate.suffix == "":
        with_html = candidate.with_suffix(".html")
        if with_html.exists():
            return with_html
    return candidate


def check_dist(dist: pathlib.Path, external: bool, timeout: float = 10.0) -> tuple[int, int, list[str]]:
    pages = html_files(dist)
    if not pages:
        return 0, 0, [f"no HTML files found under {dist}"]
    errors: list[str] = []
    internal_count = 0
    external_count = 0
    for page in pages:
        parser = LinkParser()
        parser.feed(page.read_text(encoding="utf-8"))
        for raw in parser.links:
            value = raw.strip()
            if not value or value.startswith(("#", "mailto:", "tel:", "javascript:", "data:")):
                continue
            parsed = urllib.parse.urlsplit(value)
            if parsed.scheme in {"http", "https"} or parsed.netloc or value.startswith("//"):
                external_count += 1
                if not external:
                    continue
                request_value = value if parsed.scheme else f"https:{value}"
                try:
                    request = urllib.request.Request(request_value, method="HEAD", headers={"User-Agent": "siduri-linkcheck/1"})
                    with urllib.request.urlopen(request, timeout=timeout) as response:
                        if response.status >= 400:
                            raise urllib.error.HTTPError(request_value, response.status, "HTTP error", response.headers, None)
                except urllib.error.HTTPError as exc:
                    if exc.code in {405, 501}:
                        try:
                            request = urllib.request.Request(request_value, method="GET", headers={"Range": "bytes=0-0", "User-Agent": "siduri-linkcheck/1"})
                            with urllib.request.urlopen(request, timeout=timeout) as response:
                                if response.status >= 400:
                                    raise urllib.error.HTTPError(request_value, response.status, "HTTP error", response.headers, None)
                        except Exception as fallback_exc:  # pragma: no cover - network-dependent
                            errors.append(f"{page}: external link {value!r}: {fallback_exc}")
                    else:
                        errors.append(f"{page}: external link {value!r}: {exc}")
                except Exception as exc:  # pragma: no cover - network-dependent
                    errors.append(f"{page}: external link {value!r}: {exc}")
                continue

            path = parsed.path or "/"
            if path in KNOWN_PENDING:
                reason = KNOWN_PENDING[path]
                target = internal_target(dist, page, value)
                if target is not None and target.exists() and target.is_file():
                    errors.append(
                        f"{page}: known pending route {path!r} now resolves to {target}; "
                        f"remove the stale allowlist entry ({reason})"
                    )
                else:
                    print(f"linkcheck: known pending route skipped: {path} ({reason})")
                continue
            internal_count += 1
            target = internal_target(dist, page, value)
            if target is None or not target.exists() or not target.is_file():
                errors.append(f"{page}: internal link {value!r} resolves to missing {target}")

    mode = "including external HTTP checks" if external else "internal only; external checks skipped (use --external)"
    print(f"linkcheck: {internal_count} internal links, {external_count} external links, {len(errors)} failures, {mode}")
    return internal_count, external_count, errors


def selftest() -> int:
    with tempfile.TemporaryDirectory(prefix="siduri-linkcheck-") as raw:
        dist = pathlib.Path(raw)
        (dist / "present").mkdir()
        present_page = dist / "present/index.html"
        present_page.write_text('<a href="/present/">present</a><a href="/services/">pending</a>', encoding="utf-8")
        (dist / "index.html").write_text('<a href="/present/">present</a>', encoding="utf-8")
        _, external, errors = check_dist(dist, external=False)
        assert external == 0 and not errors, errors
        (dist / "services").mkdir()
        (dist / "services/index.html").write_text("services", encoding="utf-8")
        _, _, errors = check_dist(dist, external=False)
        assert any("known pending route '/services/' now resolves" in error for error in errors), errors
        present_page.write_text('<a href="/present/">present</a>', encoding="utf-8")
        (dist / "index.html").write_text('<a href="/dead/">dead</a>', encoding="utf-8")
        _, _, errors = check_dist(dist, external=False)
        assert any("internal link '/dead/'" in error for error in errors), errors
    print("linkcheck selftest: internal success, expiring allowlist, dead-link failure, and external opt-in pass")
    return 0


def main(argv: list[str]) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--dist", type=pathlib.Path, default=pathlib.Path("dist"))
    parser.add_argument("--external", action="store_true", help="also request HTTP(S) links")
    parser.add_argument("--timeout", type=float, default=10.0)
    parser.add_argument("--selftest", action="store_true")
    args = parser.parse_args(argv)
    if args.selftest:
        return selftest()
    _, _, errors = check_dist(args.dist, args.external, args.timeout)
    if errors:
        print("Link failures:", file=__import__("sys").stderr)
        for error in errors:
            print(f"  {error}", file=__import__("sys").stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main(__import__("sys").argv[1:]))
