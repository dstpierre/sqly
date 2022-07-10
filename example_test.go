package sqly_test

import (
	"fmt"

	"github.com/dstpierre/sqly"
)

func ExampleQuery() {
	query := "SELECT * FROM test WHERE id > ?;"
	rows, err := sqly.Query(query, scan, 5)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rows[0].FirstName)
	// Output: fname_5
}

func ExampleQueryRow() {
	query := "SELECT * FROM test WHERE email = ?"
	p, err := sqly.QueryRow(query, scan, "email_3")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(p.FirstName)
	// Output: fname_3
}

func ExampleExecuteDynamicQuery() {
	query := "SELECT fname || ' ' || lname as full_name FROM test WHERE id < ?;"
	results, err := sqly.ExecuteDynamicQuery(query, 3)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(results[0]["full_name"], results[1]["full_name"])
	// Output: fname_0 lname_0 fname_1 lname_1
}
