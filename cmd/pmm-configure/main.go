package main

import (
	"github.com/Percona-Lab/pmm-manage/configurator/config"
	"github.com/Percona-Lab/pmm-manage/configurator/user"
	log "github.com/Sirupsen/logrus"
)

var c config.PMMConfig

func main() {
	c = config.ParseConfig()
	user.PMMConfig = c
	errorCounter := 0

	for _, userMap := range c.Users {
		rl := log.WithFields(log.Fields{"action": "CreateUser", "user": userMap["username"]})

		result, err := user.CreateUser(user.PMMUser{Username: userMap["username"], Password: userMap["password"]})
		if result == "success" {
			rl.Info("User was created successfully")
		} else if err != nil {
			errorCounter++
			rl.WithFields(log.Fields{"error": err}).Error(result)
		} else {
			errorCounter++
			rl.Error(result)
		}
	}

	if errorCounter == 0 {
		log.Info("PMM Server is configured correctly")
	} else {
		log.Fatal("PMM Server is not configured correctly")
	}
}
