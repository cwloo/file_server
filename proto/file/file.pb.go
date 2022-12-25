// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.7.0
// source: file/file.proto

package pb_file

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type NodeInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pid        int32         `protobuf:"varint,1,opt,name=pid,proto3" json:"pid,omitempty"`
	Name       string        `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Id         int32         `protobuf:"varint,3,opt,name=id,proto3" json:"id,omitempty"`
	NumOfLoads int32         `protobuf:"varint,4,opt,name=NumOfLoads,proto3" json:"NumOfLoads,omitempty"`
	Ip         string        `protobuf:"bytes,5,opt,name=ip,proto3" json:"ip,omitempty"`
	Port       int32         `protobuf:"varint,6,opt,name=port,proto3" json:"port,omitempty"`
	Domain     string        `protobuf:"bytes,7,opt,name=domain,proto3" json:"domain,omitempty"`
	Rpc        *NodeInfo_Rpc `protobuf:"bytes,8,opt,name=rpc,proto3" json:"rpc,omitempty"`
}

func (x *NodeInfo) Reset() {
	*x = NodeInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_file_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeInfo) ProtoMessage() {}

func (x *NodeInfo) ProtoReflect() protoreflect.Message {
	mi := &file_file_file_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeInfo.ProtoReflect.Descriptor instead.
func (*NodeInfo) Descriptor() ([]byte, []int) {
	return file_file_file_proto_rawDescGZIP(), []int{0}
}

func (x *NodeInfo) GetPid() int32 {
	if x != nil {
		return x.Pid
	}
	return 0
}

func (x *NodeInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *NodeInfo) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *NodeInfo) GetNumOfLoads() int32 {
	if x != nil {
		return x.NumOfLoads
	}
	return 0
}

func (x *NodeInfo) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *NodeInfo) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *NodeInfo) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *NodeInfo) GetRpc() *NodeInfo_Rpc {
	if x != nil {
		return x.Rpc
	}
	return nil
}

type RouterReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Md5 string `protobuf:"bytes,1,opt,name=md5,proto3" json:"md5,omitempty"`
}

func (x *RouterReq) Reset() {
	*x = RouterReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_file_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouterReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouterReq) ProtoMessage() {}

func (x *RouterReq) ProtoReflect() protoreflect.Message {
	mi := &file_file_file_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouterReq.ProtoReflect.Descriptor instead.
func (*RouterReq) Descriptor() ([]byte, []int) {
	return file_file_file_proto_rawDescGZIP(), []int{1}
}

func (x *RouterReq) GetMd5() string {
	if x != nil {
		return x.Md5
	}
	return ""
}

type RouterResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Md5     string    `protobuf:"bytes,1,opt,name=md5,proto3" json:"md5,omitempty"`
	Node    *NodeInfo `protobuf:"bytes,2,opt,name=node,proto3" json:"node,omitempty"`
	ErrCode int32     `protobuf:"varint,3,opt,name=errCode,proto3" json:"errCode,omitempty"`
	ErrMsg  string    `protobuf:"bytes,4,opt,name=errMsg,proto3" json:"errMsg,omitempty"`
}

func (x *RouterResp) Reset() {
	*x = RouterResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_file_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RouterResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RouterResp) ProtoMessage() {}

func (x *RouterResp) ProtoReflect() protoreflect.Message {
	mi := &file_file_file_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RouterResp.ProtoReflect.Descriptor instead.
func (*RouterResp) Descriptor() ([]byte, []int) {
	return file_file_file_proto_rawDescGZIP(), []int{2}
}

func (x *RouterResp) GetMd5() string {
	if x != nil {
		return x.Md5
	}
	return ""
}

func (x *RouterResp) GetNode() *NodeInfo {
	if x != nil {
		return x.Node
	}
	return nil
}

func (x *RouterResp) GetErrCode() int32 {
	if x != nil {
		return x.ErrCode
	}
	return 0
}

func (x *RouterResp) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

type NodeInfo_Rpc struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ip   string `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	Port int32  `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *NodeInfo_Rpc) Reset() {
	*x = NodeInfo_Rpc{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_file_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeInfo_Rpc) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeInfo_Rpc) ProtoMessage() {}

func (x *NodeInfo_Rpc) ProtoReflect() protoreflect.Message {
	mi := &file_file_file_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeInfo_Rpc.ProtoReflect.Descriptor instead.
func (*NodeInfo_Rpc) Descriptor() ([]byte, []int) {
	return file_file_file_proto_rawDescGZIP(), []int{0, 0}
}

func (x *NodeInfo_Rpc) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *NodeInfo_Rpc) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

var File_file_file_proto protoreflect.FileDescriptor

var file_file_file_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x66, 0x69, 0x6c, 0x65, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x22, 0xed, 0x01, 0x0a, 0x08, 0x4e, 0x6f, 0x64, 0x65,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x03, 0x70, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x4e, 0x75,
	0x6d, 0x4f, 0x66, 0x4c, 0x6f, 0x61, 0x64, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a,
	0x4e, 0x75, 0x6d, 0x4f, 0x66, 0x4c, 0x6f, 0x61, 0x64, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f,
	0x72, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x24, 0x0a, 0x03, 0x72, 0x70, 0x63, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x49,
	0x6e, 0x66, 0x6f, 0x2e, 0x52, 0x70, 0x63, 0x52, 0x03, 0x72, 0x70, 0x63, 0x1a, 0x29, 0x0a, 0x03,
	0x52, 0x70, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x1d, 0x0a, 0x09, 0x52, 0x6f, 0x75, 0x74, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x64, 0x35, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6d, 0x64, 0x35, 0x22, 0x74, 0x0a, 0x0a, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x64, 0x35, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6d, 0x64, 0x35, 0x12, 0x22, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x4e, 0x6f, 0x64, 0x65,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x72,
	0x72, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x65, 0x72, 0x72,
	0x43, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x72, 0x72, 0x4d, 0x73, 0x67, 0x32, 0x36, 0x0a, 0x04,
	0x46, 0x69, 0x6c, 0x65, 0x12, 0x2e, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x52, 0x6f, 0x75, 0x74, 0x65,
	0x72, 0x12, 0x0f, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x1a, 0x10, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x42, 0x1d, 0x5a, 0x1b, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x3b, 0x70, 0x62, 0x5f, 0x66,
	0x69, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_file_file_proto_rawDescOnce sync.Once
	file_file_file_proto_rawDescData = file_file_file_proto_rawDesc
)

func file_file_file_proto_rawDescGZIP() []byte {
	file_file_file_proto_rawDescOnce.Do(func() {
		file_file_file_proto_rawDescData = protoimpl.X.CompressGZIP(file_file_file_proto_rawDescData)
	})
	return file_file_file_proto_rawDescData
}

var file_file_file_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_file_file_proto_goTypes = []interface{}{
	(*NodeInfo)(nil),     // 0: file.NodeInfo
	(*RouterReq)(nil),    // 1: file.RouterReq
	(*RouterResp)(nil),   // 2: file.RouterResp
	(*NodeInfo_Rpc)(nil), // 3: file.NodeInfo.Rpc
}
var file_file_file_proto_depIdxs = []int32{
	3, // 0: file.NodeInfo.rpc:type_name -> file.NodeInfo.Rpc
	0, // 1: file.RouterResp.node:type_name -> file.NodeInfo
	1, // 2: file.File.GetRouter:input_type -> file.RouterReq
	2, // 3: file.File.GetRouter:output_type -> file.RouterResp
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_file_file_proto_init() }
func file_file_file_proto_init() {
	if File_file_file_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_file_file_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_file_file_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouterReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_file_file_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RouterResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_file_file_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeInfo_Rpc); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_file_file_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_file_file_proto_goTypes,
		DependencyIndexes: file_file_file_proto_depIdxs,
		MessageInfos:      file_file_file_proto_msgTypes,
	}.Build()
	File_file_file_proto = out.File
	file_file_file_proto_rawDesc = nil
	file_file_file_proto_goTypes = nil
	file_file_file_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// FileClient is the client API for File service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type FileClient interface {
	GetRouter(ctx context.Context, in *RouterReq, opts ...grpc.CallOption) (*RouterResp, error)
}

type fileClient struct {
	cc grpc.ClientConnInterface
}

func NewFileClient(cc grpc.ClientConnInterface) FileClient {
	return &fileClient{cc}
}

func (c *fileClient) GetRouter(ctx context.Context, in *RouterReq, opts ...grpc.CallOption) (*RouterResp, error) {
	out := new(RouterResp)
	err := c.cc.Invoke(ctx, "/file.File/GetRouter", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileServer is the server API for File service.
type FileServer interface {
	GetRouter(context.Context, *RouterReq) (*RouterResp, error)
}

// UnimplementedFileServer can be embedded to have forward compatible implementations.
type UnimplementedFileServer struct {
}

func (*UnimplementedFileServer) GetRouter(context.Context, *RouterReq) (*RouterResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRouter not implemented")
}

func RegisterFileServer(s *grpc.Server, srv FileServer) {
	s.RegisterService(&_File_serviceDesc, srv)
}

func _File_GetRouter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RouterReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).GetRouter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.File/GetRouter",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).GetRouter(ctx, req.(*RouterReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _File_serviceDesc = grpc.ServiceDesc{
	ServiceName: "file.File",
	HandlerType: (*FileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRouter",
			Handler:    _File_GetRouter_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "file/file.proto",
}
