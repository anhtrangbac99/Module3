// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        (unknown)
// source: Proto/SignIn/protobuf.proto

package protobuf

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type AuthorRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=Username,json=username,proto3" json:"Username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=Password,json=password,proto3" json:"Password,omitempty"`
}

func (x *AuthorRequest) Reset() {
	*x = AuthorRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Proto_SignIn_protobuf_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthorRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthorRequest) ProtoMessage() {}

func (x *AuthorRequest) ProtoReflect() protoreflect.Message {
	mi := &file_Proto_SignIn_protobuf_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthorRequest.ProtoReflect.Descriptor instead.
func (*AuthorRequest) Descriptor() ([]byte, []int) {
	return file_Proto_SignIn_protobuf_proto_rawDescGZIP(), []int{0}
}

func (x *AuthorRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *AuthorRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type AuthorRespone struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsExisted  int64 `protobuf:"varint,1,opt,name=IsExisted,json=isExisted,proto3" json:"IsExisted,omitempty"`
	User_Id    int64 `protobuf:"varint,2,opt,name=User_Id,json=userId,proto3" json:"User_Id,omitempty"`
	Authorized int64 `protobuf:"varint,3,opt,name=Authorized,json=authorized,proto3" json:"Authorized,omitempty"`
}

func (x *AuthorRespone) Reset() {
	*x = AuthorRespone{}
	if protoimpl.UnsafeEnabled {
		mi := &file_Proto_SignIn_protobuf_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthorRespone) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthorRespone) ProtoMessage() {}

func (x *AuthorRespone) ProtoReflect() protoreflect.Message {
	mi := &file_Proto_SignIn_protobuf_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthorRespone.ProtoReflect.Descriptor instead.
func (*AuthorRespone) Descriptor() ([]byte, []int) {
	return file_Proto_SignIn_protobuf_proto_rawDescGZIP(), []int{1}
}

func (x *AuthorRespone) GetIsExisted() int64 {
	if x != nil {
		return x.IsExisted
	}
	return 0
}

func (x *AuthorRespone) GetUser_Id() int64 {
	if x != nil {
		return x.User_Id
	}
	return 0
}

func (x *AuthorRespone) GetAuthorized() int64 {
	if x != nil {
		return x.Authorized
	}
	return 0
}

var File_Proto_SignIn_protobuf_proto protoreflect.FileDescriptor

var file_Proto_SignIn_protobuf_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x53, 0x69, 0x67, 0x6e, 0x49, 0x6e, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0x47, 0x0a, 0x0d, 0x41, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x22, 0x66, 0x0a, 0x0d, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x65, 0x12, 0x1c, 0x0a, 0x09, 0x49, 0x73, 0x45, 0x78, 0x69, 0x73, 0x74, 0x65, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x69, 0x73, 0x45, 0x78, 0x69, 0x73, 0x74, 0x65, 0x64, 0x12,
	0x17, 0x0a, 0x07, 0x55, 0x73, 0x65, 0x72, 0x5f, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x41, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x69, 0x7a, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x61, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x64, 0x32, 0x4a, 0x0a, 0x06, 0x53, 0x69, 0x67, 0x6e,
	0x49, 0x6e, 0x12, 0x40, 0x0a, 0x0a, 0x55, 0x73, 0x65, 0x72, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x65, 0x22, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_Proto_SignIn_protobuf_proto_rawDescOnce sync.Once
	file_Proto_SignIn_protobuf_proto_rawDescData = file_Proto_SignIn_protobuf_proto_rawDesc
)

func file_Proto_SignIn_protobuf_proto_rawDescGZIP() []byte {
	file_Proto_SignIn_protobuf_proto_rawDescOnce.Do(func() {
		file_Proto_SignIn_protobuf_proto_rawDescData = protoimpl.X.CompressGZIP(file_Proto_SignIn_protobuf_proto_rawDescData)
	})
	return file_Proto_SignIn_protobuf_proto_rawDescData
}

var file_Proto_SignIn_protobuf_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_Proto_SignIn_protobuf_proto_goTypes = []interface{}{
	(*AuthorRequest)(nil), // 0: protobuf.AuthorRequest
	(*AuthorRespone)(nil), // 1: protobuf.AuthorRespone
}
var file_Proto_SignIn_protobuf_proto_depIdxs = []int32{
	0, // 0: protobuf.SignIn.UserAuthor:input_type -> protobuf.AuthorRequest
	1, // 1: protobuf.SignIn.UserAuthor:output_type -> protobuf.AuthorRespone
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_Proto_SignIn_protobuf_proto_init() }
func file_Proto_SignIn_protobuf_proto_init() {
	if File_Proto_SignIn_protobuf_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_Proto_SignIn_protobuf_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthorRequest); i {
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
		file_Proto_SignIn_protobuf_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthorRespone); i {
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
			RawDescriptor: file_Proto_SignIn_protobuf_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_Proto_SignIn_protobuf_proto_goTypes,
		DependencyIndexes: file_Proto_SignIn_protobuf_proto_depIdxs,
		MessageInfos:      file_Proto_SignIn_protobuf_proto_msgTypes,
	}.Build()
	File_Proto_SignIn_protobuf_proto = out.File
	file_Proto_SignIn_protobuf_proto_rawDesc = nil
	file_Proto_SignIn_protobuf_proto_goTypes = nil
	file_Proto_SignIn_protobuf_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// SignInClient is the client API for SignIn service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SignInClient interface {
	UserAuthor(ctx context.Context, in *AuthorRequest, opts ...grpc.CallOption) (*AuthorRespone, error)
}

type signInClient struct {
	cc grpc.ClientConnInterface
}

func NewSignInClient(cc grpc.ClientConnInterface) SignInClient {
	return &signInClient{cc}
}

func (c *signInClient) UserAuthor(ctx context.Context, in *AuthorRequest, opts ...grpc.CallOption) (*AuthorRespone, error) {
	out := new(AuthorRespone)
	err := c.cc.Invoke(ctx, "/protobuf.SignIn/UserAuthor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SignInServer is the server API for SignIn service.
type SignInServer interface {
	UserAuthor(context.Context, *AuthorRequest) (*AuthorRespone, error)
}

// UnimplementedSignInServer can be embedded to have forward compatible implementations.
type UnimplementedSignInServer struct {
}

func (*UnimplementedSignInServer) UserAuthor(context.Context, *AuthorRequest) (*AuthorRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserAuthor not implemented")
}

func RegisterSignInServer(s *grpc.Server, srv SignInServer) {
	s.RegisterService(&_SignIn_serviceDesc, srv)
}

func _SignIn_UserAuthor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SignInServer).UserAuthor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.SignIn/UserAuthor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SignInServer).UserAuthor(ctx, req.(*AuthorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SignIn_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.SignIn",
	HandlerType: (*SignInServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserAuthor",
			Handler:    _SignIn_UserAuthor_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Proto/SignIn/protobuf.proto",
}
