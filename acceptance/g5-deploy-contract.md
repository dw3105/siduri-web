# W1 · G5 — deploy path and agent contract

wave: W1
tree: 7667b8d5c067d047bffa402cdb58fc92550a3ccc

## Criterion 12 — Contact form works with and without JavaScript, and mail lands in Gmail within 60s.

verdict: deferred
owner: phase P2, gated on the postal address at docs/site-requirements.md section 12
ran: date -u +%Y-%m-%dT%H:%M:%SZ; git grep -nE '<form|bytes.Contains\(.*<form|pre_p3.py dist' -- Makefile internal/site/pages_a5_test.go mk/A9.mk tools/pre_p3.py
at: 2026-08-27T18:24:20Z
saw: form rejection is in Makefile:34, internal/site/pages_a5_test.go:41, and tools/pre_p3.py:21; mk/A9.mk:10 runs the latter against dist.
red proof: set +e; date -u +%Y-%m-%dT%H:%M:%SZ; if grep -rIE '/services|([€$][[:space:]]*[0-9])|<form([[:space:]>]|$)' <(printf '<form>\n'); then echo 'check: pre-P3 sales or form artifact found in dist/'; status=1; else status=0; fi; echo exit $status; exit 0 — saw the fixture and `check: pre-P3 sales or form artifact found in dist/`, exit 1.
notes: P0/P1 intentionally expose only a mailto contact link. FR-16 places the form in P2, and the P2 gate is a postal address: “No list, no form and no comment intake before it does.” The P2 audience phase owns this; JavaScript/no-JavaScript delivery and Gmail timing cannot be tested before that gate without changing the guards.

## Criterion 13 — Pull requests produce a working preview URL with `noindex`.

verdict: deferred
owner: the P0 deployment lane, which owes the first real upload, the URL surfacing and the noindex verification
ran: date -u +%Y-%m-%dT%H:%M:%SZ; git merge-base --is-ancestor lane/preview HEAD; status=$?; if test $status -eq 0; then echo 'lane/preview is merged into HEAD'; else echo 'lane/preview is not merged into HEAD'; fi; git show lane/preview:mk/a8.mk | sed -n '107,113p'; git show lane/preview:.github/workflows/ci.yml | sed -n '54,63p'; echo exit $status
at: 2026-08-27T18:24:01Z
saw: lane/preview is not merged into HEAD; its preview command is `$(WRANGLER) versions upload --env preview --message "Siduri pull request preview"`, and CI passes the two Cloudflare secret names to `make preview`.
red proof: set +e; date -u +%Y-%m-%dT%H:%M:%SZ; if git grep -nE 'GITHUB_OUTPUT|echo.*(https?://|[Uu][Rr][Ll])' -- .github/workflows/ci.yml mk/a8.mk; then status=0; else echo 'no preview URL output step is present'; status=1; fi; echo exit $status; exit 0 — saw `no preview URL output step is present`, exit 1.
notes: The upload fix exists on lane/preview, is unmerged, and has never been executed against Cloudflare. The current workflow names CLOUDFLARE_API_TOKEN and CLOUDFLARE_ACCOUNT_ID but provides no repository evidence that they are configured. The preview lane/P0 deployment step owes the first real upload, URL surfacing, and noindex verification.

## Criterion 14 — `make rollback` restores the previous deployment and has been tested at least once.

verdict: deferred
owner: the P0 deployment lane, which owes the first deployment and an honest rollback test
ran: set +e; date -u +%Y-%m-%dT%H:%M:%SZ; make rollback A8_DRY_RUN=1; status=$?; echo exit $status; exit 0
at: 2026-08-27T18:24:01Z
saw: `dry run: wrangler rollback --name siduri-web --message "Siduri rollback"`, exit 0.
red proof: set +e; date -u +%Y-%m-%dT%H:%M:%SZ; make a8-rollback; status=$?; echo exit $status; exit 0 — saw `rollback: wrangler is required (not installed in this lane)` and exit 2.
notes: The real path shells out to Wrangler, while the dry-run path only prints the command. There is no deployment or rollback-test artifact in this repository, and invoking a real rollback without a known deployment would be an external mutation. The P0 deployment lane owes the first deployment and an honest rollback test.

## Criterion 15 — `AGENTS.md` lets a fresh agent session publish a post without asking a question.

verdict: deferred
owner: the P0 deployment lane, which owes make publish and a wired paths_exist check
ran: set +e; date -u +%Y-%m-%dT%H:%M:%SZ; make publish; status=$?; echo exit $status; exit 0
at: 2026-08-27T18:24:01Z
saw: `publish: Wave P1 fills this target`; `make: *** [Makefile:37: publish] Error 1`; exit 2.
red proof: the same `make publish` input is the red proof: the target is still a P1 stub and exits 2.
notes: `python3 tools/paths_exist.py AGENTS-FILLABLE.md AGENTS-PROHIBITIONS.md` passes with “2 file(s), every named path resolves,” but `git grep -n 'paths_exist' -- Makefile mk .github tools` finds no Make or CI invocation, only the tool itself. A real test would start an isolated fresh session with only AGENTS.md, have it create a valid post using the documented paths/frontmatter, run the implemented P1 publish command without interaction, and assert the published output; that test cannot run against this stub. AR-8 still makes actual publication human-approved.

## Criterion 18 — Comments are automatically closed on posts older than 12 months (FR-19).

verdict: failed
ran: set +e; date -u +%Y-%m-%dT%H:%M:%SZ; go test ./internal/site -run 'TestA7CommentFreezeUsesFixedReferenceDate|TestA7FrozenThreadShowsOneLineAndNeverShowsAForm' -count=1; status=$?; echo exit $status; exit 0
at: 2026-08-27T18:24:01Z
saw: `ok github.com/dw3105/siduri-web/internal/site`; exit 0.
red proof: set +e; date -u +%Y-%m-%dT%H:%M:%SZ; python3 -c "from datetime import date; from pathlib import Path; p=Path('internal/site/comments_a7.go').read_text(); expected='const commentFreezeReferenceDate = \"' + date.today().isoformat() + '\"'; assert expected in p, expected"; status=$?; echo exit $status; exit 0 — saw `AssertionError: const commentFreezeReferenceDate = "2026-08-27"`, exit 1.
notes: internal/site/comments_a7.go quotes `const commentFreezeReferenceDate = "2026-08-25"`; it is two calendar days behind UTC today, 2026-08-27. The test is green because it asserts behavior against that same fixed literal. This is an active failure of “automatically”; the P2 comments implementation must use the current date while preserving deterministic build behavior.
