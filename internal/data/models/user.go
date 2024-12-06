package models

import "strings"

type User struct {
	Role
	Username string
}

func (u User) ID() string {
	return strings.ToLower(u.Username)
}
