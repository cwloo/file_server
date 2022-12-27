// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.7.0
// source: gate/gate.proto

package pb_gate

import (
	context "context"
	public "github.com/cwloo/uploader/proto/public"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_gate_gate_proto protoreflect.FileDescriptor

var file_gate_gate_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x67, 0x61, 0x74, 0x65, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x67, 0x61, 0x74, 0x65, 0x1a, 0x11, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x2f,
	0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x74, 0x0a, 0x04, 0x67, 0x61,
	0x74, 0x65, 0x12, 0x32, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x12,
	0x11, 0x2e, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x1a, 0x12, 0x2e, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x2e, 0x52, 0x6f, 0x75, 0x74,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x12, 0x38, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x4e, 0x6f, 0x64,
	0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x13, 0x2e, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x2e, 0x4e,
	0x6f, 0x64, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x1a, 0x14, 0x2e, 0x70, 0x75, 0x62,
	0x6c, 0x69, 0x63, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70,
	0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63,
	0x77, 0x6c, 0x6f, 0x6f, 0x2f, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x3b, 0x70, 0x62, 0x5f, 0x67, 0x61, 0x74, 0x65,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_gate_gate_proto_goTypes = []interface{}{
	(*public.RouterReq)(nil),    // 0: public.RouterReq
	(*public.NodeInfoReq)(nil),  // 1: public.NodeInfoReq
	(*public.RouterResp)(nil),   // 2: public.RouterResp
	(*public.NodeInfoResp)(nil), // 3: public.NodeInfoResp
}
var file_gate_gate_proto_depIdxs = []int32{
	0, // 0: gate.gate.GetRouter:input_type -> public.RouterReq
	1, // 1: gate.gate.GetNodeInfo:input_type -> public.NodeInfoReq
	2, // 2: gate.gate.GetRouter:output_type -> public.RouterResp
	3, // 3: gate.gate.GetNodeInfo:output_type -> public.NodeInfoResp
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_gate_gate_proto_init() }
func file_gate_gate_proto_init() {
	if File_gate_gate_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gate_gate_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gate_gate_proto_goTypes,
		DependencyIndexes: file_gate_gate_proto_depIdxs,
	}.Build()
	File_gate_gate_proto = out.File
	file_gate_gate_proto_rawDesc = nil
	file_gate_gate_proto_goTypes = nil
	file_gate_gate_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// GateClient is the client API for Gate service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GateClient interface {
	GetRouter(ctx context.Context, in *public.RouterReq, opts ...grpc.CallOption) (*public.RouterResp, error)
	GetNodeInfo(ctx context.Context, in *public.NodeInfoReq, opts ...grpc.CallOption) (*public.NodeInfoResp, error)
}

type gateClient struct {
	cc grpc.ClientConnInterface
}

func NewGateClient(cc grpc.ClientConnInterface) GateClient {
	return &gateClient{cc}
}

func (c *gateClient) GetRouter(ctx context.Context, in *public.RouterReq, opts ...grpc.CallOption) (*public.RouterResp, error) {
	out := new(public.RouterResp)
	err := c.cc.Invoke(ctx, "/gate.gate/GetRouter", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gateClient) GetNodeInfo(ctx context.Context, in *public.NodeInfoReq, opts ...grpc.CallOption) (*public.NodeInfoResp, error) {
	out := new(public.NodeInfoResp)
	err := c.cc.Invoke(ctx, "/gate.gate/GetNodeInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GateServer is the server API for Gate service.
type GateServer interface {
	GetRouter(context.Context, *public.RouterReq) (*public.RouterResp, error)
	GetNodeInfo(context.Context, *public.NodeInfoReq) (*public.NodeInfoResp, error)
}

// UnimplementedGateServer can be embedded to have forward compatible implementations.
type UnimplementedGateServer struct {
}

func (*UnimplementedGateServer) GetRouter(context.Context, *public.RouterReq) (*public.RouterResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRouter not implemented")
}
func (*UnimplementedGateServer) GetNodeInfo(context.Context, *public.NodeInfoReq) (*public.NodeInfoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNodeInfo not implemented")
}

func RegisterGateServer(s *grpc.Server, srv GateServer) {
	s.RegisterService(&_Gate_serviceDesc, srv)
}

func _Gate_GetRouter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(public.RouterReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GateServer).GetRouter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gate.gate/GetRouter",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GateServer).GetRouter(ctx, req.(*public.RouterReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Gate_GetNodeInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(public.NodeInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GateServer).GetNodeInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gate.gate/GetNodeInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GateServer).GetNodeInfo(ctx, req.(*public.NodeInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Gate_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gate.gate",
	HandlerType: (*GateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRouter",
			Handler:    _Gate_GetRouter_Handler,
		},
		{
			MethodName: "GetNodeInfo",
			Handler:    _Gate_GetNodeInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gate/gate.proto",
}
