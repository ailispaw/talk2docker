# Talk2Docker

Talk2Docker is a simple Docker client to talk to a Docker daemon.

Contributions and suggestions would be appreciated, though it's aimed at my learning Go and Docker Remote API.

# Requirments

- [go](http://golang.org/)
- [godep](https://github.com/tools/godep)

# References

- https://docs.docker.com/reference/api/docker_remote_api/
- https://github.com/samalba/dockerclient
- https://github.com/codegangsta/cli
- https://github.com/olekukonko/tablewriter

# How to Build

```
$ git clone https://github.com/YungSang/talk2docker.git
$ make
```

# How to Use

```
$ ./talk2docker
NAME:
   talk2docker - A simple Docker client to talk to a Docker daemon

USAGE:
   talk2docker [global options] command [command options] [arguments...]

VERSION:
   0.1.0+git

AUTHOR:
  YungSang - <yungsang@gmail.com>

COMMANDS:
   ps   List containers
   images List images
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host, -H 'unix:///var/run/docker.sock'    Location of the Docker socket [$DOCKER_HOST]
   --help, -h                                  show help
   --version, -v                               print the version

```
