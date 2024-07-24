package database

import (
	"context"
	"strings"

	"github.com/k1nky/tookhook/internal/entity"
)

// Database is adapter to database.
type Database interface {
	// Open connection to database.
	Open(ctx context.Context) (err error)
	// Close connection to database.
	Close() error
	GetRules(ctx context.Context) (*entity.Rules, error)
}

type logger interface {
	Errorf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Debugf(template string, args ...interface{})
}

// New is factory of database connections.
func New(dsn string, log logger) Database {
	if strings.HasPrefix(dsn, "file://") {
		return NewFileStore(dsn, log)
	}
	return NewFileStore(dsn, log)
}
