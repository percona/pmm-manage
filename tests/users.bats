#!/usr/bin/env bats

[ -z "$SUT" ] && SUT='http://127.0.0.1:7777' || :
[ -z "$URL_PREFIX" ] && URL_PREFIX='configurator' || :

setup() {
    export USERNAME=user1name
    export OUTPUT="{\"username\":\"${USERNAME}\",\"password\":\"********\"}"

    export PASSWORD1=random-password
    export PASSWORD2=pass1word
    export INPUT1="{\"username\": \"${USERNAME}\", \"password\": \"${PASSWORD1}\"}"
    export INPUT2="{\"username\": \"${USERNAME}\", \"password\": \"${PASSWORD2}\"}"

    mkdir -p "${BATS_TMPDIR}" || :
}

@test "create user" {
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

@test "get user" {
    run curl \
        -s \
        -X GET \
        --insecure \
        -d '' \
        --user "${USERNAME}:${PASSWORD1}" \
        "${SUT}/${URL_PREFIX}/v1/users/${USERNAME}"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" = "$OUTPUT" ]]
}

@test "get all users" {
    run curl \
        -s \
        -X GET \
        --insecure \
        -d '' \
        --user "${USERNAME}:${PASSWORD1}" \
        "${SUT}/${URL_PREFIX}/v1/users"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" =~ "$OUTPUT" ]]
}

@test "check grafana user 1" {
    if [ -n "${REMOTE}" ]; then
        skip "can be checked only locally"
    fi

    run sqlite3 "${BATS_TEST_DIRNAME}/sandbox/grafana.db" "SELECT salt FROM user WHERE login='${USERNAME}';"
    [[ "$status" -eq 0 ]]
    local SALT="$output"
    echo "SALT: $SALT" >&2

    run php "${BATS_TEST_DIRNAME}/password-helper.php" "${PASSWORD1}" "${SALT}"
    [[ "$status" -eq 0 ]]
    local HASH="$output"
    echo "HASH: $HASH" >&2

    run sqlite3 "${BATS_TEST_DIRNAME}/sandbox/grafana.db" "SELECT password FROM user WHERE login='${USERNAME}';"
    echo "$output" >&2
    [[ "$status" -eq 0 ]]
    [[ "$output" = "$HASH" ]]
}

@test "check prometheus user 1" {
    if [ -n "${REMOTE}" ]; then
        skip "can be checked only locally"
    fi

    grep "^      username: " "${BATS_TEST_DIRNAME}/sandbox/prometheus.yml" >&2
    run grep "^      username: ${USERNAME}" "${BATS_TEST_DIRNAME}/sandbox/prometheus.yml"
    [[ "$status" -eq 0 ]]

    grep "^      password: " "${BATS_TEST_DIRNAME}/sandbox/prometheus.yml" >&2
    run grep "^      password: ${PASSWORD1}" "${BATS_TEST_DIRNAME}/sandbox/prometheus.yml"
    [[ "$status" -eq 0 ]]
}

@test "check http user 1" {
    if [ -n "${REMOTE}" ]; then
        skip "can be checked only locally"
    fi

    run htpasswd -csb "${BATS_TMPDIR}/htpasswd" "${USERNAME}" "${PASSWORD1}"
    echo "$output" >&2
    [[ "$status" -eq 0 ]]

    run grep "$(cat ${BATS_TMPDIR}/htpasswd)" "${BATS_TEST_DIRNAME}/sandbox/htpasswd"
    echo "$output" >&2
    [[ "$status" -eq 0 ]]

    rm -rf "${BATS_TMPDIR}/htpasswd"
}

@test "check config" {
    if [ -n "${REMOTE}" ]; then
        skip "can be checked only locally"
    fi

    grep "^  username: " "${BATS_TEST_DIRNAME}/sandbox/config.yml" >&2
    run grep "^  username: ${USERNAME}" "${BATS_TEST_DIRNAME}/sandbox/config.yml"
    [[ "$status" -eq 0 ]]

    grep "^- password: " "${BATS_TEST_DIRNAME}/sandbox/config.yml" >&2
    run grep "^- password: ${PASSWORD1}" "${BATS_TEST_DIRNAME}/sandbox/config.yml"
    [[ "$status" -eq 0 ]]
}

@test "update user" {
    run curl \
        -s \
        -X PATCH \
        --insecure \
        -d "${INPUT2}" \
        --user "${USERNAME}:${PASSWORD1}" \
        "${SUT}/${URL_PREFIX}/v1/users/${USERNAME}"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" = "$OUTPUT" ]]
}

@test "check updated grafana user" {
    if [ -n "${REMOTE}" ]; then
        skip "can be checked only locally"
    fi

    run sqlite3 "${BATS_TEST_DIRNAME}/sandbox/grafana.db" "SELECT salt FROM user WHERE login='${USERNAME}';"
    [[ "$status" -eq 0 ]]
    local SALT="$output"
    echo "SALT: $SALT" >&2

    run php "${BATS_TEST_DIRNAME}/password-helper.php" "${PASSWORD2}" "${SALT}"
    [[ "$status" -eq 0 ]]
    local HASH="$output"
    echo "HASH: $HASH" >&2

    run sqlite3 "${BATS_TEST_DIRNAME}/sandbox/grafana.db" "SELECT password FROM user WHERE login='${USERNAME}';"
    echo "$output" >&2
    [[ "$status" -eq 0 ]]
    [[ "$output" = "$HASH" ]]
}

@test "check updated prometheus user" {
    if [ -n "${REMOTE}" ]; then
        skip "can be checked only locally"
    fi

    grep "^      username: " "${BATS_TEST_DIRNAME}/sandbox/prometheus.yml" >&2
    run grep "^      username: ${USERNAME}" "${BATS_TEST_DIRNAME}/sandbox/prometheus.yml"
    [[ "$status" -eq 0 ]]

    grep "^      password: " "${BATS_TEST_DIRNAME}/sandbox/prometheus.yml" >&2
    run grep "^      password: ${PASSWORD2}" "${BATS_TEST_DIRNAME}/sandbox/prometheus.yml"
    [[ "$status" -eq 0 ]]
}


@test "check updated http user" {
    if [ -n "${REMOTE}" ]; then
        skip "can be checked only locally"
    fi

    run htpasswd -csb "${BATS_TMPDIR}/htpasswd" "${USERNAME}" "${PASSWORD2}"
    echo "$output" >&2
    [[ "$status" -eq 0 ]]

    run grep "$(cat ${BATS_TMPDIR}/htpasswd)" "${BATS_TEST_DIRNAME}/sandbox/htpasswd"
    echo "$output" >&2
    [[ "$status" -eq 0 ]]

    rm -rf "${BATS_TMPDIR}/htpasswd"
}

@test "check updated config" {
    if [ -n "${REMOTE}" ]; then
        skip "can be checked only locally"
    fi

    run grep "^  username: ${USERNAME}" "${BATS_TEST_DIRNAME}/sandbox/config.yml"
    [[ "$status" -eq 0 ]]

    run grep "^- password: ${PASSWORD2}" "${BATS_TEST_DIRNAME}/sandbox/config.yml"
    [[ "$status" -eq 0 ]]
}

@test "delete user" {
    run curl \
        -s \
        -X DELETE \
        --insecure \
        -d '' \
        --user "${USERNAME}:${PASSWORD2}" \
        "${SUT}/${URL_PREFIX}/v1/users/${USERNAME}"
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" = '{"code":200,"status":"OK"}' ]]
}

@test "check deleted grafana user" {
    if [ -n "${REMOTE}" ]; then
        skip "can be checked only locally"
    fi

    run sqlite3 "${BATS_TEST_DIRNAME}/sandbox/grafana.db" "SELECT login FROM user WHERE login='${USERNAME}';"
    echo "$output" >&2
    [[ "$status" -eq 0 ]]
    [[ -z "$output" ]]
}

@test "check deleted prometheus user" {
    if [ -n "${REMOTE}" ]; then
        skip "can be checked only locally"
    fi

    grep "^      username: " "${BATS_TEST_DIRNAME}/sandbox/prometheus.yml" >&2
    run grep "^      username: pmm" "${BATS_TEST_DIRNAME}/sandbox/prometheus.yml"
    [[ "$status" -eq 0 ]]

    grep "^      password: " "${BATS_TEST_DIRNAME}/sandbox/prometheus.yml" >&2
    run grep "^      password: pmm" "${BATS_TEST_DIRNAME}/sandbox/prometheus.yml"
    [[ "$status" -eq 0 ]]
}


@test "check deleted http user" {
    if [ -n "${REMOTE}" ]; then
        skip "can be checked only locally"
    fi

    if [ -f "${BATS_TEST_DIRNAME}/sandbox/htpasswd" ]; then
        run grep "${USERNAME}:" "${BATS_TEST_DIRNAME}/sandbox/htpasswd"
        echo "$output" >&2
        [[ "$status" -eq 1 ]]
    fi
}

@test "check deleted config" {
    if [ -n "${REMOTE}" ]; then
        skip "can be checked only locally"
    fi

    run grep "^  username: ${USERNAME}" "${BATS_TEST_DIRNAME}/sandbox/config.yml"
    [[ "$status" -ne 0 ]]

    run grep "^- password: ${PASSWORD2}" "${BATS_TEST_DIRNAME}/sandbox/config.yml"
    [[ "$status" -ne 0 ]]
}
