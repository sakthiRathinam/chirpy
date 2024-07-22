# Define the binary paths if they are not in your PATH
AIR_BIN = $(shell which air)
GOFMT_BIN = $(shell which gofmt)

# Default target
.PHONY: all
all: help

# Help target to list available commands
.PHONY: help
help:
	@echo "Makefile for Go project"
	@echo
	@echo "Usage:"
	@echo "  make [target]"
	@echo
	@echo "Targets:"
	@echo "  watch          Run air to watch for file changes and rebuild"
	@echo "  lint           Run gofmt on all Go files to check formatting"
	@echo "  docker-watch   Run docker-compose up --build"
	@echo "  docker-build   Run docker-compose build"
	@echo "  help           Display this help message"

# Target to watch files with air
.PHONY: watch
watch:
	@if [ -z "$(AIR_BIN)" ]; then \
		echo "Air is not installed. Please install it first."; \
		exit 1; \
	fi
	@$(AIR_BIN)

# Target to lint Go files with gofmt
.PHONY: lint
lint:
	@if [ -z "$(GOFMT_BIN)" ]; then \
		echo "gofmt is not installed. Please install it first."; \
		exit 1; \
	fi
	@gofmt -l -w .

# Target to run docker-compose up --build
.PHONY: docker-watch
docker-watch:
	@docker-compose up --build

# Target to run docker-compose build
.PHONY: docker-build
docker-build:
	@docker-compose build

# Ensuring air and gofmt are installed
.PHONY: check-tools
check-tools:
	@command -v air >/dev/null 2>&1 || { echo >&2 "Air is not installed. Aborting."; exit 1; }
	@command -v gofmt >/dev/null 2>&1 || { echo >&2 "gofmt is not installed. Aborting."; exit 1; }
