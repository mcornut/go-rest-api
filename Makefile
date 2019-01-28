.PHONY: test

PACKAGE := github.com/mcornut/go-rest-api

EXECUTABLE ?= go-rest-api

$(EXECUTABLE): $(shell find . -type f -print | grep -v vendor | grep "\.go")
	@echo "Building..."
	@go build

build: $(EXECUTABLE)

run: build
	@./$(EXECUTABLE) -config=local.toml

test:
	@go test -cover -coverprofile ./coverage.out ./... -count=1

cover: test
	@go tool cover -func ./coverage.out
	@go tool cover -html=./coverage.out