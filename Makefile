.DEFAULT_TARGET: help

.PHONY: help
help: ## Display this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: infra
infra: ## Start necessary infrastructure
	@docker-compose up -d


.PHONY: run
run: ## Run application
	@go run cmd/main.go

.PHONY: lint
lint: ## Execute syntatic analysis in the code and autofix minor problems
	@golangci-lint run --fix

.PHONY: test
test: ## Execute the tests
	@go test ./...

.PHONY: coverage
coverage: ## Generate test coverage in the development environment
	@go test ./... -coverprofile=/tmp/coverage.out -coverpkg=./...
	@go tool cover -html=/tmp/coverage.out

.PHONY: ci
ci: lint test ## Execute the tests and lint commands