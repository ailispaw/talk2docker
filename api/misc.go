package api

import (
	"encoding/json"
	"fmt"
)

func (client *DockerClient) Auth(auth *AuthConfig) error {
	data, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("/v%s/auth", API_VERSION)
	_, err = client.doRequest("POST", uri, data, nil)
	return err
}

func (client *DockerClient) Info() (*Info, error) {
	uri := fmt.Sprintf("/v%s/info", API_VERSION)
	data, err := client.doRequest("GET", uri, nil, nil)
	if err != nil {
		return nil, err
	}

	info := &Info{}
	if err := json.Unmarshal(data, info); err != nil {
		return nil, err
	}
	return info, nil
}

func (client *DockerClient) Version() (*Version, error) {
	uri := fmt.Sprintf("/v%s/version", API_VERSION)
	data, err := client.doRequest("GET", uri, nil, nil)
	if err != nil {
		return nil, err
	}

	version := &Version{}
	if err := json.Unmarshal(data, version); err != nil {
		return nil, err
	}
	return version, nil
}
