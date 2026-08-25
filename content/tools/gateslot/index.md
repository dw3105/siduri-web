---
title: Gateslot
slug: gateslot
date: 2026-08-25
summary: A small approval gate for changes that become public.
language: Go
status: active
repository: https://github.com/dw3105/gateslot
install: git clone https://github.com/dw3105/gateslot && go install ./cmd/gateslot
example: |
  gateslot check --path ./draft.md --require human-approval
---
Gateslot exists to make the last human decision explicit in an otherwise automated workflow. It records what is ready, what still needs judgement, and what must remain a draft.
