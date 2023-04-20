package server

import (
	"github.com/hphphp123321/mahjong-server/app/model/room"
)

type LoginRequest struct {
	Name string
}

type LoginReply struct {
	ID string
}

type CreateRoomRequest struct {
	RoomName string
}

type CreateRoomReply struct {
	RoomID string
}

type JoinRoomRequest struct {
	RoomID string
}

type JoinRoomReply struct {
	RoomName   string
	RoomInfo   *room.Info
	Seat       int
	PlayerName string
}

type RefreshRoomRequest struct {
	RoomNameFilter string
}

type RefreshRoomReply struct {
	RoomInfos []*room.Info
}

type AddRobotRequest struct {
	RobotSeat int
	RobotType string
}

type AddRobotReply struct {
	RobotSeat int
	RobotType string
}

type RemovePlayerRequest struct {
	Seat int
}

type GetReadyReply struct {
	Seat       int
	PlayerName string
}

type CancelReadyReply struct {
	Seat       int
	PlayerName string
}

type PlayerLeaveReply struct {
	Seat       int
	OwnerSeat  int
	PlayerName string
}

type ListRobotReply struct {
	RobotTypes []string
}
