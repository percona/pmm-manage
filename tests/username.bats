#!/usr/bin/env bats

setup() {
    [ -z "$SUT" ] \
        && export SUT='http://127.0.0.1:7777' || :
    [ -z "$URL_PREFIX" ] \
        && export URL_PREFIX='configurator' || :

    export FAIL_OUTPUT='{"code":403,"status":"Forbidden","title":"User name is limited to 255 bytes and may not include colon and hash symbols"}'

    mkdir -p "${BATS_TMPDIR}" || :
}

get_input() {
    local USERNAME=$1
    local PASSWORD=$2
    echo -n "{\"username\": \"${USERNAME}\", \"password\": \"${PASSWORD}\"}" >&1
    echo "input: {\"username\": \"${USERNAME}\", \"password\": \"${PASSWORD}\"}" >&2
}

@test ": symbol in username" {
    run curl \
        -s \
        -X POST \
        --insecure \
        -d "$(get_input 'random:user' 'pass!word')" \
        "${SUT}/${URL_PREFIX}/v1/users"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" = "$FAIL_OUTPUT" ]]
}

@test "# symbol in username" {
    run curl \
        -s \
        -X POST \
        --insecure \
        -d "$(get_input 'random#user' 'pass!word')" \
        "${SUT}/${URL_PREFIX}/v1/users"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" = "$FAIL_OUTPUT" ]]
}

@test "0 lenght in username" {
    run curl \
        -s \
        -X POST \
        --insecure \
        -d "$(get_input '' 'pass!word')" \
        "${SUT}/${URL_PREFIX}/v1/users"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" = "$FAIL_OUTPUT" ]]
}

@test "0 lenght in password" {
    local FAIL_OUTPUT='{"code":403,"status":"Forbidden","title":"Password is limited to 255 bytes"}'

    run curl \
        -s \
        -X POST \
        --insecure \
        -d "$(get_input 'random!user' '')" \
        "${SUT}/${URL_PREFIX}/v1/users"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" = "$FAIL_OUTPUT" ]]
}

@test ">255 lenght in username" {
    run curl \
        -s \
        -X POST \
        --insecure \
        -d "$(get_input '1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111' 'pass!word')" \
        "${SUT}/${URL_PREFIX}/v1/users"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" = "$FAIL_OUTPUT" ]]
}

@test ">255 lenght in password" {
    local FAIL_OUTPUT='{"code":403,"status":"Forbidden","title":"Password is limited to 255 bytes"}'

    run curl \
        -s \
        -X POST \
        --insecure \
        -d "$(get_input 'random!user' '1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111')" \
        "${SUT}/${URL_PREFIX}/v1/users"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" = "$FAIL_OUTPUT" ]]
}

@test "! the first symbol in username" {
    local FAIL_OUTPUT='{"code":403,"status":"Forbidden","title":"User name should start with a letter or number"}'

    run curl \
        -s \
        -X POST \
        --insecure \
        -d "$(get_input '!random!user' 'pass!word')" \
        "${SUT}/${URL_PREFIX}/v1/users"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" = "$FAIL_OUTPUT" ]]
}

@test "! the first symbol in password" {
    local FAIL_OUTPUT='{"code":403,"status":"Forbidden","title":"Password should start with a letter or number"}'

    run curl \
        -s \
        -X POST \
        --insecure \
        -d "$(get_input 'random!user' '!pass!word')" \
        "${SUT}/${URL_PREFIX}/v1/users"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" = "$FAIL_OUTPUT" ]]
}
