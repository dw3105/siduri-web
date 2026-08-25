# 0003 · The P1 publication gate is struck

**Status** · accepted · 2026-08-25 · operator's word, given directly

## What changed

`site-requirements.md` §12 carried, between P0 and P1:

> **Gate: post 1 is published on a real domain before any P1 work begins.** This
> is the single most important line in the document.

Struck in place, not deleted. P0, P1 and P2 are built as one programme.

## Why it was asked

The gate and the instruction collided. The instruction was to build everything
before commercialization — P0, P1 and P2 — with up to eight agents in parallel.
Read off §4's priority column, `/tags/`, `/search/`, `/stack/` and `/llms.txt`
are P2 and `/tools/` is P1, so a single programme puts most of the work behind a
gate no agent can open: publication to a live domain.

The gate named the risk precisely, and named agents as the thing that makes it
worse: *the characteristic death of this project is six weeks of beautiful
scaffolding and zero published posts — and building it with agents makes
scaffolding more seductive, not less.*

Raised by `siduri_reviewer`, re-derived here against §4 and §12 before it went to
the operator. He struck it. It is his document.

## What now carries the risk

The gate's reasoning is not withdrawn. It is relocated, because a struck gate
with nothing behind it is the failure the gate predicted.

**Every wave's Done-when is end-to-end and performed by a person.** Not a file,
not a suite, not a green pipeline. For this repo the line is: *open the built
site in a browser and read post 1, with JavaScript disabled* —
`site-requirements.md:513` verbatim, which the contract already required.

This is the `SS-10` shape: five packages, five sessions, one night, 1854 tests
passed uncached, 39 of 40 mutations caught, CI green through nineteen steps,
image pushed — and nothing served the page. Every package's Done-when named its
own file and every one was met honestly. Eight lanes is that with three more.

The tell, run on the lane tasks before they are handed over: read all the
Done-whens and ask which names the product. If they all name files, nobody is
testing composition.

**Post 1 stays P0 scope.** §12 P0 is unamended on this point: *ships with the
site; the site is not live without it.* An agent may draft it. `AR-8` forbids an
agent publishing it.

## What was not struck

- **§12:497**, the P1→P2 gate — *Gateslot is usable by someone who is not the
  operator, and post 2 is published.* Not raised, not ruled on, still standing.
  Open in `docs/DECISIONS.md`. It blocks wave B, not wave 0 or wave A.
- **§12:502**, *a postal address exists. No list, no form and no comment intake
  before it does.* Standing, and it should. It is a shipping gate on collecting
  personal data, not a build gate — `LR-7` and `FR-18`. Building the intake code
  is not turning intake on. This is what a build-time phase gate is legitimately
  for, and the only thing it is for.

## Residual risk, accepted and recorded

Nothing now stops eight agents scaffolding for weeks except the end-to-end
Done-when and the operator's attention. That is a weaker mechanism than a gate,
and it is weaker in exactly the direction the struck text warned about. Recorded
here rather than argued away.
