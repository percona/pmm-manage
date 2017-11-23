package user

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var credentialsRegexp = regexp.MustCompile(`(username|password|job_name): ([^\s]+)`)
var replaceNeeded = map[string]bool{
	"linux":    true,
	"proxysql": true,
	"mongodb":  true,
	"mysql-hr": true,
	"mysql-lr": true,
	"mysql-mr": true,
}

// TODO: should be fully reworked, implemented as very quick workaround for v1.1.0
func replacePrometheusUser(newUser PMMUser) error {
	input, err := ioutil.ReadFile(PMMConfig.PrometheusConfPath)
	if err != nil {
		return err
	}

	var job string
	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.Contains(line, "job_name: ") {
			job = credentialsRegexp.FindStringSubmatch(line)[2]
		}
		if strings.Contains(line, "username: ") && replaceNeeded[job] {
			lines[i] = credentialsRegexp.ReplaceAllString(line, "$1: "+newUser.Username)
		}
		if strings.Contains(line, "password: ") && replaceNeeded[job] {
			lines[i] = credentialsRegexp.ReplaceAllString(line, "$1: "+newUser.Password)
		}
	}
	output := strings.Join(lines, "\n")

	if err = ioutil.WriteFile(PMMConfig.PrometheusConfPath, []byte(output), 0644); err != nil {
		return err
	}

	if PMMConfig.SkipPrometheusReload != "true" {
		req, err := http.NewRequest("POST", "http://127.0.0.1:9090/prometheus/-/reload", nil)
		if err != nil {
			return err
		}

		client := &http.Client{}
		if _, err := client.Do(req); err != nil {
			return err
		}
	}

	return nil
}

func resetPrometheusUser() error {
	return replacePrometheusUser(PMMUser{Username: "pmm", Password: "pmm"})
}
