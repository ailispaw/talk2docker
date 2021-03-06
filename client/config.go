package client

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Default    string     `yaml:"default"`
	Hosts      []Host     `yaml:"hosts"`
	Registries []Registry `yaml:"registries,omitempty"`
}

type Host struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	Description string `yaml:"description,omitempty"`
	TLS         bool   `yaml:"tls,omitempty"`
	TLSCaCert   string `yaml:"tls-ca-cert,omitempty"`
	TLSCert     string `yaml:"tls-cert,omitempty"`
	TLSKey      string `yaml:"tls-key,omitempty"`
	TLSVerify   bool   `yaml:"tls-verify,omitempty"`
}

type Registry struct {
	Registry    string `yaml:"registry"`
	Username    string `yaml:"username"`
	Email       string `yaml:"email"`
	Credentials string `yaml:"credentials,omitempty"`
}

const (
	INDEX_SERVER  = "https://index.docker.io/v1/"
	DOCKER_SOCKET = "unix:///var/run/docker.sock"
)

func getDefaultConfig() *Config {
	var config Config

	config.Default = "default"
	config.Hosts = []Host{}

	var host Host
	host.Name = config.Default
	url := os.Getenv("DOCKER_HOST")
	if url == "" {
		url = DOCKER_SOCKET
	}
	host.URL = url

	certPath := os.Getenv("DOCKER_CERT_PATH")
	if certPath != "" {
		host.TLS = true
		host.TLSCaCert = filepath.Join(certPath, "ca.pem")
		host.TLSCert = filepath.Join(certPath, "cert.pem")
		host.TLSKey = filepath.Join(certPath, "key.pem")
		host.TLSVerify = true
	}

	config.Hosts = append(config.Hosts, host)

	return &config
}

func LoadConfig(path string) (*Config, error) {
	path = os.ExpandEnv(path)

	data, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		config := getDefaultConfig()
		return config, config.SaveConfig(path)
	}
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (config *Config) GetDefaultHost() (*Host, error) {
	return config.GetHost(config.Default)
}

func (config *Config) GetHost(name string) (*Host, error) {
	if name == "" {
		name = config.Default
	}
	for _, host := range config.Hosts {
		if host.Name == name {
			return &host, nil
		}
	}
	return nil, fmt.Errorf("\"%s\" not found in the config", name)
}

func (config *Config) SaveConfig(path string) error {
	path = os.ExpandEnv(path)

	os.Remove(path + ".new")
	os.Mkdir(filepath.Dir(path), 0700)
	file, err := os.Create(path + ".new")
	if err != nil {
		return err
	}

	defer file.Close()
	defer os.Remove(path + ".new")

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		return err
	}

	os.Remove(path + ".bak")
	if err := os.Link(path, path+".bak"); err != nil && !os.IsNotExist(err) {
		return err
	}

	file.Close()
	os.Remove(path)
	return os.Rename(path+".new", path)
}
