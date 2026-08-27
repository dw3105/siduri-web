#!/usr/bin/env python3
"""Check that a requirements-file amendment carries an auditable record.

    python3 tools/amendcheck.py BASE [HEAD]
    python3 tools/amendcheck.py --selftest

The checker is deliberately a shape check.  It cannot establish who approved
an amendment; it establishes that an added ADR records it and that changed
contract text remains struck in place.

Stdlib only.  Exit 1 for an unexplained or malformed amendment and exit 0 for
an accepted amendment or no requirements-file diff.
"""
from __future__ import annotations

import argparse
import collections
import dataclasses
import datetime
import pathlib
import re
import subprocess
import sys
import tempfile
from typing import Iterable


CONTRACT_FILES = (
    "docs/site-requirements.md",
    "docs/comments-requirements.md",
)
DIFF_CONTEXT = "-U3"
ADR_PATH = re.compile(r"^docs/adr/(?P<number>\d{4})-[^/]+\.md$")
HEADING = re.compile(r"^(?P<marks>#{1,6})\s+(?P<text>.+?)\s*$")
STATUS_PREFIX = re.compile(r"^\s*\*\*Status\*\*\s*·\s*(?P<body>.+?)\s*$")
DATE = re.compile(r"(?<!\d)(?P<date>\d{4}-\d{2}-\d{2})(?!\d)")
STRIKE = re.compile(r"~~(.*?)~~")


class AmendcheckError(Exception):
    """The checker could not inspect the requested Git range."""


@dataclasses.dataclass
class DiffLine:
    kind: str
    text: str
    old_line: int | None
    new_line: int | None


@dataclasses.dataclass
class Hunk:
    old_start: int
    new_start: int
    lines: list[DiffLine] = dataclasses.field(default_factory=list)

    @property
    def removed(self) -> list[DiffLine]:
        return [line for line in self.lines if line.kind == "-"]

    @property
    def added(self) -> list[DiffLine]:
        return [line for line in self.lines if line.kind == "+"]


@dataclasses.dataclass
class FileDiff:
    path: str
    changed: bool = False
    hunks: list[Hunk] = dataclasses.field(default_factory=list)


@dataclasses.dataclass
class Heading:
    line_number: int
    level: int
    text: str
    name: str


@dataclasses.dataclass
class AddedADR:
    path: str
    number: str
    text: str
    status_line: str | None
    date_error: str | None


def git_text(repo: pathlib.Path, arguments: list[str]) -> str:
    result = subprocess.run(
        ["git", *arguments],
        cwd=repo,
        capture_output=True,
        check=False,
    )
    if result.returncode:
        detail = result.stderr.decode("utf-8", errors="replace").strip()
        command = "git " + " ".join(arguments)
        raise AmendcheckError(f"{command} failed: {detail or 'exit ' + str(result.returncode)}")
    return result.stdout.decode("utf-8", errors="replace")


def git_bytes(repo: pathlib.Path, arguments: list[str]) -> bytes:
    result = subprocess.run(
        ["git", *arguments],
        cwd=repo,
        capture_output=True,
        check=False,
    )
    if result.returncode:
        detail = result.stderr.decode("utf-8", errors="replace").strip()
        command = "git " + " ".join(arguments)
        raise AmendcheckError(f"{command} failed: {detail or 'exit ' + str(result.returncode)}")
    return result.stdout


def repository_root() -> pathlib.Path:
    result = subprocess.run(
        ["git", "rev-parse", "--show-toplevel"],
        capture_output=True,
        check=False,
        text=True,
    )
    if result.returncode:
        raise AmendcheckError("amendcheck must run inside a Git repository")
    return pathlib.Path(result.stdout.strip()).resolve()


def parse_hunk_header(line: str) -> tuple[int, int]:
    match = re.match(r"^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@", line)
    if not match:
        raise AmendcheckError(f"cannot parse Git hunk header: {line}")
    old_start = int(match.group(1))
    new_start = int(match.group(3))
    return old_start, new_start


def parse_contract_diff(diff_text: str) -> dict[str, FileDiff]:
    files = {path: FileDiff(path) for path in CONTRACT_FILES}
    current: FileDiff | None = None
    hunk: Hunk | None = None
    old_line = 0
    new_line = 0

    for raw_line in diff_text.splitlines():
        if raw_line.startswith("diff --git "):
            match = re.match(r"^diff --git a/(\S+) b/(\S+)$", raw_line)
            current = None
            hunk = None
            if match:
                path = match.group(2) if match.group(2) != "/dev/null" else match.group(1)
                current = files.get(path)
                if current is not None:
                    current.changed = True
            continue
        if current is None:
            continue
        if raw_line.startswith("@@ "):
            old_start, new_start = parse_hunk_header(raw_line)
            hunk = Hunk(old_start, new_start)
            current.hunks.append(hunk)
            old_line = old_start
            new_line = new_start
            continue
        if hunk is None or not raw_line:
            continue
        marker = raw_line[0]
        if marker == " ":
            hunk.lines.append(DiffLine(" ", raw_line[1:], old_line, new_line))
            old_line += 1
            new_line += 1
        elif marker == "-":
            hunk.lines.append(DiffLine("-", raw_line[1:], old_line, None))
            old_line += 1
        elif marker == "+":
            hunk.lines.append(DiffLine("+", raw_line[1:], old_line, new_line))
            new_line += 1
        # Git's "\\ No newline at end of file" marker is intentionally ignored.

    return files


def read_revision_file(repo: pathlib.Path, revision: str, path: str) -> str:
    return git_bytes(repo, ["show", f"{revision}:{path}"]).decode("utf-8", errors="replace")


def added_adr_paths(repo: pathlib.Path, base: str, head: str) -> list[str]:
    output = git_bytes(
        repo,
        ["diff", "--name-only", "--diff-filter=A", "-z", base, head, "--", "docs/adr"],
    )
    return [item.decode("utf-8", errors="replace") for item in output.split(b"\0") if item]


def existing_adr_numbers(repo: pathlib.Path, base: str) -> set[str]:
    paths = git_text(repo, ["ls-tree", "-r", "--name-only", base, "--", "docs/adr"]).splitlines()
    numbers: set[str] = set()
    for path in paths:
        match = ADR_PATH.match(path)
        if match:
            numbers.add(match.group("number"))
    return numbers


def status_line(text: str) -> tuple[str | None, str | None]:
    for line in text.splitlines():
        prefix = STATUS_PREFIX.match(line)
        if prefix is None:
            continue
        date_match = DATE.search(prefix.group("body"))
        if date_match is None:
            return line.strip(), "missing an ISO date"
        try:
            datetime.date.fromisoformat(date_match.group("date"))
        except ValueError:
            return line.strip(), f"invalid date {date_match.group('date')}"
        status = prefix.group("body")[: date_match.start()].strip(" ·—-:*()")
        if not re.search(r"[A-Za-z]", status):
            return line.strip(), "missing a status before the date"
        return line.strip(), None
    return None, "missing a **Status** · line"


def inspect_added_adrs(repo: pathlib.Path, base: str, head: str) -> list[AddedADR]:
    taken = existing_adr_numbers(repo, base)
    result: list[AddedADR] = []
    for path in added_adr_paths(repo, base, head):
        match = ADR_PATH.match(path)
        if match is None:
            continue
        text = read_revision_file(repo, head, path)
        line, date_error = status_line(text)
        number = match.group("number")
        if number in taken:
            date_error = f"ADR number {number} is already taken"
        result.append(AddedADR(path, number, text, line, date_error))

    counts = collections.Counter(item.number for item in result)
    for item in result:
        if counts[item.number] > 1:
            item.date_error = f"ADR number {item.number} is duplicated in the added range"
    return result


def normalise_whitespace(value: str) -> str:
    return re.sub(r"\s+", " ", value).strip()


def strike_spans(lines: Iterable[DiffLine]) -> list[str]:
    spans: list[str] = []
    for line in lines:
        for match in STRIKE.finditer(line.text):
            content = normalise_whitespace(match.group(1))
            if content:
                spans.append(content)
    return spans


def injective_matches(removed: list[str], spans: list[str]) -> int:
    """Return a maximum matching of removed lines to added strike spans."""
    owners = [-1] * len(spans)
    edges = [
        [span_index for span_index, span in enumerate(spans) if span in removed_line]
        for removed_line in removed
    ]

    def visit(removed_index: int, seen: set[int]) -> bool:
        for span_index in edges[removed_index]:
            if span_index in seen:
                continue
            seen.add(span_index)
            if owners[span_index] == -1 or visit(owners[span_index], seen):
                owners[span_index] = removed_index
                return True
        return False

    return sum(visit(index, set()) for index in range(len(removed)))


def check_hunk(hunk: Hunk, number: int) -> tuple[bool, str]:
    removed = [normalise_whitespace(line.text) for line in hunk.removed]
    spans = strike_spans(hunk.added)
    if not removed:
        return True, f"hunk {number}: PASS — removed 0 line(s); no trace required"
    matches = injective_matches(removed, spans)
    if len(spans) < len(removed):
        return (
            False,
            f"hunk {number}: FAIL — removed {len(removed)} line(s), added "
            f"{len(spans)} non-empty strike span(s); injectivity count is too small",
        )
    if matches < len(removed):
        return (
            False,
            f"hunk {number}: FAIL — only {matches}/{len(removed)} removed line(s) "
            "have an injective struck substring in the same hunk",
        )
    return (
        True,
        f"hunk {number}: PASS — {len(removed)} removed line(s), {len(spans)} "
        f"strike span(s), {matches} injective match(es)",
    )


def heading_name(text: str) -> str | None:
    match = HEADING.match(text)
    if match is None:
        return None
    heading_text = match.group("text").strip()
    clause = re.match(r"^(?P<id>[A-Za-z][A-Za-z0-9-]*)\s+·\s+", heading_text)
    return clause.group("id") if clause else heading_text


def headings(text: str) -> list[Heading]:
    result: list[Heading] = []
    for line_number, line in enumerate(text.splitlines(), 1):
        match = HEADING.match(line)
        if match:
            text_value = match.group("text").strip()
            result.append(
                Heading(
                    line_number,
                    len(match.group("marks")),
                    text_value,
                    heading_name(line) or text_value,
                )
            )
    return result


def containing_heading(all_headings: list[Heading], line_number: int) -> str | None:
    candidates = [heading for heading in all_headings if heading.line_number < line_number]
    return candidates[-1].name if candidates else None


def touched_heading_names(file_diff: FileDiff, base_text: str) -> set[str]:
    all_headings = headings(base_text)
    names: set[str] = set()
    for hunk in file_diff.hunks:
        for line in hunk.lines:
            if line.kind not in {"-", "+"}:
                continue
            direct = heading_name(line.text)
            if direct is not None:
                names.add(direct)
            elif line.old_line is not None:
                containing = containing_heading(all_headings, line.old_line)
                if containing is not None:
                    names.add(containing)
    return names


def rule1_message(
    has_contract_diff: bool,
    added_adrs: list[AddedADR],
) -> tuple[bool, str]:
    if not has_contract_diff:
        return True, "no contract diff; no amendment record is required"
    valid = [item for item in added_adrs if item.date_error is None]
    if valid:
        details = "; ".join(f"{item.path} ({item.status_line})" for item in valid)
        return True, f"added ADR with a fresh number and dated status: {details}"
    if not added_adrs:
        return False, "no added docs/adr/NNNN-*.md file in the range"
    details = "; ".join(
        f"{item.path}: {item.date_error or 'invalid'}" for item in added_adrs
    )
    return False, f"no added ADR satisfies the number and status checks: {details}"


def rule2_messages(
    has_contract_diff: bool,
    file_diffs: dict[str, FileDiff],
) -> tuple[bool, list[str]]:
    messages: list[str] = []
    passed = True
    for path in CONTRACT_FILES:
        file_diff = file_diffs[path]
        if not file_diff.changed:
            messages.append(f"{path}: PASS — no diff")
            continue
        if not file_diff.hunks:
            messages.append(f"{path}: FAIL — changed without textual hunks")
            passed = False
            continue
        file_passed = True
        for number, hunk in enumerate(file_diff.hunks, 1):
            hunk_passed, message = check_hunk(hunk, number)
            messages.append(f"{path}: {message}")
            file_passed = file_passed and hunk_passed
        passed = passed and file_passed
    if not has_contract_diff:
        return True, messages
    return passed, messages


def rule3_messages(
    has_contract_diff: bool,
    file_diffs: dict[str, FileDiff],
    base_texts: dict[str, str],
    added_adrs: list[AddedADR],
) -> tuple[bool, list[str]]:
    messages: list[str] = []
    adr_text = "\n".join(item.text for item in added_adrs)
    passed = True
    for path in CONTRACT_FILES:
        file_diff = file_diffs[path]
        if not file_diff.changed:
            messages.append(f"{path}: PASS — no diff")
            continue
        names = sorted(touched_heading_names(file_diff, base_texts[path]))
        if not names:
            messages.append(f"{path}: PASS — no heading section was touched")
            continue
        missing = [name for name in names if name not in adr_text]
        if missing:
            messages.append(
                f"{path}: FAIL — touched section name(s) missing from added ADR text: "
                + ", ".join(missing)
            )
            passed = False
        else:
            messages.append(f"{path}: PASS — added ADR names touched section(s): {', '.join(names)}")
    if not has_contract_diff:
        return True, messages
    return passed, messages


def assess(repo: pathlib.Path, base: str, head: str) -> tuple[bool, str]:
    diff_text = git_text(
        repo,
        [
            "diff",
            "--no-ext-diff",
            "--no-color",
            DIFF_CONTEXT,
            base,
            head,
            "--",
            *CONTRACT_FILES,
        ],
    )
    file_diffs = parse_contract_diff(diff_text)
    has_contract_diff = any(file_diff.changed for file_diff in file_diffs.values())
    added_adrs = inspect_added_adrs(repo, base, head) if has_contract_diff else []
    base_texts = {
        path: read_revision_file(repo, base, path)
        for path in CONTRACT_FILES
        if file_diffs[path].changed
    }
    one_passed, one_detail = rule1_message(has_contract_diff, added_adrs)
    two_passed, two_details = rule2_messages(has_contract_diff, file_diffs)
    three_passed, three_details = rule3_messages(
        has_contract_diff,
        file_diffs,
        base_texts,
        added_adrs,
    )

    output: list[str] = [f"amendcheck: base={base} head={head} diff-context=3"]
    output.append(f"rule 1 (added ADR): {'PASS' if one_passed else 'FAIL'} — {one_detail}")
    output.append(f"rule 2 (struck traces): {'PASS' if two_passed else 'FAIL'}")
    output.extend(f"  {detail}" for detail in two_details)
    output.append(f"rule 3 (named sections): {'PASS' if three_passed else 'FAIL'}")
    output.extend(f"  {detail}" for detail in three_details)
    all_passed = one_passed and two_passed and three_passed
    if not has_contract_diff:
        output.append("decision: no contract diff; accepted")
    elif all_passed:
        output.append("decision: legitimate amendment; all three rules pass")
    else:
        failed = ", ".join(
            name
            for name, passed in (("rule 1", one_passed), ("rule 2", two_passed), ("rule 3", three_passed))
            if not passed
        )
        output.append(f"decision: rejected amendment; failed {failed}")
    return all_passed, "\n".join(output)


def run_git_fixture(repo: pathlib.Path, arguments: list[str]) -> str:
    result = subprocess.run(
        ["git", *arguments],
        cwd=repo,
        capture_output=True,
        check=True,
        text=True,
    )
    return result.stdout.strip()


FIXTURE_SITE = """# Website — Requirements

## 12 · Phasing

### P0 · Foundation

- [ ] Foundation clause remains visible.
- [ ] Bulk line 01 remains visible.
- [ ] Bulk line 02 remains visible.
- [ ] Bulk line 03 remains visible.
- [ ] Bulk line 04 remains visible.
- [ ] Bulk line 05 remains visible.
- [ ] Bulk line 06 remains visible.
- [ ] Bulk line 07 remains visible.
- [ ] Bulk line 08 remains visible.
- [ ] Bulk line 09 remains visible.
- [ ] Bulk line 10 remains visible.

keep A 01
keep A 02
keep A 03
keep A 04
keep A 05
keep A 06
keep A 07
keep A 08

- [ ] Silent line 01 remains visible.
- [ ] Silent line 02 remains visible.
- [ ] Silent line 03 remains visible.
- [ ] Silent line 04 remains visible.
- [ ] Silent line 05 remains visible.
- [ ] Silent line 06 remains visible.
- [ ] Silent line 07 remains visible.
- [ ] Silent line 08 remains visible.
- [ ] Silent line 09 remains visible.
- [ ] Silent line 10 remains visible.

## 13 · Acceptance criteria

- [ ] Criterion under a heading without a clause id.

#### OQ-5 · Domain and name

| Norse | Huginn | Existing name | Rejected |
"""

FIXTURE_COMMENTS = """# Comment System — Requirements

## 11 · Open questions

1. **Comment freeze.**
"""


def make_fixture(root: pathlib.Path) -> str:
    (root / "docs/adr").mkdir(parents=True)
    (root / "docs/site-requirements.md").write_text(FIXTURE_SITE, encoding="utf-8")
    (root / "docs/comments-requirements.md").write_text(FIXTURE_COMMENTS, encoding="utf-8")
    (root / "docs/adr/0001-existing.md").write_text(
        "# 0001 · Existing\n\n**Status** · accepted · 2026-08-01\n",
        encoding="utf-8",
    )
    run_git_fixture(root, ["init", "-q"])
    run_git_fixture(root, ["config", "user.email", "selftest@example.invalid"])
    run_git_fixture(root, ["config", "user.name", "amendcheck selftest"])
    run_git_fixture(root, ["add", "."])
    run_git_fixture(root, ["commit", "-qm", "fixture base"])
    return run_git_fixture(root, ["rev-parse", "HEAD"])


def add_adr(root: pathlib.Path, number: str, name: str, mentions: str) -> None:
    (root / f"docs/adr/{number}-{name}.md").write_text(
        f"# {number} · {name}\n\n**Status** · accepted · 2026-08-27\n\n"
        f"## What changed\n\nThe amendment names {mentions}.\n",
        encoding="utf-8",
    )


def replace_once(path: pathlib.Path, old: str, new: str) -> None:
    text = path.read_text(encoding="utf-8")
    assert text.count(old) == 1, f"fixture text did not contain one {old!r}"
    path.write_text(text.replace(old, new), encoding="utf-8")


def run_selftest_case(
    name: str,
    mutate,
    expected_code: int,
    expected_marker: str,
) -> None:
    with tempfile.TemporaryDirectory(prefix="siduri-amendcheck-") as raw:
        root = pathlib.Path(raw)
        base = make_fixture(root)
        mutate(root)
        run_git_fixture(root, ["add", "."])
        staged_changes = subprocess.run(
            ["git", "diff", "--cached", "--quiet"],
            cwd=root,
            check=False,
        )
        if staged_changes.returncode:
            run_git_fixture(root, ["commit", "-qm", name])
        head = run_git_fixture(root, ["rev-parse", "HEAD"])
        result = subprocess.run(
            [sys.executable, str(pathlib.Path(__file__).resolve()), base, head],
            cwd=root,
            capture_output=True,
            text=True,
            check=False,
        )
        combined = result.stdout + result.stderr
        assert result.returncode == expected_code, combined
        assert expected_marker in combined, combined
        selected = [
            line
            for line in result.stdout.splitlines()
            if line.startswith("rule ") or line.startswith("decision:") or "hunk 2:" in line
        ]
        print(f"selftest {name}: exit {result.returncode}")
        for line in selected:
            print(f"  {line}")


def selftest() -> int:
    # Red proofs come first: rule 2, then rule 1, then rule 3.
    run_selftest_case(
        "delete-outright",
        lambda root: (
            replace_once(
                root / "docs/site-requirements.md",
                "- [ ] Foundation clause remains visible.\n",
                "",
            ),
            add_adr(root, "0002", "foundation-deletion", "P0"),
        ),
        1,
        "rule 2 (struck traces): FAIL",
    )
    run_selftest_case(
        "strike-without-adr",
        lambda root: replace_once(
            root / "docs/site-requirements.md",
            "- [ ] Foundation clause remains visible.",
            "- ~~[ ] Foundation clause remains visible.~~",
        ),
        1,
        "rule 1 (added ADR): FAIL",
    )
    run_selftest_case(
        "strike-wrong-section-name",
        lambda root: (
            replace_once(
                root / "docs/site-requirements.md",
                "- [ ] Foundation clause remains visible.",
                "- ~~[ ] Foundation clause remains visible.~~",
            ),
            add_adr(root, "0002", "wrong-name", "OQ-5"),
        ),
        1,
        "rule 3 (named sections): FAIL",
    )
    # This is the real :591 shape: one table cell is struck while the rest of
    # the row is rewritten.  Whole-line equality would reject this case.
    run_selftest_case(
        "fragment-rewrite-real-591",
        lambda root: (
            replace_once(
                root / "docs/site-requirements.md",
                "| Norse | Huginn | Existing name | Rejected |",
                "| Norse | ~~Huginn~~ | Existing name | **Unusable** — rejected |",
            ),
            add_adr(root, "0002", "fragment-rewrite", "OQ-5"),
        ),
        0,
        "decision: legitimate amendment; all three rules pass",
    )
    # This is the real :526 shape: the nearest section is ## 13, with no
    # clause identifier for rule 3 to extract.
    run_selftest_case(
        "unidentified-heading-real-526",
        lambda root: (
            replace_once(
                root / "docs/site-requirements.md",
                "- [ ] Criterion under a heading without a clause id.",
                "- ~~[ ] Criterion under a heading without a clause id.~~",
            ),
            add_adr(root, "0002", "wrong-heading", "P0"),
        ),
        1,
        "rule 3 (named sections): FAIL",
    )
    run_selftest_case(
        "injectivity-one-span-for-ten-lines",
        lambda root: (
            replace_once(
                root / "docs/site-requirements.md",
                "- [ ] Bulk line 01 remains visible.\n"
                "- [ ] Bulk line 02 remains visible.\n"
                "- [ ] Bulk line 03 remains visible.\n"
                "- [ ] Bulk line 04 remains visible.\n"
                "- [ ] Bulk line 05 remains visible.\n"
                "- [ ] Bulk line 06 remains visible.\n"
                "- [ ] Bulk line 07 remains visible.\n"
                "- [ ] Bulk line 08 remains visible.\n"
                "- [ ] Bulk line 09 remains visible.\n"
                "- [ ] Bulk line 10 remains visible.\n",
                "- ~~the~~\n",
            ),
            add_adr(root, "0002", "bulk-loss", "P0"),
        ),
        1,
        "rule 2 (struck traces): FAIL",
    )
    run_selftest_case(
        "injectivity-per-hunk",
        lambda root: (
            replace_once(
                root / "docs/site-requirements.md",
                "- [ ] Foundation clause remains visible.",
                "- ~~Foundation clause remains visible.~~ ~~P0~~ ~~trace~~ ~~one~~ ~~two~~ ~~three~~ ~~four~~ ~~five~~ ~~six~~ ~~seven~~",
            ),
            replace_once(
                root / "docs/site-requirements.md",
                "- [ ] Silent line 01 remains visible.\n"
                "- [ ] Silent line 02 remains visible.\n"
                "- [ ] Silent line 03 remains visible.\n"
                "- [ ] Silent line 04 remains visible.\n"
                "- [ ] Silent line 05 remains visible.\n"
                "- [ ] Silent line 06 remains visible.\n"
                "- [ ] Silent line 07 remains visible.\n"
                "- [ ] Silent line 08 remains visible.\n"
                "- [ ] Silent line 09 remains visible.\n"
                "- [ ] Silent line 10 remains visible.\n",
                "",
            ),
            add_adr(root, "0002", "two-hunk-loss", "P0"),
        ),
        1,
        "hunk 2: FAIL",
    )
    run_selftest_case(
        "whole-line-strike",
        lambda root: (
            replace_once(
                root / "docs/site-requirements.md",
                "- [ ] Foundation clause remains visible.",
                "- ~~[ ] Foundation clause remains visible.~~",
            ),
            add_adr(root, "0002", "whole-line", "P0"),
        ),
        0,
        "decision: legitimate amendment; all three rules pass",
    )
    run_selftest_case(
        "no-diff",
        lambda root: None,
        0,
        "decision: no contract diff; accepted",
    )
    print("amendcheck selftest: 9 cases pass")
    return 0


def main(argv: list[str]) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("base", nargs="?")
    parser.add_argument("head", nargs="?")
    parser.add_argument("--selftest", action="store_true")
    args = parser.parse_args(argv)
    if args.selftest:
        return selftest()
    if args.base is None:
        parser.error("BASE revision is required unless --selftest is used")
    if args.head is None:
        args.head = "HEAD"
    try:
        repo = repository_root()
        passed, output = assess(repo, args.base, args.head)
    except AmendcheckError as error:
        print(f"amendcheck: ERROR — {error}", file=sys.stderr)
        return 1
    print(output)
    return 0 if passed else 1


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
