# Findings

Out-of-scope findings. Each row: where found, what breaks, what it looks like
when it breaks.

| # | Where | What breaks | What it looks like |
|---|---|---|---|
| 0001 | codex-cli 0.149.1 binary | `AGENTS.md` truncates at `project_doc_max_bytes = 32768`, dropping the end | Nothing. No warning in TUI, `exec`, or `/stats`. `gateslot-code/SKILL.md` is 32,394 bytes — 374 of headroom. Root-first, leaf-last assembly stops at the cap, so a fat root file deletes per-lane leaves entirely and "root is under the limit" reads green. |
| 0002 | `~/.codex/config.toml` | `trust_level` is keyed by absolute path, so a worktree is a different project from its parent | Anything under a worktree's `.codex/` is skipped silently. Nothing repo-side ever bound a Codex worker in this house. |
| 0003 | `git worktree add` | Path resolves relative to `-C`, so a bare name nests the worktree **inside** the repo | `?? wt-lane2/` untracked, not ignored. Gates walk it; `git add -A` stages a whole worktree. `../name` escapes and `name` nests — both relative, opposite outcomes. Absolute paths only. |
| 0004 | this box, 25 Aug 2026 | `gateslot` FRC-0009 attributes four dead Codex agents to `kernel.apparmor_restrict_unprivileged_userns=1` | The sysctl reads `0` here and `bwrap --unshare-net --dev-bind / /` exits 0. Either it moved, or the row was misattributed, or it was measured on another box. A cause that no longer reproduces makes everything near it suspect. |
| 0005 | this box | No `~/.gitconfig`, `init.defaultBranch` unset, fresh `git init` yields `master` | `fatal: invalid reference: main`, after the worktrees are cut. |
