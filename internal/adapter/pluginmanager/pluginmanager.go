package pluginmanager

import (
	"context"
	"os/exec"
	"sync"
	"time"

	hcplugin "github.com/hashicorp/go-plugin"
	"github.com/k1nky/tookhook/pkg/plugin"
)

type pluginsList map[string]*plugin.GRPCClient
type statusList map[string]bool

type Adapter struct {
	log        logger
	plugins    pluginsList
	lockStatus sync.RWMutex
	status     statusList
}

func New(log logger) *Adapter {
	return &Adapter{
		log:     log,
		plugins: make(pluginsList),
		status:  make(statusList),
	}
}

func (a *Adapter) Run(ctx context.Context) {
	a.runWathcer(ctx)
}

func (a *Adapter) Load(ctx context.Context, name string, command string) error {
	client := hcplugin.NewClient(&hcplugin.ClientConfig{
		HandshakeConfig: plugin.Handshake,
		Plugins:         plugin.PluginMap,
		Cmd:             exec.Command(command),
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
	a.plugins[name] = raw.(*plugin.GRPCClient)
	go func() {
		<-ctx.Done()
		client.Kill()
	}()
	return nil
}

func (a *Adapter) runWathcer(ctx context.Context) {
	t := time.NewTicker(time.Minute)
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			a.lockStatus.Lock()
			for k, v := range a.plugins {
				err := v.Health(ctx)
				if err != nil {
					a.log.Errorf("plugin manager watcher: %v", err)
				}
				a.status[k] = err == nil
			}
			a.lockStatus.Unlock()
		}
	}()
}

func (a *Adapter) Status() map[string]bool {
	status := make(map[string]bool)
	a.lockStatus.Lock()
	defer a.lockStatus.Unlock()
	for k, v := range a.status {
		status[k] = v
	}
	return status
}

func (a *Adapter) Get(name string) *plugin.GRPCClient {
	return a.plugins[name]
}
