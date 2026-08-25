# AGENTS.md — the agent contract for siduri-web

Required by `docs/REQUIREMENTS-SITE.md:336` (AR-7). `CLAUDE.md` is a symlink to
this file.

**Prohibitions are first on purpose.** `AR-7` lists them last. Codex assembles
this file up to `project_doc_max_bytes` — 32768 on codex-cli 0.149.1, measured in
the binary — and drops the **end** with no warning. Enforcement at the tail is
enforcement that silently disappears. The deviation from `AR-7`'s stated order is
recorded in `docs/adr/0002`, and is proposed until the operator rules on it.

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

## Architecture and commands — AR-1, AR-4

Markdown in `content/posts/` and the closed vocabulary in `content/tags.yml`
are loaded by `internal/site`. The loader validates required frontmatter and
tags, the deterministic renderer turns it into static HTML in `dist/`, and
`internal/site`'s package-level `Register(...)` seam lets later lanes add routes
from their own `init()` files. Typed templates live in `.templ` files and their
generated `_templ.go` files are committed. `cmd/siduri` is only the build and
preview entry point; it does not own route registration.

Routine operations are one command each:

    make dev · make build · make check · make publish · make comments · make deploy · make rollback

## Where things are

    docs/REQUIREMENTS-SITE.md       the contract — site
    docs/REQUIREMENTS-COMMENTS.md   the contract — comments
    docs/ARCHITECTURE.md            package boundaries
    docs/IMPLEMENTATION-PLAN.md     lanes and their scope
    docs/DECISIONS.md               decision record
    docs/FINDINGS.md                out-of-scope findings
    docs/adr/                       one file per mechanism decision

## Content conventions and voice

Posts use the required frontmatter fields `title`, `slug`, `date`, `summary`,
`plain_summary`, `tags`, and `draft`. `plain_summary` is a one-line explanation
for a non-developer. Drafts stay `draft: true` until a person approves
publication. Tags are closed and must be one of: `build-log`, `tool-release`,
`dogfooding`, `method`, or `outcome`.

Write directly, technically, specifically, and in the first person. Use real
tools, durations, costs, and failures. Avoid growth-hack language, fabricated
success, and sentences that could have been written without doing the work.

Wave 0 binds the two operational claims in this file: `make check` rejects any
change under `docs/`, and it rejects every `AGENTS*.md` or `.agents.md` file at
or above 32768 bytes. The cap is the codex-cli limit recorded in
`docs/FINDINGS.md` row 0001; a clause deposited with no guard behind it is a
clause that reads enforced and is not.
