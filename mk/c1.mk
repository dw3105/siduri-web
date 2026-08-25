.PHONY: c1-linkcheck-selftest

check: c1-linkcheck-selftest

c1-linkcheck-selftest:
	@python3 tools/linkcheck.py --selftest
