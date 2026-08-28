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

9. **Rule 1 is existential over the range, and every per-ADR check it owns
   inherits that.** `amendcheck.py:393` filters to the added ADRs with no
   `date_error` and passes if *any* survive. The number-freshness checks at
   `:259-275` — collision with an existing ADR, duplication inside the added set
   — are well built and write to that same field, so they are discarded whenever
   one added ADR is clean. Six of seven ADRs can carry a malformed, undated or
   duplicated status line and the guard reports a legitimate amendment while
   naming the one it validated, which reads exactly like set validation. Audit
   after the fact instead:

       ls docs/adr/ | sed -n 's/^\([0-9]\{4\}\)-.*/\1/p' | sort | uniq -d
       grep -L '^\*\*Status\*\* ·' docs/adr/*.md

10. **Part three of a three-part amendment is never in scope.** `CT-02` requires
    the ADR, the strike in place, and every affected check, in one commit. Rules
    1 to 3 read `docs/` only, so the guard prints *legitimate amendment* over a
    range that cannot contain a check edit — and prints the same whether that
    edit exists anywhere or nowhere. Both instances in this set are real: this
    amendment's own part three landed on another branch, in `baaface`, and
    criterion 9's landed nowhere until it was noticed by hand. The first was
    visible only because `make acceptance` happens to compare fragment headings
    against contract text.

11. **One amendment's risk assessment can be invalidated by another amendment in
    the same set, and nothing sees it.** ADR 0007 retains `LR-1`'s direct-contact
    obligation; ADR 0011 strikes the clause requiring the only implementation of
    it. Each passes the guard on its own range, and a per-commit read cannot
    catch it either — it took reading the whole set. **The remedy is not a future
    guard.** Rules 1 to 4 read `docs/` and one range, and nothing cheap
    cross-references two ADRs' risk sections against each other. A person reads
    the set, or nobody does.

The old guard may also have been vacuous in CI: `Makefile:30` fell back to a
base expression that collapses to HEAD in a shallow checkout, while
`actions/checkout@v4` defaults to `fetch-depth: 1` and the old `ci.yml` set no
depth. The current workflow now requests full history, but that change proves
the mechanism, not whether the old failure occurred in a real CI run.
