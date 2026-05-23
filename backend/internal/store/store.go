// Package store is the data-access layer: it owns the connection pool and
// every SQL query. Higher layers depend on this package, never the reverse.
package store

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a requested row does not exist.
var ErrNotFound = errors.New("not found")

// Store provides access to all persisted data.
type Store struct {
	pool *pgxpool.Pool
}

// New wraps a connection pool in a Store.
func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

// notFound translates pgx's no-rows error into the package sentinel.
func notFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
