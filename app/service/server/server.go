package server

import (
	"context"
	"github.com/hphphp123321/mahjong-server/app/errs"
	"github.com/hphphp123321/mahjong-server/app/global"
	"github.com/hphphp123321/mahjong-server/app/model/player"
	"github.com/hphphp123321/mahjong-server/app/model/room"
	"google.golang.org/grpc/metadata"
	"regexp"
)

type ImplServer struct {
	players map[string]*player.Player
	rooms   map[string]*room.Room
}

func (i ImplServer) GetID(ctx context.Context) (string, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errs.ErrMetaDataNotFound
	}
	if id, ok := headers["id"]; !ok {
		return "", errs.ErrHeaderIDNotFound
	} else if _, ok = i.players[id[0]]; !ok {
		return "", errs.ErrPlayerNotFound
	} else {
		return id[0], nil
	}
}

func (i ImplServer) GetName(ctx context.Context) (string, error) {
	if id, err := i.GetID(ctx); err != nil {
		return "", err
	} else {
		return i.players[id].Name, nil
	}
}

func (i ImplServer) Login(ctx context.Context, request *LoginRequest) (reply *LoginReply, err error) {
	id, err := global.IDGenerator.GeneratePlayerID()
	if err != nil {
		return nil, err
	}
	i.players[id] = player.NewPlayer(id, request.Name)
	return &LoginReply{
		ID: id,
	}, nil
}

func (i ImplServer) Logout(ctx context.Context) error {
	if id, err := i.GetID(ctx); err != nil {
		return err
	} else {
		delete(i.players, id)
		return nil
	}
}

func (i ImplServer) CreateRoom(ctx context.Context, request *CreateRoomRequest) (reply *CreateRoomReply, err error) {
	pid, err := i.GetID(ctx)
	if err != nil {
		return nil, err
	}
	rid, err := global.IDGenerator.GenerateRoomID()
	if err != nil {
		return nil, err
	}
	r := room.NewRoom(i.players[pid], request.RoomName, rid)
	i.rooms[rid] = r
	return &CreateRoomReply{
		RoomID: rid,
	}, nil
}

func (i ImplServer) JoinRoom(ctx context.Context, request *JoinRoomRequest) (joinReply *JoinRoomReply, err error) {
	pid, err := i.GetID(ctx)
	if err != nil {
		return nil, err
	}
	p, ok := i.players[pid]
	if !ok {
		return nil, errs.ErrPlayerNotFound
	}
	r, ok := i.rooms[request.RoomID]
	if !ok {
		return nil, errs.ErrRoomNotFound
	}
	if err := r.Join(i.players[pid]); err != nil {
		return nil, err
	}
	return &JoinRoomReply{
		RoomInfo:   r.GetInfo(),
		Seat:       p.Seat,
		PlayerName: p.Name,
	}, nil
}

func (i ImplServer) ListRooms(ctx context.Context, request *ListRoomsRequest) (reply *ListRoomsReply, err error) {
	roomInfos := make([]*room.Info, 0)
	filter := request.RoomNameFilter
	for _, r := range i.rooms {
		matched, err := regexp.MatchString(filter, r.Name)
		if err != nil && !matched {
			continue
		}
		roomInfos = append(roomInfos, r.GetInfo())
	}
	return &ListRoomsReply{
		RoomInfos: roomInfos,
	}, nil
}

func (i ImplServer) GetReady(ctx context.Context) (reply *GetReadyReply, err error) {
	pid, err := i.GetID(ctx)
	if err != nil {
		return nil, err
	}
	p, ok := i.players[pid]
	if !ok {
		return nil, errs.ErrPlayerNotFound
	}
	if err := p.GetReady(); err != nil {
		return nil, err
	}
	return &GetReadyReply{
		Seat:       p.Seat,
		PlayerName: p.Name,
	}, nil
}

func (i ImplServer) CancelReady(ctx context.Context) (reply *CancelReadyReply, err error) {
	pid, err := i.GetID(ctx)
	if err != nil {
		return nil, err
	}
	p, ok := i.players[pid]
	if !ok {
		return nil, errs.ErrPlayerNotFound
	}
	if err := p.CancelReady(); err != nil {
		return nil, err
	}
	return &CancelReadyReply{
		Seat:       p.Seat,
		PlayerName: p.Name,
	}, nil
}

func (i ImplServer) LeaveRoom(ctx context.Context) (reply *PlayerLeaveReply, err error) {
	pid, err := i.GetID(ctx)
	if err != nil {
		return nil, err
	}
	p, ok := i.players[pid]
	if !ok {
		return nil, errs.ErrPlayerNotFound
	}
	if p.RoomID == "" {
		return nil, errs.ErrPlayerNotInRoom
	}
	r, ok := i.rooms[p.RoomID]
	if !ok {
		return nil, errs.ErrRoomNotFound
	}
	if err := r.Leave(p); err != nil {
		return nil, err
	}
	if r.IsEmpty() {
		// TODO log DEBUG
		delete(i.rooms, p.RoomID)
	}
	return &PlayerLeaveReply{
		Seat:       p.Seat,
		PlayerName: p.Name,
	}, nil
}

func (i ImplServer) AddRobot(ctx context.Context, request *AddRobotRequest) (reply *AddRobotReply, err error) {
	//TODO implement me
	panic("implement me")
}

func (i ImplServer) RemovePlayer(ctx context.Context, request *RemovePlayerRequest) (reply *PlayerLeaveReply, err error) {
	//TODO implement me
	panic("implement me")
}

func (i ImplServer) ListRobots(ctx context.Context) (reply *ListRobotsReply, err error) {
	//TODO implement me
	panic("implement me")
}

func (i ImplServer) ListPlayerIDs(ctx context.Context) (reply *ListPlayerIDsReply, err error) {
	//TODO implement me
	panic("implement me")
}