# 0006 · design feedback arrives as annotations, not screenshots

**Status** · accepted · 2026-08-27 · operator's word, given directly

Design feedback on a rendered page arrives as a pasted annotation block, not as
a screenshot. A screenshot carries no way to say *this element* rather than
*that region of an image*; an annotation carries the element's selector, its
computed styles, its text, an intent tag and the comment.

The tool is `agentation-vanilla` — MIT, `github.com/mearnest-dev/agentation-vanilla`,
one file, no build step. It is not the `agentation` package on npm, which is a
React 18 component under `PolyForm-Shield-1.0.0` and cannot mount on a page this
generator writes. Two products wear one name and only one of them is usable here.

The working rule is the `screen-annotations` skill, built in `hitl` and held
there rather than owned there.

## The mechanism

The operator pastes a pinned copy of the tool into the devtools console against
`make dev`, marks up the page, and pastes the exported block back.

**Nothing enters the repository.** `AGENTS-PROHIBITIONS.md` puts a third-party
script behind explicit per-instance approval, and a console paste leaves that
prohibition untouched. It does not leave the *risk* untouched, and this ADR does
not pretend otherwise — 37561 bytes of someone else's JavaScript run in a
browser that holds live sessions. The file has zero `fetch`, zero
`XMLHttpRequest`, zero `WebSocket` and zero `sendBeacon`, read by two sessions
independently, and it stores in `localStorage` under a per-pathname key. That is
a claim about a file, not a guarantee about a supply chain of eight stars.

The copy is **pinned at sha256
`c40e131c8e17a0ca19dc0dab8ce71395befee16d4da6065cfaaf34e824b12c75`**, 37561
bytes. Its published URL tracks `main` rather than a tag, so pasting "the file at
that URL" is pasting whatever shipped this morning. A pinned paste is auditable;
a bookmarklet that fetches at annotate time is not, which is why the bookmarklet
was refused.

**`agentation-mcp` is never run.** It exists on npm at `1.2.0`, binds port 4747
and accepts annotations from any open tab with an origin policy of `*`. It is one
`npx` away, so the refusal is written down rather than left implied.

## The capture half only

We take annotations and read them by hand. We do not build a reader that
resolves a pasted block to a screen automatically.

Not because we lack the enumeration a reader needs. `dist/` **is** a total map
from URL to screen, bound to what renders by construction — the paths are the
URLs, and it cannot drift from the render because it is the render. That is a
stronger binding than the ledger or route table such a reader usually stands on.

The reason is volume: 31 screens, one operator, one wave. Hand-reading costs
nothing at that size and a resolver is machinery built ahead of its need.

## The selector field is a hint, not an identity

`**Text:**` and `**Selected text:**` are read off the element directly and are
sound. `**Selector:**` is not, on our markup, and must be confirmed against the
render before it is trusted.

Two defects in the tool, both measured against all 31 rendered pages and recorded
as findings 0060 and 0061: the selector builder computes an element's position
among **same-tag** siblings and emits it as `:nth-child`, which CSS counts over
all siblings; and it never roots the chain at `body`, so a bare tag selector is
shadowed by a deeper element of the same tag. Of 849 elements, 81 resolve to a
**different element** and 169 resolve to nothing.

The 81 are the ones that matter. They do not fail — they name a different
paragraph, equally plausible, and the exported block reads correct.

The fallback has its own limit, measured rather than assumed. `**Text:**` is
truncated to the element's first 80 characters, and across our 31 pages **82
pairs of elements share those 80 characters without one containing the other** —
almost all of them a navigation link appearing in both the header and the
footer, or a nav link against the page heading of the same name. So text
identifies the element everywhere except the case where the same words are
deliberately repeated, which is exactly where a design comment on navigation
lands. Where text and selector are both ambiguous, ask which one was meant
rather than pick.

## What it does not do

An annotation is an observation, never a verdict. A render reflects the source in
hand; the running system reflects the code deployed, plus its data, plus its
environment. A comment on a render cannot speak to anything deploy, data or
environment decides, and a judgement names the build it was made against.

The exported block carries `**Page:**` and `**Date:**` and **no revision**, so
the build is not in the evidence and has to be recorded beside it.

Annotations expire after seven days and the tool rewrites its own storage when it
drops one — finding 0062. Annotate and paste in the same sitting; this is not a
store to accumulate across a wave.
