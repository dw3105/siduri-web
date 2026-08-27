# W1 · G4 — budget and phase guards

wave: W1
tree: 7667b8d5c067d047bffa402cdb58fc92550a3ccc

## Criterion 4 — "The performance budget in NFR-1 passes in CI on a real article with code blocks and images."

verdict: deferred
ran: make check A8_EXTERNAL_LINKS=1
at: 2026-08-27T18:23:09Z
saw: exit 0; the budget reported 0.00 KB JavaScript, 2.52 KB CSS, 7.40 KB article HTML, 0 third-party requests, and a 0.21 s full build; external links ran and axe found 0 violations on 31 pages
red proof: `python3 tools/budget.py --selftest` at 2026-08-27T18:24:04Z printed JavaScript, CSS, HTML, third-party, and full-build breach cases; the remote-script case also ran `python3 tools/budget.py --dist dist --build-seconds 1` at 2026-08-27T18:22:09Z and exited 1 with `third-party request budget breached: 1 request(s)`
notes: The next NFR-1 performance-gate phase, after the P1 representative build-log fixture exists, owes the missing real-article preconditions and measurements. The current article has no code fences, no Markdown images, and no rendered image tags; `dist/` has no woff2 artefact.

| NFR-1 row | State in this checkout | Evidence and scope |
|---|---|---|
| LCP, mid-tier mobile on 4G | unmeasurable today | No LCP or real-device/browser performance probe is invoked by the repository. |
| CLS | unmeasurable today | No CLS probe is invoked. |
| INP | unmeasurable today | No INP probe is invoked. |
| JS per article page | measured | `tools/budget.py` reported 0.00 KB gzipped against 30.00 KB; this checkout has no JS asset. |
| CSS | measured | `tools/budget.py` reported 2.52 KB gzipped against 15.00 KB. |
| HTML per article, pre-images | measured, but not representative | One article page measured at 7.40 KB against 60.00 KB; the article contains neither code blocks nor images. |
| Web fonts | unmeasurable today | No font count/format/size/preload check exists, and `find dist -type f -name '*.woff2'` found 0 files. |
| Third-party requests, first load | measured | `tools/budget.py` reported 0 against the 0-request budget; its parser covers external HTML resource URLs and CSS `url()` values. |
| Lighthouse, all four categories | unmeasurable today | No Lighthouse command or result exists; axe is an accessibility scan, not Lighthouse. |
| Full build, 200 posts | measured at the wrong scale | The budget target measured a 0.21 s build, but only the current one-post tree was built; no 200-post fixture or scale test exists. |

The current article-shape probe was `date -u +%Y-%m-%dT%H:%M:%SZ; printf 'code fences in content/posts/hello-siduri.md: '; grep -c '^```' content/posts/hello-siduri.md; printf 'Markdown images in content/posts/hello-siduri.md: '; grep -c -F '![' content/posts/hello-siduri.md; printf 'rendered img tags in article: '; grep -o '<img' dist/journal/hello-siduri/index.html | wc -l; printf 'woff2 files in dist: '; find dist -type f -name '*.woff2' | wc -l; printf 'article pages budget sees: '; find dist/journal -mindepth 2 -maxdepth 2 -type f -name index.html | wc -l; printf 'HTML files in dist: '; find dist -type f -name '*.html' | wc -l` at 2026-08-27T18:24:29Z. It saw 0 code fences, 0 Markdown images, 0 rendered image tags, 0 woff2 files, 1 measured article page, and 31 HTML files.

## Criterion 10 — "No cookie banner exists, and no non-essential client-side storage is written."

verdict: open
ran: `if rg -n -i 'cookie|localStorage|sessionStorage|indexedDB|document\.cookie|setItem\(|removeItem\(' Makefile mk tools internal --glob '*_test.go'; then echo 'unexpected cookie/storage check reference found'; exit 1; else echo 'no matches in Makefile, mk/, tools/, or Go tests'; fi`
at: 2026-08-27T18:24:12Z
saw: no matches in Makefile, mk/, tools/, or Go tests
red proof: none exists; no cookie/storage check is present that can be made to exit non-zero honestly, so a regression cannot currently turn a gate red
notes: The built-output probe `if rg -n -i 'cookie.?banner|cookie.?consent|consent.?banner|localStorage|sessionStorage|indexedDB|document\.cookie|setItem\(|removeItem\(' dist; then echo 'cookie banner or client-side storage marker found'; exit 1; else echo 'no banner/storage markers in dist'; fi` at 2026-08-27T18:24:12Z saw no banner or storage markers. `rg -l -i 'cookie' dist` found only `dist/datenschutz/index.html`, where the legal text explains that non-essential cookies are absent. No `.js` or `.mjs` files exist in `dist`. Those are current-state absences, not regression coverage; the requirement is open because nothing in Makefile, mk/, tools/, or the Go tests checks it.

## Criterion 16 — "No `/services` page, price, or sales copy exists anywhere in `dist/` before P3."

verdict: held
ran: make check A8_EXTERNAL_LINKS=1
at: 2026-08-27T18:23:09Z
saw: exit 0; the pre-P3 guard measured 1 article page and the final `dist/` grep emitted no forbidden artifact
red proof: With temporary `dist/g4-red-proof.html` containing `<a href="/services/">`, the exact Makefile guard `if test -d dist && grep -rIE '/services|([€$][[:space:]]*[0-9])|<form([[:space:]>]|$)' dist; then echo 'check: pre-P3 sales or form artifact found in dist/' >&2; exit 1; fi` ran at 2026-08-27T18:21:47Z, exited 1, and printed the `/services/` line plus `check: pre-P3 sales or form artifact found in dist/`. `python3 tools/pre_p3.py dist` at 2026-08-27T18:21:50Z also exited 1 with `pre-P3 guard found a forbidden artifact: dist/g4-red-proof.html: '/services'`.
notes: The Makefile grep and `tools/pre_p3.py` use a bare `/services` string, so any `/services/` link trips them; they also detect currency-price patterns but do not semantically classify arbitrary sales prose. The temporary fixture was removed and `make build` ran at 2026-08-27T18:22:19Z; the restored `python3 tools/pre_p3.py dist` at 2026-08-27T18:22:20Z saw 1 article page and passed. The contract conflict is unresolved: FR-9 requires every article footer to route one contextual line to `/services`, while criterion 16 rejects any `/services/` occurrence before P3. The current output holds criterion 16 only because that FR-9 footer line is not present.

## Criterion 17 — "`dist/` contains no form element and no third-party script before P2."

verdict: held
ran: make check A8_EXTERNAL_LINKS=1
at: 2026-08-27T18:23:09Z
saw: exit 0; the budget reported 0 third-party first-load requests, the pre-P3 guard passed, and no remote script source or form element exists in the restored output
red proof: With temporary `dist/g4-red-proof.html` containing `<form>`, the exact Makefile guard ran at 2026-08-27T18:22:02Z, exited 1, and printed the form line plus `check: pre-P3 sales or form artifact found in dist/`; `python3 tools/pre_p3.py dist` at 2026-08-27T18:22:05Z exited 1 with `pre-P3 guard found a forbidden artifact: dist/g4-red-proof.html: '<form>'`. The same fixture contained `<script src="https://third-party.example/g4.js"></script>`; `python3 tools/budget.py --dist dist --build-seconds 1` at 2026-08-27T18:22:09Z exited 1 with `Third-party first-load requests: 1` and `third-party request budget breached: 1 request(s)`.
notes: Form elements are covered by both the Makefile grep and `tools/pre_p3.py`; external script/resource URLs are covered by the budget parser, which scans HTML resource attributes and CSS `url()` values. The restored output has no `.js` or `.mjs` files and `rg -n -i '<script[^>]+src=' dist` emitted no matches. Inline JSON-LD `<script>` elements are present, but they have no `src` and are not third-party scripts. The guards prove the current static output within those scopes; they do not test future runtime behavior outside `dist/`.
