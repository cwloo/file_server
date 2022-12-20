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

type FileServerReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Md5 string `protobuf:"bytes,1,opt,name=md5,proto3" json:"md5,omitempty"`
}

func (x *FileServerReq) Reset() {
	*x = FileServerReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_file_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileServerReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileServerReq) ProtoMessage() {}

func (x *FileServerReq) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use FileServerReq.ProtoReflect.Descriptor instead.
func (*FileServerReq) Descriptor() ([]byte, []int) {
	return file_file_file_proto_rawDescGZIP(), []int{0}
}

func (x *FileServerReq) GetMd5() string {
	if x != nil {
		return x.Md5
	}
	return ""
}

type FileServerResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Md5        string `protobuf:"bytes,1,opt,name=md5,proto3" json:"md5,omitempty"`
	ServerAddr string `protobuf:"bytes,2,opt,name=serverAddr,proto3" json:"serverAddr,omitempty"`
}

func (x *FileServerResp) Reset() {
	*x = FileServerResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_file_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileServerResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileServerResp) ProtoMessage() {}

func (x *FileServerResp) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use FileServerResp.ProtoReflect.Descriptor instead.
func (*FileServerResp) Descriptor() ([]byte, []int) {
	return file_file_file_proto_rawDescGZIP(), []int{1}
}

func (x *FileServerResp) GetMd5() string {
	if x != nil {
		return x.Md5
	}
	return ""
}

func (x *FileServerResp) GetServerAddr() string {
	if x != nil {
		return x.ServerAddr
	}
	return ""
}

var File_file_file_proto protoreflect.FileDescriptor

var file_file_file_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x66, 0x69, 0x6c, 0x65, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x22, 0x21, 0x0a, 0x0d, 0x46, 0x69, 0x6c, 0x65, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x64, 0x35, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x64, 0x35, 0x22, 0x42, 0x0a, 0x0e, 0x46, 0x69,
	0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x12, 0x10, 0x0a, 0x03,
	0x6d, 0x64, 0x35, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x64, 0x35, 0x12, 0x1e,
	0x0a, 0x0a, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x32, 0x42,
	0x0a, 0x04, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x3a, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c,
	0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x13, 0x2e, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x71, 0x1a, 0x14, 0x2e, 0x66,
	0x69, 0x6c, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x42, 0x1d, 0x5a, 0x1b, 0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x72, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x3b, 0x70, 0x62, 0x5f, 0x66, 0x69, 0x6c,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_file_file_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_file_file_proto_goTypes = []interface{}{
	(*FileServerReq)(nil),  // 0: file.FileServerReq
	(*FileServerResp)(nil), // 1: file.FileServerResp
}
var file_file_file_proto_depIdxs = []int32{
	0, // 0: file.File.GetFileServer:input_type -> file.FileServerReq
	1, // 1: file.File.GetFileServer:output_type -> file.FileServerResp
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_file_file_proto_init() }
func file_file_file_proto_init() {
	if File_file_file_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_file_file_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileServerReq); i {
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
			switch v := v.(*FileServerResp); i {
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
			NumMessages:   2,
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
	GetFileServer(ctx context.Context, in *FileServerReq, opts ...grpc.CallOption) (*FileServerResp, error)
}

type fileClient struct {
	cc grpc.ClientConnInterface
}

func NewFileClient(cc grpc.ClientConnInterface) FileClient {
	return &fileClient{cc}
}

func (c *fileClient) GetFileServer(ctx context.Context, in *FileServerReq, opts ...grpc.CallOption) (*FileServerResp, error) {
	out := new(FileServerResp)
	err := c.cc.Invoke(ctx, "/file.File/GetFileServer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileServer is the server API for File service.
type FileServer interface {
	GetFileServer(context.Context, *FileServerReq) (*FileServerResp, error)
}

// UnimplementedFileServer can be embedded to have forward compatible implementations.
type UnimplementedFileServer struct {
}

func (*UnimplementedFileServer) GetFileServer(context.Context, *FileServerReq) (*FileServerResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileServer not implemented")
}

func RegisterFileServer(s *grpc.Server, srv FileServer) {
	s.RegisterService(&_File_serviceDesc, srv)
}

func _File_GetFileServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FileServerReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServer).GetFileServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/file.File/GetFileServer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServer).GetFileServer(ctx, req.(*FileServerReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _File_serviceDesc = grpc.ServiceDesc{
	ServiceName: "file.File",
	HandlerType: (*FileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFileServer",
			Handler:    _File_GetFileServer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "file/file.proto",
}
