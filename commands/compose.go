package commands

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"

	"github.com/ailispaw/talk2docker/api"
	"github.com/ailispaw/talk2docker/client"
)

var composeFlags = struct {
	Name string // --name

	Ports   []string // --publish
	Volumes []string // --volume

	Hostname     string   // --hostname
	Domainname   string   // --domain
	User         string   // --user
	Memory       int64    // --memory
	MemorySwap   int64    // --memory-swap
	CpuShares    int64    // --cpu-shares
	Cpuset       string   // --cpuset
	ExposedPorts []string // --expose
	Tty          bool     // --tty
	OpenStdin    bool     // --interactive
	Env          []string // --env
	Cmd          []string // --cmd
	WorkingDir   string   // --workdir
	Entrypoint   string   // --entrypoint
	MacAddress   string   // --mac-address

	Privileged      bool     // --privileged
	Links           []string // --link
	PublishAllPorts bool     // --publish-all
	Dns             []string // --dns
	DnsSearch       []string // --dns-search
	ExtraHosts      []string // --add-host
	VolumesFrom     []string // --volumes-from
	Devices         []string // --device
	NetworkMode     string   // --net
	IpcMode         string   // --ipc
	PidMode         string   // --pid
	CapAdd          []string // --cap-add
	CapDrop         []string // --cap-drop
	RestartPolicy   string   // --restart
	SecurityOpt     []string // --security-opt
	ReadonlyRootfs  bool     // --read-only
}{}

var cmdCompose = &cobra.Command{
	Use:     "compose <PATH/TO/YAML> [NAME...]",
	Aliases: []string{"fig", "create"},
	Short:   "Create containers",
	Long:    APP_NAME + " compose - Create containers",
	Run:     composeContainers,
}

var cmdComposeContainers = &cobra.Command{
	Use:     "compose <PATH/TO/YAML> [NAME...]",
	Aliases: []string{"fig", "create"},
	Short:   "Create containers",
	Long:    APP_NAME + " container compose - Create containers",
	Run:     composeContainers,
}

func init() {
	for _, flags := range []*pflag.FlagSet{cmdCompose.Flags(), cmdComposeContainers.Flags()} {
		flags.StringVar(&composeFlags.Name, "name", "", "Override the name of the container")

		flags.StringSliceVarP(&composeFlags.Ports, "publish", "p", nil, "Publish a container's port to the host")
		flags.StringSliceVarP(&composeFlags.Volumes, "volume", "v", nil, "Bind mount volume(s)")

		flags.StringVar(&composeFlags.Hostname, "hostname", "", "Hostname of the container")
		flags.StringVar(&composeFlags.Domainname, "domain", "", "Domain name of the container")
		flags.StringVarP(&composeFlags.User, "user", "u", "", "Username or UID")
		flags.Int64VarP(&composeFlags.Memory, "memory", "m", 0, "Memory limit")
		flags.Int64Var(&composeFlags.MemorySwap, "memory-swap", 0, "Total memory (memory + swap), '-1' to disable swap")
		flags.Int64Var(&composeFlags.CpuShares, "cpu-shares", 0, "CPU shares (relative weight)")
		flags.StringVar(&composeFlags.Cpuset, "cpuset", "", "CPUs in which to allow execution (0-3, 0,1)")
		flags.StringSliceVar(&composeFlags.ExposedPorts, "expose", nil, "Expose a port or a range of ports without publishing")
		flags.BoolVarP(&composeFlags.Tty, "tty", "t", false, "Allocate a pseudo-TTY")
		flags.BoolVarP(&composeFlags.OpenStdin, "interactive", "i", false, "Keep STDIN open even if not attached")
		flags.StringSliceVarP(&composeFlags.Env, "env", "e", nil, "Set environment variable(s)")
		flags.StringSliceVarP(&composeFlags.Cmd, "cmd", "c", nil, "Command line to execute")
		flags.StringVarP(&composeFlags.WorkingDir, "workdir", "w", "", "Working directory inside the container")
		flags.StringVar(&composeFlags.Entrypoint, "entrypoint", "", "Overwrite the default ENTRYPOINT of the image")
		flags.StringVar(&composeFlags.MacAddress, "mac-address", "", "Assign a MAC address to the container")

		flags.BoolVar(&composeFlags.Privileged, "privileged", false, "Give extended privileges to the container")
		flags.StringSliceVar(&composeFlags.Links, "link", nil, "Add link to another container in the form of NAME:ALIAS")
		flags.BoolVar(&composeFlags.PublishAllPorts, "publish-all", false, "Publish all exposed ports to random ports")
		flags.StringSliceVar(&composeFlags.Dns, "dns", nil, "Set custom DNS servers")
		flags.StringSliceVar(&composeFlags.DnsSearch, "dns-search", nil, "Set custom DNS search domains")
		flags.StringSliceVar(&composeFlags.ExtraHosts, "add-host", nil, "Add a custom host-to-IP mapping (host:ip)")
		flags.StringSliceVar(&composeFlags.VolumesFrom, "volumes-from", nil, "Mount volumes from the specified container(s)")
		flags.StringSliceVar(&composeFlags.Devices, "device", nil, "Add a host device to the container")
		flags.StringVar(&composeFlags.NetworkMode, "net", "", "Set the Network mode for the container")
		flags.StringVar(&composeFlags.IpcMode, "ipc", "", "IPC namespace to use")
		flags.StringVar(&composeFlags.PidMode, "pid", "", "PID namespace to use")
		flags.StringSliceVar(&composeFlags.CapAdd, "cap-add", nil, "Add Linux capabilities")
		flags.StringSliceVar(&composeFlags.CapDrop, "cap-drop", nil, "Drop Linux capabilities")
		flags.StringVar(&composeFlags.RestartPolicy, "restart", "", "Restart policy to apply when a container exits (no, on-failure[:MAX-RETRY], always)")
		flags.StringSliceVar(&composeFlags.SecurityOpt, "security-opt", nil, "Security options")
		flags.BoolVar(&composeFlags.ReadonlyRootfs, "read-only", false, "Mount the container's root filesystem as read only")
	}

	cmdContainer.AddCommand(cmdComposeContainers)
}

type Composer struct {
	Name string

	Build string `yaml:"build"`

	Ports   []string `yaml:"ports"`
	Volumes []string `yaml:"volumes"`

	// api.Config
	Hostname     string   `yaml:"hostname"`
	Domainname   string   `yaml:"domainname"`
	User         string   `yaml:"user"`
	Memory       int64    `yaml:"mem_limit"`
	MemorySwap   int64    `yaml:"mem_swap"`
	CpuShares    int64    `yaml:"cpu_shares"`
	Cpuset       string   `yaml:"cpuset"`
	ExposedPorts []string `yaml:"expose"`
	Tty          bool     `yaml:"tty"`
	OpenStdin    bool     `yaml:"stdin_open"`
	Env          []string `yaml:"environment"`
	Cmd          []string `yaml:"command"`
	Image        string   `yaml:"image"`
	WorkingDir   string   `yaml:"working_dir"`
	Entrypoint   string   `yaml:"entrypoint"`
	MacAddress   string   `yaml:"mac_address"`

	// api.HostConfig
	Privileged      bool     `yaml:"privileged"`
	Links           []string `yaml:"links"`
	ExternalLinks   []string `yaml:"external_links"`
	PublishAllPorts bool     `yaml:"publish_all"`
	Dns             []string `yaml:"dns"`
	DnsSearch       []string `yaml:"dns_search"`
	ExtraHosts      []string `yaml:"add_host"`
	VolumesFrom     []string `yaml:"volumes_from"`
	Devices         []string `yaml:"device"`
	NetworkMode     string   `yaml:"net"`
	IpcMode         string   `yaml:"ipc"`
	PidMode         string   `yaml:"pid"`
	CapAdd          []string `yaml:"cap_add"`
	CapDrop         []string `yaml:"cap_drop"`
	RestartPolicy   string   `yaml:"restart"`
	SecurityOpt     []string `yaml:"security_opt"`
	ReadonlyRootfs  bool     `yaml:"read_only"`
}

func composeContainers(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ErrorExit(ctx, "Needs an argument <PATH/TO/YAML> to compose containers")
	}

	path := filepath.Clean(args[0])
	root := filepath.Dir(path)

	var names []string
	if len(args) > 1 {
		names = args[1:]
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var composers map[string]Composer
	if err := yaml.Unmarshal(data, &composers); err != nil {
		log.Fatal(err)
	}

	var gotError = false

	if len(names) == 0 {
		for name, composer := range composers {
			composer.Name = name
			composer = mergeComposeFlags(ctx, composer)
			if cid, err := composeContainer(ctx, root, composer); err != nil {
				log.Error(err)
				gotError = true
			} else {
				ctx.Println(cid)
			}
		}
	}

	for _, name := range names {
		if composer, ok := composers[name]; ok {
			composer.Name = name
			composer = mergeComposeFlags(ctx, composer)
			if cid, err := composeContainer(ctx, root, composer); err != nil {
				log.Error(err)
				gotError = true
			} else {
				ctx.Println(cid)
			}
		}
	}

	if gotError {
		log.Fatal("Error: failed to compose one or more containers")
	}
}

func composeContainer(ctx *cobra.Command, root string, composer Composer) (string, error) {
	var (
		config     api.Config
		hostConfig api.HostConfig

		localVolumes   = make(map[string]struct{})
		bindVolumes    []string
		exposedPorts   = make(map[string]struct{})
		portBindings   = make(map[string][]api.PortBinding)
		links          []string
		deviceMappings []api.DeviceMapping
	)

	if composer.Image != "" {
		r, n, t, err := client.ParseRepositoryName(composer.Image)
		if err != nil {
			return "", err
		}
		composer.Image = n + ":" + t
		if r != "" {
			composer.Image = r + "/" + composer.Image
		}
	}

	if (composer.WorkingDir != "") && !filepath.IsAbs(composer.WorkingDir) {
		return "", fmt.Errorf("Invalid working directory: it must be absolute.")
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		return "", err
	}

	if composer.Build != "" {
		if !filepath.IsAbs(composer.Build) {
			composer.Build = filepath.Join(root, composer.Build)
		}
		message, err := docker.BuildImage(composer.Build, composer.Image, false)
		if err != nil {
			return "", err
		}
		if composer.Image == "" {
			if _, err := fmt.Sscanf(message, "Successfully built %s", &composer.Image); err != nil {
				return "", err
			}
		}
	}

	for _, port := range composer.Ports {
		var (
			rawPort       = port
			hostIp        = ""
			hostPort      = ""
			containerPort = ""
			proto         = "tcp"
		)

		if i := strings.LastIndex(port, "/"); i != -1 {
			proto = strings.ToLower(port[i+1:])
			port = port[:i]
		}

		parts := strings.Split(port, ":")
		switch len(parts) {
		case 1:
			containerPort = parts[0]
		case 2:
			hostPort = parts[0]
			containerPort = parts[1]
		case 3:
			hostIp = parts[0]
			hostPort = parts[1]
			containerPort = parts[2]
		default:
			return "", fmt.Errorf("Invalid port specification: %s", rawPort)
		}

		port := fmt.Sprintf("%s/%s", containerPort, proto)
		if _, exists := exposedPorts[port]; !exists {
			exposedPorts[port] = struct{}{}
		}

		portBinding := api.PortBinding{
			HostIp:   hostIp,
			HostPort: hostPort,
		}
		bslice, exists := portBindings[port]
		if !exists {
			bslice = []api.PortBinding{}
		}
		portBindings[port] = append(bslice, portBinding)
	}

	for _, port := range composer.ExposedPorts {
		var (
			rawPort       = port
			containerPort = ""
			proto         = "tcp"
		)

		parts := strings.Split(containerPort, "/")
		switch len(parts) {
		case 1:
			containerPort = parts[0]
		case 2:
			containerPort = parts[0]
			proto = strings.ToLower(parts[1])
		default:
			return "", fmt.Errorf("Invalid port specification: %s", rawPort)
		}

		port := fmt.Sprintf("%s/%s", containerPort, proto)
		if _, exists := exposedPorts[port]; !exists {
			exposedPorts[port] = struct{}{}
		}
	}

	for _, volume := range composer.Volumes {
		if arr := strings.Split(volume, ":"); len(arr) > 1 {
			if arr[1] == "/" {
				return "", fmt.Errorf("Invalid bind mount: destination can't be '/'")
			}
			if !filepath.IsAbs(arr[0]) {
				return "", fmt.Errorf("Invalid bind mount: the host path must be absolute.")
			}
			bindVolumes = append(bindVolumes, volume)
		} else if volume == "/" {
			return "", fmt.Errorf("Invalid volume: path can't be '/'")
		} else {
			localVolumes[volume] = struct{}{}
		}
	}

	for _, link := range append(composer.Links, composer.ExternalLinks...) {
		arr := strings.Split(link, ":")
		if len(arr) < 2 {
			links = append(links, arr[0]+":"+arr[0])
		} else {
			links = append(links, link)
		}
	}

	for _, device := range composer.Devices {
		src := ""
		dst := ""
		permissions := "rwm"
		arr := strings.Split(device, ":")
		switch len(arr) {
		case 3:
			permissions = arr[2]
			fallthrough
		case 2:
			dst = arr[1]
			fallthrough
		case 1:
			src = arr[0]
		default:
			return "", fmt.Errorf("Invalid device specification: %s", device)
		}

		if dst == "" {
			dst = src
		}

		deviceMapping := api.DeviceMapping{
			PathOnHost:        src,
			PathInContainer:   dst,
			CgroupPermissions: permissions,
		}
		deviceMappings = append(deviceMappings, deviceMapping)
	}

	parts := strings.Split(composer.RestartPolicy, ":")
	restartPolicy := api.RestartPolicy{}
	restartPolicy.Name = parts[0]
	if (restartPolicy.Name == "on-failure") && (len(parts) == 2) {
		count, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", err
		}
		restartPolicy.MaximumRetryCount = count
	}

	config.Hostname = composer.Hostname
	config.Domainname = composer.Domainname
	config.User = composer.User
	config.Memory = composer.Memory
	config.MemorySwap = composer.MemorySwap
	config.CpuShares = composer.CpuShares
	config.Cpuset = composer.Cpuset
	config.ExposedPorts = exposedPorts
	config.Tty = composer.Tty
	config.OpenStdin = composer.OpenStdin
	config.Env = composer.Env
	config.Cmd = composer.Cmd
	config.Image = composer.Image
	config.Volumes = localVolumes
	config.WorkingDir = composer.WorkingDir
	if composer.Entrypoint != "" {
		config.Entrypoint = strings.Split(composer.Entrypoint, ",")
	}
	config.MacAddress = composer.MacAddress

	hostConfig.Binds = bindVolumes
	hostConfig.Privileged = composer.Privileged
	hostConfig.PortBindings = portBindings
	hostConfig.Links = links
	hostConfig.PublishAllPorts = composer.PublishAllPorts
	hostConfig.Dns = composer.Dns
	hostConfig.DnsSearch = composer.DnsSearch
	hostConfig.ExtraHosts = composer.ExtraHosts
	hostConfig.VolumesFrom = composer.VolumesFrom
	hostConfig.Devices = deviceMappings
	hostConfig.NetworkMode = composer.NetworkMode
	hostConfig.IpcMode = composer.IpcMode
	hostConfig.PidMode = composer.PidMode
	hostConfig.CapAdd = composer.CapAdd
	hostConfig.CapDrop = composer.CapDrop
	hostConfig.RestartPolicy = restartPolicy
	hostConfig.SecurityOpt = composer.SecurityOpt
	hostConfig.ReadonlyRootfs = composer.ReadonlyRootfs

	var cid string
	cid, err = docker.CreateContainer(composer.Name, config, hostConfig)
	if err != nil {
		if apiErr, ok := err.(api.Error); ok && (apiErr.StatusCode == 404) {
			if _, err := docker.PullImage(config.Image); err != nil {
				return "", err
			}

			cid, err = docker.CreateContainer(composer.Name, config, hostConfig)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	return cid, nil
}

func mergeComposeFlags(ctx *cobra.Command, composer Composer) Composer {
	if composeFlags.Name != "" {
		composer.Name = composeFlags.Name
	}

	if len(composeFlags.Ports) > 0 {
		composer.Ports = composeFlags.Ports
	}
	if len(composeFlags.Volumes) > 0 {
		composer.Volumes = composeFlags.Volumes
	}

	if composeFlags.Hostname != "" {
		composer.Hostname = composeFlags.Hostname
	}
	if composeFlags.Domainname != "" {
		composer.Domainname = composeFlags.Domainname
	}
	if composeFlags.User != "" {
		composer.User = composeFlags.User
	}
	if composeFlags.Memory != 0 {
		composer.Memory = composeFlags.Memory
	}
	if composeFlags.MemorySwap != 0 {
		composer.MemorySwap = composeFlags.MemorySwap
	}
	if composeFlags.CpuShares != 0 {
		composer.CpuShares = composeFlags.CpuShares
	}
	if composeFlags.Cpuset != "" {
		composer.Cpuset = composeFlags.Cpuset
	}
	if len(composeFlags.ExposedPorts) > 0 {
		composer.ExposedPorts = composeFlags.ExposedPorts
	}
	if ctx.Flags().Lookup("tty").Changed {
		composer.Tty = composeFlags.Tty
	}
	if ctx.Flags().Lookup("interactive").Changed {
		composer.OpenStdin = composeFlags.OpenStdin
	}
	if len(composeFlags.Env) > 0 {
		composer.Env = composeFlags.Env
	}
	if len(composeFlags.Cmd) > 0 {
		composer.Cmd = composeFlags.Cmd
	}
	if composeFlags.WorkingDir != "" {
		composer.WorkingDir = composeFlags.WorkingDir
	}
	if composeFlags.Entrypoint != "" {
		composer.Entrypoint = composeFlags.Entrypoint
	}
	if composeFlags.MacAddress != "" {
		composer.MacAddress = composeFlags.MacAddress
	}

	if ctx.Flags().Lookup("privileged").Changed {
		composer.Privileged = composeFlags.Privileged
	}
	if len(composeFlags.Links) > 0 {
		composer.Links = composeFlags.Links
	}
	if ctx.Flags().Lookup("publish-all").Changed {
		composer.PublishAllPorts = composeFlags.PublishAllPorts
	}
	if len(composeFlags.Dns) > 0 {
		composer.Dns = composeFlags.Dns
	}
	if len(composeFlags.DnsSearch) > 0 {
		composer.DnsSearch = composeFlags.DnsSearch
	}
	if len(composeFlags.ExtraHosts) > 0 {
		composer.ExtraHosts = composeFlags.ExtraHosts
	}
	if len(composeFlags.VolumesFrom) > 0 {
		composer.VolumesFrom = composeFlags.VolumesFrom
	}
	if len(composeFlags.Devices) > 0 {
		composer.Devices = composeFlags.Devices
	}
	if composeFlags.NetworkMode != "" {
		composer.NetworkMode = composeFlags.NetworkMode
	}
	if composeFlags.IpcMode != "" {
		composer.IpcMode = composeFlags.IpcMode
	}
	if composeFlags.PidMode != "" {
		composer.PidMode = composeFlags.PidMode
	}
	if len(composeFlags.CapAdd) > 0 {
		composer.CapAdd = composeFlags.CapAdd
	}
	if len(composeFlags.CapDrop) > 0 {
		composer.CapDrop = composeFlags.CapDrop
	}
	if composeFlags.RestartPolicy != "" {
		composer.RestartPolicy = composeFlags.RestartPolicy
	}
	if len(composeFlags.SecurityOpt) > 0 {
		composer.SecurityOpt = composeFlags.SecurityOpt
	}
	if ctx.Flags().Lookup("read-only").Changed {
		composer.ReadonlyRootfs = composeFlags.ReadonlyRootfs
	}

	return composer
}
