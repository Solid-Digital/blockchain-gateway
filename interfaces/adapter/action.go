package adapter

import (
	plugin "github.com/hashicorp/go-plugin"
)

func StartAction(action Action) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: ActionHandshake,
		Plugins: map[string]plugin.Plugin{
			"action": &ActionPlugin{Impl: action},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
