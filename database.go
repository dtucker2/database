package database

import (
	"database/sql"
	"reflect"

	"github.com/pkg/errors"

	"github.com/dtucker2/database/query"
)

// Database wraps a sql.DB object and provides the ability to easily insert and update rows in a MySQL database using only a pointer to a struct.
type Database struct {
	*sql.DB
	*query.QueryBuilder
}

// NewDatabase returns a pointer to a new instance of the Database struct.
func NewDatabase(db *sql.DB) *Database {
	return &Database{
		DB:           db,
		QueryBuilder: query.NewQueryBuilder(),
	}
}

// Insert constructs and executes an insert query on the database using only the passed pointer to a struct.
func (db *Database) Insert(object interface{}) error {
	query, args := db.BuildInsertQuery(object)
	_, err := db.DB.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "Failed to execute query.")
	}
	return nil
}

// Update constructs and executes an update query on the database using only the passed pointer to a struct.
func (db *Database) Update(object interface{}) error {
	query, args := db.BuildUpdateQuery(object)
	_, err := db.DB.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "Failed to execute query.")
	}
	return nil
}

// Delete constructs and executes a delete query on the database using only the passed pointer to a struct.
// The structs primary key field must be populated as this populates the 'WHERE' clause of the query.
func (db *Database) Delete(object interface{}) error {
	query, args := db.BuildDeleteQuery(object)
	_, err := db.DB.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "Failed to execute query.")
	}
	return nil
}

// Select constructs and executes a select query on the database using only the passed pointer to a struct.
// The structs primary key field must be populated as this populates the 'WHERE' clause of the query.
// The resulting row will be returned by reference in the passed struct.
func (db *Database) Select(object interface{}) error {
	query, args := db.BuildSelectQuery(object)
	row := db.DB.QueryRow(query, args...)
	if err := row.Scan(db.getFieldPointers(object)...); err != nil {
		return errors.Wrap(err, "Failed to execute query.")
	}
	return nil
}

func (db *Database) getFieldPointers(object interface{}) []interface{} {
	val := reflect.ValueOf(object).Elem()
	ptrs := make([]interface{}, 0)
	for i := 0; i < val.NumField(); i++ {
		ptrs = append(ptrs, val.Field(i).Addr().Interface())
	}
	return ptrs
}
