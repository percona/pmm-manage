package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getSSHKeyHandler(w http.ResponseWriter, req *http.Request) {
	parsedSSHKey, result, err := SSHKey.Read()
	if result != "success" {
		returnError(w, req, http.StatusInternalServerError, result, err)
	} else {
		json.NewEncoder(w).Encode(parsedSSHKey) // nolint: errcheck
	}
}

func setSSHKeyHandler(w http.ResponseWriter, req *http.Request) {
	parsedSSHKey, result, err := SSHKey.Write(req.Body)
	if result != "success" {
		returnError(w, req, http.StatusInternalServerError, result, err)
	} else {
		location := fmt.Sprintf("http://%s%s", req.Host, req.URL.String())
		w.Header().Set("Location", location)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(parsedSSHKey) // nolint: errcheck
	}
}
