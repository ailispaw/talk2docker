package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/ailispaw/talk2docker/api"
	"github.com/ailispaw/talk2docker/client"
)

var cmdCompose = &cobra.Command{
	Use:     "compose <PATH/TO/YAML> [NAME...]",
	Aliases: []string{"fig", "create"},
	Short:   "Compose containers",
	Long:    APP_NAME + " compose - Compose containers",
	Run:     composeContainers,
}

var cmdComposeContainers = &cobra.Command{
	Use:     "compose <PATH/TO/YAML> [NAME...]",
	Aliases: []string{"fig", "create"},
	Short:   "Compose containers",
	Long:    APP_NAME + " container compose - Compose containers",
	Run:     composeContainers,
}

func init() {
	cmdContainer.AddCommand(cmdComposeContainers)
}

type Composer struct {
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
		ctx.Println("Needs an argument <PATH/TO/YAML> to compose containers")
		ctx.Usage()
		return
	}

	path := os.ExpandEnv(args[0])

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

	inArray := func(a string, list []string) bool {
		for _, b := range list {
			if a == b {
				return true
			}
		}
		return false
	}

	var gotError = false
	for name, composer := range composers {
		if (len(names) == 0) || inArray(name, names) {
			if cid, err := composeContainer(ctx, name, composer); err != nil {
				log.Println(err)
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

func composeContainer(ctx *cobra.Command, name string, composer Composer) (string, error) {
	var (
		config     api.Config
		hostConfig api.HostConfig

		localVolumes   = make(map[string]struct{})
		bindVolumes    []string
		exposedPorts   = make(map[string]struct{})
		portBindings   = make(map[string][]api.PortBinding)
		deviceMappings []api.DeviceMapping
	)

	r, n, t, err := client.ParseRepositoryName(composer.Image)
	if err != nil {
		return "", err
	}
	composer.Image = n + ":" + t
	if r != "" {
		composer.Image = r + "/" + composer.Image
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		return "", err
	}

	if composer.Build != "" {
		message, err := docker.BuildImage(composer.Build, composer.Image, true)
		if err != nil {
			return "", err
		}
		if composer.Image == "" {
			if _, err := fmt.Sscanf(message, "Successfully built %s", &composer.Image); err != nil {
				return "", err
			}
		}
	}

	for _, rawPort := range composer.Ports {
		var (
			hostPort, containerPort string
		)

		if !strings.Contains(rawPort, ":") {
			hostPort = ""
			containerPort = rawPort
		} else {
			parts := strings.Split(rawPort, ":")
			hostPort = parts[0]
			containerPort = parts[1]
		}

		port := fmt.Sprintf("%s/%s", containerPort, "tcp")
		if _, exists := exposedPorts[port]; !exists {
			exposedPorts[port] = struct{}{}
		}

		portBinding := api.PortBinding{
			HostPort: hostPort,
		}
		bslice, exists := portBindings[port]
		if !exists {
			bslice = []api.PortBinding{}
		}
		portBindings[port] = append(bslice, portBinding)
	}

	for _, containerPort := range composer.ExposedPorts {
		port := fmt.Sprintf("%s/%s", containerPort, "tcp")
		if _, exists := exposedPorts[port]; !exists {
			exposedPorts[port] = struct{}{}
		}
	}

	for _, volume := range composer.Volumes {
		if arr := strings.Split(volume, ":"); len(arr) > 1 {
			if arr[1] == "/" {
				return "", fmt.Errorf("Invalid bind mount: destination can't be '/'")
			}
			bindVolumes = append(bindVolumes, volume)
		} else if volume == "/" {
			return "", fmt.Errorf("Invalid volume: path can't be '/'")
		} else {
			localVolumes[volume] = struct{}{}
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
		config.Entrypoint = []string{composer.Entrypoint}
	}
	config.MacAddress = composer.MacAddress

	hostConfig.Binds = bindVolumes
	hostConfig.Privileged = composer.Privileged
	hostConfig.PortBindings = portBindings
	hostConfig.Links = append(composer.Links, composer.ExternalLinks...)
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
	cid, err = docker.CreateContainer(name, config, hostConfig)
	if err != nil {
		if apiErr, ok := err.(api.Error); ok && (apiErr.StatusCode == 404) {
			if _, err := docker.PullImage(config.Image); err != nil {
				return "", err
			}

			cid, err = docker.CreateContainer(name, config, hostConfig)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	return cid, nil
}
