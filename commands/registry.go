package commands

import (
	"bufio"
	"os"

	"github.com/howeyc/gopass"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yungsang/tablewriter"

	"github.com/ailispaw/talk2docker/api"
	"github.com/ailispaw/talk2docker/client"
)

var cmdRegistry = &cobra.Command{
	Use:     "registry [command]",
	Aliases: []string{"reg"},
	Short:   "Manage registries",
	Long:    APP_NAME + " registry - Manage registries",
	Run: func(ctx *cobra.Command, args []string) {
		ctx.Help()
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
	Use:     "login [REGISTRY]",
	Aliases: []string{"in"},
	Short:   "Log in to a Docker registry",
	Long:    APP_NAME + " registry login - Log in to a Docker registry",
	Run:     loginRegistry,
}

var cmdLogoutRegistry = &cobra.Command{
	Use:     "logout [REGISTRY]",
	Aliases: []string{"out"},
	Short:   "Log out from a Docker registry",
	Long:    APP_NAME + " registry logout - Log out from a Docker registry",
	Run:     logoutRegistry,
}

func init() {
	flags := cmdListRegistries.Flags()
	flags.BoolVarP(&boolQuiet, "quiet", "q", false, "Only display numeric IDs")
	flags.BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")
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
			ctx.Println(registry.Registry)
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
			registry.Registry,
			registry.Username,
			registry.Email,
			FormatBool(registry.Credentials != "", "IN", ""),
		}
		items = append(items, out)
	}

	header := []string{
		"Registry",
		"Username",
		"Email",
		"Logged",
	}

	PrintInTable(ctx.Out(), header, items, 0, tablewriter.ALIGN_DEFAULT)
}

func loginRegistry(ctx *cobra.Command, args []string) {
	reg := client.INDEX_SERVER
	if len(args) > 0 {
		reg = args[0]
	}

	config, err := client.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	ctx.Printf("Log in to a Docker registry at %s\n", reg)

	registry, _ := config.GetRegistry(reg)

	authConfig := api.AuthConfig{
		ServerAddress: registry.Registry,
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

	ctx.Print("Password: ")
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

	credentials, err := docker.Auth(&authConfig)
	if err != nil {
		log.Fatal(err)
	}

	registry.Username = authConfig.Username
	registry.Email = authConfig.Email
	registry.Credentials = credentials

	config.SetRegistry(registry)

	if err := config.SaveConfig(configPath); err != nil {
		log.Fatal(err)
	}

	ctx.Println("Login Succeeded!")

	listRegistries(ctx, args)
}

func logoutRegistry(ctx *cobra.Command, args []string) {
	reg := client.INDEX_SERVER
	if len(args) > 0 {
		reg = args[0]
	}

	config, err := client.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	registry, notFound := config.GetRegistry(reg)
	if (notFound != nil) || (registry.Credentials == "") {
		log.Fatalf("Not logged in to a Docker registry at %s", reg)
	}

	config.LogoutRegistry(reg)

	if err := config.SaveConfig(configPath); err != nil {
		log.Fatal(err)
	}

	ctx.Printf("Removed login credentials for a Docker registry at %s\n\n", reg)

	listRegistries(ctx, args)
}
