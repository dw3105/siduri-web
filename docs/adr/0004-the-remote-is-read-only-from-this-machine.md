# 0004 · The remote is read-only from this machine

**Status** · accepted · 2026-08-25

`origin` is `git@github-siduri-ro:dw3105/siduri-web.git` for fetch and the literal
string `DISABLED` for push. The alias resolves through `~/.ssh/config` to
`~/.ssh/siduri_ro`, an ed25519 key registered on the repository as a
**read-only deploy key**.

This machine commits and fetches. It cannot push, and the failure mode if it
tries is a refused credential rather than a rewritten branch.

## Why asymmetric rather than careful

`AR-8` forbids an agent force-pushing or rewriting history. A prohibition in a
document is checkable only after the fact. A credential that cannot push is
checkable before it.

The three sibling repositories on this machine already do this — `hitl` and
`gateslot` through `github-<repo>-ro` aliases, `legalcopilot` through a push URL
of `DISABLED`. This repository takes both halves: the alias *and* the disabled
push URL, so the refusal happens at the git layer even if the alias is bypassed.

## What this costs

Publishing is two machines. This one commits; the operator pushes. Nothing here
can complete a publish alone, which is the point and is also the friction.

The initial history reached GitHub as a bundle rather than a push, for the same
reason: giving this machine a write credential once would have made the asymmetry
a habit rather than a property.

## What it does not defend against

The deploy key is read-only on **this** repository only. It says nothing about
what any other credential on this machine can do, and it is not a sandbox. An
agent here can still write anything it likes to the local tree and to
`content/`; what it cannot do is make that public without a human.
