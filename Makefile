all: build test check

.PHONY: build
build: modules
	bin/go-build "cmd/converter" "bin/converter" converter
	bin/go-build "cmd/executor" "bin/executor" executor

.PHONY: modules
modules:
	go mod tidy

.PHONY: test
test:
	go test ./...

.PHONY: check
check:
	golangci-lint run
