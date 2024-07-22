// Package shared contains shared data between the host and plugins.
package adapter

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/unchainio/interfaces/adapter/proto"
)

// TriggerHandshake is a common handshake that is shared by plugin and host.
var TriggerHandshake = plugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "ADAPTER_PLUGIN",
	MagicCookieValue: TriggerComponent,
}

// TriggerPluginMap is the map of plugins we can dispense.
var TriggerPluginMap = map[string]plugin.Plugin{
	TriggerComponent: &TriggerPlugin{},
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type TriggerPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	// Concrete implementation, written in Go. This is  only used for plugins
	// that are written in Go.
	Impl Trigger
}

func (p *TriggerPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterTriggerServer(s, &GRPCTriggerServer{
		Impl:   p.Impl,
		broker: broker,
	})
	return nil
}

func (p *TriggerPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCTriggerClient{
		client: proto.NewTriggerClient(c),
		broker: broker,
	}, nil
}

var _ plugin.GRPCPlugin = &TriggerPlugin{}
