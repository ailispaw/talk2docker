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
	Long:    APP_NAME + " version - Show the version information",
	Run:     showVersion,
}

func init() {
	cmdVersion.Flags().BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")
	cmdVersion.Flags().BoolVarP(&boolJSON, "json", "j", false, "Output in JSON format")
}

func showVersion(ctx *cobra.Command, args []string) {
	data := map[string]api.Version{}

	data[APP_NAME] = api.Version{
		Version:    version.APP_VERSION,
		ApiVersion: api.API_VERSION,
		GoVersion:  runtime.Version(),
		GitCommit:  version.GIT_COMMIT,
	}

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

		data["Docker Server"] = *dockerVersion
	}

Display:
	if boolJSON {
		err = PrintInJSON(data)
		if err != nil {
			log.Fatal(err)
		}
		if e != nil {
			log.Fatal(e)
		}
		return
	}

	var items [][]string

	for key, value := range data {
		out := []string{
			key,
			value.Version,
			value.ApiVersion,
			value.GoVersion,
			value.GitCommit,
		}
		items = append(items, out)
	}

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
