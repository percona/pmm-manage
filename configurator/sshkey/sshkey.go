package sshkey

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"

	"github.com/percona/pmm-manage/configurator/config"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/sys/unix"
)

// PMMConfig pass configuration via global variable :'(
var PMMConfig config.PMMConfig

func RunSSHKeyChecks() {
	sshKeyUser, err := user.Lookup(PMMConfig.SSHKeyOwner)
	if err != nil {
		log.Fatal(err)
	}
	if PMMConfig.SSHKeyPath == "" {
		PMMConfig.SSHKeyPath = path.Join(sshKeyUser.HomeDir, ".ssh/authorized_keys")
	}

	sshKeyDir := filepath.Dir(PMMConfig.SSHKeyPath)
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

func ParseSSHKey(authorizedKey []byte) (SSHKey, error) {
	pubKey, comment, _, _, err := ssh.ParseAuthorizedKey(authorizedKey)
	if err != nil {
		return SSHKey{}, err
	}
	return SSHKey{
		Type:        pubKey.Type(),
		Comment:     comment,
		Fingerprint: ssh.FingerprintSHA256(pubKey),
	}, err
}
