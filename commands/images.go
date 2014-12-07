package commands

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/olekukonko/tablewriter"
	docker "github.com/yungsang/dockerclient"
)

func CommandImages(ctx *cli.Context) {
	client, err := docker.NewDockerClient(ctx.GlobalString("host"), nil)
	if err != nil {
		log.Fatal(err)
	}

	images, err := client.ListImages(ctx.Bool("all"))
	if err != nil {
		log.Fatal(err)
	}

	var items [][]string
	for _, image := range images {
		out := []string{
			Truncate(image.Id, 12),
			strings.Join(image.RepoTags, ", "),
			FormatDateTime(time.Unix(image.Created, 0)),
			fmt.Sprintf("%.4g MB", float64(image.VirtualSize)/1000000.0),
		}
		items = append(items, out)
	}

	var header = []string{
		"ID",
		"Name:Tags",
		"Created at",
		"Virtual Size",
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetBorder(false)
	table.AppendBulk(items)
	table.Render()
}
