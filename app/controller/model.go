package controller

import (
	"github.com/hphphp123321/go-common"
	"github.com/hphphp123321/mahjong-go/mahjong"
	pb "github.com/hphphp123321/mahjong-server/app/api/v1/mahjong"
	"github.com/hphphp123321/mahjong-server/app/model/player"
	"github.com/hphphp123321/mahjong-server/app/model/room"
	"github.com/hphphp123321/mahjong-server/app/service/server"
)

func ToPbLoginReply(reply *server.LoginReply) *pb.LoginReply {
	return &pb.LoginReply{
		Message:       LoginMsg(reply.ID),
		Id:            reply.ID,
		ReconnectInfo: nil,
	}
}

func ToPbLogoutReply() *pb.LogoutReply {
	return &pb.LogoutReply{
		Message: LogoutMsg(),
	}
}

func ToPbCreateRoomReply(reply *server.CreateRoomReply) *pb.CreateRoomReply {
	return &pb.CreateRoomReply{
		Message: CreateMsg(reply.RoomInfo.ID),
		Room:    MapToRoomInfo(reply.RoomInfo),
	}
}

func MapToPlayerInfo(playerInfo *player.Info) *pb.PlayerInfo {
	return &pb.PlayerInfo{
		PlayerName: playerInfo.Name,
		PlayerSeat: int32(playerInfo.Seat),
		IsReady:    playerInfo.Ready,
	}
}

func ToPbJoinRoomReply(reply *server.JoinRoomReply) *pb.JoinRoomReply {
	return &pb.JoinRoomReply{
		Message: JoinRoomMsg(reply.RoomInfo.Name),
		Room:    MapToRoomInfo(reply.RoomInfo),
		Seat:    int32(reply.Seat),
	}
}

func MapToRoomInfo(roomInfo *room.Info) *pb.RoomInfo {
	return &pb.RoomInfo{
		RoomID:    roomInfo.ID,
		RoomName:  roomInfo.Name,
		OwnerSeat: int32(roomInfo.OwnerSeat),
		Players:   common.MapSlice(roomInfo.PlayerInfos, MapToPlayerInfo),
	}
}

func ToPbListRoomsReply(reply *server.ListRoomsReply) *pb.ListRoomsReply {
	var roomNames []string
	for _, roomInfo := range reply.RoomInfos {
		roomNames = append(roomNames, roomInfo.Name)
	}
	return &pb.ListRoomsReply{
		Message: ListRoomsMsg(roomNames),
		Rooms:   common.MapSlice(reply.RoomInfos, MapToRoomInfo),
	}
}

func ToPbListRobotsReply(reply *server.ListRobotsReply) *pb.ListRobotsReply {
	return &pb.ListRobotsReply{
		RobotTypes: reply.RobotTypes,
	}
}

func ToPbRegisterRobotReply(reply *server.RegisterRobotReply) *pb.RegisterRobotReply {
	return &pb.RegisterRobotReply{
		Message: RegisterRobotMsg(reply.RobotName),
	}
}

func ToServerStartStreamRequest(actionChan chan *mahjong.Call) *server.StreamRequest {
	return &server.StreamRequest{
		Call: actionChan,
	}
}
