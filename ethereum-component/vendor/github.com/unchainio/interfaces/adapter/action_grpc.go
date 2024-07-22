package adapter

import (
	"encoding/json"

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
	brokerID, closer := SetupStubServer(stub, m.broker)

	_ = closer
	//defer closer()

	_, err := m.client.Init(context.Background(), &proto.InitActionRequest{
		StubServer: brokerID,
		Config:     cfg,
	})

	return err
}

func (m *GRPCActionClient) Invoke(inputMessage []byte) (outputMessage []byte, err error) {
	msg, err := m.client.Invoke(context.Background(), &proto.InvokeRequest{
		Message: inputMessage,
	})

	if err != nil {
		return nil, err
	}

	return msg.Message, nil
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

	_ = closer
	//defer closer()

	return &proto.InitActionResponse{}, m.Impl.Init(stub, req.Config)
}

func (m *GRPCActionServer) Invoke(ctx context.Context, req *proto.InvokeRequest) (*proto.InvokeResponse, error) {
	message := make(map[string]interface{})

	err := json.Unmarshal(req.Message, message)

	if err != nil {
		return nil, err
	}

	omsg, err := m.Impl.Invoke(message)

	omsgBytes, err := json.Marshal(omsg)

	if err != nil {
		return nil, err
	}

	return &proto.InvokeResponse{
		Message: omsgBytes,
	}, err
}
