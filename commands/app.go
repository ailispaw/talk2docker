package commands

import (
	"github.com/spf13/cobra"
)

var appName = "Talk2Docker"

var app = &cobra.Command{
	Use:   "talk2docker",
	Short: "A simple Docker client to talk to Docker daemon",
	Long:  appName + " - A simple Docker client to talk to Docker daemon",
}

var configPath, hostName string

func init() {
	app.PersistentFlags().StringVar(&configPath, "config", "$HOME/.talk2docker/config", "Path to the configuration file")
	app.PersistentFlags().StringVar(&hostName, "host", "", "Hostname to use its config (runtime only)")
}

func Execute() {
	app.AddCommand(cmdPs)
	app.AddCommand(cmdIs)
	app.AddCommand(cmdImage)
	app.AddCommand(cmdHost)
	app.AddCommand(cmdHosts)
	app.AddCommand(cmdConfig)
	app.AddCommand(cmdVersion)

	app.Execute()
}
