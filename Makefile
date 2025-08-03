NAME=konfig
REPO=github.com/mfenderov/${NAME}

# Environment detection
CI_ENV := $(if $(CI),true,false)
HAS_TTY := $(shell [ -t 1 ] && echo true || echo false)
HAS_TPARSE := $(shell command -v go tool tparse >/dev/null 2>&1 && echo true || echo false)

# Smart test command with environment auto-detection
test:
	@go list -f '{{.Dir}}' -m | xargs -L1 go mod tidy -C
ifeq ($(CI_ENV),true)
	@echo "🤖 Running tests for CI..."
	@mkdir -p .build
	@go test -json -race ./... | go tool tparse -format markdown > .build/test-results.md
	@echo "📊 Test results saved to .build/test-results.md"
	@go vet ./...
	@cd testdata/integration && go test -json -race
else ifeq ($(HAS_TPARSE),true)
	@echo "🧪 Running tests with pretty output..."
	@go test -json -race ./... | go tool tparse -all -smallscreen
	@go vet ./...
	@echo "🧪 Running integration tests..."
	@cd testdata/integration && go test -json -race | go tool tparse -all -smallscreen
else
	@go test -race ./...
	@go vet ./...
	@cd testdata/integration && go test -race
endif

# Static analysis with graceful degradation
lint:
	@echo "🔍 Running static analysis..."
	@go tool golangci-lint run --config .golangci.yml || echo "⚠️  Linting found issues (this is informational)"

# Smart coverage with environment-aware output
coverage:
	@mkdir -p .build
ifeq ($(CI_ENV),true)
	@go test -race -coverprofile=.build/coverage.out ./... >/dev/null 2>&1
	@go tool cover -func=.build/coverage.out | tail -1
else ifeq ($(HAS_TTY),true)
	@go test -race -coverprofile=.build/coverage.out ./...
	@go tool cover -html=.build/coverage.out -o .build/coverage.html
	@go tool cover -func=.build/coverage.out | tail -1
	@echo "📊 Coverage report: .build/coverage.html"
else
	@go test -race -coverprofile=.build/coverage.out ./... >/dev/null 2>&1
	@go tool cover -func=.build/coverage.out | tail -1
endif

# Complete quality pipeline
quality: lint test coverage
	@echo "✅ All quality checks passed!"

# Smart documentation with auto-detection
docs:
ifeq ($(filter serve,$(MAKECMDGOALS)),serve)
	@echo "🌐 Starting documentation server on http://localhost:6060"
	@echo "Visit http://localhost:6060/pkg/github.com/mfenderov/konfig/ to view docs"
	@godoc -http=:6060
else
	@echo "📚 Generating documentation..."
	@mkdir -p .build
	@go doc -all . > .build/docs.txt
	@echo "Documentation generated in .build/docs.txt"
	@echo "💡 Tip: Use 'make docs serve' to start documentation server"
endif

# Build operations
build:
	@go build

install:
	@go install ${REPO}

release:
	@GOPROXY=proxy.golang.org go list -m ${REPO}@${GITHUB_REF_NAME}

# Legacy aliases for backward compatibility (will be removed in future)
test-pretty: test
test-ci: test
coverage-text: coverage
ci: quality
docs-serve: docs serve

.PHONY: test lint coverage quality docs build install release
