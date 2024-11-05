package dbbolt

import (
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/data/db"
	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"go.etcd.io/bbolt"
)

type BoltCrudTable[E db.Entity[string]] struct {
	name string
	db   *bbolt.DB
}

func NewBoltCrudTable[E db.Entity[string]](
	name string,
	db *bbolt.DB,
) *BoltCrudTable[E] {
	return &BoltCrudTable[E]{
		name: name,
		db:   db,
	}
}

func (t BoltCrudTable[E]) ID() string {
	return t.name
}

func (t BoltCrudTable[E]) All() ([]E, error) {
	vs := []E{}
	err := t.readonly(func(b *bbolt.Bucket) error {
		return b.ForEach(func(k, v []byte) error {
			e, err := models.Decode[E](v)
			if err != nil {
				logging.LogError("couldn't decode bucket value, %v", err)
				return nil
			}

			vs = append(vs, e)
			return nil
		})
	})

	return vs, err
}

func (t BoltCrudTable[E]) Delete(e E) error {
	return t.writable(func(b *bbolt.Bucket) error {
		pk := []byte(e.PK())
		return b.Delete(pk)
	})
}

func (t BoltCrudTable[E]) Insert(e E) error {
	return t.writable(func(b *bbolt.Bucket) error {
		k := []byte(e.PK())
		v, err := models.Encode(e)
		if err != nil {
			return err
		}

		return b.Put(k, v)
	})
}

func (t BoltCrudTable[E]) Lookup(pk string) (E, bool, error) {
	var d E
	var ok bool
	var err error

	err = t.readonly(func(b *bbolt.Bucket) error {
		v := b.Get([]byte(pk))
		if v == nil {
			ok = false
			return nil
		}
		d, err = models.Decode[E](v)
		ok = true

		return err
	})

	if !ok && err == nil {
		return d, ok, err
	}

	return d, err == nil, err
}

func (t BoltCrudTable[E]) Update(e E) error {
	return t.Insert(e)
}

func (t BoltCrudTable[E]) readonly(cb func(b *bbolt.Bucket) error) error {
	return t.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(t.name))
		if b == nil {
			return fmt.Errorf("bucket %s does not exist", t.name)
		}

		return cb(b)
	})
}

func (t BoltCrudTable[E]) writable(cb func(b *bbolt.Bucket) error) error {
	return t.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(t.name))
		if err != nil {
			return err
		}

		return cb(b)
	})
}
