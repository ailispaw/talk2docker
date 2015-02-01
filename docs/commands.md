# Global Options
- --config (:=$HOME/.talk2docker/config)  
	Path to the configuration file
- --host  
	Hostname to use its configuration (runtime only)
- --yaml  
	Output in YAML format
- --json  
	Output in JSON format
- --verbose (-v)  
	Print verbose messages
- --debug  
	Print debug messages
- --help (-h)  
	Print help messages about the command

# Commands

### ps (containers)  
Shortcut to `container list` command

### ls (images)  
Shortcut to `image list` command

### vs (volumes)  
Shortcut to `volume list` command

### hosts  
Shortcut to `host list` command

### build  
Shortcut to `image build` command

### compose (fig, create)  
Shortcut to `container compose` command

### commit  
Shortcut to `container commit` command

### version (v)  
Show the version information

### container (ctn)
- compose (fig, create)  
	Create containers from a YAML file like Docker Compose (formerly fig)
- list (ls)  
	List containers
- inspect (ins, info)  
	Show containers' information
- start (up)  
	Start stopped containers
- stop (down)  
	Stop running containers by sending SIGTERM and then SIGKILL after a grace period
- restart  
	Restart running containers
- kill  
	Kill running containers using SIGKILL or a specified signal
- pause (suspend)  
	Pause all processes within containers
- unpause (resume)  
	Unpause all processes within containers
- wait  
	Block until containers stop
- remove (rm)  
	Remove containers
- logs  
	Stream outputs(STDOUT/STDERR) from a container
- diff  
	Show changes on a container's filesystem from the base image
- export  
	Stream the contents of a container as a tar archive to STDOUT
- top (ps)  
	List the running processes of a container
- commit  
		Create a new image from a container

### image (img)
- list (ls)  
	List images
- build  
	Build an image from a Dockerfile
- pull  
	Pull an image from a registry
- tag  
	Tag an image
- history (hist)  
	Show the history of an image
- inspect (ins, info)  
	Show images' information
- push  
	Push an image into a registry
- remove (rm)  
	Remove images
- search  
	Search for images on a registry

### volume (vol)
- list (ls)  
	List volumes
- inspect (ins, info)  
	Show volumes' information
- remove (rm)  
	Remove volumes
- export  
	Stream the contents of a volume as a tar archive to STDOUT

### host (hst)
- list (ls)  
	List hosts
- switch (sw)  
	Switch the default host
- info  
	Show the host's information
- add  
	Add a new host into the configuration file
- remove (rm)  
	Remove a host from the configuration file

### registry (reg)
- list (ls)  
	List registries
- login (in)  
	Log in to a Docker registry
- logout (out)  
	Log out from a Docker registry

### config (cfg)
- cat (ls)  
	Show the contents of the configuration file
- edit (ed)  
	Edit the configuration file

### help  
Print help messages about the command
