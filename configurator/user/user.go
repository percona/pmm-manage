package user

import (
	"regexp"
	"strings"

	"github.com/fatih/structs"
	"github.com/percona/pmm-manage/configurator/config"
)

// PMMConfig pass configuration via global variable :'(
var PMMConfig config.PMMConfig

// CreateUser in .htpasswd file, Prometheus config and Grafana database
func CreateUser(newUser PMMUser) (string, error) { // nolint: gocyclo
	if strings.ContainsAny(newUser.Username, ":#") || len(newUser.Username) == 0 || len(newUser.Username) > 255 {
		return "User name is limited to 255 bytes and may not include colon and hash symbols", nil
	}

	if len(newUser.Password) == 0 || len(newUser.Password) > 255 {
		return "Password is limited to 255 bytes", nil
	}

	isAlphaNum := regexp.MustCompile(`^[A-Za-z0-9]`).MatchString
	if !isAlphaNum(newUser.Username) {
		return "User name should start with a letter or number", nil
	}

	if !isAlphaNum(newUser.Password) {
		return "Password should start with a letter or number", nil
	}

	if err := createGrafanaUser(newUser); err != nil {
		return "Cannot set Grafana password", err
	}

	if err := replacePrometheusUser(newUser); err != nil {
		return "Cannot set Prometheus password", err
	}

	if err := createHTTPUser(newUser); err != nil {
		return "Cannot set HTTP password", err
	}

	if err := PMMConfig.AddUser(structs.Map(newUser)); err != nil {
		return "Cannot save configuration file", err
	}

	return "success", nil
}

// DeleteUser from Grafana and .htpasswd
func DeleteUser(username string) (string, error) {
	if err := deleteGrafanaUser(username); err != nil {
		return "Cannot remove Grafana user", err
	}

	if err := resetPrometheusUser(); err != nil {
		return "Cannot reset Prometheus user", err
	}

	if err := deleteHTTPUser(username); err != nil {
		return "Cannot remove HTTP user", err
	}

	if err := PMMConfig.DeleteUser(username); err != nil {
		return "Cannot save configuration file", err
	}

	return "success", nil
}
