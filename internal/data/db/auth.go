package db

type CredentialsTable CrudTable[CredentialsEntity, string]
type CredentialsEntity struct {
	Username string
	Password string
}

func NewCredentialsEntity(
	username string,
	password string,
) CredentialsEntity {
	return CredentialsEntity{
		Username: username,
		Password: password,
	}
}

func (e CredentialsEntity) PK() string {
	return e.Username
}
