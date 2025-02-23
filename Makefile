BENCH_DIR ?= bench
BUILD_DIR ?= build
SRC_FILES := $(shell find . -name "*.go")

all: lint cover build
.PHONY: all

bench: $(BENCH_DIR)/Marshal.txt $(BENCH_DIR)/Unmarshal.txt
.PHONY: bench

$(BENCH_DIR)/%.txt: $(SRC_FILES)
	@mkdir -p $(BENCH_DIR)
	go test \
		-bench=Benchmark$* \
		-run=^$$ \
		-benchmem \
		-cpuprofile=$(BENCH_DIR)/$*.cpuprof \
		-memprofile=$(BENCH_DIR)/$*.memprof \
		-trace=$(BENCH_DIR)/$*.trace \
	| tee $@
	@echo
	benchstat $@
	@echo
	@echo "View benchmark profiles with:"
	@echo "  go tool pprof -http=:8081 $(BENCH_DIR)/$*.cpuprof"
	@echo "  go tool pprof -http=:8082 $(BENCH_DIR)/$*.memprof"

build: \
	$(BUILD_DIR)/roundtrip \
	$(BUILD_DIR)/unmarshal

$(BUILD_DIR)/%: $(SRC_FILES)
	go build -trimpath -o $@ ./cmd/$*

lint:
	golangci-lint run ./...
.PHONY: lint

cover:
	go test -coverprofile=coverage.out ./...
.PHONY: cover

coverhtml: cover
	go tool cover -html=coverage.out
.PHONY: coverhtml

test:
	go test ./...
.PHONY: test

file2fuzz:
	go tool file2fuzz -o testdata/fuzz/Fuzz/ testdata/input/*
.PHONY: file2fuzz

fuzz:
	go test -fuzz Fuzz
.PHONY: fuzz
