package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"strconv"

	"golang.org/x/crypto/ssh"
	"github.com/percona/pmm-manage/configurator/sshkey"
)

func parseSSHKey(authorizedKey []byte) (sshkey.SSHKey, error) {
	pubKey, comment, _, _, err := ssh.ParseAuthorizedKey(authorizedKey)
	if err != nil {
		return sshkey.SSHKey{}, err
	}
	return sshkey.SSHKey{
		Type:        pubKey.Type(),
		Comment:     comment,
		Fingerprint: ssh.FingerprintSHA256(pubKey),
	}, err
}

func getSSHKeyHandler(w http.ResponseWriter, req *http.Request) {
	authorizedKey, err := ioutil.ReadFile(c.SSHKeyPath)
	if err != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot read ssh key", err)
		return
	}
	sshKey, err := parseSSHKey(authorizedKey)
	if err != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot parse ssh key", err)
		return
	}
	json.NewEncoder(w).Encode(sshKey) // nolint: errcheck
}

func setSSHKeyHandler(w http.ResponseWriter, req *http.Request) {
	var newSSHKey sshkey.SSHKey
	if err := json.NewDecoder(req.Body).Decode(&newSSHKey); err != nil {
		returnError(w, req, http.StatusBadRequest, "Cannot parse json", err)
		return
	}

	parsedSSHKey, err := parseSSHKey([]byte(newSSHKey.Key))
	if err != nil {
		returnError(w, req, http.StatusBadRequest, "Cannot parse ssh key", err)
		return
	}

	if err = ioutil.WriteFile(c.SSHKeyPath, []byte(newSSHKey.Key), 0600); err != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot create authorized_keys file", err)
		return
	}

	sshKeyUser, err := user.Lookup(c.SSHKeyOwner)
	if err != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot lookup owner for authorized_keys file", err)
	}
	uid, _ := strconv.Atoi(sshKeyUser.Uid)
	gid, _ := strconv.Atoi(sshKeyUser.Gid)
	if err := os.Chown(c.SSHKeyPath, uid, gid); err != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot change owner for authorized_keys file", err)
	}

	location := fmt.Sprintf("http://%s%s", req.Host, req.URL.String())
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(parsedSSHKey) // nolint: errcheck
}
