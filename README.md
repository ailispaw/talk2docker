# Talk2Docker

Talk2Docker is a simple Docker client to talk to Docker daemon.

Contributions and suggestions would be appreciated, though it's aimed at my learning Go and Docker Remote API.

## Requirements

- [go](http://golang.org/)
- [godep](https://github.com/tools/godep)

## References

- https://docs.docker.com/reference/api/docker_remote_api/
- https://github.com/samalba/dockerclient
- https://github.com/spf13/cobra
- https://github.com/olekukonko/tablewriter

## How to Build

```
$ git clone https://github.com/YungSang/talk2docker.git
$ make
```

## How to Use

```
$ talk2docker
Talk2Docker - A simple Docker client to talk to Docker daemon

Usage:
  talk2docker [command]

Available Commands:
  ps                        List containers
  ls [NAME[:TAG]]           List images
  image [command]           Manage images
  host [command]            Manage hosts
  hosts                     Shortcut to list hosts
  version                   Show the version information
  help [command]            Help about any command

 Available Flags:
      --config="$HOME/.talk2docker/config": Path to the configuration file
  -h, --help=false: help for talk2docker
      --host="": Hostname to use its config (runtime only)

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
  tls-ca-cert: /Users/yungsang/.boot2docker/certs/boot2docker-vm/ca.pem
  tls-cert: /Users/yungsang/.boot2docker/certs/boot2docker-vm/cert.pem
  tls-key: /Users/yungsang/.boot2docker/certs/boot2docker-vm/key.pem
  tls-verify: true
```

```
$ talk2docker version
$ talk2docker --host=boot2docker version
```
