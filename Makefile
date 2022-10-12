help:
	@echo "targets:"
	@echo "    run          run server"
	@echo "    test         run all tests"
	@echo "    generate     regenerate all generated files"
	@echo "    tidy         tidy go mod"

# Load builder make environment
TOP?=$(shell cd .. && git rev-parse --show-toplevel)
-include $(TOP)/mk/base.mk

run:
	go run ./cmd/mongolia

test:
	go test -cover ./...

generate:
	@echo "generate..."
	go generate ./...

tidy:
	@echo "tidy..."
	go mod tidy
