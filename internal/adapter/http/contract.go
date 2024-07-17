package http

import "context"

type logger interface {
	Errorf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Debugf(template string, args ...interface{})
}

//go:generate mockgen -source=contract.go -destination=mock/hooker.go -package=mock hookService
type hookService interface {
	Forward(ctx context.Context, name string, data []byte) error
	Reload(ctx context.Context) error
}
