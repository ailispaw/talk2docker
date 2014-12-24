package api

import (
	"time"
)

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

type ImageHistory struct {
	Created   int64
	CreatedBy string
	Id        string
	Size      int64
	Tags      []string
}

type ImageInfo struct {
	Id              string
	Parent          string
	Comment         string
	Created         time.Time
	Container       string
	ContainerConfig RunConfig
	DockerVersion   string
	Author          string
	Config          RunConfig
	Architecture    string
	Os              string
	Size            int64
	VirtualSize     int64
	Checksum        string
}

type RunConfig struct {
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

type Version struct {
	ApiVersion string
	Version    string
	GitCommit  string
	GoVersion  string
}
