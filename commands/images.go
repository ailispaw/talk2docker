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

	var matchName = ""
	if len(ctx.Args()) > 0 {
		matchName = ctx.Args()[0]
	}

	matchImageByName := func(tags []string, name string) bool {
		_name := strings.Split(name, ":")

		for _, tag := range tags {
			_tag := strings.Split(tag, ":")
			if _tag[0] == _name[0] {
				if (len(_name) < 2) || (_tag[1] == _name[1]) {
					return true
				}
			}
		}

		return false
	}

	if ctx.Bool("quiet") {
		for _, image := range images {
			if (matchName == "") || matchImageByName(image.RepoTags, matchName) {
				fmt.Println(Truncate(image.Id, 12))
			}
		}
		return
	}

	var items [][]string
	for _, image := range images {
		if (matchName == "") || matchImageByName(image.RepoTags, matchName) {
			out := []string{
				Truncate(image.Id, 12),
				strings.Join(image.RepoTags, ", "),
				FormatDateTime(time.Unix(image.Created, 0)),
				fmt.Sprintf("%.4g MB", float64(image.VirtualSize)/1000000.0),
			}
			items = append(items, out)
		}
	}

	var header = []string{
		"ID",
		"Name:Tags",
		"Created at",
		"Virtual Size",
	}

	table := tablewriter.NewWriter(os.Stdout)
	if !ctx.Bool("no-header") {
		table.SetHeader(header)
	}
	table.SetBorder(false)
	table.AppendBulk(items)
	table.Render()
}
