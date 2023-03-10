EXECUTABLE=gql-lint
LINUX=$(EXECUTABLE)_linux_amd64
DARWIN_AMD=$(EXECUTABLE)_darwin_amd64
DARWIN_ARM=$(EXECUTABLE)_darwin_arm64
UNIVERSAL=$(EXECUTABLE)_darwin_universal
VERSION=$(shell git describe --tags --always --long --dirty)

.PHONY: all test clean

all: test build ## Build and run tests

test: ## Run unit tests
	go test ./...

build: linux darwin-amd darwin-arm universal ## Build binaries
	@echo version: $(VERSION)

linux: $(LINUX) ## Build for Linux

darwin-amd: $(DARWIN_AMD) ## Build for Darwin AMD (intel macOS)

darwin-arm: $(DARWIN_ARM) ## Build for Darwin ARM (m1 macOS)

universal: darwin-amd darwin-arm $(UNIVERSAL) ## Build for Darwin Universal (intel and m1 macOS)

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build -v -o dist/$(LINUX) -ldflags="-s -w -X main.version=$(VERSION)"  ./cmd/gql-lint.go

$(DARWIN_AMD):
	env GOOS=darwin GOARCH=amd64 go build -v -o dist/$(DARWIN_AMD) -ldflags="-s -w -X main.version=$(VERSION)"  ./cmd/gql-lint.go

$(DARWIN_ARM):
	env GOOS=darwin GOARCH=arm64 go build -v -o dist/$(DARWIN_ARM) -ldflags="-s -w -X main.version=$(VERSION)"  ./cmd/gql-lint.go

$(UNIVERSAL):
	lipo -create dist/$(DARWIN_AMD) dist/$(DARWIN_ARM) -output dist/$(UNIVERSAL)

clean: ## Remove previous build
	rm -rf dist

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
