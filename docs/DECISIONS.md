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

## Open

- **§12:497**, the P1→P2 gate — *Gateslot is usable by someone who is not the
  operator, and post 2 is published.* Not raised with him, not ruled on, still
  standing. Same class as the struck 490. Blocks wave B; blocks nothing before it.
- **§12:502**, the postal-address gate. Standing, and it should — it gates
  collecting personal data, not building the code that would.
- **ADR 0002**, `AGENTS.md` ordering. Deferred to after wave 0 by his ruling.


## Reviewer scope, as of 2026-08-25

Checked: environment claims, spec citations, `AGENTS.md` against the codex
binary, worktree and identity mechanics, both scope assumptions.

Not checked: the operator's "everything before commercialization" sentence
(never seen it), `REQUIREMENTS-COMMENTS.md` beyond two cited lines, the eight
lane definitions for `RV-23` disjointness — they do not exist yet.
