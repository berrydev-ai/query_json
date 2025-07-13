# Makefile for query_json - JSON query tool
# Generated from .claude_commands.json commands

.PHONY: build build-release build-all test test-coverage run-example run-help clean fmt vet lint mod-tidy mod-download mod-update deps size install benchmark profile create-example demo validate release-prep check-deps init doc security cross-test perf git-setup

# Build the application for development (current platform)
build:
	go build -o query_json

# Build optimized binary for release (current platform)
build-release:
	go build -ldflags='-s -w' -o query_json

# Build for all platforms (requires build.sh)
build-all:
	chmod +x build.sh && ./build.sh

# Run all tests
test:
	go test -v ./...

# Run tests with coverage report
test-coverage:
	go test -v -cover ./...

# Run with example data (requires examples/test.json)
run-example:
	./query_json examples/test.json --query '$.users[0].name'

# Show application help
run-help:
	./query_json --help

# Clean build artifacts
clean:
	rm -rf builds/ query_json query_json.exe query_json-*

# Format Go code
fmt:
	go fmt ./...

# Run go vet for static analysis
vet:
	go vet ./...

# Run golint (install with: go install golang.org/x/lint/golint@latest)
lint:
	golint ./...

# Clean up go.mod and go.sum files
mod-tidy:
	go mod tidy

# Download all dependencies
mod-download:
	go mod download

# Update all dependencies to latest versions
mod-update:
	go get -u ./... && go mod tidy

# Show dependency graph
deps:
	go mod graph

# Show binary size (requires built binary)
size:
	ls -lh query_json* 2>/dev/null || echo 'Binary not found. Run make build first.'

# Install binary to $GOPATH/bin
install:
	go install

# Run benchmarks (if any exist)
benchmark:
	go test -bench=. -benchmem

# Build with profiling enabled
profile:
	go build -ldflags='-s -w' -o query_json && echo 'Run with: ./query_json [args] && go tool pprof query_json cpu.prof'

# Create example test data file
create-example:
	mkdir -p examples && cat > examples/test.json << 'EOF'
{
  "users": [
    {
      "name": "Alice Johnson",
      "age": 30,
      "email": "alice@example.com",
      "department": "Engineering"
    },
    {
      "name": "Bob Smith",
      "age": 25,
      "email": "bob@example.com",
      "department": "Marketing"
    }
  ],
  "company": "Tech Corp"
}
EOF
	echo 'Created examples/test.json'

# Run a series of demo queries (requires examples/test.json)
demo:
	echo '=== Demo Queries ===' && echo '1. Get all users:' && ./query_json examples/test.json --query '$.users[*].name' && echo -e '\n2. Get user by age filter:' && ./query_json examples/test.json --query '$.users[?(@.age > 25)]' && echo -e '\n3. Raw email output:' && ./query_json examples/test.json --query '$.users[*].email' --raw

# Run full validation suite (format, vet, test)
validate:
	echo 'Running validation...' && go fmt ./... && go vet ./... && go test ./... && echo 'All validation passed!'

# Prepare for release (validate, build all platforms)
release-prep:
	echo 'Preparing release...' && go fmt ./... && go vet ./... && go test ./... && go mod tidy && chmod +x build.sh && ./build.sh && echo 'Release preparation complete!'

# Check for outdated dependencies
check-deps:
	go list -u -m all

# Initialize project (run after cloning)
init:
	go mod download && mkdir -p builds examples && echo 'Project initialized!'

# Generate and serve documentation
doc:
	godoc -http=:6060 && echo 'Documentation available at http://localhost:6060'

# Run security checks (requires gosec: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)
security:
	gosec ./...

# Test cross-compilation without building
cross-test:
	echo 'Testing cross-compilation...' && GOOS=windows GOARCH=amd64 go build -o /dev/null && GOOS=darwin GOARCH=amd64 go build -o /dev/null && GOOS=linux GOARCH=amd64 go build -o /dev/null && echo 'Cross-compilation test passed!'

# Performance test with large JSON file
perf:
	echo 'Creating large test file...' && python3 -c "import json; data={'users': [{'name': f'User{i}', 'age': 20+i%40, 'email': f'user{i}@example.com'} for i in range(10000)]}; open('large_test.json', 'w').write(json.dumps(data))" && echo 'Testing performance...' && time ./query_json large_test.json --query '$.users[?(@.age > 30)]' > /dev/null && rm large_test.json

# Set up git hooks and configuration
git-setup:
	echo 'Setting up git hooks...' && mkdir -p .git/hooks && echo '#!/bin/bash\ngo fmt ./...\ngo vet ./...' > .git/hooks/pre-commit && chmod +x .git/hooks/pre-commit && echo 'Pre-commit hook installed!'

# Default target
all: build

# Help target
help:
	@echo "Available targets:"
	@echo "  build          - Build the application for development"
	@echo "  build-release  - Build optimized binary for release"
	@echo "  build-all      - Build for all platforms"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  run-example    - Run with example data"
	@echo "  run-help       - Show application help"
	@echo "  clean          - Clean build artifacts"
	@echo "  fmt            - Format Go code"
	@echo "  vet            - Run go vet for static analysis"
	@echo "  lint           - Run golint"
	@echo "  mod-tidy       - Clean up go.mod and go.sum files"
	@echo "  mod-download   - Download all dependencies"
	@echo "  mod-update     - Update all dependencies"
	@echo "  deps           - Show dependency graph"
	@echo "  size           - Show binary size"
	@echo "  install        - Install binary to GOPATH/bin"
	@echo "  benchmark      - Run benchmarks"
	@echo "  profile        - Build with profiling enabled"
	@echo "  create-example - Create example test data file"
	@echo "  demo           - Run demo queries"
	@echo "  validate       - Run full validation suite"
	@echo "  release-prep   - Prepare for release"
	@echo "  check-deps     - Check for outdated dependencies"
	@echo "  init           - Initialize project"
	@echo "  doc            - Generate and serve documentation"
	@echo "  security       - Run security checks"
	@echo "  cross-test     - Test cross-compilation"
	@echo "  perf           - Performance test with large JSON file"
	@echo "  git-setup      - Set up git hooks and configuration"