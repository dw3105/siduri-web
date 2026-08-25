#!/usr/bin/env python3
"""Find credentials and raw email addresses before they enter the tree.

    python3 tools/secretscan.py .
    python3 tools/secretscan.py --pre-commit
    python3 tools/secretscan.py --install-hook
    python3 tools/secretscan.py --selftest

The CI invocation scans the tree.  The optional pre-commit hook scans staged
paths, because hooks are not cloned; ``make install-hooks`` installs it into
the current checkout.  Raw emails in the two legal Markdown files are
intentionally scoped out: those pages must carry the operator's legal contact
address.  The public contact template is likewise an explicit, reviewed
allowlist.  A raw email anywhere else, including another ``content/`` file,
is a failure.  Lines containing ``email_hash`` are excluded from both pattern
families because hashes are stored there by design, not credentials.
"""
from __future__ import annotations

import argparse
import pathlib
import re
import stat
import subprocess
import sys
import tempfile


EMAIL = re.compile(r"(?<![A-Za-z0-9._%+-])([A-Za-z0-9][A-Za-z0-9._%+-]*@[A-Za-z0-9.-]+\.[A-Za-z]{2,})(?![A-Za-z0-9.-])")
SECRET_PATTERNS = (
    ("AWS access key", re.compile(r"\bAKIA[0-9A-Z]{16}\b")),
    ("GitHub token", re.compile(r"\b(?:gh[pousr]_[A-Za-z0-9_]{20,}|github_pat_[A-Za-z0-9_]{20,})\b")),
    ("Slack token", re.compile(r"\bxox[baprs]-[0-9A-Za-z-]{20,}\b")),
    ("private key", re.compile(r"-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----")),
    ("secret assignment", re.compile(r"(?i)\b(?:api[_-]?key|secret|password|token)\s*[:=]\s*[\"']?[A-Za-z0-9/+_=-]{16,}")),
)
IGNORED_DIRS = {".git", "dist", ".dev-dist", ".preview-dist", "node_modules", ".wrangler", "__pycache__"}
LEGAL_EMAIL_FILES = {"content/legal/impressum.md", "content/legal/datenschutz.md"}
ALLOWED_CONTACT = {"hello" + "@siduri.ai"}
#: A8 guessed these filenames before A5 existed and named two files no lane
#: ever created; A5 put every page in one template. The allowlist stays
#: file-scoped rather than pattern-scoped deliberately -- a pattern would
#: permit the address anywhere it matched.
ALLOWED_CONTACT_FILES = {
    "internal/site/contact.templ",
    "internal/site/contact_templ.go",
    "internal/site/pages_a5.templ",
    "internal/site/pages_a5_templ.go",
}
EMAIL_DOCUMENTATION_PREFIX = "docs/"


def relative_name(path: pathlib.Path) -> str:
    try:
        return path.resolve().relative_to(pathlib.Path.cwd().resolve()).as_posix()
    except ValueError:
        return path.as_posix()


def is_legal_email_file(name: str) -> bool:
    return name in LEGAL_EMAIL_FILES or any(name.endswith(f"/{item}") for item in LEGAL_EMAIL_FILES)


def is_allowed_contact_file(name: str) -> bool:
    return name in ALLOWED_CONTACT_FILES or any(name.endswith(f"/{item}") for item in ALLOWED_CONTACT_FILES)


def files_under(root: pathlib.Path) -> list[pathlib.Path]:
    if root.is_file():
        return [root]
    files: list[pathlib.Path] = []
    for path in root.rglob("*"):
        if not path.is_file() or any(part in IGNORED_DIRS for part in path.parts):
            continue
        files.append(path)
    return sorted(files)


def staged_files() -> list[pathlib.Path]:
    result = subprocess.run(
        ["git", "diff", "--cached", "--name-only", "--diff-filter=ACMR", "-z"],
        check=True,
        capture_output=True,
    )
    return [pathlib.Path(raw.decode("utf-8")) for raw in result.stdout.split(b"\0") if raw]


def scan(paths: list[pathlib.Path]) -> list[str]:
    findings: list[str] = []
    for path in paths:
        name = relative_name(path)
        if path.name in {".env", ".env.local", ".env.production"} or path.suffix in {".pem", ".key"}:
            findings.append(f"{name}: secret-bearing filename")
        try:
            text = path.read_text(encoding="utf-8")
        except (UnicodeDecodeError, OSError):
            continue
        for line_number, line in enumerate(text.splitlines(), 1):
            if "email_hash" in line.lower():
                continue
            for label, pattern in SECRET_PATTERNS:
                if pattern.search(line):
                    findings.append(f"{name}:{line_number}: {label}")
            for match in EMAIL.finditer(line):
                address = match.group(1)
                # Requirements and decision records contain illustrative email
                # addresses. They are contract documentation, not content or
                # an artifact that can be served, so only their emails are
                # excluded; secrets in docs are still scanned.
                if name.startswith(EMAIL_DOCUMENTATION_PREFIX) or is_legal_email_file(name):
                    continue
                if is_allowed_contact_file(name) and address in ALLOWED_CONTACT:
                    continue
                findings.append(f"{name}:{line_number}: raw email address ({address})")
    return findings


def hook_text() -> str:
    return "#!/bin/sh\nset -eu\nexec python3 tools/secretscan.py --pre-commit\n"


def install_hook() -> int:
    result = subprocess.run(["git", "rev-parse", "--git-path", "hooks"], check=True, capture_output=True, text=True)
    hooks = pathlib.Path(result.stdout.strip())
    hooks.mkdir(parents=True, exist_ok=True)
    hook = hooks / "pre-commit"
    hook.write_text(hook_text(), encoding="utf-8")
    hook.chmod(hook.stat().st_mode | stat.S_IXUSR | stat.S_IXGRP | stat.S_IXOTH)
    print(f"secretscan: installed {hook}")
    return 0


def selftest() -> int:
    with tempfile.TemporaryDirectory(prefix="siduri-secrets-") as raw:
        root = pathlib.Path(raw)
        (root / "content").mkdir()
        (root / "content/clean.md").write_text("email_hash: 0123456789abcdef\n", encoding="utf-8")
        (root / "content/legal").mkdir()
        (root / "content/legal/impressum.md").write_text("operator" + "@example.com\n", encoding="utf-8")
        assert not scan(files_under(root))
        (root / "content/leak.md").write_text("token: " + "AK" + "IA1234567890ABCDEF\n" + "leak" + "@example.com\n", encoding="utf-8")
        findings = scan(files_under(root))
        assert any("AWS access key" in finding for finding in findings), findings
        assert any("raw email address" in finding for finding in findings), findings
        (root / "content/leak.md").unlink()
        assert not scan(files_under(root))
    assert "--pre-commit" in hook_text()
    print("secretscan selftest: email_hash/legal allowlist, planted AKIA key and content email, and clean-up pass")
    return 0


def main(argv: list[str]) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("paths", nargs="*", type=pathlib.Path, default=[pathlib.Path(".")])
    parser.add_argument("--pre-commit", action="store_true")
    parser.add_argument("--install-hook", action="store_true")
    parser.add_argument("--selftest", action="store_true")
    args = parser.parse_args(argv)
    if args.selftest:
        return selftest()
    if args.install_hook:
        return install_hook()
    paths = staged_files() if args.pre_commit else [path for arg in args.paths for path in files_under(arg)]
    findings = scan(paths)
    if findings:
        print("secretscan: prohibited material found", file=sys.stderr)
        for finding in findings:
            print(f"  {finding}", file=sys.stderr)
        return 1
    print(f"secretscan: scanned {len(paths)} file(s); no secrets or unsanctioned raw emails")
    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
