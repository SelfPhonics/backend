
##################
# Dependencies
##################
CACHE_BASE := .cache
UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)

CACHE ?= $(CACHE_BASE)/$(UNAME_OS)/$(UNAME_ARCH)
# The location where binary dependencies will be installed.
CACHE_BIN := $(CACHE)/bin
# Marker files are put into this directory to denote the current version of binaries that are installed.
CACHE_VERSIONS := $(CACHE)/versions

GOOSE_VERSION := v3.19.1

# GOOSE points to the marker file for the installed version.
#
# If GOOSE_VERSION is changed, the binary will be re-downloaded.
GOOSE := $(CACHE_VERSIONS)/gooose/$(GOOSE_VERSION)
$(GOOSE):
	GOBIN=$(abspath $(CACHE_BIN)) go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)
	@rm -rf $(dir $(GOOSE))
	@mkdir -p $(dir $(GOOSE))
	@touch $(GOOSE)

deps: $(GOOSE)

versions: deps
	@PATH=$(abspath $(CACHE_BIN)):$(GOPATH) goose --version

##################
# Running
##################

.PHONY: run
run:
	go run .

##################
# Building
##################

##################
# Publishing
##################
