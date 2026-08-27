# `acceptance` is deliberately NOT a prerequisite of `check`.
# Wiring it now deadlocks the bootstrap: the acceptance is not accepted (three
# open, two failed), so `check` would go red on the very pull request that
# carries the report, and the report could never land. Wire it the day the
# banner reads `accepted`, and not before. Until then nothing invokes this
# target automatically, which is a gap and is recorded as one.

.PHONY: acceptance acceptance-selftest

acceptance-selftest:
	@python3 tools/acceptance.py --selftest

acceptance: acceptance-selftest
	@python3 tools/acceptance.py
