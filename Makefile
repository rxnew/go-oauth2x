SHELL=/bin/bash -o pipefail
.DEFAULT_GOAL := help

PROJ := go-oauth2x
PKG := github.com/rxnew/$(PROJ)

GO := go
GOBIN := $(abspath .bin)
GOTEST := $(or $(GOTEST),$(GO) test)

export PATH := $(GOBIN):${PATH}
export GO111MODULE := on

GO_TOOLS_SRC := tools.go
GO_TOOLS_PKG := $(shell cat $(notdir $(GO_TOOLS_SRC)) | awk -F'"' '/_/ {print $$2}')

# go install is responsible for not re-building when the code hasn't changed
# usage $(call make-go-tools-dependency,<GOPKG>)
define make-go-tools-dependency
$(GOBIN)/$(notdir $1): go.mod go.sum Makefile
	GOBIN=$(GOBIN) go install $1
endef
$(foreach pkg,$(GO_TOOLS_PKG),$(eval $(call make-go-tools-dependency,$(pkg))))

.PHONY: deps
deps: ## Download dependencies
	$(GO) mod download

.PHONY: test
test: unit-test ## Exec all tests

.PHONY: unit-test
unit-test: ## Exec unit tests
	$(GOTEST) -failfast -count=1 -race ./...

.PHONY: clean
clean: ## Clean bin and cache
	rm -rf ./bin/*
	$(GO) clean -testcache

.PHONY: fmt
fmt: $(GOBIN)/goimports
fmt: ## Format sources
	goimports -local $(PKG) -w .
	$(GO) mod tidy

.PHONY: go-generate
go-generate: $(GOBIN)/mockgen $(GOBIN)/goimports
go-generate: ## Exec go generate
	$(GO) generate ./...
	goimports -local $(PKG) -w .

.PHONY: lint
lint: golangci-lint ## Exec all linters

.PHONY: golangci-lint
golangci-lint: $(GOBIN)/golangci-lint
golangci-lint: ## Exec golangci-lint
	golangci-lint run ./...

.PHONY: install-tools
install-tools: $(GO_TOOLS_PKG) ## Install tools dependencies

# https://gist.github.com/tadashi-aikawa/da73d277a3c1ec6767ed48d1335900f3
.PHONY: $(shell grep -h -E '^[a-zA-Z_-]+:' $(MAKEFILE_LIST) | sed 's/://')

# https://postd.cc/auto-documented-makefile/
help: ## Show help message
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
