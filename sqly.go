// Package sqly is a simple helper library that help limiting the amount of
// repeating code necessary to execute query against a database/sql provider.
//
// It does not write the scans for you, but does not use reflection for known
// queries. Here's an example:
//
//    type Person struct {
//      FirstName string
//      LastName string
//    }
//
//    func scan(row sqly.Scanner, p *Person) error {
//      return row.Scan(&p.FirstName, &p.LastName)
//    }
//
//    func getAllPeople() ([]Person, error) {
//      query := "SELECT fname, lname FROM people;"
//      return sqly.Query(query, scan)
//    }
//
// The library also allow the run on-the-fly / dynamic query and receive
// a []map[string]any as rows, where the key of the map is the field names.
//
//    query := "SELECT fname || ' ' || lname as full_name FROM people"
//    rows, err := sqly.ExecuteDynamicQuery(query)
//    fmt.Println(rows[0]["full_name"])
package sqly

import (
	"context"
	"database/sql"
	"reflect"
)

var (
	// DB need to be initialize before calling non sql.Stmt functions
	DB *sql.DB
)

// Scanner enables sql.Rows and sql.Row to share the scans code.
type Scanner interface {
	Scan(dest ...any) error
}

// ScanCallback represents the type of function that will be call for each row
type ScanCallback[T any] func(row Scanner, entity *T) error

// Query executes the SQL query against the DB variable
func Query[T any](query string, scan ScanCallback[T], args ...any) ([]T, error) {
	return QueryContext(context.Background(), query, scan, args...)
}

// QueryContext executes the SQL query with a context against the DB variable
func QueryContext[T any](ctx context.Context, query string, scan ScanCallback[T], args ...any) (results []T, err error) {
	rows, err := DB.QueryContext(ctx, query, args...)
	if err != nil {
		return
	}
	return execQuery(rows, scan)
}

// QueryStatement executes the statement
func QueryStatement[T any](stmt *sql.Stmt, scan ScanCallback[T], args ...any) ([]T, error) {
	return QueryStatementContext(stmt, context.Background(), scan, args...)
}

// QueryStatementContext executes the statement with a context
func QueryStatementContext[T any](stmt *sql.Stmt, ctx context.Context, scan ScanCallback[T], args ...any) (results []T, err error) {
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return
	}
	return execQuery(rows, scan)
}

func execQuery[T any](rows *sql.Rows, scan ScanCallback[T]) (results []T, err error) {
	defer rows.Close()

	for rows.Next() {
		var entity T
		if err = scan(rows, &entity); err != nil {
			return
		}

		results = append(results, entity)
	}

	err = rows.Err()
	return
}

// QueryRow executes the SQL query to get one row against the DB variable
func QueryRow[T any](query string, scan ScanCallback[T], args ...any) (*T, error) {
	return QueryRowContext(context.Background(), query, scan, args...)
}

// QueryRowContext executes the SQL query to get one row with a context against
// the DB variable
func QueryRowContext[T any](ctx context.Context, query string, scan ScanCallback[T], args ...any) (*T, error) {
	row := DB.QueryRowContext(ctx, query, args...)
	return execQueryRow(row, scan)
}

// QueryStatement executes the statement to get one row
func QueryRowStatement[T any](stmt *sql.Stmt, scan ScanCallback[T], args ...any) (entity *T, err error) {
	return QueryRowStatementContext(stmt, context.Background(), scan, args...)
}

// QueryRowStatementContext executes the statement with a context to get one row
func QueryRowStatementContext[T any](stmt *sql.Stmt, ctx context.Context, scan ScanCallback[T], args ...any) (entity *T, err error) {
	row := stmt.QueryRowContext(ctx, args...)
	return execQueryRow(row, scan)
}

func execQueryRow[T any](row *sql.Row, scan ScanCallback[T]) (entity *T, err error) {
	entity = new(T)
	err = scan(row, entity)
	return
}

// ExecuteDynamicQuery executes the SQL query and returns a slice of ma[string]any
// Where the key of the map is the field name and the value is the interface{}
// value of the database field.
//
// Please note that this function uses reflection
func ExecuteDynamicQuery(query string, args ...any) (res []map[string]any, err error) {
	rows, err := DB.Query(query, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return
	}

	values := make([]any, len(columnNames))
	for i := range columnNames {
		var v any
		values[i] = &v
	}

	for rows.Next() {
		if err = rows.Scan(values...); err != nil {
			return
		}

		row := make(map[string]any)
		for i, col := range columnNames {
			//TODO: Is there not a better way just to get the value of
			// the pointer to the interface
			row[col] = reflect.Indirect(reflect.ValueOf(values[i])).Interface()
		}

		res = append(res, row)
	}

	err = rows.Err()
	return
}
