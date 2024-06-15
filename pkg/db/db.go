// Package dbcontext provides DB transaction support for transactions tha span method calls of multiple
// repositories and services.
package dbcontext

import (
	"github.com/jackc/pgx/v5"
)

// DB represents a DB connection that can be used to run SQL queries.
type DB struct {
	db *pgx.Conn
}

// New returns a new DB connection that wraps the given dbx.DB instance.
func New(db *pgx.Conn) *DB {
	return &DB{db}
}

// DB returns the dbx.DB wrapped by this object.
func (db *DB) DB() *pgx.Conn {
	return db.db
}
