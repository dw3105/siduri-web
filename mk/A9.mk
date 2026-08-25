.PHONY: a9-tools-selftest a9-pre-p3 a9-html-structure

check: a9-tools-selftest a9-pre-p3 a9-html-structure

a9-tools-selftest:
	@python3 tools/pre_p3.py --selftest
	@python3 tools/html_structure.py --selftest

a9-pre-p3: build
	@python3 tools/pre_p3.py dist

a9-html-structure: build
	@python3 tools/html_structure.py dist
