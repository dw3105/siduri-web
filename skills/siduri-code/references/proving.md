# proving - check, breach, and what count as caught

**v1 - 25 Aug 2026.** Read before writing any check, planting a breach, or
calling something caught. `PR` rules live in `SKILL.md`.

---

## Establish red before counting green

Every check in this repo is planted before it count. Not because the check look
weak - because the two that looked strongest were the two that were green at
baseline.

Measured 25 Aug 2026, spine lane. Five breach planted, each went red, each
restored:

    docs/FINDINGS.md touched        → check: docs/ is contract-owned and must not change
    AGENTS.md grown to 43448 byte   → check: ./AGENTS.md is at or above 32768 bytes
    post with no plain_summary      → build: missing required frontmatter field
    post tagged nonsense            → build: unknown tag "nonsense"
    draft true → false → true       → slug in 2 file, then 0

That is the standard. Report which input fired, never that the check exist.

## Check that delete the workaround beat check that stand beside it

Demo route proving a property stand beside the defect and can pass while the
defect stand. Check asserting the workaround is **gone** cannot.

`grep -c ArticlePage internal/site/build.go` → `0`. Red at `1` today. The loop
removing itself is the proof.

## `git diff` see no untracked file

Proof counting new file must use `git status --porcelain`. `git diff --name-only`
report zero for untracked file and report tracked generated file it was never
asked about.

Count in a Done-when is ownership, never a number. Templ lane is three file -
`.templ`, `.go`, and the tracked generated `_templ.go`. Go-only lane is two.

## Defect invisible to one run need a run that repeat

Package-global state survive a build. One build see nothing; second build panic.
Gate that build once cannot see it at any effort level.

Same shape: orchestration with no test at all. Three test in the spine package
cover component render and content loading. **Zero invoke the build function**,
so content filtering, draft exclusion, sort, route iteration, output path and
write are all untested and the gate stay green through a broken build.

Rewriting an untested orchestration is the shape to name before it start.

## Document is not machine fact

Friction row three day old named a sysctl as cause of four dead worker. Probed:
sysctl read `0`, and the shape that killed them exit clean. Row promoted to
current machine fact, in a plan, by the same session that had spent the day
holding others to probe-this-turn.

Cause that no longer reproduce make everything near it suspect at the same
moment. Strike the cause, keep the row, record why both.

## Determinism

Two build, recursive checksum over the tree, compared. Never by eye, never by
file count.
