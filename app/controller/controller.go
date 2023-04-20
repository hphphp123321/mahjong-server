package controller

import (
	"context"
	pb "github.com/hphphp123321/mahjong-server/app/api/v1"
	"github.com/hphphp123321/mahjong-server/app/service/playerid"
	"github.com/hphphp123321/mahjong-server/app/service/server"
)

type MahjongServer struct {
	pb.UnimplementedMahjongServer
	idService playerid.IDService
	s         server.Server
}

func (m MahjongServer) Ping(ctx context.Context, empty *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (m MahjongServer) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginReply, error) {
	reply, err := m.s.Login(ctx, &server.LoginRequest{
		Name: request.PlayerName,
	})
	if err != nil {
		return nil, err
	}
	return ToPbLoginReply(reply), nil
}

func (m MahjongServer) Logout(ctx context.Context, empty *pb.Empty) (*pb.LogoutReply, error) {
	if err := m.s.Logout(ctx); err != nil {
		return nil, err
	}
	return &pb.LogoutReply{
		Message: "",
	}, nil
}

func (m MahjongServer) CreateRoom(ctx context.Context, request *pb.CreateRoomRequest) (*pb.CreateRoomReply, error) {
	reply, err := m.s.CreateRoom(ctx, &server.CreateRoomRequest{
		RoomName: request.RoomName,
	})
	if err != nil {
		return nil, err
	}
	return ToPbCreateRoomReply(reply), nil
}

func (m MahjongServer) JoinRoom(ctx context.Context, request *pb.JoinRoomRequest) (*pb.JoinRoomReply, error) {
	reply, err := m.s.JoinRoom(ctx, &server.JoinRoomRequest{
		RoomID: request.RoomID,
	})
	if err != nil {
		return nil, err
	}
	// TODO BoardCast
	return ToPbJoinRoomReply(reply), nil
}

func (m MahjongServer) RefreshRoom(ctx context.Context, request *pb.RefreshRoomRequest) (*pb.RefreshRoomReply, error) {
	reply, err := m.s.RefreshRoom(ctx, &server.RefreshRoomRequest{
		RoomNameFilter: *request.RoomName,
	})
	if err != nil {
		return nil, err
	}
	return ToPbRefreshRoomReply(reply), nil
}

func (m MahjongServer) Ready(server pb.Mahjong_ReadyServer) error {
	//TODO implement me
	panic("implement me")
}

func (m MahjongServer) Start(server pb.Mahjong_StartServer) error {
	//TODO implement me
	panic("implement me")
}
