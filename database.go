package database

import (
	"database/sql"

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
		return errors.Wrap(err, "Failed to insert object.")
	}
	return nil
}

// Update constructs and executes an update query on the database using only the passed pointer to a struct.
func (db *Database) Update(object interface{}) error {
	query, args := db.BuildUpdateQuery(object)
	_, err := db.DB.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "Failed to insert object.")
	}
	return nil
}
