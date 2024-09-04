package logger

type Blackhole struct{}

func (l *Blackhole) Debug(args ...interface{})                   {}
func (l *Blackhole) Debugf(template string, args ...interface{}) {}
func (l *Blackhole) Info(args ...interface{})                    {}
func (l *Blackhole) Infof(template string, args ...interface{})  {}
func (l *Blackhole) Warn(args ...interface{})                    {}
func (l *Blackhole) Warnf(template string, args ...interface{})  {}
func (l *Blackhole) Error(args ...interface{})                   {}
func (l *Blackhole) Errorf(template string, args ...interface{}) {}
func (l *Blackhole) Fatal(args ...interface{})                   {}
