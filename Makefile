SHELL = /bin/bash

IMPORT_PATH = github.com/rwstauner/ynetd

VERSION_VAR = main.Version
VERSION = $(shell git describe --tags --long --always --match 'v[0-9]*' | sed -e 's/-/./')
BUILD_ARGS = -tags netgo -ldflags '-w -extldflags "-static" -X $(VERSION_VAR)=$(VERSION)'

SRC_VOL = /go/src/$(IMPORT_PATH)
DOCKER_IMAGE = rwstauner/golang-release
DOCKER_RUN   = docker run --rm -v $(PWD):$(SRC_VOL) -w $(SRC_VOL) $(DOCKER_IMAGE)
DOCKER_MAKE  = $(DOCKER_RUN) make

.PHONY: all test
all: test

test:
	$(DOCKER_MAKE) _test
_test:
	go test
	go build $(BUILD_ARGS) -o build/ynetd
	YNETD=build/ynetd bats test
