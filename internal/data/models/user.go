package models

import "strings"

type User struct {
	Role
	Username string
}

func NewAdminUser(
	username string,
) User {
	return User{
		Username: username,
		Role:     AdminRole,
	}
}

func (u User) ID() string {
	return strings.ToLower(u.Username)
}

func (u User) IsAdmin() bool {
	return u.Role == AdminRole
}
