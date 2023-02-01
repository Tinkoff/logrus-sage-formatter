lint:
	@golangci-lint run --timeout=5m ./... -v

test:
	@go test --race --count=1 ./... -v