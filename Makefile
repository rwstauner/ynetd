SHELL = /bin/bash

IMPORT_PATH = github.com/rwstauner/ynetd

VERSION_VAR = main.Version
VERSION = $(shell git describe --tags --long --always --match 'v[0-9]*' | sed -e 's/-/./')
BUILD_ARGS = -tags netgo -ldflags '-w -extldflags "-static" -X $(VERSION_VAR)=$(VERSION)'

TEST_BUILD_ARGS = $(shell test -s /etc/alpine-release || echo "-race")

TESTS ?= test

SRC_VOL = /go/src/$(IMPORT_PATH)
DOCKER_IMAGE = rwstauner/golang-release
DOCKER_VARS  = -e TESTS="$(TESTS)"
DOCKER_RUN   = docker run --rm -v $(PWD)/tmp/gosrc:/go/src -v $(PWD):$(SRC_VOL) -w $(SRC_VOL) $(DOCKER_VARS) $(DOCKER_IMAGE)
DOCKER_MAKE  = $(DOCKER_RUN) make

CLI_DEP = if ! which $(1); then go get $(2); fi

PACKAGE_DIRS = $(shell find -name '*_test.go' -a -not -path './tmp/*' | sed -r 's,[^/]+$$,,' | sort | uniq)

.PHONY: all build test
all: build

dbuild:
	$(DOCKER_MAKE) build

build: clean test build_frd

build_one: ytester
	go build $(BUILD_ARGS) -o build/ynetd

build_frd: deps
	gox $(BUILD_ARGS) -output "build/ynetd-{{.OS}}-{{.Arch}}/ynetd" -verbose
	(cd build && for i in ynetd-*/; do (cd $$i && zip "../$${i%/}.zip" ynetd*); done)

clean:
	rm -rf build

install:
	go install $(BUILD_ARGS)

gotest:
	go test $(PACKAGE_DIRS)

deps:
	$(call CLI_DEP,gox,github.com/mitchellh/gox)
	$(call CLI_DEP,golint,golang.org/x/lint/golint)

test_build: deps ytester
	go get
	golint -set_exit_status $(PACKAGE_DIRS)
	go vet $(PACKAGE_DIRS)
	go test $(PACKAGE_DIRS)
	go build $(TEST_BUILD_ARGS) $(BUILD_ARGS) -o build/ynetd

ytester:
	go build -o build/ytester test/ytester.go

dtest:
	$(DOCKER_MAKE) test
test: test_build
	YNETD=build/ynetd YTESTER=build/ytester bats $(TESTS)
