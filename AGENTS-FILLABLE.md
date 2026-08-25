## Architecture and commands — AR-1, AR-4

Markdown in `content/posts/` and the closed vocabulary in `content/tags.yml`
are loaded by `internal/site`. The loader validates required frontmatter and
tags, the deterministic renderer turns it into static HTML in `dist/`, and
`internal/site`'s package-level `Register(...)` seam lets later lanes add routes
from their own `init()` files. Typed templates live in `.templ` files and their
generated `_templ.go` files are committed. `cmd/siduri` is only the build and
preview entry point; it does not own route registration.

Lane ownership is strict: a lane may add its own `.go`, its own `.templ`, its
generated `_templ.go`, and its own `mk/*.mk` file, and registers from its own
`init()`. A lane may not edit `routes.go`, `build.go`, the shared document
template, or another lane's file.

Routine operations are one command each:

    make dev · make build · make check · make publish · make comments · make deploy · make rollback

## Where things are

    docs/site-requirements.md       the contract — site
    docs/comments-requirements.md   the contract — comments
    docs/DECISIONS.md               decision record
    docs/FINDINGS.md                out-of-scope findings
    docs/adr/                       one file per mechanism decision

The seam a lane registers against lives in `internal/site/routes.go` and
`internal/site/article_sections.go`: `Register` for a fixed page,
`RegisterContent` for routes computed from loaded content, `RegisterHead` for a
`<head>` fragment, `RegisterArticleSection` for an additive article component
(`ArticleSectionBodyAside`, `ArticleSectionAfterBody`, or
`ArticleSectionFooter`), and `RegisterBodyRenderer` for the one owner allowed
to replace Markdown-to-HTML rendering. `PageOutput` / `ByteOutput` are what an
expansion produces. Article section render functions receive the post and
build's `PageData`; section names are sorted within their slot.

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
