package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

var htpasswdFile string
var listenAddress string
var pathPrefix string

func main() {
	parseFlag()

	router := mux.NewRouter().PathPrefix(pathPrefix).Subrouter()
	router.HandleFunc("/v1/users", getUserListHandler).Methods("GET")
	router.HandleFunc("/v1/users", createUserHandler).Methods("POST")
	router.HandleFunc("/v1/users/{username}", createUserHandler).Methods("PATCH")
	router.HandleFunc("/v1/users/{username}", getUserHandler).Methods("GET")
	router.HandleFunc("/v1/users/{username}", deletePersonEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(listenAddress, router))
}

func parseFlag() {
	flag.StringVar(
		&htpasswdFile,
		"htpasswd",
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
	flag.Parse()
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

	w.WriteHeader(httpStatus)
	io.WriteString(w, string(responceJSON)+"\n")
}
