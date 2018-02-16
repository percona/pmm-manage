package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

var pidRegexp = regexp.MustCompile(`PID: (\d+)`)
var resultRegexp = regexp.MustCompile(`localhost .* failed=0\s`)
var timeRegexp = regexp.MustCompile(`__(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}).log`)
var fromVersionRegexp = regexp.MustCompile(`> # v(\d+\.\d+\.\d+)\b`)
var toVersionRegexp = regexp.MustCompile(`< # v(\d+\.\d+\.\d+)\b`)
var currentVersionRegexp = regexp.MustCompile(`^# v(\d+\.\d+\.\d+)\b`)
var releaseNotesRegexp = regexp.MustCompile(`:Date: (.+)\n`)
var playbookRegexp = regexp.MustCompile(`1 plays in (.+)\n`)
var logTaskRegexp = regexp.MustCompile(`TASK`)
var playbookTaskRegexp = regexp.MustCompile(`- name:`)
var releaseNotesUrl = `https://raw.githubusercontent.com/percona/pmm/master/doc/source/release-notes/%s.rst`

func isPidAlive(pid int) bool {
	if err := syscall.Kill(pid, syscall.Signal(0x0)); err == nil {
		return true
	}
	return false
}

func isUpdateDisabled() bool {
	_, lockFileErr := os.Stat(path.Join(c.UpdateDirPath, "DISABLE_UPDATES"))
	disableUpdates, _ := strconv.ParseBool(os.Getenv("DISABLE_UPDATES"))

	// lock file exists or env variable is true
	if lockFileErr == nil || disableUpdates {
		return true
	}

	return false
}

func runCheckUpdateHandler(w http.ResponseWriter, req *http.Request) {
	pidFile := path.Join(c.UpdateDirPath, "pmm-update.pid")
	if _, err := os.Stat(pidFile); err == nil {
		timestamp, pid, err := getCurrentUpdate()
		if err != nil {
			returnError(w, req, http.StatusInternalServerError, "Cannot find update log", err)
			return
		}
		if isPidAlive(pid) {
			// update is going
			returnLog(w, req, timestamp, http.StatusOK)
			return
		}
	}

	// check if update is disabled
	if isUpdateDisabled() {
		returnError(w, req, http.StatusNotFound, "Updating of PMM is disabled.", nil)
		return
	}

	// check for update
	if cmdOutput, err := exec.Command("pmm-update-check").CombinedOutput(); err != nil {
		from, to := parseOutput(string(cmdOutput))
		json.NewEncoder(w).Encode(updateResponce{ // nolint: errcheck
			Code:   http.StatusOK,
			Status: http.StatusText(http.StatusOK),
			Title:  "A new PMM version is available.",
			From:   from,
			To:     to,
		})
		return
	}

	// no update
	returnError(w, req, http.StatusNotFound, "Your PMM version is up-to-date.", nil)
}

func parseOutput(output string) (string, string) {
	from := "unknown"
	to := "unknown"

	match := fromVersionRegexp.FindStringSubmatch(output)
	if len(match) == 2 {
		from = fetchReleaseDate(match[1])
	}

	match = toVersionRegexp.FindStringSubmatch(output)
	if len(match) == 2 {
		to = fetchReleaseDate(match[1])
	}

	return from, to
}

func fetchReleaseDate(version string) string {
	resp, err := http.Get(fmt.Sprintf(releaseNotesUrl, version))
	if err != nil {
		return version
	}
	defer resp.Body.Close() // nolint: errcheck

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return version
	}

	match := releaseNotesRegexp.FindStringSubmatch(string(body))
	if len(match) != 2 {
		return version
	}
	return version + " (" + match[1] + ")"
}

func readUpdateList() (map[string]string, error) {
	result := make(map[string]string)

	logPath := path.Join(c.UpdateDirPath, "log")
	files, err := ioutil.ReadDir(logPath)
	if err != nil {
		return result, err
	}

	for _, f := range files {
		if match := timeRegexp.FindStringSubmatch(f.Name()); len(match) == 2 {
			result[match[1]] = f.Name()
		}
	}

	return result, nil
}

func getUpdateListHandler(w http.ResponseWriter, req *http.Request) {
	updateList, err := readUpdateList()
	if err != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot read list of updates", err)
		return
	}

	keys := make([]string, 0, len(updateList))
	for k := range updateList {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	json.NewEncoder(w).Encode(keys) // nolint: errcheck
}

func getUpdateHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	returnLog(w, req, params["timestamp"], http.StatusOK)
}

func getCurrentVersionHandler(w http.ResponseWriter, req *http.Request) {
	fileContent, err := ioutil.ReadFile("/srv/update/main.yml")
	if err != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot read current version", err)
		return
	}
	match := currentVersionRegexp.FindSubmatch(fileContent)
	if len(match) == 2 {
		version := fetchReleaseDate(string(match[1]))
		json.NewEncoder(w).Encode(jsonResponce{ // nolint: errcheck
			Code:   http.StatusOK,
			Status: http.StatusText(http.StatusOK),
			Title:  version,
			Detail: version,
		})
	} else {
		returnError(w, req, http.StatusInternalServerError, "Cannot parse current version", err)
	}
}

func returnLog(w http.ResponseWriter, req *http.Request, timestamp string, httpStatus int) {
	updateList, err := readUpdateList()
	if err != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot read list of updates", err)
		return
	}

	logFile := updateList[timestamp]
	if logFile == "" {
		returnError(w, req, http.StatusNotFound, "Cannot find update", nil)
		return
	}

	filename := path.Join(c.UpdateDirPath, "log", logFile)
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot read update log", err)
		return
	}

	match := pidRegexp.FindStringSubmatch(string(fileContent))
	if len(match) != 2 {
		returnError(w, req, http.StatusInternalServerError, "Cannot find PID in update log", nil)
		return
	}

	pidInt, err := strconv.Atoi(match[1])
	if err != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot find PID in update log", nil)
		return
	}

	var updateState string
	if isPidAlive(pidInt) {
		updateState = "running"
	} else {
		if resultRegexp.MatchString(string(fileContent)) {
			updateState = "succeeded"
		} else {
			updateState = "failed"
		}
	}

	stepInfo := getStepsInfo(fileContent)

	location := fmt.Sprintf("%s/v1/updates/%s", c.PathPrefix, timestamp)
	w.Header().Set("Location", location)
	w.WriteHeader(httpStatus)

	json.NewEncoder(w).Encode(updateResponce{ // nolint: errcheck
		Code:   httpStatus,
		Status: http.StatusText(httpStatus),
		Title:  updateState,
		Detail: string(fileContent),
		Step:   stepInfo,
	})
}

func getStepsInfo(fileContent []byte) string {
	indexes := logTaskRegexp.FindAllIndex(fileContent, -1)
	currentStep := len(indexes)

	totalSteps := 0
	if playbookPathMatch := playbookRegexp.FindSubmatch(fileContent); len(playbookPathMatch) == 2 {
		if playbookContent, err := ioutil.ReadFile(string(playbookPathMatch[1])); err == nil {
			indexes := playbookTaskRegexp.FindAllIndex(playbookContent, -1)
			totalSteps = len(indexes) + 1 // add mandatory "Gathering Facts" task
		}
	}

	return fmt.Sprintf("%v/%v", currentStep, totalSteps)
}

func runUpdateHandler(w http.ResponseWriter, req *http.Request) {
	// check if update is disabled
	if isUpdateDisabled() {
		returnError(w, req, http.StatusNotFound, "Updating of PMM is disabled.", nil)
		return
	}

	if err := exec.Command("screen", "-d", "-m", "/usr/bin/pmm-update").Run(); err != nil { // nolint: gas
		returnError(w, req, http.StatusInternalServerError, "Cannot run update", err)
		return
	}

	// Advanced Sleep Programming :)
	time.Sleep(2 * time.Second)

	timestamp, _, err := getCurrentUpdate()
	if timestamp == "" || err != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot find update log", err)
		return
	}

	returnLog(w, req, timestamp, http.StatusAccepted)
}

func getCurrentUpdate() (string, int, error) {
	pidFile := path.Join(c.UpdateDirPath, "pmm-update.pid")
	pid, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return "", -1, err
	}

	pidStr := string(pid[:len(pid)-1])
	pidInt, err := strconv.Atoi(pidStr)
	if err != nil {
		return "", -1, err
	}

	pattern := fmt.Sprintf("PID: %s$", pidStr)
	logPath := path.Join(c.UpdateDirPath, "log/*.log")
	logs, err := filepath.Glob(logPath)
	if err != nil {
		return "", -1, err
	}

	args := append([]string{pattern}, logs...)
	currentLogOutput, err := exec.Command("grep", args...).Output() // nolint: gas
	if err != nil {
		return "", -1, err
	}

	match := timeRegexp.FindStringSubmatch(string(currentLogOutput))
	if len(match) != 2 {
		return "", -1, err
	}
	return match[1], pidInt, nil
}

func deleteUpdateHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	updateList, err := readUpdateList()
	if err != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot read list of updates", err)
		return
	}

	logFile := updateList[params["timestamp"]]
	if logFile == "" {
		returnError(w, req, http.StatusNotFound, "Cannot find update", nil)
		return
	}

	filename := path.Join(c.UpdateDirPath, "log", logFile)
	if err = os.Remove(filename); err != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot remove update log", nil)
		return
	}
	returnSuccess(w)
}
