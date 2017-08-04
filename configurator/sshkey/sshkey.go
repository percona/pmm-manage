package sshkey

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sys/unix"

	"github.com/percona/pmm-manage/configurator/config"
)

func Init(c config.PMMConfig) Handler {
	return Handler{
		KeyPath:  c.SSHKeyPath,
		KeyOwner: c.SSHKeyOwner,
	}
}

func (c *Handler) RunSSHKeyChecks() {
	sshKeyUser, err := user.Lookup(c.KeyOwner)
	if err != nil {
		log.Fatal(err)
	}
	if c.KeyPath == "" {
		c.KeyPath = path.Join(sshKeyUser.HomeDir, ".ssh/authorized_keys")
	}

	sshKeyDir := filepath.Dir(c.KeyPath)
	logger := log.WithField("dir", sshKeyDir)
	if dir, err := os.Stat(sshKeyDir); err != nil || !dir.IsDir() {
		if err = os.MkdirAll(sshKeyDir+"/", 0700); err != nil {
			logger.WithField("error", err).Fatal("Cannot create ssh directory")
		}
		uid, _ := strconv.Atoi(sshKeyUser.Uid)
		gid, _ := strconv.Atoi(sshKeyUser.Gid)
		if err = os.Chown(sshKeyDir, uid, gid); err != nil {
			logger.WithField("error", err).Fatal("Cannot change owner of ssh directory")
		}
	}
	if err = unix.Access(sshKeyDir, unix.W_OK); err != nil {
		logger.WithField("error", err).Fatal("Cannot write to ssh directory")
	}
}

func parse(authorizedKey []byte) (*Key, error) {
	pubKey, comment, _, _, err := ssh.ParseAuthorizedKey(authorizedKey)
	if err != nil {
		return nil, err
	}
	return &Key{
		Type:        pubKey.Type(),
		Comment:     comment,
		Fingerprint: ssh.FingerprintSHA256(pubKey),
	}, err
}

// Read reads ssh key from c.KeyPath and parses it, returns sshKey, result and error.
// sshKey is struct with parsed key.
// result is human readable message, which can be shown to the end user (can be localized).
// error is debug information which is passed for debug needs, shouldn't be shown to user.
func (c *Handler) Read() (*Key, string, error) {
	authorizedKey, err := ioutil.ReadFile(c.KeyPath)
	if err != nil {
		return nil, "Cannot read ssh key", err
	}
	sshKey, err := parse(authorizedKey)
	if err != nil {
		return sshKey, "Cannot parse ssh key", err
	}
	return sshKey, "success", nil
}

// Write parses ssh key from json and writes to c.KeyPath and returns sshKey, result and error.
// sshKey is struct with parsed key.
// result is human readable message, which can be shown to the end user (can be localized).
// error is debug information which is passed for debug needs, shouldn't be shown to user.
func (c *Handler) Write(body io.ReadCloser) (*Key, string, error) {
	var newSSHKey Key
	if err := json.NewDecoder(body).Decode(&newSSHKey); err != nil {
		return &newSSHKey, "Cannot parse json", err
	}

	parsedSSHKey, err := parse([]byte(newSSHKey.Key))
	if err != nil {
		return parsedSSHKey, "Cannot parse ssh key", err
	}

	if err = ioutil.WriteFile(c.KeyPath, []byte(newSSHKey.Key), 0600); err != nil {
		return parsedSSHKey, "Cannot create authorized_keys file", err
	}

	sshKeyUser, err := user.Lookup(c.KeyOwner)
	if err != nil {
		return parsedSSHKey, "Cannot lookup owner for authorized_keys file", err
	}
	uid, _ := strconv.Atoi(sshKeyUser.Uid)
	gid, _ := strconv.Atoi(sshKeyUser.Gid)
	if err = os.Chown(c.KeyPath, uid, gid); err != nil {
		return parsedSSHKey, "Cannot change owner for authorized_keys file", err
	}

	return parsedSSHKey, "success", nil
}
