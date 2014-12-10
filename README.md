# Talk2Docker

Talk2Docker is a simple Docker client to talk to Docker daemon.

Contributions and suggestions would be appreciated, though it's aimed at my learning Go and Docker Remote API.

## Requirments

- [go](http://golang.org/)
- [godep](https://github.com/tools/godep)

## References

- https://docs.docker.com/reference/api/docker_remote_api/
- https://github.com/samalba/dockerclient
- https://github.com/codegangsta/cli
- https://github.com/olekukonko/tablewriter

## How to Build

```
$ git clone https://github.com/YungSang/talk2docker.git
$ make
```

## How to Use

```
$ ./talk2docker
NAME:
   talk2docker - A simple Docker client to talk to Docker daemon

USAGE:
   talk2docker [global options] command [options] [arguments...]

VERSION:
   0.2.0-dev

AUTHOR:
  YungSang - <yungsang@gmail.com>

COMMANDS:
   ps           List containers
   images       List images
   version      Show the version information
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config, -C '$HOME/.talk2docker/config'     Path to the configuration file
   --host, -H                                   Hostname to use its config (runtime only)
   --help, -h                                   show help
   --version, -v                                print the version

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
  host: tcp://localhost:2375
- name: boot2docker
  host: tcp://192.168.59.104:2376
  tls: true
  tls-ca-cert: /Users/yungsang/.boot2docker/certs/boot2docker-vm/ca.pem
  tls-cert: /Users/yungsang/.boot2docker/certs/boot2docker-vm/cert.pem
  tls-key: /Users/yungsang/.boot2docker/certs/boot2docker-vm/key.pem
  tls-verify: true
```

```
$ talk2docker version
$ talk2docker --host boot2docker version
```
