package repositories

import (
	"strings"

	"github.com/guackamolly/zero-monitor/internal/data/db"
)

type DatabaseUserRepository struct {
	userTable db.UserTable
}

func NewDatabaseUserRepository(
	userTable db.UserTable,
) *DatabaseUserRepository {
	return &DatabaseUserRepository{
		userTable: userTable,
	}
}

func (r DatabaseUserRepository) AdminExists() (bool, error) {
	users, err := r.userTable.All()
	if err != nil && !strings.HasSuffix(err.Error(), "does not exist") {
		return false, err
	}

	if len(users) == 0 {
		return false, nil
	}

	for _, u := range users {
		if u.IsAdmin() {
			return true, nil
		}
	}

	return false, nil
}
