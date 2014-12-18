package commands

import (
	"log"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	api "github.com/yungsang/dockerclient"
	"github.com/yungsang/tablewriter"
	"github.com/yungsang/talk2docker/client"
	v "github.com/yungsang/talk2docker/version"
)

var cmdVersion = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Show the version information",
	Long:    appName + " version - Show the version information",
	Run:     version,
}

func init() {
	cmdVersion.Flags().BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")
}

func version(ctx *cobra.Command, args []string) {
	var items [][]string

	out := []string{
		"Talk2Docker",
		"v" + v.Version,
		api.APIVersion,
		runtime.Version(),
		v.GITCOMMIT,
	}
	items = append(items, out)

	var e error

	docker, err := client.GetDockerClient(configPath, hostName)
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
