---
name: siduri-code
description: "Govern Siduri lane ownership, content-derived registration, merged-tree checks, dependency boundaries, and pull-request shipping."
---

## skill check

siduri-code: FO-03 FO-04 FO-07 FO-12 FO-13 FO-14 FO-15 FO-17 SH-01a SH-02 SH-12 govern lane ownership, registration timing, gate ownership, seam scans, integrator files, dependencies, merged-tree prediction, and pull-request shipping.

Read first: `codex-tasks`.

## Lane and shipping rules

- **FO-03** **Lane own only new file. File two lane need belong to integrator.** Merge conflict-free by construction, never by care. [bound-with-gap]
- **FO-04** **Registration happen when data exist, never before.** Route registered at `init()` cannot know tag, tool or thread. Symptom look like output cardinality; cause is lifecycle. [bound-with-gap]
- **FO-07** **Worker run no gate.** Eight full gate on four core is box thrashing, and gate cannot see conflict from inside one lane anyway. Integrator run it once, on merged tree. [unbound]
- **FO-12** **Filename named by two lane task is seam defect, and it is grep.** `tools/lane_overlap.py` over task set before any worktree cut. Non-empty output mean wave not ready. Fix seam, never route lane around it. [bound-with-gap]
- **FO-13** **Integrator-owned file get threshold of one, never two.** `go.mod`, `go.sum`, `Makefile`, `routes.go`, `build.go`, `content.go`, `markdown.go`. Task naming one already claim file lane do not own. Suppressing them read cuttable on wave that collide four way. [bound]
- **FO-14** **Task text is declared intent. Branch is behaviour.** Lane writing file it never named is invisible to any regex. `lane_overlap.py --branches` before merge, pairwise over `git diff --name-only`. First prevent, second detect, neither substitute. [bound-with-gap]
- **FO-15** **Every dependency resolved into `go.mod` before any branch cut. No lane run `go get`.** Eight lane adding library is real conflict, never accidental, and `go.sum` land on integrator. Dependency claimed on unmerged branch is not claimed. [unbound]
- **FO-17** **Merged tree say nothing about either branch alone.** `CX-22` converse. Gating C1 and C3 together then predicting each PR gave two wrong prediction, opposite direction. [unbound]
- **SH-01a** **Nothing merge locally. `main` advance only through pull request, no bypass.** Lane commit on own branch, operator push branch, CI `check` gate merge. Lead's `main` track `origin/main` and stop being where merge happen. [bound-with-gap]
- **SH-02** **`make check` before every commit on `main`.** [unbound]
- **SH-12** **Publish path is target with refusal, never block pasted into terminal.** Bootstrap have no target so one block unavoidable; every block after that is debt with shell prompt. [bound-with-gap]
