# W1 · G2 — content guards

wave: W1
tree: 7667b8d5c067d047bffa402cdb58fc92550a3ccc

## Criterion 6 — Every post has a `plain_summary`; the build fails on a post that doesn't.

verdict: held
ran: `make build`
at: 2026-08-27T18:18:50Z
saw: `content/posts/g2-missing-plain-summary.md: missing required frontmatter field "plain_summary"`; `make: *** [Makefile:22: build] Error 1`; exit 2.
red proof: The temporary `content/posts/g2-missing-plain-summary.md` omitted `plain_summary`; the publish build rejected that input and exited 2.
notes: The fixture was removed. The loader's required-frontmatter validation is the check that binds this criterion.

## Criterion 7 — An unknown tag fails the build.

verdict: held
ran: `make build`
at: 2026-08-27T18:19:02Z
saw: `content/posts/g2-unknown-tag.md: unknown tag "g2-not-a-real-tag"`; `make: *** [Makefile:22: build] Error 1`; exit 2.
red proof: The temporary post carried `g2-not-a-real-tag`, which is absent from `content/tags.yml`; the publish build rejected that input and exited 2.
notes: The fixture was removed. The closed-vocabulary tag validation is the check that binds this criterion.

## Criterion 8 — A draft post appears nowhere in `dist/`, verified by grep.

verdict: open
ran: `make build`; `grep -rIn 'g2-draft-page-absent-20260827' dist`; `make check` with the draft fixture present
at: 2026-08-27T18:19:19Z
saw: (build 18:19:19Z, grep 18:19:25Z, make check 18:19:41Z) The temporary `g2-draft-page-absent-20260827` post built successfully; the direct `grep` returned no output and exit 1; `make check` nevertheless exited 0 with the draft still planted. `dist/journal/` contained only `hello-siduri/index.html`.
red proof: No honest red run exists for the named page-tree guard: `rg -n -i 'draft.*dist/|dist/.*draft|grep.*dist/|dist/.*grep' Makefile mk tools internal --glob '!*.sum'` at 2026-08-27T18:22:55Z found only the unrelated pre-P3 `dist/` guard. The build's draft filtering held, but the criterion's automated grep check does not exist.
notes: The manual output condition is held for this distinctive draft, but the criterion is open as written because no repository check verifies the page tree and no owning step or phase is named. Existing Go tests cover feed and metadata draft exclusion, not the generated page tree. The fixture was removed.

## Criterion 11 — `grep -rE '[a-z0-9._%+-]+@[a-z0-9.-]+' content/` returns only intended addresses.

verdict: held
ran: `grep -rE '[a-z0-9._%+-]+@[a-z0-9.-]+' content/`
at: 2026-08-27T18:20:49Z
saw: No output; exit 1, which is the non-match result accepted by the Makefile's `if` guard.
red proof: A truthful red run cannot be planted under `content/` because AR-8 prohibits writing a raw email address there without per-instance human approval. The current content produced no matches.
notes: This is held only for the guard's `content/` scope. A reader-facing address is rendered from `internal/site/contact.templ` into `dist/contact/index.html`: `grep -RIn 'the `hello@` mailbox on the siduri.ai domain' dist/ internal/site/contact.templ` at 2026-08-27T18:20:58Z showed both locations. The Makefile guard does not inspect that template or `dist/`.

## Final verification

- **ran** — `make check`
- **at** — 2026-08-27T18:21:42Z
- **saw** — exit 0; 31 HTML pages validated by the accessibility gate, 0 linkcheck failures in the published build, secretscan passed, and all Go tests passed.

## Restoration

- **ran** — `git status --porcelain`
- **at** — 2026-08-27T18:23:09Z
- **saw** — (no output)
