# Configuration Parameters for Compose in YAML

## Container Name

The top of the structure is the name of a new container to create.  
It is followed by its parameters in an indented block.

```yaml
<container-name-1>:
	image:
	command:
	.
	.
	.

<container-name-2>:
	build:
	command:
	.
	.
	.
.
.
.
```

## Parameters for a Container

Most parameters are optional as same as the options of `docker create` command, except `image` or `build`.

### image (string)

An image name:tag or ID for the base image to create the container

```yaml
	image: busybox:latest
```

### build (string)

A path to a Dockerfile to create the base image of the container.  
If `image` is specified with `build`, `image` is used as the tag of the base image.

```yaml
	build: Dockerfile
```

### command (array of string), --cmd

Command and its arguments to execute in the container

```yaml
	command: ["bundle", "exec", "thin", "-p", "3000"]
```

### hostname (string), --hostname

Host name for the container

```yaml
	hostname: foo
```

### domainname (string), --domain

Domain name for the container

```yaml
	domainname: bar.com
```

### user (string), --user

Username or UID

```yaml
	user: postgresql
```

### mem_limit (integer), --memory

Memory limit

```yaml
	mem_limit: 1000000000
```

### mem_swap (integer), --memory-swap

Total memory usage (memory + swap)

```yaml
	mem_swap: 2000000000
```

### cpu_shares (integer), --cpu-shares

CPU shares (relative weight)

```yaml
	cpu_shares: 73
```

### cpuset (string), --cpuset

CPUs in which to allow execution

### ports (array of string), --publish

Expose ports. Either specify both ports (HOST:CONTAINER), or just the container port (a random host port will be chosen).

```yaml
	ports:
		- 3000
		- "8000:8000"
		- "49100:22"
```

### expose (array of string), --expose

Expose a port or a range of ports from the container without publishing it to your host. They'll only be accessible to linked services.

```yaml
	expose:
		- 3000
		- 8000
```

### tty (boolean), --tty

Allocate a pseudo-TTY

```yaml
	tty: true
```

### stdin_open (boolean), --interactive

Keep STDIN open even if not attached

```yaml
	stdin_open: true
```

### environment (array of string), --env

Set environment variables

```yaml
	environment:
		- RACK_ENV=development
		- SESSION_SECRET
```

### working_dir (string), --workdir

Working directory inside the container. It must be an absolute path.

```yaml
	working_dir: /code
```

### entrypoint (string), --entrypoint

Overwrite the default ENTRYPOINT of the image

```yaml
	entrypoint: /code/entrypoint.sh
```

### mac_address (string), --mac-address

Container MAC address (e.g. 92:d0:c6:0a:29:33)

```yaml
	mac_address: "92:d0:c6:0a:29:33"
```

### privileged (boolean), --privileged

Give extended privileges to the container

```yaml
	privileged: true
```

### links (array of string), --link

Add link to another container in the form of NAME:ALIAS.  
If ALIAS is omitted, NAME will be used as ALIAS, too.

```yaml
	links:
		- db:database
		- redis
```

### publish_all (boolean), --publish-all

Publish all exposed ports to the host interfaces

```yaml
	publish_all: true
```

### dns (array of string), --dns

Set custom DNS servers

```yaml
	dns:
		- 8.8.8.8
		- 8.8.4.4
```

### dns_search (array of string), --dns-search

Set custom DNS search domains

```yaml
	dns_search:
		- dc1.example.com
		- dc2.example.com
```

### add_host (array of string), --add-host

Add a custom host-to-IP mapping (HOST:IP)

### volumes (array of string), --volume

Mount paths as volumes, optionally specifying a path on the host machine (HOST:CONTAINER), or an access mode (HOST:CONTAINER:ro).  
HOST must be an absolute path.

```yaml
	volumes:
		- /var/lib/mysql
		- /tmp/cache/:/tmp/cache:ro
```


### volumes_from (array of string), --volumes-from

Mount volumes from the specified container(s)

```yaml
	volumes_from:
		- data
```

### device (array of string), --device

Add a host device to the container

```yaml
	device:
		- /dev/sdc:/dev/xvdc:rwm
```

### net (string), --net

Set the Network mode for the container, 'bridge' is the default.

- 'bridge': creates a new network stack for the container on the docker bridge
- 'none': no networking for this container
- 'container:<NAME|ID>': reuses another container network stack
- 'host': use the host network stack inside the container.  
	Note: the host mode gives the container full access to local system services such as D-bus and is therefore considered insecure.
	
```yaml
	net: host
```

### ipc (string), --ipc

Default is to create a private IPC namespace (POSIX SysV IPC) for the container

- 'container:<NAME|ID>': reuses another container shared memory, semaphores and message queues
- 'host': use the host shared memory,semaphores and message queues inside the container.  
	Note: the host mode gives the container full access to local shared memory and is therefore considered insecure.
	
```yaml
	ipc: host
```

### pid (string), --pid

Default is to create a private PID namespace for the container

- 'host': use the host PID namespace inside the container.  
	Note: the host mode gives the container full access to processes on the system and is therefore considered insecure.

```yaml
	pid: host
```

### cap_add (array of string), --cap-add

Add Linux capabilities

```yaml
	cap_add:
		- ALL
```

### cap_drop (array of string), --cap-drop

Drop Linux capabilities

```yaml
	cap_drop:
		- NET_ADMIN
		- SYS_ADMIN
```

### restart (string), --restart

Restart policy to apply when a container exits

 - 'no'
 - 'on-failure[:MAX-RETRY]'
 - 'always'

```yaml
	restart: always
```

### security_opt (array of string), --security-opt

Security Options

### read_only (boolean), --read-only

Mount the container's root filesystem as read only

```yaml
	read_only: true
```
