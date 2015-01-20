package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func (client *DockerClient) ListImages(all bool, filters map[string][]string) ([]Image, error) {
	v := url.Values{}
	if all {
		v.Set("all", "1")
	}
	if (filters != nil) && (len(filters) > 0) {
		buf, err := json.Marshal(filters)
		if err == nil {
			v.Set("filters", string(buf))
		}
	}

	uri := fmt.Sprintf("/v%s/images/json?%s", API_VERSION, v.Encode())
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

func (client *DockerClient) PullImage(name string) error {
	v := url.Values{}
	v.Set("fromImage", name)

	uri := fmt.Sprintf("/v%s/images/create?%s", API_VERSION, v.Encode())

	return client.doStreamRequest("POST", uri, nil, nil)
}

func (client *DockerClient) GetImageHistory(name string) (ImageHistories, error) {
	uri := fmt.Sprintf("/v%s/images/%s/history", API_VERSION, name)
	data, err := client.doRequest("GET", uri, nil, nil)
	if err != nil {
		return nil, err
	}

	images := ImageHistories{}
	if err := json.Unmarshal(data, &images); err != nil {
		return nil, err
	}
	return images, nil
}

func (client *DockerClient) TagImage(name, repo, tag string, force bool) error {
	v := url.Values{}
	v.Set("repo", repo)
	v.Set("tag", tag)
	if force {
		v.Set("force", "1")
	}

	uri := fmt.Sprintf("/v%s/images/%s/tag?%s", API_VERSION, name, v.Encode())
	_, err := client.doRequest("POST", uri, nil, nil)
	return err
}

func (client *DockerClient) InspectImage(name string) (*ImageInfo, error) {
	uri := fmt.Sprintf("/v%s/images/%s/json", API_VERSION, name)
	data, err := client.doRequest("GET", uri, nil, nil)
	if err != nil {
		return nil, err
	}

	imageInfo := &ImageInfo{}
	if err := json.Unmarshal(data, imageInfo); err != nil {
		return nil, err
	}
	return imageInfo, nil
}

func (client *DockerClient) PushImage(name, tag, credentials string) error {
	v := url.Values{}
	v.Set("tag", tag)

	uri := fmt.Sprintf("/v%s/images/%s/push?%s", API_VERSION, name, v.Encode())

	headers := map[string]string{}
	headers["X-Registry-Auth"] = credentials

	return client.doStreamRequest("POST", uri, nil, headers)
}

func (client *DockerClient) RemoveImage(name string, force, noprune bool) error {
	v := url.Values{}
	if force {
		v.Set("force", "1")
	}
	if noprune {
		v.Set("noprune", "1")
	}

	uri := fmt.Sprintf("/v%s/images/%s?%s", API_VERSION, name, v.Encode())
	data, err := client.doRequest("DELETE", uri, nil, nil)
	if err != nil {
		return err
	}

	messages := []map[string]string{}
	if err := json.Unmarshal(data, &messages); err != nil {
		return err
	}

	for _, message := range messages {
		_, isDeleted := message["Deleted"]
		if isDeleted {
			fmt.Fprintf(client.out, "Deleted: %s\n", message["Deleted"])
		} else {
			fmt.Fprintf(client.out, "Untagged: %s\n", message["Untagged"])
		}
	}

	return nil
}

func (client *DockerClient) SearchImages(term string) (ImageSearchResults, error) {
	v := url.Values{}
	v.Set("term", term)

	uri := fmt.Sprintf("/v%s/images/search?%s", API_VERSION, v.Encode())
	data, err := client.doRequest("GET", uri, nil, nil)
	if err != nil {
		return nil, err
	}

	images := ImageSearchResults{}
	if err := json.Unmarshal(data, &images); err != nil {
		return nil, err
	}
	return images, nil
}
