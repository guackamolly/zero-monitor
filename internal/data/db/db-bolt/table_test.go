//go:build integration
// +build integration

package dbbolt_test

import (
	"fmt"
	"testing"

	dbbolt "github.com/guackamolly/zero-monitor/internal/data/db/db-bolt"
	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type TestEntity struct {
	ID    string
	Value byte
}

func (e TestEntity) PK() string {
	return e.ID
}

func TestReadonlyOperationsReturnErrorIfBucketHasntBeenCreatedYet(t *testing.T) {
	bucket := "test.bucket"
	db := createDb(t)
	tbl := dbbolt.NewBoltCrudTable[TestEntity](bucket, db.DB)

	testCases := []struct {
		desc string
		op   func() error
	}{
		{
			desc: "all method should return error since bucket does not exist",
			op: func() error {
				_, err := tbl.All()
				return err
			},
		},
		{
			desc: "lookup method should return error since bucket does not exist",
			op: func() error {
				_, _, err := tbl.Lookup("pk")
				return err
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if err := tC.op(); err == nil {
				t.Error("expected operation to error since bucket does not exists")
			}
		})
	}
}

func TestWriteOperationsCreateBucketIfItDoesntExistYet(t *testing.T) {
	bucket := "test.bucket"
	db := createDb(t)
	tbl := dbbolt.NewBoltCrudTable[TestEntity](bucket, db.DB)

	entity := TestEntity{ID: models.UUID()}

	testCases := []struct {
		desc string
		op   func() error
	}{
		{
			desc: "insert method should not error if bucket does not exist",
			op: func() error {
				return tbl.Insert(entity)
			},
		},
		{
			desc: "update method should not error if bucket does not exist",
			op: func() error {
				return tbl.Update(entity)
			},
		},
		{
			desc: "delete method should not error if bucket does not exist",
			op: func() error {
				return tbl.Delete(entity)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if err := tC.op(); err != nil {
				t.Errorf("expected operation to not error, but got %v", err)
			}
		})
	}
}

func TestWriteOperationsCreateBucketError(t *testing.T) {
	db := createDb(t)

	testCases := []struct {
		desc   string
		bucket string
	}{
		{
			desc:   "cannot create bucket whose name is blank",
			bucket: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tbl := dbbolt.NewBoltCrudTable[TestEntity](tC.bucket, db.DB)
			err := tbl.Insert(TestEntity{})

			if err == nil {
				t.Error("expected insert to error on create bucket")
			}
		})
	}
}

func TestCRUD(t *testing.T) {
	bucket := "test.bucket"
	db := createDb(t)
	tbl := dbbolt.NewBoltCrudTable[TestEntity](bucket, db.DB)

	entity := TestEntity{ID: models.UUID()}
	updatedEntity := TestEntity{ID: entity.ID, Value: 155}

	testCases := []struct {
		desc string
		op   func() error
	}{
		{
			desc: "1. Insert",
			op: func() error {
				return tbl.Insert(entity)
			},
		},
		{
			desc: "2. Update",
			op: func() error {
				return tbl.Update(updatedEntity)
			},
		},
		{
			desc: "3. Lookup",
			op: func() error {
				entity, ok, err := tbl.Lookup(entity.ID)
				if err != nil {
					return err
				}

				fmt.Printf("entity: %v\n", entity)

				if !ok {
					return fmt.Errorf("no entity found")
				}

				if entity != updatedEntity {
					return fmt.Errorf("%v != %v", entity, updatedEntity)
				}

				return nil
			},
		},
		{
			desc: "4. Delete",
			op: func() error {
				return tbl.Delete(entity)
			},
		},
		{
			desc: "5. All",
			op: func() error {
				all, err := tbl.All()
				if err != nil {
					return err
				}

				if len(all) != 0 {
					return fmt.Errorf("expected All() to return no entities")
				}

				return nil
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if err := tC.op(); err != nil {
				t.Errorf("expected operation to not error, but got %v", err)
			}
		})
	}
}

func TestLookupReturnsFalseAndNilErrorIfEntityDoesNotExist(t *testing.T) {
	bucket := "test.bucket"
	db := createDb(t)
	tbl := dbbolt.NewBoltCrudTable[TestEntity](bucket, db.DB)

	entity := TestEntity{ID: models.UUID()}
	err := tbl.Insert(entity)
	if err != nil {
		t.Fatalf("didn't expect Insert() to fail, %v", err)
	}

	_, ok, _ := tbl.Lookup("noop")
	if ok {
		t.Error("expected ok to be false")
	}
}
