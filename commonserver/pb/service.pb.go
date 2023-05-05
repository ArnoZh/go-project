// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: proto/service.proto

package pb

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
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

//数据流
type StreamData struct {
	//msg数据
	Msg []byte `protobuf:"bytes,1,opt,name=Msg,proto3" json:"Msg,omitempty"`
	//生成的时间 毫秒
	GenTs                int64    `protobuf:"varint,2,opt,name=GenTs,proto3" json:"GenTs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StreamData) Reset()         { *m = StreamData{} }
func (m *StreamData) String() string { return proto.CompactTextString(m) }
func (*StreamData) ProtoMessage()    {}
func (*StreamData) Descriptor() ([]byte, []int) {
	return fileDescriptor_c33392ef2c1961ba, []int{0}
}
func (m *StreamData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *StreamData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_StreamData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *StreamData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamData.Merge(m, src)
}
func (m *StreamData) XXX_Size() int {
	return m.Size()
}
func (m *StreamData) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamData.DiscardUnknown(m)
}

var xxx_messageInfo_StreamData proto.InternalMessageInfo

func (m *StreamData) GetMsg() []byte {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (m *StreamData) GetGenTs() int64 {
	if m != nil {
		return m.GenTs
	}
	return 0
}

type UpdateLoginInfoReq struct {
	// 账号id
	AccountId int64 `protobuf:"varint,1,opt,name=AccountId,proto3" json:"AccountId,omitempty"`
	// 角色id
	RoleId int64 `protobuf:"varint,2,opt,name=RoleId,proto3" json:"RoleId,omitempty"`
	// 游戏服id
	GameId               int32    `protobuf:"varint,3,opt,name=GameId,proto3" json:"GameId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateLoginInfoReq) Reset()         { *m = UpdateLoginInfoReq{} }
func (m *UpdateLoginInfoReq) String() string { return proto.CompactTextString(m) }
func (*UpdateLoginInfoReq) ProtoMessage()    {}
func (*UpdateLoginInfoReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_c33392ef2c1961ba, []int{1}
}
func (m *UpdateLoginInfoReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UpdateLoginInfoReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UpdateLoginInfoReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UpdateLoginInfoReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateLoginInfoReq.Merge(m, src)
}
func (m *UpdateLoginInfoReq) XXX_Size() int {
	return m.Size()
}
func (m *UpdateLoginInfoReq) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateLoginInfoReq.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateLoginInfoReq proto.InternalMessageInfo

func (m *UpdateLoginInfoReq) GetAccountId() int64 {
	if m != nil {
		return m.AccountId
	}
	return 0
}

func (m *UpdateLoginInfoReq) GetRoleId() int64 {
	if m != nil {
		return m.RoleId
	}
	return 0
}

func (m *UpdateLoginInfoReq) GetGameId() int32 {
	if m != nil {
		return m.GameId
	}
	return 0
}

type UpdateLoginInfoRes struct {
	ErrorCode            int32    `protobuf:"varint,1,opt,name=ErrorCode,proto3" json:"ErrorCode,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateLoginInfoRes) Reset()         { *m = UpdateLoginInfoRes{} }
func (m *UpdateLoginInfoRes) String() string { return proto.CompactTextString(m) }
func (*UpdateLoginInfoRes) ProtoMessage()    {}
func (*UpdateLoginInfoRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_c33392ef2c1961ba, []int{2}
}
func (m *UpdateLoginInfoRes) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UpdateLoginInfoRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UpdateLoginInfoRes.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UpdateLoginInfoRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateLoginInfoRes.Merge(m, src)
}
func (m *UpdateLoginInfoRes) XXX_Size() int {
	return m.Size()
}
func (m *UpdateLoginInfoRes) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateLoginInfoRes.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateLoginInfoRes proto.InternalMessageInfo

func (m *UpdateLoginInfoRes) GetErrorCode() int32 {
	if m != nil {
		return m.ErrorCode
	}
	return 0
}

func init() {
	proto.RegisterType((*StreamData)(nil), "pb.StreamData")
	proto.RegisterType((*UpdateLoginInfoReq)(nil), "pb.UpdateLoginInfoReq")
	proto.RegisterType((*UpdateLoginInfoRes)(nil), "pb.UpdateLoginInfoRes")
}

func init() { proto.RegisterFile("proto/service.proto", fileDescriptor_c33392ef2c1961ba) }

var fileDescriptor_c33392ef2c1961ba = []byte{
	// 308 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xd1, 0x4e, 0xc2, 0x30,
	0x14, 0x86, 0x29, 0x0b, 0x24, 0x9e, 0x28, 0x9a, 0x4a, 0x08, 0x21, 0x66, 0x21, 0xbb, 0xda, 0xd5,
	0x66, 0x90, 0x07, 0x10, 0xc4, 0x10, 0x12, 0xbd, 0x29, 0xf2, 0x00, 0xeb, 0x5a, 0xe6, 0x12, 0xd9,
	0x99, 0x6d, 0xe1, 0x59, 0x7c, 0x24, 0x2f, 0x79, 0x04, 0xc5, 0x17, 0x31, 0xed, 0x30, 0x24, 0xc8,
	0xdd, 0xf9, 0xfe, 0xf4, 0xfc, 0x7f, 0xff, 0x16, 0xae, 0x4b, 0x85, 0x06, 0x63, 0x2d, 0xd5, 0x26,
	0x4f, 0x65, 0xe4, 0x88, 0xd6, 0x4b, 0xde, 0x1b, 0x6e, 0x64, 0x21, 0x50, 0xc5, 0x59, 0x6e, 0x5e,
	0xd7, 0x3c, 0x4a, 0x71, 0x15, 0x67, 0x98, 0x61, 0xec, 0x4e, 0xf0, 0xf5, 0xd2, 0x51, 0xb5, 0x6c,
	0xa7, 0x6a, 0x33, 0x18, 0x02, 0xcc, 0x8d, 0x92, 0xc9, 0x6a, 0x92, 0x98, 0x84, 0x5e, 0x81, 0xf7,
	0xac, 0xb3, 0x2e, 0xe9, 0x93, 0xf0, 0x9c, 0xd9, 0x91, 0xb6, 0xa1, 0x31, 0x95, 0xc5, 0x8b, 0xee,
	0xd6, 0xfb, 0x24, 0xf4, 0x58, 0x05, 0x01, 0x07, 0xba, 0x28, 0x45, 0x62, 0xe4, 0x13, 0x66, 0x79,
	0x31, 0x2b, 0x96, 0xc8, 0xe4, 0x3b, 0xbd, 0x81, 0xb3, 0x51, 0x9a, 0xe2, 0xba, 0x30, 0x33, 0xe1,
	0x3c, 0x3c, 0x76, 0x10, 0x68, 0x07, 0x9a, 0x0c, 0xdf, 0xe4, 0x4c, 0xec, 0xad, 0xf6, 0x64, 0xf5,
	0x69, 0xb2, 0xb2, 0xba, 0xd7, 0x27, 0x61, 0x83, 0xed, 0x29, 0x18, 0x9c, 0xc8, 0xd0, 0x36, 0xe3,
	0x51, 0x29, 0x54, 0x0f, 0x28, 0xa4, 0xcb, 0x68, 0xb0, 0x83, 0x30, 0xb8, 0x87, 0x16, 0x93, 0x22,
	0xd7, 0x93, 0xf1, 0xbc, 0x7a, 0x1f, 0x1a, 0x41, 0xb3, 0xea, 0x47, 0x5b, 0x51, 0xc9, 0xa3, 0x43,
	0xd7, 0xde, 0x11, 0x07, 0xb5, 0x90, 0xdc, 0x92, 0x01, 0x83, 0x8b, 0x85, 0x96, 0x6a, 0xc2, 0xff,
	0x0c, 0x46, 0x70, 0x79, 0x74, 0x0d, 0xda, 0xb1, 0x9b, 0xff, 0xfb, 0xf7, 0x4e, 0xeb, 0x7a, 0xdc,
	0xde, 0x7e, 0xfb, 0xb5, 0xcf, 0x9d, 0x4f, 0xb6, 0x3b, 0x9f, 0x7c, 0xed, 0x7c, 0xf2, 0xf1, 0xe3,
	0xd7, 0x78, 0xd3, 0x7d, 0xc0, 0xdd, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb8, 0xfe, 0xb5, 0x5d,
	0xd1, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RedisDBServiceClient is the client API for RedisDBService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RedisDBServiceClient interface {
	Stream(ctx context.Context, opts ...grpc.CallOption) (RedisDBService_StreamClient, error)
}

type redisDBServiceClient struct {
	cc *grpc.ClientConn
}

func NewRedisDBServiceClient(cc *grpc.ClientConn) RedisDBServiceClient {
	return &redisDBServiceClient{cc}
}

func (c *redisDBServiceClient) Stream(ctx context.Context, opts ...grpc.CallOption) (RedisDBService_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &_RedisDBService_serviceDesc.Streams[0], "/pb.RedisDBService/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &redisDBServiceStreamClient{stream}
	return x, nil
}

type RedisDBService_StreamClient interface {
	Send(*StreamData) error
	Recv() (*StreamData, error)
	grpc.ClientStream
}

type redisDBServiceStreamClient struct {
	grpc.ClientStream
}

func (x *redisDBServiceStreamClient) Send(m *StreamData) error {
	return x.ClientStream.SendMsg(m)
}

func (x *redisDBServiceStreamClient) Recv() (*StreamData, error) {
	m := new(StreamData)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RedisDBServiceServer is the server API for RedisDBService service.
type RedisDBServiceServer interface {
	Stream(RedisDBService_StreamServer) error
}

// UnimplementedRedisDBServiceServer can be embedded to have forward compatible implementations.
type UnimplementedRedisDBServiceServer struct {
}

func (*UnimplementedRedisDBServiceServer) Stream(srv RedisDBService_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}

func RegisterRedisDBServiceServer(s *grpc.Server, srv RedisDBServiceServer) {
	s.RegisterService(&_RedisDBService_serviceDesc, srv)
}

func _RedisDBService_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RedisDBServiceServer).Stream(&redisDBServiceStreamServer{stream})
}

type RedisDBService_StreamServer interface {
	Send(*StreamData) error
	Recv() (*StreamData, error)
	grpc.ServerStream
}

type redisDBServiceStreamServer struct {
	grpc.ServerStream
}

func (x *redisDBServiceStreamServer) Send(m *StreamData) error {
	return x.ServerStream.SendMsg(m)
}

func (x *redisDBServiceStreamServer) Recv() (*StreamData, error) {
	m := new(StreamData)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _RedisDBService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.RedisDBService",
	HandlerType: (*RedisDBServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _RedisDBService_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/service.proto",
}

// UserDbServiceClient is the client API for UserDbService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UserDbServiceClient interface {
	UpdateLoginInfo(ctx context.Context, in *UpdateLoginInfoReq, opts ...grpc.CallOption) (*UpdateLoginInfoRes, error)
}

type userDbServiceClient struct {
	cc *grpc.ClientConn
}

func NewUserDbServiceClient(cc *grpc.ClientConn) UserDbServiceClient {
	return &userDbServiceClient{cc}
}

func (c *userDbServiceClient) UpdateLoginInfo(ctx context.Context, in *UpdateLoginInfoReq, opts ...grpc.CallOption) (*UpdateLoginInfoRes, error) {
	out := new(UpdateLoginInfoRes)
	err := c.cc.Invoke(ctx, "/pb.UserDbService/UpdateLoginInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserDbServiceServer is the server API for UserDbService service.
type UserDbServiceServer interface {
	UpdateLoginInfo(context.Context, *UpdateLoginInfoReq) (*UpdateLoginInfoRes, error)
}

// UnimplementedUserDbServiceServer can be embedded to have forward compatible implementations.
type UnimplementedUserDbServiceServer struct {
}

func (*UnimplementedUserDbServiceServer) UpdateLoginInfo(ctx context.Context, req *UpdateLoginInfoReq) (*UpdateLoginInfoRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLoginInfo not implemented")
}

func RegisterUserDbServiceServer(s *grpc.Server, srv UserDbServiceServer) {
	s.RegisterService(&_UserDbService_serviceDesc, srv)
}

func _UserDbService_UpdateLoginInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateLoginInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserDbServiceServer).UpdateLoginInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.UserDbService/UpdateLoginInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserDbServiceServer).UpdateLoginInfo(ctx, req.(*UpdateLoginInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserDbService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.UserDbService",
	HandlerType: (*UserDbServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateLoginInfo",
			Handler:    _UserDbService_UpdateLoginInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/service.proto",
}

func (m *StreamData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *StreamData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *StreamData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.GenTs != 0 {
		i = encodeVarintService(dAtA, i, uint64(m.GenTs))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Msg) > 0 {
		i -= len(m.Msg)
		copy(dAtA[i:], m.Msg)
		i = encodeVarintService(dAtA, i, uint64(len(m.Msg)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *UpdateLoginInfoReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UpdateLoginInfoReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UpdateLoginInfoReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.GameId != 0 {
		i = encodeVarintService(dAtA, i, uint64(m.GameId))
		i--
		dAtA[i] = 0x18
	}
	if m.RoleId != 0 {
		i = encodeVarintService(dAtA, i, uint64(m.RoleId))
		i--
		dAtA[i] = 0x10
	}
	if m.AccountId != 0 {
		i = encodeVarintService(dAtA, i, uint64(m.AccountId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *UpdateLoginInfoRes) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UpdateLoginInfoRes) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UpdateLoginInfoRes) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.ErrorCode != 0 {
		i = encodeVarintService(dAtA, i, uint64(m.ErrorCode))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintService(dAtA []byte, offset int, v uint64) int {
	offset -= sovService(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *StreamData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Msg)
	if l > 0 {
		n += 1 + l + sovService(uint64(l))
	}
	if m.GenTs != 0 {
		n += 1 + sovService(uint64(m.GenTs))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *UpdateLoginInfoReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.AccountId != 0 {
		n += 1 + sovService(uint64(m.AccountId))
	}
	if m.RoleId != 0 {
		n += 1 + sovService(uint64(m.RoleId))
	}
	if m.GameId != 0 {
		n += 1 + sovService(uint64(m.GameId))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *UpdateLoginInfoRes) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ErrorCode != 0 {
		n += 1 + sovService(uint64(m.ErrorCode))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovService(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozService(x uint64) (n int) {
	return sovService(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *StreamData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: StreamData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StreamData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthService
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Msg = append(m.Msg[:0], dAtA[iNdEx:postIndex]...)
			if m.Msg == nil {
				m.Msg = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GenTs", wireType)
			}
			m.GenTs = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GenTs |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *UpdateLoginInfoReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: UpdateLoginInfoReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UpdateLoginInfoReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AccountId", wireType)
			}
			m.AccountId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AccountId |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RoleId", wireType)
			}
			m.RoleId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RoleId |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GameId", wireType)
			}
			m.GameId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GameId |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *UpdateLoginInfoRes) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: UpdateLoginInfoRes: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UpdateLoginInfoRes: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ErrorCode", wireType)
			}
			m.ErrorCode = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ErrorCode |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipService(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowService
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowService
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowService
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthService
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupService
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthService
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthService        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowService          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupService = fmt.Errorf("proto: unexpected end of group")
)
