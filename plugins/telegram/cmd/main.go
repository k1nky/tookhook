package main

import (
	"context"
	"os"

	hcplugin "github.com/hashicorp/go-plugin"
	"github.com/k1nky/tookhook/pkg/logger"
	"github.com/k1nky/tookhook/pkg/plugin"
	"github.com/k1nky/tookhook/plugins/telegram/internal/options"
	"github.com/k1nky/tookhook/plugins/telegram/internal/telegram"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DefaultLogLevel = "debug"
)

type Plugin struct {
	log *logger.Logger
}

func (p Plugin) Validate(ctx context.Context, r plugin.Receiver) error {
	opts, err := options.New(r.Options)
	if err != nil {
		p.log.Errorf("validate: %v", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if err := opts.Validate(); err != nil {
		p.log.Errorf("validate: %v", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}

func (p Plugin) Health(ctx context.Context) error {
	return nil
}

func (p Plugin) Forward(ctx context.Context, r plugin.Receiver, data []byte) ([]byte, error) {
	opts, err := options.New(r.Options)
	if err != nil {
		p.log.Errorf("forward: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	tg := telegram.New()
	m := telegram.Message{
		ChatID: opts.Chat,
		Text:   string(data),
	}
	response, err := tg.Send(opts.Token, m)
	p.log.Debugf("forward to %s with response: %s", opts.Chat, string(response))
	return response, err
}

func main() {
	log := newLogger()
	hcplugin.Serve(&hcplugin.ServeConfig{
		HandshakeConfig: plugin.Handshake,
		Plugins: map[string]hcplugin.Plugin{
			"grpc": &plugin.GRPCPlugin{Impl: &Plugin{
				log: log,
			}},
		},

		GRPCServer: hcplugin.DefaultGRPCServer,
	})
}

func newLogger() *logger.Logger {
	logLevel := os.Getenv("TOOKHOK_PLUGIN_TELEGRAM_LOG_LEVEL")
	if logLevel == "" {
		logLevel = DefaultLogLevel
	}
	l := logger.New("telegram")
	if err := l.SetLevel(DefaultLogLevel); err != nil {
		l.Errorf("invalid log level: %v", err)
	}
	return l
}
