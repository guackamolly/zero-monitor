package repositories

import (
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/data/db"
	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type DatabaseAuthenticationRepository struct {
	authTable db.CredentialsTable
	userTable db.UserTable
}

func NewDatabaseAuthenticationRepository(
	authTable db.CredentialsTable,
	userTable db.UserTable,
) *DatabaseAuthenticationRepository {
	return &DatabaseAuthenticationRepository{
		authTable: authTable,
		userTable: userTable,
	}
}

func (r DatabaseAuthenticationRepository) SignIn(username string, password string) (models.User, error) {
	credsEntity, ok, err := r.authTable.Lookup(username)
	if !ok || err != nil {
		return models.User{}, fmt.Errorf("user credentials not found")
	}

	if credsEntity.Password != password {
		return models.User{}, fmt.Errorf("user credentials don't match")
	}

	userEntity, ok, err := r.userTable.Lookup(username)
	if !ok || err != nil {
		return models.User{}, fmt.Errorf("user not found")
	}

	return userEntity.User, nil
}
