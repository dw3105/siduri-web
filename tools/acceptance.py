#!/usr/bin/env python3
"""Check that acceptance fragments cover the current site contract.

The contract is the source of truth for the criterion set.  Acceptance
fragments may say that a criterion is not met, but they must identify every
criterion exactly once and include enough evidence for the verdict to be
audited.

    python3 tools/acceptance.py
    python3 tools/acceptance.py --operator
    python3 tools/acceptance.py --selftest
    python3 tools/acceptance.py --operator-selftest

The checker uses only the standard library.  The selftest passes copies of the
contract and synthetic fragments to a subprocess, so it never plants a bad
fixture in the repository's real ``docs/`` or ``acceptance/`` directories.
"""
from __future__ import annotations

import argparse
from collections import Counter
from dataclasses import dataclass
import datetime as dt
import pathlib
import re
import subprocess
import sys
import tempfile
import unicodedata


CONTRACT_HEADING = "## 13 · Acceptance criteria"
CHECKBOX = re.compile(r"^- \[ \](.*)$")
CRITERION_HEADING = re.compile(
    r"^## Criterion\s+(\d+)\s*(?:—|–|-)\s*(.*?)\s*$", re.MULTILINE
)
FIELD = re.compile(r"^(verdict|ran|at|saw|notes):[ \t]*(.*?)\s*$", re.MULTILINE)
TIMESTAMP = re.compile(r"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$")
VERDICTS = ("held", "deferred", "open", "failed")
COUNT_LABELS = (*VERDICTS, "unrecognised")
OPERATOR_REPORTS = (
    pathlib.Path("acceptance/OPERATOR.md"),
    pathlib.Path("acceptance/W2/OPERATOR-2.md"),
)
OPERATOR_PATH = re.compile(r"(?<![\w])`?/(?:[\w.-]+/)*`?")
OPERATOR_VERDICT = re.compile(
    r"^\s*(?:\*\*)?(held|deferred|open|failed)(?:\*\*)?(?=$|[^\w])",
    re.IGNORECASE,
)
OPERATOR_CROSS_REFERENCE = re.compile(
    r"^see\s+`?([^`#\s]+)`?\s*#\s*(\d+)\s*$", re.IGNORECASE
)
COMMAND_WORDS = {
    "awk",
    "bc",
    "bash",
    "cat",
    "chmod",
    "command",
    "cp",
    "curl",
    "diff",
    "echo",
    "find",
    "for",
    "git",
    "go",
    "grep",
    "head",
    "make",
    "mv",
    "node",
    "npm",
    "npx",
    "python",
    "python3",
    "printf",
    "rg",
    "rm",
    "sed",
    "sh",
    "sort",
    "tail",
    "test",
    "true",
    "uniq",
    "wc",
}
NUMBER_WORDS = {
    "zero": 0,
    "one": 1,
    "two": 2,
    "three": 3,
    "four": 4,
    "five": 5,
    "six": 6,
    "seven": 7,
    "eight": 8,
    "nine": 9,
    "ten": 10,
    "eleven": 11,
    "twelve": 12,
    "thirteen": 13,
    "fourteen": 14,
    "fifteen": 15,
    "sixteen": 16,
    "seventeen": 17,
    "eighteen": 18,
    "nineteen": 19,
    "twenty": 20,
    "thirty": 30,
    "forty": 40,
    "fifty": 50,
    "sixty": 60,
    "seventy": 70,
    "eighty": 80,
    "ninety": 90,
}


class AcceptanceInputError(ValueError):
    """The contract or an acceptance fragment cannot be read as a report."""


@dataclass(frozen=True)
class Criterion:
    number: int
    text: str


@dataclass
class Row:
    number: int
    heading_text: str
    fields: dict[str, str]
    body: str
    path: pathlib.Path


@dataclass
class Report:
    criteria: list[Criterion]
    rows: list[Row]
    errors: list[str]

    @property
    def counts(self) -> Counter[str]:
        counts: Counter[str] = Counter()
        for row in self.rows:
            verdict = row.fields.get("verdict", "").strip()
            counts[verdict if verdict in VERDICTS else "unrecognised"] += 1
        return counts

    @property
    def accepted(self) -> bool:
        if self.errors or not self.rows:
            return False
        return all(row.fields.get("verdict", "").strip() in {"held", "deferred"} for row in self.rows)


@dataclass
class OperatorRow:
    path: pathlib.Path
    line: int
    section: str
    number: int
    values: dict[str, str]
    raw_verdict: str
    verdict_head: str | None
    non_verdict: bool


@dataclass
class OperatorBlock:
    path: pathlib.Path
    line: int
    heading: str
    references: list[tuple[str, int]]
    fields: dict[str, str]
    verdict_head: str | None


@dataclass
class SweepRow:
    path: pathlib.Path
    line: int
    page: str
    result: str


@dataclass
class OperatorReport:
    path: pathlib.Path
    rows: list[OperatorRow]
    blocks: list[OperatorBlock]
    sweep_rows: list[SweepRow]
    errors: list[str]


def normalise(text: str) -> str:
    """Compare quoted Markdown without treating wrapping whitespace as drift."""
    return re.sub(r"\s+", " ", unicodedata.normalize("NFC", text).strip())


def read_text(path: pathlib.Path) -> str:
    try:
        return path.read_text(encoding="utf-8")
    except (OSError, UnicodeError) as error:
        raise AcceptanceInputError(f"cannot read {path}: {error}") from error


def derive_criteria(contract: pathlib.Path) -> list[Criterion]:
    """Derive the numbered criterion list from section 13, not line numbers."""
    lines = read_text(contract).splitlines()
    try:
        start = next(index for index, line in enumerate(lines) if line.rstrip() == CONTRACT_HEADING)
    except StopIteration as error:
        raise AcceptanceInputError(f"contract heading not found: {CONTRACT_HEADING}") from error

    criteria: list[Criterion] = []
    for line in lines[start + 1 :]:
        if line.startswith("## "):
            break
        match = CHECKBOX.match(line)
        if match:
            text = match.group(1).strip()
            if not text:
                raise AcceptanceInputError("acceptance criterion has no text")
            criteria.append(Criterion(len(criteria) + 1, text))

    if not criteria:
        raise AcceptanceInputError(f"no unchecked criteria found below {CONTRACT_HEADING}")

    seen: set[str] = set()
    for criterion in criteria:
        key = normalise(criterion.text)
        if key in seen:
            raise AcceptanceInputError(f"contract repeats criterion {criterion.number}: {criterion.text}")
        seen.add(key)
    return criteria


def fragment_rows(fragments_dir: pathlib.Path) -> tuple[list[Row], list[str]]:
    """Read every direct ``acceptance/*.md`` fragment in stable order."""
    if not fragments_dir.is_dir():
        return [], [f"acceptance fragments directory does not exist: {fragments_dir}"]

    # Operator reports have their own grammar and are checked by --operator.
    # Keeping them out of this path preserves the criterion report's stable
    # 18-row contract and prevents prose tables from becoming criteria.
    paths = sorted(
        path for path in fragments_dir.glob("*.md") if path.name != "OPERATOR.md"
    )
    if not paths:
        return [], [f"no acceptance fragments found in {fragments_dir}"]

    rows: list[Row] = []
    errors: list[str] = []
    for path in paths:
        text = read_text(path)
        matches = list(CRITERION_HEADING.finditer(text))
        for line in text.splitlines():
            if line.startswith("## Criterion ") and not CRITERION_HEADING.match(line):
                errors.append(f"{path}: malformed criterion heading: {line}")
        if not matches:
            errors.append(f"{path}: no criterion rows found")
            continue

        for index, match in enumerate(matches):
            end = matches[index + 1].start() if index + 1 < len(matches) else len(text)
            body = text[match.end() : end]
            fields: dict[str, str] = {}
            for field_match in FIELD.finditer(body):
                name, value = field_match.groups()
                if name in fields:
                    errors.append(f"{path}: criterion {match.group(1)} has multiple {name} fields")
                fields[name] = value.strip()
            rows.append(
                Row(
                    number=int(match.group(1)),
                    heading_text=match.group(2).strip(),
                    fields=fields,
                    body=body,
                    path=path,
                )
            )
    return rows, errors


def named_owner(body: str) -> bool:
    """Return whether a deferred row names an owner in prose or an owner field."""
    if re.search(r"\bno\s+(?:named\s+)?owner\b", body, re.IGNORECASE):
        return False

    patterns = (
        r"\bowner\s*[:=]\s*([A-Za-z0-9][A-Za-z0-9_./-]*(?:[ \t]+[A-Za-z0-9][A-Za-z0-9_./-]*)*)",
        r"\bowner[ \t]+(?:is[ \t]+)?([A-Za-z0-9][A-Za-z0-9_./-]*(?:[ \t]+[A-Za-z0-9][A-Za-z0-9_./-]*)*)",
        r"\bowned[ \t]+by[ \t]+([A-Za-z0-9][A-Za-z0-9_./-]*(?:[ \t]+[A-Za-z0-9][A-Za-z0-9_./-]*)*)",
    )
    unnamed = {"none", "tbd", "unknown", "unassigned", "n/a", "na", "-"}
    for pattern in patterns:
        match = re.search(pattern, body, re.IGNORECASE)
        if match and match.group(1).strip().rstrip(".,;:").lower() not in unnamed:
            return True
    return False


def valid_timestamp(value: str) -> bool:
    if not TIMESTAMP.fullmatch(value):
        return False
    try:
        dt.datetime.strptime(value, "%Y-%m-%dT%H:%M:%SZ")
    except ValueError:
        return False
    return True


def operator_verdict(value: str) -> str | None:
    """Return the four-word verdict head, leaving its qualifier intact."""
    match = OPERATOR_VERDICT.match(value)
    return match.group(1).lower() if match else None


def operator_non_verdict(value: str) -> bool:
    """Whether a table cell deliberately declines to judge a site row."""
    return value.strip().lower() == "n/a" or bool(OPERATOR_CROSS_REFERENCE.fullmatch(value.strip()))


def operator_table_cells(line: str) -> list[str]:
    value = line.strip()
    if value.startswith("|"):
        value = value[1:]
    if value.endswith("|"):
        value = value[:-1]
    return [cell.strip() for cell in value.split("|")]


def operator_separator(line: str) -> bool:
    cells = operator_table_cells(line)
    return bool(cells) and all(re.fullmatch(r":?-{3,}:?", cell) for cell in cells)


def operator_section_name(title: str | None) -> str:
    if title is None:
        return "(preamble)"
    if title.strip().lower() == "rows":
        return "Rows"
    if title.lstrip("`").startswith("/"):
        return re.split(r"\s+—\s+", title.strip(), maxsplit=1)[0].strip("`")
    return title.strip()


def operator_h2_title(line: str) -> str | None:
    match = re.fullmatch(r"##\s+(.+?)\s*", line)
    return match.group(1).strip() if match else None


def operator_number_references(text: str) -> list[int]:
    references: list[int] = []
    for match in re.finditer(r"(?<!\d)(\d+)(?:\s*([–—-])\s*(\d+))?", text):
        start = int(match.group(1))
        end = int(match.group(3) or match.group(1))
        if end < start:
            continue
        references.extend(range(start, end + 1))
    return references


def operator_heading_references(
    heading: str, row_sections: set[str]
) -> tuple[list[tuple[str, int]], list[str]]:
    """Read `(section, number)` references from a probe heading.

    Page paths make the W1 headings explicit.  W2's `Rows ...` headings use
    the sole `Rows` table as their implicit section.  A heading with no row
    phrase is intentionally a valid zero-row absence block.
    """
    references: list[tuple[str, int]] = []
    errors: list[str] = []
    path_matches = list(OPERATOR_PATH.finditer(heading))
    if path_matches:
        for index, path_match in enumerate(path_matches):
            end = path_matches[index + 1].start() if index + 1 < len(path_matches) else len(heading)
            segment = heading[path_match.end() : end]
            row_match = re.search(r"\brows?\b(.*?)(?=\s+—|$)", segment, re.IGNORECASE)
            if not row_match:
                continue
            for number in operator_number_references(row_match.group(1)):
                references.append((path_match.group(0).strip("`"), number))
        return references, errors

    if not re.search(r"\brows?\b", heading, re.IGNORECASE):
        return references, errors
    if "Rows" in row_sections:
        section = "Rows"
    elif len(row_sections) == 1:
        section = next(iter(row_sections))
    else:
        errors.append(f"{heading!r} has row references but no unambiguous table section")
        return references, errors
    row_match = re.search(r"\brows?\b(.*?)(?=\s+—|$)", heading, re.IGNORECASE)
    if row_match:
        references.extend((section, number) for number in operator_number_references(row_match.group(1)))
    return references, errors


def operator_command(value: str) -> bool:
    """Recognise a shell command in a possibly comma-separated probe."""
    tokens = re.findall(r"(?<![A-Za-z0-9_.-])([A-Za-z][A-Za-z0-9_-]*)(?=\s|$)", value)
    return any(token.lower() in COMMAND_WORDS for token in tokens)


def operator_count_value(token: str) -> int | None:
    token = token.strip().strip("*`.,").lower()
    if token.isdigit():
        return int(token)
    pieces = token.replace("-", " ").split()
    if not pieces or any(piece not in NUMBER_WORDS for piece in pieces):
        return None
    return sum(NUMBER_WORDS[piece] for piece in pieces)


COUNT_TOKEN = r"\d+|[A-Za-z]+(?:-[A-Za-z]+)*"


def close_count(
    summary: str, pattern: str, label: str, errors: list[str]
) -> int | None:
    match = re.search(pattern, summary, re.IGNORECASE)
    if not match:
        errors.append(f"close count {label!r} is unparseable")
        return None
    value = operator_count_value(match.group("count"))
    if value is None:
        errors.append(f"close count {label!r} value {match.group('count')!r} is unparseable")
    return value


def operator_sweep_annotation_count(result: str, label: str, errors: list[str]) -> int:
    lowered = result.lower()
    if "clean" in lowered or "unlinked" in lowered or "no annotation" in lowered:
        return 0
    if re.search(r"\bok\b", result, re.IGNORECASE):
        return 0
    if not result.strip():
        errors.append(f"{label} sweep result is unparseable")
        return 0
    # The sweep has one prose result per page.  Only the tools row records two
    # separate annotations in that result; its `tool card; filter nav` split
    # is deliberate.  Other semicolons are punctuation inside one finding.
    return 2 if "tool card" in lowered and "filter nav" in lowered else 1


def operator_close_errors(
    path: pathlib.Path,
    text: str,
    rows: list[OperatorRow],
    sweep_rows: list[SweepRow],
) -> list[str]:
    errors: list[str] = []
    close_match = re.search(r"^##\s+Close\s*$", text, re.MULTILINE)
    if not close_match:
        return [f"{path}: missing ## Close section"]
    body = text[close_match.end() :]
    summary = body.split("\n\n", 1)[0]
    # The old hand-count is deliberately retained as historical evidence in
    # OPERATOR.md.  Only the current close line is arithmetic input.
    summary = summary.split("The close line read", 1)[0]

    verdict_rows = [row for row in rows if row.verdict_head is not None]
    if re.search(r"pages walked", summary, re.IGNORECASE):
        match = re.search(
            rf"(?P<count>{COUNT_TOKEN})\s+of\s+(?P<total>{COUNT_TOKEN})\s+pages\s+walked",
            summary,
            re.IGNORECASE,
        )
        if not match:
            errors.append(f"{path}: close count 'pages walked' is unparseable")
        else:
            walked = operator_count_value(match.group("count"))
            total = operator_count_value(match.group("total"))
            page_sections = {
                row.section
                for row in verdict_rows
                if row.section not in {"Rows", "(preamble)"}
            }
            expected = len(sweep_rows) + len(page_sections)
            if walked is None or total is None:
                errors.append(f"{path}: close count 'pages walked' value is unparseable")
            elif walked != expected or total != expected:
                errors.append(
                    f"{path}: close count 'pages walked' is {walked} of {total}, expected {expected} of {expected}"
                )

        annotation_count = close_count(
            summary,
            rf"(?P<count>{COUNT_TOKEN})\s+annotations\b",
            "annotations",
            errors,
        )
        # A duplicate `see ...` entry is not a verdict row for coverage, but
        # it remains an operator annotation row and therefore belongs in the
        # close-out's annotation count.
        expected_annotations = len(rows)
        for sweep_row in sweep_rows:
            expected_annotations += operator_sweep_annotation_count(
                sweep_row.result,
                f"{path}:{sweep_row.line}",
                errors,
            )
        if annotation_count is not None and annotation_count != expected_annotations:
            errors.append(
                f"{path}: close count 'annotations' is {annotation_count}, expected {expected_annotations}"
            )

        clean_count = close_count(
            summary,
            rf"(?P<count>{COUNT_TOKEN})\s+clean\b",
            "clean",
            errors,
        )
        expected_clean = sum(1 for row in sweep_rows if "clean" in row.result.lower())
        if clean_count is not None and clean_count != expected_clean:
            errors.append(f"{path}: close count 'clean' is {clean_count}, expected {expected_clean}")

        unlinked_count = close_count(
            summary,
            rf"(?P<count>{COUNT_TOKEN})\s+unlinked\b",
            "unlinked",
            errors,
        )
        expected_unlinked = sum(1 for row in sweep_rows if "unlinked" in row.result.lower())
        if unlinked_count is not None and unlinked_count != expected_unlinked:
            errors.append(
                f"{path}: close count 'unlinked' is {unlinked_count}, expected {expected_unlinked}"
            )
        return errors

    if re.search(r"templates walked", summary, re.IGNORECASE):
        pair = re.search(
            rf"(?P<count>{COUNT_TOKEN})\s+of\s+(?P<total>{COUNT_TOKEN})\s+templates\s+walked",
            summary,
            re.IGNORECASE,
        )
        pages = {
            row.values.get("page", "").strip()
            for row in verdict_rows
            if row.values.get("page", "").strip()
        }
        clean_rows = [
            row for row in verdict_rows if "clean" in row.values.get("note", "").lower()
        ]
        if not pair:
            errors.append(f"{path}: close count 'templates walked' is unparseable")
        else:
            walked = operator_count_value(pair.group("count"))
            total = operator_count_value(pair.group("total"))
            expected_walked = len(pages)
            expected_total = len(verdict_rows) - len(clean_rows)
            if walked is None or total is None:
                errors.append(f"{path}: close count 'templates walked' value is unparseable")
            elif walked != expected_walked or total != expected_total:
                errors.append(
                    f"{path}: close count 'templates walked' is {walked} of {total}, expected {expected_walked} of {expected_total}"
                )

        row_count = close_count(
            summary,
            rf"(?P<count>{COUNT_TOKEN})\s+rows\b",
            "rows",
            errors,
        )
        if row_count is not None and row_count != len(rows):
            errors.append(f"{path}: close count 'rows' is {row_count}, expected {len(rows)}")

        no_annotation = re.search(
            rf"(?P<count>{COUNT_TOKEN})\s+of\s+the\s+(?P<total>{COUNT_TOKEN})\s+pages\s+drew\s+no\s+annotation",
            summary,
            re.IGNORECASE,
        )
        if not no_annotation:
            errors.append(f"{path}: close count 'pages with no annotation' is unparseable")
        else:
            clean_count = operator_count_value(no_annotation.group("count"))
            page_count = operator_count_value(no_annotation.group("total"))
            if clean_count is None or page_count is None:
                errors.append(f"{path}: close count 'pages with no annotation' value is unparseable")
            elif clean_count != len(clean_rows) or page_count != len(pages):
                errors.append(
                    f"{path}: close count 'pages with no annotation' is {clean_count} of {page_count}, expected {len(clean_rows)} of {len(pages)}"
                )
        return errors

    return [f"{path}: ## Close has no parseable walked-count form"]


def operator_report(path: pathlib.Path) -> OperatorReport:
    try:
        text = read_text(path)
    except AcceptanceInputError as error:
        return OperatorReport(path, [], [], [], [str(error)])

    lines = text.splitlines()
    rows: list[OperatorRow] = []
    sweep_rows: list[SweepRow] = []
    errors: list[str] = []
    current_title: str | None = None
    index = 0
    while index < len(lines):
        title = operator_h2_title(lines[index])
        if title is not None:
            current_title = title
            index += 1
            continue
        if (
            index + 1 < len(lines)
            and lines[index].lstrip().startswith("|")
            and operator_separator(lines[index + 1])
        ):
            headers = operator_table_cells(lines[index])
            lowered = [header.strip().lower() for header in headers]
            verdict_index = next((i for i, header in enumerate(lowered) if header == "verdict"), None)
            is_verdict_table = bool(headers) and lowered[0] == "#" and verdict_index is not None
            result_index = next((i for i, header in enumerate(lowered) if header == "result"), None)
            is_sweep_table = (
                bool(headers)
                and lowered[0] == "page"
                and result_index is not None
                and current_title is not None
                and current_title.lower().startswith("pages")
            )
            row_index = index + 2
            while row_index < len(lines) and lines[row_index].lstrip().startswith("|"):
                values = operator_table_cells(lines[row_index])
                if is_verdict_table:
                    if len(values) <= verdict_index:
                        errors.append(f"{path}:{row_index + 1}: operator table row has no verdict cell")
                    elif values and re.fullmatch(r"\d+", values[0]):
                        columns = {
                            header: values[column]
                            for column, header in enumerate(lowered)
                            if column < len(values) and header
                        }
                        raw_verdict = values[verdict_index]
                        verdict_head = operator_verdict(raw_verdict)
                        non_verdict = operator_non_verdict(raw_verdict)
                        if verdict_head is None and not non_verdict:
                            errors.append(
                                f"{path}:{row_index + 1}: section {operator_section_name(current_title)!r} row {values[0]} has invalid verdict {raw_verdict!r}"
                            )
                        rows.append(
                            OperatorRow(
                                path=path,
                                line=row_index + 1,
                                section=operator_section_name(current_title),
                                number=int(values[0]),
                                values=columns,
                                raw_verdict=raw_verdict,
                                verdict_head=verdict_head,
                                non_verdict=non_verdict,
                            )
                        )
                elif is_sweep_table and values and len(values) > result_index:
                    sweep_rows.append(
                        SweepRow(path, row_index + 1, values[0].strip(), values[result_index].strip())
                    )
                row_index += 1
            index = row_index
            continue
        index += 1

    row_sections = {row.section for row in rows}
    blocks: list[OperatorBlock] = []
    probe_start: int | None = None
    for line_index, line in enumerate(lines):
        title = operator_h2_title(line)
        if title is not None and title.lower() == "probes":
            probe_start = line_index + 1
            break
    if probe_start is None:
        errors.append(f"{path}: missing ## Probes section")
    else:
        probe_end = len(lines)
        for line_index in range(probe_start, len(lines)):
            if operator_h2_title(lines[line_index]) is not None:
                probe_end = line_index
                break
        block_starts = [
            line_index
            for line_index in range(probe_start, probe_end)
            if lines[line_index].startswith("### ")
        ]
        if not block_starts:
            errors.append(f"{path}: ## Probes has no blocks")
        for block_index, start in enumerate(block_starts):
            end = block_starts[block_index + 1] if block_index + 1 < len(block_starts) else probe_end
            heading = lines[start][4:].strip()
            body = "\n".join(lines[start + 1 : end])
            fields: dict[str, str] = {}
            for field_match in FIELD.finditer(body):
                name, value = field_match.groups()
                if name in fields:
                    errors.append(f"{path}:{start + 1}: probe {heading!r} has multiple {name} fields")
                fields[name] = value.strip()
            references, reference_errors = operator_heading_references(heading, row_sections)
            errors.extend(f"{path}:{start + 1}: {error}" for error in reference_errors)
            blocks.append(
                OperatorBlock(
                    path=path,
                    line=start + 1,
                    heading=heading,
                    references=references,
                    fields=fields,
                    verdict_head=operator_verdict(fields.get("verdict", "")),
                )
            )

    by_address: dict[tuple[str, int], OperatorRow] = {}
    for row in rows:
        address = (row.section, row.number)
        if address in by_address:
            errors.append(
                f"{path}:{row.line}: duplicate table row address ({row.section!r}, {row.number})"
            )
        else:
            by_address[address] = row

        if not row.raw_verdict.strip():
            errors.append(f"{path}:{row.line}: section {row.section!r} row {row.number} is missing a verdict")

    for row in rows:
        if row.non_verdict and row.raw_verdict.strip().lower().startswith("see"):
            cross_reference = OPERATOR_CROSS_REFERENCE.fullmatch(row.raw_verdict.strip())
            if cross_reference:
                target = (cross_reference.group(1), int(cross_reference.group(2)))
                target_row = by_address.get(target)
                if target_row is None or target_row.verdict_head is None:
                    errors.append(
                        f"{path}:{row.line}: row ({row.section!r}, {row.number}) cross-reference names no verdict row: {target}"
                    )
        elif row.non_verdict is False and row.verdict_head is None:
            # The detailed grammar error was emitted while reading the table.
            pass

    coverage: dict[tuple[str, int], list[str]] = {
        (row.section, row.number): [] for row in rows if row.verdict_head is not None
    }
    for block in blocks:
        label = f"{path}:{block.line} probe {block.heading!r}"
        fields = block.fields
        if block.verdict_head is None:
            errors.append(f"{label} has invalid or missing verdict")
        for required in ("ran", "at", "saw"):
            value = fields.get(required, "").strip()
            if not value:
                errors.append(f"{label} is missing {required}")
            elif required == "at" and not valid_timestamp(value):
                errors.append(f"{label} at value {value!r} is not a timestamp")
            elif required == "ran" and not operator_command(value):
                errors.append(f"{label} ran value {value!r} is not a command")

        seen_in_block: set[tuple[str, int]] = set()
        for address in block.references:
            if address in seen_in_block:
                errors.append(f"{label} names row {address} more than once")
            seen_in_block.add(address)
            row = by_address.get(address)
            if row is None:
                errors.append(f"{label} names no table row at address {address}")
                continue
            if row.verdict_head is not None:
                coverage[address].append(label)
                if block.verdict_head is not None and block.verdict_head != row.verdict_head:
                    errors.append(
                        f"{label} verdict {block.verdict_head!r} disagrees with row {address} verdict {row.verdict_head!r}"
                    )
            elif row.non_verdict and row.raw_verdict.strip().lower().startswith("see"):
                cross_reference = OPERATOR_CROSS_REFERENCE.fullmatch(row.raw_verdict.strip())
                if cross_reference and block.verdict_head is not None:
                    target = (cross_reference.group(1), int(cross_reference.group(2)))
                    target_row = by_address.get(target)
                    if target_row and target_row.verdict_head and block.verdict_head != target_row.verdict_head:
                        errors.append(
                            f"{label} verdict {block.verdict_head!r} disagrees with cross-referenced row {target} verdict {target_row.verdict_head!r}"
                        )

    for address, covered_by in coverage.items():
        if not covered_by:
            errors.append(f"{path}: verdict row {address} is uncovered")
        elif len(covered_by) > 1:
            errors.append(f"{path}: verdict row {address} is covered by multiple probe blocks")

    errors.extend(operator_close_errors(path, text, rows, sweep_rows))
    return OperatorReport(path, rows, blocks, sweep_rows, errors)


def operator_check(paths: list[pathlib.Path]) -> int:
    reports = [operator_report(path) for path in paths]
    errors = [error for report in reports for error in report.errors]
    print("accepted" if not errors else "not accepted")
    print(f"operator reports: {len(reports)}")
    print(f"operator verdict rows: {sum(sum(row.verdict_head is not None for row in report.rows) for report in reports)}")
    if errors:
        print("operator reports: refused:", file=sys.stderr)
        for error in errors:
            print(f"  {error}", file=sys.stderr)
        return 1
    return 0


def inspect(contract: pathlib.Path, fragments_dir: pathlib.Path) -> Report:
    criteria = derive_criteria(contract)
    rows, errors = fragment_rows(fragments_dir)
    by_number = {criterion.number: criterion for criterion in criteria}
    coverage: dict[int, list[pathlib.Path]] = {criterion.number: [] for criterion in criteria}

    for row in rows:
        criterion = by_number.get(row.number)
        label = f"criterion {row.number} in {row.path}"
        if criterion is None:
            errors.append(f"{label} is not in contract")
        else:
            coverage[row.number].append(row.path)
            if normalise(row.heading_text) != normalise(criterion.text):
                errors.append(
                    f"{label} heading does not match contract: {row.heading_text!r}"
                )

        verdict = row.fields.get("verdict", "").strip()
        if verdict not in VERDICTS:
            if verdict:
                errors.append(f"{label} has unrecognised verdict {verdict!r}")
            else:
                errors.append(f"{label} row missing verdict")

        for required in ("ran", "at", "saw"):
            value = row.fields.get(required, "").strip()
            if not value:
                errors.append(f"{label} row missing {required}")
            elif required == "at" and not valid_timestamp(value):
                errors.append(f"{label} at value {value!r} is not a timestamp")

        if verdict == "deferred" and not named_owner(row.body):
            errors.append(f"{label} deferred row has no named owner")

    for criterion in criteria:
        paths = coverage[criterion.number]
        if not paths:
            errors.append(
                f"criterion {criterion.number} is not covered by any fragment: {criterion.text}"
            )
        elif len(paths) > 1:
            listed = ", ".join(str(path) for path in paths)
            errors.append(f"criterion {criterion.number} is covered by multiple fragments: {listed}")

    return Report(criteria=criteria, rows=rows, errors=errors)


def print_report(report: Report) -> None:
    print("accepted" if report.accepted else "not accepted")
    print(f"criteria: {len(report.criteria)}")
    print(f"rows: {len(report.rows)}")
    print("counts: " + ", ".join(f"{label}={report.counts.get(label, 0)}" for label in COUNT_LABELS))
    if report.errors:
        print("acceptance: refused:", file=sys.stderr)
        for error in report.errors:
            print(f"  {error}", file=sys.stderr)


def check(contract: pathlib.Path, fragments_dir: pathlib.Path) -> int:
    try:
        report = inspect(contract, fragments_dir)
    except AcceptanceInputError as error:
        print("not accepted")
        print("acceptance: refused:", file=sys.stderr)
        print(f"  {error}", file=sys.stderr)
        return 1
    print_report(report)
    if report.errors:
        return 1  # the report is malformed
    return 0 if report.accepted else 2  # valid, and not accepted


def criterion_block(text: str, number: int) -> str:
    matches = list(CRITERION_HEADING.finditer(text))
    for index, match in enumerate(matches):
        if int(match.group(1)) != number:
            continue
        end = matches[index + 1].start() if index + 1 < len(matches) else len(text)
        return text[match.start() : end]
    raise AssertionError(f"synthetic report has no criterion {number}")


def without_criterion_block(text: str, number: int) -> str:
    matches = list(CRITERION_HEADING.finditer(text))
    for index, match in enumerate(matches):
        if int(match.group(1)) != number:
            continue
        end = matches[index + 1].start() if index + 1 < len(matches) else len(text)
        return text[: match.start()] + text[end:]
    raise AssertionError(f"synthetic report has no criterion {number}")


def change_field(text: str, number: int, field: str, replacement: str | None) -> str:
    matches = list(CRITERION_HEADING.finditer(text))
    for index, match in enumerate(matches):
        if int(match.group(1)) != number:
            continue
        end = matches[index + 1].start() if index + 1 < len(matches) else len(text)
        block = text[match.start() : end]
        pattern = re.compile(rf"^{re.escape(field)}:.*$", re.MULTILINE)
        if replacement is None:
            changed, count = pattern.subn("", block, count=1)
        else:
            changed, count = pattern.subn(f"{field}: {replacement}", block, count=1)
        if count != 1:
            raise AssertionError(f"synthetic report criterion {number} has no {field} field")
        return text[: match.start()] + changed + text[end:]
    raise AssertionError(f"synthetic report has no criterion {number}")


def synthetic_report(criteria: list[Criterion]) -> str:
    sections = ["# Synthetic acceptance report", "", "wave: W1", ""]
    for criterion in criteria:
        sections.extend(
            [
                f"## Criterion {criterion.number} — {criterion.text}",
                "",
                "verdict: held",
                "ran: selftest",
                "at: 2026-08-27T00:00:00Z",
                "saw: exit 0",
                "red proof: synthetic bad input was refused",
                "notes: synthetic complete report",
                "",
            ]
        )
    return "\n".join(sections)


def make_case(
    root: pathlib.Path,
    label: str,
    contract_text: str,
    fragments: dict[str, str],
) -> tuple[pathlib.Path, pathlib.Path]:
    case = root / label
    contract = case / "docs" / "site-requirements.md"
    fragments_dir = case / "acceptance"
    contract.parent.mkdir(parents=True)
    fragments_dir.mkdir()
    contract.write_text(contract_text, encoding="utf-8")
    for name, text in fragments.items():
        (fragments_dir / name).write_text(text, encoding="utf-8")
    return contract, fragments_dir


def run_subprocess(contract: pathlib.Path, fragments_dir: pathlib.Path) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        [
            sys.executable,
            str(pathlib.Path(__file__).resolve()),
            "--contract",
            str(contract),
            "--fragments-dir",
            str(fragments_dir),
        ],
        capture_output=True,
        text=True,
        check=False,
    )



def _section_lines(text: str, heading: str) -> list[str]:
    """Lines under `heading`, up to the next `## ` heading. Selftest cross-check."""
    out, inside = [], False
    for line in text.split("\n"):
        if line.startswith(heading):
            inside = True
            continue
        if inside and line.startswith("## "):
            break
        if inside:
            out.append(line)
    return out


def selftest() -> int:
    real_contract = pathlib.Path(__file__).resolve().parents[1] / "docs" / "site-requirements.md"
    real_contract_text = read_text(real_contract)
    criteria = derive_criteria(real_contract)
    # Derive the count a second way, by a different method, and compare.
    # Pinning a literal here would reintroduce one layer up exactly the coupling
    # the heading anchor removed: the guard would go red the day the operator
    # legitimately amends the contract.
    naive = sum(
        1
        for line in _section_lines(real_contract_text, CONTRACT_HEADING)
        if line.startswith("- [ ]")
    )
    assert len(criteria) > 0, "the contract anchor matched nothing"
    assert len(criteria) == naive, (
        f"parser derived {len(criteria)} criteria, naive count of the section saw {naive}"
    )
    baseline = synthetic_report(criteria)

    with tempfile.TemporaryDirectory(prefix="siduri-acceptance-") as raw:
        root = pathlib.Path(raw)
        cases: list[tuple[str, str, dict[str, str], str, str]] = []

        cases.append(
            (
                "uncovered-criterion-1",
                real_contract_text,
                {"complete.md": without_criterion_block(baseline, 1)},
                "criterion 1 is not covered by any fragment",
                "removed criterion 1 row from the copied fragments",
            )
        )
        cases.append(
            (
                "duplicate-criterion-1",
                real_contract_text,
                {"complete.md": baseline, "duplicate.md": criterion_block(baseline, 1)},
                "criterion 1 is covered by multiple fragments",
                "added criterion 1 in a second copied fragment",
            )
        )
        cases.append(
            (
                "unrecognised-verdict",
                real_contract_text,
                {"complete.md": change_field(baseline, 1, "verdict", "maybe")},
                "unrecognised verdict 'maybe'",
                "changed criterion 1 verdict to maybe",
            )
        )
        cases.append(
            (
                "deferred-without-owner",
                real_contract_text,
                {"complete.md": change_field(change_field(baseline, 1, "verdict", "deferred"), 1, "notes", "deferred until phase P2")},
                "deferred row has no named owner",
                "made criterion 1 deferred without an owner",
            )
        )
        for field in ("ran", "at", "saw"):
            cases.append(
                (
                    f"missing-{field}",
                    real_contract_text,
                    {"complete.md": change_field(baseline, 1, field, None)},
                    f"row missing {field}",
                    f"removed criterion 1 {field} field",
                )
            )
        cases.append(
            (
                "invalid-timestamp",
                real_contract_text,
                {"complete.md": change_field(baseline, 1, "at", "not-a-timestamp")},
                "at value 'not-a-timestamp' is not a timestamp",
                "changed criterion 1 at value to not-a-timestamp",
            )
        )

        # This is a separate contract-copy probe: a stale fragment must not
        # silently become a report for a smaller, accidentally edited set.
        deleted_line = f"- [ ] {criteria[-1].text}"
        shortened_contract = real_contract_text.replace(deleted_line, "", 1)
        assert shortened_contract != real_contract_text
        cases.append(
            (
                "deleted-contract-criterion",
                shortened_contract,
                {"complete.md": baseline},
                "criterion 18 in",
                "deleted criterion 18 from the copied contract",
            )
        )

        for label, contract_text, fragments, expected, input_description in cases:
            contract, fragments_dir = make_case(root, label, contract_text, fragments)
            result = run_subprocess(contract, fragments_dir)
            output = (result.stdout + result.stderr).strip()
            assert result.returncode != 0, f"{label} unexpectedly passed:\n{output}"
            assert expected in output, f"{label} did not name the refusal:\n{output}"
            one_line = " | ".join(line.strip() for line in output.splitlines() if line.strip())
            print(f"acceptance selftest: {label}: {input_description}; {one_line}")

        complete_dir = root / "real-contract-complete" / "acceptance"
        complete_dir.mkdir(parents=True)
        (complete_dir / "synthetic.md").write_text(baseline, encoding="utf-8")
        result = run_subprocess(real_contract, complete_dir)
        output = (result.stdout + result.stderr).strip()
        assert result.returncode == 0, f"complete synthetic report failed:\n{output}"
        assert output.splitlines()[0] == "accepted", output
        assert f"criteria: {len(criteria)}" in output, output
        print(f"acceptance selftest: real contract derived {len(criteria)} criteria, naive count agrees; synthetic complete report: accepted")

    print("acceptance selftest: six refusal classes, all required-field variants, and completeness invariant pass")
    return 0


OPERATOR_FIXTURE = """# Operator fixture

## Rows

| # | Page | Note | Verdict |
|---|---|---|---|
| 1 | `/` | one qualified finding | held (via fixture) |
| 2 | `/` | clean | held |
| 3 | `/` | channel annotation, not a judgement | n/a |

## Probes

### Rows 1

verdict: held (via fixture)
ran: printf fixture
at: 2026-08-28T00:00:00Z
saw: fixture output
notes: A qualified verdict is intentional.

### Rows 2

verdict: held
ran: printf clean
at: 2026-08-28T00:00:00Z
saw: no annotation
notes: The clean page is recorded as held.

### Pages 4–15 — no annotations found

verdict: held
ran: printf absence
at: 2026-08-28T00:00:00Z
saw: no rows
notes: This is a valid zero-row absence block.

## Close

One of one templates walked. Three rows. One of the one pages drew no annotation.
"""


def make_operator_case(root: pathlib.Path, label: str, text: str) -> pathlib.Path:
    report_dir = root / label / "acceptance"
    report_dir.mkdir(parents=True)
    path = report_dir / "OPERATOR.md"
    path.write_text(text, encoding="utf-8")
    return path


def run_operator_subprocess(paths: list[pathlib.Path]) -> subprocess.CompletedProcess[str]:
    arguments = [
        sys.executable,
        str(pathlib.Path(__file__).resolve()),
        "--operator",
    ]
    for path in paths:
        arguments.extend(("--operator-report", str(path)))
    return subprocess.run(arguments, capture_output=True, text=True, check=False)


def operator_selftest() -> int:
    """Exercise each operator rule with a passing and a failing fixture."""
    with tempfile.TemporaryDirectory(prefix="siduri-operator-acceptance-") as raw:
        root = pathlib.Path(raw)
        passing = [
            ("coverage-pass", "coverage rule passing fixture"),
            ("agreement-pass", "verdict agreement passing fixture"),
            ("fields-pass", "field validation passing fixture"),
            ("verdict-pass", "verdict grammar passing fixture"),
            ("close-pass", "close arithmetic passing fixture"),
        ]
        for label, description in passing:
            path = make_operator_case(root, label, OPERATOR_FIXTURE)
            result = run_operator_subprocess([path])
            output = (result.stdout + result.stderr).strip()
            assert result.returncode == 0, f"{label} unexpectedly failed:\n{output}"
            print(f"operator selftest: {label}: {description}; accepted")

        breaches = [
            (
                "coverage-breach",
                OPERATOR_FIXTURE.replace("### Rows 1\n", "### Rows 9\n", 1),
                "uncovered",
                "changed a probe block to name no existing row",
            ),
            (
                "agreement-breach",
                OPERATOR_FIXTURE.replace(
                    "### Rows 1\n\nverdict: held (via fixture)",
                    "### Rows 1\n\nverdict: open",
                    1,
                ),
                "disagrees",
                "changed a block verdict without changing its table row",
            ),
            (
                "fields-breach",
                OPERATOR_FIXTURE.replace("ran: printf fixture", "ran: checked the fixture", 1),
                "not a command",
                "replaced ran with a description rather than a command",
            ),
            (
                "verdict-breach",
                OPERATOR_FIXTURE.replace(
                    "| 1 | `/` | one qualified finding | held (via fixture) |",
                    "| 1 | `/` | one qualified finding | maybe |",
                    1,
                ),
                "invalid verdict",
                "changed the table verdict to an unrecognised word",
            ),
            (
                "close-breach",
                OPERATOR_FIXTURE.replace("Three rows.", "Four rows.", 1),
                "close count 'rows'",
                "changed the close count without changing the table",
            ),
        ]
        for label, fixture, expected, description in breaches:
            path = make_operator_case(root, label, fixture)
            result = run_operator_subprocess([path])
            output = (result.stdout + result.stderr).strip()
            assert result.returncode != 0, f"{label} unexpectedly passed:\n{output}"
            assert expected in output, f"{label} did not name the refusal:\n{output}"
            one_line = " | ".join(line.strip() for line in output.splitlines() if line.strip())
            print(f"operator selftest: {label}: {description}; {one_line}")

    print("operator selftest: five passing fixtures and five named breach fixtures pass")
    return 0


def main(argv: list[str]) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--contract", type=pathlib.Path, default=pathlib.Path("docs/site-requirements.md"))
    parser.add_argument(
        "--fragments-dir",
        "--acceptance-dir",
        dest="fragments_dir",
        type=pathlib.Path,
        default=pathlib.Path("acceptance"),
    )
    parser.add_argument("--selftest", action="store_true")
    parser.add_argument("--operator", action="store_true")
    parser.add_argument("--operator-selftest", action="store_true")
    parser.add_argument("--operator-report", action="append", type=pathlib.Path)
    args = parser.parse_args(argv)
    if args.operator_selftest:
        try:
            return operator_selftest()
        except (AssertionError, AcceptanceInputError) as error:
            print(f"operator selftest: FAILED: {error}", file=sys.stderr)
            return 1
    if args.operator:
        return operator_check(args.operator_report or list(OPERATOR_REPORTS))
    if args.selftest:
        try:
            return selftest()
        except (AssertionError, AcceptanceInputError) as error:
            print(f"acceptance selftest: FAILED: {error}", file=sys.stderr)
            return 1
    return check(args.contract, args.fragments_dir)


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
