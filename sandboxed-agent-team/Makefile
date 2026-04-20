# Makefile for sandboxed-agent-team
#
# Targets:
#   make build    - build the binary for the current host
#   make release  - cross-compile binaries for all supported platforms
#   make test     - run unit tests
#   make vet      - go vet
#   make check    - vet + test + build (what CI would run)
#   make clean    - remove build output

BINARY      := sandboxed-agent-team
RELEASE_DIR := release

.PHONY: build
build:
	go build -o $(BINARY) .

.PHONY: release
release:
	rm -rf $(RELEASE_DIR)
	mkdir -p $(RELEASE_DIR)
	GOOS=darwin  GOARCH=arm64 go build -o $(RELEASE_DIR)/$(BINARY)-darwin-arm64  .
	GOOS=darwin  GOARCH=amd64 go build -o $(RELEASE_DIR)/$(BINARY)-darwin-amd64  .
	GOOS=linux   GOARCH=amd64 go build -o $(RELEASE_DIR)/$(BINARY)-linux-amd64   .
	GOOS=windows GOARCH=amd64 go build -o $(RELEASE_DIR)/$(BINARY)-windows-amd64.exe .
	@echo "--- release binaries ---"
	@ls -1 $(RELEASE_DIR)

.PHONY: test
test:
	go test ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: check
check: vet test build

.PHONY: clean
clean:
	rm -f $(BINARY)
	rm -rf $(RELEASE_DIR)
