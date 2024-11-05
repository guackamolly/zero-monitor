//go:build integration
// +build integration

package dbbolt_test

import (
	"os"
	"testing"

	"github.com/guackamolly/zero-monitor/internal/data/db"
	dbbolt "github.com/guackamolly/zero-monitor/internal/data/db/db-bolt"
)

var testDir = os.TempDir()

func TestOpenErrorsIfTryingToOpenCorruptedDatabaseFile(t *testing.T) {
	fspath := createTempFile(t, 0600)
	err := os.WriteFile(fspath, []byte("corrupted database"), 0600)
	if err != nil {
		t.Fatalf("didn't expect writing to file to fail, %v", err)
	}

	db := dbbolt.NewBoltDatabase(fspath)

	err = db.Open()
	if err == nil {
		t.Errorf("expected open to fail since %s is a corrupted database file", fspath)
	}
}

func TestOpenErrorsIfItDoesNotHavePermissionsForReadingTheDatabaseFile(t *testing.T) {
	fspath := createTempFile(t, 0)
	db := dbbolt.NewBoltDatabase(fspath)

	err := db.Open()
	if err == nil {
		t.Errorf("expected open to fail since %s host does not have permissions to read database file", fspath)
	}
}

func TestOpenErrorsIfItDoesNotHavePermissionsForCreatingTheDatabaseFile(t *testing.T) {
	fspath := "/test.db"
	db := dbbolt.NewBoltDatabase(fspath)

	err := db.Open()
	if err == nil {
		t.Errorf("expected open to fail since %s host does not have permissions to read database file", fspath)
	}
}

func TestTable(t *testing.T) {
	bdb := createDb(t)

	testCases := []struct {
		desc   string
		input  string
		output bool
	}{
		{
			desc:   "returns not ok if table does not exist",
			input:  "this table does not exist I can assure you!",
			output: false,
		},
		{
			desc:   "returns ok if table exists",
			input:  db.TableSpeedtest,
			output: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if _, ok := bdb.Table(tC.input); ok != tC.output {
				t.Errorf("expected ok to be %v", ok)
			}
		})
	}
}

func createTempFile(t *testing.T, m os.FileMode) string {
	t.Helper()

	f, err := os.CreateTemp(testDir, "")
	if err != nil {
		t.Fatalf("didn't expect create temp file to fail, %v", err)
	}

	err = f.Chmod(m)
	if err != nil {
		t.Fatalf("didn't expect change modifiers to fail, %v", err)
	}

	t.Cleanup(func() {
		os.Remove(f.Name())
	})

	return f.Name()
}

func createDb(t *testing.T) *dbbolt.BoltDatabase {
	db := dbbolt.NewBoltDatabase(createTempFile(t, 0600))
	err := db.Open()
	if err != nil {
		t.Fatalf("didn't expect open to fail, %v", err)
	}

	t.Cleanup(func() { db.Close() })

	return db
}
