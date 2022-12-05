.DEFAULT_GOAL = help

-include .env
export

run: ## run server
	@go run ./cmd/mongolia

build:

test: ## run tests
	@go test ./... -v

tidy: ## sync go.mod to source files
	@go mod tidy

help: ## print this help
	@grep -E '^[0-9a-zA-Z%_-]+:.*## .*$$' $(firstword $(MAKEFILE_LIST)) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: test
