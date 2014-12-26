package commands

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"github.com/yungsang/tablewriter"
	"github.com/yungsang/talk2docker/api"
	"github.com/yungsang/talk2docker/client"
)

var cmdRegistry = &cobra.Command{
	Use:     "registry [command]",
	Aliases: []string{"reg"},
	Short:   "Manage registries",
	Long:    APP_NAME + " registry - Manage registries",
	Run: func(ctx *cobra.Command, args []string) {
		ctx.Usage()
	},
}

var cmdListRegistries = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List registries",
	Long:    APP_NAME + " registry list - List registries",
	Run:     listRegistries,
}

var cmdLoginRegistry = &cobra.Command{
	Use:     "login [SERVER]",
	Aliases: []string{"in"},
	Short:   "Log in to a Docker registry server through the host",
	Long:    APP_NAME + " registry login - Log in to a Docker registry server through the host",
	Run:     loginRegistry,
}

var cmdLogoutRegistry = &cobra.Command{
	Use:     "logout [SERVER]",
	Aliases: []string{"out"},
	Short:   "Log out from a Docker registry server through the host",
	Long:    APP_NAME + " registry logout - Log out from a Docker registry server through the host",
	Run:     logoutRegistry,
}

func init() {
	cmdListRegistries.Flags().BoolVarP(&boolQuiet, "quiet", "q", false, "Only display numeric IDs")
	cmdListRegistries.Flags().BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")

	cmdRegistry.AddCommand(cmdListRegistries)
	cmdRegistry.AddCommand(cmdLoginRegistry)
	cmdRegistry.AddCommand(cmdLogoutRegistry)
}

func listRegistries(ctx *cobra.Command, args []string) {
	config, err := client.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if boolQuiet {
		for _, registry := range config.Registries {
			ctx.Println(registry.URL)
		}
		return
	}

	if boolYAML || boolJSON {
		if err := FormatPrint(ctx.Out(), config.Registries); err != nil {
			log.Fatal(err)
		}
		return
	}

	var items [][]string
	for _, registry := range config.Registries {
		out := []string{
			registry.URL,
			registry.Username,
			registry.Email,
			FormatBool(registry.Auth != "", "IN", ""),
		}
		items = append(items, out)
	}

	header := []string{
		"URL",
		"Username",
		"Email",
		"Logged",
	}

	PrintInTable(ctx.Out(), header, items, 0, tablewriter.ALIGN_DEFAULT)
}

func loginRegistry(ctx *cobra.Command, args []string) {
	url := client.INDEX_SERVER
	if len(args) > 0 {
		url = args[0]
	}

	config, err := client.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	ctx.Printf("Log in to a Docker registry at %s\n", url)

	registry, _ := config.GetRegistry(url)

	authConfig := api.AuthConfig{
		ServerAddress: registry.URL,
	}

	promptDefault := func(prompt string, configDefault string) {
		if configDefault == "" {
			ctx.Printf("%s: ", prompt)
		} else {
			ctx.Printf("%s (%s): ", prompt, configDefault)
		}
	}

	readInput := func() string {
		reader := bufio.NewReader(os.Stdin)
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		return string(line)
	}

	promptDefault("Username", registry.Username)
	authConfig.Username = readInput()
	if authConfig.Username == "" {
		authConfig.Username = registry.Username
	}

	ctx.Printf("Password: ")
	authConfig.Password = string(gopass.GetPasswdMasked())

	promptDefault("Email", registry.Email)
	authConfig.Email = readInput()
	if authConfig.Email == "" {
		authConfig.Email = registry.Email
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	err = docker.Auth(&authConfig)
	if err != nil {
		log.Fatal(err)
	}

	registry.Username = authConfig.Username
	registry.Email = authConfig.Email
	registry.Auth = authConfig.Encode()

	config.SetRegistry(registry)

	err = config.SaveConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	ctx.Println("Login Succeeded!")
}

func logoutRegistry(ctx *cobra.Command, args []string) {
	url := client.INDEX_SERVER
	if len(args) > 0 {
		url = args[0]
	}

	config, err := client.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	registry, notFound := config.GetRegistry(url)
	if (notFound != nil) || (registry.Auth == "") {
		log.Fatal(fmt.Sprintf("Not logged in to a Docker registry at %s", url))
	}

	config.LogoutRegistry(url)

	err = config.SaveConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	ctx.Printf("Removed login credentials for a Docker registry at %s\n", url)
}
