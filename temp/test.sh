#!/bin/bash

_gvm="gvm"

function gvm() {
    if [ "$2" == "use" ]; then
        export GOROOT="$GOHOME/$3"
    fi
}
