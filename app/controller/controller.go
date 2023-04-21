package controller

import (
	"context"
	pb "github.com/hphphp123321/mahjong-server/app/api/v1"
	"github.com/hphphp123321/mahjong-server/app/component/idservice"
	"github.com/hphphp123321/mahjong-server/app/service/server"
)

type MahjongServer struct {
	pb.UnimplementedMahjongServer
	idService idservice.IDService
	s         server.Server

	clients map[string]*Client
}

func (m MahjongServer) Ping(ctx context.Context, empty *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (m MahjongServer) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginReply, error) {
	reply, err := m.s.Login(ctx, &server.LoginRequest{
		Name: request.PlayerName,
	})
	cid := string(reply.ID)
	m.clients[cid] = NewClient(cid)
	if err != nil {
		return nil, err
	}
	return ToPbLoginReply(reply), nil
}

func (m MahjongServer) Logout(ctx context.Context, empty *pb.Empty) (*pb.LogoutReply, error) {
	if err := m.s.Logout(ctx); err != nil {
		return nil, err
	}
	return ToPbLogoutReply(), nil
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
	if err := BoardCastReadyReply(ctx, &m, ToPbPlayerJoinReply(reply)); err != nil {
		//TODO LOG
	}
	return ToPbJoinRoomReply(reply), nil
}

func (m MahjongServer) ListRooms(ctx context.Context, request *pb.ListRoomsRequest) (*pb.ListRoomsReply, error) {
	reply, err := m.s.ListRooms(ctx, &server.ListRoomsRequest{
		RoomNameFilter: *request.RoomName,
	})
	if err != nil {
		return nil, err
	}
	return ToPbListRoomsReply(reply), nil
}

func (m MahjongServer) Ready(stream pb.Mahjong_ReadyServer) error {
	//TODO implement me
	ctx := stream.Context()
	if err := AddReadyStream(ctx, stream, &m); err != nil {
		return err
	}
	done, replyChan := StartReadyStream(ctx, stream, &m)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	case reply := <-replyChan:
		if err := BoardCastReadyReply(ctx, &m, reply); err != nil {
			//TODO LOG
		}
	}
	return nil
}

func (m MahjongServer) Game(server pb.Mahjong_GameServer) error {
	//TODO implement me
	panic("implement me")
}
