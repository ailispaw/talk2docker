package commands

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ailispaw/talk2docker/client"
)

var cmdConfig = &cobra.Command{
	Use:     "config [command]",
	Aliases: []string{"cfg"},
	Short:   "Manage the configuration file",
	Long:    APP_NAME + " config - Manage the configuration file",
	Run: func(ctx *cobra.Command, args []string) {
		ctx.Help()
	},
}

var cmdCatConfig = &cobra.Command{
	Use:     "cat",
	Aliases: []string{"ls"},
	Short:   "Cat the configuration file",
	Long:    APP_NAME + " config cat - Cat the configuration file",
	Run:     catConfig,
}

var cmdEditConfig = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"ed"},
	Short:   "Edit the configuration file",
	Long:    APP_NAME + " config edit - Edit the configuration file",
	Run:     editConfig,
}

func init() {
	cmdConfig.AddCommand(cmdCatConfig)

	cmdConfig.AddCommand(cmdEditConfig)
}

func catConfig(ctx *cobra.Command, args []string) {
	config, err := client.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := FormatPrint(ctx.Out(), config); err != nil {
		log.Fatal(err)
	}
}

func editConfig(ctx *cobra.Command, args []string) {
	path := os.ExpandEnv(configPath)

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
