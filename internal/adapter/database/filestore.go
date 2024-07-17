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
	DSN   string
	log   logger
	rules *entity.Rules
}

func NewFileStore(dsn string, log logger) *FileStore {
	return &FileStore{
		DSN:   strings.TrimPrefix(dsn, "file://"),
		log:   log,
		rules: &entity.Rules{},
	}
}

func (fs *FileStore) Open(ctx context.Context) (err error) {
	return fs.ReadRules(ctx)
}

func (fs *FileStore) ReadRules(ctx context.Context) error {
	f, err := os.Open(fs.DSN)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	rules := &entity.Rules{}
	if err := yaml.Unmarshal(data, rules); err != nil {
		return err
	}
	if err := rules.Validate(); err != nil {
		return err
	}
	fs.rules = rules
	return nil
}

func (fs *FileStore) GetIncomeHookByName(ctx context.Context, name string) (*entity.Hook, error) {
	for _, v := range fs.rules.Hooks {
		if v.Income == name {
			return &v, nil
		}
	}
	return nil, nil
}

func (fs *FileStore) Close() error {
	return nil
}
