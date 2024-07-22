package adapter

import "github.com/hashicorp/go-plugin"

func StartTrigger(trigger Trigger) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: TriggerHandshake,
		Plugins: map[string]plugin.Plugin{
			"trigger": &TriggerPlugin{Impl: trigger},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
