XC_OS="linux darwin windows"
XC_ARCH="amd64 arm64"
BIN="./bin"
SRC=$(shell find . -name "*.go")
VERSION="r1.0.0"
CLI_NAME="ristioctl"
DATE_TIME=$(shell date)

LDFLAGS="-s -w -X 'google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore' -X 'main.version=$(VERSION)' -X 'main.time=$(DATE_TIME)'"

default: all

all: generate check vet 

generate:
	go generate ./...
	$(MAKE) check

.PHONY: tidy
tidy:
	go mod tidy
	go mod vendor

.PHONY: check
check:
	go fmt ./...
	go vet ./...

	$(MAKE) tidy


.PHONY: vet
vet:
	$(info ******************** vetting ********************)
	go vet ./...

.PHONY: clean
clean:
	$(info ******************** cleaning compiled binaries ********************)
	rm -rf $(BIN)


.PHONY: build
build:
	$(info ******************** compiling binaries ********************)
	gox \
		-osarch="linux/amd64" \
		-osarch="linux/arm64" \
		-osarch="darwin/amd64" \
		-osarch="darwin/arm64" \
		-osarch="windows/amd64" \
		-osarch="windows/arm64" \
		-output=$(BIN)/$(CLI_NAME)_{{.OS}}_{{.Arch}} \
		-ldflags=$(LDFLAGS) \
		;
