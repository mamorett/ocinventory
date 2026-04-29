BINARY   := ocinventory
MODULE   := github.com/mamorett/ocinventory
CMD      := ./cmd/ocinventory
DIST     := dist
VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS  := -ldflags "-X main.version=$(VERSION) -s -w"

.PHONY: all build tidy clean linux-amd64 linux-arm64 darwin-arm64

## all: build for all three targets
all: linux-amd64 linux-arm64 darwin-arm64

## build: same as all
build: all

## linux-amd64: build for Linux x86-64
linux-amd64:
	@mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build $(LDFLAGS) -o $(DIST)/$(BINARY)-linux-amd64 $(CMD)
	@echo "→ $(DIST)/$(BINARY)-linux-amd64"

## linux-arm64: build for Linux ARM64 (e.g. AWS Graviton, OCI Ampere A1)
linux-arm64:
	@mkdir -p $(DIST)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
		go build $(LDFLAGS) -o $(DIST)/$(BINARY)-linux-arm64 $(CMD)
	@echo "→ $(DIST)/$(BINARY)-linux-arm64"

## darwin-arm64: build for macOS Apple Silicon
darwin-arm64:
	@mkdir -p $(DIST)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 \
		go build $(LDFLAGS) -o $(DIST)/$(BINARY)-darwin-arm64 $(CMD)
	@echo "→ $(DIST)/$(BINARY)-darwin-arm64"

## tidy: download dependencies and tidy go.sum
tidy:
	go mod tidy

## clean: remove compiled binaries
clean:
	rm -rf $(DIST)

## help: list available targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## //'
