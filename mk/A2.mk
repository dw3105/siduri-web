.PHONY: a2

a2:
	@go test ./internal/site -run A2 -count=1
