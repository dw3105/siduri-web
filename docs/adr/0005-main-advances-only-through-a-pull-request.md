# 0005 · main advances only through a pull request

**Status** · accepted · 2026-08-25 · operator's word, given directly

`main` carries a ruleset with no bypass actors: deletion blocked, non-fast-forward
blocked, a pull request required, and the `check` status required to pass with
strict up-to-date enforcement. The ruleset lives at `.github/ruleset-main.json`
so the protection is reviewable in the same place as the thing it protects.

**No bypass, including the owner.** A bypass held by the only person who can push
makes the rule advisory for the only person it would bind.

## What this changes

Publishing was: this machine merges to `main` locally, the operator pushes.
That path no longer exists — `main` on the remote advances only through a merged
pull request.

The loop becomes: a lane commits on its own branch here; the operator pushes the
branch; CI runs; the pull request merges when `check` is green. This machine's
`main` tracks `origin/main` and stops being where merges happen.

## Why the strongest option was taken

`AR-6` at `docs/site-requirements.md:332` requires every pull request to get a
real preview URL with `noindex` and drafts included. The preview job runs on
`pull_request` and on nothing else, so with no pull requests that requirement was
**dead** — implemented, wired, and never fired. Requiring a pull request is what
brings it to life.

`AR-5` at `:329` says green means deployable. A required status check is that
sentence made enforceable rather than believed.

## What it costs

Ceremony, on a solo repository, for every change including a typo fix. That is
the honest cost and it was accepted rather than argued away.

It also removes the operator's ability to push a fix directly during an incident.
The escape hatch is to set `enforcement` to `evaluate` or `disabled` on the
ruleset — a deliberate, logged act rather than a quiet push.

## What it does not do

It does not stop this machine rewriting local history, and it does not review
anything: `required_approving_review_count` is `0`, because a personal repository
does not let an author approve their own pull request. The gate is CI, not a
second pair of eyes.
