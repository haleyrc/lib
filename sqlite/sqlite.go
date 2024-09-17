// Package sqlite provides an opinionated wrapper around some libraries for
// working with SQLite databases.
package sqlite

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// DB provides a pool of connections to an SQLite database.
type DB struct {
	*sqlx.DB
}

// Open connects to the SQLite database at the given path. If the database does
// not already exist, it will be created the first time a method is called.
func Open(path string) (*DB, error) {
	sqlxDB, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("sqlite: open: %w", err)
	}

	db := &DB{
		DB: sqlxDB,
	}

	return db, nil
}
