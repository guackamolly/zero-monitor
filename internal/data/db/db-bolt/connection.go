package dbbolt

import "os"

const (
	boltDbPathKey = "bolt_db_path"
)

var (
	boltDbPath = os.Getenv(boltDbPathKey)
)

func Path() string {
	if len(boltDbPath) == 0 {
		return "master.db"
	}

	return boltDbPath
}
