.PHONY: a8-gates a8-budget a8-headers a8-links a8-html a8-accessibility a8-secrets a8-workflow install-hooks preview-build preview preview-deploy a8-rollback

A8_PYTHON ?= python3
A8_EXTERNAL_LINKS ?= 0
A8_AXE_PORT ?= 8765

# The original check target already builds, runs the golden tests, vets and
# tests the Go tree.  This prerequisite adds the A8 contract without editing
# the shared Makefile.
check: a8-gates

a8-gates: a8-budget a8-headers a8-links a8-html a8-accessibility a8-secrets a8-workflow

a8-budget:
	@set -eu; \
	start=$$(date +%s%N); \
	$(MAKE) --no-print-directory build >/dev/null; \
	end=$$(date +%s%N); \
	seconds=$$(awk -v start="$$start" -v end="$$end" 'BEGIN { printf "%.6f", (end-start)/1000000000 }'); \
	$(A8_PYTHON) tools/budget.py --dist dist --build-seconds "$$seconds"

a8-headers:
	@test -f _headers
	@grep -Fq 'Content-Security-Policy:' _headers
	@grep -Fq "script-src 'self'" _headers
	@grep -Fq "frame-ancestors 'none'" _headers
	@if grep -Fq 'unsafe-inline' _headers; then echo 'a8-headers: unsafe-inline is forbidden' >&2; exit 1; fi
	@grep -Eq 'Strict-Transport-Security:.*max-age=[0-9]+.*includeSubDomains' _headers
	@grep -Fq 'X-Content-Type-Options: nosniff' _headers
	@grep -Fq 'Referrer-Policy: strict-origin-when-cross-origin' _headers
	@grep -Eq 'Permissions-Policy:.*[A-Za-z-]+=\(\)' _headers
	@echo 'a8-headers: required security directives verified'

a8-links:
	@if test "$(A8_EXTERNAL_LINKS)" = 1; then \
		$(A8_PYTHON) tools/linkcheck.py --dist dist --external; \
	else \
		$(A8_PYTHON) tools/linkcheck.py --dist dist; \
	fi

# These validators are transient CI tools, not application dependencies. The
# versions stay on their major lines so a fresh checkout needs no package file.
a8-html:
	@files="$$(find dist -type f -name '*.html' -print)"; \
	test -n "$$files" || { echo 'a8-html: no HTML files found' >&2; exit 1; }; \
	npx --yes html-validate@8 --rule doctype-style:off $$files

a8-accessibility:
	@major=$$(node -p 'process.versions.node.split(".")[0]' 2>/dev/null || echo 0); \
	if [ "$$major" -lt 20 ]; then \
	  echo "a8-accessibility: SKIPPED on node $$major -- axe-core needs >= 20." >&2; \
	  echo "a8-accessibility: this stops measuring WCAG 2.2 AA entirely (NFR-2). CI pins 20 and runs it." >&2; \
	  exit 0; \
	fi; \
	set -eu; \
	(port=$$($(A8_PYTHON) -c 'import os; print(os.environ.get("A8_AXE_PORT", "8765"))'); \
	cd dist; $(A8_PYTHON) -m http.server "$$port" >/tmp/siduri-a8-http.log 2>&1) & \
	server=$$!; \
	trap 'kill "$$server" 2>/dev/null || true' EXIT HUP INT TERM; \
	port="$(A8_AXE_PORT)"; \
	for attempt in 1 2 3 4 5 6 7 8 9 10; do \
		if curl -fsS "http://127.0.0.1:$$port/index.html" >/dev/null; then break; fi; \
		sleep 1; \
	done; \
	curl -fsS "http://127.0.0.1:$$port/index.html" >/dev/null; \
	urls="$$(find dist -type f -name '*.html' -print | sed 's#^dist#http://127.0.0.1:'"$$port"'#')"; \
	npx --yes @axe-core/cli@4.10.2 --exit $$urls

a8-secrets:
	@$(A8_PYTHON) tools/secretscan.py .

# The transient Node YAML parser catches malformed workflow YAML in a fresh
# checkout without adding a YAML package to this repository.
a8-workflow:
	@npx --yes -p yaml@2 node -e 'const fs=require("fs"), path=require("path"), yaml=require("yaml"); const files=fs.readdirSync(".github/workflows").filter(f => /\.ya?ml$$/.test(f)); if (!files.length) throw Error("no workflow files"); for (const file of files) yaml.parse(fs.readFileSync(path.join(".github/workflows", file), "utf8")); console.log("workflow YAML: valid")'

install-hooks:
	@$(A8_PYTHON) tools/secretscan.py --install-hook

# cmd/siduri's dev command is the existing draft-inclusive build. It also
# serves noindex locally; copying the metadata into the upload directory makes
# the same noindex policy apply to the Cloudflare preview version.
preview-build:
	@set -eu; \
	rm -rf .preview-dist; \
	setsid sh -c 'PATH="$$(go env GOPATH)/bin:$$PATH" exec go run ./cmd/siduri dev --output .preview-dist' >/tmp/siduri-preview.log 2>&1 & \
	server=$$!; \
	trap 'kill -- -"$$server" 2>/dev/null || true' EXIT HUP INT TERM; \
	for attempt in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20; do \
		if test -f .preview-dist/index.html; then break; fi; \
		sleep 1; \
	done; \
	test -f .preview-dist/index.html; \
	cp _headers .preview-dist/_headers; \
	{ echo; echo '/*'; echo '  X-Robots-Tag: noindex'; } >> .preview-dist/_headers; \
	cp _routes.json .preview-dist/_routes.json; \
	kill -- -"$$server" 2>/dev/null || true; \
	wait "$$server" 2>/dev/null || true; \
	echo 'preview-build: drafts included and noindex metadata staged in .preview-dist'

preview-deploy: preview-build
	@wrangler versions upload .preview-dist --env preview --message "Siduri pull request preview"

preview: preview-deploy

# Cloudflare's rollback command selects the version uploaded before the latest
# one when no ID is supplied. It creates a new deployment from that version;
# it does not rebuild this checkout. A8_DRY_RUN=1 prints the exact command.
a8-rollback:
	@if test "$(A8_DRY_RUN)" = 1; then \
		echo 'dry run: wrangler rollback --name siduri-web --message "Siduri rollback"'; \
	else \
		command -v wrangler >/dev/null || { echo 'rollback: wrangler is required (not installed in this lane)' >&2; exit 1; }; \
		wrangler rollback --name siduri-web --message 'Siduri rollback'; \
	fi

# Makefile's shared placeholder remains later in parse order. The real work is
# the prerequisite above. A target-specific shell makes that known placeholder
# a no-op while preserving normal failure propagation from a8-rollback.
rollback: a8-rollback
a8-rollback: SHELL := /bin/sh
rollback: SHELL := /bin/true
