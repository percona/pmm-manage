package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var pathPrefix string
var sshKeyPath string
var sshKeyOwner string
var htpasswdPath string
var grafanaDBPath string
var listenAddress string
var prometheusConfPath string

func main() {
	parseFlag()

	router := mux.NewRouter().PathPrefix(pathPrefix).Subrouter()
	router.HandleFunc("/v1/sshkey", getSSHKeyHandler).Methods("GET")
	router.HandleFunc("/v1/sshkey", setSSHKeyHandler).Methods("POST")
	router.HandleFunc("/v1/users", getUserListHandler).Methods("GET")
	router.HandleFunc("/v1/users", createUserHandler).Methods("POST")
	router.HandleFunc("/v1/users/{username}", getUserHandler).Methods("GET")
	router.HandleFunc("/v1/users/{username}", deleteUserHandler).Methods("DELETE")

	// TODO: create separate handler with old password verification
	router.HandleFunc("/v1/users/{username}", createUserHandler).Methods("PATCH")

	log.Printf("PMM Configurator is started on %s address", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, router))
}

func parseFlag() {
	flag.StringVar(
		&htpasswdPath,
		"htpasswd-path",
		"/srv/nginx/.htpasswd",
		"htpasswd file location",
	)
	flag.StringVar(
		&listenAddress,
		"listen-address",
		"127.0.0.1:7777",
		"Address and port to listen on: [ip_address]:port",
	)
	flag.StringVar(
		&pathPrefix,
		"url-prefix",
		"/configurator",
		"Prefix for the internal routes of web endpoints",
	)
	flag.StringVar(
		&sshKeyPath,
		"ssh-key-path",
		"",
		"Path for SSH key",
	)
	flag.StringVar(
		&sshKeyOwner,
		"ssh-key-owner",
		"admin",
		"Owner of SSH key",
	)
	flag.StringVar(
		&grafanaDBPath,
		"grafana-db-path",
		"/srv/grafana/grafana.db",
		"grafana database location",
	)
	flag.StringVar(
		&prometheusConfPath,
		"prometheus-conf-path",
		"/etc/prometheus.yml",
		"prometheus configuration file location",
	)
	flag.Parse()

	runSSHKeyChecks()
}

func returnSuccess(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(jsonResponce{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
	})
}

func returnError(w http.ResponseWriter, req *http.Request, httpStatus int, title string, err error) {
	responce := jsonResponce{
		Code:   httpStatus,
		Status: http.StatusText(httpStatus),
		Title:  title,
	}
	if err != nil {
		responce.Detail = err.Error()
	}

	responceJSON, _ := json.Marshal(responce)
	log.Printf("%s %s: %s", req.Method, req.URL.String(), responceJSON)

	http.Error(w, string(responceJSON)+"\n", httpStatus)
}
