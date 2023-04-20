package controller

import (
	"github.com/hphphp123321/go-common"
	pb "github.com/hphphp123321/mahjong-server/app/api/v1"
	"github.com/hphphp123321/mahjong-server/app/model/player"
	"github.com/hphphp123321/mahjong-server/app/model/room"
	"github.com/hphphp123321/mahjong-server/app/service/server"
)

func ToPbLoginReply(reply *server.LoginReply) *pb.LoginReply {
	return &pb.LoginReply{
		Message:       "",
		Id:            reply.ID,
		ReconnectInfo: nil,
	}
}

func ToPbCreateRoomReply(reply *server.CreateRoomReply) *pb.CreateRoomReply {
	return &pb.CreateRoomReply{
		Message: "",
		Room: &pb.RoomInfo{
			RoomID: reply.RoomID,
		},
	}
}

func MapToPlayerInfo(playerInfo *player.Info) *pb.PlayerInfo {
	return &pb.PlayerInfo{
		PlayerName: playerInfo.Name,
		PlayerSeat: int32(playerInfo.Seat),
		IsOwner:    playerInfo.IsOwner,
	}
}

func ToPbJoinRoomReply(reply *server.JoinRoomReply) *pb.JoinRoomReply {
	return &pb.JoinRoomReply{
		Message: "",
		Room: &pb.RoomInfo{
			RoomName: reply.RoomName,
			Players:  common.MapSlice(reply.RoomInfo.PlayerInfos, MapToPlayerInfo),
		},
	}
}

func MapToRoomInfo(roomInfo *room.Info) *pb.RoomInfo {
	return &pb.RoomInfo{
		RoomID:   roomInfo.ID,
		RoomName: roomInfo.Name,
		Players:  common.MapSlice(roomInfo.PlayerInfos, MapToPlayerInfo),
	}
}

func ToPbRefreshRoomReply(reply *server.RefreshRoomReply) *pb.RefreshRoomReply {
	return &pb.RefreshRoomReply{
		Message: "",
		Rooms:   common.MapSlice(reply.RoomInfos, MapToRoomInfo),
	}
}
