.PHONY: test

# Load builder make environment
#TOP?=$(shell cd .. && git rev-parse --show-toplevel)
#-include $(TOP)/mk/base.mk

# Load local environment
-include .env
export

run: ## run server
	@go run ./cmd/mongolia

test: ## run tests
	@go test ./... -v

tidy: ## sync go.mod to source files
	@go mod tidy
