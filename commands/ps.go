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

func CommandPs(ctx *cli.Context) {
	client, err := docker.NewDockerClient(ctx.GlobalString("host"), nil)
	if err != nil {
		log.Fatal(err)
	}

	var filters = ""
	if ctx.Bool("latest") {
		filters += "&limit=1"
	}
	if ctx.Bool("size") {
		filters += "&size=1"
	}

	containers, err := client.ListContainers(ctx.Bool("all"), ctx.Bool("size"), filters)
	if err != nil {
		log.Fatal(err)
	}

	if ctx.Bool("quiet") {
		for _, container := range containers {
			fmt.Println(Truncate(container.Id, 12))
		}
		return
	}

	trimNamePrefix := func(ss []string) []string {
		for i, s := range ss {
			ss[i] = strings.TrimPrefix(s, "/")
		}
		return ss
	}

	formatPorts := func(ports []docker.Port) string {
		result := []string{}
		for _, p := range ports {
			result = append(result, fmt.Sprintf("%s:%d->%d/%s",
				p.IP, p.PublicPort, p.PrivatePort, p.Type))
		}
		return strings.Join(result, ", ")
	}

	var items [][]string
	for _, container := range containers {
		out := []string{
			Truncate(container.Id, 12),
			strings.Join(trimNamePrefix(container.Names), ", "),
			container.Image,
			Truncate(container.Command, 20),
			FormatDateTime(time.Unix(container.Created, 0)),
			container.Status,
			formatPorts(container.Ports),
		}
		if ctx.Bool("size") {
			out = append(out, fmt.Sprintf("%.4g MB", float64(container.SizeRw)/1000000.0))
		}
		items = append(items, out)
	}

	var header = []string{
		"ID",
		"Names",
		"Image",
		"Command",
		"Created at",
		"Status",
		"Ports",
	}
	if ctx.Bool("size") {
		header = append(header, "Size")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetBorder(false)
	table.AppendBulk(items)
	table.Render()
}
