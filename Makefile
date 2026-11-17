.DEFAULT_GOAL := help

.PHONY: build test lint fmt clean

help: ## Show available commands
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building..."
	go build -o bin/server ./cmd/server

test: ## Run tests
	@echo "Running tests..."
	go test -race -count=1 -v ./...

lint: ## Run golangci-lint
	@echo "Running linter..."
	golangci-lint run ./...

fmt: ## Go fmt
	@echo "Formatting..."
	go fmt ./...

clean: ## Remove go build artifacts
	@echo "Cleaning..."
	go clean -cache -testcache
	rm -rf cover.out bin/
