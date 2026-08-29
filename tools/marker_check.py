#!/usr/bin/env python3
"""Check enforcement-marker vocabulary and coverage in the three Siduri skills.

This is a shape check only.  It checks that every rule has exactly one marker
and that the marker belongs to the vocabulary for its skill.  It cannot check
whether a marker is true; that remains a judgement for a person.

Stdlib only.  Exit 1 on a missing, duplicate, or disallowed marker.
"""
from __future__ import annotations

import pathlib
import re
import sys
from dataclasses import dataclass


PERMANENT_MARKERS = frozenset({"bound", "bound-with-gap", "self-reported", "unbound"})
PENDING_MARKERS = frozenset({"pending-centralisation"})
SKILLS = (
    ("skills/siduri-code/SKILL.md", PERMANENT_MARKERS),
    ("skills/siduri-contract/SKILL.md", PERMANENT_MARKERS),
    ("skills/siduri-pending-tasks/SKILL.md", PENDING_MARKERS),
)
RULE_START = re.compile(
    r"^\s*(?:-\s+)?(?:\*\*)?(?P<rule>[A-Z]{2}-\d+[a-z]*)(?:\*\*)?\b"
)
MARKER = re.compile(r"\[([^\[\]\s]+)\]")


@dataclass(frozen=True)
class Violation:
    path: str
    line: int
    rule: str
    message: str

    def render(self) -> str:
        return f"  {self.path}:{self.line}: {self.rule}: {self.message}"


def check(root: pathlib.Path) -> tuple[int, list[Violation]]:
    rules = 0
    violations: list[Violation] = []
    for relative, allowed in SKILLS:
        path = root / relative
        try:
            lines = path.read_text(encoding="utf-8").splitlines()
        except OSError as error:
            violations.append(Violation(relative, 0, "file", f"cannot read skill: {error}"))
            continue

        for line_number, line in enumerate(lines, 1):
            rule_match = RULE_START.match(line)
            if rule_match is None:
                continue
            rules += 1
            rule = rule_match.group("rule")
            markers = MARKER.findall(line)
            if not markers:
                violations.append(
                    Violation(relative, line_number, rule, "has no enforcement marker; exactly one is required")
                )
                continue
            if len(markers) > 1:
                found = ", ".join(f"[{marker}]" for marker in markers)
                violations.append(
                    Violation(
                        relative,
                        line_number,
                        rule,
                        f"has {len(markers)} enforcement markers ({found}); exactly one is required",
                    )
                )
            for marker in markers:
                if marker not in allowed:
                    expected = ", ".join(f"[{value}]" for value in sorted(allowed))
                    violations.append(
                        Violation(
                            relative,
                            line_number,
                            rule,
                            f"marker [{marker}] is outside this skill's vocabulary (allowed: {expected})",
                        )
                    )
    return rules, violations


def main() -> int:
    root = pathlib.Path(__file__).resolve().parents[1]
    print(
        "marker-check: RESOLUTION ONLY — checks vocabulary and coverage; "
        "it cannot check whether any marker is true (that is a judgement)."
    )
    rules, violations = check(root)
    if violations:
        print(f"marker-check: FAIL — {len(violations)} marker violation(s) in {rules} rule line(s)", file=sys.stderr)
        for violation in violations:
            print(violation.render(), file=sys.stderr)
        return 1
    print(f"marker-check: PASS — {rules} rule line(s) carry exactly one permitted marker")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
