package pluginmanager

import "github.com/hashicorp/go-hclog"

type logger interface {
	Errorf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Debugf(template string, args ...interface{})
	AsHCLogger() hclog.Logger
}
