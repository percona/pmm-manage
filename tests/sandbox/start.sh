#!/bin/bash

ROOT_DIR=$(pwd -P)

exec $ROOT_DIR/pmm-configurator \
    -ssh-key-owner $USER \
    -config $ROOT_DIR/tests/sandbox/config.yml \
    "$@"
