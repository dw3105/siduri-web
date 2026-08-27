# W2 · A6 — a watched contract path must exist on both sides

base: 742e376
lane: A6 / lane/w2-rename

## A6 — “a watched contract path must exist at the base revision AND at HEAD; if either is false, refuse”

did: added rule 4, which checks every watched contract path with `git cat-file -e` at both revisions, reports the missing path and side, and makes the four-rule decision include that result (`tools/amendcheck.py:189`, `:199`, `:480`, `:505`). Pinned both checker diffs with `--no-renames` (`tools/amendcheck.py:216`, `:469`). Added safe handling for a watched path absent at the base so rule 4 reports the failure instead of attempting to read a nonexistent file (`tools/amendcheck.py:482`).

ran: `python3 tools/amendcheck.py --selftest`

saw: `amendcheck selftest: 13 cases pass` — all nine A4 cases plus `rename-watched-path-away`, `delete-watched-path-outright`, `add-watched-path-new`, and `both-watched-paths-present`.

red proof: the planted temporary cases `rename-watched-path-away`, `delete-watched-path-outright`, and `add-watched-path-new` exited 1. Their rule-4 output named, respectively, `docs/site-requirements.md missing at HEAD`, `docs/site-requirements.md missing at HEAD`, and `docs/site-requirements.md missing at base`. The fixtures were temporary-directory inputs and were removed after each case; the positive `both-watched-paths-present` case then printed `rule 4 (watched paths): PASS — all watched contract paths exist at base and HEAD`.

## e110d0b — historical rename replay

ran: `python3 tools/amendcheck.py e110d0b^ e110d0b`

saw: exit 1. The real amendment was refused with:

```text
rule 4 (watched paths): FAIL — watched path check failed: docs/site-requirements.md missing at base (e110d0b^); docs/comments-requirements.md missing at base (e110d0b^)
decision: rejected amendment; failed rule 3, rule 4
```

The missing base paths are the rename `{REQUIREMENTS-SITE.md => site-requirements.md}` and `{REQUIREMENTS-COMMENTS.md => comments-requirements.md}` in `e110d0b`; rule 4 rejects the rename without trying to understand it.

## Branch with no contract diff

ran: `python3 tools/amendcheck.py 742e376 HEAD`

saw: exit 0.

```text
amendcheck: base=742e376 head=HEAD diff-context=3
rule 1 (added ADR): PASS — no contract diff; no amendment record is required
rule 2 (struck traces): PASS
  docs/site-requirements.md: PASS — no diff
  docs/comments-requirements.md: PASS — no diff
rule 3 (named sections): PASS
  docs/site-requirements.md: PASS — no diff
  docs/comments-requirements.md: PASS — no diff
rule 4 (watched paths): PASS — all watched contract paths exist at base and HEAD
decision: no contract diff; accepted
```

## Verification

ran: `go build ./...`

saw: exit 0, no output.

ran: `go test ./internal/site/`

saw: `ok  	github.com/dw3105/siduri-web/internal/site	0.942s`.

ran: `make build && find dist -type f -exec sha256sum {} + | sort > /tmp/siduri-a6-b1 && find dist -depth -type f -delete && find dist -depth -type d -empty -delete && make build && find dist -type f -exec sha256sum {} + | sort > /tmp/siduri-a6-b2 && diff /tmp/siduri-a6-b1 /tmp/siduri-a6-b2 && echo DETERMINISTIC`

saw: `DETERMINISTIC`.

notes: rule 4 refuses a rename rather than understanding one. A legitimate rename of a contract file is therefore impossible without changing the checker; that is the intended tradeoff. The checker still validates amendment shape, not human authority. The full `make check` gate was not run, per COMMON.md; the integrator runs it on the merged tree.
