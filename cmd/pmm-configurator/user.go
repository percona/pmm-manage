package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func returnUser(w http.ResponseWriter, req *http.Request, username string) {
	users := readHTTPUsers()

	for _, item := range users {
		if item.Username == username {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	returnError(w, req, http.StatusNotFound, "User is not found", nil)
}

func getUserListHandler(w http.ResponseWriter, req *http.Request) {
	users := readHTTPUsers()
	json.NewEncoder(w).Encode(users)
}

func getUserHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	returnUser(w, req, params["username"])
}

func createUserHandler(w http.ResponseWriter, req *http.Request) {
	var newUser htuser
	if err := json.NewDecoder(req.Body).Decode(&newUser); err != nil {
		returnError(w, req, http.StatusBadRequest, "Cannot parse json", err)
		return
	}

	if strings.ContainsAny(newUser.Username, ":") || len(newUser.Username) == 0 || len(newUser.Username) > 255 {
		returnError(w, req, http.StatusForbidden, "Usernames are limited to 255 bytes and may not include the colon symbol", nil)
		return
	}

	if len(newUser.Password) == 0 || len(newUser.Password) > 255 {
		returnError(w, req, http.StatusForbidden, "Passwords are limited to 255 bytes", nil)
		return
	}

	if err := createHTTPUser(newUser); err != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot set http password", err)
		return
	}

	location := fmt.Sprintf("http://%s%s/%s", req.Host, req.URL.String(), newUser.Username)
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
	returnUser(w, req, newUser.Username)
}

func deleteUserHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if err := deleteHTTPUser(params["username"]); err != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot remove http user", err)
		return
	}
	returnSuccess(w)
}
