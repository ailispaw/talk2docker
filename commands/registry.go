package commands

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"code.google.com/p/gopass"
	"github.com/spf13/cobra"
	"github.com/yungsang/tablewriter"
	"github.com/yungsang/talk2docker/api"
	"github.com/yungsang/talk2docker/client"
)

var cmdRegistry = &cobra.Command{
	Use:     "registry [command]",
	Aliases: []string{"reg"},
	Short:   "Manage registries",
	Long:    appName + " registry - Manage registries",
	Run: func(ctx *cobra.Command, args []string) {
		ctx.Usage()
	},
}

var cmdListRegistries = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List registries",
	Long:    appName + " registry list - List registries",
	Run:     listRegistries,
}

var cmdLoginRegistry = &cobra.Command{
	Use:     "login [SERVER]",
	Aliases: []string{"in"},
	Short:   "Log in to a Docker registry server through the host",
	Long:    appName + " registry login - Log in to a Docker registry server through the host",
	Run:     loginRegistry,
}

var cmdLogoutRegistry = &cobra.Command{
	Use:     "logout [SERVER]",
	Aliases: []string{"out"},
	Short:   "Log out from a Docker registry server through the host",
	Long:    appName + " registry logout - Log out from a Docker registry server through the host",
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
			fmt.Println(registry.URL)
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

	table := tablewriter.NewWriter(os.Stdout)
	if !boolNoHeader {
		table.SetHeader(header)
	} else {
		table.SetBorder(false)
	}
	table.AppendBulk(items)
	table.Render()
}

func loginRegistry(ctx *cobra.Command, args []string) {
	url := client.INDEXSERVER
	if len(args) > 0 {
		url = args[0]
	}

	config, err := client.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Log in to a Docker registry at %s\n", url)

	registry, _ := config.GetRegistry(url)

	authConfig := api.AuthConfig{
		ServerAddress: registry.URL,
	}

	promptDefault := func(prompt string, configDefault string) {
		if configDefault == "" {
			fmt.Printf("%s: ", prompt)
		} else {
			fmt.Printf("%s (%s): ", prompt, configDefault)
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

	authConfig.Password, err = gopass.GetPass("Password: ")

	promptDefault("Email", registry.Email)
	authConfig.Email = readInput()
	if authConfig.Email == "" {
		authConfig.Email = registry.Email
	}

	docker, err := client.NewDockerClient(configPath, hostName)
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

	fmt.Println("Login Succeeded!")
}

func logoutRegistry(ctx *cobra.Command, args []string) {
	url := client.INDEXSERVER
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

	fmt.Printf("Removed login credentials for a Docker registry at %s\n", url)
}
