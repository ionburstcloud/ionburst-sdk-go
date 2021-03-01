VERSION := $(shell cat build_info | grep RELEASE_VERSION | cut -d "=" -f 2)
API_VERSION := $(shell cat build_info | grep RELEASE_VERSION | cut -d "=" -f 2)
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/vendor:$(GOBASE)
GOBIN := $(GOBASE)/bin
GOFILES := "./ioncli"

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -X=main.APIVersion=$(API_VERSION)"

## install: Install missing dependencies. Runs `go get` internally. e.g; make install get=github.com/foo/bar
deps:
	@echo "  >  Checking if there is any missing dependencies..."
	go mod vendor

## exec: Run given command, wrapped with custom GOPATH. e.g; make exec run="go test ./..."
exec:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GO111MODULE=on $(run)

## clean: Clean build files. Runs `go clean` internally.
clean:
	@-rm $(GOBIN)/ioncli 2> /dev/null
	@-$(MAKE) go-clean

compile: go-get go-build

go-build:
	@echo "  >  Building binary..."
	go build $(LDFLAGS) -o $(GOBIN)/ioncli $(GOFILES)

go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	go get $(get)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

go-cross-compile:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(GOBIN)/ioncli.exe $(GOFILES)
	tar cvfj ioncli-$(VERSION)-win64.tar.bz2 -C $(GOBIN) ioncli.exe
	rm -f $(GOBIN)/ioncli.exe
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(GOBIN)/ioncli $(GOFILES)
	tar cvfj ioncli-$(VERSION)-linux64.tar.bz2 -C $(GOBIN) ioncli
	rm -f $(GOBIN)/ioncli
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(GOBIN)/ioncli $(GOFILES)
	tar cvfj ioncli-$(VERSION)-linux-arm64.tar.bz2 -C bin ioncli
	rm -f $(GOBIN)/ioncli
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(GOBIN)/ioncli $(GOFILES)
	tar cvfj ioncli-$(VERSION)-darwin64.tar.bz2 -C bin ioncli
	rm -f $(GOBIN)/ioncli
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(GOBIN)/ioncli $(GOFILES)
	tar cvfj ioncli-$(VERSION)-darwin-arm64.tar.bz2 -C bin ioncli
	rm -f $(GOBIN)/ioncli

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command to run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo