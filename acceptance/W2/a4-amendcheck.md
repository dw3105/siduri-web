# W2 · A4 — the guard learns what a legitimate amendment looks like

base: 742e376
lane: A4 / lane/w2-amendcheck

## A4 — “the guard learns the shape of a legitimate amendment instead of refusing all of them”

did: added a stdlib-only Git-range checker with fresh-numbered ADR and dated-status validation, per-hunk `-U3` diff parsing, injective struck-fragment matching, all-level heading attribution, and the synthetic repository selftest (`tools/amendcheck.py:28`, `:140`, `:226`, `:261`, `:338`, `:431`, `:618`). Added this report as the only other owned file.

ran: `python3 tools/amendcheck.py --selftest`

saw: all nine cases passed. The selftest printed these decisive lines, in red-before-green order:

| case | exit | decisive output |
|---|---:|---|
| `delete-outright` | 1 | `rule 2 (struck traces): FAIL` |
| `strike-without-adr` | 1 | `rule 1 (added ADR): FAIL — no added docs/adr/NNNN-*.md file in the range` (rule 3 also fails because there is no ADR text) |
| `strike-wrong-section-name` | 1 | `rule 3 (named sections): FAIL` |
| `fragment-rewrite-real-591` | 0 | `decision: legitimate amendment; all three rules pass` |
| `unidentified-heading-real-526` | 1 | `rule 3 (named sections): FAIL` |
| `injectivity-one-span-for-ten-lines` | 1 | `rule 2 (struck traces): FAIL` |
| `injectivity-per-hunk` | 1 | `docs/site-requirements.md: hunk 2: FAIL — removed 10 line(s), added 0 non-empty strike span(s); injectivity count is too small` |
| `whole-line-strike` | 0 | `decision: legitimate amendment; all three rules pass` |
| `no-diff` | 0 | `decision: no contract diff; accepted` |

The fourth case is the real `:591` shape: it strikes `Huginn` and rewrites the rest of the table row. The fifth is the real `:526` shape under `## 13 · Acceptance criteria`; an ADR naming only `P0` is rejected because the full heading text `13 · Acceptance criteria` is the touched section name.

red proof: the first three planted inputs were an outright deletion with a valid ADR, a struck clause with no ADR, and a struck clause whose ADR names `OQ-5` instead of `P0`. They fired rules 2, 1, and 3 respectively. The two injectivity inputs then fired rule 2: one span for ten removed lines, and a separate second hunk with ten silent deletions despite ten spans in hunk 1. The per-hunk output above shows hunk 2 firing.

## Branch with no contract diff

ran: `python3 tools/amendcheck.py 742e376 HEAD`

saw:

```text
amendcheck: base=742e376 head=HEAD diff-context=3
rule 1 (added ADR): PASS — no contract diff; no amendment record is required
rule 2 (struck traces): PASS
  docs/site-requirements.md: PASS — no diff
  docs/comments-requirements.md: PASS — no diff
rule 3 (named sections): PASS
  docs/site-requirements.md: PASS — no diff
  docs/comments-requirements.md: PASS — no diff
decision: no contract diff; accepted
```

## Scratch branch: valid amendment and three breaks

The disposable clone used a scratch branch, restored the real Phrontis cell as its base, then struck that cell and added ADR 0008. The valid range was `d948d4c47ddcb8b135de8938d061200aa861d189..2c6cfe34b91c60c4fcc8f0c1a4aa3d07816aa437`.

ran: `python3 tools/amendcheck.py d948d4c47ddcb8b135de8938d061200aa861d189 2c6cfe34b91c60c4fcc8f0c1a4aa3d07816aa437`

saw:

```text
amendcheck: base=d948d4c47ddcb8b135de8938d061200aa861d189 head=2c6cfe34b91c60c4fcc8f0c1a4aa3d07816aa437 diff-context=3
rule 1 (added ADR): PASS — added ADR with a fresh number and dated status: docs/adr/0008-real-clause-amendment.md (**Status** · accepted · 2026-08-27 · operator's word, given directly)
rule 2 (struck traces): PASS
  docs/site-requirements.md: hunk 1: PASS — 1 removed line(s), 1 strike span(s), 1 injective match(es)
  docs/comments-requirements.md: PASS — no diff
rule 3 (named sections): PASS
  docs/site-requirements.md: PASS — added ADR names touched section(s): OQ-5
  docs/comments-requirements.md: PASS — no diff
decision: legitimate amendment; all three rules pass
```

ran: `python3 tools/amendcheck.py d948d4c47ddcb8b135de8938d061200aa861d189 HEAD` on each break branch.

saw, break rule 1 (ADR removed):

```text
amendcheck: base=d948d4c47ddcb8b135de8938d061200aa861d189 head=HEAD diff-context=3
rule 1 (added ADR): FAIL — no added docs/adr/NNNN-*.md file in the range
rule 2 (struck traces): PASS
  docs/site-requirements.md: hunk 1: PASS — 1 removed line(s), 1 strike span(s), 1 injective match(es)
  docs/comments-requirements.md: PASS — no diff
rule 3 (named sections): FAIL
  docs/site-requirements.md: FAIL — touched section name(s) missing from added ADR text: OQ-5
  docs/comments-requirements.md: PASS — no diff
decision: rejected amendment; failed rule 1, rule 3
```

saw, break rule 2 (struck line deleted):

```text
amendcheck: base=d948d4c47ddcb8b135de8938d061200aa861d189 head=HEAD diff-context=3
rule 1 (added ADR): PASS — added ADR with a fresh number and dated status: docs/adr/0008-real-clause-amendment.md (**Status** · accepted · 2026-08-27 · operator's word, given directly)
rule 2 (struck traces): FAIL
  docs/site-requirements.md: hunk 1: FAIL — removed 1 line(s), added 0 non-empty strike span(s); injectivity count is too small
  docs/comments-requirements.md: PASS — no diff
rule 3 (named sections): PASS
  docs/site-requirements.md: PASS — added ADR names touched section(s): OQ-5
  docs/comments-requirements.md: PASS — no diff
decision: rejected amendment; failed rule 2
```

saw, break rule 3 (ADR names `P0` instead of `OQ-5`):

```text
amendcheck: base=d948d4c47ddcb8b135de8938d061200aa861d189 head=HEAD diff-context=3
rule 1 (added ADR): PASS — added ADR with a fresh number and dated status: docs/adr/0008-real-clause-amendment.md (**Status** · accepted · 2026-08-27 · operator's word, given directly)
rule 2 (struck traces): PASS
  docs/site-requirements.md: hunk 1: PASS — 1 removed line(s), 1 strike span(s), 1 injective match(es)
  docs/comments-requirements.md: PASS — no diff
rule 3 (named sections): FAIL
  docs/site-requirements.md: FAIL — touched section name(s) missing from added ADR text: OQ-5
  docs/comments-requirements.md: PASS — no diff
decision: rejected amendment; failed rule 3
```

The scratch branches `scratch/w2-a4`, `scratch/w2-a4-break-rule1`, `scratch/w2-a4-break-rule2`, `scratch/w2-a4-break-rule3`, and `scratch/w2-a4-precedents` were deleted; the disposable clone was removed.

## Existing five strikes replayed

Each replay was a separate one-commit range in the disposable clone. All exited 0. `-U3` produced one hunk in each case.

| existing strike | checker range | quoted result |
|---|---|---|
| `docs/site-requirements.md:492` — P1 publication gate | `ac5c0e820f1a9d74b675089106acfea9f8ba5a34..0ef814eaa3cb5b5b8962ddb0a09739d868f716fc` | `docs/site-requirements.md: PASS — added ADR names touched section(s): P0`; `decision: legitimate amendment; all three rules pass` |
| `docs/site-requirements.md:568` — Phrontis cell | `ae6c879f101be8a053430f1d9c75caad6468dc3b..d510930bafc3061679d867e8da9d9aee8c7125d9` | `docs/site-requirements.md: PASS — added ADR names touched section(s): OQ-5`; `decision: legitimate amendment; all three rules pass` |
| `docs/site-requirements.md:591` — Huginn cell with the rest of the row rewritten | `d510930bafc3061679d867e8da9d9aee8c7125d9..0b75f2639aec1292c73d248301823b3442d0acc3` | `docs/site-requirements.md: hunk 1: PASS — 1 removed line(s), 1 strike span(s), 1 injective match(es)`; `decision: legitimate amendment; all three rules pass` |
| `docs/site-requirements.md:592` — Mímir and Forseti fragments | `0b75f2639aec1292c73d248301823b3442d0acc3..37ecc56b918b653a6de19651e9067dff7b9d858d` | `docs/site-requirements.md: PASS — 1 removed line(s), 2 strike span(s), 1 injective match(es)`; `decision: legitimate amendment; all three rules pass` |
| `docs/comments-requirements.md:333` — comment freeze | `37ecc56b918b653a6de19651e9067dff7b9d858d..ac5c0e820f1a9d74b675089106acfea9f8ba5a34` | `docs/comments-requirements.md: PASS — added ADR names touched section(s): 11 · Open questions`; `decision: legitimate amendment; all three rules pass` |

## Makefile handoff

The current line 30 behavior being replaced is:

```text
git diff --quiet "$$base" HEAD -- $$contract && git diff --quiet -- $$contract && git diff --cached --quiet -- $$contract
```

That first range comparison rejects every committed requirements diff. The two final comparisons reject a dirty working tree and dirty index, respectively.

Paste-ready replacement for the recipe line (not applied by this lane):

```make
@base="$$(git merge-base HEAD main 2>/dev/null || git rev-list --max-parents=0 HEAD)"; contract='docs/site-requirements.md docs/comments-requirements.md'; python3 tools/amendcheck.py "$$base" HEAD && git diff --quiet -- $$contract && git diff --cached --quiet -- $$contract || { echo 'check: the two requirements files are the contract and only the operator amends them' >&2; exit 1; }
```

`python3 tools/amendcheck.py` replaces only the range comparison and prints the three-rule decision. The working-tree and index checks remain in the same `&&` chain, so an accepted committed amendment with unstaged or staged contract edits still fails. The `AGENTS-PROHIBITIONS.md` guard on line 31 is untouched and still rejects modifications found in the merge-base range, working tree, or index. No Makefile or `docs/` file was changed.

## Limits

notes: this checker validates shape, never authority. An agent can add an ADR as easily as an operator; approval remains a human decision at the pull request.

notes: rule 2 ensures every removed line has its own struck span in its own `-U3` hunk, but it does not require a proportionate trace. Ten removed lines can each be answered by `~~the~~` if ten such spans are supplied. There is no cheap threshold that separates a small dishonest strike from a small honest one: this repository's accepted `:591` precedent is itself a one-word strike. The guard therefore enforces “no removed line is silent,” not “nothing was lost.”

notes: rule 1 checks that a new ADR has a fresh number and a `**Status** ·` line with a status and date; it does not resolve what that status line cites. ADR 0002 already demonstrates that a status can point at a file name that does not exist here.

notes: the checker pins `-U3` in its own `git diff` invocation. Without that pin, a user's `diff.context` configuration can move hunk boundaries, changing which strike is considered “in place” and making the same amendment's verdict irreproducible.

## Verification

ran: `go build ./...`
saw: exit 0, no output.

ran: `go test ./internal/site/`
saw: `ok github.com/dw3105/siduri-web/internal/site 0.770s`.

ran: `make build && find dist -type f -exec sha256sum {} + | sort > /tmp/siduri-a4-b1 && find dist -depth -type f -delete && find dist -depth -type d -empty -delete && make build && find dist -type f -exec sha256sum {} + | sort > /tmp/siduri-a4-b2 && diff /tmp/siduri-a4-b1 /tmp/siduri-a4-b2 && echo DETERMINISTIC`

saw: `DETERMINISTIC`.

ran: `git diff --check` and `git status --porcelain=v1` before the report was added.
saw: no diff-check output; the only pre-commit status entry was `?? tools/amendcheck.py`.

notes: the full `make check` gate was not run, per COMMON.md's three-lane instruction; the integrator runs it once on the merged tree.
