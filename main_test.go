package sqly_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/dstpierre/sqly"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

func TestMain(m *testing.M) {
	os.Remove("test.db")

	var err error

	db, err = sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqly.DB = db

	if err = createTables(db); err != nil {
		log.Fatal(err)
	}

	if err = insertRows(); err != nil {
		log.Fatal(err)
	}

	res := m.Run()
	os.Exit(res)
}

func createTables(db *sql.DB) error {
	// let's create a test table
	qry := `CREATE TABLE test (
					id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
					fname TEXT,
					lname TEXT,
					email TEXT
	);`

	if _, err := db.Exec(qry); err != nil {
		return err
	}
	return nil
}

func insertRows() error {
	// we used a prepared statement in here as we're running same query 10 times
	qry := `
					INSERT INTO test(fname, lname, email)
					VALUES(?, ?, ?);
	`

	ps, err := db.Prepare(qry)
	if err != nil {
		return err
	}
	defer ps.Close()

	for i := 0; i < 10; i++ {
		fname := fmt.Sprintf("fname_%d", i)
		lname := fmt.Sprintf("lname_%d", i)
		email := fmt.Sprintf("email_%d", i)

		if _, err := ps.Exec(fname, lname, email); err != nil {
			return err
		}
	}

	return nil
}
