# W1 · markup, accessibility, no-JS

wave: W1
tree: 7667b8d5c067d047bffa402cdb58fc92550a3ccc

## Criterion 3 — “Every page validates, passes axe, and is fully keyboard operable.”

verdict: open
ran: make check
at: 2026-08-27T18:28:20Z
saw: exit 0; the automated validation, axe, and Go checks passed, but no manual keyboard pass was performed
red proof: the three parts have separate probes below; the keyboard requirement has no honest automated red fixture
notes: 3a is held for the current build, 3b is held, and 3c is open with owner operator. The overall verdict is open because keyboard operability is unmeasured; no automated check substitutes for the required manual keyboard pass.

### Criterion 3a — HTML validation

verdict: held
ran: npx --yes html-validate@8 --preset recommended --rule doctype-style:off --rule void-style:off $(find dist -type f -name '*.html' -print)
at: 2026-08-27T18:25:16Z
saw: exit 0; no diagnostics for all 31 built HTML pages
red proof: npx --yes html-validate@8 --preset recommended --rule doctype-style:off --rule void-style:off .g3-scratch/malformed.html at 2026-08-27T18:26:01Z exited 1 and printed `element-required-attributes` plus `no-dup-attr` (2 errors)
notes: The checked-in `make a8-html` command also exited 0 at 2026-08-27T18:21:18Z, but its `--rule doctype-style:off` override has an effective config containing only that rule; `npx --yes html-validate@8 --print-config --rule doctype-style:off dist/index.html` showed `extends: []` and only `doctype-style: off`. The explicit recommended-preset probe above is the substantive validation; it disables only the two generated void/doctype style conventions. `npx` may need network access to fetch transient html-validate@8 on a cold npm cache; parsing the built pages itself needs no network.

### Criterion 3b — axe

verdict: held
ran: make a8-accessibility
at: 2026-08-27T18:21:25Z
saw: Node v20.20.2, axe-core 4.13.0, `dist/` served on port 37141, `wait_serving: dist/ confirmed on 37141`, all 31 pages tested, 0 violations, exit 0
red proof: a temporary isolated page scanned with `npx --yes @axe-core/cli@4 --exit http://127.0.0.1:53407/axe-bad.html` at 2026-08-27T18:22:10Z exited 1 with 6 violations, including `button-name`, `html-has-lang`, `image-alt`, and `landmark-one-main`
notes: Node 20 mode did not skip the gate. Target verification was substantive: `wait_serving.py` fetched `/index.html` from the chosen port, required HTTP 200, and required the response to contain `Siduri` before axe received any URLs. The temporary axe fixture and scratch directory were removed. `npx` may need network access to fetch the transient CLI on a cold npm cache.

### Criterion 3c — keyboard operability

verdict: open
ran: make check
at: 2026-08-27T18:28:20Z
saw: exit 0; no keyboard walkthrough was performed
red proof: cannot be made honestly by automation; the contract assigns this to a manual operator release pass
notes: Owner: operator. Axe and HTML validation do not establish tab order, focus visibility, keyboard activation, or escape behavior.

## Criterion 5 — “The site is fully readable with JavaScript disabled, including navigation and article content.”

verdict: held
ran: date -u +%Y-%m-%dT%H:%M:%SZ && find dist -type f \( -name '*.js' -o -name '*.mjs' \) -print && find dist -type f -perm /111 -print && rg -o '<script[^>]*>' dist | sort | uniq -c && python3 tools/budget.py --dist dist --build-seconds 1
at: 2026-08-27T18:26:39Z
saw: no `.js` or `.mjs` files and no executable files under `dist/`; 32 script tags were present across the pages and every one was `type="application/ld+json"`; the budget reported JavaScript 0.00 KB gzipped
red proof: no built-tree no-JavaScript invariant exists to turn red for a small added script; `go test ./internal/site -run 'TestA1BodyRendererIsBuildTimeAndKeepsPlainMarkdown|TestA6NoJavaScriptIndexListsEveryTool' -count=1` at 2026-08-27T18:22:41Z passed, but those are component checks only
notes: The JSON-LD data blocks are non-executable. The two component tests reject a script in the article body and tools index. The budget code counts `*.js` bytes but would allow a small file under its 30 KB budget; no check scans all rendered HTML for executable script tags or asserts that the built tree contains no JavaScript. The current static navigation and article content are therefore readable without JavaScript, with that future-regression gap recorded.

## Criterion 9 — “Impressum and Datenschutz are reachable in two clicks from every page.”

verdict: held
ran: go test ./internal/site -run '^TestA5RequiredPagesAndLegalLinkGraph$' -count=1
at: 2026-08-27T18:23:37Z
saw: exit 0; the A5 graph test checked both legal targets from every built page within two clicks
red proof: the graph walk reported `clean_depth=2`; after removing every `/impressum/` edge in an in-memory copy it reported `red_proof_after_removing_all_impressum_edges=None` at 2026-08-27T18:24:16Z
notes: A direct-link walk over `dist/` at 2026-08-27T18:23:14Z found `pages=31`, `direct_impressum=7`, and `direct_datenschutz=7`. From the article, the measured paths are `journal/hello-siduri/index.html -> about/index.html -> impressum/index.html` and `journal/hello-siduri/index.html -> about/index.html -> datenschutz/index.html`: article → About → each legal page. The seven direct links are expected because the footer carries About/Contact only; the graph evidence does not establish visual prominence or keyboard operability.
