package main

import (
	"github.com/Percona-Lab/pmm-manage/configurator/config"
	"github.com/Percona-Lab/pmm-manage/configurator/user"
	"log"
)

var c config.PMMConfig

func main() {
	c = config.ParseConfig()
	user.PMMConfig = c

	for _, userMap := range c.Users {
		log.Printf("CreateUser: %s\n", userMap["username"])
		result, err := user.CreateUser(user.PMMUser{Username: userMap["username"], Password: userMap["password"]})
		if result != "success" && err != nil {
			log.Printf("CreateUser: %s: %s\n", result, err)
		} else if result != "success" {
			log.Printf("CreateUser: %s\n", result)
		}
	}
}
