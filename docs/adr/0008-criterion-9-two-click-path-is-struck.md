# 0008 · Criterion 9's two-click path is struck

**Status** · accepted · 2026-08-28 · operator's word, given directly

## What changed

Under `13 · Acceptance criteria`, criterion 9 carried:

> Impressum and Datenschutz are reachable in two clicks from every page.

The criterion is struck in place. It follows the LR-1 amendment in ADR 0007; it does not amend LR-2, whose body specifies the contents of the Datenschutzerklärung and contains no link requirement.

The acceptance fragment heading at `acceptance/g3-markup-a11y.md:51` was changed to carry the same `~~ ~~` span. This is part of the amendment: `tools/acceptance.py:207` compares the normalised heading literally, and `normalise` at `:81-83` collapses whitespace and applies NFC but does not strip strike markers.

## Why it was asked

The operator's annotations were “Drop, add later” for both the `footer > a` “Impressum” and `footer > a` “Datenschutz”, followed by “Remove the footer for now”. The acceptance criterion was the check that encoded the path those footer links supplied, so it cannot remain as a held criterion after the operator's LR-1 ruling.

The required heading text is deliberately literal: the touched section is `13 · Acceptance criteria`, including the middle dot, rather than the shorthand “criterion 9”.

## What now carries the risk

I grepped `internal/site/pages_a5_test.go:46-53` and found the current positive two-click graph walk; `acceptance/g3-markup-a11y.md` records its prior green result. The shell lane's `TestA5RequiredPagesAndOperatorLinkGraph` is the corresponding inverted check and asserts that legal links are absent. `internal/site/pages_a5_test.go:21-34` still checks that both legal pages are built and routed. No check now carries a two-click acceptance obligation; the operator must ensure the legal pages and their surviving contents are suitable for publication.

The acceptance checker itself carries only the heading-consistency risk. Its red proof was a temporary strike without the heading update: `make acceptance` reported `criterion 9 in acceptance/g3-markup-a11y.md heading does not match contract`. The heading update removes that criterion mismatch. The command still reports the pre-existing forbidden `acceptance/OPERATOR.md: no criterion rows found`; that file is outside this lane's ownership and remains unedited.

`docs/FINDINGS.md:37`, row 0030, cites both LR-1's two-click path and NFR-3's orphan-page condition. Its stated grounds are stale after ADRs 0007 and 0009; the operator, who owns that findings record, owes its remeasurement or retirement. This lane does not edit `docs/FINDINGS.md`.
