package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"time"
)

func checkInstanceHandler(w http.ResponseWriter, req *http.Request) {
	var passedInstance instance
	if err := json.NewDecoder(req.Body).Decode(&passedInstance); err != nil {
		returnError(w, req, http.StatusBadRequest, "Cannot parse json", err)
		return
	}

	if result, err := checkInstance(passedInstance.ID); result != "success" {
		returnError(w, req, http.StatusForbidden, result, err)
		return
	}
	returnSuccess(w)
}

func checkInstance(instanceID string) (string, error) {
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	var rightInstanceID string
	instanceFile := path.Join(c.UpdateDirPath, "INSTANCE_ID")
	if _, err := os.Stat(instanceFile); err == nil {
		content, _ := ioutil.ReadFile(instanceFile)
		rightInstanceID = string(content)
	} else {
		resp, err := client.Get("http://169.254.169.254/latest/meta-data/instance-id")
		if _, isNetError := err.(net.Error); isNetError {
			return "success", nil
		}
		if err != nil {
			return "Cannot fetch instance meta-data", err
		}

		// ignore 404 error in non-AWS environments, like travis-ci
		if resp.StatusCode == 404 {
			return "success", nil
		}
		defer resp.Body.Close() // nolint: errcheck
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "Cannot parse instance meta-data", err
		}
		rightInstanceID = string(body)
	}

	if rightInstanceID == instanceID {
		return "success", nil
	}
	return "Wrong Instance ID", nil
}
