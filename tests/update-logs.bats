#!/usr/bin/env bats

[ -z "$SUT" ] && SUT='http://127.0.0.1:7777' || :
[ -z "$URL_PREFIX" ] && URL_PREFIX='configurator' || :

@test "prepare" {
    mkdir -p "${BATS_TEST_DIRNAME}/sandbox/log" || :
    rm -rf "${BATS_TEST_DIRNAME}"/sandbox/log/*.log || :
    printf "PID: $$\n" \
        > "${BATS_TEST_DIRNAME}/sandbox/log/00-running__0000-00-00T00:00:00.log"
    printf "PID: 77777\nlocalhost                  : ok=18   changed=10   unreachable=0    failed=0\n" \
        > "${BATS_TEST_DIRNAME}/sandbox/log/01-success__0000-00-00T00:00:01.log"
    printf "PID: 77777\nlocalhost                  : ok=18   changed=10   unreachable=0    failed=1\n" \
        > "${BATS_TEST_DIRNAME}/sandbox/log/02-failed__0000-00-00T00:00:02.log"
}

@test "list updates" {
    if [ -n "${REMOTE}" ]; then
        skip "can be checked only locally"
    fi

    run curl \
        -s \
        -X GET \
        --insecure \
        "${SUT}/${URL_PREFIX}/v1/updates"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" = '["0000-00-00T00:00:00","0000-00-00T00:00:01","0000-00-00T00:00:02"]' ]]
}

@test "get running" {
    if [ -n "${REMOTE}" ]; then
        skip "can be checked only locally"
    fi

    printf "PID: $$\n" \
        > "${BATS_TEST_DIRNAME}/sandbox/log/00-running__0000-00-00T00:00:00.log"

    run curl \
        -s \
        -X GET \
        --insecure \
        "${SUT}/${URL_PREFIX}/v1/updates/0000-00-00T00:00:00"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" =~ '"title":"running"' ]]
}

@test "get succeeded" {
    if [ -n "${REMOTE}" ]; then
        skip "can be checked only locally"
    fi

    run curl \
        -s \
        -X GET \
        --insecure \
        "${SUT}/${URL_PREFIX}/v1/updates/0000-00-00T00:00:01"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" =~ '"title":"succeeded"' ]]
}

@test "get failed" {
    if [ -n "${REMOTE}" ]; then
        skip "can be checked only locally"
    fi

    run curl \
        -s \
        -X GET \
        --insecure \
        "${SUT}/${URL_PREFIX}/v1/updates/0000-00-00T00:00:02"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" =~ '"title":"failed"' ]]
}

@test "cleanup" {
    rm -rf "${BATS_TEST_DIRNAME}"/sandbox/log/*.log || :
    rmdir "${BATS_TEST_DIRNAME}/sandbox/log"
}
