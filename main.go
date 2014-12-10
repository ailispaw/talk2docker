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

	// images command
	var cmdImages = &cobra.Command{
		Use:     "ls [NAME[:TAG]]",
		Aliases: []string{"images"},
		Short:   "List images",
		Long:    appName + " ls - List images",
		Run:     commands.ListImages,
	}
	cmdImages.Flags().BoolP(
		"all", "a", false, "Show all images. Only named/taged and leaf images are shown by default.")
	cmdImages.Flags().BoolP(
		"quiet", "q", false, "Only display numeric IDs")
	cmdImages.Flags().BoolP(
		"no-header", "n", false, "Omit the header")
	app.AddCommand(cmdImages)

	// image command
	var cmdImage = &cobra.Command{
		Use:     "image [command]",
		Aliases: []string{"img"},
		Short:   "Manage images",
		Long:    appName + " image - Manage images",
		Run: func(ctx *cobra.Command, args []string) {
			ctx.Usage()
		},
	}

	var cmdListImages = &cobra.Command{
		Use:     "list [NAME[:TAG]]",
		Aliases: []string{"ls"},
		Short:   "List images",
		Long:    appName + " image list - List images",
		Run:     commands.ListImages,
	}
	cmdListImages.Flags().BoolP(
		"all", "a", false, "Show all images. Only named/taged and leaf images are shown by default.")
	cmdListImages.Flags().BoolP(
		"quiet", "q", false, "Only display numeric IDs")
	cmdListImages.Flags().BoolP(
		"no-header", "n", false, "Omit the header")
	cmdImage.AddCommand(cmdListImages)

	app.AddCommand(cmdImage)

	// host command
	var cmdHost = &cobra.Command{
		Use:   "host [command]",
		Short: "Manage hosts",
		Long:  appName + " host - Manage hosts",
		Run: func(ctx *cobra.Command, args []string) {
			ctx.Usage()
		},
	}

	var cmdListHosts = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List hosts",
		Long:    appName + " host list - List hosts",
		Run:     commands.ListHosts,
	}
	cmdListHosts.Flags().BoolP(
		"quiet", "q", false, "Only display numeric IDs")
	cmdListHosts.Flags().BoolP(
		"no-header", "n", false, "Omit the header")
	cmdHost.AddCommand(cmdListHosts)

	var cmdSwitchHost = &cobra.Command{
		Use:     "switch <NAME>",
		Aliases: []string{"sw"},
		Short:   "Switch the default host",
		Long:    appName + " host switch - Switch the default host",
		Run:     commands.SwitchHost,
	}
	cmdHost.AddCommand(cmdSwitchHost)

	var cmdAddHost = &cobra.Command{
		Use:   "add <NAME> <URL> [DESCRIPTION]",
		Short: "Add a new host into the config file",
		Long:  appName + " host add - Add a new host into the config",
		Run:   commands.AddHost,
	}
	cmdHost.AddCommand(cmdAddHost)

	var cmdRmHost = &cobra.Command{
		Use:   "rm <NAME>",
		Short: "Rmove a host from the config file",
		Long:  appName + " host rm - Rmove a host from the config file",
		Run:   commands.RmHost,
	}
	cmdHost.AddCommand(cmdRmHost)

	var cmdEditHost = &cobra.Command{
		Use:     "edit <NAME>",
		Aliases: []string{"ed"},
		Short:   "Edit the config file",
		Long:    appName + " host edit - Edit the config file",
		Run:     commands.EditHosts,
	}
	cmdHost.AddCommand(cmdEditHost)

	app.AddCommand(cmdHost)

	// hosts command
	var cmdHosts = &cobra.Command{
		Use:   "hosts",
		Short: "Shortcut to list hosts",
		Long:  appName + " hosts - List hosts",
		Run:   commands.ListHosts,
	}
	cmdHosts.Flags().BoolP(
		"quiet", "q", false, "Only display numeric IDs")
	cmdHosts.Flags().BoolP(
		"no-header", "n", false, "Omit the header")
	app.AddCommand(cmdHosts)

	// version command
	var cmdVersion = &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Show the version information",
		Long:    appName + " version - Show the version information",
		Run:     commands.Version,
	}
	app.AddCommand(cmdVersion)

	app.Execute()
}
