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
- ~~Handle volumes, inspired by [cpuguy83/docker-volumes](https://github.com/cpuguy83/docker-volumes)~~
- Support uploading a file/folder to a containter
- Output in JSON or YAML format as well
- Organize commands by category

## Limitations

- Only work with Docker Remote API v1.16 or later

## Tested with

- boot2docker v1.4.1 (aufs, devicemapper)
- CoreOS Beta 557.2.0 (btrfs)
- CoreOS Alpha 575.0.0 (overlay)
- only-docker v0.8.0 (overlay)
- Rancheros Lite (overlay)
- DockerRoot (overlay)

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

## Cross-compiling talk2docker

### Requirements

- [VirtualBox](https://www.virtualbox.org/)
- [Vagrant](https://www.vagrantup.com/)

### How to build

```
$ git clone https://github.com/ailispaw/talk2docker.git
$ cd talk2docker
$ vagrant up
```

## Usage

```
$ talk2docker
Talk2Docker - Yet Another Docker Client to talk to Docker daemon

Usage:
  talk2docker [flags]
  talk2docker [command]

Available Commands:
  ps          List containers
  ls          List images
  vs          List volumes
  hosts       list hosts
  build       Build an image from a Dockerfile
  compose     Create containers
  commit      Create a new image from a container
  version     Show the version information
  container   Manage containers
  image       Manage images
  host        Manage hosts
  registry    Manage registries
  config      Manage the configuration file
  docker      Execute the original docker cli

Flags:
  -C, --config string   Path to the configuration file (default "$HOME/.talk2docker/config")
  -D, --debug           Print debug messages
  -h, --help            help for talk2docker
  -H, --host string     Docker hostname to use its config (runtime only)
  -J, --json            Output in JSON format
  -V, --verbose         Print verbose messages
  -v, --version         Print version information
  -Y, --yaml            Output in YAML format

Use "talk2docker [command] --help" for more information about a command.
```

You can find more details in [docs](https://github.com/ailispaw/talk2docker/tree/master/docs) and [examples](https://github.com/ailispaw/talk2docker/tree/master/examples).

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
