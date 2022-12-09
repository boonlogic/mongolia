.DEFAULT_GOAL = help

-include .env
export

.PHONY: test build format format-check go-check

UNAME := $(shell uname)

# load top-level makefile variables
TOP?=$(shell cd .. && git rev-parse --show-toplevel)
CWD=$(shell pwd)
-include $(TOP)/mk/base.mk

export GOOS
export DYLD_LIBRARY_PATH=$(INSTALL_ROOT)/lib
export LD_LIBRARY_PATH=$(INSTALL_ROOT)/lib

$(info GOPATH=$(GOPATH))
$(info INSTALL_ROOT=$(INSTALL_ROOT))

run: ## run server
	@go run ./cmd/mongolia

build: go-check
	@go build -a -modcacherw -o $(INSTALL_ROOT)/bin/mongolia cmd/mongolia/main.go

test: ## run tests
	@go test ./... -v

tidy: ## sync go.mod to source files
	@go mod tidy

help: ## print this help
	@grep -E '^[0-9a-zA-Z%_-]+:.*## .*$$' $(firstword $(MAKEFILE_LIST)) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

format: ## Run the formatter on go code
	@cd $(GOPATH)/src/$(DEVPATH)/expert/mongolia && go fmt ./...

format-check: format ## Run the formatter and perform diff (for pipeline)
	@git diff --exit-code; if [ $$? -ne 0 ]; then echo "format-check failed"; exit 1; fi; \
	echo "*** format-check passed"

go-check:
ifndef GOPATH
	$(error GOPATH is undefined)
endif
