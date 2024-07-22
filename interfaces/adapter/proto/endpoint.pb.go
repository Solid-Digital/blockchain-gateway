// Code generated by protoc-gen-go. DO NOT EDIT.
// source: endpoint.proto

package proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type InitEndpointRequest struct {
	StubServer           uint32   `protobuf:"varint,1,opt,name=stub_server,json=stubServer,proto3" json:"stub_server,omitempty"`
	Config               []byte   `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InitEndpointRequest) Reset()         { *m = InitEndpointRequest{} }
func (m *InitEndpointRequest) String() string { return proto.CompactTextString(m) }
func (*InitEndpointRequest) ProtoMessage()    {}
func (*InitEndpointRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_endpoint_7600ba83ec0d42ef, []int{0}
}
func (m *InitEndpointRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InitEndpointRequest.Unmarshal(m, b)
}
func (m *InitEndpointRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InitEndpointRequest.Marshal(b, m, deterministic)
}
func (dst *InitEndpointRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InitEndpointRequest.Merge(dst, src)
}
func (m *InitEndpointRequest) XXX_Size() int {
	return xxx_messageInfo_InitEndpointRequest.Size(m)
}
func (m *InitEndpointRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InitEndpointRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InitEndpointRequest proto.InternalMessageInfo

func (m *InitEndpointRequest) GetStubServer() uint32 {
	if m != nil {
		return m.StubServer
	}
	return 0
}

func (m *InitEndpointRequest) GetConfig() []byte {
	if m != nil {
		return m.Config
	}
	return nil
}

type InitEndpointResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InitEndpointResponse) Reset()         { *m = InitEndpointResponse{} }
func (m *InitEndpointResponse) String() string { return proto.CompactTextString(m) }
func (*InitEndpointResponse) ProtoMessage()    {}
func (*InitEndpointResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_endpoint_7600ba83ec0d42ef, []int{1}
}
func (m *InitEndpointResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InitEndpointResponse.Unmarshal(m, b)
}
func (m *InitEndpointResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InitEndpointResponse.Marshal(b, m, deterministic)
}
func (dst *InitEndpointResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InitEndpointResponse.Merge(dst, src)
}
func (m *InitEndpointResponse) XXX_Size() int {
	return xxx_messageInfo_InitEndpointResponse.Size(m)
}
func (m *InitEndpointResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_InitEndpointResponse.DiscardUnknown(m)
}

var xxx_messageInfo_InitEndpointResponse proto.InternalMessageInfo

type SendRequest struct {
	StubServer           uint32          `protobuf:"varint,1,opt,name=stub_server,json=stubServer,proto3" json:"stub_server,omitempty"`
	Message              *AdapterMessage `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *SendRequest) Reset()         { *m = SendRequest{} }
func (m *SendRequest) String() string { return proto.CompactTextString(m) }
func (*SendRequest) ProtoMessage()    {}
func (*SendRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_endpoint_7600ba83ec0d42ef, []int{2}
}
func (m *SendRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendRequest.Unmarshal(m, b)
}
func (m *SendRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendRequest.Marshal(b, m, deterministic)
}
func (dst *SendRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendRequest.Merge(dst, src)
}
func (m *SendRequest) XXX_Size() int {
	return xxx_messageInfo_SendRequest.Size(m)
}
func (m *SendRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SendRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SendRequest proto.InternalMessageInfo

func (m *SendRequest) GetStubServer() uint32 {
	if m != nil {
		return m.StubServer
	}
	return 0
}

func (m *SendRequest) GetMessage() *AdapterMessage {
	if m != nil {
		return m.Message
	}
	return nil
}

type SendResponse struct {
	Response             *AdapterMessage `protobuf:"bytes,1,opt,name=response,proto3" json:"response,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *SendResponse) Reset()         { *m = SendResponse{} }
func (m *SendResponse) String() string { return proto.CompactTextString(m) }
func (*SendResponse) ProtoMessage()    {}
func (*SendResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_endpoint_7600ba83ec0d42ef, []int{3}
}
func (m *SendResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendResponse.Unmarshal(m, b)
}
func (m *SendResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendResponse.Marshal(b, m, deterministic)
}
func (dst *SendResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendResponse.Merge(dst, src)
}
func (m *SendResponse) XXX_Size() int {
	return xxx_messageInfo_SendResponse.Size(m)
}
func (m *SendResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SendResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SendResponse proto.InternalMessageInfo

func (m *SendResponse) GetResponse() *AdapterMessage {
	if m != nil {
		return m.Response
	}
	return nil
}

type ReceiveRequest struct {
	StubServer           uint32   `protobuf:"varint,1,opt,name=stub_server,json=stubServer,proto3" json:"stub_server,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReceiveRequest) Reset()         { *m = ReceiveRequest{} }
func (m *ReceiveRequest) String() string { return proto.CompactTextString(m) }
func (*ReceiveRequest) ProtoMessage()    {}
func (*ReceiveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_endpoint_7600ba83ec0d42ef, []int{4}
}
func (m *ReceiveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReceiveRequest.Unmarshal(m, b)
}
func (m *ReceiveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReceiveRequest.Marshal(b, m, deterministic)
}
func (dst *ReceiveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReceiveRequest.Merge(dst, src)
}
func (m *ReceiveRequest) XXX_Size() int {
	return xxx_messageInfo_ReceiveRequest.Size(m)
}
func (m *ReceiveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ReceiveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ReceiveRequest proto.InternalMessageInfo

func (m *ReceiveRequest) GetStubServer() uint32 {
	if m != nil {
		return m.StubServer
	}
	return 0
}

type ReceiveResponse struct {
	Message              *TaggedAdapterMessage `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *ReceiveResponse) Reset()         { *m = ReceiveResponse{} }
func (m *ReceiveResponse) String() string { return proto.CompactTextString(m) }
func (*ReceiveResponse) ProtoMessage()    {}
func (*ReceiveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_endpoint_7600ba83ec0d42ef, []int{5}
}
func (m *ReceiveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReceiveResponse.Unmarshal(m, b)
}
func (m *ReceiveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReceiveResponse.Marshal(b, m, deterministic)
}
func (dst *ReceiveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReceiveResponse.Merge(dst, src)
}
func (m *ReceiveResponse) XXX_Size() int {
	return xxx_messageInfo_ReceiveResponse.Size(m)
}
func (m *ReceiveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ReceiveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ReceiveResponse proto.InternalMessageInfo

func (m *ReceiveResponse) GetMessage() *TaggedAdapterMessage {
	if m != nil {
		return m.Message
	}
	return nil
}

type AckRequest struct {
	StubServer           uint32          `protobuf:"varint,1,opt,name=stub_server,json=stubServer,proto3" json:"stub_server,omitempty"`
	Tag                  uint64          `protobuf:"varint,2,opt,name=tag,proto3" json:"tag,omitempty"`
	Response             *AdapterMessage `protobuf:"bytes,3,opt,name=response,proto3" json:"response,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *AckRequest) Reset()         { *m = AckRequest{} }
func (m *AckRequest) String() string { return proto.CompactTextString(m) }
func (*AckRequest) ProtoMessage()    {}
func (*AckRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_endpoint_7600ba83ec0d42ef, []int{6}
}
func (m *AckRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AckRequest.Unmarshal(m, b)
}
func (m *AckRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AckRequest.Marshal(b, m, deterministic)
}
func (dst *AckRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AckRequest.Merge(dst, src)
}
func (m *AckRequest) XXX_Size() int {
	return xxx_messageInfo_AckRequest.Size(m)
}
func (m *AckRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AckRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AckRequest proto.InternalMessageInfo

func (m *AckRequest) GetStubServer() uint32 {
	if m != nil {
		return m.StubServer
	}
	return 0
}

func (m *AckRequest) GetTag() uint64 {
	if m != nil {
		return m.Tag
	}
	return 0
}

func (m *AckRequest) GetResponse() *AdapterMessage {
	if m != nil {
		return m.Response
	}
	return nil
}

type AckResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AckResponse) Reset()         { *m = AckResponse{} }
func (m *AckResponse) String() string { return proto.CompactTextString(m) }
func (*AckResponse) ProtoMessage()    {}
func (*AckResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_endpoint_7600ba83ec0d42ef, []int{7}
}
func (m *AckResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AckResponse.Unmarshal(m, b)
}
func (m *AckResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AckResponse.Marshal(b, m, deterministic)
}
func (dst *AckResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AckResponse.Merge(dst, src)
}
func (m *AckResponse) XXX_Size() int {
	return xxx_messageInfo_AckResponse.Size(m)
}
func (m *AckResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AckResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AckResponse proto.InternalMessageInfo

type NackRequest struct {
	StubServer           uint32   `protobuf:"varint,1,opt,name=stub_server,json=stubServer,proto3" json:"stub_server,omitempty"`
	Tag                  uint64   `protobuf:"varint,2,opt,name=tag,proto3" json:"tag,omitempty"`
	Error                string   `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NackRequest) Reset()         { *m = NackRequest{} }
func (m *NackRequest) String() string { return proto.CompactTextString(m) }
func (*NackRequest) ProtoMessage()    {}
func (*NackRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_endpoint_7600ba83ec0d42ef, []int{8}
}
func (m *NackRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NackRequest.Unmarshal(m, b)
}
func (m *NackRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NackRequest.Marshal(b, m, deterministic)
}
func (dst *NackRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NackRequest.Merge(dst, src)
}
func (m *NackRequest) XXX_Size() int {
	return xxx_messageInfo_NackRequest.Size(m)
}
func (m *NackRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NackRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NackRequest proto.InternalMessageInfo

func (m *NackRequest) GetStubServer() uint32 {
	if m != nil {
		return m.StubServer
	}
	return 0
}

func (m *NackRequest) GetTag() uint64 {
	if m != nil {
		return m.Tag
	}
	return 0
}

func (m *NackRequest) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type NackResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NackResponse) Reset()         { *m = NackResponse{} }
func (m *NackResponse) String() string { return proto.CompactTextString(m) }
func (*NackResponse) ProtoMessage()    {}
func (*NackResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_endpoint_7600ba83ec0d42ef, []int{9}
}
func (m *NackResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NackResponse.Unmarshal(m, b)
}
func (m *NackResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NackResponse.Marshal(b, m, deterministic)
}
func (dst *NackResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NackResponse.Merge(dst, src)
}
func (m *NackResponse) XXX_Size() int {
	return xxx_messageInfo_NackResponse.Size(m)
}
func (m *NackResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_NackResponse.DiscardUnknown(m)
}

var xxx_messageInfo_NackResponse proto.InternalMessageInfo

type CloseRequest struct {
	StubServer           uint32   `protobuf:"varint,1,opt,name=stub_server,json=stubServer,proto3" json:"stub_server,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CloseRequest) Reset()         { *m = CloseRequest{} }
func (m *CloseRequest) String() string { return proto.CompactTextString(m) }
func (*CloseRequest) ProtoMessage()    {}
func (*CloseRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_endpoint_7600ba83ec0d42ef, []int{10}
}
func (m *CloseRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CloseRequest.Unmarshal(m, b)
}
func (m *CloseRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CloseRequest.Marshal(b, m, deterministic)
}
func (dst *CloseRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CloseRequest.Merge(dst, src)
}
func (m *CloseRequest) XXX_Size() int {
	return xxx_messageInfo_CloseRequest.Size(m)
}
func (m *CloseRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CloseRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CloseRequest proto.InternalMessageInfo

func (m *CloseRequest) GetStubServer() uint32 {
	if m != nil {
		return m.StubServer
	}
	return 0
}

type CloseResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CloseResponse) Reset()         { *m = CloseResponse{} }
func (m *CloseResponse) String() string { return proto.CompactTextString(m) }
func (*CloseResponse) ProtoMessage()    {}
func (*CloseResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_endpoint_7600ba83ec0d42ef, []int{11}
}
func (m *CloseResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CloseResponse.Unmarshal(m, b)
}
func (m *CloseResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CloseResponse.Marshal(b, m, deterministic)
}
func (dst *CloseResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CloseResponse.Merge(dst, src)
}
func (m *CloseResponse) XXX_Size() int {
	return xxx_messageInfo_CloseResponse.Size(m)
}
func (m *CloseResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CloseResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CloseResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*InitEndpointRequest)(nil), "proto.InitEndpointRequest")
	proto.RegisterType((*InitEndpointResponse)(nil), "proto.InitEndpointResponse")
	proto.RegisterType((*SendRequest)(nil), "proto.SendRequest")
	proto.RegisterType((*SendResponse)(nil), "proto.SendResponse")
	proto.RegisterType((*ReceiveRequest)(nil), "proto.ReceiveRequest")
	proto.RegisterType((*ReceiveResponse)(nil), "proto.ReceiveResponse")
	proto.RegisterType((*AckRequest)(nil), "proto.AckRequest")
	proto.RegisterType((*AckResponse)(nil), "proto.AckResponse")
	proto.RegisterType((*NackRequest)(nil), "proto.NackRequest")
	proto.RegisterType((*NackResponse)(nil), "proto.NackResponse")
	proto.RegisterType((*CloseRequest)(nil), "proto.CloseRequest")
	proto.RegisterType((*CloseResponse)(nil), "proto.CloseResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EndpointClient is the client API for Endpoint service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EndpointClient interface {
	Init(ctx context.Context, in *InitEndpointRequest, opts ...grpc.CallOption) (*InitEndpointResponse, error)
	Send(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*SendResponse, error)
	Receive(ctx context.Context, in *ReceiveRequest, opts ...grpc.CallOption) (*ReceiveResponse, error)
	Ack(ctx context.Context, in *AckRequest, opts ...grpc.CallOption) (*AckResponse, error)
	Nack(ctx context.Context, in *NackRequest, opts ...grpc.CallOption) (*NackResponse, error)
	Close(ctx context.Context, in *CloseRequest, opts ...grpc.CallOption) (*CloseResponse, error)
}

type endpointClient struct {
	cc *grpc.ClientConn
}

func NewEndpointClient(cc *grpc.ClientConn) EndpointClient {
	return &endpointClient{cc}
}

func (c *endpointClient) Init(ctx context.Context, in *InitEndpointRequest, opts ...grpc.CallOption) (*InitEndpointResponse, error) {
	out := new(InitEndpointResponse)
	err := c.cc.Invoke(ctx, "/proto.Endpoint/Init", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointClient) Send(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*SendResponse, error) {
	out := new(SendResponse)
	err := c.cc.Invoke(ctx, "/proto.Endpoint/Send", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointClient) Receive(ctx context.Context, in *ReceiveRequest, opts ...grpc.CallOption) (*ReceiveResponse, error) {
	out := new(ReceiveResponse)
	err := c.cc.Invoke(ctx, "/proto.Endpoint/Receive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointClient) Ack(ctx context.Context, in *AckRequest, opts ...grpc.CallOption) (*AckResponse, error) {
	out := new(AckResponse)
	err := c.cc.Invoke(ctx, "/proto.Endpoint/Ack", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointClient) Nack(ctx context.Context, in *NackRequest, opts ...grpc.CallOption) (*NackResponse, error) {
	out := new(NackResponse)
	err := c.cc.Invoke(ctx, "/proto.Endpoint/Nack", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointClient) Close(ctx context.Context, in *CloseRequest, opts ...grpc.CallOption) (*CloseResponse, error) {
	out := new(CloseResponse)
	err := c.cc.Invoke(ctx, "/proto.Endpoint/Close", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EndpointServer is the server API for Endpoint service.
type EndpointServer interface {
	Init(context.Context, *InitEndpointRequest) (*InitEndpointResponse, error)
	Send(context.Context, *SendRequest) (*SendResponse, error)
	Receive(context.Context, *ReceiveRequest) (*ReceiveResponse, error)
	Ack(context.Context, *AckRequest) (*AckResponse, error)
	Nack(context.Context, *NackRequest) (*NackResponse, error)
	Close(context.Context, *CloseRequest) (*CloseResponse, error)
}

func RegisterEndpointServer(s *grpc.Server, srv EndpointServer) {
	s.RegisterService(&_Endpoint_serviceDesc, srv)
}

func _Endpoint_Init_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitEndpointRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointServer).Init(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Endpoint/Init",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointServer).Init(ctx, req.(*InitEndpointRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Endpoint_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Endpoint/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointServer).Send(ctx, req.(*SendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Endpoint_Receive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReceiveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointServer).Receive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Endpoint/Receive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointServer).Receive(ctx, req.(*ReceiveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Endpoint_Ack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointServer).Ack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Endpoint/Ack",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointServer).Ack(ctx, req.(*AckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Endpoint_Nack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointServer).Nack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Endpoint/Nack",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointServer).Nack(ctx, req.(*NackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Endpoint_Close_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CloseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointServer).Close(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Endpoint/Close",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointServer).Close(ctx, req.(*CloseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Endpoint_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Endpoint",
	HandlerType: (*EndpointServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Init",
			Handler:    _Endpoint_Init_Handler,
		},
		{
			MethodName: "Send",
			Handler:    _Endpoint_Send_Handler,
		},
		{
			MethodName: "Receive",
			Handler:    _Endpoint_Receive_Handler,
		},
		{
			MethodName: "Ack",
			Handler:    _Endpoint_Ack_Handler,
		},
		{
			MethodName: "Nack",
			Handler:    _Endpoint_Nack_Handler,
		},
		{
			MethodName: "Close",
			Handler:    _Endpoint_Close_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "endpoint.proto",
}

func init() { proto.RegisterFile("endpoint.proto", fileDescriptor_endpoint_7600ba83ec0d42ef) }

var fileDescriptor_endpoint_7600ba83ec0d42ef = []byte{
	// 402 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x52, 0x4d, 0x6f, 0xda, 0x40,
	0x10, 0x15, 0xdf, 0x74, 0x6c, 0x43, 0xbb, 0x50, 0x84, 0xcc, 0xa1, 0xc8, 0x27, 0x0e, 0x15, 0x08,
	0xaa, 0x4a, 0xbd, 0x55, 0x56, 0x55, 0xa9, 0x3d, 0x84, 0x83, 0x89, 0x72, 0x45, 0xc6, 0x9e, 0x58,
	0x16, 0x89, 0xed, 0xec, 0x1a, 0xfe, 0x63, 0xfe, 0x55, 0xe4, 0xdd, 0xb1, 0xb1, 0x23, 0x90, 0x2c,
	0xe5, 0x64, 0xef, 0xdb, 0x99, 0x79, 0x6f, 0xdf, 0x3c, 0x18, 0x60, 0xe4, 0x27, 0x71, 0x18, 0xa5,
	0xcb, 0x84, 0xc7, 0x69, 0xcc, 0x3a, 0xf2, 0x63, 0x1a, 0xcf, 0x28, 0x84, 0x1b, 0xa0, 0x42, 0xad,
	0x2d, 0x8c, 0xfe, 0x47, 0x61, 0xfa, 0x97, 0x6a, 0x1d, 0x7c, 0x39, 0xa1, 0x48, 0xd9, 0x37, 0xd0,
	0x44, 0x7a, 0x3a, 0xec, 0x05, 0xf2, 0x33, 0xf2, 0x69, 0x63, 0xde, 0x58, 0x18, 0x0e, 0x64, 0xd0,
	0x4e, 0x22, 0x6c, 0x02, 0x5d, 0x2f, 0x8e, 0x1e, 0xc3, 0x60, 0xda, 0x9c, 0x37, 0x16, 0xba, 0x43,
	0x27, 0x6b, 0x02, 0xe3, 0xea, 0x3c, 0x91, 0xc4, 0x91, 0x40, 0x6b, 0x0f, 0xda, 0x0e, 0x23, 0xbf,
	0xf6, 0xfc, 0x15, 0xf4, 0x48, 0xa8, 0x24, 0xd0, 0x36, 0x5f, 0x95, 0xe0, 0xa5, 0xed, 0xbb, 0x49,
	0x8a, 0xfc, 0x4e, 0x5d, 0x3a, 0x79, 0x95, 0x65, 0x83, 0xae, 0x08, 0x14, 0x21, 0x5b, 0x43, 0x9f,
	0xd3, 0xbf, 0x1c, 0x7f, 0x73, 0x42, 0x51, 0x66, 0xad, 0x61, 0xe0, 0xa0, 0x87, 0xe1, 0x19, 0xeb,
	0xca, 0xb4, 0xfe, 0xc1, 0xb0, 0x68, 0x21, 0xe2, 0x9f, 0x17, 0xe5, 0x8a, 0x77, 0x46, 0xbc, 0xf7,
	0x6e, 0x10, 0xa0, 0x7f, 0x4b, 0x3f, 0x07, 0xb0, 0xbd, 0x63, 0x6d, 0x7f, 0x3e, 0x43, 0x2b, 0x75,
	0x95, 0xf9, 0x6d, 0x27, 0xfb, 0xad, 0x3c, 0xb8, 0x55, 0xef, 0xc1, 0x06, 0x68, 0x92, 0x93, 0x8e,
	0x0f, 0xa0, 0x6d, 0xdd, 0x0f, 0x69, 0x18, 0x43, 0x07, 0x39, 0x8f, 0xb9, 0x14, 0xf0, 0xc9, 0x51,
	0x07, 0x6b, 0x00, 0xba, 0x9a, 0x4b, 0x3c, 0x2b, 0xd0, 0xff, 0x3c, 0xc5, 0xa2, 0xbe, 0xcb, 0x43,
	0x30, 0xa8, 0x41, 0x4d, 0xd8, 0xbc, 0x36, 0xa1, 0x9f, 0x47, 0x8c, 0xfd, 0x86, 0x76, 0x16, 0x39,
	0x66, 0xd2, 0x73, 0xaf, 0xe4, 0xd9, 0x9c, 0x5d, 0xbd, 0xa3, 0x8d, 0xad, 0xa0, 0x9d, 0x45, 0x87,
	0x31, 0x2a, 0x2a, 0x05, 0xd5, 0x1c, 0x55, 0x30, 0x6a, 0xf8, 0x05, 0x3d, 0xda, 0x3a, 0xcb, 0x3d,
	0xae, 0x06, 0xc7, 0x9c, 0xbc, 0x87, 0xa9, 0xf3, 0x3b, 0xb4, 0x6c, 0xef, 0xc8, 0xbe, 0xe4, 0x9b,
	0x29, 0xdc, 0x36, 0x59, 0x19, 0xba, 0x08, 0xcb, 0x8c, 0x2b, 0x84, 0x95, 0xb6, 0x53, 0x08, 0x2b,
	0x3b, 0xcb, 0x36, 0xd0, 0x91, 0x46, 0xb1, 0xfc, 0xb6, 0xec, 0xb3, 0x39, 0xae, 0x82, 0xaa, 0xe7,
	0xd0, 0x95, 0xe0, 0x8f, 0xb7, 0x00, 0x00, 0x00, 0xff, 0xff, 0x70, 0x34, 0x4b, 0xa4, 0x30, 0x04,
	0x00, 0x00,
}