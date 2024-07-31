package builtin

type logger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}
