package server

import "context"

type Server interface {
	Login(ctx context.Context, request *LoginRequest) (reply *LoginReply, err error)
	Logout(ctx context.Context) error
	CreateRoom(ctx context.Context, request *CreateRoomRequest) (reply *CreateRoomReply, err error)
	JoinRoom(ctx context.Context, request *JoinRoomRequest) (joinReply *JoinRoomReply, err error)
	RefreshRoom(ctx context.Context, request *RefreshRoomRequest) (reply *RefreshRoomReply, err error)
	GetReady(ctx context.Context) (reply *GetReadyReply, err error)
	CancelReady(ctx context.Context) (reply *CancelReadyReply, err error)
	LeaveRoom(ctx context.Context) (reply *PlayerLeaveReply, err error)
	AddRobot(ctx context.Context) (reply *AddRobotReply, err error)
	RemovePlayer(ctx context.Context) (reply *PlayerLeaveReply, err error)
	ListRobot(ctx context.Context) (reply *ListRobotReply, err error)
}
