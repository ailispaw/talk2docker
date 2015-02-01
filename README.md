# Talk2Docker

Talk2Docker is a simple Docker client to talk to Docker daemon.

Contributions and suggestions would be appreciated, though it's aimed at my learning Go and Docker Remote API.

## Features

- Handle multiple Docker daemons
- Support multiple Dockerfiles to build in the top folder of your project
- Create containers from a YAML file like Docker Compose (formerly fig)
- Display a tree of all images, which Docker deprecates
- Display a history of an image, modeled on Dockerfile
- Show uploaded files to build in verbose mode and debug mode
- Handle volumes, inspired by [cpuguy83/docker-volumes](https://github.com/cpuguy83/docker-volumes)
- Output in JSON or YAML format as well
- Organize commands by category

## Building talk2docker

### Requirements

- [go](http://golang.org/)
- [godep](https://github.com/tools/godep)

### How to build

```
$ git clone https://github.com/ailispaw/talk2docker.git
$ cd talk2docker
$ make
```

## Usage

```
$ talk2docker
Talk2Docker - A simple Docker client to talk to Docker daemon

Usage:
  talk2docker [flags]
  talk2docker [command]

Available Commands:
  ps                               List containers
  ls [NAME[:TAG]]                  List images
  vs                               List volumes
  hosts                            list hosts
  build [PATH/TO/DOCKERFILE]       Build an image from a Dockerfile
  compose <PATH/TO/YAML> [NAME...] Compose containers
  commit <NAME|ID> <NAME[:TAG]>    Create a new image from a container
  version                          Show the version information
  container [command]              Manage containers
  image [command]                  Manage images
  volume [command]                 Manage volumes
  host [command]                   Manage hosts
  registry [command]               Manage registries
  config [command]                 Manage the configuration file
  help [command]                   Help about any command

 Available Flags:
      --config="$HOME/.talk2docker/config": Path to the configuration file
      --debug=false: Print debug messages
  -h, --help=false: help for talk2docker
      --host="": Hostname to use its config (runtime only)
      --json=false: Output in JSON format
  -v, --verbose=false: Print verbose messages
      --yaml=false: Output in YAML format

Use "talk2docker help [command]" for more information about that command.

```

You can find more examples of usage in [examples](https://github.com/ailispaw/talk2docker/tree/master/examples).

### Configuration

Talk2Docker uses a YAML file to configure a connection to Docker daemon.  
It locates `$HOME/.talk2docker/config` by default.
If it doesn't exist, it will be created automatically as below.

```yaml
default: default
hosts:
- name: default
  url: unix:///var/run/docker.sock
```

You can edit/add multiple hosts where Docker daemon runs, as below.

```yaml
default: vagrant
hosts:
- name: default
  url: unix:///var/run/docker.sock
- name: vagrant
  url: tcp://localhost:2375
- name: boot2docker
  url: tcp://192.168.59.103:2376
  description: on boot2docker-vm managed by boot2docker
  tls: true
  tls-ca-cert: /Users/ailis/.boot2docker/certs/boot2docker-vm/ca.pem
  tls-cert: /Users/ailis/.boot2docker/certs/boot2docker-vm/cert.pem
  tls-key: /Users/ailis/.boot2docker/certs/boot2docker-vm/key.pem
  tls-verify: true
```

```
$ talk2docker version
$ talk2docker --host=boot2docker version
```

## References

- https://github.com/docker/docker ([Apache License Version 2.0](https://github.com/docker/docker/blob/master/LICENSE))
- https://docs.docker.com/reference/api/docker_remote_api/
- https://github.com/samalba/dockerclient ([Apache License Version 2.0](https://github.com/samalba/dockerclient/blob/master/LICENSE))
- https://github.com/spf13/cobra ([Apache License Version 2.0](https://github.com/spf13/cobra/blob/master/LICENSE.txt))
- https://gopkg.in/yaml.v2 ([LGPLv3](https://github.com/go-yaml/yaml/blob/v2/LICENSE))
- https://github.com/howeyc/gopass ([MIT License](https://github.com/howeyc/gopass/blob/master/LICENSE.txt))
- https://github.com/sirupsen/logrus ([MIT License](https://github.com/Sirupsen/logrus/blob/master/LICENSE))
- https://github.com/yungsang/tablewriter/tree/talk2docker ([MIT License](https://github.com/olekukonko/tablewriter/blob/master/LICENCE.md))
- https://github.com/cpuguy83/docker-volumes
