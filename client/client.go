package client

import (
	"os"

	"github.com/spf13/cobra"
	api "github.com/yungsang/dockerclient"
)

func GetDockerClient(ctx *cobra.Command) (*api.DockerClient, error) {
	path := os.ExpandEnv(ctx.Flag("config").Value.String())

	config, err := LoadConfig(path)
	if err != nil {
		return nil, err
	}

	hostConfig, err := config.GetHostConfig(ctx.Flag("host").Value.String())
	if err != nil {
		return nil, err
	}

	tlsConfig, err := getTLSConfig(hostConfig)
	if err != nil {
		return nil, err
	}

	docker, err := api.NewDockerClient(hostConfig.Host, tlsConfig)
	if err != nil {
		return nil, err
	}

	err = config.SaveConfig(path)
	if err != nil {
		return nil, err
	}

	return docker, nil
}
