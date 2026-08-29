.PHONY: docsguard

check: docsguard

# Contract edits are intentionally delegated to the existing amendcheck and
# Makefile working-tree guard.  This target makes that authority run before
# this guard, while docsguard owns the rest of docs/.
docsguard: w2-amendcheck
	@set -eu; \
	  head="$$(git rev-parse --verify HEAD 2>/dev/null || true)"; \
	  test -n "$$head" || { echo 'docsguard: FAIL — cannot resolve HEAD' >&2; exit 1; }; \
	  base="$$(git merge-base "$$head" main 2>/dev/null || true)"; \
	  if test -z "$$base" || test "$$base" = "$$head"; then \
	    base="$$(git rev-parse --verify "$$head^" 2>/dev/null || true)"; \
	  fi; \
	  test -n "$$base" && test "$$base" != "$$head" || { echo "docsguard: FAIL — cannot resolve a base distinct from HEAD ($$head)" >&2; exit 1; }; \
	  python3 tools/docsguard.py "$$base" "$$head"
