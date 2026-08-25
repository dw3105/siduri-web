# contract - requirement, amendment, ADR

**v1 - 25 Aug 2026.** Read when a requirement is touched, an ADR written, or a
requirement look wrong. `CT` rules live in `SKILL.md`.

---

## Two file, copied verbatim, amended by one person

`docs/site-requirements.md` and `docs/comments-requirements.md` are the operator's
own drafts, byte-identical to source. They carry numbered decision, functional
requirement, non-functional budget, agent rule, legal rule, acceptance criteria
and open question. Build make them true; build does not redesign them.

## Rename inside a contract dangle every reference to it

First commit renamed both spec on the way in. Four reference inside the site spec
pointed at the old name - front matter, the public-repo consequence routing to the
no-commenter-email rule, the comment summary, and the cost table.

The contract cannot be edited to repair it, because only the operator amend it. So
the rename created a defect inside the one directory nobody but him may touch.

**Repair was the rename reversed.** Filename is not a clause: zero contract byte
changed, four reference resolved, no amendment spent.

## Amendment is three part, and the third is the one forgotten

ADR naming clause, reason, replacement. Clause struck **in place**, never deleted.
Every affected check edited.

Struck gate get its risk relocated in the same ADR, named. A gate removed with
nothing behind it is the failure the gate predicted. When the P1 publication gate
was struck, the risk it named moved onto an end-to-end Done-when performed by a
person - and the residual, that this is weaker than a gate in exactly the
direction the struck text warned about, was recorded rather than argued away.

Record also what was **not** struck. Two sibling gate stood, and a reader seeing
one struck will assume the set went.

## Contract can forbid the order the operator asked for

Instruction was build everything before commercialization, eight lane parallel.
The contract's own most emphatic line forbade beginning the next phase before
publication.

That is not a thing to resolve by reading harder. It went to him as a
contradiction, named, with the option to strike. He struck it. It is his document.

## Guard bind lane branch, not the integrator

`docs/` guard compare merge-base to HEAD. On `main` those are the same commit, so
the committed half is a no-op there. Working-tree and staged half still bind
everywhere.

That is defensible - the integrator carry the operator's amendments - but it must
be said in the ADR. Next reader see green `make check` on `main` and read it as
the contract being guarded there.
