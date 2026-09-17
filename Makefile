SHELL:=/bin/bash

GO = go
GO_FMT = $(GO) tool gofumpt
GO_IMPORTS = $(GO) tool goimports
STATIC_CHECK = $(GO) tool staticcheck

GO_VET_OPTS = -v
GO_TEST_OPTS = -v -race
GO_FMT_OPTS = -l -w
GO_IMPORTS_OPTS = -w -local github.com/chez-shanpu/typesafeai

.PHONY: fmt
fmt:
	$(GO_FMT) $(GO_FMT_OPTS) .
	$(GO_IMPORTS) $(GO_IMPORTS_OPTS) .

.PHONY: fix
fix:
	$(GO) fix ./...

.PHONY: mod
mod:
	$(GO) mod tidy

.PHONY: check-diff
verify-diff: mod fmt fix
	git diff --exit-code --name-only

.PHONY: vet
vet:
	$(GO) vet $(GO_VET_OPTS) ./...

.PHONY: test
test:
	$(STATIC_CHECK) ./...
	$(GO) test $(GO_TEST_OPTS) ./...

.PHONY: clean
clean:
	-$(GO) clean

.PHONY: check
verify: vet verify-diff test

.PHONY: all
all: verify

.DEFAULT_GOAL=all