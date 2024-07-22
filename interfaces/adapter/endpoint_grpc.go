package adapter

import (
	"errors"

	"github.com/hashicorp/go-plugin"
	"github.com/unchainio/interfaces/adapter/proto"
	"golang.org/x/net/context"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCEndpointClient struct {
	broker *plugin.GRPCBroker
	client proto.EndpointClient
}

func (m *GRPCEndpointClient) Init(stub Stub, cfg []byte) error {
	brokerID, closer := SetupStubServer(stub, m.broker)
	defer closer()

	_, err := m.client.Init(context.Background(), &proto.InitEndpointRequest{
		StubServer: brokerID,
		Config:     cfg,
	})

	return err
}

func (m *GRPCEndpointClient) Send(stub Stub, message *Message) (*Message, error) {
	brokerID, closer := SetupStubServer(stub, m.broker)
	defer closer()

	r, err := m.client.Send(context.Background(), &proto.SendRequest{
		StubServer: brokerID,
		Message: &proto.AdapterMessage{
			Body:       message.Body,
			Attributes: message.Attributes,
		},
	})

	if err != nil {
		return nil, err
	}

	return &Message{
		Body:       r.Response.Body,
		Attributes: r.Response.Attributes,
	}, nil
}

func (m *GRPCEndpointClient) Receive(stub Stub) (*TaggedMessage, error) {
	brokerID, closer := SetupStubServer(stub, m.broker)
	defer closer()

	r, err := m.client.Receive(context.Background(), &proto.ReceiveRequest{
		StubServer: brokerID,
	})

	if err != nil {
		return nil, err
	}

	return &TaggedMessage{
		Tag: r.Message.Tag,
		Message: &Message{
			Body:       r.Message.Message.Body,
			Attributes: r.Message.Message.Attributes,
		},
	}, nil
}

func (m *GRPCEndpointClient) Ack(stub Stub, tag uint64, response *Message) error {
	brokerID, closer := SetupStubServer(stub, m.broker)
	defer closer()

	_, err := m.client.Ack(context.Background(), &proto.AckRequest{
		StubServer: brokerID,
		Tag:        tag,
		Response: &proto.AdapterMessage{
			Body:       response.Body,
			Attributes: response.Attributes,
		},
	})

	return err
}

func (m *GRPCEndpointClient) Nack(stub Stub, tag uint64, responseError error) error {
	brokerID, closer := SetupStubServer(stub, m.broker)
	defer closer()

	_, err := m.client.Nack(context.Background(), &proto.NackRequest{
		StubServer: brokerID,
		Tag:        tag,
		Error:      responseError.Error(),
	})

	return err
}

func (m *GRPCEndpointClient) Close(stub Stub) error {
	brokerID, closer := SetupStubServer(stub, m.broker)
	defer closer()

	_, err := m.client.Close(context.Background(), &proto.CloseRequest{
		StubServer: brokerID,
	})

	return err
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCEndpointServer struct {
	// This is the real implementation
	Impl   Endpoint
	broker *plugin.GRPCBroker
}

func (m *GRPCEndpointServer) Init(ctx context.Context, req *proto.InitEndpointRequest) (*proto.InitEndpointResponse, error) {
	stub, closer, err := SetupStubClient(m.broker, req.StubServer)

	if err != nil {
		return nil, err
	}

	defer closer()

	return &proto.InitEndpointResponse{}, m.Impl.Init(stub, req.Config)
}

func (m *GRPCEndpointServer) Send(ctx context.Context, req *proto.SendRequest) (*proto.SendResponse, error) {
	stub, closer, err := SetupStubClient(m.broker, req.StubServer)

	if err != nil {
		return nil, err
	}

	defer closer()

	r, err := m.Impl.Send(stub, &Message{
		Body:       req.Message.Body,
		Attributes: req.Message.Attributes,
	})

	if err != nil {
		return nil, err
	}

	return &proto.SendResponse{
		Response: &proto.AdapterMessage{
			Body:       r.Body,
			Attributes: r.Attributes,
		},
	}, nil
}

func (m *GRPCEndpointServer) Receive(ctx context.Context, req *proto.ReceiveRequest) (*proto.ReceiveResponse, error) {
	stub, closer, err := SetupStubClient(m.broker, req.StubServer)

	if err != nil {
		return nil, err
	}

	defer closer()

	r, err := m.Impl.Receive(stub)

	if err != nil {
		return nil, err
	}

	return &proto.ReceiveResponse{
		Message: &proto.TaggedAdapterMessage{
			Tag: r.Tag,
			Message: &proto.AdapterMessage{
				Body:       r.Body,
				Attributes: r.Attributes,
			},
		},
	}, nil
}

func (m *GRPCEndpointServer) Ack(ctx context.Context, req *proto.AckRequest) (*proto.AckResponse, error) {
	stub, closer, err := SetupStubClient(m.broker, req.StubServer)

	if err != nil {
		return nil, err
	}

	defer closer()

	return &proto.AckResponse{}, m.Impl.Ack(stub, req.Tag, &Message{
		Body:       req.Response.Body,
		Attributes: req.Response.Attributes,
	})
}

func (m *GRPCEndpointServer) Nack(ctx context.Context, req *proto.NackRequest) (*proto.NackResponse, error) {
	stub, closer, err := SetupStubClient(m.broker, req.StubServer)

	if err != nil {
		return nil, err
	}

	defer closer()

	return &proto.NackResponse{}, m.Impl.Nack(stub, req.Tag, errors.New(req.Error))
}

func (m *GRPCEndpointServer) Close(ctx context.Context, req *proto.CloseRequest) (*proto.CloseResponse, error) {
	stub, closer, err := SetupStubClient(m.broker, req.StubServer)

	if err != nil {
		return nil, err
	}

	defer closer()

	return &proto.CloseResponse{}, m.Impl.Close(stub)
}
