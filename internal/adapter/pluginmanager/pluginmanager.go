package pluginmanager

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	hcplugin "github.com/hashicorp/go-plugin"
	"github.com/k1nky/tookhook/pkg/plugin"
)

type PluginsList map[string]*plugin.GRPCClient

type Adapter struct {
	plugins PluginsList
}

func New() *Adapter {
	return &Adapter{
		plugins: make(PluginsList),
	}
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
	go func() {
		t := time.NewTicker(1 * time.Minute)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if err := rpcClient.Ping(); err != nil {
					fmt.Println(err)
				}
			}
		}
	}()
	return nil
}

func (svc *Adapter) Get(name string) *plugin.GRPCClient {
	return svc.plugins[name]
}
