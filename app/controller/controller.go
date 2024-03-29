package controller

import (
	"context"
	"github.com/hphphp123321/mahjong-server/app/api/middleware/interceptor/authorization"
	"strings"
	"sync"

	mahjong "github.com/hphphp123321/mahjong-server/app/api/v1/mahjong"
	"github.com/hphphp123321/mahjong-server/app/global"
	"github.com/hphphp123321/mahjong-server/app/service/robot/remote"
	"github.com/hphphp123321/mahjong-server/app/service/server"
)

type MahjongServer struct {
	mahjong.UnimplementedMahjongServer
	s server.Server

	clients sync.Map

	methodsAuthExclude []string
}

func NewMahjongServer(s server.Server) *MahjongServer {
	return &MahjongServer{
		s: s,
		// clients: make(map[string]*Client),
		methodsAuthExclude: []string{
			"Login",
			"Logout",
			"Register",
			"Ping",
			"RegisterRobot",
		},
	}
}

func (m *MahjongServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	for _, method := range m.methodsAuthExclude {
		if strings.HasSuffix(fullMethodName, method) { // 不需要验证token的方法
			return ctx, nil
		}
	}

	return authorization.AuthInterceptor(ctx) // 验证token
}

func (m *MahjongServer) Ping(ctx context.Context, empty *mahjong.Empty) (*mahjong.Empty, error) {
	return &mahjong.Empty{}, nil
}

func (m *MahjongServer) Login(ctx context.Context, request *mahjong.LoginRequest) (*mahjong.LoginReply, error) {
	reply, err := m.s.Login(ctx, &server.LoginRequest{
		Name:     request.PlayerName,
		Password: request.Password,
	})
	if err != nil {
		return nil, err
	}
	cid := reply.ID
	c := NewClient(cid)
	m.clients.Store(cid, c)
	c.Login()
	if err != nil {
		return nil, err
	}
	// 生成JWT TOKEN
	token, err := authorization.GenerateToken(cid)
	return ToPbLoginReply(reply, token), nil
}

func (m *MahjongServer) Register(ctx context.Context, request *mahjong.RegisterRequest) (*mahjong.RegisterReply, error) {
	reply, err := m.s.Register(ctx, &server.RegisterRequest{
		Name:     request.PlayerName,
		Password: request.Password,
	})
	if err != nil {
		return nil, err
	}
	cid := reply.ID
	c := NewClient(cid)
	m.clients.Store(cid, c)
	c.Login()
	if err != nil {
		return nil, err
	}
	// 生成JWT TOKEN
	token, err := authorization.GenerateToken(cid)
	return ToPbRegisterReply(reply, token), nil
}

func (m *MahjongServer) Logout(ctx context.Context, empty *mahjong.Empty) (*mahjong.LogoutReply, error) {
	id, err := m.s.GetID(ctx)
	if err != nil {
		return nil, err
	}
	if err := m.s.Logout(ctx); err != nil {
		return nil, err
	}
	c, ok := m.clients.Load(id)
	if !ok {
		global.Log.Warnf("client %s not found", id)
	} else {
		c.(*Client).Logout()
	}
	return ToPbLogoutReply(), nil
}

func (m *MahjongServer) CreateRoom(ctx context.Context, request *mahjong.CreateRoomRequest) (*mahjong.CreateRoomReply, error) {
	reply, err := m.s.CreateRoom(ctx, &server.CreateRoomRequest{
		RoomName: request.RoomName,
	})
	if err != nil {
		return nil, err
	}
	return ToPbCreateRoomReply(reply), nil
}

func (m *MahjongServer) JoinRoom(ctx context.Context, request *mahjong.JoinRoomRequest) (*mahjong.JoinRoomReply, error) {
	reply, err := m.s.JoinRoom(ctx, &server.JoinRoomRequest{
		RoomID: request.RoomID,
	})
	if err != nil {
		return nil, err
	}
	if err := BoardCastReadyReply(ctx, m, ToPbPlayerJoinReply(reply)); err != nil {
		global.Log.Infoln("board cast error: ", err)
	}
	return ToPbJoinRoomReply(reply), nil
}

func (m *MahjongServer) ListRooms(ctx context.Context, request *mahjong.ListRoomsRequest) (*mahjong.ListRoomsReply, error) {
	reply, err := m.s.ListRooms(ctx, &server.ListRoomsRequest{
		RoomNameFilter: *request.RoomName,
	})
	if err != nil {
		return nil, err
	}
	return ToPbListRoomsReply(reply), nil
}

func (m *MahjongServer) ListRobots(ctx context.Context, empty *mahjong.Empty) (*mahjong.ListRobotsReply, error) {
	reply, err := m.s.ListRobots(ctx)
	if err != nil {
		return nil, err
	}
	return ToPbListRobotsReply(reply), nil
}

func (m *MahjongServer) RegisterRobot(ctx context.Context, request *mahjong.RegisterRobotRequest) (*mahjong.RegisterRobotReply, error) {

	reply, err := m.s.RegisterRobot(ctx, &server.RegisterRobotRequest{
		RobotName: request.RobotName,
		RobotType: remote.GrpcRobotType(request.RobotType),
		Port:      request.Port,
	})

	if err != nil {
		return nil, err
	}
	return ToPbRegisterRobotReply(reply), nil
}

func (m *MahjongServer) UnregisterRobot(ctx context.Context, request *mahjong.UnregisterRobotRequest) (*mahjong.Empty, error) {
	err := m.s.UnRegisterRobot(ctx, request.RobotName)
	if err != nil {
		return nil, err
	}
	return &mahjong.Empty{}, nil
}

func (m *MahjongServer) Ready(stream mahjong.Mahjong_ReadyServer) (err error) {
	ctx := stream.Context()
	if err := AddReadyStream(ctx, stream, m); err != nil {
		return err
	}
	done, replyChan := StartReadyStream(ctx, stream, m)
	for {
		select {
		case <-ctx.Done():
			goto TagBreak
		case err = <-done:
			if err != nil {
				global.Log.Warnln("ready stream done: ", err)
			}
			goto TagBreak
		case reply := <-replyChan:
			if err = BoardCastReadyReply(ctx, m, reply); err != nil {
				global.Log.Warnf("board cast error: %v", err)
			}
		}
	}
TagBreak:
	return RemoveReadyStream(ctx, m)
}

func (m *MahjongServer) Game(stream mahjong.Mahjong_GameServer) error {
	ctx := stream.Context()
	if err := AddGameStream(ctx, stream, m); err != nil {
		return err
	}
	recvDone, actionChan := StartGameRecvStream(ctx, stream, m)
	r, err := m.s.StartStream(ctx, ToServerStartStreamRequest(actionChan))
	if err != nil {
		return err
	}
	sendDone := StartGameSendStream(ctx, stream, r)
	select {
	case err = <-recvDone:
		if err != nil {
			global.Log.Warnln("game recv stream done: ", err)
		}
	case err = <-sendDone:
		if err != nil {
			global.Log.Warnln("game send stream done: ", err)
		}
	}
	close(actionChan)
	return RemoveGameStream(ctx, m)
}
