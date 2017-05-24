package config

import (
	"strings"
)

func (c *PMMConfig) AddUser(newUser map[string]interface{}) error {
	user := make(map[string]string)
	for key, value := range newUser {
		key = strings.ToLower(key)
		user[key] = value.(string)
	}

	if err := c.DeleteUser(user["username"]); err != nil {
		return err
	}
	c.Users = append(c.Users, user)

	return c.Save()
}

func (c *PMMConfig) DeleteUser(username string) error {
	var fixedUsers []map[string]string

	for _, userMap := range c.Users {
		if userMap["username"] != username {
			fixedUsers = append(fixedUsers, userMap)
		}
	}

	c.Users = fixedUsers
	return c.Save()
}
