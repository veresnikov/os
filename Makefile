all: build test check

.PHONY: build
build: modules
	bin/go-build "cmd" "bin/converter" converter

.PHONY: modules
modules:
	go mod tidy

.PHONY: test
test:
	go test ./...

.PHONY: check
check:
	golangci-lint run
