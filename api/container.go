package api

import (
	"encoding/binary"
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

func (client *DockerClient) CreateContainer(name string, config Config, hostConfig HostConfig) (string, error) {
	v := url.Values{}
	if name != "" {
		v.Set("name", name)
	}

	buf, err := json.Marshal(ConfigAndHostConfig{
		config,
		hostConfig,
	})
	if err != nil {
		return "", err
	}

	uri := fmt.Sprintf("/v%s/containers/create?%s", API_VERSION, v.Encode())
	data, err := client.doRequest("POST", uri, buf, nil)
	if err != nil {
		return "", err
	}

	var result struct {
		Id       string
		Warnings []string
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}

	for _, warning := range result.Warnings {
		fmt.Fprintf(client.out, "WARNING: %s\n", warning)
	}

	return result.Id, nil
}

func (client *DockerClient) InspectContainer(name string) (*ContainerInfo, error) {
	uri := fmt.Sprintf("/v%s/containers/%s/json", API_VERSION, name)
	data, err := client.doRequest("GET", uri, nil, nil)
	if err != nil {
		return nil, err
	}

	containerInfo := &ContainerInfo{}
	if err := json.Unmarshal(data, containerInfo); err != nil {
		return nil, err
	}
	return containerInfo, nil
}

func (client *DockerClient) StartContainer(name string) error {
	uri := fmt.Sprintf("/v%s/containers/%s/start", API_VERSION, name)
	if _, err := client.doRequest("POST", uri, nil, nil); err != nil {
		return err
	}
	return nil
}

func (client *DockerClient) StopContainer(name string, timeToWait int) error {
	v := url.Values{}
	if timeToWait > 0 {
		v.Set("t", strconv.Itoa(timeToWait))
	}

	uri := fmt.Sprintf("/v%s/containers/%s/stop?%s", API_VERSION, name, v.Encode())
	if _, err := client.doRequest("POST", uri, nil, nil); err != nil {
		return err
	}

	return nil
}

func (client *DockerClient) RestartContainer(name string, timeToWait int) error {
	v := url.Values{}
	if timeToWait > 0 {
		v.Set("t", strconv.Itoa(timeToWait))
	}

	uri := fmt.Sprintf("/v%s/containers/%s/restart?%s", API_VERSION, name, v.Encode())
	if _, err := client.doRequest("POST", uri, nil, nil); err != nil {
		return err
	}

	return nil
}

func (client *DockerClient) KillContainer(name, signal string) error {
	v := url.Values{}
	if signal != "" {
		v.Set("signal", signal)
	}

	uri := fmt.Sprintf("/v%s/containers/%s/kill?%s", API_VERSION, name, v.Encode())
	if _, err := client.doRequest("POST", uri, nil, nil); err != nil {
		return err
	}

	return nil
}

func (client *DockerClient) PauseContainer(name string) error {
	uri := fmt.Sprintf("/v%s/containers/%s/pause", API_VERSION, name)
	if _, err := client.doRequest("POST", uri, nil, nil); err != nil {
		return err
	}
	return nil
}

func (client *DockerClient) UnpauseContainer(name string) error {
	uri := fmt.Sprintf("/v%s/containers/%s/unpause", API_VERSION, name)
	if _, err := client.doRequest("POST", uri, nil, nil); err != nil {
		return err
	}
	return nil
}

func (client *DockerClient) WaitContainer(name string) (int, error) {
	uri := fmt.Sprintf("/v%s/containers/%s/wait", API_VERSION, name)
	data, err := client.doRequest("POST", uri, nil, nil)
	if err != nil {
		return 0, err
	}

	var result struct {
		StatusCode int
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return 0, err
	}

	return result.StatusCode, nil
}

func (client *DockerClient) RemoveContainer(name string, force bool) error {
	v := url.Values{}
	if force {
		v.Set("force", "1")
	}

	uri := fmt.Sprintf("/v%s/containers/%s?%s", API_VERSION, name, v.Encode())
	if _, err := client.doRequest("DELETE", uri, nil, nil); err != nil {
		return err
	}

	return nil
}

func (client *DockerClient) GetContainerLogs(name string, follow, stdout, stderr, timestamps bool, tail int) ([]string, error) {
	containerInfo, err := client.InspectContainer(name)
	if err != nil {
		return nil, err
	}

	v := url.Values{}
	if stdout {
		v.Set("stdout", "1")
	}
	if stderr {
		v.Set("stderr", "1")
	}
	if timestamps {
		v.Set("timestamps", "1")
	}
	if tail > 0 {
		v.Set("tail", strconv.Itoa(tail))
	}

	uri := fmt.Sprintf("/v%s/containers/%s/logs?%s", API_VERSION, name, v.Encode())
	data, err := client.doRequest("GET", uri, nil, nil)
	if err != nil {
		return nil, err
	}

	var logs []string
	if containerInfo.Config.Tty {
		logs = append(logs, string(data))
		logs = append(logs, "")
	} else {
		logs = getStreams(data)
	}
	return logs, nil
}

const (
	STREAM_HEADER_LENGTH = 8
	STREAM_TYPE_INDEX    = 0
	STREAM_SIZE_INDEX    = 4

	STREAM_TYPE_STDIN  = 0
	STREAM_TYPE_STDOUT = 1
	STREAM_TYPE_STDERR = 2
)

func getStreams(src []byte) []string {
	var (
		streams = make([]string, 2)
		size    = 0
	)

	for i := 0; len(src[i:]) > STREAM_HEADER_LENGTH; i += (STREAM_HEADER_LENGTH + size) {
		size = int(binary.BigEndian.Uint32(src[(i + STREAM_SIZE_INDEX):(i + STREAM_SIZE_INDEX + 4)]))

		buf := src[(i + STREAM_HEADER_LENGTH):(i + STREAM_HEADER_LENGTH + size)]

		switch src[i+STREAM_TYPE_INDEX] {
		case STREAM_TYPE_STDIN:
			fallthrough
		case STREAM_TYPE_STDOUT:
			streams[0] += string(buf)
		case STREAM_TYPE_STDERR:
			streams[1] += string(buf)
		default:
			continue
		}
	}

	return streams
}

func (client *DockerClient) ExportContainer(name string) error {
	uri := fmt.Sprintf("/v%s/containers/%s/export", API_VERSION, name)
	if _, err := client.doStreamRequest("GET", uri, nil, nil, true); err != nil {
		return err
	}
	return nil
}
