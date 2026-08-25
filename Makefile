SHELL := /bin/sh

GO_BIN := $(shell go env GOPATH)/bin
TEMPL := $(GO_BIN)/templ
DEV_DIR := .dev-dist
AGENT_PARTS := AGENTS-PROHIBITIONS.md AGENTS-FILLABLE.md
MK_FILES := $(sort $(wildcard mk/*.mk))

.PHONY: dev build assemble-agents check publish comments deploy rollback

-include $(MK_FILES)

assemble-agents: $(AGENT_PARTS)
	@tmp=$$(mktemp AGENTS.md.XXXXXX); trap 'rm -f "$$tmp"' EXIT HUP INT TERM; cat $(AGENT_PARTS) > "$$tmp" && { cmp -s "$$tmp" AGENTS.md || cp "$$tmp" AGENTS.md; }

dev:
	@PATH="$(GO_BIN):$$PATH" $(TEMPL) generate
	@go run ./cmd/siduri dev --output $(DEV_DIR)

build: assemble-agents
	@PATH="$(GO_BIN):$$PATH" $(TEMPL) generate
	@go run ./cmd/siduri build --output dist

check: build
	@test -z "$$(gofmt -l $$(rg --files -g '*.go' -g '!vendor/**'))" || { echo 'check: gofmt found unformatted files' >&2; exit 1; }
	@go vet ./...
	@PATH="$(GO_BIN):$$PATH" $(TEMPL) fmt -fail .
	@go test ./internal/site -run '^TestGolden' -count=1
	@go test ./...
	@base="$$(git merge-base HEAD main 2>/dev/null || git rev-list --max-parents=0 HEAD)"; contract='docs/site-requirements.md docs/comments-requirements.md'; git diff --quiet "$$base" HEAD -- $$contract && git diff --quiet -- $$contract && git diff --cached --quiet -- $$contract || { echo 'check: the two requirements files are the contract and only the operator amends them' >&2; exit 1; }
	@base="$$(git merge-base HEAD main 2>/dev/null || git rev-list --max-parents=0 HEAD)"; protected="$$(git diff --name-only --diff-filter=MD "$$base" HEAD -- AGENTS-PROHIBITIONS.md; git diff --name-only --diff-filter=MD -- AGENTS-PROHIBITIONS.md; git diff --name-only --diff-filter=MD --cached -- AGENTS-PROHIBITIONS.md)"; test -z "$$protected" || { echo 'check: AGENTS-PROHIBITIONS.md is contract-owned and must not change' >&2; exit 1; }
	@find . -type f \( -name 'AGENTS.md' -o -name 'AGENTS*.md' -o -name '.agents.md' \) -print | while IFS= read -r file; do test "$$(wc -c < "$$file")" -lt 32768 || { echo "check: $$file is at or above 32768 bytes" >&2; exit 1; }; done
	@if grep -rE '[a-z0-9._%+-]+@[a-z0-9.-]+' content/; then echo 'check: raw email address found under content/' >&2; exit 1; fi
	@if test -d dist && grep -rIE '/services|([€$$][[:space:]]*[0-9])|<form([[:space:]>]|$$)' dist; then echo 'check: pre-P3 sales or form artifact found in dist/' >&2; exit 1; fi

publish:
	@echo 'publish: Wave P1 fills this target' >&2; exit 1

comments:
	@echo 'comments: Wave P2 fills this target' >&2; exit 1

deploy:
	@echo 'deploy: the deployment lane in Wave P0 fills this target' >&2; exit 1

rollback:
	@echo 'rollback: the deployment lane in Wave P0 fills this target' >&2; exit 1
