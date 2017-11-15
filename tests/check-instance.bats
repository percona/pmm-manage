#!/usr/bin/env bats

[ -z "$SUT" ] && SUT='http://127.0.0.1:7777' || :
[ -z "$URL_PREFIX" ] && URL_PREFIX='configurator' || :
[ -z "$INSTANCE_ID" ] && INSTANCE_ID='i-00000000000000000' || :

setup() {
    export USERNAME=user1name
    export OUTPUT="{\"username\":\"${USERNAME}\",\"password\":\"********\"}"

    export PASSWORD1=random-password
    export INPUT1="{\"username\": \"${USERNAME}\", \"password\": \"${PASSWORD1}\", \"instance\": \"${INSTANCE_ID}\"}"
    export INPUT2="{\"username\": \"${USERNAME}\", \"password\": \"${PASSWORD1}\", \"instance\": \"xxx${INSTANCE_ID}\"}"
}

@test "check instance - ok" {
    echo -n $INSTANCE_ID > ${BATS_TEST_DIRNAME}"/sandbox/INSTANCE_ID"
    run curl \
        -s \
        -X POST \
        --insecure \
        -d "{\"InstanceID\": \"$INSTANCE_ID\"}" \
        "${SUT}/${URL_PREFIX}/v1/check-instance"
    echo "$output" >&2
    rm -rf ${BATS_TEST_DIRNAME}"/sandbox/INSTANCE_ID"

    [[ "$status" -eq 0 ]]
    [[ "$output" = '{"code":200,"status":"OK"}' ]]
}

@test "check instance - wrong" {
    echo -n wrong id > ${BATS_TEST_DIRNAME}"/sandbox/INSTANCE_ID"
    run curl \
        -s \
        -X POST \
        --insecure \
        -d "{\"InstanceID\": \"xxx$INSTANCE_ID\"}" \
        "${SUT}/${URL_PREFIX}/v1/check-instance"
    echo "$output" >&2
    rm -rf ${BATS_TEST_DIRNAME}"/sandbox/INSTANCE_ID"

    [[ "$status" -eq 0 ]]
    [[ "$output" = '{"code":403,"status":"Forbidden","title":"Wrong Instance ID"}' ]]
}

@test "check instance - not aws" {
    if [ -n "${REMOTE}" ]; then
        skip "can be checked only locally"
    fi

    run curl \
        -s \
        -X POST \
        --insecure \
        -d "{\"InstanceID\": \"$INSTANCE_ID\"}" \
        "${SUT}/${URL_PREFIX}/v1/check-instance"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" = '{"code":200,"status":"OK"}' ]]
}

@test "create user - ok" {
    echo -n $INSTANCE_ID > ${BATS_TEST_DIRNAME}"/sandbox/INSTANCE_ID"
    run curl \
        -s \
        -X POST \
        -d "${INPUT1}" \
        --user "${USERNAME}:${PASSWORD1}" \
        --insecure \
        "${SUT}/${URL_PREFIX}/v1/users"
    echo "$output" >&2
    rm -rf ${BATS_TEST_DIRNAME}"/sandbox/INSTANCE_ID"

    [[ "$status" -eq 0 ]]
    [[ "$output" = "$OUTPUT" ]]
}

@test "create user - wrong" {
    echo -n wrong id > ${BATS_TEST_DIRNAME}"/sandbox/INSTANCE_ID"
    run curl \
        -s \
        -X POST \
        -d "${INPUT2}" \
        --user "${USERNAME}:${PASSWORD1}" \
        --insecure \
        "${SUT}/${URL_PREFIX}/v1/users"
    echo "$output" >&2
    rm -rf ${BATS_TEST_DIRNAME}"/sandbox/INSTANCE_ID"

    [[ "$status" -eq 0 ]]
    [[ "$output" = '{"code":403,"status":"Forbidden","title":"Wrong Instance ID"}' ]]
}

@test "create user - not aws" {
    if [ -n "${REMOTE}" ]; then
        skip "can be checked only locally"
    fi

    run curl \
        -s \
        -X POST \
        -d "${INPUT1}" \
        --insecure \
        "${SUT}/${URL_PREFIX}/v1/users"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" = "$OUTPUT" ]]
}

@test "delete user" {
    run curl \
        -s \
        -X DELETE \
        --insecure \
        -d '' \
        --user "${USERNAME}:${PASSWORD1}" \
        "${SUT}/${URL_PREFIX}/v1/users/${USERNAME}"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" = '{"code":200,"status":"OK"}' ]]
}
