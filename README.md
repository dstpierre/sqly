## sqly 

A simple `database/sql` helper for making queries less repetitive.

It does not use reflection, basically just using generics to make `Query` and 
`QueryRow` feel a bit less repetitive.

You may run ad-hoc / dynamic queries. For those reflection is use to 
return a slice of `map[string]any` for each row.

### Why, there's multiple libraries already

I published this 
[podcast episode recently](https://go.transistor.fm/episodes/007-is-gos-database-sql-verbosity-that-bad) 
talking about database in Go from using the standard library to an ORM like Gorm.

Since that time I was reflecting on what I truly find tedious with the stdlib.

### Install

```sh
$> go get github.com/dstpierre/sqly
```
### Usage

You connect to your database as usual and pass the `*sql.DB` reference to sqly.

```go
package main

import (
	//...

	"github.com/dstpierre/sqly"
)

func main() {
	db, err := sql.Open(...)

	sqly.DB = db

	// from there you may run all functions
}
```

You continue to define your types and your scans function as usual, for instance:

```go
type Person struct {
	ID int64
	FirstName string
	LastName string
	Email string
}
```

Scan functions are where you turn a row into a struct (or any type you want):

```go
func scanPerson(row sqly.Scanner, p *Person) error {
	return row.Scan(
		&p.ID,
		&p.FirstName,
		&p.LastName,
		&p.Email,
	)
}
```

And you can now call `Query` or `QueryRow` like this

```go
func getAllPeople() ([]Person, error) {
	return sqly.Query("SELECT * FROM people;", scanPerson)
}
```

The `Query` and `QueryRow` add a simple scan function callback compare 
to the standard library `database/sql` functions.

If you don't want to have a scan function separately, you may inline it:

```go
func getAllPeople() (results []Person, err error) {
	sqly.Query("SELECT * FROM people WHERE id > ?;", func(row sqly.Scanner, p *Person) {
		return row.Scan(...)
	}, 5)
	return
}
```

When you're expecting one result you may call the `QueryRow` function:

```go
func getPersonByEmail(email string) (*Person, error) {
	query := "SELECT * FROM people WHERE email = ?"
	return sqly.QueryRow(query, scanPerson, email)
}
```

You may create the scan function inline for query result that are not re-used:

```go
func getFullNames() (names []string, err error) {
	query := "SELECT fname + ' '  lname FROM people;"
	scan := func(row sqly.Scanner, s *string) error {
		return row.Scan(s)
	}

	return sqly.Query(query, scan)
}
```

### Dynamic / ad-hoc queries

Sometimes you might want to quickly query the database for dynamic queries.

Here's how you can perform the same as above without handling the scan:

```go
func execDynamicQuery() (names []string, err error) {
	query := "SELECT fname + ' '  lname as full_name FROM people;"
	rows, err := sqly.ExecuteDynamicQuery(query)
	if err != nil {
		return
	}

	for _, row := range rows {
		names = append(names, row["full_name])
	}
	return
}
```

The `ExecuteDynamicQuery` returns a `[]map[string]any` where the key of the map 
is the field name and the interface{} its value.

