export APP_CMD_NAME = statemachines

all: build test check

.PHONY: build
build: modules
	bin/go-build "cmd" "bin/$(APP_CMD_NAME)" $(APP_CMD_NAME)

.PHONY: modules
modules:
	go mod tidy

.PHONY: test
test:
	go test ./...

.PHONY: check
check:
	golangci-lint run
