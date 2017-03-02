package user

import (
	"github.com/foomo/htpasswd"
)

// ReadHTTPUsers func read list of users from .htpasswd file
func ReadHTTPUsers() []PMMUser {
	var users []PMMUser
	if userMap, err := htpasswd.ParseHtpasswdFile(PMMConfig.HtpasswdPath); err == nil {
		for username := range userMap {
			users = append(users, PMMUser{Username: username, Password: "********"})
		}
	}
	return users
}

func createHTTPUser(newUser PMMUser) error {
	// htpasswd.HashBCrypt is better, but nginx server in CentOS 7, doesn't support it :(
	return htpasswd.SetPassword(PMMConfig.HtpasswdPath, newUser.Username, newUser.Password, htpasswd.HashSHA)
}

func deleteHTTPUser(username string) error {
	return htpasswd.RemoveUser(PMMConfig.HtpasswdPath, username)
}
