.PHONY: w2-amendcheck w2-contract-guard install-w2-hook w2-hook-status

W2_CONTRACT_FILES := docs/site-requirements.md docs/comments-requirements.md

# The range check belongs before the shared check recipe.  Its base must be a
# real ancestor: a missing main ref or a post-merge HEAD may not silently turn
# the range into HEAD..HEAD.
check: w2-amendcheck

w2-amendcheck:
	@set -eu; \
	head="$$(git rev-parse --verify HEAD 2>/dev/null || true)"; \
	test -n "$$head" || { echo 'w2-amendcheck: FAIL — cannot resolve HEAD' >&2; exit 1; }; \
	base="$$(git merge-base "$$head" main 2>/dev/null || true)"; \
	if test -z "$$base" || test "$$base" = "$$head"; then \
		base="$$(git rev-parse --verify "$$head^" 2>/dev/null || true)"; \
	fi; \
	test -n "$$base" && test "$$base" != "$$head" || { echo "w2-amendcheck: FAIL — cannot resolve a base distinct from HEAD ($$head)" >&2; exit 1; }; \
	python3 tools/amendcheck.py "$$base" "$$head"

# This is the local part of the contract guard.  The committed range is
# checked by amendcheck; the two working-tree checks keep a commit from
# carrying a contract edit that has not gone through the range check.
w2-contract-guard:
	@set -eu; \
	if git diff --quiet -- $(W2_CONTRACT_FILES) && git diff --cached --quiet -- $(W2_CONTRACT_FILES); then \
		echo 'w2-contract-guard: PASS — no unstaged or staged contract edits'; \
	else \
		echo 'w2-contract-guard: FAIL — the two requirements files are the contract and only the operator amends them' >&2; \
		exit 1; \
	fi

install-w2-hook:
	@set -eu; \
	hooks="$$(git rev-parse --git-path hooks)"; \
	mkdir -p "$$hooks"; \
	python3 -c 'import pathlib, stat, sys; path = pathlib.Path(sys.argv[1]); path.write_text("#!/bin/sh\nset -eu\nmake --no-print-directory w2-contract-guard\nmake --no-print-directory w2-amendcheck\n", encoding="utf-8"); path.chmod(path.stat().st_mode | stat.S_IXUSR | stat.S_IXGRP | stat.S_IXOTH)' "$$hooks/pre-commit"; \
	echo "install-w2-hook: installed $$hooks/pre-commit"; \
	$(MAKE) --no-print-directory w2-hook-status

w2-hook-status:
	@set -eu; \
	hook="$$(git rev-parse --git-path hooks)/pre-commit"; \
	if test -x "$$hook" && grep -Fq 'w2-contract-guard' "$$hook" && grep -Fq 'w2-amendcheck' "$$hook"; then \
		echo "w2-hook-status: present ($$hook)"; \
	else \
		echo "w2-hook-status: absent ($$hook)"; \
	fi
