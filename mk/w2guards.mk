.PHONY: w2-guards w2-guards-selftest

check: w2-guards w2-guards-selftest

w2-guards: build
	@python3 tools/draftscan.py dist
	@python3 tools/storagescan.py dist

w2-guards-selftest:
	@python3 tools/draftscan.py --selftest
	@python3 tools/storagescan.py --selftest
