package api

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

type Port struct {
	IP          string
	PrivatePort int
	PublicPort  int
	Type        string
}

type Image struct {
	Created     int64
	Id          string
	ParentId    string
	RepoTags    []string
	Size        int64
	VirtualSize int64
}

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

type AuthConfig struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	ServerAddress string `json:"serveraddress"`
}

type Version struct {
	ApiVersion string
	Version    string
	GitCommit  string
	GoVersion  string
}
