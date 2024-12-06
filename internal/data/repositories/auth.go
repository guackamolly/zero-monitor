package repositories

import "github.com/guackamolly/zero-monitor/internal/data/models"

type AuthenticationRepository interface {
	SignIn(username string, password string) (models.User, error)
	RegisterAdmin(username string, password string) (models.User, error)
}
