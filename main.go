package main

import (
	"github.com/spf13/cobra"
	"github.com/yungsang/talk2docker/commands"
)

func main() {
	var app = &cobra.Command{
		Use:  "talk2docker",
		Long: "Talk2Docker - A simple Docker client to talk to Docker daemon",
	}
	app.PersistentFlags().String(
		"config", "$HOME/.talk2docker/config", "Path to the configuration file")
	app.PersistentFlags().String(
		"host", "", "Hostname to use its config (runtime only)")

	// ps command
	var cmdPs = &cobra.Command{
		Use:   "ps",
		Short: "List containers",
		Run:   commands.CommandPs,
	}
	cmdPs.Flags().BoolP(
		"all", "a", false, "Show all containers. Only running containers are shown by default.")
	cmdPs.Flags().BoolP(
		"latest", "l", false, "Show only the latest created container, include non-running ones.")
	cmdPs.Flags().BoolP(
		"quiet", "q", false, "Only display numeric IDs")
	cmdPs.Flags().BoolP(
		"size", "s", false, "Display sizes")
	cmdPs.Flags().BoolP(
		"no-header", "n", false, "Omit the header")
	app.AddCommand(cmdPs)

	// images command
	var cmdImages = &cobra.Command{
		Use:   "images",
		Short: "List images",
		Run:   commands.CommandImages,
	}
	cmdImages.Flags().BoolP(
		"all", "a", false, "Show all images. Only named/taged and leaf images are shown by default.")
	cmdImages.Flags().BoolP(
		"quiet", "q", false, "Only display numeric IDs")
	cmdImages.Flags().BoolP(
		"no-header", "n", false, "Omit the header")
	app.AddCommand(cmdImages)

	// version command
	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Show the version information",
		Run:   commands.CommandVersion,
	}
	app.AddCommand(cmdVersion)

	app.Execute()
}
