VERSION := $(shell git describe --tags)
API_VERSION := "1.0.2"
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
	rm $(GOBIN)/ioncli-$(VERSION)-* 2> /dev/null
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(GOBIN)/ioncli-$(VERSION)-win64.exe $(GOFILES)
	zip $(GOBIN)/ioncli-$(VERSION)-win64.zip $(GOBIN)/ioncli-$(VERSION)-win64.exe
	GOOS=windows GOARCH=386 go build $(LDFLAGS) -o $(GOBIN)/ioncli-$(VERSION)-win32.exe $(GOFILES)
	zip $(GOBIN)/ioncli-$(VERSION)-win32.zip $(GOBIN)/ioncli-$(VERSION)-win32.exe
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(GOBIN)/ioncli-$(VERSION)-linux64 $(GOFILES)
	tar -czvf $(GOBIN)/ioncli-$(VERSION)-linux64.tar.gz $(GOBIN)/ioncli-$(VERSION)-linux64
	GOOS=linux GOARCH=386 go build $(LDFLAGS) -o $(GOBIN)/ioncli-$(VERSION)-linux32 $(GOFILES)
	tar -czvf $(GOBIN)/ioncli-$(VERSION)-linux32.tar.gz $(GOBIN)/ioncli-$(VERSION)-linux32
	GOOS=linux GOARCH=arm go build $(LDFLAGS) -o $(GOBIN)/ioncli-$(VERSION)-linux-arm32 $(GOFILES)
	tar -czvf $(GOBIN)/ioncli-$(VERSION)-linux-arm32.tar.gz $(GOBIN)/ioncli-$(VERSION)-linux-arm32
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(GOBIN)/ioncli-$(VERSION)-linux-arm64 $(GOFILES)
	tar -czvf $(GOBIN)/ioncli-$(VERSION)-linux-arm64.tar.gz $(GOBIN)/ioncli-$(VERSION)-linux-arm64
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(GOBIN)/ioncli-$(VERSION)-darwin64 $(GOFILES)
	tar -czvf $(GOBIN)/ioncli-$(VERSION)-darwin64.tar.gz $(GOBIN)/ioncli-$(VERSION)-darwin64

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo