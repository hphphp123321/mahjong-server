package server

import (
	"encoding/json"
	"errors"
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
	Players   [4]string      `json:"players"`         // 四个玩家的名字
	PlayerIDs [4]uint        `json:"player_ids"`      // 四个玩家的ID
	Events    mahjong.Events `json:"events"`          // 事件
	Error     error          `json:"error,omitempty"` // 错误
}

func (l *LogContent) MarshalJSON() ([]byte, error) {
	// 在函数内部定义临时结构体
	var tmp struct {
		Players   [4]string      `json:"players"`
		PlayerIDs [4]uint        `json:"player_ids"`
		Events    mahjong.Events `json:"events"`
		Error     string         `json:"error,omitempty"`
	}

	// 将LogContent的字段复制到临时结构体
	tmp.Players = l.Players
	tmp.PlayerIDs = l.PlayerIDs
	tmp.Events = l.Events

	// 特殊处理Error字段
	if l.Error != nil {
		tmp.Error = l.Error.Error()
	}

	// 序列化临时结构体
	return json.Marshal(tmp)
}

func (l *LogContent) UnmarshalJSON(data []byte) error {
	// 在函数内部定义临时结构体
	var tmp struct {
		Players   [4]string      `json:"players"`
		PlayerIDs [4]uint        `json:"player_ids"`
		Events    mahjong.Events `json:"events"`
		Error     string         `json:"error,omitempty"`
	}

	// 反序列化到临时结构体
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	// 将临时结构体的字段复制回LogContent
	l.Players = tmp.Players
	l.PlayerIDs = tmp.PlayerIDs
	l.Events = tmp.Events

	// 特殊处理Error字段
	if tmp.Error != "" {
		l.Error = errors.New(tmp.Error)
	} else {
		l.Error = nil
	}

	return nil
}
