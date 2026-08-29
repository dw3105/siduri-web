.PHONY: docsguard

check: docsguard

# Contract edits are intentionally delegated to the existing amendcheck and
# Makefile working-tree guard.  This target makes that authority run before
# this guard, while docsguard owns the rest of docs/.  The Python guard also
# excludes commits already reachable from the published trunk.
docsguard: w2-amendcheck
	@set -eu; \
	  head="$$(git rev-parse --verify HEAD 2>/dev/null || true)"; \
	  test -n "$$head" || { echo 'docsguard: FAIL — cannot resolve HEAD' >&2; exit 1; }; \
	  trunk="$$(git rev-parse --verify --quiet 'refs/remotes/origin/main^{commit}' 2>/dev/null || true)"; \
	  if test -z "$$trunk"; then \
	    trunk="$$(git rev-parse --verify --quiet 'refs/heads/main^{commit}' 2>/dev/null || true)"; \
	  fi; \
	  test -n "$$trunk" || { echo 'docsguard: FAIL — cannot resolve published trunk (origin/main or main)' >&2; exit 1; }; \
	  base="$$(git merge-base "$$head" "$$trunk" 2>/dev/null || true)"; \
	  test -n "$$base" || { echo "docsguard: FAIL — cannot resolve a base for HEAD ($$head) and published trunk ($$trunk)" >&2; exit 1; }; \
	  python3 tools/docsguard.py "$$base" "$$head"
