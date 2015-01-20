package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

func (client *DockerClient) ListContainers(all, size bool, limit int, since, before string, filters map[string][]string) ([]Container, error) {
	v := url.Values{}
	if all {
		v.Set("all", "1")
	}
	if size {
		v.Set("size", "1")
	}
	if limit > 0 {
		v.Set("limit", strconv.Itoa(limit))
	}
	if since != "" {
		v.Set("since", since)
	}
	if before != "" {
		v.Set("before", before)
	}
	if (filters != nil) && (len(filters) > 0) {
		buf, err := json.Marshal(filters)
		if err == nil {
			v.Set("filters", string(buf))
		}
	}

	uri := fmt.Sprintf("/v%s/containers/json?%s", API_VERSION, v.Encode())
	data, err := client.doRequest("GET", uri, nil, nil)
	if err != nil {
		return nil, err
	}

	containers := []Container{}
	if err := json.Unmarshal(data, &containers); err != nil {
		return nil, err
	}
	return containers, nil
}

func (client *DockerClient) RemoveContainer(name string, force bool) error {
	v := url.Values{}
	if force {
		v.Set("force", "1")
	}

	uri := fmt.Sprintf("/v%s/containers/%s?%s", API_VERSION, name, v.Encode())
	_, err := client.doRequest("DELETE", uri, nil, nil)
	if err != nil {
		return err
	}

	return nil
}
