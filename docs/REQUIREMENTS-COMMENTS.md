# Comment System — Requirements

**Project** · **Siduri** (`siduri.ai`) — vibecoding journal & services site (Go + templ + htmx, static, Cloudflare Workers)
**Component** · Reader comments with human-approved static publishing
**Principle** · Human in the loop — this component is its flagship demonstration
**Version** · 0.1 — draft for review
**Date** · 2026-08-25
**Status** · Pending sign-off

---

## 1 · Context

The site is a statically generated Go/templ site deployed to Cloudflare Workers static assets, built and shipped by GitHub Actions, and developed agentically with Claude Code / Codex.

Its business purpose is grassroots marketing of vibecoding services. Content is the journey: tools built, tools used, dogfooding, what broke. Comments are therefore **not a community feature**. They are three things, in priority order:

1. **Credibility.** A post with a real thread under it reads as lived work, not content marketing.
2. **SEO surface.** Statically rendered comments and author replies are indexable long-tail text on every article.
3. **Lead intake.** A commenter describing a problem they're stuck on is a qualified prospect, and the reply is the first sales touch.

This framing drives most of the decisions below — particularly why comments are files in git rather than rows in a database, and why the author reply is a first-class object rather than an afterthought.

The comment pipeline is also, deliberately, dogfooding: it is an agentic workflow built by an agentic workflow, and it becomes article material.

It carries a second load. Human-in-the-loop is the brand's guiding principle (site spec D-11), and this is where it runs in public: an agent triages, classifies and drafts at machine speed, and nothing reaches a reader until a person has approved it one item at a time. Every design choice below that looks like friction — per-comment approval, no batch accept, the pending state shown to its author — is the principle being load-bearing rather than decorative. If this component ever auto-publishes, the argument the business sells is refuted by its own comment section.

---

## 2 · Decisions

Locked-in choices that constrain everything downstream. Each records the alternative rejected.

#### D-1 · No third-party comment widget
Disqus, Giscus, Utterances and friends are rejected. Third-party JS, third-party tracking, third-party outage, and — for a site whose pitch is "I build my own tools" — badly off-brand.

#### D-2 · Approved comments are files in git
The repository is the single source of truth for published content. A comment becomes real when it is a committed file, not when a row flips a boolean. Rejected: a database as system of record. Consequence: publishing is a deploy, not a write.

#### D-3 · Pending comments are visible to their author only
Never render unmoderated text to the public, not even briefly. The commenter sees their own submission held in a pending state; nobody else does. Rejected: optimistic public display with retroactive takedown.

#### D-4 · Approval is per-comment and human-gated
The agent triages, classifies, drafts, and prepares. A human approves each comment individually before it is written to the repo. Rejected: auto-approve above a confidence threshold.

#### D-5 · Email is the notification channel, KV is the queue
Gmail gets a readable notification per comment. Cloudflare KV holds the authoritative pending record. Rejected: email as the only store — it makes the pending state non-durable and makes reconciliation guesswork.

#### D-6 · No accounts, no OAuth, no login
Name plus email, email never published. Rejected: any identity provider. Friction kills comment volume on a small site, and account data is liability.

#### D-7 · Commenter email addresses never enter the repository
The repo is public. Records committed to git carry a name and an opaque hash, never a raw address. Raw addresses live only in KV (with TTL) and in the Gmail inbox.

---

## 3 · Architecture

```
                 ┌──────────────────────────────────────────┐
   reader ──────▶│  Static article page (pre-rendered HTML)  │
                 │  approved comments already baked in       │
                 └───────────────┬──────────────────────────┘
                                 │ htmx POST /api/comments
                                 ▼
                 ┌──────────────────────────────────────────┐
                 │  Cloudflare Worker  (same origin)         │
                 │  Turnstile verify → honeypot → ratelimit  │
                 └───────┬──────────────────────┬───────────┘
                         │ write                │ send
                         ▼                      ▼
              ┌────────────────────┐   ┌──────────────────┐
              │ Workers KV         │   │ Resend  ──▶ Gmail│
              │ pending:<slug>:<id>│   │ structured mail  │
              └─────────┬──────────┘   └────────┬─────────┘
                        │                       │
                        │      ┌────────────────┘
                        ▼      ▼
              ┌──────────────────────────────────────┐
              │  Agent session  (Claude Code / Codex)│
              │  triage → classify → ONE AT A TIME → │
              │  ask human → approve / reject        │
              └─────────────────┬────────────────────┘
                                │ writes content/comments/**
                                ▼
              ┌──────────────────────────────────────┐
              │  git commit → push → GitHub Actions  │
              │  build → wrangler deploy → purge KV  │
              └──────────────────────────────────────┘
```

**Cost-critical detail.** Adding this component means the project is no longer assets-only: it gains a Worker script alongside the static assets. Requests that match a static asset are still served directly and remain free — the Worker is only invoked for paths with no matching asset, i.e. `/api/*`. Do not enable `run_worker_first`, or every page view starts billing against the Worker request quota.

---

## 4 · Functional requirements

### 4.1 Submission

#### FR-1 · Comment form
Every article page renders a comment form below the approved thread. Fields: display name (required, ≤60 chars), email (required, ≤120 chars, never published), website (optional, ≤200 chars), comment body (required, ≤2000 chars). Plus a hidden honeypot field and a Turnstile widget. The form works without JavaScript as a plain HTML POST; htmx is progressive enhancement only.

#### FR-2 · Submission endpoint
`POST /api/comments`, same origin, `application/x-www-form-urlencoded`. Returns an HTML fragment, not JSON — htmx swaps it directly into the page. On success the fragment is the pending-state confirmation card; on failure it is the form re-rendered with an inline error and the user's input preserved.

#### FR-3 · Body format
Comment bodies are stored as plain text. On render, a restricted Markdown subset is applied: paragraphs, line breaks, `code`, fenced blocks, bold, italic, and links. No images, no raw HTML, no headings, no tables. All links get `rel="nofollow ugc noopener"`. Maximum 2 links per comment; a third causes rejection at submit time with a clear message.

### 4.2 Spam defence

#### FR-4 · Layered filtering at the edge
Applied in order inside the Worker, cheapest first, each rejecting outright:

| Layer | Mechanism | Failure response |
|---|---|---|
| 1 | Honeypot field non-empty | Silent 200, discard |
| 2 | Form age < 3s or > 6h (HMAC-signed timestamp in the form) | Inline error |
| 3 | Cloudflare Turnstile server-side verify | Inline error |
| 4 | Rate limit: 3/hour and 10/day per hashed IP | Inline error, retry-after |
| 5 | Length caps, link count, unicode/script heuristics | Inline error |
| 6 | Static denylist of domains and phrases | Silent 200, discard |

Layers 1 and 6 discard silently so bots get no signal. Everything reaching KV is presumed *plausible*, not *good*.

#### FR-5 · Classification at the agent stage
The agent applies a final semantic pass before presenting anything to a human, labelling each comment `spam` / `low-value` / `genuine` / `lead`. This is advisory. It reorders and annotates the review queue; it never decides.

### 4.3 Pending state

#### FR-6 · Author-only pending display
On successful submission the client stores the returned comment ID in `localStorage` under `cmt_ids`. On every article page load, if IDs are present for that post, the page requests `GET /api/comments/pending?post=<slug>&ids=<csv>` and appends the returned fragments below the static thread, visually marked as awaiting review.

#### FR-7 · Pending card content
The pending card shows the comment as it will appear, a clear status label, and a plain-language explanation: comments are read by a human and appear at the next site publish. No fake ETA. If the reader clears site data or opens the page elsewhere, the pending card is gone — this is acceptable and is stated in the confirmation copy.

#### FR-8 · Deduplication on publish
Rendered static comments carry `data-comment-id`. The pending renderer drops any ID already present in the DOM. This makes the flow safe against a failed KV purge: worst case is a stale key, never a duplicate on screen.

### 4.4 Notification

#### FR-9 · One email per comment
Sent via Resend to a single Gmail address. Subject line: `[comment] <post-slug> · <display name>`. Body is human-readable — enough to judge the comment from a phone without opening anything.

#### FR-10 · Machine-parseable payload
The email body ends with a fenced YAML block containing the complete record. The agent parses this block and nothing else, so notification formatting can change freely without breaking the workflow. Gmail is a fallback ingestion path; KV is primary.

#### FR-11 · Digest fallback
If more than 20 comments arrive in a rolling 24h window, the Worker switches to a single hourly digest email to stay clear of provider daily caps. Individual records remain in KV regardless.

### 4.5 Approval and publishing

#### FR-12 · Approved comments become files
One file per comment: `content/comments/<post-slug>/<ulid>.md`, YAML frontmatter plus the body. Sorted lexically by ULID, which sorts chronologically for free. One file per comment, not one file per post — clean diffs, no merge conflicts, and each approval is an isolated, reviewable change.

#### FR-13 · Author replies
An author reply is a comment record with `author_role: site` and a `parent_id`. Threading is exactly one level deep: a reply to a reply is not supported and the UI offers no affordance for it. Replies are drafted by the agent and approved by the human on the same terms as any other comment.

#### FR-14 · Static rendering
The build reads `content/comments/**`, renders threads into each article page, and emits `schema.org/Comment` structured data. Comment count appears in the page metadata. A post with zero comments renders the form and no empty-state noise.

#### FR-15 · Queue reconciliation
After a successful deploy, and only after, CI deletes the published IDs from KV. Order is non-negotiable: publish, verify, then purge. A failed purge is a cosmetic problem (FR-8 covers it); a purge before a failed deploy loses a reader's comment.

#### FR-22 · Comment freeze
Comments close automatically on posts older than 12 months. The build omits the form on frozen posts and renders a single line explaining that the thread is closed; existing comments continue to render normally. The Worker rejects submissions for frozen slugs regardless of what the client sends, since the absent form is a UI convention and not a control.

#### FR-16 · Publish notification
When a comment goes live, optionally email the commenter a link to it, if they ticked the box at submission. Off by default. The email contains a one-click removal link with an HMAC token.

### 4.6 Data protection

*Not legal advice — validate with counsel. The operator is EU-based, so GDPR and TTDSG apply and this is a real constraint, not boilerplate.*

#### FR-17 · Lawful basis and notice
The form carries a short, plain-language notice at the point of submission: what is collected, why, how long it is kept, and that the name and comment body become public. Consent is the basis; submitting is the consent action. Link to the full privacy page.

#### FR-18 · Data minimisation
Store: display name, email, comment body, timestamp, post slug, optional website. IP addresses are **never stored raw** — only a salted hash, used for rate limiting, with the salt rotated monthly. No analytics cookies, no third-party embeds inside comment bodies.

#### FR-19 · Retention
KV records carry a 60-day TTL. Rejected and spam-classified comments are deleted at triage. Published comment files retain a name, body, timestamp and an `email_hash` — never the address itself (see D-7).

#### FR-20 · Erasure
A documented route to request removal, honoured within 30 days: delete the file, rebuild, deploy. The removal link in FR-16 is the self-service path.

#### FR-21 · Processors
Cloudflare and Resend are processors. Confirm DPAs are in place and list both on the privacy page.

---

## 5 · Data model

### KV record — `pending:<post_slug>:<ulid>`

```yaml
id: 01K3QZJ8X4YB7N2M9V0PQRSTUV   # ULID, sortable by time
post_slug: go-htmx-static-comments
author_name: Jane D.
author_email: jane@example.com    # KV + Gmail only, never git
author_website: https://example.com
body: |
  Multi-line plain text, unrendered.
parent_id: null                   # set for replies
created_at: 2026-08-25T14:03:22Z
ip_hash: 9f2a...                  # salted, rate limiting only
turnstile_ok: true
user_agent_family: firefox
notify_on_publish: false
status: pending                   # pending | approved | rejected | spam
```

### Repository record — `content/comments/<post-slug>/<ulid>.md`

```yaml
---
id: 01K3QZJ8X4YB7N2M9V0PQRSTUV
author_name: Jane D.
author_website: https://example.com
author_role: reader               # reader | site
email_hash: sha256:4c1a...        # dedup + optional avatar, not reversible to display
parent_id: null
created_at: 2026-08-25T14:03:22Z
approved_at: 2026-08-25T19:41:00Z
approved_by: human
---

Multi-line plain text body.
```

Note the asymmetry: the KV record is operational and expires; the git record is content and is permanent. The email address exists only on the left-hand side of that boundary.

---

## 6 · Agent workflow

Invoked as `/comments` in Claude Code (mirror as `AGENTS.md` for Codex). This is the human-in-the-loop core of the system and its contract is stricter than the rest of the spec.

### 6.1 Sequence

1. **Fetch.** List and read pending records from KV via `wrangler`. If KV is unavailable, fall back to parsing the YAML blocks from unread Gmail messages matching `subject:[comment]`.
2. **Enrich.** For each comment, load the article it belongs to and pull the relevant passage. A comment can't be judged without the thing it replies to.
3. **Classify.** Assign `spam` / `low-value` / `genuine` / `lead`, with a one-line reason. Sort the queue: leads first, spam last.
4. **Present one comment.** Article title, the passage being responded to, the comment body verbatim, the classification and its reasoning, link targets expanded, and — for `genuine` and `lead` — a drafted author reply.
5. **Wait.** Offer: approve · approve with edits · approve + publish reply · reject · mark spam · skip. Do not proceed until the human answers.
6. **Act.** On approval, write the comment file (and the reply file if approved separately). On reject or spam, record the ID in a local ledger so it is not re-presented, and delete the KV key.
7. **Repeat** until the queue is empty.
8. **Ship.** One commit for the session — `comments: publish N (slug, slug)` — push, watch the Actions run, report the result. On green, purge the published IDs from KV.

### 6.2 Constraints

The agent **must not**:

- Approve, publish, or commit any comment without explicit per-comment human input. Silence is not approval. A single "approve all" instruction is not per-comment input — refuse it and ask for the queue to be walked.
- Modify a comment body beyond whitespace trimming and stripping tracking parameters from URLs. Any semantic edit is shown as a diff and separately approved.
- Write a raw email address anywhere under `content/`.
- Touch any path outside `content/comments/` during a comment session.
- Delete a KV key before the corresponding deploy is confirmed green.
- Force-push, amend a pushed commit, or rewrite history.
- Fabricate a reply in the site author's voice and publish it as approved. Drafts are drafts until a human says otherwise.

### 6.3 Reply drafting

Drafted replies are the marketing surface, so tone matters. Match the site's register: direct, technical, specific, no gratitude filler. Answer the question that was actually asked. Reference the concrete tool or post where relevant. Never pitch services in a reply — the reply demonstrating competence *is* the pitch. If a comment is a clear buying signal, flag it to the human for a private email instead of a public reply.

---

## 7 · Non-functional requirements

#### NFR-1 · Cost
€0/month recurring beyond the domain. Every component sits inside a free tier at expected volume, and every component must fail closed rather than fail into billing.

#### NFR-2 · Performance
Article pages remain fully static: no blocking request, no layout shift from the pending fetch, which loads after paint into reserved space. Comment submission responds in under 500ms p95.

#### NFR-3 · Availability
A Worker outage degrades to a read-only site with a disabled form and an honest message. It never breaks article rendering, because articles do not depend on the Worker.

#### NFR-4 · Portability
Because published comments are files (D-2), the entire comment corpus survives dropping Cloudflare, Resend, or the Worker entirely. The lock-in surface is the pending queue only, which is at most 60 days of data.

#### NFR-5 · Agent legibility
Every part of this system must be operable by an agent from the repo alone: one command per operation, deterministic output, golden-file tests over rendered comment threads, and the workflow contract in `AGENTS.md`. If an agent needs tribal knowledge to run it, the design is wrong.

---

## 8 · Cost model

| Component | Free tier | Expected load | Headroom |
|---|---|---|---|
| Workers — static assets | Unmetered | All page views | ∞ |
| Workers — script requests | 100,000/day | `/api/*` only | Very large |
| Workers KV | 1,000 writes/day, 100,000 reads/day, 1GB | 1 write per comment | Very large |
| Turnstile | Free | 1 verify per submit | ∞ |
| Resend | 3,000/month, **100/day**, 1 domain | 1 email per comment | Binding constraint |
| Gmail | Free | Inbox | — |
| GitHub Actions | Unlimited public / 2,000 min private | Seconds per build | Large |
| Domain | — | — | ~€10/year |

**The binding constraint is Resend's 100/day cap**, which is why FR-11 exists. Nothing else in this system comes close to a limit at realistic volume for a personal site. Verify current tiers before launch; provider pricing moves.

---

## 9 · Acceptance criteria

- [ ] A comment submitted with JavaScript disabled is accepted and confirmed.
- [ ] A comment submitted with htmx swaps in a pending card without a page reload.
- [ ] The pending card is visible to its author on reload, and absent in a private window.
- [ ] An unapproved comment appears nowhere in `dist/`, verified by a build-time grep test.
- [ ] Turnstile failure, honeypot fill, and rate-limit breach each produce the correct response per FR-4.
- [ ] A notification email arrives in Gmail within 60s, with a YAML block that parses.
- [ ] `/comments` presents comments one at a time and refuses a blanket approve-all.
- [ ] Approving a comment produces exactly one new file under `content/comments/` and no other change.
- [ ] `grep -rE '[a-z0-9._%+-]+@[a-z0-9.-]+' content/comments/` returns nothing.
- [ ] After deploy, the comment renders statically with `schema.org/Comment` markup.
- [ ] KV keys for published comments are gone after a green deploy, and only after.
- [ ] A KV purge failure produces no duplicate comment on the page.
- [ ] Golden-file tests cover thread rendering, including a one-level reply and an empty thread.
- [ ] Removal request path deletes the file and the comment is gone after the next deploy.

---

## 10 · Out of scope

Threading deeper than one level · reactions and voting · comment editing by readers · comment search · an RSS feed of comments · webmentions and ActivityPub · a moderation web UI (the agent session *is* the UI) · avatars fetched from third parties · real-time updates · comment counts on index pages.

---

## 11 · Open questions

1. **Comment volume assumption.** The design targets fewer than 20 comments/day. If the site takes off, does the flow still hold, or does it need batch approval with sampled review?
2. **Preserve the pending state across devices?** Currently `localStorage`, so it is per-browser. A signed cookie or an emailed status link would fix it. Worth the complexity, or is the confirmation email enough?
3. **Publish rejected comments' existence?** Some sites show "1 comment held for review." Honest, but invites gaming.
4. ~~**Comment freeze.**~~ **Resolved: 12 months.** See FR-22.
5. **Is the author reply a comment or a post update?** For genuinely good questions, a reply may deserve promotion into the article body. Who decides, and does the comment stay?
6. **Turnstile vs. GDPR.** Turnstile is privacy-friendlier than reCAPTCHA but still a third-party call. Confirm it belongs under strictly-necessary or gets its own notice.
7. **Build trigger.** Does an approval push deploy immediately, or batch into a scheduled daily publish? The latter is calmer and makes "appears at the next publish" a real, predictable promise.
