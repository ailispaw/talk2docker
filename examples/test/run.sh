#!/bin/sh

pushd `dirname $0` > /dev/null
HERE=`pwd`
popd > /dev/null

cd "${HERE}"
talk2docker="../../talk2docker"

export TALK2DOCKER_CONFIG="../config.yml"

execute() {
  command="${talk2docker} ${*}"
  echo "\n\$ ${command}" >&2
  eval "${command}"
  status=$?
  if [ $status -eq 0 ]; then
    echo "\033[38;5;2m===> ${*}: OK\033[39m" >&2
  else
    echo "\033[38;5;9m===> ${*}: NG\033[39m" >&2
  fi
  return $status
}

eval ${talk2docker} host switch default

if command -v vagrant > /dev/null; then
  cd ..
  vagrant up
  cd "${HERE}"
  eval ${talk2docker} host switch vagrant
elif command -v boot2docker > /dev/null; then
  boot2docker up
  eval ${talk2docker} host switch boot2docker
fi

execute config cat

execute version

execute host info

# Cleanup containers
CONTAINERS="$(${talk2docker} container list -aq)"
CONTAINERS=$(echo ${CONTAINERS})
if [ -n "${CONTAINERS}" ]; then
  execute container remove "${CONTAINERS}" --force
fi
execute ps --all

# Cleanup volumes
VOLUMES="$(${talk2docker} volume list -aq)"
VOLUMES=$(echo ${VOLUMES})
if [ -n "${VOLUMES}" ]; then
  execute volume remove "${VOLUMES}"
fi
execute vs --all

# Cleanup images
IMAGES="$(${talk2docker} image list -q)"
IMAGES=$(echo ${IMAGES})
if [ -n "${IMAGES}" ]; then
  execute image remove "${IMAGES}"
fi
execute ls

execute image pull busybox --all
execute image list --all

execute compose compose.yml hello_world
execute container list --all

execute container start hello_world 
execute container list --all

execute build Dockerfile --tag=build1 --verbose
execute image list --all

execute compose compose.yml build2_with_compose --debug
execute image list --all
execute container list --all

execute image history build2:latest

if command -v jq > /dev/null; then
  execute container inspect build2_with_compose --json | jq '.[0].Id'
else
  execute container inspect build2_with_compose
fi

# Export
execute container export hello_world /tmp/docker | tar tv

execute volume list --all

if command -v jq > /dev/null; then
  execute volume inspect hello_world:/tmp/test --json | jq '.[0].Path'
else
  execute volume inspect hello_world:/tmp/test
fi

execute volume export hello_world:/tmp/docker | tar tv

# Commit
execute commit hello_world ailis/busybox:hello_world

execute image list --all

execute image history ailis/busybox:hello_world

# Upload
execute container export hello_world /tmp/wordpress | tar tv

execute container upload ../wordpress hello_world:/tmp/wordpress

execute container export hello_world /tmp/wordpress | tar tv

if command -v jq > /dev/null; then
  execute volume export hello_world:/tmp/test | tar tv

  VOLUME="$(${talk2docker} volume inspect hello_world:/tmp/test --json | jq '.[0].ID')"
  execute volume upload . "${VOLUME}:/" --verbose

  execute volume export hello_world:/tmp/test | tar tv
fi

# Compose with flags
execute compose compose.yml hello_world --name=hello_world2 --publish=8081:8080 --verbose

if command -v jq > /dev/null; then
  execute container inspect hello_world2 --json | jq '.[0].HostConfig.PortBindings'
else
  execute container inspect hello_world2
fi

execute container start hello_world2
execute container list --latest

execute container top hello_world2

execute container pause hello_world2
execute container list --latest

execute container unpause hello_world2
execute container list --latest

execute container restart hello_world2 --time=1
execute container list --latest

execute container stop hello_world2 --time=1
execute container list --latest

execute container start hello_world2
execute container list --latest

execute container kill hello_world2
execute container list --latest

execute container remove hello_world2
execute container list --all
