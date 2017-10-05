// Code generated by protoc-gen-go. DO NOT EDIT.
// source: user.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	user.proto

It has these top-level messages:
	User
	UserReq
	UserResp
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type User struct {
	Id   int64  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto1.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *User) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *User) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type UserReq struct {
	Id int64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
}

func (m *UserReq) Reset()                    { *m = UserReq{} }
func (m *UserReq) String() string            { return proto1.CompactTextString(m) }
func (*UserReq) ProtoMessage()               {}
func (*UserReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *UserReq) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type UserResp struct {
	User *User `protobuf:"bytes,1,opt,name=user" json:"user,omitempty"`
}

func (m *UserResp) Reset()                    { *m = UserResp{} }
func (m *UserResp) String() string            { return proto1.CompactTextString(m) }
func (*UserResp) ProtoMessage()               {}
func (*UserResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *UserResp) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func init() {
	proto1.RegisterType((*User)(nil), "proto.User")
	proto1.RegisterType((*UserReq)(nil), "proto.UserReq")
	proto1.RegisterType((*UserResp)(nil), "proto.UserResp")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for SdUser service

type SdUserClient interface {
	GetUserInfo(ctx context.Context, in *UserReq, opts ...grpc.CallOption) (*UserResp, error)
}

type sdUserClient struct {
	cc *grpc.ClientConn
}

func NewSdUserClient(cc *grpc.ClientConn) SdUserClient {
	return &sdUserClient{cc}
}

func (c *sdUserClient) GetUserInfo(ctx context.Context, in *UserReq, opts ...grpc.CallOption) (*UserResp, error) {
	out := new(UserResp)
	err := grpc.Invoke(ctx, "/proto.SdUser/GetUserInfo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for SdUser service

type SdUserServer interface {
	GetUserInfo(context.Context, *UserReq) (*UserResp, error)
}

func RegisterSdUserServer(s *grpc.Server, srv SdUserServer) {
	s.RegisterService(&_SdUser_serviceDesc, srv)
}

func _SdUser_GetUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SdUserServer).GetUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.SdUser/GetUserInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SdUserServer).GetUserInfo(ctx, req.(*UserReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _SdUser_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.SdUser",
	HandlerType: (*SdUserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserInfo",
			Handler:    _SdUser_GetUserInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}

func init() { proto1.RegisterFile("user.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 156 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x2d, 0x4e, 0x2d,
	0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x4a, 0x5a, 0x5c, 0x2c, 0xa1, 0xc5,
	0xa9, 0x45, 0x42, 0x7c, 0x5c, 0x4c, 0x99, 0x29, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xcc, 0x41, 0x4c,
	0x99, 0x29, 0x42, 0x42, 0x5c, 0x2c, 0x79, 0x89, 0xb9, 0xa9, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x9c,
	0x41, 0x60, 0xb6, 0x92, 0x24, 0x17, 0x3b, 0x48, 0x6d, 0x50, 0x6a, 0x21, 0xba, 0x72, 0x25, 0x6d,
	0x2e, 0x0e, 0x88, 0x54, 0x71, 0x81, 0x90, 0x3c, 0x17, 0x0b, 0xc8, 0x1e, 0xb0, 0x2c, 0xb7, 0x11,
	0x37, 0xc4, 0x3e, 0x3d, 0xb0, 0x34, 0x58, 0xc2, 0xc8, 0x8a, 0x8b, 0x2d, 0x38, 0x05, 0x6c, 0xab,
	0x01, 0x17, 0xb7, 0x7b, 0x6a, 0x09, 0x88, 0xe9, 0x99, 0x97, 0x96, 0x2f, 0xc4, 0x87, 0xac, 0x36,
	0xb5, 0x50, 0x8a, 0x1f, 0x85, 0x5f, 0x5c, 0xa0, 0xc4, 0x90, 0xc4, 0x06, 0x16, 0x31, 0x06, 0x04,
	0x00, 0x00, 0xff, 0xff, 0xae, 0x9b, 0x18, 0x71, 0xcb, 0x00, 0x00, 0x00,
}
