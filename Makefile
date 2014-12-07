WORKSPACE := $(shell pwd)/Godeps/_workspace

all: build

fmt:
	go fmt -x ./...

build: fmt dep-restore
	godep go build -v

install: build
	godep go install -v

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
