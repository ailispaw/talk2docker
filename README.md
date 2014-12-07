# Talk2Docker

Talk2Docker is a simple Docker client to talk to a Docker daemon.

It's aimed at my learning Go and Docker Remote API.

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
   ps		List containers
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host, -H 'unix:///var/run/docker.sock'	Location of the Docker socket [$DOCKER_HOST]
   --help, -h					show help
   --version, -v				print the version

$ ./talk2docker ps -h
NAME:
   ps - List containers

USAGE:
   command ps [command options] [arguments...]

OPTIONS:
   --all, -a	Show all containers. Only running containers are shown by default.
   --latest, -l	Show only the latest created container, include non-running ones.
   --quiet, -q	Only display numeric IDs
   --size, -s	Display sizes

$ ./talk2docker ps -a
       ID      |          NAMES          |          IMAGE          |      COMMAND      |       CREATED       |         STATUS         |        PORTS
+--------------+-------------------------+-------------------------+-------------------+---------------------+------------------------+----------------------+
  9daf160d6cf6 | compassionate_engelbart | yungsang/busybox:latest | sh                | 2014-12-06 13:27:43 | Up 4 hours             | 0.0.0.0:8080->80/tcp
  330ad99dcbf2 | condescending_wilson    | yungsang/busybox:latest | sh                | 2014-12-06 13:27:22 | Exited (0) 4 hours ago |
```
