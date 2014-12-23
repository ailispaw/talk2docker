package client

import (
	"testing"
)

func TestParseRepositoryName(t *testing.T) {
	registry, name, tag, err := ParseRepositoryName("busybox")
	var (
		expected_registry = ""
		expected_name     = "busybox"
		expected_tag      = "latest"
	)
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expected_registry) || (name != expected_name) || (tag != expected_tag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expected_registry, name, expected_name, tag, expected_tag)
	}

	registry, name, tag, err = ParseRepositoryName("yungsang/busybox")
	expected_registry = ""
	expected_name = "yungsang/busybox"
	expected_tag = "latest"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expected_registry) || (name != expected_name) || (tag != expected_tag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expected_registry, name, expected_name, tag, expected_tag)
	}

	registry, name, tag, err = ParseRepositoryName("yungsang/busybox:latest")
	expected_registry = ""
	expected_name = "yungsang/busybox"
	expected_tag = "latest"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expected_registry) || (name != expected_name) || (tag != expected_tag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expected_registry, name, expected_name, tag, expected_tag)
	}

	registry, name, tag, err = ParseRepositoryName("localhost/yungsang/busybox:tagname")
	expected_registry = "localhost"
	expected_name = "yungsang/busybox"
	expected_tag = "tagname"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expected_registry) || (name != expected_name) || (tag != expected_tag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expected_registry, name, expected_name, tag, expected_tag)
	}

	registry, name, tag, err = ParseRepositoryName("localhost:5000/yungsang/busybox")
	expected_registry = "localhost:5000"
	expected_name = "yungsang/busybox"
	expected_tag = "latest"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expected_registry) || (name != expected_name) || (tag != expected_tag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expected_registry, name, expected_name, tag, expected_tag)
	}

	registry, name, tag, err = ParseRepositoryName("localhost:5000/yungsang/busybox:tagname")
	expected_registry = "localhost:5000"
	expected_name = "yungsang/busybox"
	expected_tag = "tagname"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expected_registry) || (name != expected_name) || (tag != expected_tag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expected_registry, name, expected_name, tag, expected_tag)
	}

	registry, name, tag, err = ParseRepositoryName("quay.io/flannel")
	expected_registry = "quay.io"
	expected_name = "flannel"
	expected_tag = "latest"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expected_registry) || (name != expected_name) || (tag != expected_tag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expected_registry, name, expected_name, tag, expected_tag)
	}

	registry, name, tag, err = ParseRepositoryName("192.168.33.201:5000/yungsang/flannel")
	expected_registry = "192.168.33.201:5000"
	expected_name = "yungsang/flannel"
	expected_tag = "latest"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expected_registry) || (name != expected_name) || (tag != expected_tag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expected_registry, name, expected_name, tag, expected_tag)
	}

	registry, name, tag, err = ParseRepositoryName("index.docker.io/busybox")
	expected_registry = ""
	expected_name = "busybox"
	expected_tag = "latest"
	if err != nil {
		t.Errorf("%v", err)
	}
	if (registry != expected_registry) || (name != expected_name) || (tag != expected_tag) {
		t.Errorf("got %v\nwant %v,\n %v\nwant %v,\n %v\nwant %v",
			registry, expected_registry, name, expected_name, tag, expected_tag)
	}

	registry, name, tag, err = ParseRepositoryName("https://index.docker.io/v1/busybox")
	if err == nil {
		t.Errorf("%v", "This should be an error.")
	}
}
