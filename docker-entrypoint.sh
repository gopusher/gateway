#!/bin/bash

set -e

# first arg is `-f` or `--some-option`
# set -- 可以让 $@ 变为 gopusher "$@"
if [ "${1#-}" != "$1" ]; then
	set -- gopusher "$@"
fi

exec "$@"
