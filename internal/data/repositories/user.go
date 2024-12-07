package repositories

type UserRepository interface {
	AdminExists() (bool, error)
}
