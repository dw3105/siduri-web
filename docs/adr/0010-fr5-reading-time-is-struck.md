# 0010 · FR-5's reading-time phrase is struck

**Status** · accepted · 2026-08-28 · operator's word, given directly

## What changed

`FR-5` carried:

> Single column, ~68ch measure, no sidebar. Table of contents for posts over 1,500 words, inline on mobile and in the margin on wide screens. Reading time and publish date visible; updated date shown when it differs.

Only `Reading time and` is struck in place. The single-column measure, table-of-contents behavior, publish date, and conditional updated date remain. The tag display mentioned in the operator's annotation is not wording in `FR-5` and is outside this amendment.

## Why it was asked

The operator's design pass recorded “Date only, drop read-time and tags” on the root post card and again on the journal card. **He revised that on 28 August 2026, asked directly: “Date and tags, read time is too posh.”** Tags stay; reading time goes, article pages included. This amendment reaches the `FR-5` wording that requires reading time, which is the contract's only mention of it; tags appear nowhere in `FR-5`, so that half of the original annotation never needed an amendment and now needs none.

## What now carries the risk

I grepped `internal/site/article.templ:13`, `internal/site/journal.templ:24`, `internal/site/pages_a5.templ:65`, and `internal/site/tags_a2.templ:18`. The first three render reading time and tags with the date; the tag-page card at `tags_a2.templ:18` is already date-only. I found no direct test asserting that reading time is visible, so no live guard carries the removed presentation requirement. The remaining date and article-layout obligations stay in the contract for the owning template lanes to implement and test.

Tags remain untouched by this clause and **no tag amendment is pending**: the operator's revision keeps them. `internal/site/chrome.templ` was corrected at `f785243` so `postMeta` renders date and tags without reading time, which is why the call-site note above no longer describes `pages_a5.templ`.
