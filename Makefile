PROJECT := "github.com/yungsang/talk2docker"

WORKSPACE := $(shell pwd)/Godeps/_workspace

GITCOMMIT := $(shell git rev-parse --short HEAD)

all: build

fmt:
	go fmt -x ./...

build: fmt dep-restore
	godep go build -v -ldflags "-X $(PROJECT)/version.GITCOMMIT '$(GITCOMMIT)'"

install: build
	godep go install -v -ldflags "-X $(PROJECT)/version.GITCOMMIT '$(GITCOMMIT)'"

uninstall:
	go clean -x -i

clean:
	go clean -x
	$(RM) -rf $(WORKSPACE)

dep-save:
	godep save

dep-update:
	godep update ...

dep-restore:
	GOPATH=$(WORKSPACE) godep restore

.PHONY: all fmt build install uninstall clean dep-save dep-update dep-restore
