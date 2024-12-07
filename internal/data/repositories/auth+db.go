package repositories

import (
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/data/db"
	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/logging"
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

func (r DatabaseAuthenticationRepository) RegisterAdmin(username string, password string) (models.User, error) {
	if _, ok, _ := r.authTable.Lookup(username); ok {
		return models.User{}, fmt.Errorf("username already exists")
	}

	if _, ok, _ := r.userTable.Lookup(username); ok {
		return models.User{}, fmt.Errorf("username already exists")
	}

	user := models.NewAdminUser(username)
	err := r.userTable.Insert(db.NewUserEntity(user))
	if err != nil {
		return models.User{}, fmt.Errorf("user table insert failed, %v", err)
	}

	err = r.authTable.Insert(db.NewCredentialsEntity(username, password))
	if err == nil {
		return user, nil
	}

	delErr := r.userTable.Delete(db.NewUserEntity(user))
	if delErr != nil {
		logging.LogWarning("failed to insert user credentials, and now can't delete user from user table!")
	}

	return models.User{}, err
}
