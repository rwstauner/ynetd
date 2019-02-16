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

GPG_TO_SIGN = "$$GPG_PREFIX/gpg" --homedir "$$GPG_PREFIX/home" --batch
GPG_TO_VERIFY = gpg --homedir "$$GPG_PREFIX/verify" --batch
GPG_SIGN_OPTS = --yes --armor --detach-sign
PGP_KEYSERVER=hkp://ha.pool.sks-keyservers.net
PGP_FINGERPRINT=9791707D75D1474B6936CA216AD6ED6EA9371AED

CLI_DEP = if ! which $(1); then go get $(2); fi

PACKAGE_DIRS = $(shell find -name '*_test.go' -a -not -path './tmp/*' | sed -r 's,[^/]+$$,,' | sort | uniq)

.PHONY: all build test
all: build

dbuild:
	$(DOCKER_MAKE) build

build: clean test build_release

build_one: ytester
	go build $(BUILD_ARGS) -o build/ynetd

build_all: deps
	gox $(BUILD_ARGS) -output "build/ynetd-{{.OS}}-{{.Arch}}/ynetd" -verbose
	(cd build && for i in ynetd-*/; do (cd "$$i" && zip "../$${i%/}.zip" ynetd*) && rm -rf "$$i"; done)

build_release: build_all sign_builds

sign_builds:
	chmod 0700 .sign/home .sign/verify
	set -e; GPG_PREFIX=.sign; \
	if [[ -n "$$GPG_PASSPHRASE" ]] && [[ -r "$$GPG_PREFIX/key" ]]; then \
		$(GPG_TO_SIGN) --import "$$GPG_PREFIX/key"; \
		$(GPG_TO_VERIFY) --keyserver "$(PGP_KEYSERVER)" --recv-keys "$(PGP_FINGERPRINT)"; \
		for i in build/ynetd-*.zip; do \
			echo "$$GPG_PASSPHRASE" | $(GPG_TO_SIGN) --passphrase-fd 0 $(GPG_SIGN_OPTS) "$$i"; \
			$(GPG_TO_VERIFY) --verify "$$i.asc" "$$i"; \
		done; \
	fi

clean:
	rm -rf build

install:
	go install $(BUILD_ARGS)

gotest:
	go test $(PACKAGE_DIRS)

deps:
	$(call CLI_DEP,gox,github.com/mitchellh/gox)
	$(call CLI_DEP,golint,golang.org/x/lint/golint)

build_test: deps ytester
	go get
	golint -set_exit_status $(PACKAGE_DIRS)
	go vet $(PACKAGE_DIRS)
	go test $(PACKAGE_DIRS)
	go build $(TEST_BUILD_ARGS) $(BUILD_ARGS) -o build/ynetd

ytester:
	go build -o build/ytester test/ytester.go

dtest:
	$(DOCKER_MAKE) test
test: build_test
	YNETD=build/ynetd YTESTER=build/ytester bats $(TESTS)
