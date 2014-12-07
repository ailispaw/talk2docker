package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/yungsang/talk2docker/commands"
)

func main() {
	app := cli.NewApp()

	app.Name = "talk2docker"
	app.Usage = "A simple Docker client to talk to a Docker daemon"
	app.Version = Version
	app.Author = "YungSang"
	app.Email = "yungsang@gmail.com"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host, H",
			Value:  "unix:///var/run/docker.sock",
			Usage:  "Location of the Docker socket",
			EnvVar: "DOCKER_HOST",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "ps",
			Usage:  "List containers",
			Action: commands.CommandPs,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "Show all containers. Only running containers are shown by default.",
				},
				cli.BoolFlag{
					Name:  "latest, l",
					Usage: "Show only the latest created container, include non-running ones.",
				},
				cli.BoolFlag{
					Name:  "quiet, q",
					Usage: "Only display numeric IDs",
				},
				cli.BoolFlag{
					Name:  "size, s",
					Usage: "Display sizes",
				},
				cli.BoolFlag{
					Name:  "no-header, n",
					Usage: "Omit the header",
				},
			},
		},
		{
			Name:   "images",
			Usage:  "List images",
			Action: commands.CommandImages,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "Show all images. The intermediate image layers are filtered out by default.",
				},
			},
		},
	}

	app.Run(os.Args)
}
