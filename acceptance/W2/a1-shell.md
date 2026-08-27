Precondition: Do not merge `lane/w2-shell` until the integrator's ADRs covering the intentional contract contradictions land before or with this branch.

# W2 · A1 — one chrome, one source

base: 742e376
lane: A1

## / #9 and /journal/ #5 — “Remove the footer for now”

did: Added `chromeDocument` as the single doctype/head/body/content shell with no site footer, and made `Document` and `a5Document` thin wrappers that preserve their existing signatures (`internal/site/chrome.templ:3-23`, `internal/site/templates.templ:3-4`, `internal/site/pages_a5.templ:3-4`). Generated `chrome_templ.go`, `templates_templ.go`, and `pages_a5_templ.go` with templ.
ran: `grep -c '<footer' internal/site/templates.templ internal/site/pages_a5.templ || true`
saw: `internal/site/templates.templ:0` and `internal/site/pages_a5.templ:0`.
red proof: The pre-change `TestA5RequiredPagesAndLegalLinkGraph` was the planted input; after the footer removal it printed `cannot reach impressum/index.html in two clicks` and `cannot reach datenschutz/index.html in two clicks` for the rendered pages. The test was then restored to the operator contract below.
notes: `article.templ:28` remains untouched as the separately owned article footer.

## / #7 and / #8 — “Remove, add later”

did: Added `siteHeader` as the only shared site header/primary nav, retaining Journal and Tools and removing About and Contact (`internal/site/chrome.templ:25-33`). Both document wrappers call it through `chromeDocument`.
ran: `if grep -rn 'About\|Contact' internal/site/chrome.templ; then exit 1; else echo 'no About/Contact in chrome nav'; fi`
saw: `no About/Contact in chrome nav`.
red proof: No focused project test enforces the primary-nav vocabulary. The stale golden snapshots indirectly cover the old shell and now fail; `tools/linkcheck.py` checks that emitted links resolve, not which navigation links exist. The source grep above is the acceptance guard.
notes: About, Contact, Impressum, Datenschutz, and Stack remain built and routed. No Stack link was added.

## / #6 and /journal/ #1 — “Date only, drop read-time and tags”

did: Added `postMeta(post Post)` as the one date-only component, preserving `<time datetime>` and removing read-time/tags from the A5 card (`internal/site/chrome.templ:35-37`, `internal/site/pages_a5.templ:28-33`). The A5 home tagline at the former line 48 was descriptive copy, not post metadata, so it remains an ordinary `<p>` without the `post-meta` class. `readingTime` was not deleted.
ran: `rg -n 'post-meta|postMeta' internal/site/chrome.templ internal/site/templates.templ internal/site/pages_a5.templ`
saw: `postMeta` is called at `internal/site/pages_a5.templ:30`; the only `post-meta` markup among these owned sources is the date-only component at `internal/site/chrome.templ:36`.
red proof: The rendered `dist/index.html` card contains `<p class="post-meta"><time datetime="2026-08-25">2026-08-25</time></p>` and no `min read` or tag list. No separate metadata-only guard was added; the existing golden tests are stale for the shared-shell change and are reported below.
notes: `article.templ`, `journal.templ`, and `tags_a2.templ` still use their old metadata and belong to round B.

## Dominant finding — shared page header

did: Added `pageHeader(class, meta, title, summary, beforeTitle, afterSummary)` in `internal/site/chrome.templ:39-47`. It owns the shared header shape while allowing article updated/draft content before or after the common title/summary; round-B article and tools lanes can pass their own slots without forking the block.
ran: `go build ./...`
saw: exit 0 with no output.
red proof: No page-header consumer was changed in A1 because `article.templ` and `tools_a6.templ` are outside this lane. There is no honest red proof for a future consumer until round B calls the component.
notes: No article or tools template was edited.

## Contract contradiction — intentional operator ruling

The branch must not merge before the integrator's ADRs land. This change intentionally leaves five orphan routes: `/about/`, `/contact/`, `/impressum/`, `/datenschutz/`, and `/stack/`. The first four were ruled unlinked during the sitting. `/stack/` is the fifth orphan and the one nobody ruled: its only prior route was the removed footer, and the open question “What is the difference from Tools?” remains with design. The page was not deleted and no Stack link was added.

- `LR-1`: Impressum remains built but is no longer reachable in two clicks. Check that goes red: **none**; `tools/linkcheck.py` only validates links that exist.
- `LR-2`: Datenschutz has the same intentional reachability break. Check that goes red: **none**; no orphan/legal-link graph check exists.
- `NFR-3`: the no-orphan-pages requirement is violated by five routes. Check that goes red: **none**; `NFR-3` says no orphan pages and nothing in `make check` enforces it, so this branch orphans two mandatory legal pages and the gate stays green.
- `FR-5`: the A5 post cards no longer show reading time, by ruling. Check that goes red: **none**; no gate enforces reading time.
- `FR-16`: Contact remains built but unlinked, leaving Phase 0 with no contact path. Check that goes red: **none**; `mailto:` is skipped by the link checker and no contact-path assertion exists.
- Acceptance criterion “Impressum and Datenschutz are reachable in two clicks from every page”: Check that goes red: **none**; the owned test was changed because it asserted the behavior the operator removed.

The changed test is `TestA5RequiredPagesAndOperatorLinkGraph` in `internal/site/pages_a5_test.go`. Before this lane it was `TestA5RequiredPagesAndLegalLinkGraph` and required every generated page to reach both legal pages within two clicks. It now keeps the required outputs and asserts that neither legal URL is linked, matching the sitting’s ruling.

## CSS handoff

The integrator should paste these changes into `static/site.css`:

```css
/* Replace the existing header, main, footer grouping. */
header, main {
  margin: 0 auto;
  max-width: 68ch;
  padding: 1.25rem;
}

/* Replace the existing article-footer rule so the article footer keeps its typography. */
.article-footer {
  border-top: 1px solid var(--rule);
  color: var(--muted);
  font-size: .9rem;
  margin-top: 4rem;
  padding-top: 1rem;
}

/* Delete the generic shell-footer rule. */
/* footer { border-top: 1px solid var(--rule); color: var(--muted); font-size: .9rem; } */
```

The primary `nav` styling remains in `static/site.css:27-28` (`nav` and `nav a`). `static/article_a1.css` owns the other navigation landmarks (`.article-toc` at lines 6-30 and `.series-navigation` at lines 85-110); there is no literal primary-nav selector there. A1 did not edit that round-B-owned stylesheet.

`static/tools_a6.css` has no literal `footer` selector. It retains `.tool-card-footer` at lines 47-57, which is still used by tool cards and is not the removed document footer; A1 leaves that round-B tools decision alone.

## Checks and output

did: Built before and after the change; both builds emitted 31 HTML pages. Opened the rendered bodies of `dist/index.html` and `dist/journal/hello-siduri/index.html`.
ran: `make build && find dist -name '*.html' | wc -l`
saw: `31` before the change and `31` after the change.

ran: `grep -o '<header>.*</header>' dist/index.html`
saw: `<header><a href="/" aria-label="Siduri home"><strong>Siduri</strong></a><nav aria-label="Primary"><a href="/journal/">Journal</a> <a href="/tools/">Tools</a></nav></header>`.

ran: `go test ./internal/site -run '^TestA5RequiredPagesAndOperatorLinkGraph$' -count=1`
saw: `ok   github.com/dw3105/siduri-web/internal/site`.

ran: `python3 tools/linkcheck.py --dist dist`
saw: `373 internal links, 4 external links, 0 failures, internal only; external checks skipped (use --external)`.

ran: `python3 tools/paths_exist.py AGENTS.md`
saw: `paths_exist: 1 file(s), every named path resolves.` Neither script asserts footer links.

ran: `PATH="$(go env GOPATH)/bin:$PATH" templ generate --check`
saw: templ completed with `updates=0`; generated files are current.

ran:

```sh
make build && find dist -type f -exec sha256sum {} + | sort > /tmp/b1
dist_hold=$(mktemp -d /tmp/siduri-w2-dist.XXXXXX)
mv dist "$dist_hold/first"
make build && find dist -type f -exec sha256sum {} + | sort > /tmp/b2
diff /tmp/b1 /tmp/b2 && echo DETERMINISTIC
find dist -name '*.html' | wc -l
```

saw: `DETERMINISTIC`; the prescribed `rm -rf dist` form was rejected by the execution safety layer, so `dist/` was moved rather than deleted.

## Out-of-scope test handoff

ran: `go test ./internal/site/`
saw: FAIL only in `TestGoldenPages` for `long-post.html`, `short-post.html`, and `empty-tag.html`, plus `TestA2TagFanoutCoversVocabularyAndSortsPosts` because it compares the empty-tag component to the same stale golden. Each snapshot still contains the old About/Contact nav and footer. The snapshots are under `testdata/golden/`, outside A1’s owned file list, so they were not edited; the integrator must regenerate them after merging.

notes: `go build ./...`, the updated A5 test, linkcheck, templ freshness, and the two-build checksum pass. `static/site.css`, `static/article_a1.css`, `static/tools_a6.css`, `article.templ`, `journal.templ`, `tags_a2.templ`, `tools_a6.templ`, `testdata/golden/`, and `docs/` were not edited.
