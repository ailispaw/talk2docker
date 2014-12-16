package commands

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var cmdConfig = &cobra.Command{
	Use:   "config [command]",
	Short: "Manage the configuration file",
	Long:  appName + " config - Manage the configuration file",
	Run: func(ctx *cobra.Command, args []string) {
		ctx.Usage()
	},
}

var cmdCatConfig = &cobra.Command{
	Use:   "cat",
	Short: "Cat the configuration file",
	Long:  appName + " config cat - Cat the configuration file",
	Run:   catConfig,
}

var cmdEditConfig = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"ed"},
	Short:   "Edit the configuration file",
	Long:    appName + " config edit - Edit the configuration file",
	Run:     editConfig,
}

func init() {
	cmdConfig.AddCommand(cmdCatConfig)
	cmdConfig.AddCommand(cmdEditConfig)
}

func catConfig(ctx *cobra.Command, args []string) {
	path := os.ExpandEnv(configPath)

	cmd := exec.Command("cat", path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
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
