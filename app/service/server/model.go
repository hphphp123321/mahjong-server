package server

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
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
	RoomInfo *room.Info
}

type JoinRoomRequest struct {
	RoomID string
}

type JoinRoomReply struct {
	RoomInfo   *room.Info
	Seat       int
	PlayerName string
}

type ListRoomsRequest struct {
	RoomNameFilter string
}

type ListRoomsReply struct {
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

type ListRobotsReply struct {
	RobotTypes []string
}

type ListPlayerIDsReply struct {
	PlayerIDs []string
}

type StartGameRequest struct {
	Rule *mahjong.Rule
	Seed int64
}

type StartGameReply struct {
	SeatsOrder []int
}

type StepRequest struct {
	MapIDAction map[string]*mahjong.Call
}

type StepReply struct {
	Events       map[string]mahjong.Events
	ValidActions map[string]mahjong.Calls
}
