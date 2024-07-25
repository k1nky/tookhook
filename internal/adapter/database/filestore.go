package database

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/k1nky/tookhook/internal/entity"
	"gopkg.in/yaml.v2"
)

type FileStore struct {
	DSN string
	log logger
}

func NewFileStore(dsn string, log logger) *FileStore {
	return &FileStore{
		DSN: strings.TrimPrefix(dsn, "file://"),
		log: log,
	}
}

func (fs *FileStore) Open(ctx context.Context) (err error) {
	return nil
}

func (fs *FileStore) GetRules(ctx context.Context) (*entity.Rules, error) {
	f, err := os.Open(fs.DSN)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	rules := &entity.Rules{}
	if err := yaml.Unmarshal(data, rules); err != nil {
		return nil, err
	}
	return rules, nil
}

func (fs *FileStore) Close() error {
	return nil
}
