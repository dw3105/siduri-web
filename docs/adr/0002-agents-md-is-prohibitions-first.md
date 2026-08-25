# 0002 · AGENTS.md is prohibitions-first

**Status** · **proposed** — amends the stated order in `REQUIREMENTS-SITE.md:336`,
so it needs the operator's word · 2026-08-25

`AR-7` fixes the content order of `AGENTS.md`: architecture, the command list,
content conventions, the tag vocabulary, the tone guide, and *the prohibitions
below* — `AR-8`, last.

Codex assembles project docs up to `project_doc_max_bytes` and drops the end.
Measured on the binary that runs here, not from documentation:

    strings .../codex-linux-x64/vendor/.../bin/codex | grep -o 'project_doc_max_bytes = [0-9]*'
    project_doc_max_bytes = 32768

The key is unset in `~/.codex/config.toml`, so 32768 is in force. Truncation is
silent — no warning in the TUI, in `exec`, or in `/stats`.

So `AR-7`'s order puts `AR-8` exactly where the cap deletes it. The failure is
precise: acceptance line `523` — *`AGENTS.md` lets a fresh agent session publish
a post without asking a question* — reads **green**, because the command list is
near the top. The prohibitions are absent from every worker's context, because
they are at the bottom. In a public repository (`D-5`) that is a portfolio piece
for a human-in-the-loop brand, the file the contract calls the agent's contract
would be missing precisely its human gates.

`AR-8` carries the `D-11` gates: no post publishes without approval, agents draft
and never publish, no agent touches legal pages, dependencies, paid services or
secrets without per-instance approval.

**Decision proposed.** Prohibitions first. Everything else follows. The file
stays small and structural; the working rule lives in `skills/siduri-code/`,
which binds a Claude session and is under no byte cap.

**Also proposed**, both cheap and neither dependent on the above:

- Set `project_doc_max_bytes` explicitly rather than inherit a default that can
  move on the next `npm -g` upgrade — this one already differs from what the
  current docs claim.
- Put the byte count in `make check`, red on breach. *Keep it small* is advice
  and fires while auditing; a number fires while designing. Acceptance `523` as
  written has no red state; a byte check does.

**Residual, recorded rather than solved.** `CLAUDE.md` is a symlink to
`AGENTS.md`, so one file has two consumers with different limits. Claude Code
warns on an oversized `CLAUDE.md`; Codex does not. Whoever tests it in Claude
Code sees it work.
