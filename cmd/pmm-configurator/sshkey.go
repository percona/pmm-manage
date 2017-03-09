package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
)

func runSSHKeyChecks() {
	sshKeyUser, err := user.Lookup(c.SSHKeyOwner)
	if err != nil {
		log.Fatal(err)
	}
	if c.SSHKeyPath == "" {
		c.SSHKeyPath = path.Join(sshKeyUser.HomeDir, ".ssh/authorized_keys")
	}

	sshKeyDir := filepath.Dir(c.SSHKeyPath)
	if dir, err := os.Stat(sshKeyDir); err != nil || !dir.IsDir() {
		if err := os.MkdirAll(sshKeyDir, 0700); err != nil {
			log.WithFields(log.Fields{
				"dir":   sshKeyDir,
				"error": err,
			}).Fatal("Cannot create ssh directory")
		}
		uid, _ := strconv.Atoi(sshKeyUser.Uid)
		gid, _ := strconv.Atoi(sshKeyUser.Gid)
		if err := os.Chown(sshKeyDir, uid, gid); err != nil {
			log.WithFields(log.Fields{
				"dir":   sshKeyDir,
				"error": err,
			}).Fatal("Cannot change owner of ssh directory")
		}
	}
	if err := unix.Access(sshKeyDir, unix.W_OK); err != nil {
		log.WithFields(log.Fields{
			"dir":   sshKeyDir,
			"error": err,
		}).Fatal("Cannot write to ssh directory")
	}
}

func parseSSHKey(authorizedKey []byte) (sshkey, error) {
	pubKey, comment, _, _, err := ssh.ParseAuthorizedKey(authorizedKey)
	if err != nil {
		return sshkey{}, err
	}
	return sshkey{
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
	var newSSHKey sshkey
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
