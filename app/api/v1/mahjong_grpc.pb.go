// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: app/api/v1/mahjong.proto

package api_mahjong_v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MahjongClient is the client API for Mahjong service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MahjongClient interface {
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error)
	Logout(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*LogoutReply, error)
	CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*CreateRoomReply, error)
	JoinRoom(ctx context.Context, in *JoinRoomRequest, opts ...grpc.CallOption) (*JoinRoomReply, error)
	ListRooms(ctx context.Context, in *ListRoomsRequest, opts ...grpc.CallOption) (*ListRoomsReply, error)
	Ready(ctx context.Context, opts ...grpc.CallOption) (Mahjong_ReadyClient, error)
	Game(ctx context.Context, opts ...grpc.CallOption) (Mahjong_GameClient, error)
}

type mahjongClient struct {
	cc grpc.ClientConnInterface
}

func NewMahjongClient(cc grpc.ClientConnInterface) MahjongClient {
	return &mahjongClient{cc}
}

func (c *mahjongClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mahjong.Mahjong/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mahjongClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error) {
	out := new(LoginReply)
	err := c.cc.Invoke(ctx, "/mahjong.Mahjong/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mahjongClient) Logout(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*LogoutReply, error) {
	out := new(LogoutReply)
	err := c.cc.Invoke(ctx, "/mahjong.Mahjong/Logout", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mahjongClient) CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*CreateRoomReply, error) {
	out := new(CreateRoomReply)
	err := c.cc.Invoke(ctx, "/mahjong.Mahjong/CreateRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mahjongClient) JoinRoom(ctx context.Context, in *JoinRoomRequest, opts ...grpc.CallOption) (*JoinRoomReply, error) {
	out := new(JoinRoomReply)
	err := c.cc.Invoke(ctx, "/mahjong.Mahjong/JoinRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mahjongClient) ListRooms(ctx context.Context, in *ListRoomsRequest, opts ...grpc.CallOption) (*ListRoomsReply, error) {
	out := new(ListRoomsReply)
	err := c.cc.Invoke(ctx, "/mahjong.Mahjong/ListRooms", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mahjongClient) Ready(ctx context.Context, opts ...grpc.CallOption) (Mahjong_ReadyClient, error) {
	stream, err := c.cc.NewStream(ctx, &Mahjong_ServiceDesc.Streams[0], "/mahjong.Mahjong/Ready", opts...)
	if err != nil {
		return nil, err
	}
	x := &mahjongReadyClient{stream}
	return x, nil
}

type Mahjong_ReadyClient interface {
	Send(*ReadyRequest) error
	Recv() (*ReadyReply, error)
	grpc.ClientStream
}

type mahjongReadyClient struct {
	grpc.ClientStream
}

func (x *mahjongReadyClient) Send(m *ReadyRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *mahjongReadyClient) Recv() (*ReadyReply, error) {
	m := new(ReadyReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mahjongClient) Game(ctx context.Context, opts ...grpc.CallOption) (Mahjong_GameClient, error) {
	stream, err := c.cc.NewStream(ctx, &Mahjong_ServiceDesc.Streams[1], "/mahjong.Mahjong/Game", opts...)
	if err != nil {
		return nil, err
	}
	x := &mahjongGameClient{stream}
	return x, nil
}

type Mahjong_GameClient interface {
	Send(*GameRequest) error
	Recv() (*GameReply, error)
	grpc.ClientStream
}

type mahjongGameClient struct {
	grpc.ClientStream
}

func (x *mahjongGameClient) Send(m *GameRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *mahjongGameClient) Recv() (*GameReply, error) {
	m := new(GameReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MahjongServer is the server API for Mahjong service.
// All implementations must embed UnimplementedMahjongServer
// for forward compatibility
type MahjongServer interface {
	Ping(context.Context, *Empty) (*Empty, error)
	Login(context.Context, *LoginRequest) (*LoginReply, error)
	Logout(context.Context, *Empty) (*LogoutReply, error)
	CreateRoom(context.Context, *CreateRoomRequest) (*CreateRoomReply, error)
	JoinRoom(context.Context, *JoinRoomRequest) (*JoinRoomReply, error)
	ListRooms(context.Context, *ListRoomsRequest) (*ListRoomsReply, error)
	Ready(Mahjong_ReadyServer) error
	Game(Mahjong_GameServer) error
	mustEmbedUnimplementedMahjongServer()
}

// UnimplementedMahjongServer must be embedded to have forward compatible implementations.
type UnimplementedMahjongServer struct {
}

func (UnimplementedMahjongServer) Ping(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedMahjongServer) Login(context.Context, *LoginRequest) (*LoginReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedMahjongServer) Logout(context.Context, *Empty) (*LogoutReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedMahjongServer) CreateRoom(context.Context, *CreateRoomRequest) (*CreateRoomReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRoom not implemented")
}
func (UnimplementedMahjongServer) JoinRoom(context.Context, *JoinRoomRequest) (*JoinRoomReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinRoom not implemented")
}
func (UnimplementedMahjongServer) ListRooms(context.Context, *ListRoomsRequest) (*ListRoomsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRooms not implemented")
}
func (UnimplementedMahjongServer) Ready(Mahjong_ReadyServer) error {
	return status.Errorf(codes.Unimplemented, "method Ready not implemented")
}
func (UnimplementedMahjongServer) Game(Mahjong_GameServer) error {
	return status.Errorf(codes.Unimplemented, "method Game not implemented")
}
func (UnimplementedMahjongServer) mustEmbedUnimplementedMahjongServer() {}

// UnsafeMahjongServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MahjongServer will
// result in compilation errors.
type UnsafeMahjongServer interface {
	mustEmbedUnimplementedMahjongServer()
}

func RegisterMahjongServer(s grpc.ServiceRegistrar, srv MahjongServer) {
	s.RegisterService(&Mahjong_ServiceDesc, srv)
}

func _Mahjong_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MahjongServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mahjong.Mahjong/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MahjongServer).Ping(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mahjong_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MahjongServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mahjong.Mahjong/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MahjongServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mahjong_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MahjongServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mahjong.Mahjong/Logout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MahjongServer).Logout(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mahjong_CreateRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MahjongServer).CreateRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mahjong.Mahjong/CreateRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MahjongServer).CreateRoom(ctx, req.(*CreateRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mahjong_JoinRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MahjongServer).JoinRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mahjong.Mahjong/JoinRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MahjongServer).JoinRoom(ctx, req.(*JoinRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mahjong_ListRooms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRoomsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MahjongServer).ListRooms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mahjong.Mahjong/ListRooms",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MahjongServer).ListRooms(ctx, req.(*ListRoomsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mahjong_Ready_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MahjongServer).Ready(&mahjongReadyServer{stream})
}

type Mahjong_ReadyServer interface {
	Send(*ReadyReply) error
	Recv() (*ReadyRequest, error)
	grpc.ServerStream
}

type mahjongReadyServer struct {
	grpc.ServerStream
}

func (x *mahjongReadyServer) Send(m *ReadyReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *mahjongReadyServer) Recv() (*ReadyRequest, error) {
	m := new(ReadyRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Mahjong_Game_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MahjongServer).Game(&mahjongGameServer{stream})
}

type Mahjong_GameServer interface {
	Send(*GameReply) error
	Recv() (*GameRequest, error)
	grpc.ServerStream
}

type mahjongGameServer struct {
	grpc.ServerStream
}

func (x *mahjongGameServer) Send(m *GameReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *mahjongGameServer) Recv() (*GameRequest, error) {
	m := new(GameRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Mahjong_ServiceDesc is the grpc.ServiceDesc for Mahjong service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Mahjong_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mahjong.Mahjong",
	HandlerType: (*MahjongServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Mahjong_Ping_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Mahjong_Login_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _Mahjong_Logout_Handler,
		},
		{
			MethodName: "CreateRoom",
			Handler:    _Mahjong_CreateRoom_Handler,
		},
		{
			MethodName: "JoinRoom",
			Handler:    _Mahjong_JoinRoom_Handler,
		},
		{
			MethodName: "ListRooms",
			Handler:    _Mahjong_ListRooms_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Ready",
			Handler:       _Mahjong_Ready_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Game",
			Handler:       _Mahjong_Game_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "app/api/v1/mahjong.proto",
}
