package server

import (
	"context"
)

type Server interface {
	GetID(ctx context.Context) (string, error)
	GetName(ctx context.Context) (string, error)

	Login(ctx context.Context, request *LoginRequest) (reply *LoginReply, err error)
	Logout(ctx context.Context) error
	CreateRoom(ctx context.Context, request *CreateRoomRequest) (reply *CreateRoomReply, err error)
	JoinRoom(ctx context.Context, request *JoinRoomRequest) (joinReply *JoinRoomReply, err error)
	ListRooms(ctx context.Context, request *ListRoomsRequest) (reply *ListRoomsReply, err error)
	GetReady(ctx context.Context) (reply *GetReadyReply, err error)
	CancelReady(ctx context.Context) (reply *CancelReadyReply, err error)
	LeaveRoom(ctx context.Context) (reply *PlayerLeaveReply, err error)
	AddRobot(ctx context.Context, request *AddRobotRequest) (reply *AddRobotReply, err error)
	RemovePlayer(ctx context.Context, request *RemovePlayerRequest) (reply *PlayerLeaveReply, err error)
	ListRobots(ctx context.Context) (reply *ListRobotsReply, err error)

	ListPlayerIDs(ctx context.Context) (reply *ListPlayerIDsReply, err error) // List all player IDs in the room
}