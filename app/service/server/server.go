package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hphphp123321/mahjong-server/app/dao"
	"github.com/hphphp123321/mahjong-server/app/dao/entity"
	"github.com/hphphp123321/mahjong-server/app/dao/query"
	"net"
	"regexp"
	"strconv"
	"sync"

	"github.com/hphphp123321/mahjong-go/mahjong"
	"github.com/hphphp123321/mahjong-server/app/errs"
	"github.com/hphphp123321/mahjong-server/app/global"
	"github.com/hphphp123321/mahjong-server/app/model/player"
	"github.com/hphphp123321/mahjong-server/app/model/room"
	"github.com/hphphp123321/mahjong-server/app/service/robot/remote"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type ImplServer struct {
	Server
	players map[uint]*player.Player
	rooms   map[string]*room.Room

	lock sync.Mutex
}

func NewImplServer() *ImplServer {
	return &ImplServer{
		players: map[uint]*player.Player{},
		rooms:   map[string]*room.Room{},
	}
}

func (i *ImplServer) GetID(ctx context.Context) (string, error) {
	id, err := i.getID(ctx)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(id)), nil
}

func (i *ImplServer) getID(ctx context.Context) (uint, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errs.ErrMetaDataNotFound
	}
	id, ok := headers["id"]
	if !ok {
		return 0, errs.ErrHeaderIDNotFound
	}
	uid, err := ConvertStringToUInt(id[0])
	if err != nil {
		return 0, err
	}
	if _, ok = i.players[uid]; !ok {
		return 0, errs.ErrPlayerNotFound
	} else {
		return uid, nil
	}
}

func (i *ImplServer) GetName(ctx context.Context) (string, error) {
	if id, err := i.getID(ctx); err != nil {
		return "", err
	} else {
		return i.players[id].Name, nil
	}
}

func (i *ImplServer) GetRoomInfo(ctx context.Context) (*room.Info, error) {
	if id, err := i.getID(ctx); err != nil {
		return nil, err
	} else if p, ok := i.players[id]; !ok {
		return nil, errs.ErrPlayerNotFound
	} else if p.RoomID == "" {
		return nil, errs.ErrPlayerNotInRoom
	} else if r, ok := i.rooms[p.RoomID]; !ok {
		return nil, errs.ErrRoomNotFound
	} else {
		return r.GetInfo(), nil
	}
}

func (i *ImplServer) GetSeat(ctx context.Context) (int, error) {
	if id, err := i.getID(ctx); err != nil {
		return 0, err
	} else if p, ok := i.players[id]; !ok {
		return 0, errs.ErrPlayerNotFound
	} else if p.RoomID == "" {
		return 0, errs.ErrPlayerNotInRoom
	} else if _, ok := i.rooms[p.RoomID]; !ok {
		return 0, errs.ErrRoomNotFound
	} else {
		return p.Seat, nil
	}
}

func (i *ImplServer) GetIDBySeat(ctx context.Context, seat int) (string, error) {
	if id, err := i.getID(ctx); err != nil {
		return "", err
	} else if p, ok := i.players[id]; !ok {
		return "", errs.ErrPlayerNotFound
	} else if p.RoomID == "" {
		return "", errs.ErrPlayerNotInRoom
	} else if r, ok := i.rooms[p.RoomID]; !ok {
		return "", errs.ErrRoomNotFound
	} else if rp, err := r.GetPlayerBySeat(seat); err != nil {
		return "", err
	} else {
		return strconv.Itoa(int(rp.ID)), nil
	}
}

func (i *ImplServer) getPlayer(ctx context.Context) (*player.Player, error) {
	pid, err := i.getID(ctx)
	if err != nil {
		return nil, err
	}
	p, ok := i.players[pid]
	if !ok {
		return nil, errs.ErrPlayerNotFound
	}
	return p, nil
}

func (i *ImplServer) getRoom(ctx context.Context) (*room.Room, error) {
	p, err := i.getPlayer(ctx)
	if err != nil {
		return nil, err
	}
	r, ok := i.rooms[p.RoomID]
	if !ok {
		return nil, errs.ErrRoomNotFound
	}
	return r, nil
}

func (i *ImplServer) Login(ctx context.Context, request *LoginRequest) (reply *LoginReply, err error) {

	// 检查用户名是否存在
	user, err := query.User.WithContext(dao.DBCtx).Where(query.User.Name.Eq(request.Name)).First()
	if err != nil {
		return nil, err
	}

	// 检查密码是否正确
	if err := global.PasswordService.Compare([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, err
	}

	i.lock.Lock()
	i.players[user.ID] = player.NewPlayer(user.ID, request.Name)
	i.lock.Unlock()
	// add id to header
	var ids = strconv.Itoa(int(user.ID))
	header := metadata.New(map[string]string{"id": ids})
	ctx = metadata.NewOutgoingContext(ctx, header)
	return &LoginReply{
		ID: ids,
	}, nil
}

func (i *ImplServer) Logout(ctx context.Context) error {
	if id, err := i.getID(ctx); err != nil {
		return err
	} else {
		p := i.players[id]
		if p.RoomID != "" {
			roomID := p.RoomID
			r := i.rooms[roomID]
			err := r.Leave(p)
			if err != nil {
				return err
			}
			if r.IsEmpty() {
				global.Log.Infof("room id: %s is empty, delete", r.ID)
				i.lock.Lock()
				delete(i.rooms, roomID)
				i.lock.Unlock()
			}
		}
		i.lock.Lock()
		delete(i.players, id)
		i.lock.Unlock()
		return nil
	}
}

func (i *ImplServer) Register(ctx context.Context, request *RegisterRequest) (reply *RegisterReply, err error) {
	id, err := global.IDGenerator.GeneratePlayerID()
	if err != nil {
		return nil, err
	}

	encryptedPassword, err := global.PasswordService.Encrypt([]byte(request.Password))
	if err != nil {
		return nil, err
	}

	var newUser = &entity.User{
		ID:       uint(id),
		Name:     request.Name,
		Password: string(encryptedPassword),
		Logs:     nil,
	}
	// 检查用户名是否已存在
	count, err := query.User.WithContext(dao.DBCtx).Where(query.User.Name.Eq(request.Name)).Count()
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errs.ErrUserNameExist
	}

	// 创建用户
	if err := query.User.WithContext(dao.DBCtx).Create(newUser); err != nil {
		return nil, err
	}

	// 创建玩家
	i.lock.Lock()
	i.players[uint(id)] = player.NewPlayer(uint(id), request.Name)
	i.lock.Unlock()
	// add id to header
	uid := strconv.Itoa(int(id))
	header := metadata.New(map[string]string{"id": uid})
	ctx = metadata.NewOutgoingContext(ctx, header)

	return &RegisterReply{
		ID: uid,
	}, nil
}

func (i *ImplServer) CreateRoom(ctx context.Context, request *CreateRoomRequest) (reply *CreateRoomReply, err error) {
	pid, err := i.getID(ctx)
	if err != nil {
		return nil, err
	}
	rid, err := global.IDGenerator.GenerateRoomID()
	if err != nil {
		return nil, err
	}
	r, err := room.NewRoom(i.players[pid], request.RoomName, rid)
	if err != nil {
		return nil, err
	}
	i.lock.Lock()
	i.rooms[rid] = r
	i.lock.Unlock()
	return &CreateRoomReply{
		RoomInfo: r.GetInfo(),
	}, nil
}

func (i *ImplServer) JoinRoom(ctx context.Context, request *JoinRoomRequest) (joinReply *JoinRoomReply, err error) {
	p, err := i.getPlayer(ctx)
	if err != nil {
		return nil, err
	}
	r, ok := i.rooms[request.RoomID]
	if !ok {
		return nil, errs.ErrRoomNotFound
	}
	if err := r.Join(p); err != nil {
		return nil, err
	}
	return &JoinRoomReply{
		RoomInfo:   r.GetInfo(),
		Seat:       p.Seat,
		PlayerName: p.Name,
	}, nil
}

func (i *ImplServer) ListRooms(ctx context.Context, request *ListRoomsRequest) (reply *ListRoomsReply, err error) {
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

func (i *ImplServer) GetReady(ctx context.Context) (reply *GetReadyReply, err error) {
	p, err := i.getPlayer(ctx)
	if err != nil {
		return nil, err
	}
	if err := p.GetReady(); err != nil {
		return nil, err
	}
	return &GetReadyReply{
		Seat:       p.Seat,
		PlayerName: p.Name,
	}, nil
}

func (i *ImplServer) CancelReady(ctx context.Context) (reply *CancelReadyReply, err error) {
	p, err := i.getPlayer(ctx)
	if err != nil {
		return nil, err
	}
	if err := p.CancelReady(); err != nil {
		return nil, err
	}
	return &CancelReadyReply{
		Seat:       p.Seat,
		PlayerName: p.Name,
	}, nil
}

func (i *ImplServer) LeaveRoom(ctx context.Context) (reply *PlayerLeaveReply, err error) {
	p, err := i.getPlayer(ctx)
	seat := p.Seat
	if err != nil {
		return nil, err
	}
	if p.RoomID == "" {
		return nil, errs.ErrPlayerNotInRoom
	}
	id := p.RoomID
	r, ok := i.rooms[id]
	if !ok {
		return nil, errs.ErrRoomNotFound
	}
	if err := r.Leave(p); err != nil {
		return nil, err
	}
	if r.IsEmpty() {
		global.Log.Infof("room id: %s is empty, delete", r.ID)
		i.lock.Lock()
		delete(i.rooms, id)
		i.lock.Unlock()
	}
	return &PlayerLeaveReply{
		Seat:       seat,
		PlayerName: p.Name,
		OwnerSeat:  r.OwnerSeat,
	}, nil
}

func (i *ImplServer) AddRobot(ctx context.Context, request *AddRobotRequest) (reply *AddRobotReply, err error) {
	p, err := i.getPlayer(ctx)
	if err != nil {
		return nil, err
	}
	if p.RoomID == "" {
		return nil, errs.ErrPlayerNotInRoom
	}
	r, ok := i.rooms[p.RoomID]
	if !ok {
		return nil, errs.ErrRoomNotFound
	}
	if p.Seat != r.OwnerSeat {
		return nil, errs.ErrPlayerNotOwner
	}
	robotName, err := r.AddRobot(request.RobotType, request.RobotSeat)
	if err != nil {
		return nil, err
	}
	return &AddRobotReply{
		RobotSeat: request.RobotSeat,
		RobotName: robotName,
	}, nil
}

func (i *ImplServer) RemovePlayer(ctx context.Context, request *RemovePlayerRequest) (reply *PlayerLeaveReply, err error) {
	p, err := i.getPlayer(ctx)
	if err != nil {
		return nil, err
	}
	if p.RoomID == "" {
		return nil, errs.ErrPlayerNotInRoom
	}
	r, ok := i.rooms[p.RoomID]
	if !ok {
		return nil, errs.ErrRoomNotFound
	}
	if p.Seat != r.OwnerSeat {
		return nil, errs.ErrPlayerNotOwner
	}
	p2Remove, err := r.GetPlayerBySeat(request.Seat)
	p2Seat := p2Remove.Seat
	p2Name := p2Remove.Name
	if err != nil {
		return nil, err
	}
	if err := r.LeaveRoomBySeat(p2Seat); err != nil {
		return nil, err
	}
	return &PlayerLeaveReply{
		Seat:       p2Seat,
		OwnerSeat:  r.OwnerSeat,
		PlayerName: p2Name,
	}, nil
}

func (i *ImplServer) ListRobots(ctx context.Context) (reply *ListRobotsReply, err error) {
	_, err = i.getRoom(ctx)
	if err != nil {
		return nil, err
	}
	TestAllRobots()
	return &ListRobotsReply{
		RobotTypes: global.RobotRegistry.GetRobotTypes(),
	}, nil
}

func (i *ImplServer) RegisterRobot(ctx context.Context, request *RegisterRobotRequest) (reply *RegisterRobotReply, err error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("peer not found")
	}
	clientIP, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		return nil, err
	}
	robotAddr := fmt.Sprintf("%s:%d", clientIP, request.Port)
	global.Log.Infof("robot addr: %s", robotAddr)

	grpcRobot, err := remote.NewGrpcRobot(request.RobotName, request.RobotType, robotAddr)
	if err != nil {
		return nil, err
	}
	if err = TestRobot(grpcRobot); err != nil {
		return nil, err
	}
	if err = global.RobotRegistry.Register(grpcRobot); err != nil {
		return nil, err
	}
	global.Log.Infof("robot %s registered", grpcRobot.Name)
	return &RegisterRobotReply{
		RobotName: grpcRobot.Name,
	}, nil
}

func (i *ImplServer) UnRegisterRobot(ctx context.Context, robotName string) (err error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return errors.New("peer not found")
	}
	clientIP, _, err := net.SplitHostPort(p.Addr.String())
	if err != nil {
		return err
	}

	robot, ok := global.RobotRegistry.GetRobot(robotName)
	if !ok {
		return errs.ErrRobotNotFound
	}
	grpcRobot, ok := robot.(*remote.GrpcRobot)
	if !ok {
		return errors.New("robot type not grpc")
	}
	if grpcRobot.IP != clientIP {
		return errors.New("client ip not match")
	}

	if err = global.RobotRegistry.Unregister(robotName); err != nil {
		return err
	}
	global.Log.Infof("robot %s unregistered", robotName)
	return nil
}

func (i *ImplServer) ListPlayerIDs(ctx context.Context) (reply *ListPlayerIDsReply, err error) {
	r, err := i.getRoom(ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0)
	for _, p := range r.Players {
		if p.ID == 0 {
			continue
		}
		ids = append(ids, strconv.Itoa(int(p.ID)))
	}
	return &ListPlayerIDsReply{
		PlayerIDs: ids,
	}, nil
}

func (i *ImplServer) StartGame(ctx context.Context, request *StartGameRequest) (reply *StartGameReply, err error) {
	p, err := i.getPlayer(ctx)
	if err != nil {
		return nil, err
	}
	if p.RoomID == "" {
		return nil, errs.ErrPlayerNotInRoom
	}
	r, ok := i.rooms[p.RoomID]
	if !ok {
		return nil, errs.ErrRoomNotFound
	}
	if p.Seat != r.OwnerSeat {
		return nil, errs.ErrPlayerNotOwner
	}
	gameInfo, err := r.StartGame(request.Rule, request.Mode)
	if err != nil {
		return nil, err
	}
	i.ProcessGameResult(ctx, gameInfo) // start process game result
	return &StartGameReply{
		SeatsOrder: gameInfo.SeatsOrder,
	}, nil
}

func (i *ImplServer) StartStream(ctx context.Context, request *StreamRequest) (reply *StreamReply, err error) {
	p, err := i.getPlayer(ctx)
	if err != nil {
		return nil, err
	}
	if p.RoomID == "" {
		return nil, errs.ErrPlayerNotInRoom
	}
	r, ok := i.rooms[p.RoomID]
	if !ok {
		return nil, errs.ErrRoomNotFound
	}
	ech := r.StartGameStream(p, request.Call)
	return &StreamReply{
		Events: ech,
	}, nil
}

func (i *ImplServer) GetBoardState(ctx context.Context) (*mahjong.BoardState, error) {
	p, err := i.getPlayer(ctx)
	if err != nil {
		return nil, err
	}
	if p.RoomID == "" {
		return nil, errs.ErrPlayerNotInRoom
	}
	r, ok := i.rooms[p.RoomID]
	if !ok {
		return nil, errs.ErrRoomNotFound
	}
	return r.GetBoardState(p)
}

func (i *ImplServer) ProcessGameResult(ctx context.Context, gameInfo *room.GameInfo) {
	pRoom, err := i.getRoom(ctx)
	if err != nil {
		return
	}
	var result = gameInfo.Result
	var gameID = gameInfo.GameID
	var playerInfos = map[int]*player.Info{}
	for _, p := range pRoom.Players {
		playerInfos[p.Seat] = p.GetInfo()
	}
	go func() {
		var logContent = &LogContent{}
		select {
		case r := <-result:
			if r.Err != nil {
				global.Log.Warnf("game %s end error: %v", gameID, r.Err)
				logContent.Error = r.Err
			}
			global.Log.Debugf("game %s end", gameID)
			var playerIDs = [4]uint{0, 0, 0, 0}
			var players = [4]string{"", "", "", ""}
			for _, p := range playerInfos {
				pSeat := p.Seat
				pOrder := r.Seat2Order[pSeat]
				if p.ID != 0 {
					playerIDs[pOrder] = p.ID
				}
				players[pOrder] = p.Name
			}

			// log in database
			logContent.Players = players
			logContent.PlayerIDs = playerIDs
			logContent.Events = r.AllEvents

			logContentJson, err := json.Marshal(logContent)
			if err != nil {
				global.Log.Warnf("marshal log content error: %v", err)
				return
			}

			var entityUsers = make([]*entity.User, 0)
			for _, pID := range pRoom.ListPlayerIDs() {
				user, err := query.User.WithContext(dao.DBCtx).Where(query.User.ID.Eq(pID)).First()
				if err != nil {
					global.Log.Warnf("get user error: %v", err)
					continue
				}
				entityUsers = append(entityUsers, user)
			}
			var entityLog = &entity.Log{
				ID:      gameID,
				Content: string(logContentJson),
				Users:   entityUsers,
			}

			if err := query.Log.WithContext(dao.DBCtx).Create(entityLog); err != nil {
				global.Log.Warnf("create log error: %v", err)
				return
			}
		}
	}()
}
