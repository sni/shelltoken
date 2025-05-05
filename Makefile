#!/usr/bin/make -f

PROJECT=shelltoken
MAKE:=make
SHELL:=bash
GOVERSION:=$(shell \
    go version | \
    awk -F'go| ' '{ split($$5, a, /\./); printf ("%04d%04d", a[1], a[2]); exit; }' \
)
# also update  .github/workflows/citest.yml when changing minumum version
# find . -name go.mod
MINGOVERSION:=00010023
MINGOVERSIONSTR:=1.23
BUILD:=$(shell git rev-parse --short HEAD)
REVISION:=$(shell printf "%04d" $$( git rev-list --all --count))
# see https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md
# and https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
TOOLSFOLDER=$(shell pwd)/tools
export GOBIN := $(TOOLSFOLDER)
export PATH := $(GOBIN):$(PATH)

BUILD_FLAGS=-ldflags "-s -w -X main.Build=$(BUILD) -X main.Revision=$(REVISION)"
TEST_FLAGS=-timeout=5m $(BUILD_FLAGS)
GO=go

all: build

tools: | versioncheck
	set -e; for DEP in $(shell grep "_ " buildtools/tools.go | awk '{ print $$2 }' | grep -v go-spew); do \
		( cd buildtools && $(GO) install $$DEP@latest ) ; \
	done
	set -e; for DEP in $(shell grep "_ " buildtools/tools.go | awk '{ print $$2 }' | grep go-spew); do \
		( cd buildtools && $(GO) install $$DEP ) ; \
	done
	( cd buildtools && $(GO) mod tidy )

updatedeps: versioncheck
	$(MAKE) clean
	$(MAKE) tools
	$(GO) mod download
	$(GO) mod download
	$(MAKE) cleandeps

cleandeps:
	$(GO) mod tidy
	( cd buildtools && $(GO) mod tidy )

vendor:
	$(GO) mod download
	$(GO) mod tidy
	GOWORK=off $(GO) mod vendor

build:
	@echo "this is a library, run make test to run tests."

test: vendor
	$(GO) test -short -v $(TEST_FLAGS) .
	if grep -Irn TODO: *.go ; then exit 1; fi

# test with filter
testf: vendor
	$(GO) test -short -v $(TEST_FLAGS) . -run "$(filter-out $@,$(MAKECMDGOALS))" 2>&1 | grep -v "no test files" | grep -v "no tests to run" | grep -v "^PASS"

citest: tools vendor
	#
	# Checking gofmt errors
	#
	if [ $$(gofmt -s -l *.go | wc -l) -gt 0 ]; then \
		echo "found format errors in these files:"; \
		gofmt -s -l .; \
		exit 1; \
	fi
	#
	# Checking TODO items
	#
	if grep -rn TODO: *.go; then exit 1; fi
	#
	# Run other subtests
	#
	$(MAKE) golangci
	-$(MAKE) govulncheck
	$(MAKE) fmt
	#
	# Normal test cases
	#
	go test -v
	#
	# Benchmark tests
	#
	go test -v -bench=B\* -run=^$$ . -benchmem
	#
	# Race rondition tests
	#
	$(MAKE) racetest
	#
	# All CI tests successfull
	#

benchmark:
	$(GO) test -v -bench=B\* -run=^$$ -benchmem .

racetest:
	go test -race -v .

covertest:
	$(GO) test -v -coverprofile=cover.out .
	$(GO) tool cover -func=cover.out
	$(GO) tool cover -html=cover.out -o coverage.html

coverweb:
	$(GO) test -v -coverprofile=cover.out .
	$(GO) tool cover -html=cover.out

clean:
	rm -rf vendor
	rm -rf $(TOOLSFOLDER)

GOVET=$(GO) vet -all
fmt: tools
	$(GOVET) .
	gofmt -w -s *.go
	./tools/gofumpt -w *.go
	./tools/gci write *.go  --skip-generated
	goimports -w *.go

versioncheck:
	@[ $$( printf '%s\n' $(GOVERSION) $(MINGOVERSION) | sort | head -n 1 ) = $(MINGOVERSION) ] || { \
		echo "**** ERROR:"; \
		echo "**** shelltoken library requires at least golang version $(MINGOVERSIONSTR) or higher"; \
		echo "**** this is: $$(go version)"; \
		exit 1; \
	}

golangci: tools
	#
	# golangci combines a few static code analyzer
	# See https://github.com/golangci/golangci-lint
	#
	golangci-lint run *.go

govulncheck: tools
	govulncheck ./...

# just skip unknown make targets
.DEFAULT:
	@if [[ "$(MAKECMDGOALS)" =~ ^testf ]]; then \
		: ; \
	else \
		echo "unknown make target(s): $(MAKECMDGOALS)"; \
		exit 1; \
	fi
