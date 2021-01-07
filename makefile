SHELL=bash

BUILD=build
BIN_DIR?=.

.PHONY: build
build:
	@mkdir -p $(BUILD)/$(BIN_DIR)
	go build -o $(BUILD)/$(BIN_DIR)/books-api main.go book.go errors.go

.PHONY: debug
debug: build
	HUMAN_LOG=1 go run -race main.go book.go errors.go

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: convey
convey:
	goconvey ./...
