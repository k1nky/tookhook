package database

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
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

func (a *FileStore) Open(ctx context.Context) (err error) {
	if err := a.runWatcher(ctx); err != nil {
		return err
	}
	return a.ReadRules(ctx)
}

func (a *FileStore) ReadRules(ctx context.Context) error {
	f, err := os.Open(a.DSN)
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
	a.rules = rules
	return nil

}

func (a *FileStore) GetIncomeHookByName(ctx context.Context, name string) (*entity.Hook, error) {
	var hook *entity.Hook = new(entity.Hook)

	for _, v := range a.rules.Hooks {
		if v.Income == name {
			// TODO: check it
			*hook = v
			return hook, nil
		}
	}
	return hook, nil
}

func (a *FileStore) Close() error {
	return nil
}

func (a *FileStore) runWatcher(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if err := watcher.Add(a.DSN); err != nil {
		return err
	}
	go func() {
		defer watcher.Close()
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					if err := a.ReadRules(ctx); err != nil {
						a.log.Errorf("read rules: %v", err)
						continue
					}
					a.log.Infof("read rules: success")
				}
			case err := <-watcher.Errors:
				a.log.Errorf("rules watcher: %v", err)
			}
		}
	}()
	time.Sleep(time.Second)
	return nil
}
