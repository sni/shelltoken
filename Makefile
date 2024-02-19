#!/usr/bin/make -f

MAKE:=make
SHELL:=bash
GOVERSION:=$(shell \
    go version | \
    awk -F'go| ' '{ split($$5, a, /\./); printf ("%04d%04d", a[1], a[2]); exit; }' \
)
MINGOVERSION:=00010021
MINGOVERSIONSTR:=1.21
# see https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md
# and https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
TOOLSFOLDER=$(shell pwd)/tools
export GOBIN := $(TOOLSFOLDER)
export PATH := $(GOBIN):$(PATH)
GO=go

all: build

tools: | versioncheck vendor
	$(GO) mod download
	set -e; for DEP in $(shell grep "_ " buildtools/tools.go | awk '{ print $$2 }'); do \
		( cd buildtools && $(GO) install $$DEP ) ; \
	done
	$(GO) mod tidy
	( cd buildtools && $(GO) mod tidy )
	# pin these dependencies
	( cd buildtools && $(GO) get github.com/golangci/golangci-lint@latest )
	$(GO) mod vendor

updatedeps: versioncheck
	$(MAKE) clean
	$(MAKE) tools
	$(GO) mod download
	set -e; for DEP in $(shell grep "_ " buildtools/tools.go | awk '{ print $$2 }'); do \
		( cd buildtools && $(GO) get $$DEP ) ; \
	done
	( cd buildtools && $(GO) get -u )
	$(GO) mod download
	$(GO) get -u
	$(GO) get -t -u
	$(MAKE) cleandeps

cleandeps:
	$(GO) mod tidy

vendor:
	$(GO) mod download
	$(GO) mod tidy
	$(GO) mod vendor

build:
	@echo "this is a library, run make test to run tests."

test: vendor
	$(GO) test -v .
	if grep -Irn TODO: *.go ; then exit 1; fi

# test with filter
testf: vendor
	$(GO) test -v . -run "$(filter-out $@,$(MAKECMDGOALS))" 2>&1 | grep -v "no test files" | grep -v "no tests to run" | grep -v "^PASS"

citest: vendor
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
	rm -rf \
		vendor \
		tools \
		cover.out \

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

golangci:
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
