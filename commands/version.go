package commands

import (
	"log"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	api "github.com/yungsang/dockerclient"
	"github.com/yungsang/tablewriter"
	"github.com/yungsang/talk2docker/client"
	"github.com/yungsang/talk2docker/version"
)

func Version(ctx *cobra.Command, args []string) {
	var items [][]string

	out := []string{
		"Talk2Docker",
		"v" + version.Version,
		api.APIVersion,
		runtime.Version(),
		version.GITCOMMIT,
	}
	items = append(items, out)

	var e error

	docker, err := client.GetDockerClient(ctx)
	if err != nil {
		e = err
		goto Display
	}

	{
		version, err := docker.Version()
		if err != nil {
			e = err
			goto Display
		}

		out = []string{
			"Docker Server",
			"v" + version.Version,
			"v" + version.ApiVersion,
			version.GoVersion,
			version.GitCommit,
		}
		items = append(items, out)
	}

Display:
	var header = []string{
		"",
		"Version",
		"API Version",
		"Go Version",
		"Git commit",
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(items)
	table.Render()

	if e != nil {
		log.Fatal(e)
	}
}
