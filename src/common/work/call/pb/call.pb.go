// Code generated by protoc-gen-go. DO NOT EDIT.
// source: call.proto

package call

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Request struct {
	Method               string   `protobuf:"bytes,1,opt,name=method,proto3" json:"method,omitempty"`
	Params               string   `protobuf:"bytes,2,opt,name=params,proto3" json:"params,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_caa5955d5eab2d2d, []int{0}
}

func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *Request) GetParams() string {
	if m != nil {
		return m.Params
	}
	return ""
}

type Response struct {
	Result               string   `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_caa5955d5eab2d2d, []int{1}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

func init() {
	proto.RegisterType((*Request)(nil), "Request")
	proto.RegisterType((*Response)(nil), "Response")
}

func init() { proto.RegisterFile("call.proto", fileDescriptor_caa5955d5eab2d2d) }

var fileDescriptor_caa5955d5eab2d2d = []byte{
	// 188 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0x41, 0x8a, 0x83, 0x30,
	0x14, 0x86, 0x27, 0xc3, 0xe0, 0xe8, 0x5b, 0xa6, 0x50, 0x44, 0x28, 0x94, 0xac, 0x5c, 0x05, 0x69,
	0x17, 0xa5, 0xdb, 0x7a, 0x03, 0x3d, 0x41, 0xaa, 0x8f, 0x36, 0x10, 0x13, 0x9b, 0x3c, 0x7b, 0xbc,
	0x9e, 0xad, 0x58, 0x75, 0x5b, 0x5c, 0x85, 0xef, 0xe7, 0xfb, 0x08, 0x3c, 0x80, 0x46, 0x19, 0x23,
	0x7b, 0xef, 0xc8, 0x89, 0x33, 0xfc, 0x57, 0xf8, 0x18, 0x30, 0x10, 0xdf, 0x42, 0xd4, 0x21, 0xdd,
	0x5d, 0x9b, 0xb2, 0x3d, 0xcb, 0x93, 0x6a, 0xa6, 0x71, 0xef, 0x95, 0x57, 0x5d, 0x48, 0x7f, 0xa7,
	0x7d, 0x22, 0x21, 0x20, 0xae, 0x30, 0xf4, 0xce, 0x06, 0x1c, 0x1d, 0x8f, 0x61, 0x30, 0xb4, 0xb4,
	0x13, 0x1d, 0x5e, 0x0c, 0xfe, 0x4a, 0x65, 0x0c, 0xdf, 0xcd, 0x6f, 0x2c, 0xe7, 0xef, 0xb2, 0x44,
	0x2e, 0xb5, 0xf8, 0xe1, 0x12, 0x36, 0x35, 0xfa, 0x27, 0xfa, 0x9a, 0x3c, 0xaa, 0x4e, 0xdb, 0xdb,
	0x17, 0xbb, 0x60, 0xa3, 0x5f, 0x1a, 0x8d, 0x96, 0xd6, 0xf8, 0x39, 0xe3, 0x27, 0xc8, 0x2e, 0xba,
	0xd5, 0x1e, 0x1b, 0xd2, 0xce, 0x2a, 0xb3, 0x2e, 0x2b, 0xd8, 0x35, 0xfa, 0x9c, 0xe9, 0xf8, 0x0e,
	0x00, 0x00, 0xff, 0xff, 0x04, 0xb3, 0x14, 0x34, 0x34, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CallClient is the client API for Call service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CallClient interface {
	Call(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	ServerStreamingCall(ctx context.Context, in *Request, opts ...grpc.CallOption) (Call_ServerStreamingCallClient, error)
	ClientStreamingCall(ctx context.Context, opts ...grpc.CallOption) (Call_ClientStreamingCallClient, error)
	BidirectionalStreamingCall(ctx context.Context, opts ...grpc.CallOption) (Call_BidirectionalStreamingCallClient, error)
}

type callClient struct {
	cc *grpc.ClientConn
}

func NewCallClient(cc *grpc.ClientConn) CallClient {
	return &callClient{cc}
}

func (c *callClient) Call(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/Call/Call", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *callClient) ServerStreamingCall(ctx context.Context, in *Request, opts ...grpc.CallOption) (Call_ServerStreamingCallClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Call_serviceDesc.Streams[0], "/Call/ServerStreamingCall", opts...)
	if err != nil {
		return nil, err
	}
	x := &callServerStreamingCallClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Call_ServerStreamingCallClient interface {
	Recv() (*Response, error)
	grpc.ClientStream
}

type callServerStreamingCallClient struct {
	grpc.ClientStream
}

func (x *callServerStreamingCallClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *callClient) ClientStreamingCall(ctx context.Context, opts ...grpc.CallOption) (Call_ClientStreamingCallClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Call_serviceDesc.Streams[1], "/Call/ClientStreamingCall", opts...)
	if err != nil {
		return nil, err
	}
	x := &callClientStreamingCallClient{stream}
	return x, nil
}

type Call_ClientStreamingCallClient interface {
	Send(*Request) error
	CloseAndRecv() (*Response, error)
	grpc.ClientStream
}

type callClientStreamingCallClient struct {
	grpc.ClientStream
}

func (x *callClientStreamingCallClient) Send(m *Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *callClientStreamingCallClient) CloseAndRecv() (*Response, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *callClient) BidirectionalStreamingCall(ctx context.Context, opts ...grpc.CallOption) (Call_BidirectionalStreamingCallClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Call_serviceDesc.Streams[2], "/Call/BidirectionalStreamingCall", opts...)
	if err != nil {
		return nil, err
	}
	x := &callBidirectionalStreamingCallClient{stream}
	return x, nil
}

type Call_BidirectionalStreamingCallClient interface {
	Send(*Request) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type callBidirectionalStreamingCallClient struct {
	grpc.ClientStream
}

func (x *callBidirectionalStreamingCallClient) Send(m *Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *callBidirectionalStreamingCallClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CallServer is the server API for Call service.
type CallServer interface {
	Call(context.Context, *Request) (*Response, error)
	ServerStreamingCall(*Request, Call_ServerStreamingCallServer) error
	ClientStreamingCall(Call_ClientStreamingCallServer) error
	BidirectionalStreamingCall(Call_BidirectionalStreamingCallServer) error
}

func RegisterCallServer(s *grpc.Server, srv CallServer) {
	s.RegisterService(&_Call_serviceDesc, srv)
}

func _Call_Call_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CallServer).Call(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Call/Call",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CallServer).Call(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Call_ServerStreamingCall_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CallServer).ServerStreamingCall(m, &callServerStreamingCallServer{stream})
}

type Call_ServerStreamingCallServer interface {
	Send(*Response) error
	grpc.ServerStream
}

type callServerStreamingCallServer struct {
	grpc.ServerStream
}

func (x *callServerStreamingCallServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func _Call_ClientStreamingCall_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CallServer).ClientStreamingCall(&callClientStreamingCallServer{stream})
}

type Call_ClientStreamingCallServer interface {
	SendAndClose(*Response) error
	Recv() (*Request, error)
	grpc.ServerStream
}

type callClientStreamingCallServer struct {
	grpc.ServerStream
}

func (x *callClientStreamingCallServer) SendAndClose(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *callClientStreamingCallServer) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Call_BidirectionalStreamingCall_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CallServer).BidirectionalStreamingCall(&callBidirectionalStreamingCallServer{stream})
}

type Call_BidirectionalStreamingCallServer interface {
	Send(*Response) error
	Recv() (*Request, error)
	grpc.ServerStream
}

type callBidirectionalStreamingCallServer struct {
	grpc.ServerStream
}

func (x *callBidirectionalStreamingCallServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *callBidirectionalStreamingCallServer) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Call_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Call",
	HandlerType: (*CallServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Call",
			Handler:    _Call_Call_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ServerStreamingCall",
			Handler:       _Call_ServerStreamingCall_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ClientStreamingCall",
			Handler:       _Call_ClientStreamingCall_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "BidirectionalStreamingCall",
			Handler:       _Call_BidirectionalStreamingCall_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "call.proto",
}
