.PHONY: acceptance acceptance-selftest

acceptance-selftest:
	@python3 tools/acceptance.py --selftest

acceptance: acceptance-selftest
	@python3 tools/acceptance.py
