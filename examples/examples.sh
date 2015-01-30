#!/bin/sh

pushd `dirname $0` > /dev/null
HERE=`pwd`
popd > /dev/null

cd "${HERE}"
talk2docker="../talk2docker --config=config.yml"

if command -v boot2docker > /dev/null; then
  boot2docker up
  eval ${talk2docker} host switch boot2docker
fi
if command -v vagrant > /dev/null; then
  vagrant up
  eval ${talk2docker} host switch vagrant
fi

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

execute volume list --all
