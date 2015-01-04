package client

import (
	"testing"
)

func TestParseRepositoryName(t *testing.T) {
	registry, name, tag, err := ParseRepositoryName("busybox")
	var (
		expectedRegistry = ""
		expectedName     = "busybox"
		expectedTag      = "latest"
	)
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expectedRegistry) || (name != expectedName) || (tag != expectedTag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expectedRegistry, name, expectedName, tag, expectedTag)
	}

	registry, name, tag, err = ParseRepositoryName("ailispaw/busybox")
	expectedRegistry = ""
	expectedName = "ailispaw/busybox"
	expectedTag = "latest"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expectedRegistry) || (name != expectedName) || (tag != expectedTag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expectedRegistry, name, expectedName, tag, expectedTag)
	}

	registry, name, tag, err = ParseRepositoryName("ailispaw/busybox:latest")
	expectedRegistry = ""
	expectedName = "ailispaw/busybox"
	expectedTag = "latest"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expectedRegistry) || (name != expectedName) || (tag != expectedTag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expectedRegistry, name, expectedName, tag, expectedTag)
	}

	registry, name, tag, err = ParseRepositoryName("localhost/ailispaw/busybox:tagname")
	expectedRegistry = "localhost"
	expectedName = "ailispaw/busybox"
	expectedTag = "tagname"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expectedRegistry) || (name != expectedName) || (tag != expectedTag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expectedRegistry, name, expectedName, tag, expectedTag)
	}

	registry, name, tag, err = ParseRepositoryName("localhost:5000/ailispaw/busybox")
	expectedRegistry = "localhost:5000"
	expectedName = "ailispaw/busybox"
	expectedTag = "latest"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expectedRegistry) || (name != expectedName) || (tag != expectedTag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expectedRegistry, name, expectedName, tag, expectedTag)
	}

	registry, name, tag, err = ParseRepositoryName("localhost:5000/ailispaw/busybox:tagname")
	expectedRegistry = "localhost:5000"
	expectedName = "ailispaw/busybox"
	expectedTag = "tagname"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expectedRegistry) || (name != expectedName) || (tag != expectedTag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expectedRegistry, name, expectedName, tag, expectedTag)
	}

	registry, name, tag, err = ParseRepositoryName("quay.io/flannel")
	expectedRegistry = "quay.io"
	expectedName = "flannel"
	expectedTag = "latest"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expectedRegistry) || (name != expectedName) || (tag != expectedTag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expectedRegistry, name, expectedName, tag, expectedTag)
	}

	registry, name, tag, err = ParseRepositoryName("192.168.33.201:5000/ailispaw/flannel")
	expectedRegistry = "192.168.33.201:5000"
	expectedName = "ailispaw/flannel"
	expectedTag = "latest"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expectedRegistry) || (name != expectedName) || (tag != expectedTag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expectedRegistry, name, expectedName, tag, expectedTag)
	}

	registry, name, tag, err = ParseRepositoryName("index.docker.io/busybox")
	expectedRegistry = ""
	expectedName = "busybox"
	expectedTag = "latest"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expectedRegistry) || (name != expectedName) || (tag != expectedTag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expectedRegistry, name, expectedName, tag, expectedTag)
	}

	registry, name, tag, err = ParseRepositoryName("https://index.docker.io/v1/busybox")
	if err == nil {
		t.Errorf("%v", "This should be an error.")
	}
}
