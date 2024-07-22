package adapter

import "github.com/hashicorp/go-plugin"

func StartEndpoint(endpoint Endpoint) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: EndpointHandshake,
		Plugins: map[string]plugin.Plugin{
			"endpoint": &EndpointPlugin{Impl: endpoint},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
