package client

import (
	"errors"
	"fmt"
	"strings"
)

func (config *Config) GetRegistry(reg string) (*Registry, error) {
	for _, registry := range config.Registries {
		if registry.Registry == reg {
			return &registry, nil
		}
	}
	return &Registry{
		Registry: reg,
	}, fmt.Errorf("\"%s\" not found in the config", reg)
}

func (config *Config) SetRegistry(newRegistry *Registry) {
	for i, registry := range config.Registries {
		if registry.Registry == newRegistry.Registry {
			config.Registries[i] = *newRegistry
			return
		}
	}
	config.Registries = append(config.Registries, *newRegistry)
	return
}

func (config *Config) LogoutRegistry(reg string) {
	if reg == "" {
		return
	}
	for i, registry := range config.Registries {
		if registry.Registry == reg {
			config.Registries[i].Credentials = ""
			return
		}
	}
	return
}

func ParseRepositoryName(name string) (string, string, string, error) {
	var (
		reg = ""
		tag = "latest"
	)

	if strings.Contains(name, "://") {
		return "", "", "", errors.New("Invalid repository name with a schema")
	}

	if n := strings.LastIndex(name, ":"); n >= 0 {
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

	reg = names[0]
	name = names[1]
	if strings.Contains(reg, "index.docker.io") {
		return "", name, tag, nil
	}

	return reg, name, tag, nil
}
