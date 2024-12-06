package models

const (
	AdminRole Role = iota + 1
	GuestRole
)

type Role int
