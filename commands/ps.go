package commands

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	api "github.com/yungsang/dockerclient"
	"github.com/yungsang/tablewriter"
	"github.com/yungsang/talk2docker/client"
)

func Ps(ctx *cobra.Command, args []string) {
	docker, err := client.GetDockerClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var filters = ""
	if GetBoolFlag(ctx, "latest") {
		filters += "&limit=1"
	}
	if GetBoolFlag(ctx, "size") {
		filters += "&size=1"
	}

	containers, err := docker.ListContainers(GetBoolFlag(ctx, "all"), GetBoolFlag(ctx, "size"), filters)
	if err != nil {
		log.Fatal(err)
	}

	if GetBoolFlag(ctx, "quiet") {
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

	formatPorts := func(ports []api.Port) string {
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
		if GetBoolFlag(ctx, "size") {
			out = append(out, fmt.Sprintf("%.3f", float64(container.SizeRw)/1000000.0))
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
	if GetBoolFlag(ctx, "size") {
		header = append(header, "Size(MB)")
	}

	table := tablewriter.NewWriter(os.Stdout)
	if !GetBoolFlag(ctx, "no-header") {
		table.SetHeader(header)
	}
	table.SetBorder(false)
	table.AppendBulk(items)
	table.Render()
}
