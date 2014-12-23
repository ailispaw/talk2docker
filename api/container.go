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
		v.Add("all", "1")
	}
	if size {
		v.Add("size", "1")
	}
	if limit > 0 {
		v.Add("limit", strconv.Itoa(limit))
	}
	if since != "" {
		v.Add("since", since)
	}
	if before != "" {
		v.Add("before", before)
	}
	if (filters != nil) && (len(filters) > 0) {
		buf, err := json.Marshal(filters)
		if err == nil {
			v.Add("filters", string(buf))
		}
	}

	uri := fmt.Sprintf("/v%s/containers/json?%s", API_VERSION, v.Encode())
	data, err := client.doRequest("GET", uri, nil, nil)
	if err != nil {
		return nil, err
	}

	containers := []Container{}
	err = json.Unmarshal(data, &containers)
	if err != nil {
		return nil, err
	}
	return containers, nil
}
