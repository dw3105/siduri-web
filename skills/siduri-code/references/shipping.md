# shipping - worktree, commit, command block

**v1 - 25 Aug 2026.** Read before cutting a commit or writing a command block.
`SH` rules live in `SKILL.md`.

---

## Before any worktree exist

No `~/.gitconfig` on this box. Fresh `git init` refuse the first commit with
`Author identity unknown`, and give default branch `master`, not `main`.

Worktree share the parent's `.git/config`. One `git config --local` on the parent
cover every lane - **and only if it land first.** Unset means eight worker die at
their first commit, in eight worktree, after doing all the work.

Identity here is the brand, not the person: the repo is public by decision, and
the contract ship the first two phase anonymously. Author line on every commit is
a published disclosure.

## Worktree path is absolute

    git -C repo worktree add -b lane2 wt-lane2 <base>   → lands at repo/wt-lane2

Path resolve relative to `-C`. Bare name **nest** inside the repo, where every
gate walk it and `git add -A` stage a whole worktree. `../name` **escape**. Both
relative, opposite outcome. Absolute only.

## Launching

Interactive, one terminal per lane:

    codex --cd <absolute worktree> --dangerously-bypass-approvals-and-sandbox

Then paste a pointer to the task file, never the task. Pasted task fill the
screen, and on single-line input the first newline submit while every line after
answer whatever dialog come next.

Unattended form work on this box - probed end to end, exit 0, shell command ran,
no sandbox failure:

    tools/run_codex.sh --deadline <s> --label <label> <wt> <task.md>  # codex exec via the runner

The host-wide Claude Code `PreToolUse` guard `~/skills/guards/bash_check.py`
(release `01a0bd5`) denies the bare form.

Eight interactive launch is a multi-turn sequence. Each turn restate which box of
eight and what the last one produced, or the operator lose their place at four.

## Commit

One lane, one commit, subject assigned by the lead and used verbatim. Two `-m`
flag, subject under 72 char, body max six line.

Worktree removed and branch deleted at merge, so `git worktree list` stay a true
statement of what is outstanding.
