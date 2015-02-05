package commands

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/yungsang/tablewriter"

	"github.com/ailispaw/talk2docker/client"
)

var cmdHosts = &cobra.Command{
	Use:   "hosts",
	Short: "list hosts",
	Long:  APP_NAME + " hosts - List hosts",
	Run:   listHosts,
}

var cmdHost = &cobra.Command{
	Use:     "host [command]",
	Aliases: []string{"hst"},
	Short:   "Manage hosts",
	Long:    APP_NAME + " host - Manage hosts",
	Run: func(ctx *cobra.Command, args []string) {
		ctx.Help()
	},
}

var cmdListHosts = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List hosts",
	Long:    APP_NAME + " host list - List hosts",
	Run:     listHosts,
}

var cmdSwitchHost = &cobra.Command{
	Use:     "switch <NAME>",
	Aliases: []string{"sw"},
	Short:   "Switch the default host",
	Long:    APP_NAME + " host switch - Switch the default host",
	Run:     switchHost,
}

var cmdGetHostInfo = &cobra.Command{
	Use:   "info [NAME]",
	Short: "Show the host's information",
	Long:  APP_NAME + " host info - Show the host's information",
	Run:   getHostInfo,
}

var cmdAddHost = &cobra.Command{
	Use:   "add <NAME> <URL> [DESCRIPTION]",
	Short: "Add a new host into the configuration file",
	Long:  APP_NAME + " host add - Add a new host into the configuration file",
	Run:   addHost,
}

var cmdRemoveHost = &cobra.Command{
	Use:     "remove <NAME>",
	Aliases: []string{"rm", "delete", "del"},
	Short:   "Remove a host from the configuration file",
	Long:    APP_NAME + " host remove - Remove a host from the configuration file",
	Run:     removeHost,
}

func init() {
	for _, flags := range []*pflag.FlagSet{cmdHosts.Flags(), cmdListHosts.Flags()} {
		flags.BoolVarP(&boolQuiet, "quiet", "q", false, "Only display host names")
		flags.BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")
	}

	cmdHost.AddCommand(cmdListHosts)

	cmdHost.AddCommand(cmdSwitchHost)

	flags := cmdGetHostInfo.Flags()
	flags.BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")
	cmdHost.AddCommand(cmdGetHostInfo)

	cmdHost.AddCommand(cmdAddHost)

	cmdHost.AddCommand(cmdRemoveHost)
}

func listHosts(ctx *cobra.Command, args []string) {
	config, err := client.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if boolQuiet {
		for _, host := range config.Hosts {
			ctx.Println(host.Name)
		}
		return
	}

	if boolYAML || boolJSON {
		if err := FormatPrint(ctx.Out(), config.Hosts); err != nil {
			log.Fatal(err)
		}
		return
	}

	var items [][]string
	for _, host := range config.Hosts {
		out := []string{
			FormatBool(host.Name == config.Default, "*", ""),
			host.Name,
			host.URL,
			FormatNonBreakingString(host.Description),
			FormatBool(host.TLS, "YES", ""),
		}
		items = append(items, out)
	}

	header := []string{
		"",
		"Name",
		"URL",
		"Description",
		"TLS",
	}

	PrintInTable(ctx.Out(), header, items, 0, tablewriter.ALIGN_DEFAULT)
}

func switchHost(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument <NAME> to switch")
	}

	name := args[0]

	config, err := client.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	host, err := config.GetHost(name)
	if err != nil {
		log.Fatal(err)
	}

	config.Default = host.Name

	if err := config.SaveConfig(configPath); err != nil {
		log.Fatal(err)
	}

	listHosts(ctx, args)
}

func getHostInfo(ctx *cobra.Command, args []string) {
	if len(args) > 0 {
		hostName = args[0]
	}

	config, err := client.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	host, err := config.GetHost(hostName)
	if err != nil {
		log.Fatal(err)
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	info, err := docker.Info()
	if err != nil {
		log.Fatal(err)
	}

	if boolYAML || boolJSON {
		data := make([]interface{}, 2)
		data[0] = host
		data[1] = info
		if err := FormatPrint(ctx.Out(), data); err != nil {
			log.Fatal(err)
		}
		return
	}

	var items [][]string

	items = append(items, []string{
		"Host", host.Name,
	})
	items = append(items, []string{
		"URL", host.URL,
	})
	items = append(items, []string{
		"Description", FormatNonBreakingString(host.Description),
	})
	items = append(items, []string{
		"TLS", FormatBool(host.TLS, "Supported", "No"),
	})
	if host.TLS {
		items = append(items, []string{
			FormatNonBreakingString("  CA Certificate file"), FormatNonBreakingString(host.TLSCaCert),
		})
		items = append(items, []string{
			FormatNonBreakingString("  Certificate file"), FormatNonBreakingString(host.TLSCert),
		})
		items = append(items, []string{
			FormatNonBreakingString("  Key file"), FormatNonBreakingString(host.TLSKey),
		})
		items = append(items, []string{
			FormatNonBreakingString("  Verify"), FormatBool(host.TLSVerify, "Required", "No"),
		})
	}

	items = append(items, []string{
		"Containers", strconv.Itoa(info.Containers),
	})
	items = append(items, []string{
		"Images", strconv.Itoa(info.Images),
	})
	items = append(items, []string{
		"Storage Driver", info.Driver,
	})
	for _, pair := range info.DriverStatus {
		items = append(items, []string{
			FormatNonBreakingString(fmt.Sprintf("  %s", pair[0])), FormatNonBreakingString(fmt.Sprintf("%s", pair[1])),
		})
	}
	items = append(items, []string{
		"Execution Driver", info.ExecutionDriver,
	})
	items = append(items, []string{
		"Kernel Version", info.KernelVersion,
	})
	items = append(items, []string{
		"Operating System", FormatNonBreakingString(info.OperatingSystem),
	})
	items = append(items, []string{
		"CPUs", strconv.Itoa(info.NCPU),
	})
	items = append(items, []string{
		"Total Memory", fmt.Sprintf("%s GB", FormatFloat(float64(info.MemTotal)/1000000000)),
	})

	items = append(items, []string{
		"Index Server Address", info.IndexServerAddress,
	})

	items = append(items, []string{
		"Memory Limit", FormatBool(info.MemoryLimit != 0, "Supported", "No"),
	})
	items = append(items, []string{
		"Swap Limit", FormatBool(info.SwapLimit != 0, "Supported", "No"),
	})
	items = append(items, []string{
		"IPv4 Forwarding", FormatBool(info.IPv4Forwarding != 0, "Enabled", "Disabled"),
	})

	items = append(items, []string{
		"ID", info.ID,
	})
	items = append(items, []string{
		"Name", info.Name,
	})
	items = append(items, []string{
		"Labels", FormatNonBreakingString(strings.Join(info.Labels, ", ")),
	})

	items = append(items, []string{
		"Debug Mode", FormatBool(info.Debug != 0, "Yes", "No"),
	})
	if info.Debug != 0 {
		items = append(items, []string{
			FormatNonBreakingString("  Events Listeners"), strconv.Itoa(info.NEventsListener),
		})
		items = append(items, []string{
			FormatNonBreakingString("  Fds"), strconv.Itoa(info.NFd),
		})
		items = append(items, []string{
			FormatNonBreakingString("  Goroutines"), strconv.Itoa(info.NGoroutines),
		})

		items = append(items, []string{
			FormatNonBreakingString("  Init Path"), info.InitPath,
		})
		items = append(items, []string{
			FormatNonBreakingString("  Init SHA1"), info.InitSha1,
		})
		items = append(items, []string{
			FormatNonBreakingString("  Docker Root Dir"), info.DockerRootDir,
		})
	}

	PrintInTable(ctx.Out(), nil, items, 0, tablewriter.ALIGN_LEFT)
}

func addHost(ctx *cobra.Command, args []string) {
	if len(args) < 2 {
		ErrorExit(ctx, "Needs two arguments <NAME> and <URL> at least")
	}

	var (
		name = args[0]
		url  = args[1]
		desc = ""
	)

	if len(args) > 2 {
		desc = strings.Join(args[2:], " ")
	}

	config, err := client.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := config.GetHost(name); err == nil {
		log.Fatalf("\"%s\" already exists", name)
	}

	newHost := client.Host{
		Name:        name,
		URL:         url,
		Description: desc,
	}

	config.Default = newHost.Name
	config.Hosts = append(config.Hosts, newHost)

	if err := config.SaveConfig(configPath); err != nil {
		log.Fatal(err)
	}

	listHosts(ctx, args)
}

func removeHost(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument <NAME> to remove")
	}

	name := args[0]

	config, err := client.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if config.Default == name {
		log.Fatal("You can't remove the default host.")
	}

	if _, err := config.GetHost(name); err != nil {
		log.Fatal(err)
	}

	hosts := []client.Host{}

	for _, host := range config.Hosts {
		if host.Name != name {
			hosts = append(hosts, host)
		}
	}

	config.Hosts = hosts

	if err := config.SaveConfig(configPath); err != nil {
		log.Fatal(err)
	}

	listHosts(ctx, args)
}
