# Decisions

One row per decision. `siduri_planner` proposes, `siduri_reviewer` agrees or
dissents, the operator decides and breaks ties. Reviewer agreement is a
precondition on what may be proposed — never a substitute for approval, and
nothing reaches the tree because two sessions agree.

Agreement is recorded **with its scope**: what the reviewer checked, and what
they did not.

| # | Date | Decision | Planner | Reviewer | Operator |
|---|---|---|---|---|---|
| 1 | 2026-08-25 | Commit identity is `Siduri <siduri@siduri.ai>` | proposed, from D-5 (public repo) + LR-7:396 (P0/P1 anonymous) | agreed; flagged that `siduri.ai` may not be registered yet | **decided** |
| 2 | 2026-08-25 | Default branch is `main`, set explicitly | proposed | raised it — fresh `git init` here gives `master` | routine, taken |
| 3 | 2026-08-25 | `AGENTS.md` is prohibitions-first, against AR-7's stated order | proposed, ADR 0002 | raised it; truncation eats the tail | **deferred — decide after wave 0.** Stays prohibitions-first meanwhile; ADR 0002 stays *proposed* |
| 4 | 2026-08-25 | Build order across P0/P1/P2 vs the §12 gates | proposed a `dist/` phase gate | **refuted it** — phase gate guards `dist/`, gate 490 guards the tree | **struck 490.** One programme, P0+P1+P2. ADR 0003 |
| 5 | 2026-08-25 | Spec filenames | renamed to `REQUIREMENTS-*.md` | **refuted it** — four in-document refs dangled, incl. front matter and the D-5 consequence | renamed back; zero contract bytes changed, no amendment needed |
| 6 | 2026-08-25 | Insert lane W1 — widen the seam — before wave A | proposed: eight lanes cannot run against finding 0006 | pending | pending |
| 7 | 2026-08-25 | `siduri-code` skill, canonical in `skills/`, deployed to `~/.claude/skills/` | written, v1.0, 66 clauses, every one earned today | pending | pending |
| 8 | 2026-08-25 | `siduri-code` v1.1 — add `## What this enforces`, `ST-13`, restore `ST-03`'s tell, convert `FO-12` to a script | v1.0 shipped 66 clauses with nothing distinguishing the guarded from the unguarded | **refuted v1.0's honesty** — named the 15 that go red nowhere, and refused `ST-13` alone as making the file less honest | pending |
| 9 | 2026-08-25 | `siduri-code` v1.2 and `lane_overlap.py` v2 — two thresholds, behaviour mode, `FO-13`/`FO-14`/`FO-15` | v1 shipped as a gate with four false negatives | **refuted it** — attacked the tool they proposed, four lanes reporting cuttable | pending |
| 10 | 2026-08-25 | Insert W2 — body-extension seam — and re-run A1, A2, A7 against it | proposed: three lanes blocked, five unaffected and still running | pending | pending |
| 11 | 2026-08-25 | A2 merged at `5c55ab9`; wiring deferred to after W2 | verified by planting, not by report | pending | pending |
| 12 | 2026-08-25 | A6 merged at the second attempt after the legal link-graph fix | verified by walking links, not grepping | pending | pending |
| 13 | 2026-08-25 | **Publish post 1** — `draft: true` → `false` | proposed with the gate evidence; agents draft and never publish | not consulted on this one | **approved explicitly**, in these words: *publish post 1 - we will replace it with real post later* |
| 14 | 2026-08-25 | Wave A closed: 11 lanes merged, `main` green, FR-12 wired | verified by planting every Done-when on the merged tree | pending | pending |
| 15 | 2026-08-25 | Gateslot's repository link is dropped until the repo is public | offered three options with FR-13 and the P1 gate against each | not consulted | **decided**: drop the link until it is public |
| 16 | 2026-08-27 | **Design feedback arrives as annotations, not screenshots** — `agentation-vanilla` pasted into the console against `make dev`, capture half only, `agentation-mcp` never run. ADR 0006 | proposed after running the tool over all 31 pages: 81 of 849 selectors resolve to a *different* element | **agreed on all five points and the recommendation**; re-derived every number on a second engine and found two the planner had not — the unrooted chain (0061) and the seven-day destructive expiry (0062) | **decided**, in one word: *YES* |
| 17 | 2026-08-27 | **Wave W1 — the §13 acceptance pass.** Six Codex lanes, a `make acceptance` gate deriving its set from the contract by heading, and the operator's own pass over fifteen rendered pages | proposed pages-first after the operator's ruling; cut six lanes, not eight | agreed the cut; corrected the harness three times — heading anchor not line range, exit code split three ways, and the pinned count replaced by a two-method cross-check | **decided**: *go*. Merged as PR #6 |
| 18 | 2026-08-27 | **The legal pages stay built and unlinked, the footer is dropped, and every repeated element loads from one source** | recorded from the operator's own annotations during his pass | not consulted — taken live during the sitting | **decided**, in his words: *drop the footer, keep the pages hidden and unlinked*, and *ALL REPEATED elements must be loaded from the same source, not reinvented* |
| 19 | 2026-08-25 | Session closed with the site unreviewed; design edits become the next wave | operator's call, after looking at the running site | not consulted | **decided**: *Site will require a lot of edits, but we are stopping here* |

| 20 | 2026-08-27 | **The amendment deadlock: the guard learns the shape of a legitimate amendment.** `Makefile:30` refused every committed diff to the two requirements files, so a branch amending the contract could never pass `check` and therefore never merge — including an amendment the operator approved. Replaced by `tools/amendcheck.py`: an added ADR with a `**Status** ·` line, every removed line matched injectively to a `~~ ~~` span in the same hunk, every touched heading named in the ADR, and watched contract paths present at base and HEAD. ADR 0013 | put three options and recommended A | agreed A; found nine holes in the first implementation, including non-injective per-hunk matching and the shallow-checkout vacuity | **decided**: *A*. He then ran the blocked guard-swap commit himself as `301957a`, after the classifier refused it twice to the planner |
| 21 | 2026-08-28 | **Post-meta is date and tags. Read time is dropped.** | proposed date-only, from his first-pass annotation *Date only, drop read-time and tags* | not consulted — taken live during the sitting | **decided**, revising his own earlier ruling: *Date and tags, read time is too posh* |
| 22 | 2026-08-28 | **Comment threading goes to two levels, readers may reply anywhere, and `FR-1`/`FR-2` are amended to match.** Six clauses struck and restated, new clause `FR-13a · Reader replies`, ADR 0014 | proposed one extra level, scoped to replies addressed to the commenter | agreed; corrected three of the planner's own instructions — a scoped assertion satisfiable by the JSON-LD, a depth check that would have given unbounded depth if lifted, and a contradictory `LR-1` strike. **After merge, found the shipped implementation drops a comment silently — finding 0087** | **decided** in three steps: *allow 1 more level*, then *Readers can reply anywhere*, then *Go all the way now* |
| 23 | 2026-08-28 | **A wave ships as one branch and one pull request, not one per lane.** `lane/w3` carried five commits from three lanes | the planner had been handing over a push box per lane | not consulted | **decided**, in his words: *can't you make SINGLE PR???* |
| 24 | 2026-08-28 | **Anything the operator is asked to look at is handed over as an exact clickable URL on the host he holds, never a path and never a port this machine chose.** Recorded as `AS-09` | proposed after handing him a tunnel command forwarding a port already taken on his laptop, then telling him to open the wrong one | not consulted | **decided**: *give me exact links to click when you are asking me to look into something! save it into the skill and wherever!* |
| 25 | 2026-08-28 | **The tools filter nav is raised as its own review task rather than folded into a design lane.** The task has not been created — finding 0082 | recorded from his annotation *not sure — add the task to review later* | not consulted | **decided**: raise it separately |
| 26 | 2026-08-28 | **Session close-out must deposit artefacts, not describe duties.** Measured: 15 operator verdicts with 0 probes, no findings row since 0078, no decision row for four rulings, one skill clause existing only in the deployed copy, and `AGENTS.md` carrying a false claim about its own guard. Backfilled as rows 0079–0088 and decisions 20–26 | proposed five parts; withdrew one on finding `tools/acceptance.py` already enforces `{ran, at, saw}` with four planted breaches | **agreed the diagnosis and part 3; dissented on sequencing.** Named a destructive ordering defect — `skill-install` would overwrite the only copy of `PR-11a`; corrected the baseline argument from *prove red* to *prove it can also go green*, one passing fixture per rule; ruled guard-plus-backfill one PR under `FO-16`; collapsed the manifest into part 3 with a checked/unchecked split. Scope: re-derived all five measurements, the greps, the `docs/` guard by planting, A9 under amendcheck, and the drop by fixture. Not checked: `hitl`'s install direction, the templ and CSS diffs, `make check` on `14f5296`, and the backfill content | pending |

| 27 | 2026-08-28 | **Wave W4 — close-out that deposits artefacts.** Three Codex lanes: the silent comment drop fixed by a reachability pass; `tools/acceptance.py` extended to walk the operator's own reports with ten fixtures; `make skill-install`/`skill-check` with a provenance manifest and eight fixtures. Plus `Makefile:32` extended to `SKILL.md`, and the records backfill — findings 0079–0108, decisions 20–27, probes on both acceptance reports | proposed five parts; withdrew one on finding the mechanism already existed, and rewrote all four task files after review | **agreed and gated every step.** Found something blocking in each of the four task files before they were cut — an unbounded walk fed its own cycle input, a guard rejecting the repo's own records convention, a one-shot red proof, and a `docs/` exception contradicting `AGENTS.md`. Verified the merged tree, including four things nobody had run on it: determinism over 41 files, the acceptance invariant, fixture leakage, and B1's property test against the pre-fix implementation. Scope: not checked — B1's and B2's report prose, findings text beyond ids, and whether CI reproduces the 32 s gate time | **decided**: *fix plus guard*, one PR. Merged as `c19c51e` |
| 28 | 2026-08-28 | **The enforcement section is rewritten to four states, not two.** `FO-02` and `ST-06` are bound by a machine check, `Makefile:32`; `FO-12`/`13`/`14` are bound with a named false negative; eight are self-reported fields nothing reads; `CT-04` named a mechanism that does not exist. No total is given — the categories are the record, and a total is the figure that rots | proposed the rewrite after the audit | **audited all ten and supplied the four-state shape.** Also declined the authority the lead asserted on the operator's behalf — `RV-05`, an arrangement relayed by the lead is a claim until the operator states it directly — while noting the change did not need it, since `CT-01` scopes operator approval to the two requirements files | **pending** — see open question below |

## Open

- **Does `siduri-reviewer` gate changes on the operator's behalf, or only screen
  them before they reach him?** The lead told the reviewer that the operator had
  delegated approval of skill and rules changes to them. The operator said that
  to the lead, never to the reviewer, who declined to act on it — `RV-05`, an
  arrangement relayed by the lead is a claim until the operator states it
  directly. The reviewer's reading is that `RV-01` makes them a precondition on
  what the lead may propose, never a substitute for approval. Nothing has moved
  on the disputed reading: `CT-01` scopes operator approval to the two
  requirements files, and every skill change this wave was a skill file. The
  question is live for the next thing that touches the contract. Decision 28.

- **Should a skill edit be unable to merge until it has been deployed?**
  `skill-install` could write the canonical tree digest into
  `skills/MANIFEST.tsv` at deploy time, and CI could recompute and compare —
  catching *canonical moved and nobody redeployed*, which is the state this box
  is in whenever the skill is edited. It needs no live host and pins no foreign
  checksum. The cost is a real workflow constraint: a pull request editing a
  skill stays red until someone runs `skill-install` and commits the receipt,
  which makes *deploy before merge* mechanical rather than remembered. Proposed
  by `siduri-reviewer`, finding 0094. The other direction — the deployed copy
  edited directly, finding 0085's actual shape — is unobservable from CI and no
  mechanism helps.


- **`FR-13` is a deliberate partial** while Gateslot's repository stays private,
  and `site-requirements.md:499` — *Gateslot is usable by someone who is not the
  operator* — cannot be met until it is public. Recorded so a green tool page is
  not read as a met requirement.

- **Post 1 is a placeholder by the operator's own word** and is to be replaced
  with a real post. It is published, so replacing it is an edit to live content
  rather than a draft change.

- **Wire A2's related posts into the article page** once W2's section seam lands
  (finding 0027). One registration in an A2-owned file; no rebuild.

- **Dependencies before branches.** Every library the eight lanes need resolved
  into `go.mod` before any branch is cut; no lane runs `go get` (`FO-15`).
  Cannot be done until W1 lands, since W1 holds `go.mod`.

- **§12:497**, the P1→P2 gate — *Gateslot is usable by someone who is not the
  operator, and post 2 is published.* Not raised with him, not ruled on, still
  standing. Same class as the struck 490. Blocks wave B; blocks nothing before it.
- **§12:502**, the postal-address gate. Standing, and it should — it gates
  collecting personal data, not building the code that would.
- **ADR 0002**, `AGENTS.md` ordering. Deferred to after wave 0 by his ruling.

- **The annotation channel's selector field is a hint, not an identity.**
  `**Text:**` and `**Selected text:**` are sound; `**Selector:**` is wrong on 250
  of 849 elements across our pages (findings 0060, 0061). Two defects, both
  upstream, both one small change. Reported to the author with the operator's
  approval on 2026-08-27; if the fix lands, this open item closes and the
  selector becomes usable as written.

- **Nothing asserts the built route set** (finding 0064). `dist/` is a total map
  from URL to screen in one direction only — a screen that stops rendering
  vanishes with the gate green and no diff. Not scheduled; recorded so the
  walked tree is not read as coverage of what *should* render.

- **`AGENTS.md:90-91` overstates its own guard** (finding 0063). It claims
  `make check` rejects any change under `docs/`; `Makefile:30` guards the two
  requirements files. Either the sentence narrows to match the guard or the
  guard widens to match the sentence. `AGENTS.md` is not the contract, so this
  is the operator's call and not an amendment.

- **The acceptance is not accepted.** `make acceptance` on `main` reports
  `held=8 deferred=5 open=3 failed=2` over eighteen criteria. Criterion 1 fails —
  a clean checkout does not build. Criterion 18 fails — the twelve-month comment
  freeze tests a hardcoded date literal, so *automatically* is false. Criteria 3c,
  8 and 10 are open. Recorded so a green `make check` is not read as an accepted
  phase.

- **`siduri.ai` is not available and the build asserts it 272 times** (finding
  0066). Four unlinked sites hold the domain, one of them a security guard. The
  name decision is the operator's and blocks nothing until publish.

- **One chrome, one source** (finding 0067) is decision 18's implementation and
  has no lane yet. It subsumes the footer removal, the unlinking, the two content
  widths, and the unstyled tag pages.

- **`LR-1` and `LR-2` are not struck.** The operator chose to unlink rather than
  amend, so both clauses stand and both pages still ship. Recorded so the next
  reader does not infer an amendment that was never made.

- **The keyboard pass has still never been performed** — `NFR-2:294`, criterion
  3c, owner operator. Nothing automated substitutes for it.


## Reviewer scope, as of 2026-08-25

Checked: environment claims, spec citations, `AGENTS.md` against the codex
binary, worktree and identity mechanics, both scope assumptions.

Not checked: the operator's "everything before commercialization" sentence
(never seen it), `REQUIREMENTS-COMMENTS.md` beyond two cited lines, the eight
lane definitions for `RV-23` disjointness — they do not exist yet.

## Reviewer scope, as of 2026-08-27 — decision 16 only

Checked by running: the tool fetched fresh, size, sha256 and `LICENSE`;
network-call greps; the mount point and shadow root; storage load and save;
`getSelector` read whole; selector resolution over all 849 elements of all 31
pages in jsdom, unpatched and patched; landmark nesting across all 31 pages by
parser; npm metadata for all three package names.

Not checked: the skill's `references/measured.md`; **the tool's interactive
path** — both sessions evaluated its selector function and neither has clicked
its UI or seen a clipboard block emitted, so the export format is read at source
and nowhere else; Chrome-versus-jsdom divergence beyond the two runs agreeing on
all five numbers; whether the author takes the issue.

## Reviewer scope, as of 2026-08-27 — wave W1

Checked by running: the live domain, its title, RDAP record and MX; every
`siduri.ai` occurrence in `dist/`; criterion 1 by running `make build` with
`templ` absent; body-child structure of all 31 rendered pages by parser; the
selector algorithm's resolution over all 849 elements in jsdom, patched and
unpatched; `tools/acceptance.py` anchors, `check()` and selftest; `mk/accept.mk`;
`tools/secretscan.py` exemption logic; the criterion 11 evidence command, which
was broken by the integrator's own normalisation and caught by re-running it.

Not checked: `references/measured.md` in the annotation skill; the annotation
tool's interactive path — the selector function was evaluated, the tool's UI was
never clicked; seventeen of the eighteen fragment rows against their own `ran`
lines; whether `make check` was green on the integration branch, which the
integrator ran and the reviewer did not.
