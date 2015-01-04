PROJECT := "github.com/ailispaw/talk2docker"

WORKSPACE := `godep path`

GIT_COMMIT := `git rev-parse --short HEAD`

ifeq ($(OS),Windows_NT)
	KERNEL_VERSION := ""
else
	KERNEL_VERSION := `uname -r`
endif

all: build

get:
	godep get ./...

fmt:
	go fmt -x ./...

test:
	godep go test ./...

build: fmt restore
	godep go build -v -ldflags "-X $(PROJECT)/version.GIT_COMMIT '$(GIT_COMMIT)' -X $(PROJECT)/version.KERNEL_VERSION '$(KERNEL_VERSION)'"

install: fmt restore uninstall
	godep go install -v -ldflags "-s -w -X $(PROJECT)/version.GIT_COMMIT '$(GIT_COMMIT)' -X $(PROJECT)/version.KERNEL_VERSION '$(KERNEL_VERSION)'"

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
