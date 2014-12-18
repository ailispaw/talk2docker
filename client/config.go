package client

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Default      string        `yaml:"default"`
	Hosts        []Host        `yaml:"hosts"`
	IndexServers []IndexServer `yaml:"indexservers,omitempty"`
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

type IndexServer struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Auth     string `yaml:"auth,omitempty"`
	Email    string `yaml:"email"`
}

const INDEXSERVER = "https://index.docker.io/v1/"

func getDefaultConfig() *Config {
	var config Config

	config.Default = "default"
	config.Hosts = []Host{}

	var host Host
	host.Name = config.Default
	url := os.Getenv("DOCKER_HOST")
	if url == "" {
		url = "unix:///var/run/docker.sock"
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
	err = yaml.Unmarshal(data, &config)
	if err != nil {
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
	return nil, errors.New(fmt.Sprintf("\"%s\" not found in the config", name))
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
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	file.Close()
	return os.Rename(path+".new", path)
}

func (config *Config) GetIndexServer(url string) (*IndexServer, error) {
	if url == "" {
		url = INDEXSERVER
	}
	for _, server := range config.IndexServers {
		if server.URL == url {
			return &server, nil
		}
	}
	return &IndexServer{
		URL: INDEXSERVER,
	}, errors.New(fmt.Sprintf("\"%s\" not found in the config", url))
}

func (config *Config) SetIndexServer(newServer *IndexServer) {
	for i, server := range config.IndexServers {
		if server.URL == newServer.URL {
			config.IndexServers[i] = *newServer
			return
		}
	}
	config.IndexServers = append(config.IndexServers, *newServer)
	return
}

func (config *Config) LogoutIndexServer(url string) {
	if url == "" {
		return
	}
	for i, server := range config.IndexServers {
		if server.URL == url {
			config.IndexServers[i].Auth = ""
			return
		}
	}
	return
}

func (server *IndexServer) Encode(username, password string) string {
	authStr := username + ":" + password
	msg := []byte(authStr)
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(msg)))
	base64.StdEncoding.Encode(encoded, msg)
	return string(encoded)
}

func (server *IndexServer) Decode() (string, string, error) {
	authStr := server.Auth
	decLen := base64.StdEncoding.DecodedLen(len(authStr))
	decoded := make([]byte, decLen)
	authByte := []byte(authStr)
	n, err := base64.StdEncoding.Decode(decoded, authByte)
	if err != nil {
		return "", "", err
	}
	if n > decLen {
		return "", "", fmt.Errorf("Something went wrong decoding auth configuration")
	}
	arr := strings.SplitN(string(decoded), ":", 2)
	if len(arr) != 2 {
		return "", "", fmt.Errorf("Invalid auth configuration")
	}
	password := strings.Trim(arr[1], "\x00")
	return arr[0], password, nil
}
