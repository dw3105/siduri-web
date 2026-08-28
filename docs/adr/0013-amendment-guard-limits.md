# 0013 · The amendment guard is a shape check, not authority

**Status** · accepted · 2026-08-28 · operator's word, given directly

## What changed

Before this wave, `Makefile:30` used a range comparison that rejected every
committed requirements diff:

> `git diff --quiet "$$base" HEAD -- $$contract`

The operator chose the escape from four options on 27 August 2026. `301957a`
moved the range judgement to `mk/w2amendcheck.mk`, where the existing working
tree checks remain in `Makefile:30` and the range checker now invokes
`tools/amendcheck.py`. The mechanism has no requirements diff of its own.

## Why it was asked

`main` requires `make check`, main has no bypass, and the old range expression
made an operator-approved contract amendment impossible to merge. A legitimate
amendment therefore needed a mechanical shape check: a fresh ADR with a dated
status, struck traces in the same diff hunk, named touched sections, and both
contract files present at both ends of the range.

CT-01 is satisfied because the operator holds the merge, not because this lane
holds an exemption. His ruling covers the six strikes; the replacement wording
below each strike is this lane's wording and remains his to reject during
review.

## What now carries the risk

The pull request and the operator carry authority. `tools/amendcheck.py` carries
only the shape of an amendment, and `mk/w2amendcheck.mk` carries the range
invocation. The four-rule selftest and two independent shallow-checkout
fixtures proved the mechanism's stated behavior. The occurrence of a vacuous
old-guard pass in CI is unproved: nobody has read a CI log.

The mechanism has these limits, all recorded rather than hidden:

1. It validates shape, never authority. An agent writes an ADR as easily as an
   operator does.
2. Every removed line leaves a trace, but the trace is not proportionate. No
   cheap rule separates a small honest strike from a small dishonest one,
   because this repository's own accepted precedent is a small strike.
3. Rule 1 checks that a `**Status** ·` line is present, never that what it cites
   resolves. ADR 0002 names a file that does not exist.
4. Rules 2 and 3 do not see commit boundaries. The range is one diff, clustered
   clauses can share a hunk, and added ADRs are unioned, so splitting an
   amendment across commits makes pairing legible rather than enforced.
5. `added_adr_paths` uses `--diff-filter=A` with `--no-renames`; renaming an
   existing ADR to an unused number satisfies rule 1 without writing anything.
6. The guard cannot check that an ADR carries the reason and replacement that
   CT-02 requires. A correctly shaped ADR with no content passes.
7. Rule 2 cannot see edits to surviving text. `check_hunk` at
   `tools/amendcheck.py:314-335` builds `removed` from removed lines and
   `spans` from added lines and matches spans into removed lines by substring;
   it never compares the added line's non-struck remainder with the line it
   replaces. An amendment can silently rewrite arbitrary surviving prose while
   one struck span remains a substring of the original. This is a hole in the
   rule the guard exists for, not in what it reads, and belongs with the
   proportionality limit above.
8. `check: w2-amendcheck` lives in `mk/w2amendcheck.mk`, reached through
   `-include`, which is silent by construction. Deleting that file drops the
   prerequisite and makes `make -n check` byte-identical to the run with it
   present. Nothing asserts that `check` still depends on it. The fix is not
   part of this lane.

The old guard may also have been vacuous in CI: `Makefile:30` fell back to a
base expression that collapses to HEAD in a shallow checkout, while
`actions/checkout@v4` defaults to `fetch-depth: 1` and the old `ci.yml` set no
depth. The current workflow now requests full history, but that change proves
the mechanism, not whether the old failure occurred in a real CI run.
