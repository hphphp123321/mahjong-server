package controller

import (
	"context"
	pb "github.com/hphphp123321/mahjong-server/app/api/v1"
	"github.com/hphphp123321/mahjong-server/app/global"
	"github.com/hphphp123321/mahjong-server/app/service/server"
)

type MahjongServer struct {
	pb.UnimplementedMahjongServer
	s server.Server

	clients map[string]*Client
}

func NewMahjongServer(s server.Server) *MahjongServer {
	return &MahjongServer{
		s:       s,
		clients: make(map[string]*Client),
	}
}

func (m MahjongServer) Ping(ctx context.Context, empty *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (m MahjongServer) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginReply, error) {
	reply, err := m.s.Login(ctx, &server.LoginRequest{
		Name: request.PlayerName,
	})
	cid := reply.ID
	c := NewClient(cid)
	m.clients[cid] = c
	c.Login()
	if err != nil {
		return nil, err
	}
	return ToPbLoginReply(reply), nil
}

func (m MahjongServer) Logout(ctx context.Context, empty *pb.Empty) (*pb.LogoutReply, error) {
	id, err := m.s.GetID(ctx)
	if err != nil {
		return nil, err
	}
	if err := m.s.Logout(ctx); err != nil {
		return nil, err
	}
	c, ok := m.clients[id]
	if !ok {
		global.Log.Warnf("client %s not found\n", id)
	} else {
		c.Logout()
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
		global.Log.Infoln("board cast error: ", err)
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

func (m MahjongServer) ListRobots(ctx context.Context, empty *pb.Empty) (*pb.ListRobotsReply, error) {
	reply, err := m.s.ListRobots(ctx)
	if err != nil {
		return nil, err
	}
	return ToPbListRobotsReply(reply), nil
}

func (m MahjongServer) Ready(stream pb.Mahjong_ReadyServer) error {
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
			global.Log.Warnf("board cast error: %v\n", err)
		}
	}
	return nil
}

func (m MahjongServer) Game(server pb.Mahjong_GameServer) error {
	//TODO implement me
	panic("implement me")
}
