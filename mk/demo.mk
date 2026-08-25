.PHONY: demo

demo:
	@go test ./internal/site -run '^(TestDemoLane|TestBuildTwiceWithContentRegistration)$$' -count=1
