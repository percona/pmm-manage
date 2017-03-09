package main

import (
	"encoding/json"
	"github.com/Percona-Lab/pmm-manage/configurator/config"
	"github.com/Percona-Lab/pmm-manage/configurator/user"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

var c config.PMMConfig

func main() {
	c = config.ParseConfig()
	user.PMMConfig = c
	runSSHKeyChecks()

	router := mux.NewRouter().PathPrefix(c.PathPrefix).Subrouter()
	router.HandleFunc("/v1/sshkey", getSSHKeyHandler).Methods("GET")
	router.HandleFunc("/v1/sshkey", setSSHKeyHandler).Methods("POST")

	router.HandleFunc("/v1/check-update", runCheckUpdateHandler).Methods("GET")
	router.HandleFunc("/v1/updates", getUpdateListHandler).Methods("GET")
	router.HandleFunc("/v1/updates", runUpdateHandler).Methods("POST")
	router.HandleFunc("/v1/updates/{timestamp}", getUpdateHandler).Methods("GET")
	router.HandleFunc("/v1/updates/{timestamp}", deleteUpdateHandler).Methods("DELETE")

	router.HandleFunc("/v1/users", getUserListHandler).Methods("GET")
	router.HandleFunc("/v1/users", createUserHandler).Methods("POST")
	router.HandleFunc("/v1/users/{username}", getUserHandler).Methods("GET")
	router.HandleFunc("/v1/users/{username}", deleteUserHandler).Methods("DELETE")

	// TODO: create separate handler with old password verification
	router.HandleFunc("/v1/users/{username}", createUserHandler).Methods("PATCH")

	log.WithFields(log.Fields{
		"address": c.ListenAddress,
	}).Info("PMM Configurator is started")
	log.Fatal(http.ListenAndServe(c.ListenAddress, router))
}

func returnSuccess(w io.Writer) {
	json.NewEncoder(w).Encode(jsonResponce{ // nolint: errcheck
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
	})
}

func returnError(w http.ResponseWriter, req *http.Request, httpStatus int, title string, err error) {
	response := jsonResponce{
		Code:   httpStatus,
		Status: http.StatusText(httpStatus),
		Title:  title,
	}
	if err != nil {
		response.Detail = err.Error()
	}

	responseJSON, _ := json.Marshal(response)
	log.Errorf("%s %s: %s", req.Method, req.URL.String(), responseJSON)

	http.Error(w, string(responseJSON)+"\n", httpStatus)
}
