package db

import "github.com/guackamolly/zero-monitor/internal/data/models"

type UserTable CrudTable[UserEntity, string]
type UserEntity struct {
	models.User
}

func NewUserEntity(
	user models.User,
) UserEntity {
	return UserEntity{
		User: user,
	}
}

func (e UserEntity) PK() string {
	return e.ID()
}
