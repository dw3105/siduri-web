#!/usr/bin/env python3
"""Check that acceptance fragments cover the current site contract.

The contract is the source of truth for the criterion set.  Acceptance
fragments may say that a criterion is not met, but they must identify every
criterion exactly once and include enough evidence for the verdict to be
audited.

    python3 tools/acceptance.py
    python3 tools/acceptance.py --selftest

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

    paths = sorted(fragments_dir.glob("*.md"))
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
    args = parser.parse_args(argv)
    if args.selftest:
        try:
            return selftest()
        except (AssertionError, AcceptanceInputError) as error:
            print(f"acceptance selftest: FAILED: {error}", file=sys.stderr)
            return 1
    return check(args.contract, args.fragments_dir)


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
