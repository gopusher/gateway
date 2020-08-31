#!/bin/bash

### auto build dev environment ###

# for gin, see more: https://github.com/codegangsta/gin
GIN="gin"
if [ -n "$GOBIN" ]; then
    GIN="$GOBIN/gin"
fi

port=3000
if [ -n "$1" ] ;then
    port=$1
fi

$GIN -p="$port" -a=8080 -b="runtime/app" -i --all --excludeDir="runtime" --excludeDir=".idea" --buildArgs="-ldflags '-s -w'" --build="./app/api/app/cmd" "start" "-c" "config.yaml"
