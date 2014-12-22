package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func (client *DockerClient) ListImages(all bool, filters map[string][]string) ([]Image, error) {
	v := url.Values{}
	if all {
		v.Add("all", "1")
	}
	if (filters != nil) && (len(filters) > 0) {
		buf, err := json.Marshal(filters)
		if err == nil {
			v.Add("filters", string(buf))
		}
	}

	uri := fmt.Sprintf("/v%s/images/json?%s", APIVersion, v.Encode())
	data, err := client.doRequest("GET", uri, nil, nil)
	if err != nil {
		return nil, err
	}

	images := []Image{}
	if err := json.Unmarshal(data, &images); err != nil {
		return nil, err
	}
	return images, nil
}

func (client *DockerClient) PullImage(name, auth string) error {
	v := url.Values{}
	v.Add("fromImage", name)

	uri := fmt.Sprintf("/v%s/images/create?%s", APIVersion, v.Encode())

	headers := map[string]string{}
	if auth != "" {
		headers["X-Registry-Auth"] = auth
	}

	return client.doStreamRequest("POST", uri, nil, headers)
}
