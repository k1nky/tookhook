package database

import (
	"context"
	"strings"

	"github.com/k1nky/tookhook/internal/entity/rules"
)

// Database is adapter to database.
type Database interface {
	// Open connection to database.
	Open(ctx context.Context) (err error)
	// Close connection to database.
	Close() error
	GetRules(ctx context.Context) (*rules.Rules, error)
}

// New is factory of database connections.
func New(dsn string, log logger) Database {
	if strings.HasPrefix(dsn, "file://") {
		return NewFileStore(dsn, log)
	}
	return NewFileStore(dsn, log)
}
