.PHONY: c3

c3:
	@go test ./internal/site -run C3 -count=1
