# 0010 · FR-5's reading-time phrase is struck

**Status** · accepted · 2026-08-28 · operator's word, given directly

## What changed

`FR-5` carried:

> Single column, ~68ch measure, no sidebar. Table of contents for posts over 1,500 words, inline on mobile and in the margin on wide screens. Reading time and publish date visible; updated date shown when it differs.

Only `Reading time and` is struck in place. The single-column measure, table-of-contents behavior, publish date, and conditional updated date remain. The tag display mentioned in the operator's annotation is not wording in `FR-5` and is outside this amendment.

## Why it was asked

The operator said “Date only, drop read-time and tags” for the root post metadata and repeated “Date only — drop read-time and tags” on the journal card. This amendment reaches the `FR-5` wording that requires reading time; it does not invent a contract strike for tags that `FR-5` does not contain.

## What now carries the risk

I grepped `internal/site/article.templ:13`, `internal/site/journal.templ:24`, `internal/site/pages_a5.templ:65`, and `internal/site/tags_a2.templ:18`. The first three render reading time and tags with the date; the tag-page card at `tags_a2.templ:18` is already date-only. I found no direct test asserting that reading time is visible, so no live guard carries the removed presentation requirement. The remaining date and article-layout obligations stay in the contract for the owning template lanes to implement and test.

The ambiguity was resolved narrowly: the operator's “drop read-time and tags” annotation is two presentation findings, but only reading time is named by this clause. Tags remain untouched here and any tag-display amendment needs its own operator ruling and trace.
