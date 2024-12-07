package db

// All database tables managed by the master node.
const (
	TableSpeedtest   = "node.speedtests"
	TableCredentials = "auth.credentials"
	TableUser        = "auth.user"
)

type Table interface {
	ID() string
}

// A [Table] with all Create/Read/Update/Delete operations, for a specific [Entity].
type CrudTable[E Entity[PK], PK ~string] interface {
	Table
	Insert(E) error
	Update(E) error
	Lookup(PK) (E, bool, error)
	All() ([]E, error)
	Delete(E) error
}

// Abstract an "object" that can be persisted .
type Entity[PK ~string] interface {
	PK() PK
}
