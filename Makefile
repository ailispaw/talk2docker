VERSION := 1.3.2

PROJECT := github.com/ailispaw/talk2docker

WORKSPACE := `godep path`

GIT_COMMIT := `git rev-parse --short HEAD`

all: build

get:
	godep get ./...

fmt:
	go fmt -x ./...

test:
	godep go test ./...

build: fmt restore
	CGO_ENABLED=0 godep go build -a -installsuffix cgo -v -ldflags "-X $(PROJECT)/version.GIT_COMMIT '$(GIT_COMMIT)' -X $(PROJECT)/version.APP_VERSION '$(VERSION)'"

install: uninstall build
	cp talk2docker $(GOPATH)/bin

uninstall:
	go clean -x -i

clean:
	go clean -x
	$(RM) -rf "$(WORKSPACE)"

save:
	godep save

update:
	godep update ...
	$(RM) -rf "$(WORKSPACE)"

restore:
	GOPATH="$(WORKSPACE)" godep restore

.PHONY: all get fmt test build install uninstall clean save update restore

xc:
	$(RM) -r bin/$(VERSION)
	vagrant up --no-provision
	vagrant provision
	vagrant suspend

goxc:
	CGO_ENABLED=0 goxc -d="bin" -bc="darwin linux,!arm" -build-installsuffix="cgo" -build-ldflags="-X $(PROJECT)/version.GIT_COMMIT '$(GIT_COMMIT)' -X $(PROJECT)/version.APP_VERSION '$(VERSION)'" -pv=$(VERSION) xc archive

.PHONY: xc goxc
