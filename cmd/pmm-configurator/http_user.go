package main

import (
	"github.com/foomo/htpasswd"
)

func readHTTPUsers() []htuser {
	var users []htuser
	if userMap, err := htpasswd.ParseHtpasswdFile(c.HtpasswdPath); err == nil {
		for username := range userMap {
			users = append(users, htuser{Username: username, Password: "********"})
		}
	}
	return users
}

func createHTTPUser(newUser htuser) error {
	// htpasswd.HashBCrypt is better, but nginx server in CentOS 7, doesn't support it :(
	return htpasswd.SetPassword(c.HtpasswdPath, newUser.Username, newUser.Password, htpasswd.HashSHA)
}

func deleteHTTPUser(username string) error {
	return htpasswd.RemoveUser(c.HtpasswdPath, username)
}
