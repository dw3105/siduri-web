---
name: siduri-code
description: Working rule for the Siduri site build, repo siduri-web: two specs are the contract, eight Codex lanes run parallel and read no skill, the merged tree is the only verdict. Read before acting.
---

**v1.4 - 27 Aug 2026.** Derived from `gateslot-code` v1.13. Structure kept,
content earned in one day: repo genesis, one spine lane, one seam lane, nine
findings, four of them against this session's own claims. File caveman full.
Output stay **lite**. Change log live in `ledger.md`, never here.

**Canonical copy is `skills/siduri-code/` in `siduri-web`.**
`~/.claude/skills/siduri-code` is deployed copy. **No install target and no drift
check exist yet**, and the two copies **differ right now** - this file is ahead of
deployed by `SS-01`..`SS-07`, `SS-04a` and `PR-11a`, so clause you are reading here
may not be clause the running session loaded. Earlier version of this paragraph said
`deployed == canonical` is true today and held by nothing; that sentence was false
for as long as `PR-11a` existed in deployed and nowhere else, and stayed in the file
through the commit that repaired the drift it described. `RD-02`: probe beat
document, including this one. `hitl` bind it with `make skill-check`. Findings 0013,
0085.

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
siduri-code v1.4: two spec files are the contract and only operator amend them | Codex worker read task file and never this skill, so clause reaching worker live in task text or in guard | AGENTS.md truncate at 32768 silently and drop end, so prohibition go first and live in file worker cannot edit | lane own only new file, and file every lane need belong to integrator | registration happen when data exist, never before | worker run no gate, integrator run it once on merged tree | branch green alone say nothing about branch merged | check green at baseline certify nothing, prove it red first | agent report is claim, read tree | document is not machine fact, probe this turn | rename inside contract dangle reference the contract cannot fix | reviewer agreement is precondition on proposal, never substitute for approval | agreement carry scope, what was checked and what was not | one decision per turn, box last, Your move and Mine | route operator already use is fact like credential he hold | task naming parsed output pin every field | correctness per-part mean defect live between part, no per-part gate see it | field carrying command get re-run or carry no command | guard pinned to value it guard go red on correct work | directory doing exempted job need its exemption | gate running only in mode where it cannot fail is not gate | binding is not serving, check gate scan what it think it scan | under green-required gate, unit of work is whatever make gate green | wave not closeable until operator looked at artefact | operator ruling get its row in the turn it is given | clause written into deployed skill and not into skills/ is lost, and install target overwrite it | never end turn with uncommitted work, and destructive probe go in throwaway clone | ask what is running never what is installed | probe block verdict equal every row it name and blocks partition once, row addressed by section and number | finding withdrawn after checking get row carrying probe that withdrew it | section carry no count, because number here went stale four time while carrying its own command
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
- **ST-13** **Rule phrased as advice fire while auditing, never while designing.** Write what diff can be checked against, never what to be careful about. Read `## What this enforces` before adding clause: unbound table there already list one class of clause that fail this, and several were added by session that quoted this clause while adding them.
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
- **RD-06** **Environment table and the route operator already use are fact about machine, never advice.** Re-read at point of writing box that touch credential **or transfer**. `gh auth` was proposed for machine whose table list it under *Never holds*. Later same week, remote `claude-vm` was written in this session own memory and three replacement transfer proposed anyway - `scp`, bundle, attachment - before operator said *we did this before*.
- **RD-07** **Before handing operator command, ask what it would take to run it here.** *Cannot verify here* is almost always *did not look*. Node 20 and Chrome installed in two command after operator insisted; three defect found in ten minute after, none diagnosable from CI failure text.

## Session shape

Every clause here come from loss measured in this repo on 28 Aug 2026, and two
of them were broken by the session writing them, same hour.

- **SS-01** **Session is three phase: close carry-in, do work, deposit record.** Third phase is where every loss measured here happened. It is not closing paragraph, and it is not optional because turn ran long.
- **SS-02** **Operator ruling get its `docs/DECISIONS.md` row in turn it is given, never at close.** Four ruling in one day, zero row, recovered only because one transcript survived one `/compact`. Row is cheap in the turn and unwritable a day later.
- **SS-03** **Finding row and decision row commit unasked. Only rule change get asked.** Asking permission to record is how recording become optional.
- **SS-04** **Never end turn with uncommitted work in tree.** `git status --short` is one call. Backfill of ten finding row and seven decision row sat uncommitted six tool call and died to `git reset --hard` run beside it.
- **SS-04a** **Destructive probe run in throwaway clone, never in worktree holding work.** `git reset --hard`, `checkout -f`, `clean -fd` take whole tree, never just what probe touched. Reviewer ran identical probe in clone same morning and lost nothing.
- **SS-05** **Clause written into deployed skill and not into `skills/` is lost, and nothing report it.** `PR-11a` lived only in `~/.claude/skills/` while every session loaded it. Skill edit land in `skills/` first, deploy from there, never other direction. **Install target copy canonical over deployed, so building it before re-landing lost clause perform the loss it exist to stop.**
- **SS-06** **Number in record carry command that produced it.** Close line reading *3 clean* over table showing two survived a wave because nothing recomputed it.
- **SS-07** **Never assert where session sit in own context.** Harness summarise when conversation grow long. Guess wearing confidence of measurement justify stopping, and operator cannot check it.

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
- **FO-10a** **Wave not closeable until operator has looked at artefact.** Eleven lane merged, CI green, publish path designed and credential requested before anybody opened one page. Deliverable built here and never delivered is `OP-35` failing silently, and it cost a day of unreviewed design.
- **FO-11** **Worker report is claim. Read tree.** Not summary, not test count, not *gate passed*.
- **FO-12** **Filename named by two lane task is seam defect, and it is grep.** `tools/lane_overlap.py` over task set before any worktree cut. Non-empty output mean wave not ready. Fix seam, never route lane around it.
- **FO-13** **Integrator-owned file get threshold of one, never two.** `go.mod`, `go.sum`, `Makefile`, `routes.go`, `build.go`, `content.go`, `markdown.go`. Task naming one already claim file lane do not own. Suppressing them read cuttable on wave that collide four way.
- **FO-14** **Task text is declared intent. Branch is behaviour.** Lane writing file it never named is invisible to any regex. `lane_overlap.py --branches` before merge, pairwise over `git diff --name-only`. First prevent, second detect, neither substitute.
- **FO-15** **Every dependency resolved into `go.mod` before any branch cut. No lane run `go get`.** Eight lane adding library is real conflict, never accidental, and `go.sum` land on integrator. Dependency claimed on unmerged branch is not claimed.
- **FO-16** **Under gate requiring green, unit of work is whatever make gate green.** Splitting free only while every piece independently sufficient. Three lane each fixing part of one failure can none of them merge, and stacking not help because first in any order still red.
- **FO-18** **Task naming output another program parse pin every field that program read.** Format given as prose is format each lane invent. Six lane produced three convention and harness invented fifth; gate refused four time on first run. Lane was not wrong, spec was.
- **FO-19** **Where correctness is per-part, defect live in relation between part, and no gate examining one part see it.** Operator found two site chrome, three post-meta, two content width in five minute; each copy internally consistent and green alone. `hitl` met same shape three time same day without rendering involved - guard inert exactly while it mattered, probe attached to claim it did not measure, one link parser silently dropping anchor shape and blinding seven guard each individually correct. Rendering is one instance, never the class. This is *why* `FO-10a`, not restatement of it.
- **FO-17** **Merged tree say nothing about either branch alone.** `CX-22` converse. Gating C1 and C3 together then predicting each PR gave two wrong prediction, opposite direction.

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
- **PR-10** **Gate running only in mode where it cannot fail is not gate.** External link check ran a day with `A8_EXTERNAL_LINKS=0`, green every time, 66 failure waiting. Ask which mode gate ran in, never whether it ran.
- **PR-11** **Binding is not serving. Check gate scan what it think it scan.** Port 8765 held by unrelated service; `http.server` failed to bind, gate did not notice, axe scanned that application and reported 122 violation against markup this repo do not contain.
- **PR-13** **Field carrying command line get re-run like `ran`, or carry prose and no command.** De-quoting address to pass secret scanner turned runnable `grep` into one matching nothing, while sentence around it still claim it showed both location. Survived because harness validate field presence, never `notes`.
- **PR-14** **Guard pinned to value it guard go red on correct work. Derive target twice by different mean, compare, never name it.** Tool built to remove pinned line range pinned count `18` in own selftest - same coupling one layer up.
- **PR-15** **Directory doing job of exempted directory need its exemption, or it is record that cannot record.** `acceptance/` sit outside `docs/` so no lane edit contract; secret scanner exempt by prefix `docs/`, so fragment quoting the address it found turned `make check` red on own evidence. Scope follow role, never path.
- **PR-11a** **Ask what is running, never what is installed, and freeze cover server not only byte.** `hitl` installer report new tag from `--version` while `/proc/<pid>/exe` resolve to install path marked `(deleted)` and running image digest differ from installed file - every green true, composite false. Port answering 200 say nothing about which program answer. This session swapped the server under operator's own acceptance pass and recorded only the injected tag.
- **PR-16** **Report grouping row under one probe get two property checked, never one: block verdict equal table verdict of every row it name, and block partition row set exactly once.** Grouping several annotation of one defect under one probe is right, and it is also where verdict hide. Both property are mechanical and found real contradiction by hand in a minute - one annotation carrying three verdict across two report, each defensible as of different date, none saying which was current.
- **PR-16a** **Row address is (section, number), never number.** Report with table per page repeat `row 1` per page - three time in one file here. Rule implemented on bare integer report mismatch that is artefact of own flattening, which is `PR-12`: check on check need same standard as check. Measured while running `PR-16` for first time.
- **PR-12** **Check on check need same standard as check.** Grep for `AKIA` where tool print `AWS access key` report working scan broken. Plant landing after frontmatter delimiter report working guard absent. Name input, read actual output, confirm plant landed where aimed.

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
- **SH-01a** **Nothing merge locally. `main` advance only through pull request, no bypass.** Lane commit on own branch, operator push branch, CI `check` gate merge. Lead's `main` track `origin/main` and stop being where merge happen.
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
- **SH-12** **Publish path is target with refusal, never block pasted into terminal.** Bootstrap have no target so one block unavoidable; every block after that is debt with shell prompt.

## Ask

- **AS-01** **One decision per turn. Size is not unit.**
- **AS-02** **Blocker asked at once, never batched. Unblocked work continue meanwhile.**
- **AS-03** **Question carry its own weakness, and rejected option shown ruled out.**
- **AS-04** **Recommendation mandatory, with its reason.**
- **AS-05** **Operator hold no context.** Turn carry every fact needed, in that message.
- **AS-06** **Frame what break, never what changed.**
- **AS-07** **Decision takeable without opening laptop.**
- **AS-08** **Question that contradict operator's own instruction go to him as contradiction, named.** Contract he wrote can forbid order he asked for.
- **AS-09** **Asking operator to look at rendered thing mean giving exact URL on host they are using, never path and never page name.** They read turn on phone, in another session, hours later, and reconstructing base URL is work they should not do. Port is theirs, not yours: tunnel that forward local 18080 make every 8080 link wrong. Ask once what base they hold, then paste full link per page, one per line, clickable.

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
- **CL-05** **Finding withdrawn after checking get row, and row carry probe that withdrew it.** Row exist or it do not, so this is checkable; *act of looking is evidence* is not, and would join unbound table rather than bind anything. Withdrawn finding leaving no trace get re-opened by next reader with no record anybody looked. Requires `docs/FINDINGS.md`'s negative-result row shape, added same day.
- **CL-03** **Claim later action erase carry its date and that action.**
- **CL-04** **Turn end `Your move:` / `Mine:`. Box go last.**

## What this enforces

**No count appear in this section, deliberately.** Number here went stale four
time - *fifteen*, *seventy-two*, *ninety-nine*, then stale again on very next
commit that added clause. Fourth one broke inside commit about not losing
record, with instruction *re-run it rather than trust it* one line away and its
command printed beside it.

**That is `SS-06` failing on its own terms.** Clause say number carry command
that produced it. This number carried its command, verbatim, one line away, and
rotted anyway - because carrying command is not running it, and instruction to
re-derive fire while auditing, never while adding. `ST-13`, one level up from
where it usually apply. Clause was satisfied and defect recurred, four time.

So numeral is gone rather than corrected fifth time. Count when you need it:

    grep -oE '\*\*[A-Z]{2}-[0-9]+[a-z]?\*\*' SKILL.md | sort -u | wc -l

Two caveat for whoever run it. Suffixed id - `PR-11a`, `SS-04a`, `PR-16a` - need
`[a-z]?` or they vanish, which is how *seventy-two* survived. And `CT-01` and
`RD-01` are bold paragraph rather than list item, so `^- `-anchored command give
two fewer and read as two clause lost.

**Property that replace the number, and it is greppable:** every clause added is
named in bound list or unbound table below, or it sit in unaudited remainder by
not being named at all. Test is one grep for the id. `CL-05` was unclassified
for one commit and no number would have caught that - `PR-14`, derive target,
never name it, demonstrated four time.

Nothing in formatting tell bound from unbound, and `skill check` reciting every
clause read as coverage. This section exist so reader can tell which is which - and reader who
most need it is this session, in three week, under load.

**Bound-pending, and by what.** `SS-04` by `git status --short`, one call, run before turn end. `SS-05` and `SS-06` bind once `skill-check` and `close-check` exist and are unbound until then - listed here as **pending**, never as bound, because a clause bound by a target nobody has built is the exact shape this section exist to expose. `PR-16` and `PR-16a` bind by partition-and-uniformity over report row, same condition.

**Bound today, and by what.** `RD-01` and `RD-05` by *requirements read whole with line
ranges* in pre-flight. `PR-01` by *baseline: every check proved red before counted
green*. `PR-06` by *breaches n planted, n red, n restored*. `PR-07` by determinism
row. `FO-02` by byte count. `FO-10` by *end-to-end, page opened*. `FO-11` by *each
lane's tree read, not its report*. `RV-02` by *decisions n rows with scope*.
`CT-04` by path check in `make check`. `FO-12`, `FO-13` and `FO-14` by `tools/lane_overlap.py`, ten selftest case.

**Unbound, named rather than left to look enforced:**

| Clause | Why nothing catch it | What would |
|---|---|---|
| `SS-01` | Phase structure. No diff show a session skipped a phase | Nothing cheap |
| `SS-02` | *In the turn it is given* is turn-scoped, and turn text unread by any tool | `close-check` refusing a wave whose operator ruling count exceed its `DECISIONS.md` row count |
| `SS-03` | *Commit unasked* bind intent, never artefact | Same guard as `SS-02` |
| `SS-04a` | *Throwaway clone* is where command ran, not what it produced | Nothing cheap |
| `SS-07` | Bind turn text | Nothing cheap |
| `CL-05` | Row exist or it do not, so it is checkable - but nothing walk `docs/FINDINGS.md` looking for withdrawn finding, and nothing know a finding was withdrawn | `close-check` extended to `FINDINGS.md`: refuse row whose *What breaks* is `nothing` and whose *What it looks like* carry no probe |
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
phrased as advice fire while auditing* plus every advice clause in unbound table break own rule
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
