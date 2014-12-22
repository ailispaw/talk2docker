package commands

import (
	"log"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/yungsang/tablewriter"
	"github.com/yungsang/talk2docker/api"
	"github.com/yungsang/talk2docker/client"
	"github.com/yungsang/talk2docker/version"
)

var cmdVersion = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Show the version information",
	Long:    appName + " version - Show the version information",
	Run:     showVersion,
}

func init() {
	cmdVersion.Flags().BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")
}

func showVersion(ctx *cobra.Command, args []string) {
	var items [][]string

	out := []string{
		"Talk2Docker",
		version.Version,
		api.APIVersion,
		runtime.Version(),
		version.GITCOMMIT,
	}
	items = append(items, out)

	var e error

	docker, err := client.NewDockerClient(configPath, hostName)
	if err != nil {
		e = err
		goto Display
	}

	{
		dockerVersion, err := docker.Version()
		if err != nil {
			e = err
			goto Display
		}

		out = []string{
			"Docker Server",
			dockerVersion.Version,
			dockerVersion.ApiVersion,
			dockerVersion.GoVersion,
			dockerVersion.GitCommit,
		}
		items = append(items, out)
	}

Display:
	header := []string{
		"",
		"Version",
		"API Version",
		"Go Version",
		"Git commit",
	}

	table := tablewriter.NewWriter(os.Stdout)
	if !boolNoHeader {
		table.SetHeader(header)
	} else {
		table.SetBorder(false)
	}
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(items)
	table.Render()

	if e != nil {
		log.Fatal(e)
	}
}
