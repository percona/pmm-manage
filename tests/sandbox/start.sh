#!/bin/bash

ROOT_DIR=$(pwd -P)
PATH="$PATH:$(dirname $0)"

$(which gsed || which sed) -i "/ssh-key-path/assh-key-owner: $USER" ${ROOT_DIR}/tests/sandbox/config.yml
export TEST_CONFIG=${ROOT_DIR}/tests/sandbox/config.yml

pushd ${ROOT_DIR}
    go test \
        -coverpkg="github.com/percona/pmm-manage/..." \
        -c -tags testrunmain -o ./pmm-configurator.test \
        ./cmd/pmm-configurator
popd

exec ${ROOT_DIR}/pmm-configurator.test \
    -test.run "^TestRunMain$" \
    -test.coverprofile=coverage.txt
