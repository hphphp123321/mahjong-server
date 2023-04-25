package controller

import (
	"context"
	"fmt"
	pb "github.com/hphphp123321/mahjong-server/app/api/v1"
	"github.com/hphphp123321/mahjong-server/app/errs"
	"github.com/hphphp123321/mahjong-server/app/global"
	"io"
)

func AddReadyStream(ctx context.Context, stream pb.Mahjong_ReadyServer, server *MahjongServer) error {
	cid, err := server.s.GetID(ctx)
	if err != nil {
		return err
	}
	server.clients[cid].readyStream = stream
	return nil
}

func RemoveReadyStream(ctx context.Context, server *MahjongServer) error {
	cid, err := server.s.GetID(ctx)
	if err != nil {
		return err
	}
	server.clients[cid].readyStream = nil
	return nil
}

func BoardCastReadyReply(ctx context.Context, server *MahjongServer, reply *pb.ReadyReply) error {
	r, err := server.s.ListPlayerIDs(ctx)
	if err != nil {
		return err
	}
	for _, pid := range r.PlayerIDs {
		c, ok := server.clients[pid]
		if !ok {
			return errs.ErrClientNotFound
		}
		if c.readyStream == nil {
			continue
		}
		if err = c.readyStream.Send(reply); err != nil {
			global.Log.Warnf("send ready reply to %s failed: %v", pid, err)
			continue
		}
	}
	return err
}

func SendBack(ctx context.Context, server *MahjongServer, reply *pb.ReadyReply) error {
	cid, err := server.s.GetID(ctx)
	if err != nil {
		return err
	}
	c, ok := server.clients[cid]
	if !ok {
		return errs.ErrClientNotFound
	}
	if c.readyStream == nil {
		return errs.ErrStreamNotFound
	}
	return c.readyStream.Send(reply)
}

func StartReadyStream(ctx context.Context, stream pb.Mahjong_ReadyServer, server *MahjongServer) (done chan error, replyChan chan *pb.ReadyReply) {
	done = make(chan error)
	replyChan = make(chan *pb.ReadyReply)
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				global.Log.Infoln("ready stream closed: ", err)
				done <- nil
				return
			} else if err != nil {
				global.Log.Warnln(err)
				done <- err
				return
			}
			switch in.GetRequest().(type) {
			case *pb.ReadyRequest_RefreshRoom:
				reply, err := handleRefreshRoom(ctx, server, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_GetReady:
				reply, err := handleGetReady(ctx, server, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_CancelReady:
				reply, err := handleCancelReady(ctx, server, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_LeaveRoom:
				reply, err := handleLeaveRoom(ctx, server, in)
				sendChan(replyChan, done, reply, err)
				if reply == nil && err == nil {
					done <- nil
					return
				}
			case *pb.ReadyRequest_RemovePlayer:
				reply, err := handleRemovePlayer(ctx, server, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_AddRobot:
				reply, err := handleAddRobot(ctx, server, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_Chat:
				reply, err := handleChat(ctx, server, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_StartGame:
				reply, err := handleStartGame(ctx, server, in)
				sendChan(replyChan, done, reply, err)
			}
		}
	}()
	return done, replyChan
}

func handleRefreshRoom(ctx context.Context, server *MahjongServer, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	info, err := server.s.GetRoomInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("handle refresh room: %s", err)
	}
	reply = &pb.ReadyReply{
		Message: "",
		Reply:   &pb.ReadyReply_RefreshRoomReply{RefreshRoomReply: MapToRoomInfo(info)},
	}
	if err := SendBack(ctx, server, reply); err != nil {
		return nil, fmt.Errorf("handle refresh room: %s", err)
	}
	return nil, nil
}

func handleGetReady(ctx context.Context, server *MahjongServer, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	r, err := server.s.GetReady(ctx)
	if err != nil {
		return nil, fmt.Errorf("handle get ready: %s", err)
	}
	return ToPbGetReadyReply(r), nil
}

func handleCancelReady(ctx context.Context, server *MahjongServer, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	r, err := server.s.CancelReady(ctx)
	if err != nil {
		return nil, fmt.Errorf("handle cancel ready: %s", err)
	}
	return ToPbCancelReadyReply(r), nil
}

func handleLeaveRoom(ctx context.Context, server *MahjongServer, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	re, err := server.s.ListPlayerIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("handle leave room: %s", err)
	}
	var pids []string
	id, _ := server.s.GetID(ctx)
	for _, pid := range re.PlayerIDs {
		pids = append(pids, pid)
	}
	r, err := server.s.LeaveRoom(ctx)
	if err != nil {
		return nil, err
	}
	reply = ToPbLeaveRoomReply(r)
	for _, pid := range pids {
		c, ok := server.clients[pid]
		if !ok {
			return nil, errs.ErrClientNotFound
		}
		if c.readyStream == nil {
			return nil, errs.ErrStreamNotFound
		}
		if err = c.readyStream.Send(reply); err != nil {
			global.Log.Warnf("send leave room reply to %s failed: %v", pid, err)
			continue
		}
		if pid == id {
			c.readyStream = nil
		}
	}
	return nil, nil
}

func handleRemovePlayer(ctx context.Context, server *MahjongServer, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	re, err := server.s.ListPlayerIDs(ctx)
	reSeat := in.GetRemovePlayer().PlayerSeat
	reId, err := server.s.GetIDBySeat(ctx, int(reSeat))
	if err != nil {
		return nil, fmt.Errorf("handle remove player: %s", err)
	}
	var pids []string
	for _, pid := range re.PlayerIDs {
		pids = append(pids, pid)
	}
	r, err := server.s.RemovePlayer(ctx, ToServerRemovePlayerRequest(in))
	if err != nil {
		return nil, fmt.Errorf("handle remove player: %s", err)
	}
	for _, pid := range pids {
		c, ok := server.clients[pid]
		if !ok {
			return nil, errs.ErrClientNotFound
		}
		if c.readyStream == nil {
			return nil, errs.ErrStreamNotFound
		}
		if err = c.readyStream.Send(ToPbRemovePlayerReply(r)); err != nil {
			global.Log.Warnf("send remove player reply to %s failed: %v", pid, err)
			continue
		}
		if pid == reId {
			c.readyStream = nil
		}
	}
	return nil, nil
}

func handleAddRobot(ctx context.Context, server *MahjongServer, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	r, err := server.s.AddRobot(ctx, ToServerAddRobotRequest(in))
	if err != nil {
		return nil, fmt.Errorf("handle add robot: %s", err)
	}
	return ToPbAddRobotReply(r), nil
}

//func handleListRobots(ctx context.Context, server *MahjongServer, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
//	r, err := server.s.ListRobots(ctx)
//	if err != nil {
//		return nil, err
//	}
//	return ToPbListRobotsReply(r), nil
//}

func handleChat(ctx context.Context, server *MahjongServer, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	name, err := server.s.GetName(ctx)
	if err != nil {
		return nil, fmt.Errorf("handle chat: %s", err)
	}
	seat, err := server.s.GetSeat(ctx)
	if err != nil {
		return nil, fmt.Errorf("handle chat: %s", err)
	}
	return ToPbChatReply(in, name, seat), nil
}

func handleStartGame(ctx context.Context, server *MahjongServer, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	return nil, nil
	// TODO
}

func sendChan(replyChan chan<- *pb.ReadyReply, done chan<- error, reply *pb.ReadyReply, err error) {
	if err != nil {
		done <- err
	} else if reply != nil {
		replyChan <- reply
	}
}
