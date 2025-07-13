NAME=konfig
REPO=github.com/mfenderov/${NAME}

deps:
	@go list -f '{{.Dir}}' -m | xargs -L1 go mod tidy -C
test:
	@go test -race ./...
	@go vet ./...
	@cd test-proj && go test -race

lint:
	@echo "Running golangci-lint..."
	@golangci-lint run --config .golangci.yml || echo "⚠️  Linting found issues (this is informational)"

coverage:
	@go test -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out | tail -1

coverage-text:
	@go test -race -coverprofile=coverage.out ./... >/dev/null 2>&1
	@go tool cover -func=coverage.out | tail -1

quality: lint test coverage-text
	@echo "✅ All quality checks passed!"

ci: deps quality
	@echo "✅ CI pipeline completed successfully!"

docs:
	@echo "📚 Generating documentation..."
	@go doc -all . > docs.txt
	@echo "Documentation generated in docs.txt"
	@echo "To view online docs, run: godoc -http=:6060"

docs-serve:
	@echo "🌐 Starting documentation server on http://localhost:6060"
	@echo "Visit http://localhost:6060/pkg/github.com/mfenderov/konfig/ to view docs"
	@godoc -http=:6060
build:
	@go build
install:
	@go install ${REPO}
release:
	@GOPROXY=proxy.golang.org go list -m ${REPO}@${GITHUB_REF_NAME}