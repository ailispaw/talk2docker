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

	registry, name, tag, err = ParseRepositoryName("yungsang/busybox")
	expectedRegistry = ""
	expectedName = "yungsang/busybox"
	expectedTag = "latest"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expectedRegistry) || (name != expectedName) || (tag != expectedTag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expectedRegistry, name, expectedName, tag, expectedTag)
	}

	registry, name, tag, err = ParseRepositoryName("yungsang/busybox:latest")
	expectedRegistry = ""
	expectedName = "yungsang/busybox"
	expectedTag = "latest"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expectedRegistry) || (name != expectedName) || (tag != expectedTag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expectedRegistry, name, expectedName, tag, expectedTag)
	}

	registry, name, tag, err = ParseRepositoryName("localhost/yungsang/busybox:tagname")
	expectedRegistry = "localhost"
	expectedName = "yungsang/busybox"
	expectedTag = "tagname"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expectedRegistry) || (name != expectedName) || (tag != expectedTag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expectedRegistry, name, expectedName, tag, expectedTag)
	}

	registry, name, tag, err = ParseRepositoryName("localhost:5000/yungsang/busybox")
	expectedRegistry = "localhost:5000"
	expectedName = "yungsang/busybox"
	expectedTag = "latest"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expectedRegistry) || (name != expectedName) || (tag != expectedTag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expectedRegistry, name, expectedName, tag, expectedTag)
	}

	registry, name, tag, err = ParseRepositoryName("localhost:5000/yungsang/busybox:tagname")
	expectedRegistry = "localhost:5000"
	expectedName = "yungsang/busybox"
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

	registry, name, tag, err = ParseRepositoryName("192.168.33.201:5000/yungsang/flannel")
	expectedRegistry = "192.168.33.201:5000"
	expectedName = "yungsang/flannel"
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
