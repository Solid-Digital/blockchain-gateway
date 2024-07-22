package adapter

import (
	"github.com/hashicorp/go-plugin"
	"github.com/unchainio/interfaces/adapter/proto"
	"golang.org/x/net/context"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCActionClient struct {
	broker *plugin.GRPCBroker
	client proto.ActionClient
}

func (m *GRPCActionClient) Init(stub Stub, cfg []byte) error {

	_, err := m.client.Init(context.Background(), &proto.InitActionRequest{
		Config:     cfg,
	})

	return err
}


func (m *GRPCActionClient) Invoke(stub Stub, message *Message) error {
	brokerID, closer := SetupStubServer(stub, m.broker)
	defer closer()

	imsg, err := m.client.Invoke(context.Background(), &proto.InvokeRequest{
		StubServer: brokerID,
		Message: &proto.AdapterMessage{
			Body:       message.Body,
			Attributes: message.Attributes,
		},
	})

	if err != nil {
		return err
	}

	message.Body = imsg.Message.Body
	message.Attributes = imsg.Message.Attributes

	return nil
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCActionServer struct {
	// This is the real implementation
	Impl   Action
	broker *plugin.GRPCBroker
}

func (m *GRPCActionServer) Init(ctx context.Context, req *proto.InitActionRequest) (*proto.InitActionResponse, error) {
	stub, closer, err := SetupStubClient(m.broker, req.StubServer)

	if err != nil {
		return nil, err
	}

	defer closer()

	return &proto.InitActionResponse{}, m.Impl.Init(stub, req.Config)
}

func (m *GRPCActionServer) Invoke(ctx context.Context, req *proto.InvokeRequest) (*proto.InvokeResponse, error) {
	stub, closer, err := SetupStubClient(m.broker, req.StubServer)

	if err != nil {
		return nil, err
	}

	defer closer()

	msg := &Message{
		Body:       req.Message.Body,
		Attributes: req.Message.Attributes,
	}

	err = m.Impl.Invoke(stub, msg)

	return &proto.InvokeResponse{
		Message: &proto.AdapterMessage{
			Body:       msg.Body,
			Attributes: msg.Attributes,
		},
	}, err
}
