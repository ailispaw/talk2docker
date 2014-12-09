package client

import (
	"github.com/codegangsta/cli"
	api "github.com/yungsang/dockerclient"
)

func GetDockerClient(ctx *cli.Context) (*api.DockerClient, error) {
	tlsConfig, err := getTLSConfig(ctx)
	if err != nil {
		return nil, err
	}
	return api.NewDockerClient(ctx.GlobalString("host"), tlsConfig)
}
