package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hphphp123321/go-common"
	"github.com/hphphp123321/mahjong-common/player"
	"github.com/hphphp123321/mahjong-common/robots"
	_ "github.com/hphphp123321/mahjong-common/robots/simple"
	"github.com/hphphp123321/mahjong-common/room"
	pb "github.com/hphphp123321/mahjong-common/services/mahjong/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"io"
	"strings"
	"sync"
	"time"
)

type MahjongServer struct {
	pb.UnimplementedMahjongServer
	clients    map[uuid.UUID]*client
	clientMu   sync.RWMutex
	maxClients int

	rooms  map[uuid.UUID]*room.Room
	roomMu sync.RWMutex
}

func NewMahjongServer(maxClients int) *MahjongServer {
	return &MahjongServer{
		clients:    make(map[uuid.UUID]*client),
		rooms:      make(map[uuid.UUID]*room.Room),
		maxClients: maxClients,
	}
}

func (s *MahjongServer) Ping(ctx context.Context, in *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *MahjongServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginReply, error) {
	s.clientMu.Lock()
	defer s.clientMu.Unlock()
	for t, p := range s.clients {
		if p.p.PlayerName == in.PlayerName {
			p.online = true
			return &pb.LoginReply{
				Message: "login success",
				Token:   t.String(),
			}, nil
		}
	}
	if len(s.clients) >= s.maxClients {
		return nil, errors.New("too many clients")
	}
	token := uuid.New()
	s.clients[token] = newClient(in.PlayerName, token)
	log.WithFields(log.Fields{
		"Event":      "Login",
		"PlayerName": in.PlayerName,
		"UUID":       token.String(),
	}).Info("player login success")
	return &pb.LoginReply{
		Message: "login success",
		Token:   token.String(),
	}, nil
}

func (s *MahjongServer) Logout(ctx context.Context, in *pb.Empty) (*pb.LogoutReply, error) {
	c, err := s.getClient(ctx)
	if err != nil {
		return nil, err
	}
	if c.p.RoomID != uuid.Nil {
		err = s.LeaveRoom(c)
		if err != nil {
			return nil, err
		}
	}
	c.readyStream = nil
	close(c.done)
	s.removeClient(c)
	log.WithFields(log.Fields{
		"Event":      "Logout",
		"PlayerName": c.p.PlayerName,
		"UUID":       c.p.Token.String(),
	}).Info("player logout success")
	return &pb.LogoutReply{
		Message: "logout success",
	}, nil
}

func (s *MahjongServer) RefreshRoom(ctx context.Context, in *pb.RefreshRoomRequest) (*pb.RefreshRoomReply, error) {
	_, err := s.getClient(ctx)
	if err != nil {
		return nil, err
	}
	roomSlice := make([]*pb.Room, 0)
	rName := ""
	if in.RoomName != nil {
		rName = *in.RoomName
	}
	for id, r := range s.rooms {
		if strings.Contains(r.RoomName, rName) {
			roomSlice = append(roomSlice, &pb.Room{
				RoomID:      id.String(),
				RoomName:    r.RoomName,
				PlayerCount: int32(r.PlayerCount),
				OwnerName:   r.Owner.PlayerName,
			})
		}
	}
	log.WithFields(log.Fields{
		"Event": "ListRooms",
	}).Debug("refresh room success")
	return &pb.RefreshRoomReply{
		Message: "refresh room success, room count: " + fmt.Sprint(len(roomSlice)),
		Rooms:   roomSlice,
	}, nil
}

func (s *MahjongServer) CreateRoom(ctx context.Context, in *pb.CreateRoomRequest) (*pb.CreateRoomReply, error) {
	c, err := s.getClient(ctx)
	if err != nil {
		return nil, err
	}
	if c.p.RoomID != uuid.Nil {
		return nil, errors.New("already in room")
	}
	s.clientMu.Lock()
	defer s.clientMu.Unlock()
	roomId := uuid.New()
	newRoom := room.NewRoom(roomId, in.RoomName, c.p)
	err = newRoom.AddPlayer(c.p)
	if err != nil {
		return nil, err
	}
	s.rooms[roomId] = newRoom
	c.p.RoomID = newRoom.RoomID

	log.WithFields(log.Fields{
		"Event":      "CreateRoom",
		"PlayerName": c.p.PlayerName,
		"RoomName":   in.RoomName,
		"UUID":       roomId.String(),
	}).Info("create room success")
	return &pb.CreateRoomReply{
		Message: fmt.Sprintf("Create Room Success! Room UUID: %s", roomId.String()),
		Room: &pb.Room{
			RoomID:      roomId.String(),
			RoomName:    in.RoomName,
			PlayerCount: 1,
			OwnerName:   c.p.PlayerName,
		}}, nil
}

func (s *MahjongServer) JoinRoom(ctx context.Context, in *pb.JoinRoomRequest) (*pb.JoinRoomReply, error) {
	c, err := s.getClient(ctx)
	if err != nil {
		return nil, err
	}
	if c.p.RoomID != uuid.Nil {
		return nil, errors.New("already in room")
	}
	s.clientMu.Lock()
	defer s.clientMu.Unlock()
	roomId, err := uuid.Parse(in.RoomID)
	if err != nil {
		return nil, err
	}
	joinRoom, ok := s.rooms[roomId]
	if !ok {
		return nil, errors.New("room not found")
	}
	if joinRoom.IsFull() {
		return nil, errors.New("room is full")
	}
	if err = joinRoom.AddPlayer(c.p); err != nil {
		return nil, err
	}
	c.p.RoomID = joinRoom.RoomID
	seat, err := joinRoom.GetSeat(c.p)
	if err != nil {
		return nil, err
	}

	rep := &pb.ReadyReply{
		Message: fmt.Sprintf("player: %s, join room", c.p.PlayerName),
		Reply: &pb.ReadyReply_PlayerJoin{PlayerJoin: &pb.PlayerJoinReply{
			Seat:       int32(seat),
			PlayerName: c.p.PlayerName,
		}},
	}
	err = s.readyBoardCast(c, rep, false)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"Event":      "JoinRoom",
		"PlayerName": c.p.PlayerName,
		"RoomName":   joinRoom.RoomName,
		"Seat":       seat,
	}).Info("join room success")
	return &pb.JoinRoomReply{
		Message: fmt.Sprintf("Join Room Success! Room UUID: %s", roomId.String()),
		Seat:    int32(seat),
		Room: &pb.Room{
			RoomID:      roomId.String(),
			RoomName:    joinRoom.RoomName,
			PlayerCount: int32(joinRoom.PlayerCount),
			OwnerName:   joinRoom.Owner.PlayerName,
		}}, nil
}

func (s *MahjongServer) Ready(stream pb.Mahjong_ReadyServer) error {
	ctx := stream.Context()
	c, err := s.getClient(ctx)
	if err != nil {
		return err
	}
	if c.p.RoomID == uuid.Nil {
		return errors.New("not in room")
	}
	if c.readyStream != nil {
		return errors.New("already has ready stream")
	}
	c.readyStream = stream
	log.Infof("Start new ReadyStream for player: %s", c.p.PlayerName)
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				c.done <- nil
				return
			}
			if err != nil {
				log.Warningf("receive error %v", err)
				c.done <- errors.New("failed to receive request")
				return
			}
			switch in.GetRequest().(type) {
			case *pb.ReadyRequest_GetReady:
				err = s.handleGetReadyRequest(c, in)
				if err != nil {
					return
				}
			case *pb.ReadyRequest_CancelReady:
				err = s.handleCancelReadyRequest(c, in)
				if err != nil {
					return
				}
			case *pb.ReadyRequest_LeaveRoom:
				err = s.handleLeaveRoomRequest(c, in)
				if err != nil {
					return
				}
			case *pb.ReadyRequest_RemovePlayer:
				err = s.handleRemovePlayerRequest(c, in)
				if err != nil {
					return
				}
			case *pb.ReadyRequest_AddRobot:
				err = s.handleAddRobotRequest(c, in)
				if err != nil {
					return
				}
			case *pb.ReadyRequest_Chat:
				rep := &pb.ReadyReply{
					Message: fmt.Sprintf("player: %s, chat: %s", c.p.PlayerName, in.GetChat().Message),
					Reply: &pb.ReadyReply_Chat{Chat: &pb.ChatReply{
						Message:    in.GetChat().Message,
						PlayerName: c.p.PlayerName,
					}},
				}
				err = s.readyBoardCast(c, rep, true)
				if err != nil {
					c.done <- err
					return
				}
			}
		}
	}()
	var doneError error
	select {
	case <-ctx.Done():
		doneError = ctx.Err()
	case <-c.done:
		log.Info("ReadyStream done for player: ", c.p.PlayerName)
	}
	if doneError != nil {
		return doneError
	}
	return nil
}

func (s *MahjongServer) readyBoardCast(c *client, resp *pb.ReadyReply, includeSelf bool) error {
	if c.p.RoomID == uuid.Nil {
		return errors.New("not in room")
	}
	r, err := s.getRoomByClient(c)
	if err != nil {
		return err
	}
	for _, p := range r.Players {
		if p.Token == c.p.Token {
			if !includeSelf {
				continue
			}
		}
		if s.clients[p.Token].readyStream == nil {
			return errors.New("don't have ready stream")
		}
		if err := s.clients[p.Token].readyStream.Send(resp); err != nil {
			return err
		}
	}
	return nil
}

func (s *MahjongServer) getToken(ctx context.Context) (uuid.UUID, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.UUID{}, errors.New("no metadata in context")
	}
	token, err := uuid.Parse(headers["token"][0])
	_, ok = s.clients[token]
	if err != nil || !ok {
		return uuid.UUID{}, errors.New("invalid token")
	}
	return token, nil
}

func (s *MahjongServer) getClient(ctx context.Context) (*client, error) {
	s.clientMu.RLock()
	defer s.clientMu.RUnlock()
	token, err := s.getToken(ctx)
	if err != nil {
		return nil, err
	}
	c, ok := s.clients[token]
	c.lastTime = time.Now()
	if !ok {
		return nil, errors.New("invalid token")
	}
	return c, nil
}

func (s *MahjongServer) removeClient(c *client) {
	s.clientMu.Lock()
	delete(s.clients, c.p.Token)
	s.clientMu.Unlock()
}

func (s *MahjongServer) getRoomByClient(c *client) (*room.Room, error) {
	s.roomMu.RLock()
	defer s.roomMu.RUnlock()
	if c.p.RoomID == uuid.Nil {
		return nil, errors.New("not in room")
	}
	cRoom, ok := s.rooms[c.p.RoomID]
	if !ok {
		return nil, errors.New("room not exist")
	}
	return cRoom, nil
}

// LeaveRoom detect when client disconnects or leaves room
func (s *MahjongServer) LeaveRoom(c *client) error {
	r, err := s.getRoomByClient(c)
	if err != nil {
		return err
	}
	if r == nil {
		return nil
	}
	log.Debugf("LeaveRoom: PlayerName: %s, RoomName: %s", c.p.PlayerName, r.RoomName)

	roomID := r.RoomID
	if err := r.RemovePlayer(c.p); err != nil {
		return err
	}
	if r.IsEmpty() {
		s.roomMu.Lock()
		delete(s.rooms, roomID)
		s.roomMu.Unlock()
		log.Printf("Room %s is empty, delete", roomID.String())
	} else {
		rep := &pb.ReadyReply{Message: fmt.Sprintf("player: %s, leave room", c.p.PlayerName),
			Reply: &pb.ReadyReply_PlayerLeave{PlayerLeave: &pb.PlayerLeaveReply{
				Seat:       int32(c.p.Seat),
				OwnerSeat:  int32(r.Owner.Seat),
				PlayerName: c.p.PlayerName,
			}}}
		if err := s.readyBoardCast(c, rep, true); err != nil {
			return err
		}
	}
	log.WithFields(log.Fields{
		"PlayerName": c.p.PlayerName,
		"RoomName":   r.RoomName,
	}).Debug("LeaveRoom success")
	r = nil
	c.readyStream = nil
	c.p.RoomID = uuid.Nil
	return nil
}

func (s *MahjongServer) handleGetReadyRequest(c *client, in *pb.ReadyRequest) error {
	var err error
	log.Debugf("GetReady Req: PlayerName: %s,  request: %s", c.p.PlayerName, in.GetGetReady().String())
	if c.p.Ready {
		log.Warning("Player already ready")
		return nil
	}
	c.p.SetReady(true)

	rep := &pb.ReadyReply{
		Message: fmt.Sprintf("player: %s, get ready success", c.p.PlayerName),
		Reply: &pb.ReadyReply_GetReady{GetReady: &pb.GetReadyReply{
			Seat:       int32(c.p.Seat),
			PlayerName: c.p.PlayerName,
		}},
	}
	err = s.readyBoardCast(c, rep, true)
	if err != nil {
		c.done <- err
		return err
	}
	log.WithFields(log.Fields{
		"Event":      "GetReady",
		"PlayerName": c.p.PlayerName,
		"Seat":       c.p.Seat,
	}).Info("Player Get Ready Success")
	return nil
}

func (s *MahjongServer) handleCancelReadyRequest(c *client, in *pb.ReadyRequest) error {
	var err error
	log.Debugf("CancelReady Req: PlayerName: %s,  request: %s", c.p.PlayerName, in.GetCancelReady().String())
	c.p.SetReady(false)

	rep := &pb.ReadyReply{
		Message: fmt.Sprintf("player: %s, cancel ready success", c.p.PlayerName),
		Reply: &pb.ReadyReply_CancelReady{CancelReady: &pb.CancelReadyReply{
			Seat:       int32(c.p.Seat),
			PlayerName: c.p.PlayerName,
		}},
	}
	err = s.readyBoardCast(c, rep, true)
	if err != nil {
		c.done <- err
		return err
	}
	log.WithFields(log.Fields{
		"Event":      "CancelReady",
		"PlayerName": c.p.PlayerName,
		"Seat":       c.p.Seat,
	}).Info("Player Cancel Ready success")
	return nil
}

func (s *MahjongServer) handleLeaveRoomRequest(c *client, in *pb.ReadyRequest) error {
	var err error
	log.Debugf("LeaveRoom Req: PlayerName: %s,  request: %s", c.p.PlayerName, in.GetLeaveRoom().String())
	err = s.LeaveRoom(c)
	if err != nil {
		c.done <- err
		return err
	}
	log.WithFields(log.Fields{
		"Event":      "LeaveRoom",
		"PlayerName": c.p.PlayerName,
		"Seat":       c.p.Seat,
	}).Info("Player Leave Room success")
	return nil
}

func (s *MahjongServer) handleRemovePlayerRequest(c *client, in *pb.ReadyRequest) error {
	var err error
	r, err := s.getRoomByClient(c)
	if err != nil {
		return err
	}
	if r.Owner != c.p {
		err = c.sendReadyMessage(fmt.Sprintf("player: %s, is not owner, can't remove player", c.p.PlayerName))
		if err != nil {
			return err
		}
	}
	log.Debugf("RemovePlayer Req: PlayerName: %s, RoomName: %s, request: %s", c.p.PlayerName, r.RoomName, in.GetRemovePlayer().String())
	seat := int(in.GetRemovePlayer().PlayerSeat)
	playerToRemove, err := r.GetPlayerBySeat(seat)
	if err != nil {
		return err
	}
	if playerToRemove == c.p {
		err = c.sendReadyMessage(fmt.Sprintf("player: %s, can't remove self", c.p.PlayerName))
		if err != nil {
			return err
		}
	}
	err = r.RemovePlayer(playerToRemove)
	if err != nil {
		c.done <- err
		return err
	}
	rep := &pb.ReadyReply{
		Message: fmt.Sprintf("player: %s, remove player: %s success", c.p.PlayerName, playerToRemove.PlayerName),
		Reply: &pb.ReadyReply_PlayerLeave{PlayerLeave: &pb.PlayerLeaveReply{
			Seat:       int32(playerToRemove.Seat),
			PlayerName: playerToRemove.PlayerName,
			OwnerSeat:  int32(r.Owner.Seat),
		}},
	}
	err = s.readyBoardCast(c, rep, true)
	if err != nil {
		c.done <- err
		return err
	}
	log.WithFields(log.Fields{
		"Event":             "RemovePlayer",
		"PlayerName":        c.p.PlayerName,
		"RemovedPlayerName": playerToRemove.PlayerName,
		"RoomName":          r.RoomName,
		"Seat":              seat,
	}).Info("Player Remove Player success")
	return nil
}

func (s *MahjongServer) handleAddRobotRequest(c *client, in *pb.ReadyRequest) error {
	var err error
	r, err := s.getRoomByClient(c)
	if err != nil {
		return err
	}
	if r.Owner != c.p {
		err = c.sendReadyMessage(fmt.Sprintf("player: %s, is not owner, can't add robot", c.p.PlayerName))
		if err != nil {
			return err
		}
	}
	log.Debugf("AddRobot Req: PlayerName: %s, RoomName: %s, request: %s", c.p.PlayerName, r.RoomName, in.GetAddRobot().String())
	seat := int(in.GetAddRobot().RobotSeat)
	if !common.SliceContain(r.IdleSeats, seat) {
		err = c.sendReadyMessage(fmt.Sprintf("seat %v not valid", seat))
		if err != nil {
			return err
		}
	}
	level := in.GetAddRobot().RobotLevel
	robot, err := robots.GetRobot(level)
	if err != nil {
		err = c.sendReadyMessage(fmt.Sprintf("robot level %s not valid", level))
		if err != nil {
			return err
		}
	}
	robotPlayer := player.NewRobot(level, seat, robot)
	err = r.AddRobot(robotPlayer)
	if err != nil {
		err = c.sendReadyMessage(err.Error())
		if err != nil {
			return err
		}
	}
	rep := &pb.ReadyReply{
		Message: fmt.Sprintf("player: %s, add robot: %s success", c.p.PlayerName, robotPlayer.PlayerName),
		Reply: &pb.ReadyReply_AddRobot{AddRobot: &pb.AddRobotReply{
			RobotSeat:  int32(robotPlayer.Seat),
			RobotLevel: level,
		}},
	}
	err = s.readyBoardCast(c, rep, true)
	log.WithFields(log.Fields{
		"Event":      "AddRobot",
		"PlayerName": c.p.PlayerName,
		"RoomName":   r.RoomName,
		"RobotLevel": level,
		"Seat":       seat,
	}).Info("Player Add Robot success")
	if err != nil {
		c.done <- err
		return err
	}
	return nil
}

//func (s *MahjongServer) CheckClients() {
//	for {
//		time.Sleep(30 * time.Second)
//		s.clientMu.Lock()
//		for token, c := range s.clients {
//			if c.online {
//				if time.Since(c.lastTime) > 60*time.Second {
//					c.online = false
//
//				}
//			} else {
//				if time.Since(c.lastTime) > 60*time.Second {
//					if c.room != nil {
//						if err := c.room.RemovePlayer(c.p); err != nil {
//							log.Printf("CheckClients: %s", err)
//						}
//						if c.room.PlayerCount == 0 {
//							delete(s.rooms, c.room.RoomID)
//							log.Printf("Room %s is empty, delete", c.room.RoomID.String())
//						}
//						c.room = nil
//					}
//					delete(s.clients, token)
//				}
//			}
//		}
//		s.clientMu.Unlock()
//	}
//}
