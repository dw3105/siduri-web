# Acceptance W2 — the operator's second pass

tree `main` @ `57641f5` · surface `.dev-dist` served on 127.0.0.1:8080, tunnelled
to the operator's 18080, with the pinned annotation tool injected as one
`<script>` before `</body>`.

| # | Page | Element | Note | Verdict | Owner |
|---|---|---|---|---|---|
| 1 | `/` | `header > a > strong` "Siduri" | `siduri.ai` is taken — brainstorm a new free name | **open** — amends `D-8`, `OQ-5`, `:423` | operator |
| 2 | `/journal/` | `p.post-meta` | "I asked you to remove time to read" | **failed** — his ruling of 28 Aug, not implemented on this page | integrator |
| 3 | `/tools/` | `article.tool-card` | Should look like console | open, repeat of pass 1 | design |
| 4 | `/tools/` | `nav.tool-filters` | Not sure — add the task to review later | **deferred** — he has now ruled: raise it as its own task | integrator |
| 5 | `/journal/hello-siduri/` | `footer.article-footer > p:nth-child(2)` | "Newsletter capture is intentionally not part of this first build." — WTF | **fixed** — sentence removed; `FR-9`'s newsletter element stays unmet and is now a finding rather than a line on the page | integrator |
| 6 | `/journal/hello-siduri/` | `footer` "About · Contact" | Remove footer | **cannot reproduce** — 0 footers of that shape in the build, 0 `href="/about/"` anywhere | integrator |
| 7 | `/journal/hello-siduri/` | `header.article-header` | Layout of text inside is AWFUL | open — `.article-header` has 1 CSS rule | design |
| 8 | `/tags/method/` | `ol.post-list` | numbering looks like shit | open — `.post-list` has 0 CSS rules | design |
| 9 | `/tools/gateslot/` | `header.tool-header` | Text inside this looks TERRIBLE | open — `.tool-header` has 0 CSS rules | design |
| 10 | `/tools/status/abandoned/` | `p.tool-empty` | Looks unstyled | open — 1 rule, and it only sets `grid-column` | design |
| 11 | `/journal/hello-siduri/` | `p.comment-refused` | Allow one more level so a commenter can answer a question addressed to them, without going too deep | **blocked** — contradicts `docs/comments-requirements.md:156` and `:324`; needs an amendment before it can be built | operator |
| 12 | `/tools/` | `nav.tool-filters > a` "Go" | Either "All" everywhere, or nowhere | open, ruled — Status has `All`, Language does not | design |
| 13 | `/tools/` | `article.tool-card` | Stylize it like console — monotype fonts | open, ruled — first pass said "should look like console" with no direction | design |

## Fixed during the sitting

- Read time removed from `/journal/` and the article page, per his 28 Aug ruling.
- The newsletter-apology sentence removed from the article footer.
- `static/site.css:26`'s bare `header { display: flex }` scoped to `.site-header`.
  It was matching every `<header>` on every page, which is why the article header,
  the tool header and the filter nav were all laid out as shell rows. One line,
  and the root cause of three of his annotations.
| 14 | `/tags/outcome/` | — | **clean** | held | — |
| 15 | `/about/` | — | **clean** | held | — |

## Close

Nine of thirteen templates walked: `/`, `/journal/`, `/journal/hello-siduri/`,
`/tags/method/`, `/tags/outcome/`, `/tools/`, `/tools/status/abandoned/`,
`/tools/gateslot/`, `/about/`. Two clean.

Not walked: `/contact/`, `/stack/`, `/impressum/`, `/datenschutz/`, `/404.html`,
`/tools/status/active/language/go/`. Four of those six are unlinked by his own
ruling and he declined to spend time on them in the first pass.

**Dominant finding: the three "layout is TERRIBLE" annotations were one bare
element selector.** `static/site.css:26` styled every `<header>` on the page as
the site's top bar. Not a design failure in three components — one shell rule
leaking into all of them, which is the same class as the first pass's *one
chrome, six reinventions* and was invisible to every gate for the same reason.

**Still undesigned after the fix**, with rule counts measured rather than judged:
`.article-header` 1, `.tool-header` 0, `.post-list` 0, `.tool-empty` 1 and
layout-only.
