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

func (l *Logger) Debug(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	l.hcLogger.Debug(args[0].(string), args[1:]...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.hcLogger.Debug(fmt.Sprintf(template, args...))
}

func (l *Logger) Info(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	l.hcLogger.Info(args[0].(string), args[1:]...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.hcLogger.Info(fmt.Sprintf(template, args...))
}

func (l *Logger) Error(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	l.hcLogger.Error(args[0].(string), args[1:]...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.hcLogger.Error(fmt.Sprintf(template, args...))
}

func (l *Logger) Warn(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	l.hcLogger.Warn(args[0].(string), args[1:]...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.hcLogger.Warn(fmt.Sprintf(template, args...))
}

func (l *Logger) Fatal(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	l.hcLogger.Error(args[0].(string), args[1:]...)
}

func (l *Logger) AsHCLogger() hclog.Logger {
	return l.hcLogger
}

func (l *Logger) Sub(name string) *Logger {
	return &Logger{
		hcLogger: l.hcLogger.Named(name),
	}
}
