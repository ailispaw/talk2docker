package main

import (
	"github.com/spf13/cobra"
	"github.com/yungsang/talk2docker/commands"
)

func main() {
	var appName = "Talk2Docker"

	var app = &cobra.Command{
		Use:   "talk2docker",
		Short: "A simple Docker client to talk to Docker daemon",
		Long:  appName + " - A simple Docker client to talk to Docker daemon",
	}
	app.PersistentFlags().String(
		"config", "$HOME/.talk2docker/config", "Path to the configuration file")
	app.PersistentFlags().String(
		"host", "", "Hostname to use its config (runtime only)")

	// ps command
	var cmdPs = &cobra.Command{
		Use:   "ps",
		Short: "List containers",
		Long:  appName + " ps - List containers",
		Run:   commands.Ps,
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

	// image command
	var cmdImage = &cobra.Command{
		Use:   "image [command]",
		Short: "Manage images",
		Long:  appName + " image - Manage images",
		Run: func(ctx *cobra.Command, args []string) {
			ctx.Usage()
		},
	}

	var cmdListImages = &cobra.Command{
		Use:   "list [NAME[:TAG]]",
		Short: "List images",
		Long:  appName + " image list - List images",
		Run:   commands.ListImages,
	}
	cmdListImages.Flags().BoolP(
		"all", "a", false, "Show all images. Only named/taged and leaf images are shown by default.")
	cmdListImages.Flags().BoolP(
		"quiet", "q", false, "Only display numeric IDs")
	cmdListImages.Flags().BoolP(
		"no-header", "n", false, "Omit the header")
	cmdImage.AddCommand(cmdListImages)

	app.AddCommand(cmdImage)

	// images command
	var cmdImages = &cobra.Command{
		Use:   "images [NAME[:TAG]]",
		Short: "Shortcut to list images",
		Long:  appName + " images - List images",
		Run:   commands.ListImages,
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
		Long:  appName + " version - Show the version information",
		Run:   commands.Version,
	}
	app.AddCommand(cmdVersion)

	app.Execute()
}
