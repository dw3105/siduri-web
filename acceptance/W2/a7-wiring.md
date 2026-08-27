# W2 · A7 — the guards get invoked, and the checkout gets deep

base: 742e376
lane: A7 / lane/w2-wiring

## Wiring — “Nothing invokes `amendcheck`”

did: added `mk/w2amendcheck.mk` with `check: w2-amendcheck`; the target resolves a base and invokes the inherited four-rule checker. The clean gate invokes it before the shared `check` recipe.

ran: `make w2-amendcheck` in a disposable clone after a committed contract edit without an ADR or struck trace

saw:

```text
amendcheck: base=3e7d58e93ace890d594a3496ff892645f1246b33 head=47c32317e51423f396c6c4628ca7e70cbba5b100 diff-context=3
rule 1 (added ADR): FAIL — no added docs/adr/NNNN-*.md file in the range
rule 2 (struck traces): FAIL
  docs/site-requirements.md: hunk 1: FAIL — removed 1 line(s), added 0 non-empty strike span(s); injectivity count is too small
rule 3 (named sections): FAIL
  docs/site-requirements.md: FAIL — touched section name(s) missing from added ADR text: The single job
rule 4 (watched paths): PASS — all watched contract paths exist at base and HEAD
decision: rejected amendment; failed rule 1, rule 2, rule 3
make: *** [mk/w2amendcheck.mk:11: w2-amendcheck] Error 1
exit=2
```

red proof: the planted input was a committed rewrite of the phase-gate sentence in `docs/site-requirements.md`; the target fired rules 1, 2, and 3. `make check` on the restored clean branch exited 0 and printed all four PASS rules with `base=742e376368f67335329a810d0ec47d8a3b3dc087`.

notes: the root `Makefile:30` still has its old range guard because that file is not owned by this lane. The paste-ready replacement is below.

## Base — “a guard that cannot find its base has not passed”

did: the caller reads `HEAD`, tries `git merge-base HEAD main`, treats a missing result or `HEAD` itself as unusable, then tries `HEAD^`; it refuses if the result is empty or still equals `HEAD`. This permits the post-merge `main` check to compare the merge result with its first parent while refusing a depth-1 checkout with no parent.

ran: `make -C "$test_dir" w2-amendcheck` in a `--depth 1` clone with no local `main` ref

saw:

```text
commits=1
has main ref: no
w2-amendcheck: FAIL — cannot resolve a base distinct from HEAD (3e7d58e93ace890d594a3496ff892645f1246b33)
make: *** [mk/w2amendcheck.mk:11: w2-amendcheck] Error 1
make: Leaving directory '/tmp/siduri-w2-a7-depth.rCqlqf'
exit=2
```

red proof: the depth-1/no-`main` input had neither a merge base nor `HEAD^`; it refused instead of passing `HEAD..HEAD`. The three disposable clones and their planted changes were removed; final `git status --porcelain` was checked.

## CI — “Set `fetch-depth: 0` on both checkouts”

did: set `fetch-depth: 0` on the `check` checkout (`.github/workflows/ci.yml:17`) and the `preview` checkout (`.github/workflows/ci.yml:49`).

ran: `python3 tools/workflow_check.py`

saw: `workflow YAML: valid (1 file(s))`

notes: the base resolver and both workflow checkouts must remain aligned; changing checkout depth without retaining the distinct-base refusal would restore a silent vacuous pass.

## Hook — “a fast local warning, not a second mechanical end”

did: added `install-w2-hook`, which writes an executable `.git/hooks/pre-commit` running `w2-contract-guard` and `w2-amendcheck`, plus `w2-hook-status`, which prints `present` only when that hook is executable and contains both calls. The contract guard checks staged and unstaged requirements files; amendcheck checks the committed range.

ran: `make install-w2-hook`

saw:

```text
install-w2-hook: installed .git/hooks/pre-commit
w2-hook-status: present (.git/hooks/pre-commit)
```

ran: `git add docs/site-requirements.md && git commit -m 'hook breach proof'; code=$?; echo "commit-exit=$code"; exit 0`

saw:

```text
w2-contract-guard: FAIL — the two requirements files are the contract and only the operator amends them
make: *** [mk/w2amendcheck.mk:25: w2-contract-guard] Error 1
commit-exit=1
```

red proof: the staged phase-gate rewrite was refused by the installed hook. The fixture was restored before the clean commit probe.

ran: `git restore --staged --worktree docs/site-requirements.md && git status --porcelain && git commit --allow-empty -m 'hook clean proof'; code=$?; echo "clean-commit-exit=$code"; git status --porcelain; exit 0`

saw:

```text
w2-contract-guard: PASS — no unstaged or staged contract edits
rule 1 (added ADR): PASS — no contract diff; no amendment record is required
rule 2 (struck traces): PASS
rule 3 (named sections): PASS
rule 4 (watched paths): PASS — all watched contract paths exist at base and HEAD
decision: no contract diff; accepted
clean-commit-exit=0
```

notes: the hook is local, opt-in, absent by default, and bypassable with `--no-verify`; it is a fast local warning, not a second mechanical end. It deliberately does not run the build, Go tests or vet, templ formatting, rendered-site gates, link/accessibility checks, workflow validation, or the secret scanner. `w2-hook-status` reports `absent` until `install-w2-hook` is run in a clone.

## Makefile handoff

The behavior being replaced at `Makefile:30` is:

```text
git diff --quiet "$$base" HEAD -- $$contract && git diff --quiet -- $$contract && git diff --cached --quiet -- $$contract
```

That range comparison rejects every committed requirements diff; the final two comparisons reject dirty working-tree and index requirements edits. Paste-ready replacement for the recipe line:

```make
@set -eu; head="$$(git rev-parse --verify HEAD 2>/dev/null || true)"; test -n "$$head" || { echo 'check: cannot resolve HEAD' >&2; exit 1; }; base="$$(git merge-base "$$head" main 2>/dev/null || true)"; if test -z "$$base" || test "$$base" = "$$head"; then base="$$(git rev-parse --verify "$$head^" 2>/dev/null || true)"; fi; test -n "$$base" && test "$$base" != "$$head" || { echo "check: cannot resolve a base distinct from HEAD ($$head)" >&2; exit 1; }; contract='docs/site-requirements.md docs/comments-requirements.md'; python3 tools/amendcheck.py "$$base" "$$head" && git diff --quiet -- $$contract && git diff --cached --quiet -- $$contract || { echo 'check: the two requirements files are the contract and only the operator amends them' >&2; exit 1; }
```

## Verification

ran: `go build ./...`

saw: exit 0, no output.

ran: `go test ./internal/site/`

saw: `ok github.com/dw3105/siduri-web/internal/site 0.982s`.

ran: `python3 tools/amendcheck.py --selftest`

saw: `amendcheck selftest: 13 cases pass`.

ran: `make build && find dist -type f -exec sha256sum {} + | sort > /tmp/siduri-w2-a7-b1 && find dist -depth -type f -delete && find dist -depth -type d -empty -delete && make build && find dist -type f -exec sha256sum {} + | sort > /tmp/siduri-w2-a7-b2 && diff /tmp/siduri-w2-a7-b1 /tmp/siduri-w2-a7-b2 && echo DETERMINISTIC`

saw: `DETERMINISTIC`.

ran: `make check`

saw: exit 0; the output included the four-rule no-diff decision with the resolved base shown above.

ran: `git status --porcelain`

saw: no output after the report was added and committed.
