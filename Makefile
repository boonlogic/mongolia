help:
	@echo "targets:"
	@echo "    run          run server"
	@echo "    test         run all tests"
	@echo "    generate     regenerate all generated files"
	@echo "    tidy         tidy go mod"

# Load builder make environment
TOP?=$(shell cd .. && git rev-parse --show-toplevel)
-include $(TOP)/mk/base.mk

# Load local environment
-include .env
export

run:
	go run ./cmd/mongodm

test:
	go test -cover ./...

tidy:
	@echo "tidy..."
	go mod tidy
