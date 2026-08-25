.PHONY: a4-budget

a4-budget:
	@external=$$(gzip -n -c static/site.css | wc -c); critical=$$(gzip -n -c internal/site/assets_a4.css | wc -c); total=$$((external + critical)); echo "CSS gzip: $$total bytes (external $$external bytes + critical $$critical bytes; budget 15360 bytes)"; test "$$total" -le 15360
