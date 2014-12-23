package client

import (
	"errors"
	"fmt"
	"strings"
)

func (config *Config) GetRegistry(url string) (*Registry, error) {
	for _, registry := range config.Registries {
		if registry.URL == url {
			return &registry, nil
		}
	}
	return &Registry{
		URL: url,
	}, errors.New(fmt.Sprintf("\"%s\" not found in the config", url))
}

func (config *Config) SetRegistry(newRegistry *Registry) {
	for i, registry := range config.Registries {
		if registry.URL == newRegistry.URL {
			config.Registries[i] = *newRegistry
			return
		}
	}
	config.Registries = append(config.Registries, *newRegistry)
	return
}

func (config *Config) LogoutRegistry(url string) {
	if url == "" {
		return
	}
	for i, registry := range config.Registries {
		if registry.URL == url {
			config.Registries[i].Auth = ""
			return
		}
	}
	return
}

func ParseRepositoryName(name string) (string, string, string, error) {
	var (
		registry = ""
		tag      = "latest"
	)

	if strings.Contains(name, "://") {
		return "", "", "", errors.New("Invalid repository name with a schema")
	}

	n := strings.LastIndex(name, ":")
	if n >= 0 {
		if !strings.Contains(name[n+1:], "/") {
			tag = name[n+1:]
			name = name[:n]
		}
	}

	names := strings.SplitN(name, "/", 2)
	if (len(names) == 1) ||
		(!strings.Contains(names[0], ".") && !strings.Contains(names[0], ":") && (names[0] != "localhost")) {
		return "", name, tag, nil
	}

	registry = names[0]
	name = names[1]
	if strings.Contains(registry, "index.docker.io") {
		return "", name, tag, nil
	}

	return registry, name, tag, nil
}
