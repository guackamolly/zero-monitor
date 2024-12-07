package dbbolt

import (
	"github.com/guackamolly/zero-monitor/internal/data/db"
	"go.etcd.io/bbolt"
)

type BoltDatabase struct {
	DB     *bbolt.DB
	fspath string
	tables []db.Table
}

func NewBoltDatabase(
	fspath string,
) *BoltDatabase {
	return &BoltDatabase{
		fspath: fspath,
		tables: []db.Table{},
	}
}

func (bdb *BoltDatabase) Open() error {
	bbdb, err := bbolt.Open(bdb.fspath, 0600, nil)
	if err != nil {
		return err
	}

	stt := NewBoltCrudTable[db.SpeedtestEntity](db.TableSpeedtest, bbdb)
	crt := NewBoltCrudTable[db.CredentialsEntity](db.TableCredentials, bbdb)
	ust := NewBoltCrudTable[db.UserEntity](db.TableUser, bbdb)

	bdb.DB = bbdb
	bdb.tables = []db.Table{stt, crt, ust}

	return nil
}

func (bdb *BoltDatabase) Tables() []db.Table {
	return bdb.tables
}

func (bdb *BoltDatabase) Table(id string) (db.Table, bool) {
	for _, t := range bdb.tables {
		if t.ID() == id {
			return t, true
		}
	}

	return nil, false
}

func (bdb *BoltDatabase) Close() error {
	return bdb.DB.Close()
}
