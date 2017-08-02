#!/usr/bin/env bats

[ -z "$SUT" ] && SUT='http://127.0.0.1:7777' || :
[ -z "$URL_PREFIX" ] && URL_PREFIX='configurator' || :

@test "prepare" {
    mkdir -p "${BATS_TMPDIR}" || :
    rm -rf "${BATS_TMPDIR}"/id_rsa* || :
    ssh-keygen -t rsa -N '' -f "${BATS_TMPDIR}/id_rsa"
    ssh-keygen -y -f "${BATS_TMPDIR}/id_rsa" > "${BATS_TMPDIR}/id_rsa.pub"
}

@test "set sshkey" {
    KEY=$(cat ${BATS_TMPDIR}/id_rsa.pub)
    DIGEST=$(ssh-keygen -lf "${BATS_TMPDIR}/id_rsa.pub" | awk '{print$2}')

    run curl \
        -s \
        -X POST \
        --insecure \
        -d "{\"Key\": \"${KEY}\"}" \
        ${SUT}/${URL_PREFIX}/v1/sshkey
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" =~ '"type":"ssh-rsa"' ]]
    if [[ "$DIGEST" =~ "SHA256:" ]]; then
        [[ "$output" =~ "\"fingerprint\":\"$DIGEST\"" ]]
    fi

    if [ -z "$REMOTE" ]; then
        run diff -w "${BATS_TEST_DIRNAME}/sandbox/authorized_keys" "${BATS_TMPDIR}/id_rsa.pub"
        echo "$output" >&2
        [[ "$status" -eq 0 ]]
    fi
}

@test "get sshkey" {
    DIGEST=$(ssh-keygen -lf "${BATS_TMPDIR}/id_rsa.pub" | awk '{print$2}')

    run curl \
        -s \
        -X GET \
        --insecure \
        -d '' \
        ${SUT}/${URL_PREFIX}/v1/sshkey
    echo "$output" >&2

    [[ "$status" -eq 0 ]]
    [[ "$output" =~ '"type":"ssh-rsa"' ]]
    if [[ "$DIGEST" =~ "SHA256:" ]]; then
        [[ "$output" =~ "\"fingerprint\":\"$DIGEST\"" ]]
    fi
}

@test "cleanup" {
    rm -rf "${BATS_TEST_DIRNAME}/sandbox/authorized_keys" || :
    rm -rf "${BATS_TMPDIR}"/id_rsa* || :
}
