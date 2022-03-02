BUILD_DIR ?= build

all: lint cover build
.PHONY: all

lint:
	golangci-lint run ./...
.PHONY: lint

test:
	go test ./...
.PHONY: test

cover:
	go test -coverprofile=coverage.out ./...
.PHONY: cover

coverhtml: cover
	go tool cover -html=coverage.out
.PHONY: coverhtml

build: $(BUILD_DIR)/unmarshal

$(BUILD_DIR)/%: $(shell find . -type f -name "*.go")
	go build -trimpath -o $@ ./cmd/$*
