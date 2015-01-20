package commands

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yungsang/tablewriter"

	"github.com/ailispaw/talk2docker/api"
	"github.com/ailispaw/talk2docker/client"
)

var (
	boolLatest, boolSize bool
)

var cmdPs = &cobra.Command{
	Use:     "ps",
	Aliases: []string{"containers"},
	Short:   "List containers",
	Long:    APP_NAME + " ps - List containers",
	Run:     listContainers,
}

var cmdContainer = &cobra.Command{
	Use:     "container [command]",
	Aliases: []string{"ctn"},
	Short:   "Manage containers",
	Long:    APP_NAME + " container - Manage containers",
	Run: func(ctx *cobra.Command, args []string) {
		ctx.Help()
	},
}

var cmdListContainers = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List containers",
	Long:    APP_NAME + " container list - List containers",
	Run:     listContainers,
}

var cmdRemoveContainer = &cobra.Command{
	Use:     "remove <NAME|ID>...",
	Aliases: []string{"rm"},
	Short:   "Remove containers",
	Long:    APP_NAME + " remove - Remove containers",
	Run:     removeContainers,
}

func init() {
	flags := cmdPs.Flags()
	flags.BoolVarP(&boolAll, "all", "a", false, "Show all containers. Only running containers are shown by default.")
	flags.BoolVarP(&boolLatest, "latest", "l", false, "Show only the latest created container, include non-running ones.")
	flags.BoolVarP(&boolQuiet, "quiet", "q", false, "Only display numeric IDs")
	flags.BoolVarP(&boolSize, "size", "s", false, "Display sizes")
	flags.BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")

	flags = cmdListContainers.Flags()
	flags.BoolVarP(&boolAll, "all", "a", false, "Show all containers. Only running containers are shown by default.")
	flags.BoolVarP(&boolLatest, "latest", "l", false, "Show only the latest created container, include non-running ones.")
	flags.BoolVarP(&boolQuiet, "quiet", "q", false, "Only display numeric IDs")
	flags.BoolVarP(&boolSize, "size", "s", false, "Display sizes")
	flags.BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")
	cmdContainer.AddCommand(cmdListContainers)

	flags = cmdRemoveContainer.Flags()
	flags.BoolVarP(&boolForce, "force", "f", false, "Force the removal of a running container")
	cmdContainer.AddCommand(cmdRemoveContainer)
}

func listContainers(ctx *cobra.Command, args []string) {
	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	limit := 0
	if boolLatest {
		limit = 1
	}

	containers, err := docker.ListContainers(boolAll, boolSize, limit, "", "", nil)
	if err != nil {
		log.Fatal(err)
	}

	if boolQuiet {
		for _, container := range containers {
			ctx.Println(Truncate(container.Id, 12))
		}
		return
	}

	if boolYAML || boolJSON {
		if err := FormatPrint(ctx.Out(), containers); err != nil {
			log.Fatal(err)
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
			Truncate(container.Command, 30),
			FormatDateTime(time.Unix(container.Created, 0)),
			container.Status,
			formatPorts(container.Ports),
		}
		if boolSize {
			out = append(out, FormatFloat(float64(container.SizeRw)/1000000))
		}
		items = append(items, out)
	}

	header := []string{
		"ID",
		"Names",
		"Image",
		"Command",
		"Created at",
		"Status",
		"Ports",
	}
	if boolSize {
		header = append(header, "Size(MB)")
	}

	PrintInTable(ctx.Out(), header, items, 0, tablewriter.ALIGN_DEFAULT)
}

func removeContainers(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ctx.Println("Needs an argument <NAME|ID> at least to remove")
		ctx.Usage()
		return
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	var gotError = false
	for _, name := range args {
		if err := docker.RemoveContainer(name, boolForce); err != nil {
			log.Println(err)
			gotError = true
		} else {
			ctx.Println(name)
		}
	}
	if gotError {
		log.Fatal("Error: failed to remove one or more containers")
	}
}
