# W2 · A2 guards

base: 742e376
lane: A2

## Criterion 8 — “A draft post appears nowhere in dist/, verified by grep.”

did: Added source-aware draft loading and output scanning for slug paths, title, summary, plain_summary, and rendered body fragments in `tools/draftscan.py:116-202`; added clean and breach fixtures to its self-test at `tools/draftscan.py:224-253`; wired the real and self-test targets through `mk/w2guards.mk:3-11`.
ran: `python3 tools/draftscan.py --selftest`
saw: `draftscan selftest: clean case pass; slug path -> draftscan: /tmp/siduri-draftscan-6297x4si/dist/journal/selftest-draft-slug/index.html:1: draft slug in path leaked ('selftest-draft-slug'); title -> draftscan: /tmp/siduri-draftscan-6297x4si/dist/title.html:1: draft title leaked ('Selftest Draft Title'); summary -> draftscan: /tmp/siduri-draftscan-6297x4si/dist/summary.html:1: draft summary leaked ('Selftest draft summary'); plain_summary -> draftscan: /tmp/siduri-draftscan-6297x4si/dist/plain.html:1: draft plain_summary leaked ('Selftest draft plain summary'); body -> draftscan: /tmp/siduri-draftscan-6297x4si/dist/body.html:1: draft body text leaked ('Selftest draft body sentinel.')`
ran: `python3 tools/draftscan.py dist`
saw: `draftscan: checked 0 draft post(s) against dist; no draft content found`
red proof: Added temporary `content/posts/w2-a2-redproof.md` with `draft: true`; `make build` excluded its route. After injecting `dist/journal/w2-a2-real-draft-red-proof-20260827/index.html`, `python3 tools/draftscan.py dist` printed `draftscan: dist/journal/w2-a2-real-draft-red-proof-20260827/index.html:1: draft slug in path leaked ('w2-a2-real-draft-red-proof-20260827')`, `...:2: draft title leaked ('W2 A2 Real Draft Red Proof')`, and `...:4: draft body text leaked ('W2 A2 real draft scanner red proof body 20260827.')`; the command exited 1. The fixture and injected output were removed, then the build was restored.
notes: The guard only treats frontmatter drafts found in `content/posts/` as sources and checks their known metadata plus rendered body text; it does not ban the literal word `draft` in unrelated published content.

## Criterion 10 — “No cookie banner exists, and no non-essential client-side storage is written.”

did: Added emitted-file scanning for `localStorage`, `sessionStorage`, `indexedDB`, `document.cookie`, `Set-Cookie` in meta tags or `_headers`, and cookie-banner class/id vocabulary in `tools/storagescan.py:38-100`; added clean, JSON-LD script, and breach fixtures to its self-test at `tools/storagescan.py:114-137`; wired it beside the draft guard in `mk/w2guards.mk:3-11`.
ran: `python3 tools/storagescan.py --selftest`
saw: `storagescan selftest: clean case and script allowance pass; localStorage -> storagescan: /tmp/siduri-storagescan-sueu7we6/dist/storage.html:1: forbidden localStorage; sessionStorage -> storagescan: /tmp/siduri-storagescan-sueu7we6/dist/storage.html:1: forbidden sessionStorage; indexedDB -> storagescan: /tmp/siduri-storagescan-sueu7we6/dist/storage.html:1: forbidden indexedDB; document.cookie -> storagescan: /tmp/siduri-storagescan-sueu7we6/dist/storage.html:1: forbidden document.cookie; Set-Cookie header -> storagescan: /tmp/siduri-storagescan-sueu7we6/dist/_headers:2: forbidden Set-Cookie in _headers; Set-Cookie meta -> storagescan: /tmp/siduri-storagescan-sueu7we6/dist/meta.html:1: forbidden Set-Cookie in meta tag; cookie banner -> storagescan: /tmp/siduri-storagescan-sueu7we6/dist/banner.html:1: forbidden cookie-banner markup (gdpr-consent)`
ran: `python3 tools/storagescan.py dist`
saw: `storagescan: checked emitted files in dist; no storage or cookie-banner markers found`
red proof: The self-test planted `<script>localStorage.setItem("x", "y")</script>` in `storage.html`, `Set-Cookie: siduri=1` in `_headers`, `<meta http-equiv="Set-Cookie" ...>` in `meta.html`, and `id="gdpr-consent"` in `banner.html`; each produced the offending file and line shown above and was removed with the temporary tree.
notes: This scans only bytes emitted into `dist/`; it cannot see storage written by a script this build does not emit, and it cannot see a cookie set by the host. This build emits no runtime JavaScript today, which is why the criterion currently holds; the scanner matters when that changes. `acceptance/OPERATOR.md` records the annotation tool injected into the `.dev-dist` preview surface; `find dist -type f \( -name '*.js' -o -name '*.mjs' \) -print` printed no paths, confirming that tool is not in `dist/` on this build.

## Build evidence

did: The existing `Makefile` wildcard includes `mk/w2guards.mk`; `make -n check` showed both real scanners and both self-tests under `check`, so no root makefile edit was needed.
ran: `make w2-guards-selftest && make w2-guards && go build ./... && go test ./internal/site/`
saw: Both self-tests passed; both real scanners passed; `go test` reported `ok github.com/dw3105/siduri-web/internal/site`.
ran: `make build && find dist -type f -exec sha256sum {} + | sort > /tmp/b1`
ran: `make build && find dist -type f -exec sha256sum {} + | sort > /tmp/b2`
ran: `diff /tmp/b1 /tmp/b2 && echo DETERMINISTIC`
saw: `DETERMINISTIC`
