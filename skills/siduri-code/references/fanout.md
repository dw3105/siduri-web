# fanout - cutting lane, writing task, merging wave

**v1 - 25 Aug 2026.** Read before cutting a lane, writing a lane task, or merging
a wave. `FO` rules live in `SKILL.md` and stand alone. This file carry shape.

---

## The seam decide the lane cut, never the feature list

Lane cut by feature look clean and merge badly. What decide a cut is which file
each lane must write.

Measured, 25 Aug 2026. Eight lanes drafted against `internal/site` after the
spine landed. Drafting them surfaced four shared line nobody would have seen from
inside one lane:

- `PageData` is one struct. Tags, feeds, tools, comments each want a field.
- `templates.templ` hold fourteen template in one file. Every lane add one.
- `<head>` take feed auto-discovery, JSON-LD and font preload from three lane.
- `Makefile` take target from any lane needing one.

None of four is visible while writing a single lane. All four land on integrator.
**Writing every lane task before cutting any worktree is what surface them**, and
it cost an hour against a wave.

## Registration lifecycle beat output cardinality

First diagnosis was `Route.Output` being one string, so no lane express a
collection. True and shallow.

Real defect: every `Register` call sit in `init()`, and content load inside
`Build()`. Lane needing one page per tag cannot register at all - at `init()`
time tag do not exist. Four of eight lane need per-content route.

Second floor under it: route set is package-global and never cleared, so
registering during build panic on second build in same process with `duplicate
route`. Nothing in tree call build twice, so no check see it.

**Test:** ask whether lane can name its output before content load. Cannot →
lifecycle, not cardinality.

## Proof of disjointness pass on broken seam if written from inside

First disjointness proof: *add a lane file registering a route, touch no other
file, report the diff*. Reviewer wrote it against the un-widened seam. It worked,
because worked example already committed did exactly that.

Two reason it passed:

- Demo carry own data by closing over a literal. Real lane data is
  content-derived and exist only inside build. Demo dodge shared struct
  **because** it is demo.
- Command was `git diff --name-only`. New file is untracked, so it report zero -
  and it **would** report the generated shared file if regeneration touched it.
  Blind to what it count, sighted on what it never asked.

**Replacement shape:** check that delete workaround. `grep -c ArticlePage
build.go` returning `0` is red while loop exist and green only when mechanism
carry the case that forced the loop. Verified red at `1` before it shipped.

## Worker cannot be bound by anything it does not read

Codex worker under `-s danger-full-access` read its task file. It does not read
this skill, and it does not read `CLAUDE.md` unless that is `AGENTS.md`.

`AGENTS.md` is discovered by walking up from working directory, so it load
regardless of `trust_level` - and `trust_level` is keyed by absolute path, so a
worktree is a different project from its parent whatever the parent say.

`project_doc_max_bytes = 32768`, read out of the binary, unset in config.
Assembly is root-first leaf-last and **stop at cap**, dropping end. House form put
enforcement at end. So a skill shipped as `AGENTS.md` delete its own pre-flight
and close-out, silently, in every lane.

Per-file byte check is not the mechanism. Cap apply to combined assembly. One
file today make the two equivalent; a nested `AGENTS.md` make every file pass
alone while the leaf silently disappear.

## Worker editing the file that govern worker

Spine lane edited `AGENTS.md`, legitimately - it filled content area the contract
require. Verified by planting, not by reading, and both claim held.

Shape is still wrong: mechanism that bind worker is mechanism worker can rewrite.
Filling content area is the job. Editing prohibition is removing own leash. Path
check cannot separate them while they are one file.

**Fix is a path, never a rule:** prohibition live in own file, assembled into
`AGENTS.md` at build. Then the un-editable half is a path the existing guard
already handle.
