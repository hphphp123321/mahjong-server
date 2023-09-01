package controller

import (
	"github.com/hphphp123321/go-common"
	pb "github.com/hphphp123321/mahjong-server/app/api/v1"
	"github.com/hphphp123321/mahjong-server/app/service/server"
)

func ToPbPlayerJoinReply(r *server.JoinRoomReply) *pb.ReadyReply {
	return &pb.ReadyReply{
		Message: JoinMsg(r.PlayerName),
		Reply: &pb.ReadyReply_PlayerJoin{
			PlayerJoin: &pb.PlayerJoinReply{
				Seat:       int32(r.Seat),
				PlayerName: r.PlayerName,
			},
		},
	}
}

func ToPbGetReadyReply(r *server.GetReadyReply) *pb.ReadyReply {
	return &pb.ReadyReply{
		Message: GetReadyMsg(r.PlayerName),
		Reply: &pb.ReadyReply_GetReady{
			GetReady: &pb.GetReadyReply{
				Seat:       int32(r.Seat),
				PlayerName: r.PlayerName,
			},
		},
	}
}

func ToPbCancelReadyReply(r *server.CancelReadyReply) *pb.ReadyReply {
	return &pb.ReadyReply{
		Message: CancelReadyMsg(r.PlayerName),
		Reply: &pb.ReadyReply_CancelReady{
			CancelReady: &pb.CancelReadyReply{
				Seat:       int32(r.Seat),
				PlayerName: r.PlayerName,
			},
		},
	}
}

func ToPbLeaveRoomReply(r *server.PlayerLeaveReply) *pb.ReadyReply {
	return &pb.ReadyReply{
		Message: LeaveRoomMsg(r.PlayerName),
		Reply: &pb.ReadyReply_PlayerLeave{
			PlayerLeave: &pb.PlayerLeaveReply{
				Seat:       int32(r.Seat),
				PlayerName: r.PlayerName,
				OwnerSeat:  int32(r.OwnerSeat),
			},
		},
	}
}

func ToServerRemovePlayerRequest(r *pb.ReadyRequest) *server.RemovePlayerRequest {
	return &server.RemovePlayerRequest{
		Seat: int(r.GetRemovePlayer().PlayerSeat),
	}
}

func ToPbRemovePlayerReply(r *server.PlayerLeaveReply) *pb.ReadyReply {
	return &pb.ReadyReply{
		Message: RemovePlayerMsg(r.PlayerName),
		Reply: &pb.ReadyReply_PlayerLeave{
			PlayerLeave: &pb.PlayerLeaveReply{
				Seat:       int32(r.Seat),
				PlayerName: r.PlayerName,
				OwnerSeat:  int32(r.OwnerSeat),
			},
		},
	}
}

func ToServerAddRobotRequest(r *pb.ReadyRequest) *server.AddRobotRequest {
	return &server.AddRobotRequest{
		RobotSeat: int(r.GetAddRobot().RobotSeat),
		RobotType: r.GetAddRobot().RobotType,
	}
}

func ToPbAddRobotReply(r *server.AddRobotReply) *pb.ReadyReply {
	return &pb.ReadyReply{
		Message: AddRobotMsg(r.RobotName),
		Reply: &pb.ReadyReply_AddRobot{
			AddRobot: &pb.AddRobotReply{
				RobotSeat: int32(r.RobotSeat),
				RobotName: r.RobotName,
			},
		},
	}
}

func ToPbReadyChatReply(in *pb.ReadyRequest, playerName string, seat int) *pb.ReadyReply {
	return &pb.ReadyReply{
		Message: ChatMsg(playerName, in.GetChat().Message),
		Reply: &pb.ReadyReply_Chat{Chat: &pb.ChatReply{
			Message: in.GetChat().Message,
			Seat:    int32(seat),
		}},
	}
}

func ToServerStartGameRequest(in *pb.ReadyRequest) *server.StartGameRequest {
	return &server.StartGameRequest{
		Rule: ToMahjongGameRule(in.GetStartGame().GetGameRule()),
		Mode: int(in.GetStartGame().GetMode()),
	}
}

func ToPbStartGameReply(in *server.StartGameReply) *pb.ReadyReply {
	return &pb.ReadyReply{
		Message: StartGameMsg(),
		Reply: &pb.ReadyReply_StartGame{
			StartGame: &pb.StartGameReply{
				SeatsOrder: common.MapSlice(in.SeatsOrder, func(i int) int32 { return int32(i) }),
			},
		},
	}
}
