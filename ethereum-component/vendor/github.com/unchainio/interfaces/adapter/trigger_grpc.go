package adapter

import (
	"encoding/json"
	"errors"

	"github.com/hashicorp/go-plugin"
	"github.com/unchainio/interfaces/adapter/proto"
	"golang.org/x/net/context"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCTriggerClient struct {
	broker *plugin.GRPCBroker
	client proto.TriggerClient
}

func (m *GRPCTriggerClient) Init(stub Stub, cfg []byte) error {
	brokerID, closer := SetupStubServer(stub, m.broker)
	_ = closer
	//defer closer()

	_, err := m.client.Init(context.Background(), &proto.InitTriggerRequest{
		StubServer: brokerID,
		Config:     cfg,
	})

	return err
}

func (m *GRPCTriggerClient) Trigger() (string, []byte, error) {
	r, err := m.client.Trigger(context.Background(), &proto.TriggerRequest{})

	if err != nil {
		return "", nil, err
	}

	return r.Tag, r.Message, nil
}

func (m *GRPCTriggerClient) Respond(tag string, response []byte, responseError error) error {
	_, err := m.client.Respond(context.Background(), &proto.RespondRequest{
		Tag:      tag,
		Response: response,
		Error:    responseError.Error(),
	})

	return err
}

func (m *GRPCTriggerClient) Close() error {
	_, err := m.client.Close(context.Background(), &proto.CloseRequest{})

	return err
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCTriggerServer struct {
	// This is the real implementation
	Impl   Trigger
	broker *plugin.GRPCBroker
}

func (m *GRPCTriggerServer) Init(ctx context.Context, req *proto.InitTriggerRequest) (*proto.InitTriggerResponse, error) {
	stub, closer, err := SetupStubClient(m.broker, req.StubServer)

	if err != nil {
		return nil, err
	}

	_ = closer
	//defer closer()

	return &proto.InitTriggerResponse{}, m.Impl.Init(stub, req.Config)
}

func (m *GRPCTriggerServer) Trigger(ctx context.Context, req *proto.TriggerRequest) (*proto.TriggerResponse, error) {
	tag, r, err := m.Impl.Trigger()

	if err != nil {
		return nil, err
	}

	rBytes, err := json.Marshal(r)

	if err != nil {
		return nil, err
	}

	return &proto.TriggerResponse{
		Tag:     tag,
		Message: rBytes,
	}, nil
}

func (m *GRPCTriggerServer) Respond(ctx context.Context, req *proto.RespondRequest) (*proto.RespondResponse, error) {
	response := make(map[string]interface{})
	err := json.Unmarshal(req.Response, response)

	if err != nil {
		return nil, err
	}

	return &proto.RespondResponse{}, m.Impl.Respond(req.Tag, response, errors.New(req.Error))
}

func (m *GRPCTriggerServer) Close(ctx context.Context, req *proto.CloseRequest) (*proto.CloseResponse, error) {
	return &proto.CloseResponse{}, m.Impl.Close()
}
