// Package shared contains shared data between the host and plugins.
package adapter

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/unchainio/interfaces/adapter/proto"
)

// EndpointHandshake is a common handshake that is shared by plugin and host.
var EndpointHandshake = plugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "ADAPTER_PLUGIN",
	MagicCookieValue: "endpoint",
}

// EndpointPluginMap is the map of plugins we can dispense.
var EndpointPluginMap = map[string]plugin.Plugin{
	"endpoint": &EndpointPlugin{},
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type EndpointPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	// Concrete implementation, written in Go. This is  only used for plugins
	// that are written in Go.
	Impl Endpoint
}

func (p *EndpointPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterEndpointServer(s, &GRPCEndpointServer{
		Impl:   p.Impl,
		broker: broker,
	})
	return nil
}

func (p *EndpointPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCEndpointClient{
		client: proto.NewEndpointClient(c),
		broker: broker,
	}, nil
}

var _ plugin.GRPCPlugin = &EndpointPlugin{}
