package pluginmanager

import (
	"context"
	"os/exec"
	"sync"
	"time"

	hcplugin "github.com/hashicorp/go-plugin"
	"github.com/k1nky/tookhook/internal/entity"
	"github.com/k1nky/tookhook/pkg/plugin"
)

type pluginInstance struct {
	client     *hcplugin.Client
	command    string
	name       string
	grpcClient *plugin.GRPCClient
}

type Adapter struct {
	log     logger
	plugins sync.Map
	status  sync.Map
}

func New(log logger) *Adapter {
	return &Adapter{
		log: log,
	}
}

func (a *Adapter) Run(ctx context.Context) {
	a.runWatcher(ctx)
}

func (a *Adapter) Load(ctx context.Context, name string, command string) error {
	client := hcplugin.NewClient(&hcplugin.ClientConfig{
		HandshakeConfig: plugin.Handshake,
		Plugins:         plugin.PluginMap,
		Cmd:             exec.Command(command),
		Logger:          a.log.AsHCLogger(),
		AllowedProtocols: []hcplugin.Protocol{
			hcplugin.ProtocolGRPC},
	})

	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return err
	}
	raw, err := rpcClient.Dispense("grpc")
	if err != nil {
		client.Kill()
		return nil
	}
	// TODO: check type interface
	instance := pluginInstance{
		client:     client,
		command:    command,
		name:       name,
		grpcClient: raw.(*plugin.GRPCClient),
	}
	a.plugins.Store(name, instance)
	a.checkPluginHealth(ctx, name)

	return nil
}

func (a *Adapter) checkPluginHealth(ctx context.Context, name string) {
	p := a.Get(name)
	if p == nil {
		return
	}
	err := p.Health(ctx)
	if err != nil {
		a.log.Errorf("plugin manager watcher: %v", err)
		a.status.Store(name, entity.StatusFailed)
	} else {
		a.status.Store(name, entity.StatusOk)
	}
}

func (a *Adapter) runWatcher(ctx context.Context) {
	t := time.NewTicker(10 * time.Second)
	go func() {
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				a.plugins.Range(func(key, value any) bool {
					pluginName := key.(string)
					a.checkPluginHealth(ctx, pluginName)
					return true
				})
			}
		}
	}()
}

func (a *Adapter) Health(ctx context.Context) entity.PluginsStatus {
	status := make(entity.PluginsStatus)
	a.status.Range(func(key, value any) bool {
		status[key.(string)] = value.(entity.Status)
		return true
	})
	return status
}

func (a *Adapter) Get(name string) plugin.Plugin {
	c, ok := a.plugins.Load(name)
	if !ok {
		return nil
	}
	return c.(pluginInstance).grpcClient
}
