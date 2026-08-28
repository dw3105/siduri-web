# Acceptance W1 — the operator's pass

tree `main` @ `7667b8d5c067d047bffa402cdb58fc92550a3ccc` · surface `make dev` build in
`.dev-dist/`, served on 127.0.0.1:8080 with the pinned annotation tool injected as
one `<script>` before `</body>` (sha256 `c40e131c…`). The injection is part of the
surface judged and is recorded here for that reason.

| # | Page | Element | Note | Verdict |
|---|---|---|---|---|
| — | `/` | `main > h1` "Human in the loop." | Test — channel check, not a judgement | n/a |

## `/` — 9 annotations

| # | Element | Note | Verdict | Owner |
|---|---|---|---|---|
| 1 | `header > a > strong` "Siduri" | Name is taken; brainstorm a free one | **open** — amends `D-8`, `OQ-5`, `:423` | operator |
| 2 | `footer > a` "Impressum" | Drop, add later | **open** — contradicts `LR-1` (§ 5 DDG) | operator |
| 3 | `footer > a` "Datenschutz" | Drop, add later | **open** — contradicts `LR-2` | operator |
| 4 | `footer > a` "Stack" | What is the difference from Tools? | open | design |
| 5 | `main > h2` "Latest writing" | Sounds stupid | open | copy |
| 6 | `p.post-meta` | Date only, drop read-time and tags | open | `FR-5` |
| 7 | `header > nav > a` "About" | Remove, add later | open | design |
| 8 | `header > nav > a` "Contact" | Remove, add later | **open** — `FR-16` ships `mailto:` as P0's only contact | operator |
| 9 | `footer` | Remove the footer for now | **open** — carries the only `LR-1` path | operator |

## `/journal/hello-siduri/` — 5 annotations

| # | Element | Note | Verdict | Owner |
|---|---|---|---|---|
| 1 | `main > article > header.article-header` | Looks terrible | open | design |
| 2 | comment body `p` | `<script>alert(...)</script>` visible — WTF | **open** — escaped, not executable; a security fixture published as content | content |
| 3 | `p.comment-refused` | ULID is noise to a reader | open | `A7` |
| 4 | article footer `p` | "Newsletter capture is intentionally not part of this first build." — WTF | **open** — `FR-9` element 2 rendered as an apology | design |
| 5 | `footer` | Remove footer | see `/` #9 | operator |

## `/journal/` — 1 annotation

| # | Element | Note | Verdict | Owner |
|---|---|---|---|---|
| 1 | `main > article.post-card > p.post-meta` | Date only — drop read-time and tags | open, duplicate of `/` #6 | `FR-5` |

## Pages 4–15

| Page | Result |
|---|---|
| `/tags/method/` | `ol.post-list` numbered and **unstyled** — `.post-list` has 0 CSS rules |
| `/tags/outcome/` | empty state is a bare classless `<p>`; `<main>` is two elements, no way back |
| `/tools/` | tool card "should look like console"; filter nav flagged *review later*, no verdict |
| `/tools/status/active/language/go/` | Status has `All`, Language does not — no control clears a language filter |
| `/tools/status/abandoned/` | empty state styled (`.tool-empty`, 1 rule) but reads bare in an 864px grid |
| `/tools/gateslot/` | `.tool-header` has **0** CSS rules; structurally identical to `.article-header`, built twice |
| `/about/` | **clean** |
| `/contact/` | **unlinked**, page unchanged — operator's word |
| `/stack/` | OK. Collides with `/tools/` on one line of copy only |
| `/impressum/` | **unlinked**, left as is — operator's word |
| `/datenschutz/` | **unlinked**, left as is — operator's word |
| `/404.html` | **clean** |

## Close

15 of 15 pages walked. 22 annotations. **2 clean** — `/about/` and
`/404.html` — and 3 unlinked by decision. The close line read *3 clean*
until 2026-08-28; the table it summarises has always shown two. `AC-21`:
a count is computed from the rows, never counted by hand. Finding 0081.

**Dominant finding: one chrome, six reinventions.** header, footer, nav,
post-meta, post list and the header block are each built two or three times by
different lanes, and the second copy is usually unstyled. Root is the lane
ownership rule in `AGENTS.md` — a lane may not edit the shared document
template, so a lane needing different chrome builds its own. Every gate passes
on both copies because each is internally consistent.

**Operator decisions taken during the sitting:** drop the footer; keep
Impressum, Datenschutz and Contact built but unlinked; `siduri.ai` is not
available and the name must change; all repeated elements load from one source.

**Outstanding from this step:** the manual keyboard pass (`NFR-2:294`,
criterion 3c). Not performed.
