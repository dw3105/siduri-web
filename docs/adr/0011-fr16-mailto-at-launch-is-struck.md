# 0011 · FR-16's launch mailto sentence is struck

**Status** · accepted · 2026-08-28 · operator's word, given directly

## What changed

`FR-16` carried this P0/P1 sentence:

> Phase 0 and 1 ship a plain `mailto:` link and nothing else.

That sentence is struck in place. The adjacent prohibition on a form, Worker, data processing, and spam surface, and the condition against collecting personal data before the §9 questions are settled, remain. The phase-2 paragraph beginning “A form arrives in phase 2 alongside comments” is untouched.

## Why it was asked

The operator annotated the Contact navigation item “Remove, add later”, with the contract note that `FR-16` made a `mailto:` link P0's only contact path. The launch mailto sentence is therefore the exact P0/P1 surface he ruled on.

The scope is deliberately narrow. This is a one-line sentence strike, not a strike of the whole multi-line FR-16 clause; the phase-2 form design remains a separate later-phase decision and requirement.

## What now carries the risk

I grepped `internal/site/contact.templ:10` and `internal/site/pages_a5.templ:96` and found the current mailto implementations. I found `internal/site/pages_a5_test.go:21-34` checks that the Contact page is built, but no test requiring a Contact link or a mailto link. `tools/secretscan.py:39` and `:44-48` allow the reviewed address and rendering files; that is a secret-scan exception, not a contact-path guard.

The current templates therefore still expose the old mailto until the owning shell/navigation change lands, but FR-16 no longer requires it. In the operator's intended Phase 0 state, removing the only mailto sentence and the Contact link leaves no contact path at all, which is stronger than the page merely being unreachable. The surviving no-form/no-processing boundary and the phase-2 form paragraph carry only their stated risks; nothing in the tree replaces the removed P0 contact path.
