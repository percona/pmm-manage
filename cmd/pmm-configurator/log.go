package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-logfmt/logfmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func getLogFileHandler(w http.ResponseWriter, req *http.Request) { // nolint: gocyclo
	fileContent, err := ioutil.ReadFile(c.LogFilePath)
	if err != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot read pmm-manage log", err)
		return
	}

	// parse logfmt lines and show only fatal and error lines
	result := "ok"
	var output []string
	d := logfmt.NewDecoder(strings.NewReader(string(fileContent)))
	for d.ScanRecord() {
		line := ""
		addLine := false
		for d.ScanKeyval() {
			key := string(d.Key())
			value := string(d.Value())
			if key == "level" && value == "fatal" {
				addLine = true
				result = "fatal"
			} else if key == "level" && value == "error" {
				addLine = true
			} else if key == "msg" {
				line = value + " " + line
			} else if key != "time" && key != "action" {
				line = fmt.Sprintf(`%s %s="%s"`, line, key, value)
			}
		}
		if addLine {
			output = append(output, line)
		}
	}
	if d.Err() != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot parse pmm-manage log", d.Err())
		return
	}

	json.NewEncoder(w).Encode(jsonResponce{ // nolint: errcheck
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Title:  result,
		Detail: strings.Join(output, "\n"),
	})
}
