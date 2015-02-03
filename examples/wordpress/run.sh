#!/bin/sh

pushd `dirname $0` > /dev/null
HERE=`pwd`
popd > /dev/null

cd "${HERE}"
talk2docker="../../talk2docker --config=../config.yml --host=vagrant"

execute() {
  command="${talk2docker} ${*}"
  echo "\n\$ ${command}" >&2
  eval "${command}"
  status=$?
  if [ $status -ne 0 ]; then
    exit $status
  fi
  return $status
}

if command -v vagrant > /dev/null; then
  cd ..
  vagrant up
  cd "${HERE}"
fi

eval ${talk2docker} container remove web db --force

# https://github.com/docker/fig/blob/master/docs/wordpress.md

if [ ! -d wordpress ]; then
  curl https://wordpress.org/latest.tar.gz | tar -xzf -
fi

cp wp-config.php wordpress/
cp router.php wordpress/

execute compose compose.yml db web

execute container start db web
