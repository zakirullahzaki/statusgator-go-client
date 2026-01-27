.PHONY: test test-coverage lint fmt vet build clean

# Run all tests
test:
	go test ./... -v -race -count=1

# Run tests with coverage
test-coverage:
	go test ./... -v -race -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Run go vet
vet:
	go vet ./...

# Build
build:
	go build ./...

# Clean
clean:
	rm -f coverage.out coverage.html

# Run all checks
check: fmt vet lint test
