# W2 · A8 — the seven amendments the operator ruled

base: 742e376
lane: A8 / lane/w2-amendments

## 0007 — `LR-1` — “the two-click path only”

did: Struck only the reachability sentence in `docs/site-requirements.md:375`; restated the working-email and direct-contact obligation at `:377`; added `docs/adr/0007-lr1-two-click-path-is-struck.md`.

ran: `python3 tools/amendcheck.py 301957a 7b1e9c6`

saw:

```text
rule 1 (added ADR): PASS — added ADR with a fresh number and dated status: docs/adr/0007-lr1-two-click-path-is-struck.md (**Status** · accepted · 2026-08-28 · operator's word, given directly)
rule 2 (struck traces): PASS
  docs/site-requirements.md: hunk 1: PASS — 1 removed line(s), 1 strike span(s), 1 injective match(es)
  docs/comments-requirements.md: PASS — no diff
rule 3 (named sections): PASS
  docs/site-requirements.md: PASS — added ADR names touched section(s): LR-1
  docs/comments-requirements.md: PASS — no diff
rule 4 (watched paths): PASS — all watched contract paths exist at base and HEAD
decision: legitimate amendment; all four rules pass
```

red proof: `python3 tools/amendcheck.py --selftest` planted and rejected `delete-outright` with `rule 2 (struck traces): FAIL`, `strike-without-adr` with rule 1 failing, and `strike-wrong-section-name` with rule 3 failing; the temporary fixtures were isolated and removed by the selftest.

notes: `LR-2` was inspected and not touched: its body has no `click`, `reachab`, or `link` requirement. `internal/site/pages_a5_test.go:46-53` still has the old positive graph assertion on this checkout; `lane/w2-shell` carries the inverted absence assertion. The measured cost is 31 built pages, 7 with direct legal footer links, and all 31 with a header link to `/`; the footer is the whole two-click path. `docs/FINDINGS.md` row 0030 is now stale and remains the operator's remeasurement/retirement work.

## 0008 — criterion 9 — “the acceptance criterion follows `LR-1`”

did: Struck criterion 9 in `docs/site-requirements.md:536`, added its ADR note at `:537`, changed only the matching heading at `acceptance/g3-markup-a11y.md:51`, and added `docs/adr/0008-criterion-9-two-click-path-is-struck.md`.

ran: `make acceptance` before 0008; temporary `make acceptance` after the contract strike and before the heading fix; `make acceptance` after 0008.

saw: Before 0008 and after the heading fix, the target exits 2 with:

```text
acceptance: refused:
  acceptance/OPERATOR.md: no criterion rows found
not accepted
criteria: 18
rows: 18
counts: held=8, deferred=5, open=3, failed=2, unrecognised=0
make: *** [mk/accept.mk:14: acceptance] Error 1
```

red proof: The temporary strike before the heading fix added the second failure:

```text
  criterion 9 in acceptance/g3-markup-a11y.md heading does not match contract: 'Impressum and Datenschutz are reachable in two clicks from every page.'
```

The unowned `acceptance/OPERATOR.md` failure remains; editing it is outside A8's file list.

ran: `python3 tools/amendcheck.py 7b1e9c6 f1ee71d`

saw:

```text
rule 1 (added ADR): PASS — added ADR with a fresh number and dated status: docs/adr/0008-criterion-9-two-click-path-is-struck.md (**Status** · accepted · 2026-08-28 · operator's word, given directly)
rule 2 (struck traces): PASS
  docs/site-requirements.md: hunk 1: PASS — 1 removed line(s), 1 strike span(s), 1 injective match(es)
  docs/comments-requirements.md: PASS — no diff
rule 3 (named sections): PASS
  docs/site-requirements.md: PASS — added ADR names touched section(s): 13 · Acceptance criteria
  docs/comments-requirements.md: PASS — no diff
rule 4 (watched paths): PASS — all watched contract paths exist at base and HEAD
decision: legitimate amendment; all four rules pass
```

red proof: The heading mismatch above is the affected-check red proof. `tools/acceptance.py` does not strip `~~`; the heading was aligned rather than changing the checker.

notes: The ADR contains the literal section name `13 · Acceptance criteria`. The old A5 graph evidence remains historical; the shell lane owns the inverted implementation check. The acceptance target cannot be green without changing forbidden `acceptance/OPERATOR.md`, so this lane leaves that pre-existing blocker visible.

## 0009 — `NFR-3` — “no orphan pages”

did: Struck only `no orphan pages` in `docs/site-requirements.md:303`, added the trace note at `:305`, and added `docs/adr/0009-nfr3-no-orphan-pages-is-struck.md`.

ran: `python3 tools/amendcheck.py f1ee71d d23649f`

saw:

```text
rule 1 (added ADR): PASS — added ADR with a fresh number and dated status: docs/adr/0009-nfr3-no-orphan-pages-is-struck.md (**Status** · accepted · 2026-08-28 · operator's word, given directly)
rule 2 (struck traces): PASS
  docs/site-requirements.md: hunk 1: PASS — 1 removed line(s), 1 strike span(s), 1 injective match(es)
  docs/comments-requirements.md: PASS — no diff
rule 3 (named sections): PASS
  docs/site-requirements.md: PASS — added ADR names touched section(s): NFR-3
  docs/comments-requirements.md: PASS — no diff
rule 4 (watched paths): PASS — all watched contract paths exist at base and HEAD
decision: legitimate amendment; all four rules pass
```

red proof: No direct NFR-3 orphan-page guard exists to turn red honestly. `internal/site/pages_a5_test.go` checks required output and legal reachability, while `tools/lane_overlap.py:240` checks task-path recognition, not built-page reachability.

notes: The four ruled pages are `/about/`, `/contact/`, `/impressum/`, and `/datenschutz/`. `/stack/` is orphaned by consequence, not by the operator's decision; its open annotation is “What is the difference from Tools?”. Row 0030 in `docs/FINDINGS.md` is stale and remains the operator's work.

## 0010 — `FR-5` — “reading time”

did: Struck only `Reading time and` in `docs/site-requirements.md:198`, added the trace note at `:200`, and added `docs/adr/0010-fr5-reading-time-is-struck.md`.

ran: `python3 tools/amendcheck.py d23649f 176ff31`

saw:

```text
rule 1 (added ADR): PASS — added ADR with a fresh number and dated status: docs/adr/0010-fr5-reading-time-is-struck.md (**Status** · accepted · 2026-08-28 · operator's word, given directly)
rule 2 (struck traces): PASS
  docs/site-requirements.md: hunk 1: PASS — 1 removed line(s), 1 strike span(s), 1 injective match(es)
  docs/comments-requirements.md: PASS — no diff
rule 3 (named sections): PASS
  docs/site-requirements.md: PASS — added ADR names touched section(s): FR-5
  docs/comments-requirements.md: PASS — no diff
rule 4 (watched paths): PASS — all watched contract paths exist at base and HEAD
decision: legitimate amendment; all four rules pass
```

red proof: No direct guard asserts visible reading time, so an honest red fixture does not exist. The operator's separate tag complaint is not wording in FR-5 and was not silently folded into this strike.

notes: Grep found reading-time and tag rendering in `internal/site/article.templ:13`, `internal/site/journal.templ:24`, and `internal/site/pages_a5.templ:65`; `internal/site/tags_a2.templ:18` is already date-only. The measure, table of contents, publish date, and conditional updated date remain.

## 0011 — `FR-16` — “the `mailto:` link”

did: Struck only the P0/P1 sentence in `docs/site-requirements.md:243`, added the deliberate-scope note at `:245`, left the phase-2 paragraph at `:247` unchanged, and added `docs/adr/0011-fr16-mailto-at-launch-is-struck.md`.

ran: `python3 tools/amendcheck.py 176ff31 c4d2553`

saw:

```text
rule 1 (added ADR): PASS — added ADR with a fresh number and dated status: docs/adr/0011-fr16-mailto-at-launch-is-struck.md (**Status** · accepted · 2026-08-28 · operator's word, given directly)
rule 2 (struck traces): PASS
  docs/site-requirements.md: hunk 1: PASS — 1 removed line(s), 1 strike span(s), 1 injective match(es)
  docs/comments-requirements.md: PASS — no diff
rule 3 (named sections): PASS
  docs/site-requirements.md: PASS — added ADR names touched section(s): FR-16
  docs/comments-requirements.md: PASS — no diff
rule 4 (watched paths): PASS — all watched contract paths exist at base and HEAD
decision: legitimate amendment; all four rules pass
```

red proof: No live guard requires a P0 contact path. `internal/site/pages_a5_test.go:21-34` requires the Contact page, but not its link or mailto; `tools/secretscan.py` only allowlists the reviewed address and files. There is no honest red fixture for the removed requirement.

notes: The current templates still contain the mailto until their owning navigation change lands. The ruled Phase 0 state has no contact path at all, not merely an unreachable Contact page. The phase-2 form paragraph and no-processing boundary survive.

## 0012 — `FR-9` — “the services line”

did: Struck the whole first sentence in `docs/site-requirements.md:212`, restated the two-element requirement at `:214`, left “No sales copy inside the post body.” standing, and added `docs/adr/0012-fr9-services-line-moves-to-p3.md`.

ran: `python3 tools/amendcheck.py c4d2553 f4e2ffc`

saw:

```text
rule 1 (added ADR): PASS — added ADR with a fresh number and dated status: docs/adr/0012-fr9-services-line-moves-to-p3.md (**Status** · accepted · 2026-08-28 · operator's word, given directly)
rule 2 (struck traces): PASS
  docs/site-requirements.md: hunk 1: PASS — 1 removed line(s), 1 strike span(s), 1 injective match(es)
  docs/comments-requirements.md: PASS — no diff
rule 3 (named sections): PASS
  docs/site-requirements.md: PASS — added ADR names touched section(s): FR-9
  docs/comments-requirements.md: PASS — no diff
rule 4 (watched paths): PASS — all watched contract paths exist at base and HEAD
decision: legitimate amendment; all four rules pass
```

red proof: `python3 tools/pre_p3.py --selftest` planted `/services/` in a temporary HTML artifact and rejected it as a forbidden artifact; `python3 tools/linkcheck.py --selftest` exercised the pending `/services/` exception and the failure when that route resolves. Both temporary fixtures were removed by their selftests.

notes: `tools/pre_p3.py:21` and `Makefile:34` actively forbid `/services` in the build. `tools/linkcheck.py:26` is an inert known-pending exemption. `internal/site/article_a1.go:62-70` registers an intentionally empty contextual slot, while the current article footer still has a newsletter placeholder rather than a capture. The moved services line is P3; the no-sales-copy rule remains.

## 0013 — amendment mechanism — “the escape from the old range guard”

did: Added `docs/adr/0013-amendment-guard-limits.md`, recording the `Makefile:30` mechanism change, operator authority, the eight stated limitations, and the shallow-checkout caveat. Added this report at `acceptance/W2/a8-amendments.md`.

ran: `python3 tools/amendcheck.py f4e2ffc HEAD`

saw:

```text
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

red proof: The mechanism's red cases are the `amendcheck.py --selftest` cases recorded under 0007 and the two independent shallow-checkout fixtures recorded in ADR 0013. CI occurrence of the old vacuous base is explicitly unproved.

notes: The amendment guard is not an authority system. `mk/w2amendcheck.mk` can disappear silently through `-include`; no file in this lane is allowed to fix that gap.

## Narrow readings recorded separately

- `LR-1`: the operator ruled on the footer's two-click path, not the legal page's contents, working email address, or direct-contact means. Because those obligations share a sentence, the whole sentence was struck and the surviving obligation restated.
- `LR-2`: the task's “seven” count does not make LR-2 an amendment. Its body has no click, reachability, or link language; it remains untouched.
- Criterion 9: the touched heading is the whole `13 · Acceptance criteria` heading name, including the middle dot, and the acceptance fragment heading was changed only to carry the same strike markers.
- `NFR-3`: only the `no orphan pages` phrase was struck. The first four pages were ruled; `/stack/` was not, and its orphaning is recorded as a consequence rather than a decision.
- `FR-5`: only `Reading time and` was struck. The operator also mentioned tags, but tags are not in this clause and were not edited.
- `FR-16`: only the P0/P1 mailto sentence was struck; the phase-2 form paragraph was left standing rather than treating the multi-line clause as one target.
- `FR-9`: the entire first sentence was struck so `three` did not become `two` through an untraced in-place rewrite; the two-element wording is a new line below the trace.

## Final verification

ran: `go build ./... && go test ./internal/site/ -count=1`

saw: `go build` exit 0; `go test` exit 0 with `ok github.com/dw3105/siduri-web/internal/site`.

ran: `make build && find dist -type f -exec sha256sum {} + | sort > /tmp/b1`; `find dist -mindepth 1 -depth -delete`; `make build && find dist -type f -exec sha256sum {} + | sort > /tmp/b2`; `diff /tmp/b1 /tmp/b2 && echo DETERMINISTIC`

saw: `DETERMINISTIC`.

ran: `make check`

saw: exit 0; build, pre-P3, HTML structure, budgets, headers, internal/external-link selftests, axe on 31 pages, secretscan, workflow validation, the six-amendment range check, templ formatting, and all Go tests passed.

ran: `make acceptance`

saw: exit 2 at `mk/accept.mk:14` because `acceptance/OPERATOR.md: no criterion rows found`; criterion 9's heading mismatch is fixed. This target cannot be green within A8's ownership because editing `acceptance/OPERATOR.md` is forbidden.

ran: `git status --porcelain`

saw: no output after the seven commits.
