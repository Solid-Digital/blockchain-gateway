// Package shared contains shared data between the host and plugins.
package adapter

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/unchainio/interfaces/adapter/proto"
)

// ActionHandshake is a common handshake that is shared by plugin and host.
var ActionHandshake = plugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "ADAPTER_PLUGIN",
	MagicCookieValue: ActionComponent,
}

// ActionPluginMap is the map of plugins we can dispense.
var ActionPluginMap = map[string]plugin.Plugin{
	ActionComponent: &ActionPlugin{},
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type ActionPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	// Concrete implementation, written in Go. This is  only used for plugins
	// that are written in Go.
	Impl Action
}

func (p *ActionPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterActionServer(s, &GRPCActionServer{
		Impl:   p.Impl,
		broker: broker,
	})
	return nil
}

func (p *ActionPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCActionClient{
		client: proto.NewActionClient(c),
		broker: broker,
	}, nil
}

var _ plugin.GRPCPlugin = &ActionPlugin{}
