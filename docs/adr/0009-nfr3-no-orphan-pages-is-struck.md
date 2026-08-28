# 0009 · NFR-3's no-orphan-pages phrase is struck

**Status** · accepted · 2026-08-28 · operator's word, given directly

## What changed

`NFR-3` carried this list of SEO properties:

> Clean semantic HTML, one `h1` per page, descriptive titles under 60 characters, meta descriptions from `summary`, JSON-LD (`Article`, `Person`, `SoftwareApplication`, `Comment`), canonical URLs, generated sitemap, no orphan pages.

Only the phrase `no orphan pages` is struck in place. The semantic HTML, heading, title, description, JSON-LD, canonical URL, and sitemap requirements remain unchanged.

## Why it was asked

The operator ruled “Remove, add later” for the About and Contact navigation items and “Drop, add later” for the Impressum and Datenschutz footer links, followed by “Remove the footer for now”. Those four pages become orphans when the footer and those navigation links are removed.

The operator did not rule on `/stack/`; its annotation was the open question “What is the difference from Tools?”. `/stack/` is orphaned by consequence of the footer decision, not by a decision to remove or defer that page. This amendment deliberately strikes only the `no orphan pages` phrase, not the surrounding NFR-3 sentence.

## What now carries the risk

I grepped `internal/site/pages_a5_test.go` and found required-output checks at `:21-34` and a legal-link graph at `:46-53`, but no general orphan-page assertion. I also found `tools/lane_overlap.py:240`, which checks that integrator-owned task paths are recognised by `is_path`; it does not inspect built-page reachability. There is therefore no automated NFR-3 orphan-page guard. The operator's navigation review carries the risk for the four ruled pages `/about/`, `/contact/`, `/impressum/`, and `/datenschutz/`; `/stack/` is an orphaned consequence awaiting the separate “difference from Tools” design decision.

`docs/FINDINGS.md:37`, row 0030, records both LR-1's two-click path and NFR-3's orphan-page condition as its grounds. Those grounds are stale after ADRs 0007 and 0009; the operator, who owns that findings record, owes its remeasurement or retirement. This lane does not edit `docs/FINDINGS.md`.
