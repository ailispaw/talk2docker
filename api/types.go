package api

import (
	"time"
)

// https://github.com/docker/docker/blob/master/daemon%2Flist.go#L123
type Container struct {
	Id         string
	Names      []string
	Image      string
	Command    string
	Created    int64
	Status     string
	Ports      []Port
	SizeRw     int64
	SizeRootFs int64
}

// https://github.com/docker/docker/blob/master/daemon%2Fnetwork_settings.go#L38
type Port struct {
	IP          string
	PrivatePort int
	PublicPort  int
	Type        string
}

// https://github.com/docker/docker/blob/master/graph%2Flist.go#L72
type Image struct {
	Created     int64
	Id          string
	ParentId    string
	RepoTags    []string
	Size        int64
	VirtualSize int64
}

// https://github.com/docker/docker/blob/master/graph%2Fhistory.go#L33
type ImageHistory struct {
	Created   int64
	CreatedBy string
	Id        string
	Size      int64
	Tags      []string
}

type ImageHistories []ImageHistory

func (images ImageHistories) Len() int {
	return len(images)
}

func (images ImageHistories) Swap(i, j int) {
	images[i], images[j] = images[j], images[i]
}

func (images ImageHistories) Less(i, j int) bool {
	if images[i].Created == images[j].Created {
		return i > j
	}
	return images[i].Created < images[j].Created
}

type ImageSearchResult struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Official    bool   `json:"is_official"`
	Automated   bool   `json:"is_trusted"`
	Stars       int    `json:"star_count"`
}

type ImageSearchResults []ImageSearchResult

func (images ImageSearchResults) Len() int {
	return len(images)
}

func (images ImageSearchResults) Swap(i, j int) {
	images[i], images[j] = images[j], images[i]
}

type SortImagesByName struct {
	ImageSearchResults
}

func (by SortImagesByName) Less(i, j int) bool {
	return by.ImageSearchResults[i].Name < by.ImageSearchResults[j].Name
}

type SortImagesByStars struct {
	ImageSearchResults
}

func (by SortImagesByStars) Less(i, j int) bool {
	if by.ImageSearchResults[i].Stars == by.ImageSearchResults[j].Stars {
		return by.ImageSearchResults[i].Name > by.ImageSearchResults[j].Name
	}
	return by.ImageSearchResults[i].Stars < by.ImageSearchResults[j].Stars
}

// https://github.com/docker/docker/blob/master/daemon%2Finspect.go#L31
// https://github.com/docker/docker/blob/master/daemon%2Fcontainer.go#L53
type ContainerInfo struct {
	Id              string
	Created         time.Time
	Path            string
	Args            []string
	Config          Config
	State           State
	Image           string
	NetworkSettings NetworkSettings
	ResolvConfPath  string
	HostnamePath    string
	HostsPath       string
	Name            string
	RestartCount    int
	Driver          string
	ExecDriver      string
	MountLabel      string
	ProcessLabel    string
	Volumes         map[string]string
	VolumesRW       map[string]bool
	AppArmorProfile string
	ExecIDs         []string
	HostConfig      HostConfig
}

// https://github.com/docker/docker/blob/master/daemon%2Fstate.go#L12
type State struct {
	Running    bool
	Paused     bool
	Restarting bool
	OOMKilled  bool
	Pid        int
	ExitCode   int
	Error      string
	StartedAt  time.Time
	FinishedAt time.Time
}

// https://github.com/docker/docker/blob/master/daemon%2Fnetwork_settings.go#L11
type NetworkSettings struct {
	IPAddress              string
	IPPrefixLen            int
	MacAddress             string
	LinkLocalIPv6Address   string
	LinkLocalIPv6PrefixLen int
	GlobalIPv6Address      string
	GlobalIPv6PrefixLen    int
	Gateway                string
	IPv6Gateway            string
	Bridge                 string
	Ports                  map[string][]PortBinding
	//PortMapping            map[string]PortMapping // Deprecated
}

// https://github.com/docker/docker/blob/master/graph%2Fservice.go#L140
// https://github.com/docker/docker/blob/master/image%2Fimage.go#L24
type ImageInfo struct {
	Id              string
	Parent          string
	Comment         string
	Created         time.Time
	Container       string
	ContainerConfig Config
	DockerVersion   string
	Author          string
	Config          Config
	Architecture    string
	Os              string
	Size            int64
	VirtualSize     int64
	Checksum        string
}

// https://github.com/docker/docker/blob/master/runconfig%2Fconfig.go#L11
type Config struct {
	Hostname        string
	Domainname      string
	User            string
	Memory          int64
	MemorySwap      int64
	CpuShares       int64
	Cpuset          string
	AttachStderr    bool
	AttachStdin     bool
	AttachStdout    bool
	PortSpecs       []string
	ExposedPorts    map[string]struct{}
	Tty             bool
	OpenStdin       bool
	StdinOnce       bool
	Env             []string
	Cmd             []string
	Image           string
	Volumes         map[string]struct{}
	WorkingDir      string
	Entrypoint      []string
	NetworkDisabled bool
	MacAddress      string
	OnBuild         []string
}

// https://github.com/docker/docker/blob/master/daemon%2Finfo.go#L67
type Info struct {
	ID                 string
	Containers         int
	Images             int
	Driver             string
	DriverStatus       [][]string
	MemoryLimit        int // bool
	SwapLimit          int // bool
	IPv4Forwarding     int // bool
	Debug              int // bool
	NFd                int
	NGoroutines        int
	ExecutionDriver    string
	NEventsListener    int
	KernelVersion      string
	OperatingSystem    string
	IndexServerAddress string
	InitSha1           string
	InitPath           string
	NCPU               int
	MemTotal           int64
	DockerRootDir      string
	Name               string
	Labels             []string
}

// https://github.com/docker/docker/blob/master/builtins%2Fbuiltins.go#L61
type Version struct {
	Version       string
	ApiVersion    string
	GoVersion     string
	GitCommit     string
	Os            string
	KernelVersion string
	Arch          string
}

// https://github.com/docker/docker/blob/master/runconfig%2Fhostconfig.go#L126
type ConfigAndHostConfig struct {
	Config
	HostConfig HostConfig
}

// https://github.com/docker/docker/blob/master/runconfig%2Fhostconfig.go#L101
type HostConfig struct {
	Binds           []string
	ContainerIDFile string
	LxcConf         []KeyValuePair
	Privileged      bool
	PortBindings    map[string][]PortBinding
	Links           []string
	PublishAllPorts bool
	Dns             []string
	DnsSearch       []string
	ExtraHosts      []string
	VolumesFrom     []string
	Devices         []DeviceMapping
	NetworkMode     string
	IpcMode         string
	PidMode         string
	CapAdd          []string
	CapDrop         []string
	RestartPolicy   RestartPolicy
	SecurityOpt     []string
	ReadonlyRootfs  bool
}

// https://github.com/docker/docker/blob/master/utils%2Futils.go#L30
type KeyValuePair struct {
	Key   string
	Value string
}

// https://github.com/docker/docker/blob/master/nat%2Fnat.go#L20
type PortBinding struct {
	HostIp   string
	HostPort string
}

// https://github.com/docker/docker/blob/master/runconfig%2Fhostconfig.go#L90
type DeviceMapping struct {
	PathOnHost        string
	PathInContainer   string
	CgroupPermissions string
}

// https://github.com/docker/docker/blob/master/runconfig%2Fhostconfig.go#L96
type RestartPolicy struct {
	Name              string
	MaximumRetryCount int
}

// https://github.com/docker/docker/blob/master/pkg%2Farchive%2Fchanges.go#L28
type Change struct {
	Path string
	Kind int
}

type Processes struct {
	Titles    []string
	Processes [][]string
}
