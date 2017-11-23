package main

import (
	"encoding/json"
	"net/http"
)

func getSSHKeyHandler(w http.ResponseWriter, req *http.Request) {
	parsedSSHKey, result, err := SSHKey.Read()
	if result == "success" {
		json.NewEncoder(w).Encode(parsedSSHKey) // nolint: errcheck
	} else {
		returnError(w, req, http.StatusInternalServerError, result, err)
	}
}

func setSSHKeyHandler(w http.ResponseWriter, req *http.Request) {
	result, err := checkInstance("") // dirty hack, check if AWS EC2 instance
	if result != "success" {
		returnError(w, req, http.StatusForbidden, result, err)
		return
	}

	parsedSSHKey, result, err := SSHKey.Write(req.Body)
	if result == "success" {
		w.Header().Set("Location", req.URL.String())
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(parsedSSHKey) // nolint: errcheck
	} else {
		returnError(w, req, http.StatusInternalServerError, result, err)
	}
}
