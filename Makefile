SHELL = /bin/bash

IMPORT_PATH = github.com/rwstauner/ynetd

VERSION_VAR = main.Version
VERSION = $(shell git describe --tags --long --always --match 'v[0-9]*' | sed -e 's/-/./')
BUILD_ARGS = -tags netgo -ldflags '-w -extldflags "-static" -X $(VERSION_VAR)=$(VERSION)'

TESTS ?= test

SRC_VOL = /go/src/$(IMPORT_PATH)
DOCKER_IMAGE = rwstauner/golang-release
DOCKER_VARS  = -e TESTS="$(TESTS)"
DOCKER_RUN   = docker run --rm -v $(PWD)/tmp/gosrc:/go/src -v $(PWD):$(SRC_VOL) -w $(SRC_VOL) $(DOCKER_VARS) $(DOCKER_IMAGE)
DOCKER_MAKE  = $(DOCKER_RUN) make

.PHONY: all build test
all: build

build:
	$(DOCKER_MAKE) _build
_build: _test
	gox $(BUILD_ARGS) -output "build/ynetd-{{.OS}}-{{.Arch}}/ynetd" -verbose
	(cd build && for i in ynetd-*/; do (cd $$i && zip "../$${i%/}.zip" ynetd*); done)

_test_build:
	go get
	go test
	go build $(BUILD_ARGS) -o build/ynetd
	go build -o build/ytester test/ytester.go

test:
	$(DOCKER_MAKE) _test
_test: _test_build
	YNETD=build/ynetd YTESTER=build/ytester bats $(TESTS)
