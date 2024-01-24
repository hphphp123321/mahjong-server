package server

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
	"github.com/hphphp123321/mahjong-server/app/model/player"
	"github.com/hphphp123321/mahjong-server/app/model/room"
	"github.com/hphphp123321/mahjong-server/app/service/robot/remote"
)

type LoginRequest struct {
	Name     string
	Password string
}

type LoginReply struct {
	ID string
}

type RegisterRequest struct {
	Name     string
	Password string
}

type RegisterReply struct {
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
	RobotName string
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

type RegisterRobotRequest struct {
	RobotName string
	RobotType remote.GrpcRobotType
	Port      int32
}

type RegisterRobotReply struct {
	RobotName string
}

type ListPlayerIDsReply struct {
	PlayerIDs []string
}

type StartGameRequest struct {
	Rule *mahjong.Rule
	Mode int
}

type StartGameReply struct {
	SeatsOrder []int
}

type StreamRequest struct {
	Call chan *mahjong.Call
}

type StreamReply struct {
	Events chan *player.GameEventChannel
}

type LogContent struct {
	Players   [4]string      `json:"players"`    // 四个玩家的名字
	PlayerIDs [4]uint        `json:"player_ids"` // 四个玩家的ID
	Events    mahjong.Events `json:"events"`     // 事件
}
