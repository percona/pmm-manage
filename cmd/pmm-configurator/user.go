package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/percona/pmm-manage/configurator/user"
)

func returnUser(w http.ResponseWriter, req *http.Request, username string) {
	users := user.ReadHTTPUsers()

	for _, item := range users {
		if item.Username == username {
			json.NewEncoder(w).Encode(item) // nolint: errcheck
			return
		}
	}
	returnError(w, req, http.StatusNotFound, "User is not found", nil)
}

func getUserListHandler(w http.ResponseWriter, req *http.Request) {
	users := user.ReadHTTPUsers()
	json.NewEncoder(w).Encode(users) // nolint: errcheck
}

func getUserHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	returnUser(w, req, params["username"])
}

func createUserHandler(w http.ResponseWriter, req *http.Request) {
	var newUser user.PMMUser
	if err := json.NewDecoder(req.Body).Decode(&newUser); err != nil {
		returnError(w, req, http.StatusBadRequest, "Cannot parse json", err)
		return
	}

	result, err := user.CreateUser(newUser)
	if err != nil {
		returnError(w, req, http.StatusInternalServerError, result, err)
		return
	}
	if result == "success" {
		location := fmt.Sprintf("%s/%s", req.URL.String(), newUser.Username)
		w.Header().Set("Location", location)
		w.WriteHeader(http.StatusCreated)
		returnUser(w, req, newUser.Username)
	} else {
		returnError(w, req, http.StatusForbidden, result, nil)
	}
}

func deleteUserHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if result, err := user.DeleteUser(params["username"]); result != "success" {
		returnError(w, req, http.StatusInternalServerError, result, err)
	} else {
		returnSuccess(w)
	}
}
