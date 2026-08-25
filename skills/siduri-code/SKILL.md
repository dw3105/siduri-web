---
name: siduri-code
description: Working rule for the Siduri site build, repo siduri-web: two specs are the contract, eight Codex lanes run parallel and read no skill, the merged tree is the only verdict. Read before acting.
---

**v1.1 - 25 Aug 2026.** Derived from `gateslot-code` v1.13. Structure kept,
content earned in one day: repo genesis, one spine lane, one seam lane, nine
findings, four of them against this session's own claims. File caveman full.
Output stay **lite**. Change log live in `ledger.md`, never here.

**Canonical copy is `skills/siduri-code/` in `siduri-web`.**
`~/.claude/skills/siduri-code` is deployed copy. **No install target and no drift
check exist yet** - copy is by hand, and `deployed == canonical` is true today and
held by nothing. `hitl` bind it with `make skill-check`. Finding 0013.

**This file bind Claude session only.** Codex worker read task file. Clause meant
to reach worker go in task text or in guard, never here. `FO-01`.

**Read reference at moment it name, that turn.** Rule here is index. Failure
shape live only in reference.

| Reference | Read when |
|---|---|
| `references/contract.md` | requirement touched, ADR written, requirement look wrong |
| `references/fanout.md` | cutting lane, writing lane task, merging wave |
| `references/proving.md` | writing check, planting breach, calling something caught |
| `references/shipping.md` | cutting commit, writing command block |
| `references/closing.md` | closing wave |
| `references/ledger.md` | rule look wrong, new rule proposed |

Repo fact live in `AGENTS.md` - command, layout, prohibition. Skill govern how
work move; `AGENTS.md` govern what worker may touch. Both read.

## skill check

`skill check` → this line, nothing else:

```
siduri-code v1.1: two spec files are the contract and only operator amend them | Codex worker read task file and never this skill, so clause reaching worker live in task text or in guard | AGENTS.md truncate at 32768 silently and drop end, so prohibition go first and live in file worker cannot edit | lane own only new file, and file every lane need belong to integrator | registration happen when data exist, never before | worker run no gate, integrator run it once on merged tree | branch green alone say nothing about branch merged | check green at baseline certify nothing, prove it red first | agent report is claim, read tree | document is not machine fact, probe this turn | rename inside contract dangle reference the contract cannot fix | reviewer agreement is precondition on proposal, never substitute for approval | agreement carry scope, what was checked and what was not | one decision per turn, box last, Your move and Mine | fifteen clause here go red nowhere, and section naming them is the only thing telling reader which sixty-six are which
```

`skill check <ID>` → that rule line, verbatim. `skill check refs` → every filename
in reference table. Cannot recall → `skill check: FAILED, reattach skill`.

## Style

- **ST-01** **Two registers, and file decide which.** Turn, task text, ADR, commit body, finding row - caveman **lite**. This file and `references/` - caveman **full**. Off only "stop caveman" / "normal mode".
- **ST-02** Caveman lite is full sentence, article, technical term, code, command, error string, filename all survive.
- **ST-03** Scan before send. Cut opener, filler (just, really, basically, actually, simply), hedge, closer, tool narration. **One filler = scan not run.**
- **ST-04** **Every reference carry its content same sentence, every turn.** `FO-04` is nothing; *FO-04, registration happen when data exist* is readable.
- **ST-05** **Fact leaving session carry where it was read from.** Ref for file, command for query, unit for number.
- **ST-06** **This file stay under 32768 byte.** Same cap that eat `AGENTS.md`. Skill over cap is skill that truncate wherever it is pasted, and `gateslot-code` sit at 32394 with 374 spare.
- **ST-13** **Rule phrased as advice fire while auditing, never while designing.** Write what diff can be checked against, never what to be careful about. Read `## What this enforces` before adding clause: fifteen here already fail this.
- **ST-07** **`SKILL.md` is index. Rule is one line, one claim, under 300 char.** Shape live in reference.

## Read first

**RD-01** Never start from title. Read all four.

- `docs/site-requirements.md` - every requirement in scope, whole, by line range
- `docs/comments-requirements.md` - same, where comments in scope
- `AGENTS.md` - whole, every session
- `docs/DECISIONS.md`, `docs/FINDINGS.md`, `docs/adr/` - what is already ruled

- **RD-02** **Probe this turn beat every document.** Undated machine claim is hypothesis. Measured: friction row three day old name a sysctl as cause; sysctl read `0` and shape exit clean. Document promoted to machine fact is how it happen.
- **RD-03** **Never restate machine fact or code fact from memory.** Three sibling repo carry three different git identity; one recalled cover all three is wrong twice.
- **RD-04** **Claim that thing is built cite file and line, or say unprobed.** Grep hit is not substrate.
- **RD-05** **Read requirement whole, by range.** grep locate; grep never read.

## Contract-first

`docs/site-requirements.md` and `docs/comments-requirements.md` outrank this
skill, `AGENTS.md`, any plan, and anything said in session.

**CT-01** **Only operator amend requirement, by explicit word given directly.**
Never relayed through peer, never inferred from review. Finding produce ADR marked
*proposed* plus question. Never edit.

- **CT-02** **Amendment is three part in one commit**: ADR naming clause, reason, replacement; clause struck in place and never deleted; every affected check edited.
- **CT-03** **Struck gate get its risk relocated, named, in same ADR.** Gate removed with nothing behind it is failure that gate predicted.
- **CT-04** **No worker edit `docs/`.** Bound by path check, never by prose. Worker habit is to fix document found wrong, and that habit amend contract.
- **CT-05** **Renaming file inside contract dangle every reference to it.** Contract cannot be edited to repair it, so rename is not free. Restore name instead.
- **CT-06** **Requirement with no check that can go red is unimplemented.**

## Fan-out

Eight lane run parallel. Every rule here come from that and from nothing else.

- **FO-01** **Codex worker read task file, never skill, never `CLAUDE.md`.** Clause meant to bind worker go in task text or become guard. Clause only lead read bind nobody in eight lane.
- **FO-02** **`AGENTS.md` reach worker regardless of trust, and truncate at 32768 byte in silence.** End go first. Prohibition therefore go first and live in file worker cannot edit; guard is byte count, never intention.
- **FO-03** **Lane own only new file. File two lane need belong to integrator.** Merge conflict-free by construction, never by care.
- **FO-04** **Registration happen when data exist, never before.** Route registered at `init()` cannot know tag, tool or thread. Symptom look like output cardinality; cause is lifecycle.
- **FO-05** **Value unique across lane is allocated in task, never chosen by worker.** Counter claimed on unmerged branch is not claimed.
- **FO-06** **Task deliver as file worker read. Operator paste pointer, never task.**
- **FO-07** **Worker run no gate.** Eight full gate on four core is box thrashing, and gate cannot see conflict from inside one lane anyway. Integrator run it once, on merged tree.
- **FO-08** **Branch green alone say nothing about branch merged.** Eight branch is 28 pair.
- **FO-09** **Integration is named lane with own task, never residue.** Only clean fast-forward with green gate stay in lead's hand.
- **FO-10** **Wave get one Done-when that is end-to-end and performed by person.** Read all lane Done-when and ask which name product. All naming file mean nobody test composition.
- **FO-11** **Worker report is claim. Read tree.** Not summary, not test count, not *gate passed*.
- **FO-12** **Filename named by two lane task is seam defect, and it is grep.** `tools/lane_overlap.py` run over task set before any worktree cut; non-empty output mean wave not ready. Fix seam, never route lane around it.

## Proving

- **PR-01** **Check green at baseline certify nothing.** Prove red first, name input that turn it red. Measured: disjointness proof passed against un-widened seam, because worked example already in tree satisfied it.
- **PR-02** **Prefer check that delete workaround over check that stand beside it.** `grep -c ArticlePage build.go` returning `0` is red today and green only when mechanism real.
- **PR-03** **`git diff` see no untracked file.** Proof counting new file report zero, and report shared file it never asked about. Use `git status --porcelain`.
- **PR-04** **Count in Done-when is ownership, never fixed number.** Templ lane is three file, Go-only lane is two.
- **PR-05** **Proof land as committed check, never as demonstration then deleted.** Deleted demonstration bind memory.
- **PR-06** **Plant every breach, watch it red, put it back.** Report which input fired.
- **PR-07** **Determinism proved by recursive checksum over tree, never by eye.**
- **PR-08** **Defect invisible to single run need test that run twice.** Package-global state survive build; one build see nothing.
- **PR-09** **Silence mean broken until something prove otherwise.** Print value to compare, never `ok`.

## Reviewer

- **RV-01** **Reviewer agreement is precondition on what lead may propose, never substitute for what operator approve.** Nothing reach tree because two session agree.
- **RV-02** **Agreement carry scope: what reviewer checked, what they did not.** Agreement with no scope become laundering three week later.
- **RV-03** **Reviewer claim is claim. Re-derive load-bearing half, say which half.**
- **RV-04** **Reviewer is never channel to operator.** Finding reach operator under reviewer's own name, unresolved.
- **RV-05** **Arrangement relayed by lead is claim until operator state it to reviewer directly.**
- **RV-06** **Decision land row in `docs/DECISIONS.md`.** Message thread bind nothing.

## Shipping

Dev host commit. No remote, nothing pushed. Read `references/shipping.md` before
cutting one.

- **SH-01** **One lane, one commit.** Subject assigned by lead, used verbatim.
- **SH-02** **`make check` before every commit on `main`.**
- **SH-03** **Worktree cut with absolute path.** Relative path resolve against `-C` and nest worktree inside repo, where every gate walk it.
- **SH-04** **Identity set on parent before first worktree exist.** Worktree share parent `.git/config`. Unset mean every worker die at first commit after doing all work.
- **SH-05** **Branch name stated as fact in task.** Fresh `git init` here give `master`.
- **SH-06** **Merged worktree removed and branch deleted at merge.**
- **SH-07** **Commit message two `-m` flag.** Subject under 72 char, body max six line.
- **SH-08** **Every box carry machine line, box-n-of-m line, command, expected value as trailing comment.** No prose inside fence.
- **SH-09** **Command shipped only after running here, or after reading its `--help` on this binary.**
- **SH-10** **Watched command get no pipe, no redirect, no `tee`.**
- **SH-11** **Pipeline exit status belong to last command.** Write to file, audit by outcome.

## Ask

- **AS-01** **One decision per turn. Size is not unit.**
- **AS-02** **Blocker asked at once, never batched. Unblocked work continue meanwhile.**
- **AS-03** **Question carry its own weakness, and rejected option shown ruled out.**
- **AS-04** **Recommendation mandatory, with its reason.**
- **AS-05** **Operator hold no context.** Turn carry every fact needed, in that message.
- **AS-06** **Frame what break, never what changed.**
- **AS-07** **Decision takeable without opening laptop.**
- **AS-08** **Question that contradict operator's own instruction go to him as contradiction, named.** Contract he wrote can forbid order he asked for.

## Block pre-flight

**Before every block, each one, in turn that ship it.** Blank or `no` **stop block**.

```
block pre-flight
  machine       this box, and what it cannot run: <why>
  ran here      every command run here, or --help read on this binary (SH-09)
  worktree      absolute path, outside repo                          (SH-03)
  identity      parent carry user.name and user.email                (SH-04)
  branch        named as fact, not assumed                           (SH-05)
  paste         pointer to task file, never task text                (FO-06)
  watched       no pipe, no redirect, no tee                         (SH-10)
  verdict       no $? after pipeline                                 (SH-11)
```

## Pre-flight

Before every commit on `main`. Blank or `no` **stop it**.

```
pre-flight <commit subject>
  requirements  <ids> read whole this turn, with line ranges
  refs read     <file>:<line range> → "<phrase from section bearing on this>"
  lane          <id>, base commit, previous wave merged: yes/no
  done-when     <n> lines in task, <n> proved here, <n> deferred by name
  breaches      <n> planted, <n> went red, <n> restored          (PR-06)
  baseline      every check proved red before counted green      (PR-01)
  determinism   two builds, recursive checksum: same/differ      (PR-07)
  end-to-end    page opened and read, URL and heading quoted     (FO-10)
  docs touched  no, or the operator's amendment cited            (CT-01)
  AGENTS.md     <n> bytes, cap 32768                             (FO-02)
  make check    <before> -> <after>
```

## Close-out

In closing message. Blank or `no` **stop the close**.

```
close-out wave <N>
  lanes              <n> of <m> merged, rest named in body
  gate on merged     green, run once by integrator               (FO-07)
  each lane's tree   read, not its report                        (FO-11)
  end-to-end         performed by person, what was seen          (FO-10)
  worktrees          <n> removed, <n> branches deleted           (SH-06)
  findings           <n> raised, <n> in docs/FINDINGS.md
  decisions          <n> rows in docs/DECISIONS.md with scope    (RV-02)
  reviewer           agreed / dissented / not consulted, by name
  next wave named    no
```

- **CL-01** **Done-when is list cited from task, and every line met or deferred by name.** Answered from memory is not answered.
- **CL-02** **Finding outside scope go in `docs/FINDINGS.md`, in closing commit.** Row carry where found, what break, what it look like when it break.
- **CL-03** **Claim later action erase carry its date and that action.**
- **CL-04** **Turn end `Your move:` / `Mine:`. Box go last.**

## What this enforces

**Fifty-one of sixty-six clause bind an artefact. Fifteen bind nothing.** Nothing
in the formatting tell them apart, and `skill check` reciting sixty-six read as
coverage. This section exist so reader can tell which is which - and reader who
most need it is this session, in three week, under load.

**Bound, and by what.** `RD-01` and `RD-05` by *requirements read whole with line
ranges* in pre-flight. `PR-01` by *baseline: every check proved red before counted
green*. `PR-06` by *breaches n planted, n red, n restored*. `PR-07` by determinism
row. `FO-02` by byte count. `FO-10` by *end-to-end, page opened*. `FO-11` by *each
lane's tree read, not its report*. `RV-02` by *decisions n rows with scope*.
`CT-04` by path check in `make check`. `FO-12` by `tools/lane_overlap.py`.

**Unbound, named rather than left to look enforced:**

| Clause | Why nothing catch it | What would |
|---|---|---|
| `ST-03` | Turn text unread by any tool | Filler list is enumerable - grep it |
| `ST-04`, `ST-05` | Bind turn text | Nothing cheap |
| `ST-06` | `Makefile:24` match `AGENTS*.md` and `.agents.md`, never `SKILL.md` | One `-o -name 'SKILL.md'` |
| `PR-02` | *Prefer* is tell. No input turn it red | Heuristic, keep as one |
| `PR-09` | First half is stance | Second half - *print value, never `ok`* - is grep |
| `SH-02` | `.git/hooks` empty | Pre-commit hook, four line. Cheapest unconverted rule here |
| `SH-03` | Nothing read worktree path | `git worktree list` against repo path, two line |
| `SH-07` | Nothing read commit shape | `git log` arithmetic |
| `AS-01`..`AS-08` | Every one is turn-shape obligation | Nothing. Family is prediction, never record |

**Three of nine W1 item ship unguarded** - byte emitter, `<head>` fragment,
`AGENTS.md` assembly. Hand-planted at close by this session. Hand-planting is not
gate, and saying so here is what stop green close reading as covered.

**`ST-13` and this section arrive together, deliberately.** File carrying *rule
phrased as advice fire while auditing* plus fifteen advice clause break own rule
on same page. `ST-13` alone read as satisfied and make file less honest.

## Step checklist

**Reference read at step that use it. Never batched at commit.**

1. Read `AGENTS.md` whole, plus every requirement in scope, by range (`RD-01`)
2. Probe every machine fact this turn. Document is not fact (`RD-02`, `RD-03`)
3. Read `contract.md` before touching requirement. ADR before code (`CT-01`)
4. Read `fanout.md` before cutting lane or writing task (`FO-03`, `FO-12`)
5. Write task: what is true, what to build, what done mean (`FO-06`)
6. Read `proving.md` before writing check. Prove red at baseline (`PR-01`)
7. Worker return: read tree, plant every Done-when, watch red (`FO-11`, `PR-06`)
8. Merge, gate merged tree once, remove worktree (`FO-07`, `FO-08`, `SH-06`)
9. Read `closing.md`. Close-out block in full
10. End `Your move:` / `Mine:`, box last (`CL-04`)
