all: lint cover
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
