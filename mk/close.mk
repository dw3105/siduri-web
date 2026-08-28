.PHONY: close-check close-selftest

# Operator reports are deliberately separate from criterion fragments.  This
# gate validates both reports and runs the isolated red/green harness first.
check: close-check

close-check: close-selftest
	@python3 tools/acceptance.py --operator

close-selftest:
	@python3 tools/acceptance.py --operator-selftest
