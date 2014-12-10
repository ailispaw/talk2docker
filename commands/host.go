package commands

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yungsang/tablewriter"
	"github.com/yungsang/talk2docker/client"
)

func ListHosts(ctx *cobra.Command, args []string) {
	path := os.ExpandEnv(GetStringFlag(ctx, "config"))

	config, err := client.LoadConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	if GetBoolFlag(ctx, "quiet") {
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
	if !GetBoolFlag(ctx, "no-header") {
		table.SetHeader(header)
	}
	table.SetBorder(false)
	table.AppendBulk(items)
	table.Render()
}

func SwitchHost(ctx *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Needs an argument <NAME> to switch")
		ctx.Usage()
		return
	}

	path := os.ExpandEnv(GetStringFlag(ctx, "config"))

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

	fmt.Printf("\"%s\" is the default host from now on.\n", config.Default)
}

func AddHost(ctx *cobra.Command, args []string) {
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

	path := os.ExpandEnv(GetStringFlag(ctx, "config"))

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

	ListHosts(ctx, args)
}

func RmHost(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Needs an argument <NAME> to remove")
		ctx.Usage()
		return
	}

	name := args[0]

	path := os.ExpandEnv(GetStringFlag(ctx, "config"))

	config, err := client.LoadConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	_, err = config.GetHost(name)
	if err != nil {
		log.Fatal(fmt.Sprintf("\"%s\" doesn't exist", name))
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

	ListHosts(ctx, args)
}
