.PHONY: help build run test lint tidy
.DEFAULT_GOAL := help

BINARY := knowledge-mcp

help:
	@echo "Usage: make <target> [ROOT=<dir>]"
	@echo ""
	@echo "Targets:"
	@echo "  build        Build the binary (output: ./$(BINARY))"
	@echo "  run          Run the server (pass ROOT=<dir> for knowledge store path)"
	@echo "  test         Run all tests"
	@echo "  lint         Run golangci-lint"
	@echo "  deps         Update dependencies"
	@echo "  tidy         Tidy go.mod"

build:
	go build -o $(BINARY) ./cmd/app

run:
	go run ./cmd/app $(ROOT)

test:
	go test ./...

lint:
	golangci-lint run

deps:
	go get -u ./...

tidy:
	go mod tidy