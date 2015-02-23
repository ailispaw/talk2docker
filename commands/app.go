package commands

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	APP_NAME = "Talk2Docker"
)

var (
	configPath string
	hostName   string

	boolYAML, boolJSON, boolVerbose, boolDebug bool

	boolAll, boolQuiet, boolNoHeader, boolNoTrunc bool
)

var app = &cobra.Command{
	Use:   "talk2docker",
	Short: "Yet Another Docker Client to talk to Docker daemon",
	Long:  APP_NAME + " - Yet Another Docker Client to talk to Docker daemon",
	Run: func(ctx *cobra.Command, args []string) {
		ctx.Help()
	},
}

func init() {
	app.PersistentFlags().StringVarP(&configPath, "config", "C", "$HOME/.talk2docker/config", "Path to the configuration file")
	app.PersistentFlags().StringVarP(&hostName, "host", "H", "", "Docker hostname to use its config (runtime only)")

	app.PersistentFlags().BoolVarP(&boolYAML, "yaml", "Y", false, "Output in YAML format")
	app.PersistentFlags().BoolVarP(&boolJSON, "json", "J", false, "Output in JSON format")

	app.PersistentFlags().BoolVarP(&boolVerbose, "verbose", "V", false, "Print verbose messages")
	app.PersistentFlags().BoolVarP(&boolDebug, "debug", "D", false, "Print debug messages")

	cobra.OnInitialize(Initialize)
}

func Initialize() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.WarnLevel)

	if boolVerbose {
		log.SetFormatter(&LogFormatter{})
		log.SetLevel(log.InfoLevel)
	}
	if boolDebug {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.DebugLevel)
	}
}

func Execute() {
	app.AddCommand(cmdPs)
	app.AddCommand(cmdIs)
	app.AddCommand(cmdVs)
	app.AddCommand(cmdHosts)

	app.AddCommand(cmdBuild)
	app.AddCommand(cmdCompose)
	app.AddCommand(cmdCommit)
	app.AddCommand(cmdVersion)

	app.AddCommand(cmdContainer)
	app.AddCommand(cmdImage)
	app.AddCommand(cmdVolume)
	app.AddCommand(cmdHost)
	app.AddCommand(cmdRegistry)
	app.AddCommand(cmdConfig)

	app.SetOutput(os.Stdout)
	app.Execute()
}

func ErrorExit(ctx *cobra.Command, message string) {
	log.Error(message)
	ctx.Usage()
	os.Exit(1)
}
