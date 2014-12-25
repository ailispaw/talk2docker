package commands

import (
	"github.com/spf13/cobra"
)

const (
	APP_NAME = "Talk2Docker"
)

var (
	configPath string
	hostName   string

	boolAll, boolQuiet, boolNoHeader, boolJSON bool
)

var app = &cobra.Command{
	Use:   "talk2docker",
	Short: "A simple Docker client to talk to Docker daemon",
	Long:  APP_NAME + " - A simple Docker client to talk to Docker daemon",
}

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
	app.AddCommand(cmdRegistry)
	app.AddCommand(cmdConfig)
	app.AddCommand(cmdVersion)

	app.Execute()
}
