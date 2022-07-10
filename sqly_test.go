package sqly_test

import (
	"fmt"
	"testing"

	"github.com/dstpierre/sqly"
)

type Person struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
}

func scan(row sqly.Scanner, p *Person) error {
	return row.Scan(
		&p.ID,
		&p.FirstName,
		&p.LastName,
		&p.Email,
	)
}

func TestQuery(t *testing.T) {
	qry := "SELECT * FROM test;"

	results, err := sqly.Query(qry, scan)
	if err != nil {
		t.Fatal(err)
	} else if len(results) < 10 {
		t.Errorf("expected at least 10 rows got %d", len(results))
	} else if results[1].FirstName != "fname_1" {
		t.Errorf(`expected fname to be "fname_1" got %s`, results[0].FirstName)
	}
}

func TestQueryRow(t *testing.T) {
	qry := `SELECT * FROM test WHERE id = ?;`
	person, err := sqly.QueryRow[Person](qry, scan, 5)
	if err != nil {
		t.Fatal(err)
	} else if person.FirstName != "fname_4" {
		t.Errorf(`expected fname to be "fname_4" got %s`, person.FirstName)
	}
}

func TestDynamicQuery(t *testing.T) {
	qry := "SELECT email FROM test WHERE fname = ?;"
	rows, err := sqly.ExecuteDynamicQuery(qry, "fname_3")
	if err != nil {
		t.Fatal(err)
	} else if len(rows) != 1 {
		t.Errorf("expected 1 row got %d", len(rows))
	} else if rows[0]["email"] != "email_3" {
		t.Errorf(`expected email to be "email_3" got %v`, rows[0]["email"])
	}
}

func TestQueryStatement(t *testing.T) {
	stmt, err := db.Prepare("SELECT * FROM test WHERE id = ?")
	if err != nil {
		t.Fatal(err)
	}

	for i := 1; i < 6; i++ {
		p, err := sqly.QueryRowStatement(stmt, scan, i)
		if err != nil {
			t.Fatal(err)
		} else if n := fmt.Sprintf("fname_%d", i-1); p.FirstName != n {
			t.Errorf(`expected fname to be %s got %s`, n, p.FirstName)
		}
	}
}

func TestInlineType(t *testing.T) {
	type overview struct {
		FirstName string
		LastName  string
	}

	scanOverview := func(row sqly.Scanner, o *overview) error {
		return row.Scan(&o.FirstName, &o.LastName)
	}

	qry := "SELECT fname, lname FROM test;"
	results, err := sqly.Query(qry, scanOverview)
	if err != nil {
		t.Fatal(err)
	} else if len(results) < 5 {
		t.Errorf("expected at least 5 results, got %d", len(results))
	} else if n := results[3].FirstName; n != "fname_3" {
		t.Errorf(`expected fname to be "fname_3" got "%s"`, n)
	}
}
