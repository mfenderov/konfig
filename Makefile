NAME=konfig
REPO=github.com/mfenderov/${NAME}

dependencies:
	@go mod vendor
	@go mod tidy
test:
	@go list -f '{{.Dir}}/...' -m | xargs go test -race ./...
	@go list -f '{{.Dir}}/...' -m | xargs go vet ./...
build:
	@go build
install:
	@go install ${REPO}
release:
	@GOPROXY=proxy.golang.org go list -m ${REPO}@${GITHUB_REF_NAME}