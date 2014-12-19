PROJECT := "github.com/yungsang/talk2docker"

WORKSPACE := $(shell godep path)

GITCOMMIT := $(shell git rev-parse --short HEAD)

all: build

get:
	godep get ./...

fmt:
	go fmt -x ./...

test:
	godep go test ./...

build: fmt restore
	godep go build -v -ldflags "-X $(PROJECT)/version.GITCOMMIT '$(GITCOMMIT)'"

install: fmt restore
	godep go install -v -ldflags "-X $(PROJECT)/version.GITCOMMIT '$(GITCOMMIT)'"

uninstall:
	go clean -x -i

clean:
	go clean -x
	$(RM) -rf $(WORKSPACE)

save:
	godep save

update:
	godep update ...
	$(RM) -rf $(WORKSPACE)

restore:
	GOPATH=$(WORKSPACE) godep restore

.PHONY: all get fmt test build install uninstall clean save update restore
