# W2 · A9 — comment threading goes to two levels

base: 742e376
lane: A9 / lane/w2-threading

## Amendment — “allow 1 more level, so commenter can answer to the questions addressed to him, but without going too deep”

did: Struck and restated the `FR-1` and `FR-2` form/endpoint clauses with a per-comment `parent_id` affordance and a depth-2 endpoint boundary; struck the old depth sentence in `FR-13`, the false golden-file criterion in `9 · Acceptance criteria`, and the first item in `10 · Out of scope`; added `FR-13a · Reader replies`; added `docs/adr/0014-two-level-comment-threading.md`. The first commit is `a93cb7e` (`comments: allow two levels of threading`).
ran: `python3 tools/amendcheck.py HEAD^ HEAD`
saw:

```text
rule 1 (added ADR): PASS — added ADR with a fresh number and dated status: docs/adr/0014-two-level-comment-threading.md (**Status** · accepted · 2026-08-28 · operator's word, given directly)
rule 2 (struck traces): PASS
  docs/site-requirements.md: PASS — no diff
  docs/comments-requirements.md: hunk 1: PASS — 2 removed line(s), 3 strike span(s), 2 injective match(es)
  docs/comments-requirements.md: hunk 2: PASS — 1 removed line(s), 1 strike span(s), 1 injective match(es)
  docs/comments-requirements.md: hunk 3: PASS — 2 removed line(s), 2 strike span(s), 2 injective match(es)
rule 3 (named sections): PASS
  docs/site-requirements.md: PASS — no diff
  docs/comments-requirements.md: PASS — added ADR names touched section(s): 10 · Out of scope, 9 · Acceptance criteria, FR-1, FR-13, FR-13a, FR-2
rule 4 (watched paths): PASS — all watched contract paths exist at base and HEAD
decision: legitimate amendment; all four rules pass
```

red proof: `python3 tools/amendcheck.py --selftest` planted `delete-outright` and printed `rule 2 (struck traces): FAIL`; it also planted `strike-without-adr` and printed `rule 1 (added ADR): FAIL`, and `strike-wrong-section-name` and printed `rule 3 (named sections): FAIL`. The selftest removed its temporary fixtures.
notes: The role wording in `FR-13` is read as a definition of an author reply, not an exclusive restriction. `FR-13a` records the operator's wider ruling that readers may reply to any comment, including another reader's reply. The affordance is specified now and arrives with the P2 Worker; no usable submission path exists at P0. The comments acceptance block remains unparsed by the acceptance harness, so the struck golden-file criterion was false before and after this amendment.

## Implementation — two-level rendering

did: Re-parameterised `arrangeCommentRecords`' depth guard to reject only a reply whose parent has a parent, removed the incorrect site-only role gate, and replaced value-copy attachment with recursive bottom-up assembly that sorts every reply list before attachment (`internal/site/comments_a7.go:202-246`). Added the reader-role depth-2 fixture `content/comments/hello-siduri/01K3QZJ8X4YB7N2M9V0PQRSTUY.md`; changed the now-valid site depth-2 fixture's body; reduced nested `.comment-replies` indentation (`static/comments_a7.css:52-55`). The existing `@commentCard(reply)` call in `internal/site/comments_a7.templ:46-50` was already recursive, so neither it nor `comments_a7_templ.go` needed a change.
ran: `go test ./internal/site/ -run '^TestA7RepositoryCommentsRenderAsOneLevelThread$' -count=1`
saw: `ok github.com/dw3105/siduri-web/internal/site 0.065s`
red proof: Before adding the reader fixture and implementation, the new assertion failed with `comments_a7_test.go:51: reader depth-2 comment is missing from rendered thread markup`; the test name `TestA7RepositoryCommentsRenderAsOneLevelThread` was retained. The assertion targets `id="comment-01K3QZJ8X4YB7N2M9V0PQRSTUY"`, not an arbitrary ULID in JSON-LD. The count assertions now derive the number of visible comment elements instead of pinning `2`.

did: Added a temporary tier-3 reader fixture parented to `01K3QZJ8X4YB7N2M9V0PQRSTUY`, built it, observed the grandparent condition fire, and removed the fixture before committing. The condition is `parent.ParentID != ""` followed by `grandparent.ParentID != ""`; the role check is gone.
ran: `make build && grep -n -C1 '01K3QZJ8X4YB7N2M9V0PQRSTUZ\|replies can be attached only up to two levels deep' dist/journal/hello-siduri/index.html`
saw: `<p class="comment-refused" data-comment-refusal="Comment 01K3QZJ8X4YB7N2M9V0PQRSTUZ was not rendered: replies can be attached only up to two levels deep.">This comment was not shown because replies can be attached only up to two levels deep.</p>`; `git status --porcelain` after removal listed no `...PQRSTUZ` path.
notes: The permanent fixture set has four rendered comments and no `.comment-refused`; the disappearing refusal is intentional. Four stale records/evidence surfaces remain outside this lane's scope: `static/comments_a7.css:57-63` is dead refusal CSS until a refusal exists; `acceptance/OPERATOR.md:32` has the open “ULID is noise to a reader” annotation, now moot against the absent element; `acceptance/W2/a3-comments.md:17,26` greps and quotes the old refusal; and `acceptance/W2/a5-contrast.md:14` targets `.comment-refused` with axe. `testdata/golden/` needed nothing: its three fixtures contain no comments, and thread rendering is assertion-tested in `comments_a7_test.go`, not golden-file tested.

## Final evidence

ran: `go test ./internal/site/`
saw: `ok github.com/dw3105/siduri-web/internal/site 1.134s`

ran: `go build ./...`
saw: exit 0

ran: `make build && rg -o 'id="comment-01K3QZJ8X4YB7N2M9V0PQRSTUY".{0,120}' dist/journal/hello-siduri/index.html`
saw: `<div class="comment" id="comment-01K3QZJ8X4YB7N2M9V0PQRSTUY" data-comment-id="01K3QZJ8X4YB7N2M9V0PQRSTUY" itemscope itemtype="https://schema.org/Comment">`; the corresponding `...PQRSTUX` opening element was also present, and `comment-refused` was absent from the final fixture build.

ran: `find dist -type f -exec sha256sum {} + | sort > /tmp/b1 && find dist -mindepth 1 -delete && make build && find dist -type f -exec sha256sum {} + | sort > /tmp/b2 && diff /tmp/b1 /tmp/b2 && echo DETERMINISTIC`
saw: `DETERMINISTIC`

notes: `make check` was not run because COMMON.md directs parallel lanes not to run the full gate; the integrator must run it once on the merged tree. The depth-3 fixture was removed before the implementation commit, and the required final state is clean `git status --porcelain`.
