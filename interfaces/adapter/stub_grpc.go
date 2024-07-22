package adapter

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-plugin"
	"github.com/unchainio/interfaces/adapter/proto"
	"google.golang.org/grpc"
)

func SetupStubServer(stub Stub, broker *plugin.GRPCBroker) (brokerID uint32, close func()) {
	stubHelperServer := &GRPCStubServer{Impl: stub}

	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		proto.RegisterStubHelperServer(s, stubHelperServer)

		return s
	}

	brokerID = broker.NextId()
	go broker.AcceptAndServe(brokerID, serverFunc)

	return brokerID, func() { s.Stop() }
}

func SetupStubClient(broker *plugin.GRPCBroker, brokerID uint32) (stub Stub, close func(), err error) {
	conn, err := broker.Dial(brokerID)
	if err != nil {
		return nil, nil, err
	}

	stub = &GRPCStubHelperClient{proto.NewStubHelperClient(conn)}

	return stub, func() { conn.Close() }, nil
}

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCStubHelperClient struct{ client proto.StubHelperClient }

func (m *GRPCStubHelperClient) Debugf(format string, v ...interface{}) {
	m.client.Debugf(context.Background(), &proto.LogRequest{
		Message: fmt.Sprintf(format, v...),
	})
}

func (m *GRPCStubHelperClient) Errorf(format string, v ...interface{}) {
	m.client.Errorf(context.Background(), &proto.LogRequest{
		Message: fmt.Sprintf(format, v...),
	})
}

func (m *GRPCStubHelperClient) Fatalf(format string, v ...interface{}) {
	m.client.Fatalf(context.Background(), &proto.LogRequest{
		Message: fmt.Sprintf(format, v...),
	})
}

func (m *GRPCStubHelperClient) Panicf(format string, v ...interface{}) {
	m.client.Panicf(context.Background(), &proto.LogRequest{
		Message: fmt.Sprintf(format, v...),
	})
}

func (m *GRPCStubHelperClient) Printf(format string, v ...interface{}) {
	m.client.Printf(context.Background(), &proto.LogRequest{
		Message: fmt.Sprintf(format, v...),
	})
}

func (m *GRPCStubHelperClient) Warnf(format string, v ...interface{}) {
	m.client.Warnf(context.Background(), &proto.LogRequest{
		Message: fmt.Sprintf(format, v...),
	})
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCStubServer struct {
	// This is the real implementation
	Impl Stub
}

func (m *GRPCStubServer) Printf(ctx context.Context, req *proto.LogRequest) (*proto.LogResponse, error) {
	m.Impl.Printf("%s", req.Message)

	return &proto.LogResponse{}, nil
}

func (m *GRPCStubServer) Fatalf(ctx context.Context, req *proto.LogRequest) (*proto.LogResponse, error) {
	m.Impl.Fatalf("%s", req.Message)

	return &proto.LogResponse{}, nil
}

func (m *GRPCStubServer) Panicf(ctx context.Context, req *proto.LogRequest) (*proto.LogResponse, error) {
	m.Impl.Panicf("%s", req.Message)

	return &proto.LogResponse{}, nil
}

func (m *GRPCStubServer) Debugf(ctx context.Context, req *proto.LogRequest) (*proto.LogResponse, error) {
	m.Impl.Debugf("%s", req.Message)

	return &proto.LogResponse{}, nil
}

func (m *GRPCStubServer) Warnf(ctx context.Context, req *proto.LogRequest) (*proto.LogResponse, error) {
	m.Impl.Warnf("%s", req.Message)

	return &proto.LogResponse{}, nil
}

func (m *GRPCStubServer) Errorf(ctx context.Context, req *proto.LogRequest) (*proto.LogResponse, error) {
	m.Impl.Errorf("%s", req.Message)

	return &proto.LogResponse{}, nil
}
