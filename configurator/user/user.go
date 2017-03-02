package user

import (
	"github.com/Percona-Lab/pmm-manage/configurator/config"
	"strings"
)

// PMMConfig pass configuration via global variable :'(
var PMMConfig config.PMMConfig

// CreateUser in .htpasswd file, Prometheus config and Grafana database
func CreateUser(newUser PMMUser) (string, error) {
	if strings.ContainsAny(newUser.Username, ":#") || len(newUser.Username) == 0 || len(newUser.Username) > 255 {
		return "Usernames are limited to 255 bytes and may not include colon and hash symbols", nil
	}

	if len(newUser.Password) == 0 || len(newUser.Password) > 255 {
		return "Passwords are limited to 255 bytes", nil
	}

	if err := createGrafanaUser(newUser); err != nil {
		return "Cannot set grafana password", err
	}

	if err := replacePrometheusUser(newUser); err != nil {
		return "Cannot set prometheus password", err
	}

	if err := createHTTPUser(newUser); err != nil {
		return "Cannot set http password", err
	}

	return "success", nil
}

// DeleteUser from Grafana and .htpasswd
// TODO: check user in Prometheus and replace to default if needed
func DeleteUser(username string) (string, error) {
	if err := deleteGrafanaUser(username); err != nil {
		return "Cannot remove grafana user", err
	}

	if err := deleteHTTPUser(username); err != nil {
		return "Cannot remove http user", err
	}

	return "success", nil
}
