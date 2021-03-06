package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/yungsang/tablewriter"

	"github.com/ailispaw/talk2docker/api"
	"github.com/ailispaw/talk2docker/client"
)

var (
	boolLatest, boolSize, boolTimestamps, boolPause bool

	timeToWait, tail        int
	signal, message, author string
)

var cmdPs = &cobra.Command{
	Use:     "ps",
	Aliases: []string{"containers"},
	Short:   "List containers",
	Long:    APP_NAME + " ps - List containers",
	Run:     listContainers,
}

var cmdCommit = &cobra.Command{
	Use:   "commit <NAME|ID> <NAME[:TAG]>",
	Short: "Create a new image from a container",
	Long:  APP_NAME + " commit - Create a new image from a container",
	Run:   commitContainer,
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

var cmdInspectContainers = &cobra.Command{
	Use:     "inspect <NAME|ID>...",
	Aliases: []string{"ins", "info"},
	Short:   "Show containers' information",
	Long:    APP_NAME + " container inspect - Show containers' information",
	Run:     inspectContainers,
}

var cmdStartContainers = &cobra.Command{
	Use:     "start <NAME|ID>...",
	Aliases: []string{"up"},
	Short:   "Start containers",
	Long:    APP_NAME + " container start - Start containers",
	Run:     startContainers,
}

var cmdStopContainers = &cobra.Command{
	Use:     "stop <NAME|ID>...",
	Aliases: []string{"down"},
	Short:   "Stop containers",
	Long:    APP_NAME + " container stop - Stop containers",
	Run:     stopContainers,
}

var cmdRestartContainers = &cobra.Command{
	Use:   "restart <NAME|ID>...",
	Short: "Restart containers",
	Long:  APP_NAME + " container restart - Restart containers",
	Run:   restartContainers,
}

var cmdKillContainers = &cobra.Command{
	Use:   "kill <NAME|ID>...",
	Short: "Kill containers",
	Long:  APP_NAME + " container kill - Kill containers",
	Run:   killContainers,
}

var cmdPauseContainers = &cobra.Command{
	Use:     "pause <NAME|ID>...",
	Aliases: []string{"suspend"},
	Short:   "Pause containers",
	Long:    APP_NAME + " container pause - Pause containers",
	Run:     pauseContainers,
}

var cmdUnpauseContainers = &cobra.Command{
	Use:     "unpause <NAME|ID>...",
	Aliases: []string{"resume"},
	Short:   "Unpause containers",
	Long:    APP_NAME + " container unpause - Unpause containers",
	Run:     unpauseContainers,
}

var cmdWaitContainers = &cobra.Command{
	Use:   "wait <NAME|ID>...",
	Short: "Block until containers stop",
	Long:  APP_NAME + " container wait - Block until containers stop",
	Run:   waitContainers,
}

var cmdRemoveContainers = &cobra.Command{
	Use:     "remove <NAME|ID>...",
	Aliases: []string{"rm"},
	Short:   "Remove containers",
	Long:    APP_NAME + " container remove - Remove containers",
	Run:     removeContainers,
}

var cmdGetContainerLogs = &cobra.Command{
	Use:   "logs <NAME|ID>",
	Short: "Stream outputs(STDOUT/STDERR) from a container",
	Long:  APP_NAME + " container logs - Stream outputs(STDOUT/STDERR) from a container",
	Run:   getContainerLogs,
}

var cmdGetContainerChanges = &cobra.Command{
	Use:   "diff <NAME|ID>",
	Short: "Show changes on a container's filesystem",
	Long:  APP_NAME + " container diff - Show changes on a container's filesystem",
	Run:   getContainerChanges,
}

var cmdExportContainer = &cobra.Command{
	Use:   "export <NAME|ID> [PATH]",
	Short: "Stream the contents of a container as a tar archive",
	Long:  APP_NAME + " container export - Stream the contents of a container as a tar archive",
	Run:   exportContainer,
}

var cmdGetContainerProcesses = &cobra.Command{
	Use:     "top <NAME|ID> [PS-ARGS]",
	Aliases: []string{"ps"},
	Short:   "List the running processes of a container",
	Long:    APP_NAME + " container top - List the running processes of a container",
	Run:     getContainerProcesses,
}

var cmdCommitContainer = &cobra.Command{
	Use:   "commit <NAME|ID> <IMAGE-NAME[:TAG]>",
	Short: "Create a new image from a container",
	Long:  APP_NAME + " container commit - Create a new image from a container",
	Run:   commitContainer,
}

func init() {
	for _, flags := range []*pflag.FlagSet{cmdPs.Flags(), cmdListContainers.Flags()} {
		flags.BoolVarP(&boolAll, "all", "a", false, "Show all containers. Only running containers are shown by default.")
		flags.BoolVarP(&boolLatest, "latest", "l", false, "Show only the latest created container, include non-running ones.")
		flags.BoolVarP(&boolQuiet, "quiet", "q", false, "Only display numeric IDs")
		flags.BoolVarP(&boolSize, "size", "s", false, "Display sizes")
		flags.BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")
	}

	cmdContainer.AddCommand(cmdListContainers)

	for _, flags := range []*pflag.FlagSet{cmdCommit.Flags(), cmdCommitContainer.Flags()} {
		flags.StringVarP(&message, "message", "m", "", "Commit message")
		flags.StringVarP(&author, "author", "a", "", "Author (e.g., \"A.I. <ailis@paw.zone>\")")
		flags.BoolVarP(&boolPause, "pause", "p", true, "Pause container during commit")
	}

	cmdContainer.AddCommand(cmdInspectContainers)

	cmdContainer.AddCommand(cmdStartContainers)

	flags := cmdStopContainers.Flags()
	flags.IntVarP(&timeToWait, "time", "t", 10, "Number of seconds to wait for the container to stop before killing it. Default is 10 seconds.")
	cmdContainer.AddCommand(cmdStopContainers)

	flags = cmdRestartContainers.Flags()
	flags.IntVarP(&timeToWait, "time", "t", 10, "Number of seconds to wait for the container to stop before killing it. Default is 10 seconds.")
	cmdContainer.AddCommand(cmdRestartContainers)

	flags = cmdKillContainers.Flags()
	flags.StringVarP(&signal, "signal", "s", "SIGKILL", "Signal to send to the container")
	cmdContainer.AddCommand(cmdKillContainers)

	cmdContainer.AddCommand(cmdPauseContainers)

	cmdContainer.AddCommand(cmdUnpauseContainers)

	cmdContainer.AddCommand(cmdWaitContainers)

	flags = cmdRemoveContainers.Flags()
	flags.BoolVarP(&boolForce, "force", "f", false, "Force the removal of a running container")
	cmdContainer.AddCommand(cmdRemoveContainers)

	flags = cmdGetContainerLogs.Flags()
	flags.BoolVarP(&boolTimestamps, "timestamps", "t", false, "Show timestamps")
	flags.IntVar(&tail, "tail", 0, "Output the specified number of lines at the end of logs (0 for all)")
	cmdContainer.AddCommand(cmdGetContainerLogs)

	cmdContainer.AddCommand(cmdGetContainerChanges)

	cmdContainer.AddCommand(cmdExportContainer)

	cmdContainer.AddCommand(cmdGetContainerProcesses)

	cmdContainer.AddCommand(cmdCommitContainer)

	cmdContainer.AddCommand(cmdUploadToContainer)
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
			if p.IP == "" {
				result = append(result, fmt.Sprintf("%d/%s", p.PrivatePort, p.Type))
			} else {
				result = append(result, fmt.Sprintf("%s:%d->%d/%s",
					p.IP, p.PublicPort, p.PrivatePort, p.Type))
			}
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

func inspectContainers(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument <NAME|ID> at least to inspect")
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	var containers []api.ContainerInfo
	var gotError = false

	for _, name := range args {
		if containerInfo, err := docker.InspectContainer(name); err != nil {
			log.Error(err)
			gotError = true
		} else {
			containers = append(containers, *containerInfo)
		}
	}

	if len(containers) > 0 {
		if err := FormatPrint(ctx.Out(), containers); err != nil {
			log.Fatal(err)
		}
	}

	if gotError {
		log.Fatal("Error: failed to inspect one or more containers")
	}
}

func startContainers(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument <NAME|ID> at least to start")
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	var gotError = false
	for _, name := range args {
		if err := docker.StartContainer(name); err != nil {
			log.Error(err)
			gotError = true
		} else {
			ctx.Println(name)
		}
	}
	if gotError {
		log.Fatal("Error: failed to start one or more containers")
	}
}

func stopContainers(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument <NAME|ID> at least to stop")
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	var gotError = false
	for _, name := range args {
		if err := docker.StopContainer(name, timeToWait); err != nil {
			log.Error(err)
			gotError = true
		} else {
			ctx.Println(name)
		}
	}
	if gotError {
		log.Fatal("Error: failed to stop one or more containers")
	}
}

func restartContainers(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument <NAME|ID> at least to restart")
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	var gotError = false
	for _, name := range args {
		if err := docker.RestartContainer(name, timeToWait); err != nil {
			log.Error(err)
			gotError = true
		} else {
			ctx.Println(name)
		}
	}
	if gotError {
		log.Fatal("Error: failed to restart one or more containers")
	}
}

func killContainers(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument <NAME|ID> at least to kill")
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	var gotError = false
	for _, name := range args {
		if err := docker.KillContainer(name, signal); err != nil {
			log.Error(err)
			gotError = true
		} else {
			ctx.Println(name)
		}
	}
	if gotError {
		log.Fatal("Error: failed to kill one or more containers")
	}
}

func pauseContainers(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument <NAME|ID> at least to pause")
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	var gotError = false
	for _, name := range args {
		if err := docker.PauseContainer(name); err != nil {
			log.Error(err)
			gotError = true
		} else {
			ctx.Println(name)
		}
	}
	if gotError {
		log.Fatal("Error: failed to pause one or more containers")
	}
}

func unpauseContainers(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument <NAME|ID> at least to unpause")
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	var gotError = false
	for _, name := range args {
		if err := docker.UnpauseContainer(name); err != nil {
			log.Error(err)
			gotError = true
		} else {
			ctx.Println(name)
		}
	}
	if gotError {
		log.Fatal("Error: failed to unpause one or more containers")
	}
}

func waitContainers(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument <NAME|ID> at least to wait")
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	var gotError = false
	for _, name := range args {
		if status, err := docker.WaitContainer(name); err != nil {
			log.Error(err)
			gotError = true
		} else {
			fmt.Fprintf(ctx.Out(), "%s: %d\n", name, status)
		}
	}
	if gotError {
		log.Fatal("Error: failed to wait one or more containers")
	}
}

func removeContainers(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument <NAME|ID> at least to remove")
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	var gotError = false
	for _, name := range args {
		if err := docker.RemoveContainer(name, boolForce); err != nil {
			log.Error(err)
			gotError = true
		} else {
			ctx.Println(name)
		}
	}
	if gotError {
		log.Fatal("Error: failed to remove one or more containers")
	}
}

func getContainerLogs(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument <NAME|ID> to get logs")
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	logs, err := docker.GetContainerLogs(args[0], false, true, true, boolTimestamps, tail)
	if err != nil {
		log.Fatal(err)
	}

	if logs[0] != "" {
		fmt.Fprint(os.Stdout, logs[0])
	}
	if logs[1] != "" {
		fmt.Fprint(os.Stderr, logs[1])
	}
}

func getContainerChanges(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument <NAME|ID> to get changes")
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	changes, err := docker.GetContainerChanges(args[0])
	if err != nil {
		log.Fatal(err)
	}

	for _, change := range changes {
		var kind string
		switch change.Kind {
		case api.CHANGE_TYPE_MODIFY:
			kind = "C"
		case api.CHANGE_TYPE_ADD:
			kind = "A"
		case api.CHANGE_TYPE_DELETE:
			kind = "D"
		}
		ctx.Printf("%s %s\n", kind, change.Path)
	}
}

func exportContainer(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument <NAME|ID> to export")
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	if len(args) > 1 {
		if err := docker.CopyContainer(args[0], args[1]); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := docker.ExportContainer(args[0]); err != nil {
			log.Fatal(err)
		}
	}
}

func getContainerProcesses(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument <NAME|ID> to execute ps")
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	ps_args := ""
	if len(args) > 1 {
		ps_args = strings.Join(args[1:], " ")
	}

	ps, err := docker.GetContainerProcesses(args[0], ps_args)
	if err != nil {
		log.Fatal(err)
	}

	if boolYAML || boolJSON {
		if err := FormatPrint(ctx.Out(), ps); err != nil {
			log.Fatal(err)
		}
		return
	}

	PrintInTable(ctx.Out(), ps.Titles, ps.Processes, 100, tablewriter.ALIGN_DEFAULT)
}

func commitContainer(ctx *cobra.Command, args []string) {
	if len(args) < 2 {
		ErrorExit(ctx, "Needs two arguments to commit <CONTAINER-NAME|ID> to <IMAGE-NAME[:TAG]>")
	}

	reg, name, tag, err := client.ParseRepositoryName(args[1])
	if err != nil {
		log.Fatal(err)
	}

	if reg != "" {
		name = reg + "/" + name
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	if _, err := docker.CommitContainer(args[0], name, tag, message, author, boolPause); err != nil {
		log.Fatal(err)
	}

	ctx.Printf("Committed %s as %s:%s\n", args[0], name, tag)
}
