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
import re
import tempfile
import urllib.error
import urllib.parse
import urllib.request


KNOWN_PENDING = {
    "/services/": "the current article fixture names the later services surface before its lane exists",
}

REPO_ROOT = pathlib.Path(__file__).resolve().parent.parent
SITE_BASE_SOURCE = REPO_ROOT / "internal" / "site" / "meta_a3.go"
SITE_BASE_DECLARATION = re.compile(r'^\s*const\s+siteBaseURL\s*=\s*"([^"]+)"\s*$', re.MULTILINE)


def read_site_base_url(source: pathlib.Path = SITE_BASE_SOURCE) -> str:
    """Read the one site-origin declaration used by the Go renderer."""
    try:
        text = source.read_text(encoding="utf-8")
    except OSError as exc:
        raise ValueError(f"cannot read site base URL from {source}: {exc}") from exc
    matches = SITE_BASE_DECLARATION.findall(text)
    if len(matches) != 1:
        raise ValueError(f"expected one siteBaseURL declaration in {source}, found {len(matches)}")
    value = matches[0]
    parsed = urllib.parse.urlsplit(value)
    if parsed.scheme not in {"http", "https"} or not parsed.netloc:
        raise ValueError(f"siteBaseURL must be an absolute HTTP(S) URL, got {value!r}")
    if parsed.path not in {"", "/"} or parsed.query or parsed.fragment:
        raise ValueError(f"siteBaseURL must identify an origin, got {value!r}")
    try:
        if parsed.hostname is None:
            raise ValueError("missing hostname")
        _ = parsed.port
    except ValueError as exc:
        raise ValueError(f"siteBaseURL has an invalid origin, got {value!r}: {exc}") from exc
    return value


def normalized_origin(value: str, protocol_scheme: str | None = None) -> tuple[str, str, int | None] | None:
    parsed = urllib.parse.urlsplit(value)
    if not parsed.netloc:
        return None
    scheme = (parsed.scheme or protocol_scheme or "").lower()
    if not scheme:
        return None
    try:
        hostname = parsed.hostname
        port = parsed.port
    except ValueError:
        return None
    if hostname is None:
        return None
    hostname = hostname.lower().rstrip(".")
    if port is None and scheme in {"http", "https"}:
        port = 80 if scheme == "http" else 443
    return scheme, hostname, port


def origin_text(origin: tuple[str, str, int | None]) -> str:
    scheme, hostname, port = origin
    default_port = (scheme == "http" and port == 80) or (scheme == "https" and port == 443)
    suffix = "" if default_port or port is None else f":{port}"
    return f"{scheme}://{hostname}{suffix}"


def is_same_origin(value: str, site_origin: tuple[str, str, int | None]) -> bool:
    return normalized_origin(value, protocol_scheme=site_origin[0]) == site_origin


class LinkParser(html.parser.HTMLParser):
    def __init__(self) -> None:
        super().__init__(convert_charrefs=True)
        self.links: list[str] = []
        self.site_metadata_links: list[tuple[str, str]] = []

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        attrs_by_name = {name.lower(): value for name, value in attrs}
        tag = tag.lower()
        if tag in {"a", "area"} and attrs_by_name.get("href"):
            self.links.append(attrs_by_name["href"] or "")
        if tag == "link" and attrs_by_name.get("href"):
            rel = {part.lower() for part in (attrs_by_name.get("rel") or "").split()}
            if rel & {"alternate", "icon", "modulepreload", "preload", "stylesheet"}:
                self.links.append(attrs_by_name["href"] or "")
            for relation in ("canonical", "alternate"):
                if relation in rel:
                    self.site_metadata_links.append((relation, attrs_by_name["href"] or ""))
        if tag in {"script", "img", "iframe", "video", "audio", "source", "object", "embed"}:
            for name in ("src", "data"):
                if attrs_by_name.get(name):
                    self.links.append(attrs_by_name[name] or "")


def html_files(dist: pathlib.Path) -> list[pathlib.Path]:
    return sorted(dist.rglob("*.html")) if dist.is_dir() else []


def internal_target(
    dist: pathlib.Path,
    page: pathlib.Path,
    value: str,
    site_origin: tuple[str, str, int | None] | None = None,
) -> pathlib.Path | None:
    parsed = urllib.parse.urlsplit(value)
    if parsed.scheme or parsed.netloc:
        if site_origin is None or not is_same_origin(value, site_origin):
            return None
        path = urllib.parse.unquote(parsed.path or "/")
    else:
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


def check_dist(
    dist: pathlib.Path,
    external: bool,
    timeout: float = 10.0,
    base_source: pathlib.Path = SITE_BASE_SOURCE,
    open_url=None,
) -> tuple[int, int, list[str]]:
    pages = html_files(dist)
    if not pages:
        return 0, 0, [f"no HTML files found under {dist}"]
    try:
        site_base_url = read_site_base_url(base_source)
        site_origin = normalized_origin(site_base_url)
    except ValueError as exc:
        return 0, 0, [f"site base URL configuration error: {exc}"]
    if site_origin is None:
        return 0, 0, [f"site base URL configuration error: invalid origin {site_base_url!r}"]
    if open_url is None:
        open_url = urllib.request.urlopen
    errors: list[str] = []
    internal_count = 0
    external_count = 0
    for page in pages:
        parser = LinkParser()
        parser.feed(page.read_text(encoding="utf-8"))
        for relation, metadata_value in parser.site_metadata_links:
            metadata_origin = normalized_origin(metadata_value, protocol_scheme=site_origin[0])
            if metadata_origin is not None and metadata_origin != site_origin:
                errors.append(
                    f"{page}: site base URL mismatch: {relation} link {metadata_value!r} has origin "
                    f"{origin_text(metadata_origin)!r}, expected {origin_text(site_origin)!r} "
                    f"from {base_source}"
                )
        for raw in parser.links:
            value = raw.strip()
            if not value or value.startswith(("#", "mailto:", "tel:", "javascript:", "data:")):
                continue
            parsed = urllib.parse.urlsplit(value)
            if parsed.scheme in {"http", "https"} or parsed.netloc or value.startswith("//"):
                if is_same_origin(value, site_origin):
                    internal_count += 1
                    target = internal_target(dist, page, value, site_origin)
                    if target is None or not target.exists() or not target.is_file():
                        errors.append(f"{page}: internal link {value!r} resolves to missing {target}")
                    continue
                external_count += 1
                if not external:
                    continue
                request_value = value if parsed.scheme else f"https:{value}"
                try:
                    request = urllib.request.Request(request_value, method="HEAD", headers={"User-Agent": "siduri-linkcheck/1"})
                    with open_url(request, timeout=timeout) as response:
                        if response.status >= 400:
                            raise urllib.error.HTTPError(request_value, response.status, "HTTP error", response.headers, None)
                except urllib.error.HTTPError as exc:
                    if exc.code in {405, 501}:
                        try:
                            request = urllib.request.Request(request_value, method="GET", headers={"Range": "bytes=0-0", "User-Agent": "siduri-linkcheck/1"})
                            with open_url(request, timeout=timeout) as response:
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
        base_source = dist / "meta_a3.go"
        base_source.write_text('const siteBaseURL = "https://siduri.test"\n', encoding="utf-8")
        (dist / "present").mkdir()
        present_page = dist / "present/index.html"
        present_page.write_text(
            '<link rel="canonical" href="https://siduri.test/present/">'
            '<a href="https://siduri.test/present/">present</a>'
            '<a href="/services/">pending</a>',
            encoding="utf-8",
        )
        (dist / "index.html").write_text(
            '<link rel="canonical" href="https://siduri.test/">'
            '<a href="https://siduri.test/present/">present</a>',
            encoding="utf-8",
        )
        internal, external, errors = check_dist(dist, external=False, base_source=base_source)
        assert internal == 2 and external == 0 and not errors, errors
        (dist / "services").mkdir()
        (dist / "services/index.html").write_text("services", encoding="utf-8")
        _, _, errors = check_dist(dist, external=False, base_source=base_source)
        assert any("known pending route '/services/' now resolves" in error for error in errors), errors
        (dist / "services/index.html").unlink()
        present_page.write_text(
            '<link rel="canonical" href="https://siduri.test/present/">'
            '<a href="https://siduri.test/present/">present</a>'
            '<a href="https://siduri.test/nope/">dead self-link</a>'
            '<a href="https://example.invalid/nope">external</a>',
            encoding="utf-8",
        )
        internal, external, errors = check_dist(dist, external=False, base_source=base_source)
        assert internal == 3 and external == 1, (internal, external, errors)
        assert any("internal link 'https://siduri.test/nope/'" in error for error in errors), errors
        assert not any("external link 'https://example.invalid/nope'" in error for error in errors), errors

        def fail_external(request, timeout):
            raise urllib.error.URLError("selftest external failure")

        _, _, errors = check_dist(dist, external=True, base_source=base_source, open_url=fail_external)
        assert any("external link 'https://example.invalid/nope'" in error for error in errors), errors

        base_source.write_text('const siteBaseURL = "https://different.test"\n', encoding="utf-8")
        _, _, errors = check_dist(dist, external=False, base_source=base_source)
        assert any("site base URL mismatch" in error for error in errors), errors
    print("linkcheck selftest: same-origin internal links, broken self-link, external failure, and source/build divergence pass")
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
