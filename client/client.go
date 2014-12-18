package client

import (
	api "github.com/yungsang/dockerclient"
)

func GetDockerClient(configPath, hostName string) (*api.DockerClient, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	host, err := config.GetHost(hostName)
	if err != nil {
		return nil, err
	}

	tlsConfig, err := getTLSConfig(host)
	if err != nil {
		return nil, err
	}

	docker, err := api.NewDockerClient(host.URL, tlsConfig)
	if err != nil {
		return nil, err
	}

	return docker, nil
}
