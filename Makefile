DIST := $(PWD)/dist/

default: build

.PHONY: mocks
mocks:
	mockery --config=mockery.yaml

.PHONY: build
build:
	go build -o $(DIST) ./cmd/metrics/...

.PHONY: test
test:
	go test -v -count=1 ./...
