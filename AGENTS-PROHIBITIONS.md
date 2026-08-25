# AGENTS.md — the agent contract for siduri-web

Required by `docs/site-requirements.md:336` (AR-7). `CLAUDE.md` is a symlink to
this file.

**Prohibitions are first on purpose.** `AR-7` lists them last. Codex assembles
this file up to `project_doc_max_bytes` — 32768 on codex-cli 0.149.1, measured in
the binary — and drops the **end** with no warning. Enforcement at the tail is
enforcement that silently disappears. The deviation from `AR-7`'s stated order
is recorded in `docs/adr/0002`, and is proposed until the operator rules on it.

## Prohibitions — AR-8

Agents **must not**, without explicit per-instance human approval:

- Publish a post. Agents draft; a human publishes.
- Publish a comment. Approval is per-comment; "approve all" is refused.
- Modify `/impressum/` or `/datenschutz/`.
- Write anything to `content/` containing a raw email address.
- Add a runtime dependency, a paid service, or a third-party script.
- Force-push or rewrite history.

Two more, specific to this repo:

- **Never edit anything under `docs/`.** The two requirements files are the
  contract and only the operator amends them. A worker that finds a requirement
  wrong stops and says so. This overrides the general "fix the document you find
  wrong" habit.
- **Never merge, never touch `main`, never push.** One commit on your own branch.
