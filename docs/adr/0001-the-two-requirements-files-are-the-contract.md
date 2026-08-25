# 0001 · The two requirements files are the contract

**Status** · accepted · 2026-08-25

`docs/REQUIREMENTS-SITE.md` and `docs/REQUIREMENTS-COMMENTS.md` are copied
verbatim from the operator's drafts of 2026-08-25. They outrank `AGENTS.md`,
this repo's skill, any plan, and anything said in a session.

They are amended by the operator's explicit word, given directly. Never relayed
through a peer session, never inferred from a review. A finding against a
requirement produces an ADR marked *proposed* plus a question — never an edit.

No Codex worker edits anything under `docs/`. Eight workers with full access in
worktrees, plus the ordinary habit of fixing a document found wrong, is one
plausible edit away from an amended contract nobody authorised. The prohibition
is in `AGENTS.md` and is guarded by a path check in `make check`, because a
clause the worker never reads binds nobody.

A requirement with no passing test is unimplemented, whatever the code does.
