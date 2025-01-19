NAME=konfig
REPO=github.com/mfenderov/${NAME}

dependencies:
	@go mod vendor
	@go mod tidy
clean:
	@rm -rf build/
test:
	@go test -race ./...
build:
	@go build -o build/ -v ./...
install:
	@go install ${REPO}
release:
	@GOPROXY=proxy.golang.org go list -m ${REPO}@${GITHUB_REF_NAME}