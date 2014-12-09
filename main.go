package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/yungsang/talk2docker/commands"
	"github.com/yungsang/talk2docker/version"
)

func main() {
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} {{if .Flags}}[global options] {{end}}command{{if .Flags}} [options]{{end}} [arguments...]

VERSION:
   {{.Version}}{{if or .Author .Email}}

AUTHOR:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{if .Description}}{{.Description}}{{else}}{{.Usage}}{{end}}
   {{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`

	cli.CommandHelpTemplate = `NAME:
   {{.Name}} - {{.Description}}

USAGE:
   {{.Usage}}{{if .Flags}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{ end }}
`

	app := cli.NewApp()

	app.Name = "talk2docker"
	app.Usage = "A simple Docker client to talk to a Docker daemon"
	app.Version = version.Version
	app.Author = "YungSang"
	app.Email = "yungsang@gmail.com"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host, H",
			Value:  "unix:///var/run/docker.sock",
			Usage:  "Location of the Docker socket",
			EnvVar: "DOCKER_HOST",
		},
		cli.StringFlag{
			Name:   "tls",
			Usage:  "Path to the certificate files for TLS",
			EnvVar: "DOCKER_CERT_PATH",
		},
		cli.BoolFlag{
			Name:  "insecure-tls",
			Usage: "Skip verification of the certificates for TLS. Must verify by default.",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "ps",
			Usage:       app.Name + " [global options] ps [options]",
			Description: "List containers",
			Action:      commands.CommandPs,
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
			Name:        "images",
			Usage:       app.Name + " [global options] images [options] [NAME[:TAG]]",
			Description: "List images",
			Action:      commands.CommandImages,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "Show all images. Only named/taged and leaf images are shown by default.",
				},
				cli.BoolFlag{
					Name:  "quiet, q",
					Usage: "Only display numeric IDs",
				},
				cli.BoolFlag{
					Name:  "no-header, n",
					Usage: "Omit the header",
				},
			},
		},
		{
			Name:        "version",
			Usage:       app.Name + " [global options] version",
			Description: "Show the version information",
			Action:      commands.CommandVersion,
		},
	}

	app.Run(os.Args)
}
