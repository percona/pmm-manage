package main

import (
	"encoding/json"
	"fmt"
	"github.com/foomo/htpasswd"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

var users []user

func readUsers() {
	users = users[:0]
	if userMap, err := htpasswd.ParseHtpasswdFile(htpasswdFile); err == nil {
		for username := range userMap {
			users = append(users, user{Username: username, Password: "********"})
		}
	}
}

func returnUser(w http.ResponseWriter, req *http.Request, username string) {
	readUsers()

	for _, item := range users {
		if item.Username == username {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	returnError(w, req, http.StatusNotFound, "User is not found", nil)
}

func getUserListHandler(w http.ResponseWriter, req *http.Request) {
	readUsers()
	json.NewEncoder(w).Encode(users)
}

func getUserHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	returnUser(w, req, params["username"])
}

func createUserHandler(w http.ResponseWriter, req *http.Request) {
	var newUser user
	if err := json.NewDecoder(req.Body).Decode(&newUser); err != nil {
		returnError(w, req, http.StatusBadRequest, "Cannot parse json", err)
		return
	}

	if strings.ContainsAny(newUser.Username, ":") || len(newUser.Username) >= 255 {
		returnError(w, req, http.StatusForbidden, "Usernames are limited to 255 bytes and may not include the character :.", nil)
		return
	}

	if err := htpasswd.SetPassword(htpasswdFile, newUser.Username, newUser.Password, htpasswd.HashBCrypt); err != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot set password", err)
		return
	}

	location := fmt.Sprintf("http://%s%s/%s", req.Host, req.URL.String(), newUser.Username)
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
	returnUser(w, req, newUser.Username)
}

func deletePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if err := htpasswd.RemoveUser(htpasswdFile, params["username"]); err != nil {
		returnError(w, req, http.StatusInternalServerError, "Cannot remove the user", err)
		return
	}
	returnSuccess(w)
}
