# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
BINARY_NAME=reqparser
GOFILES=$(shell find . -type f -name '*.go')

# Tools
GOLINT=golangci-lint
GOSEC=gosec

.PHONY: all build test clean run lint fmt sec tidy coverage help

all: lint test build ## Run lint, test, and build

build: ## Build the binary
	$(GOBUILD) -o $(BINARY_NAME) -v

test: ## Run tests
	$(GOTEST) -v ./...

clean: ## Remove binary and test cache
	rm -f $(BINARY_NAME)
	$(GOCMD) clean
	rm -f coverage.out

run: build ## Build and run the binary
	./$(BINARY_NAME)

lint: ## Run linter
	@if ! command -v $(GOLINT) > /dev/null; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin; \
	fi
	$(GOLINT) run ./...

fmt: ## Format code
	$(GOFMT) ./...

sec: ## Run security check
	@if ! command -v $(GOSEC) > /dev/null; then \
		echo "Installing gosec..."; \
		$(GOGET) github.com/securego/gosec/v2/cmd/gosec; \
	fi
	$(GOSEC) ./...

tidy: ## Tidy up module dependencies
	$(GOMOD) tidy

coverage: ## Generate test coverage report
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

check: fmt lint sec test ## Run all checks (format, lint, security, tests)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Default target
.DEFAULT_GOAL := help
