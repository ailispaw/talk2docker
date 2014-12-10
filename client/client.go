package client

import (
	"os"

	"github.com/codegangsta/cli"
	api "github.com/yungsang/dockerclient"
)

func GetDockerClient(ctx *cli.Context) (*api.DockerClient, error) {
	path := os.ExpandEnv(ctx.GlobalString("config"))

	config, err := LoadConfig(path)
	if err != nil {
		return nil, err
	}

	hostConfig, err := config.GetHostConfig(ctx.GlobalString("host"))
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
