# 0014 · Two-level comment threading

**Status** · accepted · 2026-08-28 · operator's word, given directly

## What changed

The operator's 28 August ruling changes the comment contract's depth ceiling
from one level to two. The old FR-13 sentence was:

> Threading is exactly one level deep: a reply to a reply is not supported and
> the UI offers no affordance for it.

In `FR-13`, the old sentence is struck and replaced with a two-level ceiling:
a reply to a reply is allowed, a reply to a depth-2 reply is refused, and the
UI offers no reply affordance on a depth-2 comment. The old one-level golden
criterion in `9 · Acceptance criteria` is struck because those golden fixtures
do not exist; thread rendering is covered by assertion tests instead. The
first item in `10 · Out of scope` is struck, leaving its separator and the
remaining exclusions intact.

The ruling also requires the form and submission path to carry the parent
identifier. `FR-1` now specifies a reply form beside every approved depth-0 or
depth-1 comment, with no form on a depth-2 comment, and a hidden `parent_id`.
`FR-2` now requires `POST /api/comments` to accept that field and refuse a
depth-2 parent. The plain HTML POST requirement in FR-1 remains in force for
each form.

The new `FR-13a · Reader replies` clause records the ruling that readers may
reply anywhere, including to another reader's reply, subject to the same
two-level ceiling. It is a suffixed clause for legibility beside `FR-13`; no
mechanical renumbering reason is being asserted.

## Why it was asked

The operator annotated the refusal on `/journal/hello-siduri/`: “allow 1 more
level, so commenter can answer to the questions addressed to him, but without
going too deep.” When asked, he ruled twice that readers may reply to any
comment and that the amendment should change the clauses that make replying
possible as well as the rendering clause.

The `author_role: site` wording in FR-13 defines an author reply; it does not
say that only a site comment may be a reply. The old role check turned that
definition into an exclusive restriction, so its removal is a bug fix in the
implementation, not a contract amendment.

## What now carries the risk

The rendering path ships in the implementation commit that follows this
amendment: an approved depth-2 reader comment committed under
`content/comments/` will appear in the static thread. The reply affordance is
specified now and arrives with the P2 Worker; before P2, nothing a reader can
use exists at P0 because the current tree has fixtures and no submission path.
This amendment is what makes the P2 reply control part of the contract; the P2
spec previously had no reply control either.

The implementation keeps a mechanical depth check, now against a depth-2
parent, and deliberately omits a reply form on depth-2 comments. The endpoint
must enforce the same boundary. Nested views are assembled without stale value
copies, and each nested reply list is sorted before it is attached, so the
static HTML and structured data remain deterministic and agree. The assertion
test checks the reader-role depth-2 ULID in the visible `id="comment-..."`
markup rather than accepting a JSON-LD-only occurrence.

The old golden-file criterion was false before this amendment and remains
unimplemented: `testdata/golden/` contains no comment fixtures, and nothing
parses this comments acceptance block. Its removal therefore breaks no
harness. The live refusal check is part three of this amendment and is
re-parameterised in the implementation, while the four existing refusal
records become stale once the current fixtures are legal. The old
`comment-refused` CSS and acceptance evidence must not be mistaken for a
regression when no refusal remains in the built fixture set.
