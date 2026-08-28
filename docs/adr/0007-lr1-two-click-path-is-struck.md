# 0007 · LR-1's two-click path is struck

**Status** · accepted · 2026-08-28 · operator's word, given directly

## What changed

`LR-1` carried this sentence:

> Must be reachable in two clicks from every page, with a working email address and a means of direct contact.

The sentence is struck in place. The working email address and means of direct contact remain an explicit obligation in the contract; only the two-click reachability path was amended. `LR-2` was not amended: its body contains the Datenschutzerklärung content specification and no click, reachability, or link requirement.

## Why it was asked

The operator's annotation was `footer > a` “Impressum” — “Drop, add later”, and the footer annotation was “Remove the footer for now”. The ruling concerns the footer path, not the legal contents of the Impressum or the separate obligation to provide a working email address and a means of direct contact.

This is a deliberately whole-sentence strike because the reachability and contact obligations were written as one sentence. Striking only the reachability phrase would leave the surviving clause beginning with “with a working email address”; repairing that seam in place would silently rewrite surviving contract text.

## What now carries the risk

I grepped `internal/site/pages_a5_test.go:46-53` and found the current positive graph assertion that walks every built page to `impressum/index.html` and `datenschutz/index.html` within two clicks. The same file still requires both legal pages at `:21-34`. On `lane/w2-shell`, that graph assertion is already inverted to `TestA5RequiredPagesAndOperatorLinkGraph`, which checks that the legal links are absent; the shell lane owns that test and this amendment does not edit it.

The measured build has 31 pages, 7 pages with direct legal links in their footer, and all 31 pages with a header link to `/`. Therefore the footer is the whole two-click path: the other 24 pages reach the legal pages through the home page's footer, and removing that footer removes the path for all 31 pages. The legal pages remain built and routed, but no repository check carries the removed reachability obligation after the shell lane's inverted test.

**This ADR does not cite the `mailto:` as carrying the surviving direct-contact obligation, and an earlier draft did.** ADR 0011, later in this same amendment set, strikes the clause that required that `mailto:` and records that nothing in the tree replaces it. So the surviving § 5 DDG obligation to offer a means of direct contact is **unmet in the Phase 0 state this set produces**, and that is the operator's to resolve rather than a drafting matter. Part three of this amendment — the check edit `CT-02` requires — landed separately in `baaface` on `lane/w2-shell`, which deleted `a5ReachableWithin` outright rather than inverting it; the capability to test two-click reachability no longer exists, so restoring this clause would not restore its check. German legal obligations still apply; the operator must publish legally compliant pages with the surviving direct-contact details.

Three records outside this lane's grant are made stale by this amendment set, and they are named together because tracking one and not the others implies the others were checked:

- `docs/FINDINGS.md:37`, row 0030, cites both LR-1's two-click path and NFR-3's orphan-page condition as its grounds.
- `acceptance/OPERATOR.md:21` and `:40` record *Date only, drop read-time and tags*, half of which the operator withdrew on 28 August 2026.
- `acceptance/g3-markup-a11y.md`'s criterion 9 row keeps `verdict: held` with `ran: go test ./internal/site -run '^TestA5RequiredPagesAndLegalLinkGraph$'`. That test is absent on `lane/w2-shell`, and `go test -run` on a pattern matching nothing prints `ok … [no tests to run]` and exits 0. After the merge the row re-verifies itself by executing nothing, under a `held` verdict for a criterion this set strikes.

The operator owns all three. This lane edits none of them.
