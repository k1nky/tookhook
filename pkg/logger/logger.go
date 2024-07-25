package logger

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
)

// TODO: impl. hclog.Logger interface

type Logger struct {
	hcLogger hclog.Logger
}

func New(name string) *Logger {
	l := hclog.New(&hclog.LoggerOptions{
		Name: name,
		// Exclude: ,
	})
	return &Logger{
		hcLogger: l,
	}
}

func (l *Logger) SetLevel(level string) error {
	levelValue := hclog.LevelFromString(level)
	if levelValue == hclog.NoLevel {
		return fmt.Errorf("unsupported log level %s", level)
	}
	l.hcLogger.SetLevel(levelValue)
	return nil
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.hcLogger.Debug(fmt.Sprintf(template, args...))
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.hcLogger.Info(fmt.Sprintf(template, args...))
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.hcLogger.Error(fmt.Sprintf(template, args...))
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.hcLogger.Warn(fmt.Sprintf(template, args...))
}

func (l *Logger) AsHCLogger() hclog.Logger {
	return l.hcLogger
}
