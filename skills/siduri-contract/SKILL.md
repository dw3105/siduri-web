---
name: siduri-contract
description: "Governs Siduri's contract and requirement documents, and where decisions and findings are recorded."
---

## skill check

Eleven rules: read the requirements; amend only through an ADR and affected checks; keep workers out of `docs/`; record rulings and findings; keep canonical skill edits ahead of deployed copies.

## Rules

**RD-01** Never start from title. Read all four. [self-reported]
- **RD-05** **Read requirement whole, by range.** grep locate; grep never read. [self-reported]
- **CT-02** **Amendment is three part in one commit**: ADR naming clause, reason, replacement; clause struck in place and never deleted; every affected check edited. [bound-with-gap]
- **CT-03** **Struck gate get its risk relocated, named, in same ADR.** Gate removed with nothing behind it is failure that gate predicted. [unbound]
- **CT-04** **No worker edit `docs/`.** Bound by path check, never by prose. Worker habit is to fix document found wrong, and that habit amend contract. Case: acceptance/ outside docs/; tools/secretscan.py exempts docs/, so quoted address made make check red; no record could land. [bound-with-gap]
- **SS-02** **Operator ruling get its `docs/DECISIONS.md` row in turn it is given, never at close.** Four ruling in one day, zero row, recovered only because one transcript survived one `/compact`. Row is cheap in the turn and unwritable a day later. [unbound]
- **SS-03** **Finding row and decision row commit unasked. Only rule change get asked.** Asking permission to record is how recording become optional. [unbound]
- **SS-05** **Clause written into deployed skill and not into `skills/` is lost, and nothing report it.** This repo's skill-check compared canonical/deployed by byte digest; tools/skill_install.py read frontmatter with partition(":") and agreed for weeks with a file PyYAML refuses outright. [bound]
- **RV-06** **Decision land row in `docs/DECISIONS.md`.** Message thread bind nothing. [unbound]
- **CL-02** **Finding outside scope go in `docs/FINDINGS.md`, in closing commit.** Row carry where found, what break, what it look like when it break. [unbound]
- **CL-05** **Finding withdrawn after checking get row, and row carry probe that withdrew it.** Row exist or it do not, so this is checkable; *act of looking is evidence* is not, and would join unbound table rather than bind anything. Withdrawn finding leaving no trace get re-opened by next reader with no record anybody looked. Requires `docs/FINDINGS.md`'s negative-result row shape, added same day. [unbound]
