package database

import (
	"context"
	"strings"

	"github.com/k1nky/tookhook/internal/entity"
)

type Database interface {
	Open(ctx context.Context) (err error)
	Close() error
	ReadRules(ctx context.Context) error
	GetIncomeHookByName(ctx context.Context, name string) (*entity.Hook, error)
}

type logger interface {
	Errorf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Debugf(template string, args ...interface{})
}

func New(dsn string, log logger) Database {
	if strings.HasPrefix(dsn, "file://") {
		return NewFileStore(dsn, log)
	}
	return nil
}
