# Talk2Docker

Talk2Docker is a simple Docker client to talk to Docker daemon.

Contributions and suggestions would be appreciated, though it's aimed at my learning Go and Docker Remote API.

## Features

- Handle multiple Docker deamons
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
  ps                        List containers
  build [PATH]              Build an image from a Dockerfile
  ls [NAME[:TAG]]           List images
  image [command]           Manage images
  host [command]            Manage hosts
  hosts                     Shortcut to list hosts
  registry [command]        Manage registries
  config [command]          Manage the configuration file
  version                   Show the version information
  help [command]            Help about any command

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

### Configuration

Talk2Docker uses a YAML file to configure a connection to Docker daemon.  
It locates `$HOME/.talk2docker/config` by default.
If it does't exist, it will be created automatically as below.  

```yaml
default: default
hosts:
- name: default
  host: unix:///var/run/docker.sock
```

You can edit/add multiple hosts where Docker daemon runs.  

```yaml
default: default
hosts:
- name: default
  url: tcp://localhost:2375
- name: boot2docker
  url: tcp://192.168.59.104:2376
  description: on boot2docker-vm managed by boot2docker
  tls: true
  tls-ca-cert: /Users/ailispaw/.boot2docker/certs/boot2docker-vm/ca.pem
  tls-cert: /Users/ailispaw/.boot2docker/certs/boot2docker-vm/cert.pem
  tls-key: /Users/ailispaw/.boot2docker/certs/boot2docker-vm/key.pem
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
