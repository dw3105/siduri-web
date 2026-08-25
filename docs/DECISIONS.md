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
| 2 | 2026-08-25 | Default branch is `main`, set explicitly | proposed | raised it — fresh `git init` here gives `master` | pending |
| 3 | 2026-08-25 | `AGENTS.md` is prohibitions-first, against AR-7's stated order | proposed, ADR 0002 | raised it; truncation eats the tail | **pending — his, it amends a requirement** |
| 4 | 2026-08-25 | Build order across P0/P1/P2 vs the §12 gates | proposed PHASE gate; reviewer refuted | dissents — PHASE guards `dist/`, the `490` gate guards the tree | **pending — open question to him** |

## Reviewer scope, as of 2026-08-25

Checked: environment claims, spec citations, `AGENTS.md` against the codex
binary, worktree and identity mechanics, both scope assumptions.

Not checked: the operator's "everything before commercialization" sentence
(never seen it), `REQUIREMENTS-COMMENTS.md` beyond two cited lines, the eight
lane definitions for `RV-23` disjointness — they do not exist yet.
