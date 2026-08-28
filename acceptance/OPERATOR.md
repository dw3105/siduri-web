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
| 6 | `p.post-meta` | Date only, drop read-time and tags | held (via `lane/w3`) | `FR-5` |
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
| 1 | `main > article.post-card > p.post-meta` | Date only — drop read-time and tags | held (via `lane/w3`), duplicate of `/` #6 | `FR-5` |

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

## Probes

Added 2026-08-28, after the sitting, on tree `14f5296` and host
`legalcopilot-dev.europe-west3-b.c.legalcopilot-dev.internal`. The sitting itself
recorded no evidence; this is a repair, not a claim that evidence existed at the
time. Written in the convention `tools/acceptance.py:35` parses — `field:` at
line start — because the detail-block form used elsewhere in `acceptance/` is
read by nothing. Finding 0089.

Several rows are judgements of taste with no command behind them ("Sounds
stupid", "Looks terrible"). Those carry the build that produced what he read and
the rendered state he read, which is what `AC-18` asks for. A probe cannot make
an opinion measurable, and pretending otherwise is the failure `AC-36` names.

### `/` rows 1, 2, 3, 9 — the name, the legal links and the footer

verdict: open
ran: grep -rIl 'href="/impressum/"' dist/ | wc -l
at: 2026-08-28T11:58:00Z
saw: 0 of 31 pages link it; the first pass measured 7 of 31, before the footer was dropped
notes: Owner: operator. His ruling of 27 August dropped the footer and kept the legal pages built and unlinked, which is decision 18. `LR-1` (§ 5 DDG) requires reachability and nothing now reaches them, so the ruling and the clause are in open conflict. Rows 2, 3 and 9 are one defect with three annotations. Row 1, the site name, is probed in the second pass.

### `/` rows 4, 5, 7 — Stack versus Tools, the heading, the nav

verdict: open
ran: grep -oE '<nav[^>]*>.*</nav>' dist/index.html | grep -oc '<a '
at: 2026-08-28T11:58:00Z
saw: 2 nav items — Journal and Tools; About and Contact are gone, per his ruling
notes: Owner: design and copy. Row 4's question — what distinguishes Stack from Tools — is unanswered in the tree; both pages exist and neither says. Row 5 is taste. Row 7 is done.

### `/` row 6 and `/journal/` row 1 — post-meta

verdict: held (via `lane/w3`)
ran: grep -rIoE '[0-9]+ min read|min read' dist/ | wc -l
at: 2026-08-28T11:42:28Z
saw: 0
red proof: the same command on `7667b8d` returned a non-zero count
notes: He revised this ruling on 28 August to date and tags rather than date alone — decision 21. Both annotations are one defect, met.

### `/` row 8 — the Contact link

verdict: open
ran: grep -rIl 'href="/contact/"' dist/ | wc -l
at: 2026-08-28T11:58:00Z
saw: 0 — built and unlinked
notes: Owner: operator. `FR-16` ships `mailto:` as P0's only contact path and the page carrying it is now unreachable. Same shape as rows 2, 3 and 9: his ruling and a clause disagree, and the clause is his to amend.

### `/journal/hello-siduri/` rows 1–5 — the article page

verdict: open
ran: for c in article-header post-list; do grep -cE "^[^{]*\.$c[^{]*\{" static/site.css; done
at: 2026-08-28T11:42:28Z
saw: `.article-header` 2 rules, `.post-list` 0
notes: Row 1 (header layout) is open with owner design and is re-probed in the second pass. Row 2, the escaped `<script>` fixture published as content, and row 3, the ULID in the refusal line, are finding 0071 — both are test fixtures shipped as the only post's content. Row 4, the newsletter apology, is met. Row 5 is the footer, above.

### Pages 4–15 — the sweep

verdict: open
ran: make check
at: 2026-08-28T11:42:28Z
saw: exit 0; 31 pages, 0 axe violations; `/about/` and `/404.html` drew no annotation
notes: Twelve pages walked in one pass. Two clean, three unlinked by his ruling, and the rest carry the unstyled-component defects that finding 0067 groups. `/impressum/` renders 1 `NOCH EINZUTRAGEN` placeholder today; the plan recorded five before the pass, so the number went stale and the mechanism did not — `AC-38`.

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
