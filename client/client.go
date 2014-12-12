package client

import (
	"os"

	api "github.com/yungsang/dockerclient"
)

func GetDockerClient(configPath, hostName string) (*api.DockerClient, error) {
	path := os.ExpandEnv(configPath)

	config, err := LoadConfig(path)
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

	err = config.SaveConfig(path)
	if err != nil {
		return nil, err
	}

	return docker, nil
}
