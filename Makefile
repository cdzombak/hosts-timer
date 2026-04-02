SHELL:=/usr/bin/env bash

BIN_NAME:=hosts-timer
BIN_VERSION:=$(shell ./.version.sh)

default: help

# via https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## Print help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## Remove built products in ./out
	rm -rf ./out

.PHONY: build
build: clean ## Build (for the current platform & architecture) to ./out
	mkdir -p out
	go build -ldflags="-X main.version=${BIN_VERSION}" -o ./out/${BIN_NAME} .

.PHONY: install
install: ## Build & install hosts-timer to /usr/local/bin
	go build -ldflags="-X main.version=${BIN_VERSION}" -o /usr/local/bin/${BIN_NAME} .
