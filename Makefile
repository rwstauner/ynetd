SHELL = /bin/bash

DOCKER_IMAGE = rwstauner/golang-release
DOCKER_RUN   = docker run --rm -v $(PWD):/src -w /src $(DOCKER_IMAGE)
DOCKER_MAKE  = $(DOCKER_RUN) make
BUILD_ARGS   = -tags netgo -ldflags '-w -extldflags "-static"'

.PHONY: all test
all: test

test:
	$(DOCKER_MAKE) _test
_test:
	go test
	go build $(BUILD_ARGS) -o build/ynetd
	YNETD=build/ynetd bats test
