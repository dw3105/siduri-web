# Website — Requirements

**Project** · **Siduri** — vibecoding journal & services site
**Domain** · `siduri.ai`
**Principle** · Human in the loop
**Purpose** · Grassroots marketing of vibecoding services through published work
**Stack** · Go + templ + htmx · static · Cloudflare Workers · GitHub Actions
**Operator** · Solo; jurisdiction under review (LR-7)
**Version** · 0.1 — draft for review
**Date** · 2026-08-25
**Component specs** · `comments-requirements.md`

---

## 1 · Business context

### The single job

**Eventually:** convert a stranger's technical curiosity into a qualified inquiry, by making it obvious that the operator ships working software fast.

**Right now there is no offer, so the site cannot sell.** Its job at launch is narrower and must not be confused with the one above: earn enough credibility with technical readers that when an offer exists, there is already an audience to put it in front of.

This is a deliberate sequence, not a delay:

1. **Be worth reading.** The introduction post, then Gateslot and its build log.
2. **Be worth subscribing to.** A body of work and a reason to return.
3. **Sell.** `/services` is written when the offer is one sentence — not before.

The phase gate is a requirement, not a preference. A services page with nothing behind it is worse than no services page: it makes an unbacked promise, and it invites inquiries the operator cannot yet answer. Ship no sales furniture until phase 3.

Every page, every post, every feature answers to the job of the *current* phase. If a proposed feature does not make someone more likely to believe "this person ships" or more likely to come back, it does not get built yet.

### Audience

| Segment | Who | What they need to see | Buys? |
|---|---|---|---|
| Primary | Founder or ops lead at a 5–50 person company with an internal tool they can't get built | Evidence of speed, evidence of finishing, a price | Yes |
| Secondary | Technical peer, solo dev, indie hacker | Method, tooling, honest postmortems | No — amplifies |
| Tertiary | Recruiter, agency, potential collaborator | Range and range of stack | Sometimes |

The secondary segment drives traffic; the primary segment pays. The site is written for the secondary and structured for the primary.

### Success metrics

The north star changes by phase. Using the phase-3 metric during phase 1 produces despair and bad decisions, because the number is structurally zero.

| Phase | North star | Supporting | Target |
|---|---|---|---|
| 1 · Be worth reading | Posts shipped | Gateslot published and usable | 2 posts, 1 tool |
| 2 · Be worth subscribing to | Newsletter subscribers | Returning readers, comments received | 250 in 6 months |
| 3 · Sell | Qualified inquiries / month | Article → `/services` ≥ 4%, `/services` → contact ≥ 8% | 4 / month |

Explicitly **not** metrics, in any phase: pageviews, time on page, social follower counts, HN front pages. They are noise that feels like signal.

---

## 2 · Positioning

### The thesis

**Human in the loop.** Machines do the work; a person stays at the gate. That is the guiding principle of the brand, the argument the business sells, and the constraint every system on this site is built under.

It is also the name. Siduri keeps a tavern at the edge of the world. Gilgamesh arrives half-mad with grief, hunting immortality; she tells him plainly that he will never find it — go home, eat your bread, wash, hold your child — and then tells him where the ferryman is anyway. The counsel that is honest about the quest and still helps with the next step. That is the posture: not the agent that does everything, and not the sceptic who does nothing.

The site is the artifact. It is built with the method it sells — vibe-coded, agent-operated, statically compiled, free to run — and it says so. A marketing site for vibecoding services that was not itself vibe-coded is a contradiction the audience will notice. One that auto-publishes unreviewed agent output is refuted by its own homepage.

Corollary: the build process, the tooling, the mistakes and the cost model are all content. The repository is a portfolio piece.

On the `.ai` domain: there is an obvious tension in a human-in-the-loop brand living at an AI address. Name it rather than avoid it — the tension *is* the pitch, and it belongs in the tagline.

### Content pillars

| Pillar | What it is | Share | Primary reader |
|---|---|---|---|
| Build logs | A thing built, start to shipped, with the actual prompts and the actual dead ends | 40% | Peer |
| Tool releases | Something usable, with a repo | 20% | Peer |
| Dogfooding & postmortems | What broke, what it cost, what it changed | 20% | Both |
| Method | How the agentic workflow is structured, and why | 10% | Both |
| Buyer-legible pieces | The same work told as an outcome, in the language of someone who isn't a developer | 10% | Buyer |

The last row is small and it is the one that will get skipped. It is also the one the business will eventually run on. **From phase 2 onward, at least one buyer-legible post per month is a hard requirement, not an aspiration.** In phase 1 the quota is suspended — with two posts total, the mix is meaningless.

### Voice

Direct, technical, specific, first person. Real numbers, real durations, real costs. Name the tools. Show the failures — a postmortem is more persuasive than a success story, because everyone has seen fabricated success stories.

Never: growth-hack register, "in today's fast-moving landscape", em-dash-laden LinkedIn cadence, or any sentence that would survive being written by someone who had not done the work.

### Failure modes to design against

Each of these has killed a site like this before.

- **It becomes a blog and stops selling.** Mitigation: the buyer-legible quota above; `/services` in the primary nav; a quarterly review of inquiry count, not post count.
- **Posts optimize for peer approval over buyer legibility.** Mitigation: every post carries a one-line plain-language summary that a non-developer could parse. If that line can't be written, the post is unfinished.
- **The tooling becomes the project.** Yak-shaving the generator instead of writing. Mitigation: a hard time budget on infrastructure work, and the v0 ship gate in §12.
- **Comments become a free support desk.** Mitigation: replies demonstrate competence and stop; anything needing real work routes to a paid conversation.
- **Publishing stops for six weeks and the site reads as abandoned.** Mitigation: visible last-updated date; a buffer of two drafted posts at all times.

---

## 3 · Decisions

#### D-1 · Static-first, no runtime content rendering
All content HTML is generated at build time. The Worker exists only for form intake. If the Worker is down, the site is fully readable.

#### D-2 · Own the stack, no third-party widgets
No embedded analytics scripts, comment services, chat bubbles, font CDNs, cookie-consent SaaS, or social embeds. Every one of them is a tracking vector, a performance cost, and — on a site selling technical judgement — an admission.

#### D-3 · No tracking, therefore no cookie banner
Cookieless, server-side-only measurement. This is a design constraint chosen deliberately: a consent banner on a site whose pitch is technical taste is an own goal, and avoiding it is easier than complying with it.

#### D-4 · Content and code in one repository
Posts, tool pages, comments, templates and the generator live together. One `git push` is the only publishing action. No CMS, no external content store.

#### D-5 · The repository is public
The repo is the strongest single proof of the pitch. Consequences are binding, not optional: no secrets in the tree, no commenter email addresses (see `comments-requirements.md` D-7), no client names without written permission, and commit messages readable by strangers.

#### D-6 · htmx and progressive enhancement, no SPA framework
Every interaction works without JavaScript. htmx improves it. Nothing requires it.

#### D-7 · Agent-operable by construction
Any routine operation — publish, moderate, add a tool, deploy, roll back — is a single command executable by an agent from a clean checkout. This is a constraint on the design, not a documentation task bolted on afterwards.

#### D-8 · Brand: Siduri, at `siduri.ai`
The site trades under a brand, not the operator's name. **Siduri**, selected after diligence recorded in OQ-5; domain `siduri.ai`.

Rejected as a *name*: the human-in-the-loop family, which is the industry's generic label for the concept, already traded under by at least two companies, subject to two US applications on the exact phrase, and impossible to rank against wall-to-wall definitional content. It survives as the principle (D-11), which is where it does real work.

Also rejected on diligence: Phrontis, Huginn, Mímir, Forseti, Etak, and the Polytropos family. See OQ-5 for the evidence and the method.

Consequences of `.ai`: the Anguilla registry mandates a **two-year minimum registration** — no one-year term exists at any registrar — at roughly €50–80/year, billed as a block of about €100–160 up front. The name is therefore committed for at least two years, which is what the diligence was for. Consider a defensive `siduri.dev` as canonical fallback and to deny a confusable neighbour; note also that ccTLD policy is set by a single small territory and has been disrupted elsewhere.

#### D-9 · No offer at launch
The site launches with no services page, no price, and no sales copy. This is the sequencing decision from §1, recorded here because it constrains the IA, the phasing and the legal surface.

#### D-11 · Human in the loop is an operating constraint, not a slogan
Every system here that could act autonomously has a human gate, and the gate is load-bearing rather than ceremonial:

- No comment publishes without per-comment human approval (comment spec D-4, §6.2).
- No post publishes without human approval; agents draft, never publish (AR-8).
- No agent modifies legal pages, adds a dependency or a paid service, or touches secrets without explicit per-instance approval (AR-8).
- The gate is visible to readers rather than hidden: the pending state exists partly so a commenter knows a person will read it.

The constraint runs both ways. If a proposed feature cannot carry a meaningful human gate, that is evidence it should not be automated here at all. And an agent asking for approval more often than a human can absorb is a design failure, not diligence — the gate has to be where judgement is actually needed.

This is what the business sells, so the site has to be a working demonstration of it. The comment pipeline is the flagship: it is the principle, running in production, in public.

#### D-12 · English content, German legal pages
The market for these services is international; the content is English. The legal pages are German, because that is where the operator's obligations sit — see §9 and OQ-3. No i18n framework, no translation pipeline.

---

## 4 · Information architecture

| Path | Purpose | Source | Priority |
|---|---|---|---|
| `/` | Thesis, recent posts, newsletter capture | `content/home.md` + index | P0 |
| `/journal/` | Reverse-chronological index, paginated | `content/posts/**` | P0 |
| `/journal/<slug>/` | Article | `content/posts/<slug>.md` | P0 |
| `/contact/` | A `mailto:` address only — no form until phase 2 | template | P0 |
| `/services/` | The offer, the proof, the price, one action | `content/services.md` | **P3 — blocked** |
| `/about/` | Who, why, credibility — and the Siduri story as the statement of principle | `content/about.md` | P0 |
| `/impressum/` | Legal (DE) | `content/legal/impressum.md` | P0 |
| `/datenschutz/` | Legal (DE) | `content/legal/datenschutz.md` | P0 |
| `/404.html` | Not found, with a route back | template | P0 |
| `/sitemap.xml`, `/robots.txt` | Crawlers | generated | P0 |
| `/tools/` · `/tools/<slug>/` | Things built, with repos | `content/tools/**` | **P1** |
| `/tags/<tag>/` | Topic index | derived | P2 |
| `/feed.xml` · `/feed.json` | Syndication | generated | P0 |
| `/search/` | Full-text | Pagefind index | P2 |
| `/stack/` | Tools actually used, updated | `content/stack.md` | P2 |
| `/llms.txt` | Machine-readable site summary | generated | P2 |
| `/changelog/` | Site's own build log | generated from git | P3 |

Nav at launch is four items: Journal · Tools · About · Contact. Services is added as a fifth, visually distinct item in phase 3 and not before. Everything else lives in the footer.

---

## 5 · Functional requirements

### 5.1 Content pipeline

#### FR-1 · Source format
Posts are Markdown with YAML frontmatter in `content/posts/`. Required fields: `title`, `slug`, `date`, `summary`, `plain_summary`, `tags`, `draft`. The `plain_summary` field is the buyer-legible one-liner from §2 and the build fails without it.

#### FR-2 · Build command
`make build` produces `dist/` from a clean checkout with no network access and no manual steps. Output is byte-identical across runs given identical input.

#### FR-3 · Drafts
`draft: true` excludes a post from `dist/` entirely — not `noindex`, not unlinked, absent. Preview deploys build drafts under a `noindex` header so they can be reviewed at a real URL.

#### FR-4 · Assets
Images are optimized at build time to WebP/AVIF with explicit width and height attributes, content-hashed filenames, and `loading="lazy"` below the fold. No unprocessed image reaches `dist/`.

### 5.2 Article pages

#### FR-5 · Reading experience
Single column, ~68ch measure, no sidebar. Table of contents for posts over 1,500 words, inline on mobile and in the margin on wide screens. Reading time and publish date visible; updated date shown when it differs.

#### FR-6 · Code blocks
Syntax highlighting applied at build time — no client-side highlighter. Copy button, language label, optional filename caption, and horizontal scroll rather than wrapping. Given the subject matter, code blocks are primary content and get first-class typographic treatment.

#### FR-7 · Series
Posts may declare `series` and `part`. Series members render a compact navigator: position, previous, next. Build logs are naturally serial and this is how a casual reader becomes a repeat reader.

#### FR-8 · Cross-references
Posts reference tools by slug and render as linked cards. Tool pages list the posts that mention them. This bidirectional link is built from the content graph, not maintained by hand.

#### FR-9 · Post footer
Every article ends with the same three elements, in order: author line, newsletter capture, one contextual line routing to `/services`. No sales copy inside the post body.

### 5.3 Discovery

#### FR-10 · Tags
A closed, hand-maintained tag vocabulary — the build fails on an unknown tag. Free tagging on a solo site produces forty tags with one post each.

#### FR-11 · Search
Pagefind, indexed at build time, running client-side. Loaded on interaction only, never on page load. Degrades to a link to `/journal/` without JavaScript.

#### FR-12 · Related posts
Three related posts at the end of each article, computed at build time from shared tags and series membership. Deterministic, no ML, no scoring service.

### 5.4 Tools

#### FR-13 · Tool pages
One page per tool: what it does, why it exists, a runnable example or screenshot, install instructions, repository link, current status (`active` / `maintained` / `abandoned`). **`abandoned` is displayed honestly** — a graveyard of finished experiments is more credible than a wall of eternal betas.

#### FR-14 · Tool index
Grid sorted by recency, filterable by status and language, rendered statically with htmx-swapped filter fragments.

### 5.5 Conversion path

#### FR-15 · Services page — phase 3, blocked
Deferred by D-9 and blocked on OQ-1. When written it must state, above the fold: what problem gets solved, who for, and what engagement looks like. Then three links to build logs as proof — not testimonial cards, not logos, actual work. Then a price or a price band. Then one action, repeated at the bottom.

Silence on price is the single biggest conversion killer for solo consultants. A range with conditions beats "let's talk". Record the decision in OQ-2 before this page is written, not while writing it.

#### FR-16 · Contact — `mailto:` at launch
Phase 0 and 1 ship a plain `mailto:` link and nothing else. No form, no Worker, no data processing, no spam surface, and — importantly — no personal data collected before the §9 questions are settled.

A form arrives in phase 2 alongside comments, reusing the same Worker: `POST /api/contact`, Turnstile, delivered to Gmail via Resend, same layered spam defence. Fields: name, email, one message box, one optional line — "what are you trying to build". No dropdowns, no budget selector, no phone field.

#### FR-17 · Inquiry acknowledgement
From phase 2: automatic reply within seconds stating a real response time. If the promise is 48 hours, the calendar has to honour it.

#### FR-18 · Newsletter
**Buttondown**, confirmed. Owned distribution, double opt-in, tracking disabled, plain text or lightly styled HTML, CSV export verified before the first send.

**Hard prerequisite: a postal address.** Commercial email requires a valid physical postal address in the footer in essentially every jurisdiction that regulates it — US CAN-SPAM explicitly, and German practice via the Impressum. This is a property of running a mailing list, not of being in the EU, and it does not go away by relocating. It does not have to be a home address: a registered PO box or a private mailbox from a commercial mail-receiving agency satisfies it, for roughly €100–200/year.

Sequencing consequence: **no list before a mailbox.** Until then, RSS and JSON Feed are the distribution channels — they ask nothing of the operator and cost nothing to run.

Rejected: self-hosted Listmonk, which is free software but needs a server and permanent maintenance, breaking NFR-7 and adding ops for no reader-visible gain. Also rejected for now: building it on Worker + Resend — on-brand and a good post, but signup, confirmation, storage, unsubscribe and bounce handling are weeks of scope aimed at a problem that isn't the point. Revisit as a phase-3 build log, when there is a list large enough that migrating it is itself the interesting story.

Signup lives in the post footer and on `/`, nowhere else. No modals, no exit-intent, ever. Blocked on §9: no signup ships before the controller address question is settled.

### 5.6 Comments

#### FR-19 · Comments
Specified in full in `comments-requirements.md`. Summary: reader submits → held in KV → notification to Gmail → agent triage with per-comment human approval → published as files in git → rendered statically at next deploy. Author-only pending state, never public unmoderated text. Comments close automatically on posts older than 12 months (OQ-8); closed threads keep their existing comments and replace the form with a one-line explanation.

### 5.7 Syndication

#### FR-20 · Feeds
RSS (`/feed.xml`) and JSON Feed (`/feed.json`), full content, not excerpts. Auto-discovery links in `<head>`. A per-tag feed is v2.

#### FR-21 · Canonical and cross-posting
Every post carries a self-referential canonical URL. Cross-posts to dev.to, Hashnode, LinkedIn or Medium set canonical back here. Cross-posting is manual in v1 and a candidate for agent automation in v2.

#### FR-22 · Social cards
Open Graph images generated at build time by the Go binary from post metadata — title, series, date, in the site's own typography. No third-party card service. This is itself a dogfooding post.

---

## 6 · Non-functional requirements

#### NFR-1 · Performance budget
Enforced in CI. A pull request that breaches any line fails.

| Metric | Budget |
|---|---|
| LCP, mid-tier mobile on 4G | < 1.5s |
| CLS | < 0.02 |
| INP | < 100ms |
| JS per article page | ≤ 30 KB gzipped total |
| CSS | ≤ 15 KB gzipped, critical path inlined |
| HTML per article, pre-images | ≤ 60 KB |
| Web fonts | ≤ 2 files, subset, woff2, self-hosted, preloaded |
| Third-party requests, first load | 0 |
| Lighthouse, all four categories | ≥ 98 |
| Full build, 200 posts | < 15s |

#### NFR-2 · Accessibility
WCAG 2.2 AA as a quality floor: semantic landmarks, visible focus, 4.5:1 contrast, keyboard-operable everything, `prefers-reduced-motion` respected, alt text required by the build for content images. Automated axe checks in CI plus a manual keyboard pass before each release. See §9 on whether this is also a legal obligation — build to it either way.

#### NFR-3 · SEO
Clean semantic HTML, one `h1` per page, descriptive titles under 60 characters, meta descriptions from `summary`, JSON-LD (`Article`, `Person`, `SoftwareApplication`, `Comment`), canonical URLs, generated sitemap, no orphan pages. No keyword engineering — the long tail here is genuinely specific technical language, which the content produces naturally.

#### NFR-4 · Browser support
Last two versions of Chrome, Firefox, Safari, Edge, plus Safari iOS. No IE, no polyfill bundle. Baseline CSS features only where support is universal; progressive enhancement above that.

#### NFR-5 · Availability
Static assets served from Cloudflare's edge. A Worker outage disables forms and leaves the site fully readable. A GitHub outage blocks publishing but not serving. No single dependency can take the content offline.

#### NFR-6 · Portability
Content in Markdown, comments in Markdown, config in one file. Migrating off Cloudflare means changing a deploy step. Migrating off the generator means rewriting templates. Neither touches the content. Verified by a documented dry run in v2.

#### NFR-7 · Cost
€0/month recurring beyond domain registration. Every dependency sits inside a free tier at expected volume and fails closed rather than into billing. Any change that introduces a recurring cost requires an explicit decision recorded here. Note the domain line is no longer trivial: `.ai` runs roughly €50–80/year against ~€10 for a `.com`, and is billed in mandatory two-year blocks.

---

## 7 · Agentic development

The site is built and operated by Claude Code and Codex. That is a requirement on the repository, not a description of a habit.

#### AR-1 · One command per operation
`make dev` · `make build` · `make check` · `make publish` · `make comments` · `make deploy`. Every routine action is one of these. An agent that has to compose four commands will get one wrong.

#### AR-2 · Deterministic builds
Sorted map iteration, stable IDs, no timestamps in output except real content dates. A build that produces different bytes from identical input makes every diff untrustworthy and sends agents chasing ghosts.

#### AR-3 · Golden-file tests
Rendered output for a representative set — long post, short post, post with a series, post with comments, post with no comments, tool page, empty tag page — committed to `testdata/` and diffed on every run. This is the highest-leverage single investment in the repository: it converts "did the vibe-coded change break the site" into a red/green signal an agent can iterate against without supervision.

#### AR-4 · Type-safe templates
templ, so that a wrong field name is a compile error rather than a blank space in production. Compile errors are the feedback loop agents self-correct against.

#### AR-5 · CI gates
Every pull request runs: build, golden-file diff, `go vet`, link checker (internal and external), HTML validation, axe accessibility scan, performance budget check, and a grep for secrets and email addresses. Green means deployable.

#### AR-6 · Preview deploys
Every pull request gets a real URL via `wrangler versions upload`, `noindex` headers, drafts included. Review happens on the artifact, not the diff.

#### AR-7 · Agent contract
`AGENTS.md` at the repository root, with `CLAUDE.md` symlinked to it. Contains: architecture in one page, the command list, content conventions, the tag vocabulary, the tone guide, and the prohibitions below.

#### AR-8 · Prohibitions
Agents **must not**, without explicit per-instance human approval: modify `/impressum/` or `/datenschutz/`, publish a comment, publish a post (drafts only), write anything to `content/` that contains a raw email address, add a runtime dependency or a paid service, add a third-party script, force-push, or rewrite history.

#### AR-9 · Rollback
`make rollback` reverts to the previous deployment in one step, without a rebuild. Documented and tested before launch, not discovered during an incident.

---

## 8 · Security and privacy

#### FR-23 · Headers
CSP with no `unsafe-inline` for scripts, `frame-ancestors 'none'`, HSTS, `X-Content-Type-Options: nosniff`, `Referrer-Policy: strict-origin-when-cross-origin`, restrictive `Permissions-Policy`. Delivered via `_headers` and verified in CI.

#### FR-24 · Secrets
Repository secrets in GitHub Actions and Worker secrets in Cloudflare only. Nothing in the tree, no `.env` committed, a pre-commit hook and a CI grep as backstops. The repository is public (D-5), so this is a hard boundary rather than hygiene.

#### FR-25 · Data collected
Comments (name, email, body), contact form (name, email, message), newsletter (email). Nothing else. No IP addresses stored raw anywhere; only salted hashes for rate limiting, salt rotated monthly.

#### FR-26 · Third-party processors
Cloudflare, Resend, the newsletter provider, GitHub. Each listed on the privacy page with a DPA in place. Adding a processor is a decision recorded in §3, not an implementation detail.

---

## 9 · Legal (Germany)

*Not legal advice. The operator is EU-based and offering commercial services, so these are real obligations. Have a Fachanwalt review before launch — the fee is smaller than one Abmahnung.*

#### LR-1 · Impressum
Mandatory. <cite index="30-1">The obligation now sits in § 5 DDG, which replaced § 5 TMG on 14 May 2024 with the mandatory contents essentially unchanged</cite> — <cite index="29-1">the substantive requirements carried over and only the terminology shifted from "Telemedien" to "digitale Dienste".</cite> Must be reachable in two clicks from every page, with a working email address and a means of direct contact. <cite index="27-1">There is no legal requirement to cite the statute in the Impressum at all, and omitting the reference avoids the risk of an Abmahnung for citing a repealed one.</cite>

#### LR-1a · Brand is not a legal person
`Siduri` is a trading name. Whatever the Impressum and the privacy notice require, they name the natural or legal person behind it, not the brand. Registering the domain does not create an entity and does not satisfy any disclosure obligation. See OQ-3 and LR-7.

#### LR-2 · Datenschutzerklärung
Covers comments, contact form, newsletter, Turnstile, hosting and log data, plus data subject rights and the supervisory authority. Note that <cite index="30-1">the former TTDSG has been renamed TDDDG, with the cookie-consent rule now at § 25 TDDDG and unchanged in substance</cite> — references must be updated accordingly.

#### LR-3 · Consent
D-3 removes the need for a consent banner by removing non-essential storage. The only client-side storage is the pending-comment ID in `localStorage`, argued as strictly necessary for a service the reader explicitly requested. Document that reasoning; revisit if any analytics or embed is ever added.

#### LR-4 · Newsletter
Double opt-in with logged confirmation, one-click unsubscribe in every message, Impressum in the footer of every message. Non-negotiable under German practice.

#### LR-5 · Accessibility obligation
<cite index="41-1">The BFSG has applied since 28 June 2025, but purely informational or image websites without contract conclusion generally fall outside its full scope</cite>, and <cite index="42-1">microenterprises providing services — under ten employees and no more than €2M in turnover or balance sheet total, both conditions together — are exempt.</cite> <cite index="36-1">Purely B2B offerings are also outside its scope, provided it is clearly evident that nothing is sold to consumers.</cite> A solo B2B consultancy with no online booking or checkout is therefore very likely out of scope — but NFR-2 stands regardless, because the cost of building to AA from the start is near zero and retrofitting is not. Reassess if online booking or a paid digital product is ever added.

#### LR-6 · Business status
Confirm whether the activity registers as freiberuflich or gewerblich, whether Kleinunternehmerregelung applies, and whether VAT must be shown on `/services`. Deferred with the offer itself (OQ-1), since there is nothing to invoice yet. Must be settled before phase 3.

#### LR-7 · Operator identity disclosure — **resolved in principle**
Two facts settle this, and neither depends on where the operator lives.

**Jurisdiction follows the operator, not the reader.** § 5 DDG binds providers established in Germany. Relocate and it stops applying; the destination's equivalent applies instead. GDPR is the exception that reaches across borders — it binds non-EU controllers who offer goods or services to EU residents or monitor their behaviour — but a blog that sells nothing and tracks nobody is outside it. The site as scoped for P0 and P1 collects no personal data at all: static pages, feeds, and a `mailto:` link. There is nothing to disclose because nothing is collected.

**The mailing list is the trigger, not the country.** Commercial email requires a valid physical postal address in the footer under US CAN-SPAM and under German practice alike. Any competent provider will require one before the first send. This obligation attaches to running a list, wherever the list is run from, and no amount of relocating removes it.

Therefore:

- **P0 and P1 ship anonymously and legitimately.** No forms, no list, no comments, no analytics. RSS and JSON Feed carry distribution. The "just a blog" framing is accurate rather than a fiction.
- **P2 requires a mailbox, not a confession.** A registered PO box or a commercial mail-receiving agency address, €100–200/year, in whichever country the operator is in at the time. Home address stays private. Country becomes visible; street does not.
- **Business registration (OQ-3) is a separate and later question.** It attaches to invoicing, not to publishing, and can wait for P3.

Revisit only if the site starts selling to EU residents from outside the EU, at which point GDPR reattaches regardless of the operator's location.

## 10 · Measurement

#### FR-27 · Analytics
Cloudflare Web Analytics: free, cookieless, no client-side identifier, no consent implication. Aggregate only.

#### FR-28 · Conversion visibility
`/services` and `/contact` visits, and inquiry count from the Gmail label, reviewed monthly against §1 targets. Deliberately coarse: the site does not have the traffic for anything finer to be meaningful, and pretending otherwise wastes hours on noise.

#### FR-29 · Review cadence
Monthly: posts shipped, inquiries, one thing to change. Quarterly: is the content pillar mix (§2) actually holding, or has the buyer-legible quota quietly gone to zero.

---

## 11 · Cost model

Figures in EUR, converted from USD list prices where applicable; providers adjust quarterly, so verify at checkout. Ranges are wide where a plan tier is involved.

### 11.1 Kick-off — one-time

| Item | Cost | When | Note |
|---|---|---|---|
| `siduri.ai`, 2-year block | **€100–165** | Now | Registry mandates 2-year increments; Cloudflare Registrar near-wholesale |
| `siduri.dev` defensive registration | €12–15/yr | Optional, now | Canonical fallback; denies a confusable neighbour |
| Trademark search, DIY | €0 | Before the 2-year term is paid | TMview, EUIPO, DPMA, USPTO all free to search |
| Legal pages via generator | €0–60 | P0 | See LR-1, LR-2 |
| Postal mailbox setup fee | €0–50 | P1 | Some providers charge a one-off |
| Fachanwalt review | €300–800 | **Defer to P3** | Only once actually selling |
| Trademark registration | €290 DPMA (3 classes) / ~€1,050 EUIPO | **Defer to P3** | Only if the brand proves worth defending |
| Business registration | €0 freiberuflich – €60 Gewerbe | **Defer to P3** | OQ-3 |

**Minimum to be live: €100–165.** That is the domain; everything else in P0 is free or deferred.

### 11.2 Recurring — monthly

| Line | P0–P1 | P2 | P3 |
|---|---|---|---|
| Cloudflare — Workers, static assets, KV, Turnstile, Web Analytics | €0 | €0 | €0 |
| GitHub Actions (public repo, D-5) | €0 | €0 | €0 |
| Resend — 3,000/mo, 100/day | — | €0 | €0 |
| Buttondown — free ≤100 subs, ~€8 at 1,000, ~€26 at 5,000 | — | €0 → €8 → €26 | €8–26 |
| Gmail | €0 | €0 | €0 |
| Domain, amortised over the 2-year term | €4–7 | €4–7 | €4–7 |
| Postal mailbox (LR-7, FR-18) | — | €8–17 | €8–17 |
| **Infrastructure subtotal** | **€4–7** | **€12–50** | **€20–50** |
| AI coding subscription | €18–90 | €18–90 | €18–90 |
| **Total** | **€22–97** | **€30–140** | **€38–140** |

Two observations that matter more than the totals.

**The infrastructure is genuinely free.** Excluding the AI subscription — which is paid regardless of this site, for work that is not this site — Siduri costs €4–7/month in P0 and €12–50/month fully built. NFR-7 holds.

**Buttondown's free tier expires by design.** The §1 target is 250 subscribers in six months against a 100-subscriber free cap, so the €8 tier should be expected around month four. Not a problem; not a surprise either.

### 11.3 First-year cash

| Scenario | Year 1 |
|---|---|
| P0/P1 only, AI subscription already paid for other work | **€110–230** |
| Full build through P2, AI subscription already paid | €250–500 |
| Full build through P2, AI subscription attributed here | €470–1,600 |

### 11.4 Time — the actual cost

| Item | Hours |
|---|---|
| P0 build — scaffold, templ layout, CI with §6 budget gates, golden tests, legal pages, deploy | 25–40 |
| Post 1 — introduction, vision, experience | 6–12 |
| P1 — Gateslot plus post 2 as build log | 25–70, depending entirely on the tool |
| P2 — comment system per `comments-requirements.md` | 20–35 |
| **Ongoing** — 4 posts/month plus comment triage | **18–34 per month** |

At a notional €80/hour, P0 alone represents €2,000–4,000 of foregone billable time, which makes the €165 of infrastructure a rounding error.

This reframes the commitment: not a €165 project with a hosting bill, but a **100–160 hour build with a standing ~20 hour/month content obligation.** That is what should be signed off, and it is why OQ-7 — the cadence floor that survives a busy client month — is the most consequential open question in this document.

### 11.5 Cost discipline

The realistic path to an unexpected bill is adding a service casually: a form backend, an analytics SaaS, a CMS, a managed database. AR-8 forbids an agent doing it; this section is the human-facing equivalent. Any new recurring line requires a decision recorded in §3.

## 12 · Phasing

Renumbered as phases rather than versions, because the gates are content events rather than feature sets.

### P0 · Foundation
Build pipeline, article rendering, `/`, `/journal/`, `/about/`, `/contact/` as `mailto:`, legal pages, 404, sitemap, robots, RSS + JSON feed, deploy, CI with the §6 budget gates. No services page, no forms, no comments, no newsletter.

**Post 1 — introduction, vision, experience.** Who, why, what this is going to be, and what the operator has already done. The thesis post. Ships with the site; the site is not live without it.

**[AMENDED 2026-08-25 — ADR 0003, operator's word given directly]**

> ~~**Gate: post 1 is published on a real domain before any P1 work begins.** This is the single most important line in the document. The characteristic death of this project is six weeks of beautiful scaffolding and zero published posts — and building it with agents makes scaffolding more seductive, not less.~~
>
> **Gate struck.** P0, P1 and P2 are built in one programme rather than gated on
> publication. The reasoning the struck text carries is not withdrawn and is not
> restated as a rule: the risk it names — scaffolding instead of publishing, made
> worse rather than better by building with agents — is now carried by the
> end-to-end acceptance in ADR 0003 and by nothing else. Post 1 remains P0 scope
> and the site is not live without it (§12, P0).

### P1 · Gateslot
The first free tool: its own page under `/tools/`, a public repository, and **post 2** as the build log. This is the first real proof and the first thing a stranger can use without trusting anyone.

Also in P1, off the critical path: arrange a postal mailbox per LR-7 and FR-18, so P2 is not blocked on paperwork. This is needed wherever the operator is living by then, so it is not contingent on the relocation question.

> **Gate: Gateslot is usable by someone who is not the operator, and post 2 is published.**

### P2 · Audience
Newsletter (Buttondown), comments, contact form, search, tags, `/stack/`, OG image generation, `/llms.txt`. Every item here collects or serves an audience, so none of it is worth building before there is one — and none of it ships before the mailbox exists.

> **Gate: a postal address exists. No list, no form and no comment intake before it does.**

### P3 · The offer
`/services`, price, contact as a sales channel, `/changelog/`, cross-post automation, portability dry run. Opens when OQ-1 is one sentence.

## 13 · Acceptance criteria

- [ ] A clean checkout builds to `dist/` with one command, no network, no manual steps.
- [ ] Two consecutive builds from identical input produce byte-identical output.
- [ ] Every page validates, passes axe, and is fully keyboard operable.
- [ ] The performance budget in NFR-1 passes in CI on a real article with code blocks and images.
- [ ] The site is fully readable with JavaScript disabled, including navigation and article content.
- [ ] Every post has a `plain_summary`; the build fails on a post that doesn't.
- [ ] An unknown tag fails the build.
- [ ] A draft post appears nowhere in `dist/`, verified by grep.
- [ ] Impressum and Datenschutz are reachable in two clicks from every page.
- [ ] No cookie banner exists, and no non-essential client-side storage is written.
- [ ] `grep -rE '[a-z0-9._%+-]+@[a-z0-9.-]+' content/` returns only intended addresses.
- [ ] Contact form works with and without JavaScript, and mail lands in Gmail within 60s.
- [ ] Pull requests produce a working preview URL with `noindex`.
- [ ] `make rollback` restores the previous deployment and has been tested at least once.
- [ ] `AGENTS.md` lets a fresh agent session publish a post without asking a question.
- [ ] No `/services` page, price, or sales copy exists anywhere in `dist/` before P3.
- [ ] `dist/` contains no form element and no third-party script before P2.
- [ ] Comments are automatically closed on posts older than 12 months (FR-19).

---

## 14 · Out of scope

Multi-author support · i18n and translated content · a CMS or admin UI · user accounts · paid content or memberships · e-commerce or checkout · a booking calendar (would change the §9 accessibility analysis) · dark-mode toggle beyond `prefers-color-scheme` · AMP · a mobile app · a Discord or forum · client portals · podcast hosting.

---

## 15 · Open questions

Status as of this revision. Resolved items are kept rather than deleted, so the reasoning survives.

#### OQ-1 · What exactly is being sold? — **deferred by design, blocks P3 only**
No offer yet. Sequence agreed: introduction post → Gateslot (free tool) + build log → sell once there is something to sell. This no longer blocks launch, because launch no longer contains a services page (D-9). It blocks P3 and nothing else. The answer will be shaped by which parts of Gateslot people actually use and ask about — which is a better input than guessing now.

#### OQ-2 · Price, and does it appear on the page? — **deferred with OQ-1**
Recommendation stands: a starting number or a band, on the page. Decide before FR-15 is written, not during.

#### OQ-3 · Freiberuflich or gewerblich, Kleinunternehmerregelung? — **deferred with OQ-1, and possibly moot**
Nothing to invoice yet, and the operator may not be in Germany by the time there is. Business registration attaches to invoicing rather than publishing, so it waits for P3 and is answered in whichever jurisdiction applies then. Distinct from LR-7, which is now resolved on its own terms.

#### OQ-4 · Newsletter provider — **resolved: Buttondown**
See FR-18. Self-hosting rejected on cost and ops; building it rejected as scope, revisit as a P3 build log. Carries one prerequisite that is now on the critical path for P2: a postal address for the email footer.

#### OQ-5 · Domain and name — **RESOLVED: Siduri, at `siduri.ai`**
Brand rather than personal name: agreed. The human-in-the-loop family is rejected as a *name* (D-8) — occupied since 2017 by a Sofia annotation company and again by a Belgian AI-agent startup as of February 2026, weak as a mark with two US applications on the exact phrase, and unrankable against wall-to-wall definitional content from Google Cloud, GeeksforGeeks and a dozen vendors. The concept stays as the thesis and tagline, where it is accurate and does real work.

Direction chosen: classical, drawn from the *Iliad* and *Odyssey*. The thematic centre is **machine power, human steering** — which is both the architecture and the positioning.

| Name | Source | Meaning | Case | Risk |
|---|---|---|---|---|
| ~~**Phrontis**~~ | *Od.* 3 | Menelaus's helmsman; φροντίς = thought, care, attention | Best metaphor on the list — **but withdrawn, see below** | Occupied on three fronts |
| **Thole** | Homeric rowing | The pin fixing oar to boat; archaic English *to thole* = to endure | The pivot converting raw force into direction. Five letters, unsaturated | Obscure; pronunciation not obvious |
| **Moly** | *Od.* 10 | Herb from Hermes that makes Odysseus immune to Circe's transformation | Handle a transformative power without being transformed. Four letters | Reads slight; "holy moly"; molybdenum slang |
| **Scheria** | *Od.* 8 | Phaeacian home, whose ships need no helmsman — they read men's minds and sail alone; Poseidon petrifies one for it | Autonomous agents described in 800 BC, punishment included. Best About-page story available | Needs telling; opaque cold |
| **Lemnos** | *Il.* 18 | Hephaestus's forge, where he built speaking golden handmaidens and self-moving tripods | First thinking machines in Western literature, built by a craftsman working beside them | Real island; some existing commercial use |
| **Winnow** | *Od.* 11 | Teiresias sends Odysseus inland with an oar until someone calls it a winnowing fan | The reference *and* the literal act of separating signal from noise — also exactly the §2 buyer-legibility problem | Ordinary English word, hard to own |
| **Nostos** | throughout | νόστος, the homecoming; root of *nostalgia* | Arriving. Shipping. Best pure sound here | Weakest conceptual fit |
| **Weft** | *Od.* 2, 19 | Penelope's crossing thread — woven by day, unpicked by night | Iterate, refactor, iterate. Four letters | Obscure; textile associations |

**Phrontis: withdrawn.** An earlier assessment called it unclaimed in software. Diligence contradicted that on three fronts:

- **PHRONTIS LIMITED**, UK company 02915860, incorporated 6 April 1994 and still active, registered in Adderbury, Oxfordshire. A management consultancy and research practice in systems thinking, with SIC codes covering other publishing (58190), other IT service activities (62090) and management consultancy (70229). Overlapping field of use, current accounts, and it holds `phrontis.com`, which now serves a placeholder page rather than its former paper library. Not a dormant squatter and not a cheap acquisition.
- **A US brand family** already spans at least `phrontisleadnetworks.com` (AI lead scoring), `phrontisfinancial.com`, `phrontisfuel.com` and `phrontisacquisition.com`.
- **`phrontis.ai` sits beside `phronetic.ai`**, an operating agentic-AI company. Near-homophone, same category — the worst available option rather than the escape hatch. Related: `phronesis-ai.com` is parked for sale, and `github.com/phronesis-io` publishes an agent harness built on Claude Code.

Exact registration date and registrar for `phrontis.com` were not verified; .com registrant data is redacted post-GDPR, though creation date remains public at any registrar.

**Current front-runners: Scheria, Thole, Moly.** None surfaced a collision during the Phrontis diligence, but none has been vetted directly. Method note, learned here: a first-pass search under-counts occupation. Vet the chosen name against (a) national company registers, (b) trademark databases in the relevant classes, (c) the exact-string domain family beyond the obvious TLD, and (d) near-homophones in the same category — before buying anything.

**Cross-mythology longlist.** Phrontis fuses a named helmsman with a noun meaning *thought-and-care*. Analogues split along that seam; few traditions fuse both. Recorded unvetted.

| Tradition | Name | Sense | Note |
|---|---|---|---|
| Norse | ~~Huginn~~ | Odin's raven *Thought*; he fears for its return — thought plus care | **Unusable.** An existing decade-old open-source agent platform: agents that read the web, watch events and act on your behalf along a directed graph, including a HumanTaskAgent. Same metaphor, same product, 2013 |
| Norse | ~~Mímir~~ · ~~Forseti~~ | The counselling head; the one who presides over disputes | **Both rejected — see Norse note below** |
| Egyptian | **Sia** | Personified perception, riding in Ra's solar barque | Insight as crew aboard the vessel. Likely too short to own |
| Egyptian | Seshat · Mahaf | Measurement and record; the ferryman who must be answered correctly | |
| Mesopotamian | **Siduri** | The alewife at the world's edge who counsels and points to the ferryman | Six letters, clean, apparently unsaturated |
| Mesopotamian | **Urshanabi** | Utnapishtim's boatman across the Waters of Death | Long but entirely unclaimed |
| Irish | Manannán / *Sguaba Tuinne* | A boat that sails where its owner wills, no sail or oar | The autonomy motif again — pairs with Scheria |
| Brittonic | **Barinthus** | Navigator who ferried Arthur to Avalon, knowing waters and stars | Pure helmsman, highly distinctive |
| Sanskrit | Sārathi | Charioteer — Krishna holds the reins, counsels, refuses to fight | Arguably the most precise human-in-the-loop image in world myth; common personal name in India |
| Sanskrit | Buddhi · **Viveka** | The deciding faculty; discernment of real from apparent | |
| Chinese | Zhǐnán chē | South-pointing chariot, geared to hold bearing through every turn | Concept, not a name |
| Finnish | Väinämöinen | Sings a boat into being, lacks three words to finish it | Concept, not a name |
| Zoroastrian | Vohu Manah | "Good Thought" as a divine aspect | |
| Micronesian | **Etak** | Wayfinding frame: the canoe holds still, the islands move past | Conceptually superb; living tradition, borrow with care |
| Roman | Palinurus | Helmsman who sleeps at the tiller and is lost | The inverse of phrontis |

**Decision: Siduri, `siduri.ai`.** Human-in-the-loop retained as the principle (D-11), not the name. Comparison that settled it:

| | Siduri | Etak |
|---|---|---|
| Occupant | Siduri Wines (est. 1994, sold to Jackson Family Wines 2015) — a Sonoma Pinot Noir producer inside the largest wine company headquartered in Sonoma County | Etak Inc. (1983) — Stan Honey and Nolan Bushnell's digital-mapping firm; its Navigator was the first commercially viable car navigation system, on the cover of *Popular Science* in June 1985 |
| Field overlap | None. Wine | **Direct.** Navigation technology, Silicon Valley |
| Trademark | Clear in Classes 9 / 35 / 42; a Class 33 wine mark does not reach software. Jackson Family litigates marks, but only within wine | Likely lapsed, association permanent |
| Metaphor free? | Yes | **No** — Etak Inc. was itself named after the wayfinding term, described in its own early literature as navigating intuitively by the stars. The metaphor was taken in 1983 |
| Domain | Six letters, obscure referent. Plausible on `.com`, likely on `.dev` / `.studio` | Four letters — `.com` certainly gone |
| Pronunciation | si-DOO-ree; spells as it sounds in English and German | Ambiguous EE-/EH-; parses as "e-" + "tak", a 1999-era e-brand; sits near *attack* |
| Cultural risk | Low — no living tradition | Real — Carolinian wayfinding is living practice, and Etak Inc. misattributed it as Polynesian. Either inherit the error or carry the correction |
| Extensibility | Deep — Urshanabi, Utnapishtim, the whole epic | Would draw further on the same living tradition |

**Why it fits the business.** Siduri keeps a tavern at the edge of the world. Gilgamesh arrives half-mad with grief, hunting immortality; she tells him he will never find it — go home, eat, wash, hold your child — and then tells him where the ferryman is anyway. The counsellor who is honest about the quest and still helps with the next step. That maps to the consultancy better than the helmsman framing ever did, and it gives `/about` a reason to exist.

**Outstanding, non-blocking:** a Class 9 / 35 / 42 trademark search in the relevant registries, and confirmation that no software or consulting entity trades under the name. The wine mark sits in Class 33 and does not reach software. Worth completing before the two-year `.ai` term is paid, since that term is the commitment.

**Naming family now available for tools.** The epic supplies a coherent well: *Urshanabi*, the ferryman who is the only one able to carry you across, for anything that moves or migrates; *Utnapishtim* for archival or long-term storage. Optional — see OQ-9.

**Norse candidates: all three rejected on diligence.**

- **Huginn** — three independent occupants, two directly in category. The 2013 open-source agent-automation platform (agents that read the web, watch events and act on your behalf; hundreds of contributors). `huginn.sh`, live, selling "your AI coding team… persistent AI agents that remember your codebase." And Nosto's October 2025 launch of *Huginn*, an AI agent for commerce that orchestrates specialized agents under the brand's supervision. Catastrophic.
- **Mímir** — Grafana Mimir: flagship open-source horizontally scalable long-term storage for Prometheus and OpenTelemetry metrics, successor to Cortex, at v3.0 as of KubeCon NA 2025. The collision sits precisely inside the secondary audience — infrastructure-literate developers hear "metrics backend" and nothing else.
- **Forseti** — Google and Spotify's open-source GCP security toolkit, 2017–2023. Archived by Google for low community engagement; repos read-only, `forsetisecurity.org` still held by the verified org. Name technically vacated, but six years of Google Cloud documentation and blog posts are permanent. Inheriting an abandoned Google project's search results is a liability, not an opening.

**Rule derived: skip Norse entirely.** It is the most exhausted naming vein in software — Grafana alone runs both Loki and Mimir, Munin (Huginn's sibling raven) has been a monitoring tool for two decades, and Odin, Thor, Heimdall, Bifrost, Yggdrasil, Sleipnir, Fenrir, Asgard and Valhalla are all multiply spoken for. Three good candidates, three collisions, is the base rate rather than bad luck. Greek is the second most exhausted vein, which is what killed Phrontis. Prefer the under-mined traditions: Mesopotamian, Brittonic, Micronesian, and the corners of Greek that are neither gods nor heroes.

Not names, but keep: **Palinurus** and **Elpenor** as a postmortem series — both lost because nobody was watching. **Väinämöinen's three missing words** for shipping. **Wheelwright Bian**, who tells the Duke the skill in his hands cannot be written in a book, for the tacit-knowledge problem at the centre of vibecoding.

Rejected on saturation: the Polytropos/Polytropo family (`polytropo.io` is a container-storage company, plus at least three others), Metis, Argus, Talos, Aegis, Daedalus, Mentor.

Registration status was not verified here — check candidates at a registrar directly. Four-letter `.com` domains have been fully allocated for years, so any such candidate is taken regardless.

#### OQ-6 · Is client work publishable? — **deferred**
Revisit before the first client engagement, not after. Default options remain: anonymized case studies, own-tools-only, or a publication clause negotiated into the contract. The third is easiest to get agreed *before* a project starts and nearly impossible after.

#### OQ-7 · Post cadence under load — **still open**
Not answered. It matters less in P0/P1, where the plan is two posts and a tool, and a great deal in P2 onward. The question to settle: what is the monthly floor that survives a busy client month — one post? — and is the two-drafts-in-reserve buffer policy accepted?

#### OQ-8 · Comment freeze — **resolved: 12 months**
Comments close automatically on posts older than 12 months. Recorded as FR-19 here and in the comment spec. Closed threads still render their existing comments and show a short line explaining why the form is gone.

#### OQ-9 · Working name of the free tool — **still unconfirmed**
Recorded as **Gateslot** throughout. Still flagged in case it was shorthand rather than a name. Now also worth deciding whether tools take names from the epic (see OQ-5) or stay plainly descriptive. Either is defensible; mixing them arbitrarily is not. Note that "Gateslot" happens to describe an approval gate, which sits well with D-11 regardless of the source.
