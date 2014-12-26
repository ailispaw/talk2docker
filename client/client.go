package client

import (
	"io"
	"time"

	"github.com/yungsang/talk2docker/api"
)

func NewDockerClient(configPath, hostName string, out io.Writer) (*api.DockerClient, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	host, err := config.GetHost(hostName)
	if err != nil {
		return nil, err
	}

	tlsConfig, err := host.getTLSConfig()
	if err != nil {
		return nil, err
	}

	docker, err := api.NewDockerClient(host.URL, tlsConfig, 30*time.Second, out)
	if err != nil {
		return nil, err
	}

	return docker, nil
}
