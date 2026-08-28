# 0012 · FR-9's services line moves to P3

**Status** · accepted · 2026-08-28 · operator's word, given directly

## What changed

`FR-9` carried:

> Every article ends with the same three elements, in order: author line, newsletter capture, one contextual line routing to `/services`. No sales copy inside the post body.

The first sentence is struck in place and restated below it as two elements: author line and newsletter capture. The second sentence, “No sales copy inside the post body.”, remains untouched. The contextual services line moves to P3.

## Why it was asked

The operator annotated the article footer with “Newsletter capture is intentionally not part of this first build.” — “WTF”, and separately ruled “Remove the footer for now”. The launch contract cannot require a `/services` line while P3 is explicitly the phase for `/services` and the footer is being removed.

The whole first sentence is struck deliberately because changing “three” to “two” in place would silently rewrite surviving contract text. The corrected two-element requirement is a new line below the trace.

## What now carries the risk

I grepped `internal/site/article.templ:28-32` and found the author line, the current newsletter-capture placeholder, and the footer section slot. `internal/site/article_a1.go:62-70` registers that slot and documents that it is intentionally empty until the phase-3 offer exists. The implementation therefore has no contextual services line to remove today, and its newsletter text is an apology rather than a capture; the owning article lane still carries the two-element implementation risk.

The pre-P3 guard already carries the moved-route risk: `tools/pre_p3.py:21` rejects `/services`, and `Makefile:34` rejects `/services` in `dist/`. `tools/linkcheck.py:9-11` and `:25-27` carry an inert known-pending exemption for `/services/`, so a future route is not silently treated as a broken link; its selftest at `:253-266` fails when that pending route resolves unexpectedly. The no-sales-copy sentence remains guarded by the same output scans. No sales copy or services route is therefore permitted before P3, while the newsletter capture itself has no dedicated guard.
