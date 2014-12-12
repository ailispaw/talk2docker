package commands

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yungsang/tablewriter"
	"github.com/yungsang/talk2docker/client"
)

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
	Run:     listHosts,
}

var cmdSwitchHost = &cobra.Command{
	Use:     "switch <NAME>",
	Aliases: []string{"sw"},
	Short:   "Switch the default host",
	Long:    appName + " host switch - Switch the default host",
	Run:     switchHost,
}

var cmdAddHost = &cobra.Command{
	Use:   "add <NAME> <URL> [DESCRIPTION]",
	Short: "Add a new host into the config file",
	Long:  appName + " host add - Add a new host into the config",
	Run:   addHost,
}

var cmdRmHost = &cobra.Command{
	Use:   "rm <NAME>",
	Short: "Rmove a host from the config file",
	Long:  appName + " host rm - Rmove a host from the config file",
	Run:   rmHost,
}

var cmdEditHost = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"ed"},
	Short:   "Edit the config file",
	Long:    appName + " host edit - Edit the config file",
	Run:     editHosts,
}

// Define at ps.go
// var boolQuite, boolNoHeader bool

func init() {
	cmdListHosts.Flags().BoolVarP(&boolQuiet, "quiet", "q", false, "Only display numeric IDs")
	cmdListHosts.Flags().BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")

	cmdHost.AddCommand(cmdListHosts)
	cmdHost.AddCommand(cmdSwitchHost)
	cmdHost.AddCommand(cmdAddHost)
	cmdHost.AddCommand(cmdRmHost)
	cmdHost.AddCommand(cmdEditHost)
}

func listHosts(ctx *cobra.Command, args []string) {
	path := os.ExpandEnv(configPath)

	config, err := client.LoadConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	if boolQuiet {
		for _, host := range config.Hosts {
			fmt.Println(host.Name)
		}
		return
	}

	var items [][]string
	for _, host := range config.Hosts {
		out := []string{
			FormatBool(host.Name == config.Default, "*"),
			host.Name,
			host.URL,
			FormatNonBreakingString(host.Description),
			FormatBool(host.TLS, "YES"),
		}
		items = append(items, out)
	}

	var header = []string{
		"",
		"Name",
		"URL",
		"Description",
		"TLS",
	}

	table := tablewriter.NewWriter(os.Stdout)
	if !boolNoHeader {
		table.SetHeader(header)
	}
	table.SetBorder(false)
	table.AppendBulk(items)
	table.Render()
}

func switchHost(ctx *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Needs an argument <NAME> to switch")
		ctx.Usage()
		return
	}

	path := os.ExpandEnv(configPath)

	config, err := client.LoadConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	var name = args[0]

	host, err := config.GetHost(name)
	if err != nil {
		log.Fatal(err)
	}

	config.Default = host.Name

	err = config.SaveConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	listHosts(ctx, args)
}

func addHost(ctx *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("Needs two arguments <NAME> and <URL> at least")
		ctx.Usage()
		return
	}

	name := args[0]
	url := args[1]
	desc := ""
	if len(args) > 2 {
		desc = strings.Join(args[2:], " ")
	}

	path := os.ExpandEnv(configPath)

	config, err := client.LoadConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	_, err = config.GetHost(name)
	if err == nil {
		log.Fatal(fmt.Sprintf("\"%s\" already exists", name))
	}

	newHost := client.Host{
		Name:        name,
		URL:         url,
		Description: desc,
	}

	config.Default = newHost.Name
	config.Hosts = append(config.Hosts, newHost)

	err = config.SaveConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	listHosts(ctx, args)
}

func rmHost(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Needs an argument <NAME> to remove")
		ctx.Usage()
		return
	}

	name := args[0]

	path := os.ExpandEnv(configPath)

	config, err := client.LoadConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	_, err = config.GetHost(name)
	if err != nil {
		log.Fatal(err)
	}

	hosts := []client.Host{}

	for _, host := range config.Hosts {
		if host.Name != name {
			hosts = append(hosts, host)
		}
	}

	config.Hosts = hosts
	if config.Default == name {
		config.Default = hosts[0].Name
	}

	err = config.SaveConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	listHosts(ctx, args)
}

func editHosts(ctx *cobra.Command, args []string) {
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

	listHosts(ctx, args)
}
